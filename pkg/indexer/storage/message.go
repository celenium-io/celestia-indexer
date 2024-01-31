// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

func saveMessages(
	ctx context.Context,
	tx storage.Transaction,
	messages []*storage.Message,
	addrToId map[string]uint64,
) (int, error) {
	if err := tx.SaveMessages(ctx, messages...); err != nil {
		return 0, err
	}

	var (
		namespaceMsgs []storage.NamespaceMessage
		msgAddress    []storage.MsgAddress
		blobLogs      = make([]storage.BlobLog, 0)
		validators    = make([]*storage.Validator, 0)
		namespaces    = make(map[string]uint64)
		addedMsgId    = make(map[uint64]struct{})
		msgAddrMap    = make(map[string]struct{})
	)
	for i := range messages {
		for j := range messages[i].Namespace {
			nsId := messages[i].Namespace[j].Id
			key := messages[i].Namespace[j].String()
			if nsId == 0 {
				if _, ok := addedMsgId[messages[i].Id]; ok { // in case of duplication of writing to one namespace inside one messages
					continue
				}

				id, ok := namespaces[key]
				if !ok {
					continue
				}
				nsId = id
			} else {
				namespaces[key] = nsId
			}

			addedMsgId[messages[i].Id] = struct{}{}
			namespaceMsgs = append(namespaceMsgs, storage.NamespaceMessage{
				MsgId:       messages[i].Id,
				NamespaceId: nsId,
				Time:        messages[i].Time,
				Height:      messages[i].Height,
				TxId:        messages[i].TxId,
				Size:        uint64(messages[i].Namespace[j].Size),
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
			msgAddressEntity := storage.MsgAddress{
				MsgId:     messages[i].Id,
				AddressId: id,
				Type:      messages[i].Addresses[j].Type,
			}
			key := msgAddressEntity.String()
			if _, ok := msgAddrMap[key]; !ok {
				msgAddress = append(msgAddress, msgAddressEntity)
				msgAddrMap[key] = struct{}{}
			}
		}

		for j := range messages[i].BlobLogs {
			if messages[i].BlobLogs[j].Namespace == nil {
				return 0, errors.New("nil namespace in pay for blob message")
			}
			nsId, ok := namespaces[messages[i].BlobLogs[j].Namespace.String()]
			if !ok {
				return 0, errors.Errorf("can't find namespace for pay for blob message: %s", messages[i].BlobLogs[j].Namespace.String())
			}
			if messages[i].BlobLogs[j].Signer == nil {
				return 0, errors.New("nil signer address in pay for blob message")
			}
			signerId, ok := addrToId[messages[i].BlobLogs[j].Signer.Address]
			if !ok {
				return 0, errors.Errorf("can't find signer address for pay for blob message: %s", messages[i].BlobLogs[j].Signer.Address)
			}

			messages[i].BlobLogs[j].MsgId = messages[i].Id
			messages[i].BlobLogs[j].TxId = messages[i].TxId
			messages[i].BlobLogs[j].SignerId = signerId
			messages[i].BlobLogs[j].NamespaceId = nsId

			blobLogs = append(blobLogs, *messages[i].BlobLogs[j])
		}
	}

	if err := tx.SaveNamespaceMessage(ctx, namespaceMsgs...); err != nil {
		return 0, err
	}
	newValidatorsCount, err := tx.SaveValidators(ctx, validators...)
	if err != nil {
		return 0, err
	}
	if err := tx.SaveMsgAddresses(ctx, msgAddress...); err != nil {
		return 0, err
	}
	if err := tx.SaveBlobLogs(ctx, blobLogs...); err != nil {
		return 0, err
	}

	return newValidatorsCount, nil
}
