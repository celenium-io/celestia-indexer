// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
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

func (r *Rollup) Leaderboard(ctx context.Context, fltrs storage.LeaderboardFilters) (rollups []storage.RollupWithStats, err error) {
	switch fltrs.SortField {
	case timeColumn:
		fltrs.SortField = "last_time"
	case sizeColumn, blobsCountColumn, feeColumn:
	case "":
		fltrs.SortField = sizeColumn
	default:
		return nil, errors.Errorf("unknown sort field: %s", fltrs.SortField)
	}
	query := r.DB().NewSelect().
		Table(storage.ViewLeaderboard).
		ColumnExpr("leaderboard.*").
		ColumnExpr("da_change.da_pct as da_pct").
		Offset(fltrs.Offset).
		Join("left join da_change on da_change.rollup_id = leaderboard.id")

	if len(fltrs.Category) > 0 {
		query = query.Where("category IN (?)", bun.In(fltrs.Category))
	}

	if len(fltrs.Tags) > 0 {
		query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			for i := range fltrs.Tags {
				q.WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
					return sq.Where("? = ANY(tags)", fltrs.Tags[i])
				})
			}
			return q
		})
	}

	if len(fltrs.Stack) > 0 {
		query = query.Where("stack IN (?)", bun.In(fltrs.Stack))
	}

	if len(fltrs.Provider) > 0 {
		query = query.Where("provider IN (?)", bun.In(fltrs.Provider))
	}

	if len(fltrs.Type) > 0 {
		query = query.Where("type IN (?)", bun.In(fltrs.Type))
	}

	if fltrs.IsActive != nil {
		query = query.Where("is_active = ?", *fltrs.IsActive)
	}

	query = sortScope(query, fltrs.SortField, fltrs.Sort)
	query = limitScope(query, fltrs.Limit)
	err = query.Scan(ctx, &rollups)
	return
}

func (r *Rollup) LeaderboardDay(ctx context.Context, fltrs storage.LeaderboardFilters) (rollups []storage.RollupWithDayStats, err error) {
	switch fltrs.SortField {
	case "avg_size", blobsCountColumn, "total_size", "total_fee", "throughput", "namespace_count", "pfb_count", "mb_price":
	case "":
		fltrs.SortField = "throughput"
	default:
		return nil, errors.Errorf("unknown sort field: %s", fltrs.SortField)
	}

	query := r.DB().NewSelect().
		Table(storage.ViewLeaderboardDay).
		Column("avg_size", blobsCountColumn, "total_size", "total_fee", "throughput", "namespace_count", "pfb_count", "mb_price").
		ColumnExpr("rollup.*").
		Offset(fltrs.Offset).
		Join("left join rollup on rollup.id = rollup_id AND rollup.verified = true")

	if len(fltrs.Category) > 0 {
		query = query.Where("category IN (?)", bun.In(fltrs.Category))
	}

	if len(fltrs.Tags) > 0 {
		query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			for i := range fltrs.Tags {
				q.WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
					return sq.Where("? = ANY(tags)", fltrs.Tags[i])
				})
			}
			return q
		})
	}

	if len(fltrs.Stack) > 0 {
		query = query.Where("stack IN (?)", bun.In(fltrs.Stack))
	}

	if len(fltrs.Provider) > 0 {
		query = query.Where("provider IN (?)", bun.In(fltrs.Provider))
	}

	if len(fltrs.Type) > 0 {
		query = query.Where("type IN (?)", bun.In(fltrs.Type))
	}

	query = sortScope(query, fltrs.SortField, fltrs.Sort)
	query = limitScope(query, fltrs.Limit)
	err = query.Scan(ctx, &rollups)
	return
}

