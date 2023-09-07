package parser

import (
	"context"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Module struct {
	input  *modules.Input
	output *modules.Output
	log    zerolog.Logger
	g      workerpool.Group
}

const (
	name        = "parser"
	BlocksInput = "blocks"
	DataOutput  = "data"
)

func NewModule() Module {
	return Module{
		input:  modules.NewInput(BlocksInput),
		output: modules.NewOutput(DataOutput),
		log:    log.With().Str("module", name).Logger(),
		g:      workerpool.NewGroup(),
	}
}

// Name -
func (*Module) Name() string {
	return name
}

func (p *Module) Start(ctx context.Context) {
	p.log.Info().Msg("starting parser module...")
	p.g.GoCtx(ctx, p.listen)
}

func (p *Module) Close() error {
	p.log.Info().Msg("closing...")
	p.g.Wait()

	return p.input.Close()
}

func (p *Module) Output(name string) (*modules.Output, error) {
	if name != DataOutput {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}

	return p.output, nil
}

func (p *Module) Input(name string) (*modules.Input, error) {
	if name != BlocksInput {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return p.input, nil
}

func (p *Module) AttachTo(name string, input *modules.Input) error {
	output, err := p.Output(name)
	if err != nil {
		return err
	}

	output.Attach(input)
	return nil
}
