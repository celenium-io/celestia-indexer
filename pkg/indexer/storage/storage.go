package storage

import (
	"context"
	"encoding/hex"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"strconv"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InputName -
const InputName = "data"

// Module - saves received from input block to storage.
//
//	                     |----------------|
//	                     |                |
//	-- storage.Block ->  |     MODULE     |
//	                     |                |
//	                     |----------------|
type Module struct {
	storage postgres.Storage
	input   *modules.Input
	state   *storage.State
	log     zerolog.Logger
	g       workerpool.Group
}

// NewModule -
func NewModule(pg postgres.Storage, cfg config.Indexer) Module {
	m := Module{
		storage: pg,
		input:   modules.NewInput(InputName),
		state: &storage.State{
			Name: cfg.Name,
		},
		g: workerpool.NewGroup(),
	}
	m.log = log.With().Str("module", m.Name()).Logger()

	return m
}

// Name -
func (*Module) Name() string {
	return "storage"
}

func (module *Module) initState(ctx context.Context) error {
	module.log.Info().Msg("loading current state from database...")

	state, err := module.storage.State.ByName(ctx, module.state.Name)
	switch {
	case err == nil:
		module.state = &state
		module.log.Info().
			Str("indexer_name", module.state.Name).
			Uint64("height", uint64(module.state.LastHeight)).
			Str("hash", hex.EncodeToString(module.state.LastHash)).
			Time("last_updated", module.state.LastTime).
			Msg("current state")
		return nil

	case module.storage.State.IsNoRows(err):
		module.log.Info().Msg("state is not found. creating empty state...")
		return module.storage.State.Save(ctx, module.state)

	default:
		return errors.Wrap(err, "state loading")
	}
}

// Start -
func (module *Module) Start(ctx context.Context) {
	if err := module.initState(ctx); err != nil {
		module.log.Err(err).Msg("error during storage module initialization")
		return
	}

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

			if err := module.saveBlock(ctx, block); err != nil {
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
	return nil, errors.Wrap(modules.ErrUnknownOutput, name)
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

func (module *Module) updateState(block storage.Block, totalAccounts uint64) {
	if types.Level(block.Id) <= module.state.LastHeight {
		return
	}

	module.state.LastHeight = block.Height
	module.state.LastHash = block.Hash
	module.state.LastTime = block.Time
	module.state.TotalTx += block.TxCount
	module.state.TotalAccounts += totalAccounts
	module.state.TotalBlobsSize = block.BlobsSize
	module.state.TotalFee = module.state.TotalFee.Add(block.Fee)
	module.state.ChainId = block.ChainId
}

func (module *Module) saveBlock(ctx context.Context, block storage.Block) error {
	start := time.Now()
	module.log.Info().Uint64("height", uint64(block.Height)).Msg("saving block...")
	tx, err := postgres.BeginTransaction(ctx, module.storage.Transactable)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	block.Id = uint64(block.Height)
	if err := tx.Add(ctx, &block); err != nil {
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
			if _, ok := addresses[key]; !ok {
				addresses[key] = &block.Txs[i].Addresses[j].Address
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

	var namespaceMsgs []storage.NamespaceMessage
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
	}
	if err := tx.SaveNamespaceMessage(ctx, namespaceMsgs...); err != nil {
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

	module.updateState(block, totalAccounts)
	if err := tx.Update(ctx, module.state); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	module.log.Info().
		Uint64("height", block.Id).
		Uint64("block_ns_size", block.BlobsSize).
		Str("block_fee", block.Fee.String()).
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