func (r *Rollup) Namespaces(ctx context.Context, rollupId uint64, limit, offset int) (namespaceIds []uint64, err error) {
	query := r.DB().NewSelect().
		TableExpr("rollup_stats_by_month as r").
		ColumnExpr("distinct r.namespace_id").
		Join("inner join rollup_provider as rp on (rp.address_id = r.signer_id OR rp.address_id = 0) AND (rp.namespace_id = r.namespace_id OR rp.namespace_id = 0)").
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

func (r *Rollup) RollupsByNamespace(ctx context.Context, namespaceId uint64, limit, offset int) (rollups []storage.Rollup, err error) {
	subQuery := r.DB().NewSelect().
		Model((*storage.RollupProvider)(nil)).
		Column("rollup_id").
		Where("namespace_id = ?", namespaceId).
		Group("rollup_id").
		Offset(offset)

	subQuery = limitScope(subQuery, limit)

	err = r.DB().NewSelect().
		With("rollups", subQuery).
		Table("rollups").
		ColumnExpr("rollup.*").
		Join("left join rollup on rollup.id = rollups.rollup_id and rollup.verified = true").
		Scan(ctx, &rollups)
	return
}

func (r *Rollup) Series(ctx context.Context, rollupId uint64, timeframe storage.Timeframe, column string, req storage.SeriesRequest) (items []storage.HistogramItem, err error) {
	providers, err := r.Providers(ctx, rollupId)
	if err != nil {
		return nil, err
	}

	if len(providers) == 0 {
		return nil, nil
	}

	query := r.DB().NewSelect().Order("time desc").Group("time")

	switch timeframe {
	case storage.TimeframeHour:
		query = query.Table("rollup_stats_by_hour")
	case storage.TimeframeDay:
		query = query.Table("rollup_stats_by_day")
	case storage.TimeframeMonth:
		query = query.Table("rollup_stats_by_month")
	default:
		return nil, errors.Errorf("invalid timeframe: %s", timeframe)
	}

	switch column {
	case "blobs_count":
		query = query.ColumnExpr("sum(blobs_count) as value, time as bucket")
	case "size":
		query = query.ColumnExpr("sum(size) as value, time as bucket")
	case "size_per_blob":
		query = query.ColumnExpr("(sum(size) / sum(blobs_count)) as value, time as bucket")
	case "fee":
		query = query.ColumnExpr("sum(fee) as value, time as bucket")
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
	count, err := r.DB().NewSelect().Model((*storage.Rollup)(nil)).Where("verified = TRUE").Count(ctx)
	return int64(count), err
}

func (r *Rollup) Stats(ctx context.Context, rollupId uint64) (stats storage.RollupStats, err error) {
	err = r.DB().NewSelect().Table(storage.ViewLeaderboard).
		Column("blobs_count", "size", "last_time", "first_time", "fee", "size_pct", "fee_pct", "blobs_count_pct").
		Where("id = ?", rollupId).Scan(ctx, &stats)
	return
}

func (r *Rollup) BySlug(ctx context.Context, slug string) (rollup storage.RollupWithStats, err error) {
	err = r.DB().NewSelect().
		Table(storage.ViewLeaderboard).
		ColumnExpr("leaderboard.*").
		ColumnExpr("da_change.da_pct as da_pct").
		Where("slug = ?", slug).
		Limit(1).
		Join("left join da_change on da_change.rollup_id = leaderboard.id").
		Scan(ctx, &rollup)
	return
}

func (r *Rollup) ById(ctx context.Context, rollupId uint64) (rollup storage.RollupWithStats, err error) {
	err = r.DB().NewSelect().
		Table(storage.ViewLeaderboard).
		Where("id = ?", rollupId).
		ColumnExpr("leaderboard.*").
		ColumnExpr("da_change.da_pct as da_pct").
		Limit(1).
		Join("left join da_change on da_change.rollup_id = leaderboard.id").
		Scan(ctx, &rollup)
	return
}

func (r *Rollup) Distribution(ctx context.Context, rollupId uint64, series string, groupBy storage.Timeframe) (items []storage.DistributionItem, err error) {
	providers, err := r.Providers(ctx, rollupId)
	if err != nil {
		return
	}

	if len(providers) == 0 {
		return
	}

	cte := r.DB().NewSelect()

	for i := range providers {
		cte.WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			if providers[i].NamespaceId > 0 {
				return sq.
					Where("namespace_id = ?", providers[i].NamespaceId).
					Where("signer_id = ?", providers[i].AddressId)
			}
			return sq.Where("signer_id = ?", providers[i].AddressId)
		})
	}

	switch groupBy {
	case storage.TimeframeDay:
		cte = cte.Table("rollup_stats_by_day").
			ColumnExpr("extract(isodow from time) as name").
			Where("time >= ?", time.Now().AddDate(0, -3, 0).UTC())
	case storage.TimeframeHour:
		cte = cte.Table("rollup_stats_by_hour").
			ColumnExpr("extract(hour from time) as name").
			Where("time >= ?", time.Now().AddDate(0, -1, 0).UTC())
	default:
		err = errors.Errorf("invalid distribution rollup groupBy: %s", groupBy)
		return
	}

	switch series {
	case "size":
		cte = cte.ColumnExpr("size as value")
	case "blobs_count":
		cte = cte.ColumnExpr("blobs_count as value")
	case "size_per_blob":
		cte = cte.ColumnExpr("(size / blobs_count) as value")
	case "fee_per_blob":
		cte = cte.ColumnExpr("(fee / blobs_count) as value")
	default:
		err = errors.Errorf("invalid distribution rollup series: %s", groupBy)
		return
	}

	err = r.DB().NewSelect().
		TableExpr("(?) as cte", cte).
		ColumnExpr("name, avg(value) as value").
		Group("name").
		Order("name asc").
		Scan(ctx, &items)
	return
}

