// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
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
	addressQuery := a.DB().NewSelect().
		Model((*storage.Balance)(nil))

	addressQuery = addressListFilter(addressQuery, filters)

	err = a.DB().NewSelect().
		TableExpr("(?) as balance", addressQuery).
		ColumnExpr("address.*").
		ColumnExpr("balance.currency as balance__currency, balance.spendable as balance__spendable, balance.delegated as balance__delegated, balance.unbonding as balance__unbonding").
		Join("left join address on balance.id = address.id").
		Scan(ctx, &result)
	return
}

func (a *Address) Series(ctx context.Context, addressId uint64, timeframe storage.Timeframe, column string, req storage.SeriesRequest) (items []storage.HistogramItem, err error) {
	query := a.DB().NewSelect().
		Where("address_id = ?", addressId).
		Order("time desc").
		Limit(100)

	switch timeframe {
	case storage.TimeframeHour:
		query = query.Table("accounts_tx_by_hour")
	case storage.TimeframeDay:
		query = query.Table("accounts_tx_by_day")
	case storage.TimeframeMonth:
		query = query.Table("accounts_tx_by_month")
	default:
		return nil, errors.Errorf("invalid timeframe: %s", timeframe)
	}

	switch column {
	case "gas_used":
		query = query.ColumnExpr("gas_used as value, time as bucket")
	case "gas_wanted":
		query = query.ColumnExpr("gas_wanted as value, time as bucket")
	case "count":
		query = query.ColumnExpr("count as value, time as bucket")
	case "fee":
		query = query.ColumnExpr("fee as value, time as bucket")
	default:
		return nil, errors.Errorf("invalid column: %s", column)
	}

	if !req.From.IsZero() {
		query = query.Where("time >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("time < ?", req.To)
	}

	err = query.Scan(ctx, &items)

	return
}

// IdByHash -
func (a *Address) IdByHash(ctx context.Context, hash []byte) (id uint64, err error) {
	err = a.DB().NewSelect().
		Model((*storage.Address)(nil)).
		Column("id").
		Where("hash = ?", hash).
		Limit(1).
		Scan(ctx, &id)
	return
}
