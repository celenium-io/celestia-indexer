package receiver

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

func (r *Module) sync(ctx context.Context) {
	var blocksCtx context.Context
	blocksCtx, r.cancelReadBlocks = context.WithCancel(ctx)
	if err := r.readBlocks(blocksCtx); err != nil && !errors.Is(err, context.Canceled) {
		r.log.Err(err).Msg("while reading blocks")
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
				r.log.Err(err).Msg("while reading blocks by timer")
				return
			}
		}
	}
}
