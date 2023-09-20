package receiver

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (r *Module) readBlocks(ctx context.Context) error {
	headLevel, err := r.headLevel(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	for level, _ := r.Level(); level <= headLevel; level++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			if _, ok := r.taskQueue.Get(level); !ok {
				r.taskQueue.Set(level, struct{}{})
				r.pool.AddTask(level)
			}
		}
	}

	return nil
}

func (r *Module) headLevel(ctx context.Context) (types.Level, error) {
	status, err := r.api.Status(ctx)
	if err != nil {
		return 0, err
	}

	return types.Level(status.SyncInfo.LatestBlockHeight), nil
}
