// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

func saveZkIsm(
	ctx context.Context,
	tx storage.Transaction,
	zkism []*storage.ZkISM,
	addrToId map[string]uint64,
) error {
	if len(zkism) == 0 {
		return nil
	}

	for i := range zkism {
		if zkism[i].Creator != nil {
			if addrId, ok := addrToId[zkism[i].Creator.Address]; ok {
				zkism[i].CreatorId = addrId
			}
		}
	}

	return tx.SaveZkISMs(ctx, zkism...)
}

func saveZkIsmUpdates(
	ctx context.Context,
	tx storage.Transaction,
	zkism *sync.Map[string, *storage.ZkISM],
	updates []*storage.ZkISMUpdate,
	addrToId map[string]uint64,
) error {
	if len(updates) == 0 {
		return nil
	}

	for i := range updates {
		if updates[i].Signer != nil {
			if addrId, ok := addrToId[updates[i].Signer.Address]; ok {
				updates[i].SignerId = addrId
			}
		}

		if item, ok := zkism.Get(updates[i].ExternalId()); !ok {
			// ISM from a previous block: resolve its DB id now.
			dbIsm, err := tx.ZkISMById(ctx, updates[i].ZkISMExternalId)
			if err != nil {
				return errors.Wrapf(err, "can't find zk ism for update: external_id=%s", updates[i].ExternalId())
			}
			updates[i].ZkISMId = dbIsm.Id
		} else {
			updates[i].ZkISMId = item.Id
		}
	}

	return tx.SaveZkISMUpdates(ctx, updates...)
}

func saveZkIsmMessages(
	ctx context.Context,
	tx storage.Transaction,
	zkism *sync.Map[string, *storage.ZkISM],
	msgs []*storage.ZkISMMessage,
	addrToId map[string]uint64,
) error {
	if len(msgs) == 0 {
		return nil
	}

	for i := range msgs {
		if msgs[i].Signer != nil {
			if addrId, ok := addrToId[msgs[i].Signer.Address]; ok {
				msgs[i].SignerId = addrId
			}
		}

		if item, ok := zkism.Get(msgs[i].ExternalId()); !ok {
			// ISM from a previous block: resolve its DB id now.
			dbIsm, err := tx.ZkISMById(ctx, msgs[i].ZkISMExternalId)
			if err != nil {
				return errors.Wrapf(err, "can't find zk ism for update: external_id=%s", msgs[i].ExternalId())
			}
			msgs[i].ZkISMId = dbIsm.Id
		} else {
			msgs[i].ZkISMId = item.Id
		}
	}

	return tx.SaveZkISMMessages(ctx, msgs...)
}
