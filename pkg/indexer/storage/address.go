// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func saveAddresses(
	ctx context.Context,
	tx storage.Transaction,
	addresses []*storage.Address,
) (map[string]uint64, int64, error) {
	if len(addresses) == 0 {
		return nil, 0, nil
	}

	totalAccounts, err := tx.SaveAddresses(ctx, addresses...)
	if err != nil {
		return nil, 0, err
	}

	addToId := make(map[string]uint64)
	balances := make([]storage.Balance, len(addresses))
	for i := range addresses {
		addToId[addresses[i].Address] = addresses[i].Id
		addresses[i].Balance.Id = addresses[i].Id
		balances[i] = addresses[i].Balance
	}
	err = tx.SaveBalances(ctx, balances...)
	return addToId, totalAccounts, err
}

func saveSigners(
	ctx context.Context,
	tx storage.Transaction,
	addrToId map[string]uint64,
	txs []storage.Tx,
) error {
	if len(txs) == 0 || len(addrToId) == 0 {
		return nil
	}

	var txAddresses []storage.Signer
	for i := range txs {
		for j := range txs[i].Signers {
			if addrId, ok := addrToId[txs[i].Signers[j].Address]; ok {
				txAddresses = append(txAddresses, storage.Signer{
					TxId:      txs[i].Id,
					AddressId: addrId,
				})
			}
		}
	}
	return tx.SaveSigners(ctx, txAddresses...)
}
