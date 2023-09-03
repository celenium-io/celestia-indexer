package rollback

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	InputName          = "signal"
	OutputName         = "state"
	defaultIndexerName = "celestia-indexer"
)

// Module - executes rollback on signal from input (target height) and notify all subscribers about new state after rollback operation.
//
//	                     |----------------|
//	                     |                |
//	-- storage.Level ->  |     MODULE     |  -- storage.State ->
//	                     |                |
//	                     |----------------|
type Module struct {
	tx        sdk.Transactable
	state     storage.IState
	indexName string
	input     *modules.Input
	output    *modules.Output
	log       zerolog.Logger
	g         workerpool.Group
}

func NewModule(tx sdk.Transactable, state storage.IState, opts ...ModuleOption) *Module {
	module := Module{
		tx:        tx,
		state:     state,
		input:     modules.NewInput(InputName),
		output:    modules.NewOutput(OutputName),
		indexName: defaultIndexerName,
		g:         workerpool.NewGroup(),
	}
	module.log = log.With().Str("module", module.Name()).Logger()

	for i := range opts {
		opts[i](&module)
	}

	return &module
}

func (*Module) Name() string {
	return "rollback"
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
				continue
			}

			height, ok := msg.(storage.Level)
			if !ok {
				module.log.Error().Msgf("invalid input type: %T", height)
				continue
			}

			if err := module.rollback(ctx, height); err != nil {
				module.log.Err(err).Msgf("error occured")
			}
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
	if name != OutputName {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return module.output, nil
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
func (module *Module) rollback(ctx context.Context, height storage.Level) error {
	state, err := module.state.ByName(ctx, module.indexName)
	if err != nil {
		return err
	}

	for level := state.LastHeight; level > height; level-- {
		log.Info().
			Uint64("height", uint64(level)).
			Uint64("new_height", uint64(level-1)).
			Msg("start roll back from height...")

		if err := module.rollbackBlock(ctx, level); err != nil {
			return err
		}

		log.Info().
			Uint64("height", uint64(level)).
			Uint64("new_height", uint64(level-1)).
			Msg("rolled back from height")
	}

	newState, err := module.state.ByName(ctx, module.indexName)
	if err != nil {
		return err
	}
	module.output.Push(newState)

	return nil
}

func (module *Module) rollbackBlock(ctx context.Context, height storage.Level) error {
	tx, err := postgres.BeginTransaction(ctx, module.tx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	block, err := tx.RollbackBlock(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	addresses, err := tx.RollbackAddresses(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	_, err = tx.RollbackTxs(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	msgs, err := tx.RollbackMessages(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	events, err := tx.RollbackEvents(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := module.balances(ctx, events, addresses); err != nil {
		return tx.HandleError(ctx, err)
	}

	nsMsgs, err := tx.RollbackNamespaceMessages(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	ns, err := tx.RollbackNamespaces(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := module.namespaces(ctx, tx, nsMsgs, ns, msgs); err != nil {
		return tx.HandleError(ctx, errors.Wrap(err, "namespace rollback"))
	}

	newBlock, err := tx.LastBlock(ctx)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	state, err := tx.State(ctx, module.indexName)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	state.LastHeight = newBlock.Height
	state.LastTime = newBlock.Time
	state.TotalTx -= block.TxCount
	state.TotalBlobsSize -= block.BlobsSize
	state.TotalNamespaces -= uint64(len(ns))
	state.TotalAccounts -= uint64(len(addresses))
	state.TotalFee = state.TotalFee.Sub(block.Fee)

	totalSupplyDiff, err := module.totalSupplyDiff(ctx, events)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	state.TotalSupply = state.TotalSupply.Add(totalSupplyDiff)

	if err := tx.Update(ctx, &state); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}

	return nil
}
