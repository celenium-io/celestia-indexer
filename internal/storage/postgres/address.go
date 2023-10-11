// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Address -
type Address struct {
	*postgres.Table[*storage.Address]
}

// NewAddress -
func NewAddress(db *database.Bun) *Address {
	return &Address{
		Table: postgres.NewTable[*storage.Address](db),
	}
}

// ByHash -
func (a *Address) ByHash(ctx context.Context, hash []byte) (address storage.Address, err error) {
	err = a.DB().NewSelect().Model(&address).
		Where("hash = ?", hash).
		Relation("Balance").
		Scan(ctx)
	return
}

func (a *Address) ListWithBalance(ctx context.Context, fltrs storage.AddressListFilter) (result []storage.Address, err error) {
	query := a.DB().NewSelect().Model(&result).
		Offset(fltrs.Offset).
		Relation("Balance")

	query = addressListFilter(query, fltrs)

	err = query.Scan(ctx)
	return
}
