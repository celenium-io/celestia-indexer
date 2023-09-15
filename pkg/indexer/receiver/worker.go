package receiver

import (
	"context"
	"time"

	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (r *Module) worker(ctx context.Context, level types.Level) {
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
			r.Log.Err(err).
				Uint64("height", uint64(level)).
				Msg("while getting block data")

			time.Sleep(time.Second)
			continue
		}

		result = block
		break
	}

	r.Log.Info().
		Uint64("height", uint64(result.Height)).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("received block")
	r.blocks <- result
}

func (r *Module) blockData(ctx context.Context, level types.Level) (types.BlockData, error) {
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
