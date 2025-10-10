// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

type RollupProvider struct {
	db *database.Bun
}

func NewRollupProvider(db *database.Bun) *RollupProvider {
	return &RollupProvider{
		db: db,
	}
}

func (r *RollupProvider) ByRollupId(ctx context.Context, rollupId uint64) (providers []storage.RollupProvider, err error) {
	err = r.db.DB().NewSelect().
		Model(&providers).
		ColumnExpr("rollup_provider.*").
		ColumnExpr("address.address as address__address").
		ColumnExpr("namespace.namespace_id as namespace__namespace_id, namespace.version as namespace__version").
		Where("rollup_id = ?", rollupId).
		Join("left join address ON address.id = rollup_provider.address_id").
		Join("left join namespace ON namespace.id = rollup_provider.namespace_id").
		Scan(ctx)
	return
}
