package receiver

import (
	"context"
	"github.com/pkg/errors"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

func (r *Receiver) readBlocks(ctx context.Context) error {
	headLevel, err := r.headLevel(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	for level := r.level; level <= headLevel; level++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			r.pool.AddTask(level)
		}
	}

	return nil
}

func (r *Receiver) headLevel(ctx context.Context) (storage.Level, error) {
	status, err := r.api.Status(ctx)
	if err != nil {
		return 0, err
	}

	return storage.Level(status.SyncInfo.LatestBlockHeight), nil
}
