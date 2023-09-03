package receiver

import (
	"context"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (r *Receiver) worker(ctx context.Context, level storage.Level) {
	start := time.Now()

	var result types.BlockData
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		block, err := r.blockData(ctx, level)
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

	r.log.Info().
		Int64("height", int64(result.Height)).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("received block data")
	r.blocks <- result
}

func (r *Receiver) blockData(ctx context.Context, level storage.Level) (types.BlockData, error) {
	block, err := r.api.Block(ctx, level)
	if err != nil {
		return types.BlockData{}, err
	}

	blockResults, err := r.api.BlockResults(ctx, level)
	if err != nil {
		return types.BlockData{}, err
	}

	return types.BlockData{ResultBlock: block, ResultBlockResults: blockResults}, nil
}
