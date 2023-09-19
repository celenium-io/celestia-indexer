package storage

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

func saveMessages(
	ctx context.Context,
	tx storage.Transaction,
	messages []*storage.Message,
	addrToId map[string]uint64,
) error {
	if err := tx.SaveMessages(ctx, messages...); err != nil {
		return err
	}

	var (
		namespaceMsgs []storage.NamespaceMessage
		msgAddress    []storage.MsgAddress
		validators    = make([]*storage.Validator, 0)
		namespaces    = make(map[string]uint64)
		addedMsgId    = make(map[uint64]struct{})
	)
	for i := range messages {
		for _, ns := range messages[i].Namespace {
			nsId := ns.Id
			if nsId == 0 {
				if _, ok := addedMsgId[messages[i].Id]; ok { // in case of duplication of writing to one namespace inside one messages
					continue
				}

				id, ok := namespaces[ns.String()]
				if !ok {
					continue
				}
				nsId = id
			} else {
				namespaces[ns.String()] = nsId
			}

			addedMsgId[messages[i].Id] = struct{}{}
			namespaceMsgs = append(namespaceMsgs, storage.NamespaceMessage{
				MsgId:       messages[i].Id,
				NamespaceId: nsId,
				Time:        messages[i].Time,
				Height:      messages[i].Height,
				TxId:        messages[i].TxId,
			})
		}

		if messages[i].Validator != nil {
			messages[i].Validator.MsgId = messages[i].Id
			validators = append(validators, messages[i].Validator)
		}

		for j := range messages[i].Addresses {
			id, ok := addrToId[messages[i].Addresses[j].String()]
			if !ok {
				continue
			}
			msgAddress = append(msgAddress, storage.MsgAddress{
				MsgId:     messages[i].Id,
				AddressId: id,
				Type:      messages[i].Addresses[j].Type,
			})
		}
	}

	if err := tx.SaveNamespaceMessage(ctx, namespaceMsgs...); err != nil {
		return err
	}
	if err := tx.SaveValidators(ctx, validators...); err != nil {
		return err
	}
	if err := tx.SaveMsgAddresses(ctx, msgAddress...); err != nil {
		return err
	}

	return nil
}
