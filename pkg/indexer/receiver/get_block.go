// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"time"

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

		blocks, err := r.api.BlockBulkData(ctx, r.taskQueue.Keys()...)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			r.Log.Err(err).
				Msg("while getting block data")

			time.Sleep(time.Second)
			continue
		}

		for i := range blocks {
			r.Log.Info().
				Uint64("height", uint64(blocks[i].Height)).
				Int64("ms", time.Since(start).Milliseconds()).
				Msg("received block")
			r.blocks <- blocks[i]
			if blocks[i].Height > r.receivedLevel {
				r.receivedLevel = blocks[i].Height
			}
		}

		r.taskQueue.Clear()
		return
	}

}
