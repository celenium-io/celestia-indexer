// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/uptrace/bun"

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
	err = a.DB().NewSelect().Model(&address).
		Where("hash = ?", hash).
		Relation("Balance").
		Scan(ctx)
	return
}

func (a *Address) ListWithBalance(ctx context.Context, filters storage.AddressListFilter) (result []storage.Address, err error) {
	addressQuery := a.DB().NewSelect().Model((*storage.Address)(nil)).
		Offset(filters.Offset)
	addressQuery = addressListFilter(addressQuery, filters)

	err = a.DB().NewSelect().
		TableExpr("(?) as address", addressQuery).
		ColumnExpr("address.*").
		ColumnExpr("balance.currency as balance__currency, balance.total as balance__total").
		Join("left join balance on balance.id = address.id").
		Scan(ctx, &result)
	return
}

func (a *Address) Messages(ctx context.Context, id uint64, filters storage.AddressMsgsFilter) (msgs []storage.MsgAddress, err error) {
	query := a.DB().NewSelect().Model(&msgs).
		Where("address_id = ?", id).
		Offset(filters.Offset).
		Relation("Msg")

	query = addressMsgsFilter(query, filters)

	err = query.Scan(ctx)
	return
}

func addressMsgsFilter(query *bun.SelectQuery, filters storage.AddressMsgsFilter) *bun.SelectQuery {
	query = limitScope(query, filters.Limit)
	query = sortScope(query, "msg_id", filters.Sort)
	return query
}
