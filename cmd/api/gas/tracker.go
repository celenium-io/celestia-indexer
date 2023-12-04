// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package gas

import (
	"context"
	"sort"
	"sync"

	"github.com/celenium-io/celestia-indexer/cmd/api/bus"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/workerpool"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

const (
	blockCount = 100
)

var (
	percentiles = []float64{.1, .5, .99}
)

type Tracker struct {
	state    storage.IState
	stats    storage.IBlockStats
	tx       storage.ITx
	observer *bus.Observer
	log      zerolog.Logger
	mx       *sync.RWMutex
	gasState GasPrice
	q        *queue
	g        workerpool.Group
}

func NewTracker(
	state storage.IState,
	stats storage.IBlockStats,
	tx storage.ITx,
	observer *bus.Observer,
) *Tracker {
	return &Tracker{
		state:    state,
		stats:    stats,
		tx:       tx,
		observer: observer,
		mx:       new(sync.RWMutex),
		gasState: GasPrice{
			Slow:   "0",
			Median: "0",
			Fast:   "0",
		},
		log: log.With().Str("module", "gas_tracker").Logger(),
		q:   newQueue(blockCount),
		g:   workerpool.NewGroup(),
	}
}

func (tracker *Tracker) Start(ctx context.Context) {
	tracker.g.GoCtx(ctx, tracker.listen)
}

func (tracker *Tracker) Close() error {
	tracker.g.Wait()

	return nil
}

func (tracker *Tracker) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case block, ok := <-tracker.observer.Blocks():
			if !ok {
				return
			}
			if err := tracker.processBlock(ctx, block.Stats); err != nil {
				log.Err(err).Msg("new block processing")
				continue
			}
			if err := tracker.computeMetrics(); err != nil {
				log.Err(err).Msg("compute metrics")
			}
		}
	}
}

func (tracker *Tracker) Init(ctx context.Context) error {
	state, err := tracker.state.List(ctx, 1, 0, sdk.SortOrderAsc)
	if err != nil {
		return err
	}
	if len(state) == 0 {
		return nil
	}
	blockStats, err := tracker.stats.LastFrom(ctx, state[0].LastHeight, blockCount)
	if err != nil {
		return err
	}

	for i := len(blockStats) - 1; i >= 0; i-- {
		if err := tracker.processBlock(ctx, blockStats[i]); err != nil {
			return err
		}
	}

	return tracker.computeMetrics()
}

func (tracker *Tracker) processBlock(ctx context.Context, blockStat storage.BlockStats) error {
	data := info{
		Height:       uint64(blockStat.Height),
		TxCount:      blockStat.TxCount,
		GasUsed:      blockStat.GasUsed,
		GasWanted:    blockStat.GasLimit,
		Fee:          blockStat.Fee,
		GasUsedRatio: decimal.New(0, 1),
		Percentiles:  make([]decimal.Decimal, 0),
	}

	for range percentiles {
		data.Percentiles = append(data.Percentiles, decimal.New(0, 1))
	}

	if data.GasWanted > 0 {
		data.GasUsedRatio = decimal.NewFromInt(data.GasUsed).Div(decimal.NewFromInt(data.GasWanted))
	}

	if blockStat.TxCount == 0 {
		tracker.q.Push(data)
		return nil
	}

	txs, err := tracker.tx.Gas(ctx, blockStat.Height)
	if err != nil {
		return err
	}
	sort.Sort(storage.ByGasPrice(txs))

	var (
		sumGas  = txs[0].GasWanted
		txIndex = 0
	)

	for i, p := range percentiles {
		threshold := uint64(float64(blockStat.GasLimit) * p)
		for sumGas < int64(threshold) && txIndex < len(txs)-1 {
			txIndex++
			sumGas += txs[txIndex].GasWanted
		}
		data.Percentiles[i] = txs[txIndex].GasPrice.Copy()
	}

	tracker.q.Push(data)
	return nil
}

func (tracker *Tracker) State() GasPrice {
	tracker.mx.RLock()
	defer tracker.mx.RUnlock()

	return tracker.gasState
}

func (tracker *Tracker) computeMetrics() error {
	slow := decimal.New(0, 1)
	median := decimal.New(0, 1)
	fast := decimal.New(0, 1)

	err := tracker.q.Range(func(item info) (bool, error) {
		slow = slow.Add(item.Percentiles[0])
		median = median.Add(item.Percentiles[1])
		fast = fast.Add(item.Percentiles[2])
		return false, nil
	})
	if err != nil {
		return err
	}
	count := int64(tracker.q.Size())

	slow = slow.Div(decimal.NewFromInt(count))
	median = median.Div(decimal.NewFromInt(count))
	fast = fast.Div(decimal.NewFromInt(count))

	tracker.mx.Lock()
	{
		tracker.gasState.Slow = currency.StringTia(slow)
		tracker.gasState.Median = currency.StringTia(median)
		tracker.gasState.Fast = currency.StringTia(fast)
	}
	tracker.mx.Unlock()

	return nil
}
