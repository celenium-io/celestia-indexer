package rollback

import (
	"bytes"
	"context"

	"github.com/dipdup-io/celestia-indexer/pkg/node"

	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
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
	modules.BaseModule
	tx        sdk.Transactable
	state     storage.IState
	blocks    storage.IBlock
	node      node.API
	indexName string
}

var _ modules.Module = (*Module)(nil)

func NewModule(
	tx sdk.Transactable,
	state storage.IState,
	blocks storage.IBlock,
	node node.API,
	cfg config.Indexer,
) Module {
	module := Module{
		BaseModule: modules.New("rollback"),
		tx:         tx,
		state:      state,
		blocks:     blocks,
		node:       node,
		indexName:  cfg.Name,
	}

	module.CreateInput(InputName)
	module.CreateOutput(OutputName)
	module.CreateOutput(StopOutput)

	return module
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
		case _, ok := <-input.Listen():
			if !ok {
				module.Log.Warn().Msg("can't read message from input")
				return
			}

			if err := module.rollback(ctx); err != nil {
				module.Log.Err(err).Msgf("error occured")
			}
		}
	}
}

// Close -
func (module *Module) Close() error {
	module.Log.Info().Msg("closing module...")
	module.G.Wait()
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
	module.MustInput(OutputName).Push(newState)

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
