// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package bus

import (
	"context"
	"strconv"
	"sync"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/workerpool"
	jsoniter "github.com/json-iterator/go"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Dispatcher struct {
	listener storage.Listener
	blocks   storage.IBlock

	mx        *sync.RWMutex
	observers []*Observer

	g workerpool.Group
}

func NewDispatcher(
	factory storage.ListenerFactory,
	blocks storage.IBlock,
) (*Dispatcher, error) {
	if factory == nil {
		return nil, errors.New("nil listener factory")
	}
	listener := factory.CreateListener()
	return &Dispatcher{
		listener:  listener,
		blocks:    blocks,
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
	if err := d.listener.Subscribe(ctx, storage.ChannelHead, storage.ChannelBlock); err != nil {
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
	switch notification.Channel {
	case storage.ChannelHead:
		return d.handleState(ctx, notification.Extra)
	case storage.ChannelBlock:
		return d.handleBlock(ctx, notification.Extra)
	default:
		return errors.Errorf("unknown channel name: %s", notification.Channel)
	}
}

func (d *Dispatcher) handleBlock(ctx context.Context, payload string) error {
	id, err := strconv.ParseUint(payload, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "parse block id: %s", payload)
	}

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

func (d *Dispatcher) handleState(ctx context.Context, payload string) error {
	var state storage.State
	if err := jsoniter.UnmarshalFromString(payload, &state); err != nil {
		return err
	}

	d.mx.RLock()
	for i := range d.observers {
		d.observers[i].notifyState(&state)
	}
	d.mx.RUnlock()
	return nil
}
