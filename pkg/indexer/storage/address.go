// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func saveAddresses(
	ctx context.Context,
	tx storage.Transaction,
	addresses map[string]*storage.Address,
) (map[string]uint64, int64, error) {
	if len(addresses) == 0 {
		return nil, 0, nil
	}

	data := make([]*storage.Address, 0, len(addresses))
	for key := range addresses {
		data = append(data, addresses[key])
	}

	totalAccounts, err := tx.SaveAddresses(ctx, data...)
	if err != nil {
		return nil, 0, err
	}

	addToId := make(map[string]uint64)
	balances := make([]storage.Balance, len(data))
	for i := range data {
		addToId[data[i].Address] = data[i].Id
		data[i].Balance.Id = data[i].Id
		balances[i] = data[i].Balance
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
