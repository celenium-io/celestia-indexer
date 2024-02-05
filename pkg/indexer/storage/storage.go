// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
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
	storage     sdk.Transactable
	notificator storage.Notificator
	validators  map[string]uint64
	indexerName string
}

var _ modules.Module = (*Module)(nil)

// NewModule -
func NewModule(
	storage sdk.Transactable,
	notificator storage.Notificator,
	cfg config.Indexer,
) Module {
	m := Module{
		BaseModule:  modules.New("storage"),
		storage:     storage,
		notificator: notificator,
		validators:  make(map[string]uint64),
		indexerName: cfg.Name,
	}

	m.CreateInputWithCapacity(InputName, 16)
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
				module.MustOutput(StopOutput).Push(struct{}{})
				continue
			}
			block, ok := msg.(storage.Block)
			if !ok {
				module.Log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			state, err := module.saveBlock(ctx, &block)
			if err != nil {
				module.Log.Err(err).
					Uint64("height", uint64(block.Height)).
					Msg("block saving error")
				module.MustOutput(StopOutput).Push(struct{}{})
				continue
			}

			if err := module.notify(ctx, state, block); err != nil {
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

func (module *Module) saveBlock(ctx context.Context, block *storage.Block) (storage.State, error) {
	start := time.Now()
	module.Log.Info().Uint64("height", uint64(block.Height)).Msg("saving block...")
	tx, err := postgres.BeginTransaction(ctx, module.storage)
	if err != nil {
		return storage.State{}, err
	}
	defer tx.Close(ctx)

	state, err := module.processBlockInTransaction(ctx, tx, block)
	if err != nil {
		return state, tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return state, tx.HandleError(ctx, err)
	}
	module.Log.Info().
		Uint64("height", uint64(block.Height)).
		Time("block_time", block.Time).
		Int64("block_ns_size", block.Stats.BlobsSize).
		Str("block_fee", block.Stats.Fee.String()).
		Int64("ms", time.Since(start).Milliseconds()).
		Int("tx_count", len(block.Txs)).
		Msg("block saved")
	return state, nil
}

func (module *Module) processBlockInTransaction(ctx context.Context, tx storage.Transaction, block *storage.Block) (storage.State, error) {
	state, err := tx.State(ctx, module.indexerName)
	if err != nil {
		return state, err
	}
	block.Stats.BlockTime = uint64(block.Time.Sub(state.LastTime).Milliseconds())

	if len(module.validators) > 0 {
		if id, ok := module.validators[block.ProposerAddress]; ok {
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
		addresses  = make(map[string]*storage.Address, 0)
	)

	for i := range block.Addresses {
		key := block.Addresses[i].String()
		if addr, ok := addresses[key]; !ok {
			addresses[key] = &block.Addresses[i]
		} else {
			addr.Balance.Total = addr.Balance.Total.Add(block.Addresses[i].Balance.Total)
		}
	}

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

		for j := range block.Txs[i].Signers {
			key := block.Txs[i].Signers[j].String()
			if _, ok := addresses[key]; !ok {
				addresses[key] = &block.Txs[i].Signers[j]
			}
		}
	}
	if len(events) > 0 {
		if err := tx.SaveEvents(ctx, events...); err != nil {
			return state, err
		}
	}

	addrToId, totalAccounts, err := saveAddresses(ctx, tx, addresses)
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

	totalValidators, err := module.saveMessages(ctx, tx, messages, addrToId)
	if err != nil {
		return state, err
	}

	if err := module.saveBlockSignatures(ctx, tx, block.BlockSignatures, block.Height); err != nil {
		return state, err
	}

	updateState(block, totalAccounts, totalNamespaces, totalValidators, &state)
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

	blockId := strconv.FormatUint(block.Id, 10)
	if err := module.notificator.Notify(ctx, storage.ChannelBlock, blockId); err != nil {
		return err
	}

	return nil
}
