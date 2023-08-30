package receiver

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
	"time"
)

func (r *Receiver) worker(ctx context.Context, level storage.Level) {
	start := time.Now()

	var result types.ResultBlock
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		block, err := r.api.GetBlock(ctx, level)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			r.log.Err(err).Msg("get block request")
			time.Sleep(time.Second)
			continue
		}

		result = block
		break
	}

	r.log.Info().Int64("height", result.Block.Height).Int64("ms", time.Since(start).Milliseconds()).Msg("received block data")
	r.blocks <- result
}
