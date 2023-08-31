package parser

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
)

func (p *Parser) listen(ctx context.Context) {
	defer p.wg.Done()

	p.log.Info().Msg("module started")

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-p.input.Listen():
			if !ok {
				p.log.Warn().Msg("can't read message from input")
				continue
			}

			block, ok := msg.(types.ResultBlock)
			if !ok {
				p.log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			if err := p.parse(ctx, block); err != nil {
				p.log.Err(err).Msg("block parsing error")
				continue
			}
		}
	}
}
