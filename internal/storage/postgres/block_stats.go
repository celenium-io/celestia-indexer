package postgres

import (
	"context"

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
func (b *BlockStats) ByHeight(ctx context.Context, height uint64) (stats storage.BlockStats, err error) {
	err = b.db.DB().NewSelect().Model(&stats).
		Where("height = ?", height).
		Limit(1).
		Scan(ctx)
	return
}
