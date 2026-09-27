// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

// Upgrade -
type Upgrade struct {
	db bun.IDB
}

// NewUpgrade -
func NewUpgrade(db bun.IDB) storage.IUpgrade {
	return &Upgrade{
		db: db,
	}
}

func (t *Upgrade) List(ctx context.Context, filters storage.ListUpgradesFilter) (upgrades []storage.Upgrade, err error) {
	query := t.db.NewSelect().
		Model((*storage.Upgrade)(nil))

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	query = limitScope(query, filters.Limit)
	query = sortScope(query, "version", filters.Sort)

	if filters.SignerId != nil {
		query = query.Where("signer_id = ?", *filters.SignerId)
	}

	if filters.TxId != nil {
		query = query.Where("tx_id = ?", *filters.TxId)
	}

	if filters.Height > 0 {
		query = query.Where("height = ?", filters.Height)
	}

	q := t.db.NewSelect().
		TableExpr("(?) as upgrade", query).
		ColumnExpr("upgrade.*").
		ColumnExpr("signer.address as signer__address").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join address as signer on signer.id = signer_id").
		Join("left join tx on tx_id = tx.id")

	q = sortScope(q, "version", filters.Sort)
	err = q.Scan(ctx, &upgrades)
	return
}

func (t *Upgrade) ByVersion(ctx context.Context, version uint64) (upgrade storage.Upgrade, err error) {
	query := t.db.NewSelect().
		Model((*storage.Upgrade)(nil)).
		Where("version = ?", version).
		Limit(1)

	err = t.db.NewSelect().
		TableExpr("(?) as upgrade", query).
		ColumnExpr("upgrade.*").
		ColumnExpr("signer.address as signer__address").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join address as signer on signer.id = signer_id").
		Join("left join tx on tx_id = tx.id").
		Scan(ctx, &upgrade)
	return
}
