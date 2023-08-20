package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Blocks -
type Blocks struct {
	*postgres.Table[*storage.Block]
}

// NewBlocks -
func NewBlocks(db *database.Bun) *Blocks {
	return &Blocks{
		Table: postgres.NewTable[*storage.Block](db),
	}
}

// ByHeight -
func (b *Blocks) ByHeight(ctx context.Context, height uint64) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).Where("id = ?", height).Limit(1).Scan(ctx)
	return
}

// Last -
func (b *Blocks) Last(ctx context.Context) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).Order("id desc").Limit(1).Scan(ctx)
	return
}
