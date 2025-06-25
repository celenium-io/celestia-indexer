// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Redelegation -
type Redelegation struct {
	*postgres.Table[*storage.Redelegation]
}

// NewRedelegation -
func NewRedelegation(db *database.Bun) *Redelegation {
	return &Redelegation{
		Table: postgres.NewTable[*storage.Redelegation](db),
	}
}

func (d *Redelegation) ByAddress(ctx context.Context, addressId uint64, limit, offset int) (redelegations []storage.Redelegation, err error) {
	subQuery := d.DB().NewSelect().Model((*storage.Redelegation)(nil)).
		Where("address_id = ?", addressId).
		Order("amount desc")

	subQuery = limitScope(subQuery, limit)
	if offset > 0 {
		subQuery = subQuery.Offset(offset)
	}

	err = d.DB().NewSelect().
		TableExpr("(?) as redelegation", subQuery).
		ColumnExpr("redelegation.*").
		ColumnExpr("source.id as source__id, source.moniker as source__moniker, source.cons_address as source__cons_address").
		ColumnExpr("dest.id as destination__id, dest.moniker as destination__moniker, dest.cons_address as destination__cons_address").
		Join("left join validator as source on source.id = src_id").
		Join("left join validator as dest on dest.id = dest_id").
		Scan(ctx, &redelegations)

	return
}
