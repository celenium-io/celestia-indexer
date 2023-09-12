package storage

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
)

func (module *Module) saveMessages(
	ctx context.Context,
	tx postgres.Transaction,
	messages []*storage.Message,
	addrToId map[string]uint64,
) error {
	if err := tx.SaveMessages(ctx, messages...); err != nil {
		return nil
	}

	var (
		namespaceMsgs []storage.NamespaceMessage
		msgAddress    []storage.MsgAddress
		validators    = make([]*storage.Validator, 0)
	)
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
		return nil
	}
	if err := tx.SaveValidators(ctx, validators...); err != nil {
		return nil
	}
	if err := tx.SaveMsgAddresses(ctx, msgAddress...); err != nil {
		return nil
	}

	return nil
}
