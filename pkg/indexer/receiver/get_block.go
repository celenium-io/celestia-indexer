// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"os"
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

		start := time.Now()
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
			}, levels...)
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
		return
	}
}
