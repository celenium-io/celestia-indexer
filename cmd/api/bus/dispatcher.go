// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package bus

import (
	"context"
	"strconv"
	"sync"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/workerpool"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Dispatcher struct {
	listener storage.Listener
	blocks   storage.IBlock
	txs      storage.ITx

	mx        *sync.RWMutex
	observers []*Observer

	g workerpool.Group
}

func NewDispatcher(
	factory storage.ListenerFactory,
	blocks storage.IBlock,
	txs storage.ITx,
) (*Dispatcher, error) {
	if factory == nil {
		return nil, errors.New("nil listener factory")
	}
	listener := factory.CreateListener()
	return &Dispatcher{
		listener:  listener,
		blocks:    blocks,
		txs:       txs,
		observers: make([]*Observer, 0),
		mx:        new(sync.RWMutex),
		g:         workerpool.NewGroup(),
	}, nil
}

func (d *Dispatcher) Observe(channels ...string) *Observer {
	if observer := NewObserver(channels...); observer != nil {
		d.mx.Lock()
		d.observers = append(d.observers, observer)
		d.mx.Unlock()
		return observer
	}

	return nil
}

func (d *Dispatcher) Start(ctx context.Context) {
	if err := d.listener.Subscribe(ctx, storage.ChannelHead, storage.ChannelTx); err != nil {
		log.Err(err).Msg("subscribe on postgres notifications")
		return
	}
	d.g.GoCtx(ctx, d.listen)
}

func (d *Dispatcher) Close() error {
	d.g.Wait()

	d.mx.Lock()
	for i := range d.observers {
		if err := d.observers[i].Close(); err != nil {
			return err
		}
	}
	d.mx.Unlock()

	return nil
}

func (d *Dispatcher) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case notification, ok := <-d.listener.Listen():
			if !ok {
				return
			}
			if notification == nil {
				log.Warn().Str("channel", notification.Channel).Msg("nil notification")
				continue
			}
			if err := d.handleNotification(ctx, notification); err != nil {
				log.Err(err).Str("channel", notification.Channel).Msg("handle notification")
			}
		}
	}
}

func (d *Dispatcher) handleNotification(ctx context.Context, notification *pq.Notification) error {
	id, err := strconv.ParseUint(notification.Extra, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "parse block id: %s", notification.Extra)
	}

	switch notification.Channel {
	case storage.ChannelHead:
		err = d.handleBlock(ctx, id)
	case storage.ChannelTx:
		err = d.handleTx(ctx, id)
	default:
		err = errors.Errorf("unknown channel name: %s", notification.Channel)
	}

	return err
}

func (d *Dispatcher) handleBlock(ctx context.Context, id uint64) error {
	block, err := d.blocks.ByIdWithRelations(ctx, id)
	if err != nil {
		return err
	}
	d.mx.RLock()
	for i := range d.observers {
		d.observers[i].notifyBlocks(&block)
	}
	d.mx.RUnlock()
	return nil
}

func (d *Dispatcher) handleTx(ctx context.Context, id uint64) error {
	tx, err := d.txs.ByIdWithRelations(ctx, id)
	if err != nil {
		return err
	}
	d.mx.RLock()
	for i := range d.observers {
		d.observers[i].notifyTxs(&tx)
	}
	d.mx.RUnlock()
	return nil
}
