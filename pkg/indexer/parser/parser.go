package parser

import (
	"context"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type Parser struct {
	input  *modules.Input
	output *modules.Output
	log    zerolog.Logger
	wg     *sync.WaitGroup
}

const (
	name        = "parser"
	BlocksInput = "blocks"
	DataOutput  = "data"
)

func NewModule() Parser {
	return Parser{
		input:  modules.NewInput(BlocksInput),
		output: modules.NewOutput(DataOutput),
		log:    log.With().Str("module", name).Logger(),
		wg:     new(sync.WaitGroup),
	}
}

// Name -
func (*Parser) Name() string {
	return name
}

func (p *Parser) Start(ctx context.Context) {
	p.log.Info().Msg("starting parser module...")

	p.wg.Add(1)
	go p.listen(ctx)
}

func (p *Parser) Close() error {
	p.log.Info().Msg("closing...")
	p.wg.Wait()

	return p.input.Close()
}

func (p *Parser) Output(name string) (*modules.Output, error) {
	if name != DataOutput {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}

	return p.output, nil
}

func (p *Parser) Input(name string) (*modules.Input, error) {
	if name != BlocksInput {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return p.input, nil
}

func (p *Parser) AttachTo(name string, input *modules.Input) error {
	output, err := p.Output(name)
	if err != nil {
		return err
	}

	output.Attach(input)
	return nil
}
