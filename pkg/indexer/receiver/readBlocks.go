package receiver

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

func (r *Receiver) readBlocks(ctx context.Context) error {
	head, err := r.api.Head(ctx)
	if err != nil {
		return err
	}

	startLevel := storage.Level(r.cfg.Indexer.StartLevel) // TODO read from current state
	headLevel := storage.Level(head.Block.Height)

	for level := startLevel; level <= headLevel; level++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			r.pool.AddTask(level)
		}
	}

	return nil
}
