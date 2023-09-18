package genesis

import (
	"context"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
)

func (module *Module) save(ctx context.Context, data parsedData) error {
	start := time.Now()
	module.Log.Info().Uint64("height", uint64(data.block.Height)).Msg("saving block...")
	tx, err := postgres.BeginTransaction(ctx, module.storage.Transactable)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	if err := tx.SaveConstants(ctx, data.constants...); err != nil {
		return tx.HandleError(ctx, err)
	}

	for i := range data.denomMetadata {
		if err := tx.Add(ctx, &data.denomMetadata[i]); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.Add(ctx, &data.block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, &data.block.Stats); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.SaveTransactions(ctx, data.block.Txs...); err != nil {
		return tx.HandleError(ctx, err)
	}

	var (
		messages   = make([]*storage.Message, 0)
		events     = make([]any, len(data.block.Events))
		namespaces = make(map[string]*storage.Namespace, 0)
	)

	for i := range data.block.Events {
		events[i] = &data.block.Events[i]
	}

	for i := range data.block.Txs {
		for j := range data.block.Txs[i].Messages {
			data.block.Txs[i].Messages[j].TxId = data.block.Txs[i].Id
			messages = append(messages, &data.block.Txs[i].Messages[j])

			for k := range data.block.Txs[i].Messages[j].Namespace {
				key := data.block.Txs[i].Messages[j].Namespace[k].String()
				if _, ok := namespaces[key]; !ok {
					data.block.Txs[i].Messages[j].Namespace[k].PfbCount = 1
					namespaces[key] = &data.block.Txs[i].Messages[j].Namespace[k]
				}
			}
		}

		for j := range data.block.Txs[i].Events {
			data.block.Txs[i].Events[j].TxId = &data.block.Txs[i].Id
			events = append(events, &data.block.Txs[i].Events[j])
		}

		for j := range data.block.Txs[i].Signers {
			key := data.block.Txs[i].Signers[j].String()
			if addr, ok := data.addresses[key]; !ok {
				data.addresses[key] = &data.block.Txs[i].Signers[j]
			} else {
				addr.Balance.Total = addr.Balance.Total.Add(data.block.Txs[i].Signers[j].Balance.Total)
			}
		}
	}

	var totalAccounts uint64
	if len(data.addresses) > 0 {
		entities := make([]*storage.Address, 0, len(data.addresses))
		for key := range data.addresses {
			entities = append(entities, data.addresses[key])
		}

		totalAccounts, err = tx.SaveAddresses(ctx, entities...)
		if err != nil {
			return tx.HandleError(ctx, err)
		}

		balances := make([]storage.Balance, len(entities))
		for i := range entities {
			balances[i] = entities[i].Balance
		}
		if err := tx.SaveBalances(ctx, balances...); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	var totalNamespaces uint64
	if len(namespaces) > 0 {
		entities := make([]*storage.Namespace, 0, len(namespaces))
		for key := range namespaces {
			entities = append(entities, namespaces[key])
		}

		totalNamespaces, err = tx.SaveNamespaces(ctx, entities...)
		if err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.SaveMessages(ctx, messages...); err != nil {
		return tx.HandleError(ctx, err)
	}

	if len(events) > 0 {
		if err := tx.BulkSave(ctx, events); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	var namespaceMsgs []storage.NamespaceMessage
	for i := range messages {
		for j := range messages[i].Namespace {
			if messages[i].Namespace[j].Id == 0 { // in case of duplication of writing to one namespace inside one messages
				continue
			}
			namespaceMsgs = append(namespaceMsgs, storage.NamespaceMessage{
				MsgId:       messages[i].Id,
				NamespaceId: messages[i].Namespace[j].Id,
				Time:        messages[i].Time,
				Height:      messages[i].Height,
				TxId:        messages[i].TxId,
			})
		}
	}
	if err := tx.SaveNamespaceMessage(ctx, namespaceMsgs...); err != nil {
		return tx.HandleError(ctx, err)
	}

	var signers []storage.Signer
	for _, transaction := range data.block.Txs {
		for _, address := range transaction.Signers {
			signers = append(signers, storage.Signer{
				TxId:      transaction.Id,
				AddressId: address.Id,
			})
		}
	}

	if err := tx.SaveSigners(ctx, signers...); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, &storage.State{
		Name:            module.indexerName,
		LastHeight:      data.block.Height,
		LastTime:        data.block.Time,
		LastHash:        data.block.Hash,
		ChainId:         data.block.ChainId,
		TotalTx:         data.block.Stats.TxCount,
		TotalSupply:     data.block.Stats.SupplyChange,
		TotalFee:        data.block.Stats.Fee,
		TotalBlobsSize:  data.block.Stats.BlobsSize,
		TotalAccounts:   totalAccounts,
		TotalNamespaces: totalNamespaces,
	}); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	module.Log.Info().
		Uint64("height", data.block.Id).
		Uint64("block_ns_size", data.block.Stats.BlobsSize).
		Str("block_fee", data.block.Stats.Fee.String()).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block saved")
	return nil
}
