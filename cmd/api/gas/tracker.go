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
	coreTypes "github.com/tendermint/tendermint/types"
)

const (
	blockCount        = 100
	emptyBlockPercent = .90
)

var (
	percentiles  = []float64{.10, .50, .99}
	maxBlockSize = coreTypes.MaxDataBytesNoEvidence(1974272, 100)
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
		Height:         uint64(blockStat.Height),
		TxCount:        blockStat.TxCount,
		GasUsed:        blockStat.GasUsed,
		GasWanted:      blockStat.GasLimit,
		Fee:            blockStat.Fee,
		GasUsedRatio:   decimal.New(0, 1),
		Percentiles:    make([]decimal.Decimal, 0),
		BlockOccupancy: float64(blockStat.BytesInBlock) / float64(maxBlockSize),
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

	tracker.compute(txs, blockStat.GasLimit, &data)

	if data.BlockOccupancy < emptyBlockPercent {
		// If block occupancy is less than empty block threshold set all percentiles to slow.
		for i := 1; i < len(data.Percentiles); i++ {
			data.Percentiles[i] = data.Percentiles[0].Copy()
		}
	}

	tracker.q.Push(data)
	return nil
}

func (tracker *Tracker) compute(txs []storage.Gas, gasLimit int64, data *info) {
	if len(txs) == 0 {
		return
	}

	var (
		txIndex = 0
		sumGas  = txs[txIndex].GasWanted
	)

	for i, p := range percentiles {
		threshold := uint64(float64(gasLimit) * p)
		for sumGas < int64(threshold) && txIndex < len(txs)-1 {
			txIndex++
			sumGas += txs[txIndex].GasWanted
		}
		data.Percentiles[i] = txs[txIndex].GasPrice.Copy()
	}
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
