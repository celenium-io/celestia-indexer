// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"sync"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/node"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Worker struct {
	api      node.Api
	blocks   chan types.BlockData
	log      zerolog.Logger
	queue    []types.Level
	capacity int
	liveMode bool
	m        map[types.Level]struct{}
	mx       *sync.RWMutex
}

func NewWorker(api node.Api, log zerolog.Logger, blocks chan types.BlockData, capacity int) *Worker {
	if capacity == 0 {
		capacity = 10
	}
	return &Worker{
		api:      api,
		blocks:   blocks,
		log:      log,
		queue:    make([]types.Level, 0),
		m:        make(map[types.Level]struct{}),
		capacity: capacity,
		mx:       new(sync.RWMutex),
	}
}

func (worker *Worker) SetLiveMode(liveMode bool) {
	worker.mx.Lock()
	{
		worker.liveMode = liveMode
	}
	worker.mx.Unlock()
}

func (worker *Worker) Capacity() int {
	return worker.capacity
}

func (worker *Worker) Do(ctx context.Context, level types.Level) {
	if _, ok := worker.m[level]; ok {
		return
	}
	worker.queue = append(worker.queue, level)
	worker.m[level] = struct{}{}

	if len(worker.queue) < worker.capacity {
		worker.mx.RLock()
		{
			if !worker.liveMode {
				worker.mx.RUnlock()
				return
			}
		}
		worker.mx.RUnlock()
	}

	start := time.Now()

	var result []types.BlockData
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		requestTimeout, cancel := context.WithTimeout(ctx, time.Minute)
		blocks, err := worker.api.BlockBulkData(requestTimeout, worker.queue...)
		if err != nil {
			cancel()

			if errors.Is(err, context.Canceled) {
				return
			}

			worker.log.Err(err).
				Uint64("height", uint64(level)).
				Msg("while getting block data")

			time.Sleep(time.Second)
			continue
		}
		result = blocks
		cancel()
		break
	}

	for i := range result {
		worker.log.Info().
			Uint64("height", uint64(result[i].Height)).
			Int64("ms", time.Since(start).Milliseconds()).
			Msg("received block")
		worker.blocks <- result[i]
	}

	worker.queue = worker.queue[:0]
	clear(worker.m)
}
