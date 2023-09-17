package storage

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
)

func (module *Module) saveAddresses(
	ctx context.Context,
	tx postgres.Transaction,
	addresses map[string]*storage.Address,
) (map[string]uint64, uint64, error) {
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
	balances := make([]storage.Balance, 0)
	for i := range data {
		addToId[data[i].Address] = data[i].Id
		data[i].Balance.Id = data[i].Id
		balances = append(balances, data[i].Balance)
	}
	err = tx.SaveBalances(ctx, balances...)
	return addToId, totalAccounts, err
}

func (module *Module) saveSigners(
	ctx context.Context,
	tx postgres.Transaction,
	addrToId map[string]uint64,
	txs []storage.Tx,
) error {
	if len(txs) == 0 || len(addrToId) == 0 {
		return nil
	}

	var txAddresses []storage.Signer
	for _, transaction := range txs {
		for _, signer := range transaction.Signers {
			if addrId, ok := addrToId[signer.Address]; ok {
				txAddresses = append(txAddresses, storage.Signer{
					TxId:      transaction.Id,
					AddressId: addrId,
				})
			}
		}
	}
	return tx.SaveSigners(ctx, txAddresses...)
}