func (r *Rollup) AllSeries(ctx context.Context, timeframe storage.Timeframe) (items []storage.RollupHistogramItem, err error) {
	subQuery := r.DB().NewSelect().
		ColumnExpr("rp.rollup_id, sum(size) as size, sum(blobs_count) as blobs_count, sum(fee) as fee, time").
		Join("inner join rollup_provider rp on (rp.namespace_id = 0 or rp.namespace_id = stats.namespace_id) and (rp.address_id = signer_id OR rp.address_id = 0)").
		Group("rollup_id", "time")

	switch timeframe {
	case storage.TimeframeHour:
		subQuery = subQuery.TableExpr("? as stats", bun.Safe(storage.ViewRollupStatsByHour)).Where("time > now() - '24 hours'::interval")
	case storage.TimeframeDay:
		subQuery = subQuery.TableExpr("? as stats", bun.Safe(storage.ViewRollupStatsByDay)).Where("time > now() - '30 days'::interval")
	case storage.TimeframeMonth:
		subQuery = subQuery.TableExpr("? as stats", bun.Safe(storage.ViewRollupStatsByMonth)).Where("time > now() - '1 year'::interval")
	}

	err = r.DB().NewSelect().
		TableExpr("(?) as series", subQuery).
		ColumnExpr("series.time as time, series.size as size, series.blobs_count as blobs_count, series.fee as fee, rollup.name as name, rollup.logo as logo").
		Join("left join rollup on rollup.id = series.rollup_id").
		Where("rollup.verified = true").
		OrderExpr("time desc").
		Scan(ctx, &items)

	return
}

func (r *Rollup) RollupStatsGrouping(ctx context.Context, fltrs storage.RollupGroupStatsFilters) (results []storage.RollupGroupedStats, err error) {
	query := r.DB().NewSelect().Table(storage.ViewLeaderboard)

	switch fltrs.Func {
	case "sum":
		query = query.
			ColumnExpr("sum(fee) as fee").
			ColumnExpr("sum(size) as size").
			ColumnExpr("sum(blobs_count) as blobs_count")
	case "avg":
		query = query.
			ColumnExpr("avg(fee) as fee").
			ColumnExpr("avg(size) as size").
			ColumnExpr("avg(blobs_count) as blobs_count")
	default:
		return nil, errors.Errorf("unknown func field: %s", fltrs.Column)
	}

	switch fltrs.Column {
	case "stack", "type", "category", "vm", "provider":
		query = query.ColumnExpr(fltrs.Column + " as group").Group(fltrs.Column)
	default:
		return nil, errors.Errorf("unknown column field: %s", fltrs.Column)
	}

	err = query.Scan(ctx, &results)
	return
}

func (r *Rollup) Tags(ctx context.Context) (arr []string, err error) {
	err = r.DB().NewSelect().
		Model((*storage.Rollup)(nil)).
		Distinct().
		ColumnExpr("unnest(tags)").
		Where("verified = true").
		Scan(ctx, &arr)
	return
}

func (r *Rollup) Unverified(ctx context.Context) (rollups []storage.Rollup, err error) {
	err = r.DB().NewSelect().
		Model(&rollups).
		Where("verified = false").
		Scan(ctx)
	return
}
