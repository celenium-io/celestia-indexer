package rollback

import (
	"bytes"
	"context"

	"github.com/dipdup-io/celestia-indexer/pkg/node"

	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/types"

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
	InputName  = "signal"
	OutputName = "state"
	StopOutput = "stop"
)

// Module - executes rollback on signal from input and notify all subscribers about new state after rollback operation.
//
//	                |----------------|
//	                |                |
//	-- struct{} ->  |     MODULE     |  -- storage.State ->
//	                |                |
//	                |----------------|
type Module struct {
	tx        sdk.Transactable
	state     storage.IState
	blocks    storage.IBlock
	node      node.API
	indexName string
	input     *modules.Input
	outputs   map[string]*modules.Output
	log       zerolog.Logger
	g         workerpool.Group
}

func NewModule(
	tx sdk.Transactable,
	state storage.IState,
	blocks storage.IBlock,
	node node.API,
	cfg config.Indexer,
) Module {
	module := Module{
		tx:     tx,
		state:  state,
		blocks: blocks,
		node:   node,
		input:  modules.NewInput(InputName),
		outputs: map[string]*modules.Output{
			OutputName: modules.NewOutput(OutputName),
			StopOutput: modules.NewOutput(StopOutput),
		},
		indexName: cfg.Name,
		g:         workerpool.NewGroup(),
	}
	module.log = log.With().Str("module", module.Name()).Logger()

	return module
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
		case _, ok := <-module.input.Listen():
			if !ok {
				module.log.Warn().Msg("can't read message from input")
				return
			}

			if err := module.rollback(ctx); err != nil {
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

func (module *Module) rollback(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
		default:
			lastBlock, err := module.blocks.Last(ctx)
			if err != nil {
				return errors.Wrap(err, "receive last block from database")
			}

			nodeBlock, err := module.node.Block(ctx, lastBlock.Height)
			if err != nil {
				return errors.Wrapf(err, "receive block from node by height: %d", lastBlock.Height)
			}

			log.Debug().
				Uint64("height", uint64(lastBlock.Height)).
				Hex("db_block_hash", lastBlock.Hash).
				Hex("node_block_hash", nodeBlock.BlockID.Hash).
				Msg("comparing hash...")

			if bytes.Equal(lastBlock.Hash, nodeBlock.BlockID.Hash) {
				return module.finish(ctx)
			}

			log.Warn().
				Uint64("height", uint64(lastBlock.Height)).
				Hex("db_block_hash", lastBlock.Hash).
				Hex("node_block_hash", nodeBlock.BlockID.Hash).
				Msg("need rollback")

			if err := module.rollbackBlock(ctx, lastBlock.Height); err != nil {
				return errors.Wrapf(err, "rollback block: %d", lastBlock.Height)
			}
		}
	}
}

func (module *Module) finish(ctx context.Context) error {
	newState, err := module.state.ByName(ctx, module.indexName)
	if err != nil {
		return err
	}
	module.outputs[OutputName].Push(newState)

	log.Info().
		Uint64("new_height", uint64(newState.LastHeight)).
		Msg("roll backed to new height")

	return nil
}

func (module *Module) rollbackBlock(ctx context.Context, height types.Level) error {
	tx, err := postgres.BeginTransaction(ctx, module.tx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	if err := tx.RollbackBlock(ctx, height); err != nil {
		return tx.HandleError(ctx, err)
	}
	blockStats, err := tx.RollbackBlockStats(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	addresses, err := tx.RollbackAddresses(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := module.rollbackTransactions(ctx, tx, height); err != nil {
		return tx.HandleError(ctx, err)
	}

	totalNamespaces, err := module.rollbackMessages(ctx, tx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	events, err := tx.RollbackEvents(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := module.rollbackBalances(ctx, events, addresses); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.RollbackValidators(ctx, height); err != nil {
		return tx.HandleError(ctx, err)
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
	state.TotalTx -= blockStats.TxCount
	state.TotalBlobsSize -= blockStats.BlobsSize
	state.TotalNamespaces -= totalNamespaces
	state.TotalAccounts -= uint64(len(addresses))
	state.TotalFee = state.TotalFee.Sub(blockStats.Fee)
	state.TotalSupply = state.TotalSupply.Sub(blockStats.SupplyChange)

	if err := tx.Update(ctx, &state); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}

	return nil
}
