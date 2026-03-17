// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"os"
	"slices"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (r *Module) fetchBatch(ctx context.Context, levels []types.Level) {
	r.Log.Debug().
		Int("batch_size", len(levels)).
		Int64("bulk_size", r.bulkSize.Load()).
		Msg("fetch batch started")

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		bulkSize := int(r.bulkSize.Load())
		chunks := slices.Chunk(levels, len(levels))
		if len(levels) > bulkSize {
			chunkSize := len(levels) / max(1, bulkSize)
			chunks = slices.Chunk(levels, chunkSize)
		}

		start := time.Now()
		for chunk := range chunks {
			_, err := r.circuitBreaker.Execute(func() (any, error) {
				err := r.api.BlockBulkDataStream(ctx, func(block types.BlockData) error {
					r.Log.Info().
						Uint64("height", uint64(block.Height)).
						Int64("ms", time.Since(start).Milliseconds()).
						Msg("received block")
					select {
					case r.blocks <- &block:
					case <-ctx.Done():
						return ctx.Err()
					}
					return nil
				}, chunk...)
				return nil, err
			})
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}

				if os.IsTimeout(err) {
					r.adjustBulkSize(len(levels), time.Since(start))
					r.bulkSize.Store(1)
					r.Log.Info().
						Int64("bulk_size", 1).
						Msg("reset bulk size due to timeout error")
				}

				r.Log.Err(err).
					Msg("while getting block data")

				time.Sleep(time.Second)
				continue
			}

			r.adjustBulkSize(len(levels), time.Since(start))
		}
		return
	}
}
