package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
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
	err = b.DB().NewSelect().Model(&block).
		Where("block.height = ?", height).
		Relation("Stats").
		Limit(1).
		Scan(ctx)
	return
}

// Last -
func (b *Blocks) Last(ctx context.Context) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).Order("id desc").Limit(1).Scan(ctx)
	return
}

// ByHeight -
func (b *Blocks) ByHash(ctx context.Context, hash []byte) (block storage.Block, err error) {
	err = b.DB().NewSelect().
		Model(&block).
		Where("hash = ?", hash).
		Relation("Stats").
		Limit(1).
		Scan(ctx)
	return
}

// ListWithStats -
func (b *Blocks) ListWithStats(ctx context.Context, stats bool, limit, offset uint64, order sdk.SortOrder) (blocks []storage.Block, err error) {
	query := b.DB().NewSelect().Model(&blocks)
	query = postgres.Pagination(query, limit, offset, order)

	if stats {
		query = query.Relation("Stats")
	}

	err = query.Scan(ctx)
	return
}
