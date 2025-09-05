// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

const (
	InputName  = "data"
	StopOutput = "stop"
)

// Module - saves received from input block to storage.
//
//	                     |----------------|
//	                     |                |
//	-- storage.Block ->  |     MODULE     |
//	                     |                |
//	                     |----------------|
type Module struct {
	modules.BaseModule
	storage                 sdk.Transactable
	constants               storage.IConstant
	validators              storage.IValidator
	notificator             storage.Notificator
	validatorsByConsAddress map[string]uint64
	validatorsByAddress     map[string]uint64
	validatorsByDelegator   map[string]uint64

	slashingForDowntime   decimal.Decimal
	slashingForDoubleSign decimal.Decimal
	indexerName           string
}

var _ modules.Module = (*Module)(nil)

// NewModule -
func NewModule(
	tx sdk.Transactable,
	constants storage.IConstant,
	validators storage.IValidator,
	notificator storage.Notificator,
	cfg config.Indexer,
) Module {
	m := Module{
		BaseModule:              modules.New("storage"),
		storage:                 tx,
		constants:               constants,
		validators:              validators,
		notificator:             notificator,
		validatorsByConsAddress: make(map[string]uint64),
		validatorsByAddress:     make(map[string]uint64),
		validatorsByDelegator:   make(map[string]uint64),
		slashingForDowntime:     decimal.Zero,
		slashingForDoubleSign:   decimal.Zero,
		indexerName:             cfg.Name,
	}

	m.CreateInputWithCapacity(InputName, 128)
	m.CreateOutput(StopOutput)

	return m
}

// Start -
func (module *Module) Start(ctx context.Context) {
	if err := module.init(ctx); err != nil {
		panic(err)
	}
	module.G.GoCtx(ctx, module.listen)
}

func (module *Module) init(ctx context.Context) error {
	var (
		limit  = 100
		offset = 0
		end    = false
	)

	for !end {
		validators, err := module.validators.List(ctx, uint64(limit), uint64(offset), sdk.SortOrderDesc)
		if err != nil {
			return err
		}
		for i := range validators {
			module.validatorsByConsAddress[validators[i].ConsAddress] = validators[i].Id
			module.validatorsByAddress[validators[i].Address] = validators[i].Id
			module.validatorsByDelegator[validators[i].Delegator] = validators[i].Id
		}
		offset += len(validators)
		end = limit > len(validators)
	}

	return module.initConstants(ctx)
}

func (module *Module) isConstantsEmpty() bool {
	return module.slashingForDoubleSign.IsZero() || module.slashingForDowntime.IsZero()
}

