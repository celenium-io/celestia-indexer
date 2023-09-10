package genesis

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// constants
const (
	InputName  = "block"
	OutputName = "finished"
	StopOutput = "stop"
)

// Module - saves received from input genesis block to storage and notify if it was success.
//
//	                     |----------------|
//	                     |                |
//	-- storage.Block ->  |     MODULE     | -- struct{} ->
//	                     |                |
//	                     |----------------|
type Module struct {
	storage     postgres.Storage
	input       *modules.Input
	outputs     map[string]*modules.Output
	indexerName string
	log         zerolog.Logger
	g           workerpool.Group
}

// NewModule -
func NewModule(pg postgres.Storage, cfg config.Indexer) Module {
	m := Module{
		storage: pg,
		input:   modules.NewInput(InputName),
		outputs: map[string]*modules.Output{
			OutputName: modules.NewOutput(OutputName),
			StopOutput: modules.NewOutput(StopOutput),
		},
		indexerName: cfg.Name,
		g:           workerpool.NewGroup(),
	}
	m.log = log.With().Str("module", m.Name()).Logger()

	return m
}

// Name -
func (*Module) Name() string {
	return "genesis"
}

// Start -
func (module *Module) Start(ctx context.Context) {
	module.g.GoCtx(ctx, module.listen)
}

func (module *Module) listen(ctx context.Context) {
	module.log.Info().Msg("module started")

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-module.input.Listen():
			if !ok {
				module.log.Warn().Msg("can't read message from input")
				return
			}
			genesis, ok := msg.(types.Genesis)
			if !ok {
				module.log.Warn().Msgf("invalid message type: %T", msg)
				return
			}

			module.log.Info().Msg("received genesis message")

			block, err := module.parse(genesis)
			if err != nil {
				module.log.Err(err).Msgf("parsing genesis block")
				return
			}
			module.log.Info().Msg("parsed genesis message")

			if err := module.save(ctx, block); err != nil {
				module.log.Err(err).Msg("saving genesis block error")
				return
			}
			module.log.Info().Msg("saved genesis message")

			module.outputs[OutputName].Push(struct{}{})
			return
		}
	}
}

// Close -
func (module *Module) Close() error {
	module.log.Info().Msg("closing module...")
	module.g.Wait()

	return module.input.Close()
}

// Output -
func (module *Module) Output(name string) (*modules.Output, error) {
	output, ok := module.outputs[name]
	if !ok {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return output, nil
}

// Input -
func (module *Module) Input(name string) (*modules.Input, error) {
	if name != InputName {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return module.input, nil
}

// AttachTo -
func (module *Module) AttachTo(name string, input *modules.Input) error {
	output, err := module.Output(name)
	if err != nil {
		return err
	}

	output.Attach(input)
	return nil
}
