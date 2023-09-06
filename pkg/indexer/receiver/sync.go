package receiver

import (
	"context"
	"time"
)

func (r *Receiver) sync(ctx context.Context) {
	if err := r.readBlocks(ctx); err != nil {
		r.log.Err(err).Msg("while reading blocks")
		return
	}

	ticker := time.NewTicker(time.Second * time.Duration(r.cfg.BlockPeriod))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.readBlocks(ctx); err != nil {
				r.log.Err(err).Msg("while reading blocks by timer")
				return
			}
		}
	}
}
