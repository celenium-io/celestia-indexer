// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
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

	if err != nil {
		return
	}

	var msgsStats []typeCount
	err = b.db.DB().NewSelect().Model((*storage.Message)(nil)).
		ColumnExpr("message.type, count(*)").
		Where("message.height = ?", height).
		Group("message.type").
		Scan(ctx, &msgsStats)

	if err != nil {
		return
	}

	stats.MessagesCounts = make(map[storageTypes.MsgType]int64)
	for _, stat := range msgsStats {
		stats.MessagesCounts[stat.Type] = stat.Count
	}

	return
}
