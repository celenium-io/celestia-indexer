// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (r *Module) worker(ctx context.Context, level types.Level) {
	defer r.taskQueue.Delete(level)

	start := time.Now()

	var result types.BlockData
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		requestTimeout, cancel := context.WithTimeout(ctx, time.Minute)
		block, err := r.api.BlockData(requestTimeout, level)
		if err != nil {
			cancel()

			if errors.Is(err, context.Canceled) {
				return
			}

			r.Log.Err(err).
				Uint64("height", uint64(level)).
				Msg("while getting block data")

			time.Sleep(time.Second)
			continue
		}

		block.AppVersion = r.appVersion.Load()
		result = block
		cancel()
		break
	}

	r.Log.Info().
		Uint64("height", uint64(result.Height)).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("received block")
	r.blocks <- result
}
