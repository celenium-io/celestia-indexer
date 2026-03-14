// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (r *Module) getBlocks(ctx context.Context) {
	if r.taskQueue.Len() == 0 {
		return
	}

	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, err := r.circuitBreaker.Execute(func() (any, error) {
			err := r.api.BlockBulkDataStream(ctx, func(block types.BlockData) error {
				r.Log.Info().
					Uint64("height", uint64(block.Height)).
					Int64("ms", time.Since(start).Milliseconds()).
					Msg("received block")
				r.blocks <- block
				if block.Height > r.receivedLevel {
					r.receivedLevel = block.Height
				}
				return nil
			}, r.taskQueue.Keys()...)
			return nil, err
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			r.Log.Err(err).
				Msg("while getting block data")

			time.Sleep(time.Second)
			continue
		}

		r.taskQueue.Clear()
		return
	}

}
