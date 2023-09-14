package genesis

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
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
	modules.BaseModule
	storage     postgres.Storage
	indexerName string
}

var _ modules.Module = (*Module)(nil)

// NewModule -
func NewModule(pg postgres.Storage, cfg config.Indexer) Module {
	m := Module{
		BaseModule:  modules.New("genesis"),
		storage:     pg,
		indexerName: cfg.Name,
	}

	m.CreateInput(InputName)
	m.CreateOutput(OutputName)
	m.CreateOutput(StopOutput)

	return m
}

// Start -
func (module *Module) Start(ctx context.Context) {
	module.G.GoCtx(ctx, module.listen)
}

func (module *Module) listen(ctx context.Context) {
	module.Log.Info().Msg("module started")
	input := module.MustInput(InputName)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-input.Listen():
			if !ok {
				module.Log.Warn().Msg("can't read message from input")
				return
			}
			genesis, ok := msg.(types.Genesis)
			if !ok {
				module.Log.Warn().Msgf("invalid message type: %T", msg)
				return
			}

			module.Log.Info().Msg("received genesis message")

			block, err := module.parse(genesis)
			if err != nil {
				module.Log.Err(err).Msgf("parsing genesis block")
				return
			}
			module.Log.Info().Msg("parsed genesis message")

			if err := module.save(ctx, block); err != nil {
				module.Log.Err(err).Msg("saving genesis block error")
				return
			}
			module.Log.Info().Msg("saved genesis message")

			module.MustOutput(OutputName).Push(struct{}{})
			return
		}
	}
}

// Close -
func (module *Module) Close() error {
	module.Log.Info().Msg("closing module...")
	module.G.Wait()
	return nil
}
