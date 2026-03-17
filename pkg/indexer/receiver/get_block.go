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
	batchSize := len(levels)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		chunks := slices.Chunk(levels, batchSize)

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
					batchSize = 1
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
