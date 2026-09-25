// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

// BlockStats -
type BlockStats struct {
	db bun.IDB
}

// NewBlockStats -
func NewBlockStats(db bun.IDB) storage.IBlockStats {
	return &BlockStats{
		db: db,
	}
}

// ByHeight -
func (b *BlockStats) ByHeight(ctx context.Context, height pkgTypes.Level) (stats storage.BlockStats, err error) {
	err = b.db.NewSelect().Model(&stats).
		Where("height = ?", height).
		Limit(1).
		Scan(ctx)

	return
}

func (b *BlockStats) LastFrom(ctx context.Context, head pkgTypes.Level, limit int) (stats []storage.BlockStats, err error) {
	err = b.db.NewSelect().Model(&stats).
		Where("height <= ?", head).
		Limit(limit).
		Order("id desc").
		Scan(ctx)
	return
}
