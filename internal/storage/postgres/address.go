// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
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
	addressQuery := a.DB().NewSelect().
		Model((*storage.Address)(nil)).
		Where("hash = ?", hash).
		Order("id asc").
		Limit(1)

	err = a.DB().NewSelect().TableExpr("(?) as address", addressQuery).
		ColumnExpr("address.*").
		ColumnExpr("balance.currency as balance__currency, balance.spendable as balance__spendable, balance.delegated as balance__delegated, balance.unbonding as balance__unbonding").
		Join("left join balance on balance.id = address.id").
		Scan(ctx, &address)
	return
}

func (a *Address) ListWithBalance(ctx context.Context, filters storage.AddressListFilter) (result []storage.Address, err error) {
	addressQuery := a.DB().NewSelect().Model((*storage.Address)(nil)).
		Offset(filters.Offset)
	addressQuery = addressListFilter(addressQuery, filters)

	err = a.DB().NewSelect().
		TableExpr("(?) as address", addressQuery).
		ColumnExpr("address.*").
		ColumnExpr("balance.currency as balance__currency, balance.spendable as balance__spendable, balance.delegated as balance__delegated, balance.unbonding as balance__unbonding").
		Join("left join balance on balance.id = address.id").
		Scan(ctx, &result)
	return
}
