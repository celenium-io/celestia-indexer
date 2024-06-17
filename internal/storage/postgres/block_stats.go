// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// BlockStats -
type BlockStats struct {
	db *database.Bun
}

// NewBlockStats -
func NewBlockStats(db *database.Bun) *BlockStats {
	return &BlockStats{
		db: db,
	}
}

// ByHeight -
func (b *BlockStats) ByHeight(ctx context.Context, height pkgTypes.Level) (stats storage.BlockStats, err error) {
	err = b.db.DB().NewSelect().Model(&stats).
		Where("height = ?", height).
		Limit(1).
		Scan(ctx)

	return
}

func (b *BlockStats) LastFrom(ctx context.Context, head pkgTypes.Level, limit int) (stats []storage.BlockStats, err error) {
	err = b.db.DB().NewSelect().Model(&stats).
		Where("height <= ?", head).
		Limit(limit).
		Order("id desc").
		Scan(ctx)
	return
}
