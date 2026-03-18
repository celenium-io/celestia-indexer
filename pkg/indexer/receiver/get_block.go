// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func isConnectionError(err error) bool {
	if os.IsTimeout(err) {
		return true
	}
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return true
	}
	s := err.Error()
	return strings.Contains(s, "EOF") ||
		strings.Contains(s, "connection reset") ||
		strings.Contains(s, "broken pipe")
}

func (r *Module) fetchBatch(ctx context.Context, levels []types.Level) {
	remaining := levels

	for len(remaining) > 0 {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var (
			batchSize = min(int(r.bulkSize.Load()), len(remaining))
			chunk     = remaining[:batchSize]
			start     = time.Now()
		)

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

			r.Log.Err(err).
				Msg("while getting block data")

			if isConnectionError(err) {
				r.connectionErrorDecrease()
			}

			time.Sleep(time.Second)
			continue
		}

		r.adjustBulkSize(len(chunk), time.Since(start))
		remaining = remaining[batchSize:]
	}
}
