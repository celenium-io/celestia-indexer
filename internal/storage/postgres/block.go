// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
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

// ByHeightWithStats -
func (b *Blocks) ByHeightWithStats(ctx context.Context, height types.Level) (block storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&block).
		Where("block.height = ?", height).
		Limit(1)

	err = b.DB().NewSelect().
		ColumnExpr("block.id, block.height, block.time, block.version_block, block.version_app, block.message_types, block.hash, block.parent_hash, block.last_commit_hash, block.data_hash, block.validators_hash, block.next_validators_hash, block.consensus_hash, block.app_hash, block.last_results_hash, block.evidence_hash, block.proposer_id").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count, stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.gas_limit AS stats__gas_limit, stats.gas_used AS stats__gas_used, stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee, stats.bytes_in_block AS stats__bytes_in_block, stats.blobs_count AS stats__blobs_count, stats.rewards AS stats__rewards, stats.commissions AS stats__commissions, stats.square_size AS stats__square_size").
		ColumnExpr("proposer.id AS proposer__id, proposer.cons_address AS proposer__cons_address, proposer.moniker AS proposer__moniker").
		With("q", subQuery).
		TableExpr("q as block").
		Join("LEFT JOIN block_stats AS stats ON (stats.height = block.height) AND (stats.time = block.time)").
		Join("LEFT JOIN validator AS proposer ON (proposer.id = block.proposer_id)").
		Scan(ctx, &block)

	return
}

// ByIdWithRelations -
func (b *Blocks) ByIdWithRelations(ctx context.Context, id uint64) (block storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&block).
		Where("block.id = ?", id).
		Limit(1)

	err = b.DB().NewSelect().
		ColumnExpr("block.id, block.height, block.time, block.version_block, block.version_app, block.message_types, block.hash, block.parent_hash, block.last_commit_hash, block.data_hash, block.validators_hash, block.next_validators_hash, block.consensus_hash, block.app_hash, block.last_results_hash, block.evidence_hash, block.proposer_id").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count, stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.gas_limit AS stats__gas_limit, stats.gas_used AS stats__gas_used, stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee, stats.bytes_in_block AS stats__bytes_in_block, stats.blobs_count AS stats__blobs_count, stats.rewards AS stats__rewards, stats.commissions AS stats__commissions, stats.square_size AS stats__square_size").
		ColumnExpr("proposer.id AS proposer__id, proposer.cons_address AS proposer__cons_address, proposer.moniker AS proposer__moniker").
		With("q", subQuery).
		TableExpr("q as block").
		Join("LEFT JOIN block_stats AS stats ON (stats.height = block.height) AND (stats.time = block.time)").
		Join("LEFT JOIN validator AS proposer ON (proposer.id = block.proposer_id)").
		Scan(ctx, &block)

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
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count, stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.gas_limit AS stats__gas_limit, stats.gas_used AS stats__gas_used, stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee, stats.bytes_in_block AS stats__bytes_in_block, stats.blobs_count AS stats__blobs_count, stats.rewards AS stats__rewards, stats.commissions AS stats__commissions, stats.square_size AS stats__square_size").
		ColumnExpr("proposer.id AS proposer__id, proposer.cons_address AS proposer__cons_address, proposer.moniker AS proposer__moniker").
		With("q", subQuery).
		TableExpr("q as block").
		Join("LEFT JOIN block_stats AS stats ON (stats.height = block.height) AND (stats.time = block.time)").
		Join("LEFT JOIN validator AS proposer ON (proposer.id = block.proposer_id)").
		Scan(ctx, &block)

	return
}

// ListWithStats -
func (b *Blocks) ListWithStats(ctx context.Context, limit, offset uint64, order sdk.SortOrder) (blocks []*storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&blocks)
	subQuery = postgres.Pagination(subQuery, limit, offset, order)

	query := b.DB().NewSelect().
		ColumnExpr("block.*").
		ColumnExpr("v.id AS proposer__id, v.cons_address as proposer__cons_address, v.moniker as proposer__moniker").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count, stats.blobs_count as stats__blobs_count").
		ColumnExpr("stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.bytes_in_block AS stats__bytes_in_block, stats.rewards AS stats__rewards, stats.commissions AS stats__commissions").
		ColumnExpr("stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee, stats.gas_used AS stats__gas_used, stats.gas_limit AS stats__gas_limit, stats.square_size AS stats__square_size").
		TableExpr("(?) as block", subQuery).
		Join("LEFT JOIN block_stats as stats ON stats.height = block.height").
		Join("LEFT JOIN validator as v ON v.id = block.proposer_id")
	query = sortScope(query, "block.id", order)

	err = query.Scan(ctx, &blocks)
	return
}

func (b *Blocks) ByProposer(ctx context.Context, proposerId uint64, limit, offset int) (blocks []storage.Block, err error) {
	blocksQuery := b.DB().NewSelect().Model(&blocks).
		Where("proposer_id = ?", proposerId).
		Order("time desc")

	blocksQuery = limitScope(blocksQuery, limit)
	if offset > 0 {
		blocksQuery = blocksQuery.Offset(offset)
	}

	err = b.DB().NewSelect().
		ColumnExpr("block.*").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count, stats.events_count AS stats__events_count, stats.blobs_count as stats__blobs_count").
		ColumnExpr("stats.blobs_size AS stats__blobs_size, stats.block_time AS stats__block_time, stats.bytes_in_block AS stats__bytes_in_block, stats.rewards AS stats__rewards, stats.commissions AS stats__commissions").
		ColumnExpr("stats.supply_change AS stats__supply_change, stats.inflation_rate AS stats__inflation_rate, stats.fee AS stats__fee, stats.gas_used AS stats__gas_used, stats.gas_limit AS stats__gas_limit, stats.square_size AS stats__square_size").
		TableExpr("(?) as block", blocksQuery).
		Join("LEFT JOIN block_stats as stats ON stats.height = block.height").
		Order("time desc").
		Scan(ctx, &blocks)
	return
}

func (b *Blocks) Time(ctx context.Context, height pkgTypes.Level) (response time.Time, err error) {
	err = b.DB().NewSelect().Model((*storage.Block)(nil)).
		Column("time").
		Where("height = ?", height).
		Scan(ctx, &response)
	return
}
