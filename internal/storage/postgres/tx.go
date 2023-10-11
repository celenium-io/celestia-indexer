// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Tx -
type Tx struct {
	*postgres.Table[*storage.Tx]
}

// NewTx -
func NewTx(db *database.Bun) *Tx {
	return &Tx{
		Table: postgres.NewTable[*storage.Tx](db),
	}
}

func (tx *Tx) ByHash(ctx context.Context, hash []byte) (transaction storage.Tx, err error) {
	err = tx.DB().NewSelect().Model(&transaction).
		Where("hash = ?", hash).
		Scan(ctx)
	return
}

func (tx *Tx) Filter(ctx context.Context, fltrs storage.TxFilter) (txs []storage.Tx, err error) {
	query := tx.DB().NewSelect().Model(&txs).Offset(fltrs.Offset)
	query = txFilter(query, fltrs)

	err = query.Scan(ctx)
	return
}

func (tx *Tx) ByIdWithRelations(ctx context.Context, id uint64) (transaction storage.Tx, err error) {
	err = tx.DB().NewSelect().Model(&transaction).
		Where("id = ?", id).
		Relation("Messages").
		Scan(ctx)
	return
}

func (tx *Tx) ByAddress(ctx context.Context, addressId uint64, fltrs storage.TxFilter) ([]storage.Tx, error) {
	var relations []storage.Signer
	query := tx.DB().NewSelect().
		Model(&relations).
		Where("address_id = ?", addressId).
		Relation("Tx")

	query = txFilter(query, fltrs)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	transactions := make([]storage.Tx, len(relations))
	for i := range relations {
		transactions[i] = *relations[i].Tx
	}
	return transactions, nil
}

func (tx *Tx) Genesis(ctx context.Context, limit, offset int, sortOrder sdk.SortOrder) (txs []storage.Tx, err error) {
	query := tx.DB().NewSelect().Model(&txs).Offset(offset).Where("hash IS NULL")
	query = limitScope(query, limit)
	query = sortScope(query, "id", sortOrder)

	err = query.Scan(ctx)
	return
}
