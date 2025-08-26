// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// Upgrade -
type Upgrade struct {
	*database.Bun
}

// NewUpgrade -
func NewUpgrade(db *database.Bun) *Upgrade {
	return &Upgrade{
		Bun: db,
	}
}

func (t *Upgrade) List(ctx context.Context, filters storage.ListUpgradesFilter) (upgrades []storage.Upgrade, err error) {
	query := t.DB().NewSelect().
		Model((*storage.Upgrade)(nil))

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	query = limitScope(query, filters.Limit)
	query = sortScope(query, "height", filters.Sort)

	if filters.SignerId != nil {
		query = query.Where("signer_id = ?", *filters.SignerId)
	}

	if filters.TxId != nil {
		query = query.Where("tx_id = ?", *filters.TxId)
	}

	if filters.Height > 0 {
		query = query.Where("height = ?", filters.Height)
	}

	err = t.DB().NewSelect().
		TableExpr("(?) as upgrade", query).
		ColumnExpr("upgrade.*").
		ColumnExpr("signer.address as signer__address").
		Join("left join address as signer on signer.id = signer_id").
		Scan(ctx, &upgrades)

	return
}
