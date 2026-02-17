// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Forwarding -
type Forwarding struct {
	*postgres.Table[*storage.Forwarding]
}

// NewForwarding -
func NewForwarding(db *database.Bun) *Forwarding {
	return &Forwarding{
		Table: postgres.NewTable[*storage.Forwarding](db),
	}
}

func (f *Forwarding) ById(ctx context.Context, id uint64) (forwarding storage.Forwarding, err error) {
	subQuery := f.DB().NewSelect().
		Model(&forwarding).
		Where("id = ?", id).
		Limit(1)

	err = f.DB().NewSelect().
		TableExpr("(?) as forwarding", subQuery).
		ColumnExpr("forwarding.*").
		ColumnExpr("address.id as address__id, address.address as address__address").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join address on address.id = forwarding.address_id").
		Join("left join tx on tx.id = forwarding.tx_id").
		Scan(ctx, &forwarding)

	return
}

func (f *Forwarding) Filter(ctx context.Context, filters storage.ForwardingFilter) (forwardings []storage.Forwarding, err error) {
	subQuery := f.DB().NewSelect().Model((*storage.Forwarding)(nil))

	if filters.Height != nil {
		subQuery = subQuery.Where("height = ?", *filters.Height)
	}
	if filters.AddressId != nil {
		subQuery = subQuery.Where("address_id = ?", *filters.AddressId)
	}
	if filters.TxId != nil {
		subQuery = subQuery.Where("tx_id = ?", *filters.TxId)
	}

	if filters.Sort == "" {
		filters.Sort = sdk.SortOrderAsc
	}

	if !filters.From.IsZero() {
		subQuery = subQuery.Where("time >= ?", filters.From)
	}
	if !filters.To.IsZero() {
		subQuery = subQuery.Where("time < ?", filters.To)
	}

	subQuery = subQuery.OrderExpr("time ?0, id ?0", bun.Safe(filters.Sort))
	subQuery = limitScope(subQuery, filters.Limit)
	if filters.Offset > 0 {
		subQuery = subQuery.Offset(filters.Offset)
	}

	err = f.DB().NewSelect().
		TableExpr("(?) as forwarding", subQuery).
		ColumnExpr("forwarding.*").
		ColumnExpr("address.id as address__id, address.address as address__address").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join address on address.id = forwarding.address_id").
		Join("left join tx on tx.id = forwarding.tx_id").
		OrderExpr("forwarding.time ?0, forwarding.id ?0", bun.Safe(filters.Sort)).
		Scan(ctx, &forwardings)

	return
}
