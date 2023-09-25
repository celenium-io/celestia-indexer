package receiver

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"time"
)

func (r *Module) sync(ctx context.Context) {
	var blocksCtx context.Context
	blocksCtx, r.cancelReadBlocks = context.WithCancel(ctx)
	if err := r.readBlocks(blocksCtx); err != nil && !errors.Is(err, context.Canceled) {
		r.Log.Err(err).Msg("while reading blocks")
		r.stopAll()
		return
	}

	ticker := time.NewTicker(time.Second * time.Duration(r.cfg.BlockPeriod))
	defer ticker.Stop()

	for {
		r.rollbackSync.Wait()

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			blocksCtx, r.cancelReadBlocks = context.WithCancel(ctx)
			if err := r.readBlocks(blocksCtx); err != nil && !errors.Is(err, context.Canceled) {
				r.Log.Err(err).Msg("while reading blocks by timer")
				r.stopAll()
				return
			}
		}
	}
}

func (r *Module) readBlocks(ctx context.Context) error {
	headLevel, err := r.headLevel(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	level, _ := r.Level()
	level += 1

	for ; level <= headLevel; level++ {
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