func (module *Module) initConstants(ctx context.Context) error {
	doubleSign, err := module.constants.Get(ctx, types.ModuleNameSlashing, "slash_fraction_double_sign")
	if err != nil {
		if module.validators.IsNoRows(err) {
			return nil
		}
		return err
	}
	module.slashingForDoubleSign, err = decimal.NewFromString(doubleSign.Value)
	if err != nil {
		return err
	}

	downtime, err := module.constants.Get(ctx, types.ModuleNameSlashing, "slash_fraction_downtime")
	if err != nil {
		if module.validators.IsNoRows(err) {
			return nil
		}
		return err
	}
	module.slashingForDowntime, err = decimal.NewFromString(downtime.Value)
	if err != nil {
		return err
	}
	return nil
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
				module.MustOutput(StopOutput).Push(struct{}{})
				continue
			}
			decodedContext, ok := msg.(*decodeContext.Context)
			if !ok {
				module.Log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			if module.isConstantsEmpty() {
				if err := module.initConstants(ctx); err != nil {
					module.Log.Warn().Err(err).Msgf("constant initialization error")
					continue
				}
			}

			state, err := module.saveBlock(ctx, decodedContext)
			if err != nil {
				module.Log.Err(err).
					Uint64("height", uint64(decodedContext.Block.Height)).
					Msg("block saving error")
				module.MustOutput(StopOutput).Push(struct{}{})
				continue
			}

			if err := module.notify(ctx, state, *decodedContext.Block); err != nil {
				module.Log.Err(err).Msg("block notification error")
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

func (module *Module) saveBlock(ctx context.Context, dCtx *decodeContext.Context) (storage.State, error) {
	start := time.Now()
	module.Log.Info().Uint64("height", uint64(dCtx.Block.Height)).Msg("saving block...")
	tx, err := postgres.BeginTransaction(ctx, module.storage)
	if err != nil {
		return storage.State{}, err
	}
	defer tx.Close(ctx)

	state, err := module.processBlockInTransaction(ctx, tx, dCtx)
	if err != nil {
		return state, tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return state, tx.HandleError(ctx, err)
	}
	module.Log.Info().
		Uint64("height", uint64(dCtx.Block.Height)).
		Time("block_time", dCtx.Block.Time).
		Int64("block_ns_size", dCtx.Block.Stats.BlobsSize).
		Str("block_fee", dCtx.Block.Stats.Fee.String()).
		Int64("ms", time.Since(start).Milliseconds()).
		Int("tx_count", len(dCtx.Block.Txs)).
		Msg("block saved")
	return state, nil
}

func (module *Module) processBlockInTransaction(ctx context.Context, tx storage.Transaction, dCtx *decodeContext.Context) (storage.State, error) {
	block := dCtx.Block

	state, err := tx.State(ctx, module.indexerName)
	if err != nil {
		return state, err
	}

	if block.Height == 1 {
		// init after genesis block
		if err := module.init(ctx); err != nil {
			return state, err
		}
	}

	block.Stats.BlockTime = uint64(block.Time.Sub(state.LastTime).Milliseconds())

	if len(module.validatorsByConsAddress) > 0 {
		if id, ok := module.validatorsByConsAddress[block.ProposerAddress]; ok {
			block.ProposerId = id
		} else {
			return state, errors.Errorf("unknown block proposer: %s", block.ProposerAddress)
		}
	} else {
		proposerId, err := tx.GetProposerId(ctx, block.ProposerAddress)
		if err != nil {
			return state, errors.Wrap(err, "can't find block proposer")
		}
		block.ProposerId = proposerId
	}

	if err := tx.Add(ctx, block); err != nil {
		return state, err
	}

	if err := tx.Add(ctx, &block.Stats); err != nil {
		return state, err
	}

	if err := tx.SaveTransactions(ctx, block.Txs...); err != nil {
		return state, err
	}

	if err := tx.SaveEvents(ctx, block.Events...); err != nil {
		return state, err
	}

	var (
		messages   = make([]*storage.Message, 0)
		namespaces = make(map[string]*storage.Namespace, 0)
	)

	events := make([]storage.Event, 0, 10000)
	for i := range block.Txs {
		for j := range block.Txs[i].Messages {
			block.Txs[i].Messages[j].TxId = block.Txs[i].Id
			messages = append(messages, &block.Txs[i].Messages[j])
			setNamespacesFromMessage(block.Txs[i].Messages[j], namespaces)
		}

		for j := range block.Txs[i].Events {
			block.Txs[i].Events[j].TxId = &block.Txs[i].Id
		}

		events = append(events, block.Txs[i].Events...)
		if len(events) >= 10000 {
			if err := tx.SaveEvents(ctx, events...); err != nil {
				return state, err
			}
			events = make([]storage.Event, 0, 10000)
		}
	}
	if len(events) > 0 {
		if err := tx.SaveEvents(ctx, events...); err != nil {
			return state, err
		}
	}

	addrToId, totalAccounts, err := saveAddresses(ctx, tx, dCtx.GetAddresses())
	if err != nil {
		return state, err
	}

	if err := saveSigners(ctx, tx, addrToId, block.Txs); err != nil {
		return state, err
	}

	totalNamespaces, err := saveNamespaces(ctx, tx, namespaces)
	if err != nil {
		return state, err
	}

	ibcClientsCount, err := module.saveMessages(ctx, tx, messages, addrToId)
	if err != nil {
		return state, err
	}

	totalValidators, err := module.saveValidators(ctx, tx, dCtx.GetValidators(), dCtx.Jails)
	if err != nil {
		return state, err
	}

	if err := module.saveDelegations(ctx, tx, dCtx, addrToId); err != nil {
		return state, err
	}

	if err := module.saveBlockSignatures(ctx, tx, block.BlockSignatures, block.Height); err != nil {
		return state, err
	}

	totalProposals, err := module.saveProposals(ctx, tx, dCtx.Block.Height, dCtx.GetProposals(), dCtx.Votes, addrToId)
	if err != nil {
		return state, err
	}

	if err := updateState(block, totalAccounts, totalNamespaces, totalProposals, ibcClientsCount, totalValidators, &state); err != nil {
		return state, err
	}

	err = tx.Update(ctx, &state)
	return state, err
}

func (module *Module) notify(ctx context.Context, state storage.State, block storage.Block) error {
	if time.Since(block.Time) > time.Hour {
		// do not notify all about events if initial indexing is in progress
		return nil
	}

	rawState, err := jsoniter.MarshalToString(state)
	if err != nil {
		return err
	}
	if err := module.notificator.Notify(ctx, storage.ChannelHead, rawState); err != nil {
		return err
	}

	rawBlock, err := jsoniter.MarshalToString(block)
	if err != nil {
		return err
	}
	if err := module.notificator.Notify(ctx, storage.ChannelBlock, rawBlock); err != nil {
		return err
	}

	return nil
}
