package postgres

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
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
		Limit(1).
		Scan(ctx)
	return
}

type typeCount struct {
	Type  storageTypes.MsgType `bun:"type"`
	Count int64                `bun:"count"`
}

// ByHeightWithStats -
func (b *Blocks) ByHeightWithStats(ctx context.Context, height uint64) (block storage.Block, err error) {

	err = b.DB().NewSelect().Model(&block).
		Where("block.height = ?", height).
		Relation("Stats").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return
	}

	var msgsStats []typeCount
	err = b.DB().NewSelect().Model((*storage.Message)(nil)).
		ColumnExpr("message.type, count(*)").
		Where("message.height = ?", height).
		Group("message.type").
		Scan(ctx, &msgsStats)

	if err != nil {
		return
	}

	block.Stats.MessagesCounts = make(map[storageTypes.MsgType]int64)
	for _, stat := range msgsStats {
		block.Stats.MessagesCounts[stat.Type] = stat.Count
	}

	return
}

// ByIdWithRelations -
func (b *Blocks) ByIdWithRelations(ctx context.Context, id uint64) (block storage.Block, err error) {

	err = b.DB().NewSelect().Model(&block).
		Where("block.id = ?", id).
		Relation("Stats").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return
	}

	var msgsStats []typeCount
	err = b.DB().NewSelect().Model((*storage.Message)(nil)).
		ColumnExpr("message.type, count(*)").
		Where("message.height = ?", block.Height).
		Group("message.type").
		Scan(ctx, &msgsStats)

	if err != nil {
		return
	}

	block.Stats.MessagesCounts = make(map[storageTypes.MsgType]int64)
	for _, stat := range msgsStats {
		block.Stats.MessagesCounts[stat.Type] = stat.Count
	}

	return
}

// Last -
func (b *Blocks) Last(ctx context.Context) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).Order("id desc").Limit(1).Scan(ctx)
	return
}

// ByHash -
func (b *Blocks) ByHash(ctx context.Context, hash []byte) (block storage.Block, err error) {
	err = b.DB().NewSelect().
		Model(&block).
		Where("hash = ?", hash).
		Relation("Stats").
		Limit(1).
		Scan(ctx)
	return
}

type listTypeCount struct {
	Height types.Level          `bun:"height"`
	Type   storageTypes.MsgType `bun:"type"`
	Count  int64                `bun:"count"`
}

// ListWithStats -
func (b *Blocks) ListWithStats(ctx context.Context, limit, offset uint64, order sdk.SortOrder) (blocks []*storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&blocks)
	subQuery = postgres.Pagination(subQuery, limit, offset, order)

	query := b.DB().NewSelect().
		ColumnExpr("block.*").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count, stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee").
		TableExpr("(?) as block", subQuery).
		Join("LEFT JOIN block_stats as stats").
		JoinOn("stats.height = block.height")
	query = sortScope(query, "block.id", order)
	err = query.Scan(ctx, &blocks)

	if err != nil {
		return
	}

	heights := make([]uint64, len(blocks))
	blocksHeighMap := make(map[types.Level]*storage.Block)
	for i, b := range blocks {
		heights[i] = uint64(b.Height)
		blocksHeighMap[b.Height] = b
	}

	var listTypeCounts []listTypeCount
	queryMsgsCounts := b.DB().NewSelect().Model((*storage.Message)(nil)).
		ColumnExpr("message.height, message.type, count(*)").
		Where("message.height IN (?)", bun.In(heights)).
		Group("message.type").
		Group("message.height")

	queryMsgsCounts = sortScope(queryMsgsCounts, "message.height", order)
	err = queryMsgsCounts.Scan(ctx, &listTypeCounts)

	if err != nil {
		return
	}

	for _, stat := range listTypeCounts {
		if blocksHeighMap[stat.Height].Stats.MessagesCounts == nil {
			blocksHeighMap[stat.Height].Stats.MessagesCounts = make(map[storageTypes.MsgType]int64)
		}

		blocksHeighMap[stat.Height].Stats.MessagesCounts[stat.Type] = stat.Count
	}

	return
}
