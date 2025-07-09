// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
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
		Where("hash = ?", hash)

	err = a.DB().NewSelect().TableExpr("(?) as address", addressQuery).
		ColumnExpr("address.*").
		ColumnExpr("celestial.id as celestials__id, celestial.image_url as celestials__image_url").
		ColumnExpr("balance.currency as balance__currency, balance.spendable as balance__spendable, balance.delegated as balance__delegated, balance.unbonding as balance__unbonding").
		Join("left join balance on balance.id = address.id").
		Join("left join celestial on celestial.address_id = address.id and celestial.status = 'PRIMARY'").
		Scan(ctx, &address)
	return
}

func (a *Address) ListWithBalance(ctx context.Context, filters storage.AddressListFilter) (result []storage.Address, err error) {
	if filters.SortField == "last_height" || filters.SortField == "first_height" {
		addressQuery := a.DB().NewSelect().
			Model((*storage.Address)(nil))

		addressQuery = addressListFilter(addressQuery, filters)

		err = a.DB().NewSelect().
			TableExpr("(?) as address", addressQuery).
			ColumnExpr("address.*").
			ColumnExpr("celestial.id as celestials__id, celestial.image_url as celestials__image_url").
			ColumnExpr("balance.currency as balance__currency, balance.spendable as balance__spendable, balance.delegated as balance__delegated, balance.unbonding as balance__unbonding").
			Join("left join balance on balance.id = address.id").
			Join("left join celestial on celestial.address_id = address.id and celestial.status = 'PRIMARY'").
			Scan(ctx, &result)

	} else {
		addressQuery := a.DB().NewSelect().
			Model((*storage.Balance)(nil))

		addressQuery = addressListFilter(addressQuery, filters)

		err = a.DB().NewSelect().
			TableExpr("(?) as balance", addressQuery).
			ColumnExpr("address.*").
			ColumnExpr("celestial.id as celestials__id, celestial.image_url as celestials__image_url").
			ColumnExpr("balance.currency as balance__currency, balance.spendable as balance__spendable, balance.delegated as balance__delegated, balance.unbonding as balance__unbonding").
			Join("left join address on balance.id = address.id").
			Join("left join celestial on celestial.address_id = address.id and celestial.status = 'PRIMARY'").
			Scan(ctx, &result)
	}

	return
}

func (a *Address) Series(ctx context.Context, addressId uint64, timeframe storage.Timeframe, column string, req storage.SeriesRequest) (items []storage.HistogramItem, err error) {
	query := a.DB().NewSelect().
		Where("address_id = ?", addressId).
		Order("time desc")

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
	case "tx_count":
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
func (a *Address) IdByHash(ctx context.Context, hash ...[]byte) (id []uint64, err error) {
	err = a.DB().NewSelect().
		Model((*storage.Address)(nil)).
		Column("id").
		Where("hash IN (?)", bun.In(hash)).
		Scan(ctx, &id)
	return
}

// IdByAddress -
func (a *Address) IdByAddress(ctx context.Context, address string, ids ...uint64) (id uint64, err error) {
	query := a.DB().NewSelect().
		Model((*storage.Address)(nil)).
		Column("id").
		Where("address = ?", address)
	if len(ids) > 0 {
		query = query.Where("id IN (?)", bun.In(ids))
	}
	err = query.Scan(ctx, &id)
	return
}
