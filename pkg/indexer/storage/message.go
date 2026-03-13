// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

var errCantFindAddress = errors.New("can't find address")

func saveNamespaceMessages(
	ctx context.Context,
	tx storage.Transaction,
	msgs []*storage.NamespaceMessage,
) error {
	for i := range msgs {
		if msgs[i].Namespace == nil {
			return errors.New("nil namespace in namespace message")
		}
		msgs[i].NamespaceId = msgs[i].Namespace.Id
	}
	return tx.SaveNamespaceMessage(ctx, msgs...)
}

func saveAddressMessage(
	ctx context.Context,
	tx storage.Transaction,
	msgs []*storage.MsgAddress,
	addrToId map[string]uint64,
) error {
	if len(msgs) == 0 {
		return nil
	}

	for i := range msgs {
		if msgs[i].Address == nil {
			return errors.New("nil address in msg_address")
		}
		id, ok := addrToId[msgs[i].Address.String()]
		if !ok {
			return errors.Errorf("unknown address in msg_address: %s", msgs[i].Address.String())
		}
		msgs[i].AddressId = id
	}

	return tx.SaveMsgAddresses(ctx, msgs...)
}

func (module *Module) saveMessages(
	ctx context.Context,
	tx storage.Transaction,
	messages []*storage.Message,
) error {
	if len(messages) == 0 {
		return nil
	}

	if err := tx.SaveMessages(ctx, messages...); err != nil {
		return err
	}

	var (
		valMsgs []storage.MsgValidator
	)
	for i := range messages {
		if len(messages[i].Validators) > 0 {
			for _, val := range messages[i].Validators {
				id, ok := module.validatorsByAddress[val]
				if !ok {
					return errors.Errorf("validator %s not found", val)
				}

				valMsgs = append(valMsgs, storage.MsgValidator{
					Height:      messages[i].Height,
					Time:        messages[i].Time,
					MsgId:       messages[i].Id,
					ValidatorId: id,
				})
			}
		}
	}
	if err := tx.SaveMsgValidator(ctx, valMsgs...); err != nil {
		return errors.Wrap(err, "saving message validators")
	}

	return nil
}
