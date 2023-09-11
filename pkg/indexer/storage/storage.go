package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	storage     postgres.Storage
	input       *modules.Input
	output      *modules.Output
	indexerName string
	log         zerolog.Logger
	g           workerpool.Group
}

// NewModule -
func NewModule(pg postgres.Storage, cfg config.Indexer) Module {
	m := Module{
		storage:     pg,
		input:       modules.NewInput(InputName),
		output:      modules.NewOutput(StopOutput),
		indexerName: cfg.Name,
		g:           workerpool.NewGroup(),
	}
	m.log = log.With().Str("module", m.Name()).Logger()

	return m
}

// Name -
func (*Module) Name() string {
	return "storage"
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
			block, ok := msg.(storage.Block)
			if !ok {
				module.log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			if err := module.saveBlock(ctx, &block); err != nil {
				module.log.Err(err).Msg("block saving error")
				continue
			}

			if err := module.notify(ctx, block); err != nil {
				module.log.Err(err).Msg("block notification error")
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
	if name != StopOutput {
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

func (module *Module) updateState(block *storage.Block, totalAccounts uint64, state *storage.State) {
	if types.Level(block.Id) <= state.LastHeight {
		return
	}

	state.LastHeight = block.Height
	state.LastHash = block.Hash
	state.LastTime = block.Time
	state.TotalTx += block.Stats.TxCount
	state.TotalAccounts += totalAccounts
	state.TotalBlobsSize = block.Stats.BlobsSize
	state.TotalFee = state.TotalFee.Add(block.Stats.Fee)
	state.TotalSupply = state.TotalSupply.Add(block.Stats.SupplyChange)
	state.ChainId = block.ChainId
}

func (module *Module) saveBlock(ctx context.Context, block *storage.Block) error {
	start := time.Now()
	module.log.Info().Uint64("height", uint64(block.Height)).Msg("saving block...")
	tx, err := postgres.BeginTransaction(ctx, module.storage.Transactable)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	state, err := tx.State(ctx, module.indexerName)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, &block.Stats); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.SaveTransactions(ctx, block.Txs...); err != nil {
		return tx.HandleError(ctx, err)
	}

	var (
		messages   = make([]any, 0)
		events     = make([]any, len(block.Events))
		namespaces = make(map[string]*storage.Namespace, 0)
		addresses  = make(map[string]*storage.Address, 0)

		totalAccounts uint64
	)

	for i := range block.Events {
		events[i] = &block.Events[i]
	}

	for i := range block.Txs {
		for j := range block.Txs[i].Messages {
			block.Txs[i].Messages[j].TxId = block.Txs[i].Id
			messages = append(messages, &block.Txs[i].Messages[j])

			for k := range block.Txs[i].Messages[j].Namespace {
				key := block.Txs[i].Messages[j].Namespace[k].String()
				if _, ok := namespaces[key]; !ok {
					block.Txs[i].Messages[j].Namespace[k].PfbCount = 1
					namespaces[key] = &block.Txs[i].Messages[j].Namespace[k]
				}
			}
		}

		for j := range block.Txs[i].Events {
			block.Txs[i].Events[j].TxId = &block.Txs[i].Id
			events = append(events, &block.Txs[i].Events[j])
		}

		for j := range block.Txs[i].Addresses {
			key := block.Txs[i].Addresses[j].String()
			if addr, ok := addresses[key]; !ok {
				addresses[key] = &block.Txs[i].Addresses[j].Address
			} else {
				addr.Balance.Total = addr.Balance.Total.Add(block.Txs[i].Addresses[j].Address.Balance.Total)
			}
		}
	}

	if len(addresses) > 0 {
		data := make([]*storage.Address, 0, len(addresses))
		for key := range addresses {
			data = append(data, addresses[key])
		}

		if err := tx.SaveAddresses(ctx, data...); err != nil {
			return tx.HandleError(ctx, err)
		}

		balances := make([]storage.Balance, 0)
		for i := range data {
			data[i].Balance.Id = data[i].Id
			balances = append(balances, data[i].Balance)
		}
		if err := tx.SaveBalances(ctx, balances...); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if len(namespaces) > 0 {
		data := make([]*storage.Namespace, 0, len(namespaces))
		for key := range namespaces {
			data = append(data, namespaces[key])
		}

		if err := tx.SaveNamespaces(ctx, data...); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if len(messages) > 0 {
		if err := tx.BulkSave(ctx, messages); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if len(events) > 0 {
		if err := tx.BulkSave(ctx, events); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	var (
		namespaceMsgs []storage.NamespaceMessage
		validators    = make([]*storage.Validator, 0)
	)
	for i := range messages {
		msg, ok := messages[i].(*storage.Message)
		if !ok {
			continue
		}
		for j := range msg.Namespace {
			if msg.Namespace[j].Id == 0 { // in case of duplication of writing to one namespace inside one messages
				continue
			}
			namespaceMsgs = append(namespaceMsgs, storage.NamespaceMessage{
				MsgId:       msg.Id,
				NamespaceId: msg.Namespace[j].Id,
				Time:        msg.Time,
				Height:      msg.Height,
				TxId:        msg.TxId,
			})
		}

		if msg.Validator != nil {
			msg.Validator.MsgId = msg.Id
			validators = append(validators, msg.Validator)
		}
	}
	if err := tx.SaveNamespaceMessage(ctx, namespaceMsgs...); err != nil {
		return tx.HandleError(ctx, err)
	}
	if err := tx.SaveValidators(ctx, validators...); err != nil {
		return tx.HandleError(ctx, err)
	}

	var txAddresses []storage.TxAddress
	for _, transaction := range block.Txs {
		for _, address := range transaction.Addresses {
			txAddresses = append(txAddresses, storage.TxAddress{
				TxId:      transaction.Id,
				AddressId: address.Id,
				Type:      address.Type,
			})
		}
	}

	if err := tx.SaveTxAddresses(ctx, txAddresses...); err != nil {
		return tx.HandleError(ctx, err)
	}

	module.updateState(block, totalAccounts, &state)
	if err := tx.Update(ctx, &state); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	module.log.Info().
		Uint64("height", block.Id).
		Uint64("block_ns_size", block.Stats.BlobsSize).
		Str("block_fee", block.Stats.Fee.String()).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block saved")
	return nil
}

func (module *Module) notify(ctx context.Context, block storage.Block) error {
	blockId := strconv.FormatUint(block.Id, 10)
	if err := module.storage.Notificator.Notify(ctx, storage.ChannelHead, blockId); err != nil {
		return err
	}

	for i := range block.Txs {
		txId := strconv.FormatUint(block.Txs[i].Id, 10)
		if err := module.storage.Notificator.Notify(ctx, storage.ChannelTx, txId); err != nil {
			return err
		}
	}
	return nil
}
