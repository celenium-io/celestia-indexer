// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Undelegation -
type Undelegation struct {
	*postgres.Table[*storage.Undelegation]
}

// NewUndelegation -
func NewUndelegation(db *database.Bun) *Undelegation {
	return &Undelegation{
		Table: postgres.NewTable[*storage.Undelegation](db),
	}
}

func (d *Undelegation) ByAddress(ctx context.Context, addressId uint64, limit, offset int) (undelegations []storage.Undelegation, err error) {
	subQuery := d.DB().NewSelect().Model((*storage.Undelegation)(nil)).
		Where("address_id = ?", addressId).
		Order("amount desc")

	subQuery = limitScope(subQuery, limit)
	if offset > 0 {
		subQuery = subQuery.Offset(offset)
	}

	err = d.DB().NewSelect().
		TableExpr("(?) as undelegation", subQuery).
		ColumnExpr("undelegation.*").
		ColumnExpr("validator.id as validator__id, validator.moniker as validator__moniker, validator.cons_address as validator__cons_address").
		Join("left join validator on validator.id = validator_id").
		Scan(ctx, &undelegations)

	return
}
