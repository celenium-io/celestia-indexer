// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
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
func (b *Blocks) ByHeight(ctx context.Context, height types.Level) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).
		Where("block.height = ?", height).
		Relation("Proposer", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("id", "cons_address", "moniker")
		}).
		Limit(1).
		Scan(ctx)
	return
}

type typeCount struct {
	Type  storageTypes.MsgType `bun:"type"`
	Count int64                `bun:"count"`
}

// ByHeightWithStats -
func (b *Blocks) ByHeightWithStats(ctx context.Context, height types.Level) (block storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&block).
		Where("block.height = ?", height).
		Limit(1)

	err = b.DB().NewSelect().
		ColumnExpr("block.id, block.height, block.time, block.version_block, block.version_app, block.message_types, block.hash, block.parent_hash, block.last_commit_hash, block.data_hash, block.validators_hash, block.next_validators_hash, block.consensus_hash, block.app_hash, block.last_results_hash, block.evidence_hash, block.proposer_id").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count, stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.gas_limit AS stats__gas_limit, stats.gas_used AS stats__gas_used, stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee, stats.bytes_in_block AS stats__bytes_in_block").
		ColumnExpr("proposer.id AS proposer__id, proposer.cons_address AS proposer__cons_address, proposer.moniker AS proposer__moniker").
		With("q", subQuery).
		TableExpr("q as block").
		Join("LEFT JOIN block_stats AS stats ON (stats.id = block.id) AND (stats.time = block.time)").
		Join("LEFT JOIN validator AS proposer ON (proposer.id = block.proposer_id)").
		Scan(ctx, &block)

	if err != nil {
		return
	}

	block.Stats.MessagesCounts = make(map[storageTypes.MsgType]int64)

	if block.Stats.TxCount > 0 {
		var msgsStats []typeCount
		err = b.DB().NewSelect().Model((*storage.Message)(nil)).
			ColumnExpr("message.type, count(*)").
			Where("message.height = ?", height).
			Where("message.time = ?", block.Time).
			Group("message.type").
			Scan(ctx, &msgsStats)

		if err != nil {
			return
		}

		for _, stat := range msgsStats {
			block.Stats.MessagesCounts[stat.Type] = stat.Count
		}
	}

	return
}

// ByIdWithRelations -
func (b *Blocks) ByIdWithRelations(ctx context.Context, id uint64) (block storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&block).
		Where("block.id = ?", id).
		Limit(1)

	err = b.DB().NewSelect().
		ColumnExpr("block.id, block.height, block.time, block.version_block, block.version_app, block.message_types, block.hash, block.parent_hash, block.last_commit_hash, block.data_hash, block.validators_hash, block.next_validators_hash, block.consensus_hash, block.app_hash, block.last_results_hash, block.evidence_hash, block.proposer_id").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count, stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.gas_limit AS stats__gas_limit, stats.gas_used AS stats__gas_used, stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee, stats.bytes_in_block AS stats__bytes_in_block").
		ColumnExpr("proposer.id AS proposer__id, proposer.cons_address AS proposer__cons_address, proposer.moniker AS proposer__moniker").
		With("q", subQuery).
		TableExpr("q as block").
		Join("LEFT JOIN block_stats AS stats ON (stats.id = block.id) AND (stats.time = block.time)").
		Join("LEFT JOIN validator AS proposer ON (proposer.id = block.proposer_id)").
		Scan(ctx, &block)

	if err != nil {
		return
	}

	block.Stats.MessagesCounts = make(map[storageTypes.MsgType]int64)

	if block.Stats.TxCount > 0 {
		var msgsStats []typeCount
		err = b.DB().NewSelect().Model((*storage.Message)(nil)).
			ColumnExpr("message.type, count(*)").
			Where("message.height = ?", block.Height).
			Where("message.time = ?", block.Time).
			Group("message.type").
			Scan(ctx, &msgsStats)

		if err != nil {
			return
		}

		for _, stat := range msgsStats {
			block.Stats.MessagesCounts[stat.Type] = stat.Count
		}
	}

	return
}

// Last -
func (b *Blocks) Last(ctx context.Context) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).
		Relation("Proposer", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("id", "cons_address", "moniker")
		}).
		Order("id desc").
		Limit(1).
		Scan(ctx)
	return
}

// ByHash -
func (b *Blocks) ByHash(ctx context.Context, hash []byte) (block storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&block).
		Where("hash = ?", hash).
		Limit(1)

	err = b.DB().NewSelect().
		ColumnExpr("block.id, block.height, block.time, block.version_block, block.version_app, block.message_types, block.hash, block.parent_hash, block.last_commit_hash, block.data_hash, block.validators_hash, block.next_validators_hash, block.consensus_hash, block.app_hash, block.last_results_hash, block.evidence_hash, block.proposer_id").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count, stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.gas_limit AS stats__gas_limit, stats.gas_used AS stats__gas_used, stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee, stats.bytes_in_block AS stats__bytes_in_block").
		ColumnExpr("proposer.id AS proposer__id, proposer.cons_address AS proposer__cons_address, proposer.moniker AS proposer__moniker").
		With("q", subQuery).
		TableExpr("q as block").
		Join("LEFT JOIN block_stats AS stats ON (stats.id = block.id) AND (stats.time = block.time)").
		Join("LEFT JOIN validator AS proposer ON (proposer.id = block.proposer_id)").
		Scan(ctx, &block)

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
		ColumnExpr("v.id AS proposer__id, v.cons_address as proposer__cons_address, v.moniker as proposer__moniker").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count").
		ColumnExpr("stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.bytes_in_block AS stats__bytes_in_block").
		ColumnExpr("stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee, stats.gas_used AS stats__gas_used, stats.gas_limit AS stats__gas_limit").
		TableExpr("(?) as block", subQuery).
		Join("LEFT JOIN block_stats as stats ON stats.height = block.height").
		Join("LEFT JOIN validator as v ON v.id = block.proposer_id")
	query = sortScope(query, "block.id", order)
	err = query.Scan(ctx, &blocks)

	if err != nil {
		return
	}

	heights := make([]types.Level, len(blocks))
	blocksHeightMap := make(map[types.Level]*storage.Block)
	for i, b := range blocks {
		heights[i] = b.Height
		blocksHeightMap[b.Height] = b
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
		if blocksHeightMap[stat.Height].Stats.MessagesCounts == nil {
			blocksHeightMap[stat.Height].Stats.MessagesCounts = make(map[storageTypes.MsgType]int64)
		}

		blocksHeightMap[stat.Height].Stats.MessagesCounts[stat.Type] = stat.Count
	}

	return
}

func (b *Blocks) ByProposer(ctx context.Context, proposerId uint64, limit, offset int) (blocks []storage.Block, err error) {
	query := b.DB().NewSelect().Model(&blocks).
		Where("proposer_id = ?", proposerId).
		Relation("Stats").
		Order("id desc")

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx)
	return
}
