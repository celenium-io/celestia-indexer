package parser

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func (p *Module) listen(ctx context.Context) {
	p.Log.Info().Msg("module started")

	input := p.MustInput(InputName)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-input.Listen():
			if !ok {
				p.Log.Warn().Msg("can't read message from input, it was drained and closed")
				p.MustOutput(StopOutput).Push(struct{}{})
				return
			}

			block, ok := msg.(types.BlockData)
			if !ok {
				p.Log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			if err := p.parse(ctx, block); err != nil {
				p.Log.Err(err).
					Uint64("height", uint64(block.Height)).
					Msg("block parsing error")
				p.MustOutput(StopOutput).Push(struct{}{})
				continue
			}
		}
	}
}
