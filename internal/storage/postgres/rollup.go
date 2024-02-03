// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

// Rollup -
type Rollup struct {
	*postgres.Table[*storage.Rollup]
}

// NewRollup -
func NewRollup(db *database.Bun) *Rollup {
	return &Rollup{
		Table: postgres.NewTable[*storage.Rollup](db),
	}
}

func (r *Rollup) Leaderboard(ctx context.Context, sortField string, sort sdk.SortOrder, limit, offset int) (rollups []storage.RollupWithStats, err error) {
	switch sortField {
	case timeColumn:
		sortField = "last_time"
	case sizeColumn, blobsCountColumn:
	case "":
		sortField = sizeColumn
	default:
		return nil, errors.Errorf("unknown sort field: %s", sortField)
	}

	timeAggQuery := r.DB().NewSelect().Table("rollup_stats_by_month").
		ColumnExpr("sum(size) as size, sum(blobs_count) as blobs_count, max(last_time) as last_time, namespace_id, signer_id").
		Group("namespace_id", "signer_id")

	leaderboardQuery := r.DB().NewSelect().TableExpr("(?) as agg", timeAggQuery).
		ColumnExpr("sum(size) as size, sum(blobs_count) as blobs_count, max(last_time) as last_time, rollup_id").
		Join("inner join rollup_provider as rp on rp.address_id = agg.signer_id AND (rp.namespace_id = agg.namespace_id OR rp.namespace_id = 0)").
		Group("rollup_id")

	leaderboardQuery = sortScope(leaderboardQuery, sortField, sort)
	if offset > 0 {
		leaderboardQuery = leaderboardQuery.Offset(offset)
	}
	leaderboardQuery = limitScope(leaderboardQuery, limit)

	query := r.DB().NewSelect().Table("leaderboard").With("leaderboard", leaderboardQuery).
		ColumnExpr("size, blobs_count, last_time, rollup.*").
		Join("inner join rollup on rollup.id = leaderboard.rollup_id")

	query = sortScope(query, sortField, sort)
	err = query.Scan(ctx, &rollups)
	return
}

func (r *Rollup) Namespaces(ctx context.Context, rollupId uint64, limit, offset int) (namespaceIds []uint64, err error) {
	query := r.DB().NewSelect().TableExpr("rollup_stats_by_hour as r").
		ColumnExpr("distinct r.namespace_id").
		Join("inner join rollup_provider as rp on rp.address_id = r.signer_id AND (rp.namespace_id = r.namespace_id OR rp.namespace_id = 0)").
		Where("rollup_id = ?", rollupId)
	if offset > 0 {
		query = query.Offset(offset)
	}
	query = limitScope(query, limit)
	err = query.Scan(ctx, &namespaceIds)
	return
}

func (r *Rollup) Providers(ctx context.Context, rollupId uint64) (providers []storage.RollupProvider, err error) {
	err = r.DB().NewSelect().Model(&providers).
		Where("rollup_id = ?", rollupId).
		Scan(ctx)
	return
}

func (r *Rollup) Series(ctx context.Context, rollupId uint64, timeframe, column string, req storage.SeriesRequest) (items []storage.HistogramItem, err error) {
	providers, err := r.Providers(ctx, rollupId)
	if err != nil {
		return nil, err
	}

	if len(providers) == 0 {
		return nil, nil
	}

	query := r.DB().NewSelect().Order("time desc").Limit(100)

	switch timeframe {
	case "hour":
		query = query.Table("rollup_stats_by_hour")
	case "day":
		query = query.Table("rollup_stats_by_day")
	case "month":
		query = query.Table("rollup_stats_by_month")
	default:
		return nil, errors.Errorf("invalid timeframe: %s", timeframe)
	}

	switch column {
	case "blobs_count":
		query = query.ColumnExpr("blobs_count as value, time as bucket")
	case "size":
		query = query.ColumnExpr("size as value, time as bucket")
	default:
		return nil, errors.Errorf("invalid column: %s", column)
	}

	if !req.From.IsZero() {
		query = query.Where("time >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("time < ?", req.To)
	}

	if len(providers) > 0 {
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			for i := range providers {
				q.WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
					sq = sq.Where("signer_id = ?", providers[i].AddressId)
					if providers[i].NamespaceId > 0 {
						sq = sq.Where("namespace_id = ?", providers[i].NamespaceId)
					}
					return sq
				})
			}

			return q
		})
	}

	err = query.Scan(ctx, &items)

	return
}

func (r *Rollup) Count(ctx context.Context) (int64, error) {
	count, err := r.DB().NewSelect().Model((*storage.Rollup)(nil)).Count(ctx)
	return int64(count), err
}

func (r *Rollup) Stats(ctx context.Context, rollupId uint64) (stats storage.RollupStats, err error) {
	providers, err := r.Providers(ctx, rollupId)
	if err != nil {
		return
	}

	if len(providers) == 0 {
		return
	}

	query := r.DB().NewSelect().Table("rollup_stats_by_month").
		ColumnExpr("sum(blobs_count) as blobs_count, sum(size) as size, max(last_time) as last_time")

	for i := range providers {
		query.WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			if providers[i].NamespaceId > 0 {
				return sq.
					Where("namespace_id = ?", providers[i].NamespaceId).
					Where("signer_id = ?", providers[i].AddressId)
			}
			return sq.Where("signer_id = ?", providers[i].AddressId)
		})
	}

	err = query.Scan(ctx, &stats)
	return
}
