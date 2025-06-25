// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Stats struct {
	db *database.Bun
}

func NewStats(conn *database.Bun) Stats {
	return Stats{conn}
}

func (s Stats) Count(ctx context.Context, req storage.CountRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	query := s.db.DB().NewSelect().Table(req.Table).
		ColumnExpr("COUNT(*)")

	if req.From > 0 {
		query = query.Where("time >= to_timestamp(?)", req.From)
	}
	if req.To > 0 {
		query = query.Where("time < to_timestamp(?)", req.To)
	}

	var count string
	err := query.Scan(ctx, &count)
	return count, err
}

func (s Stats) Summary(ctx context.Context, req storage.SummaryRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	query := s.db.DB().NewSelect().Table(req.Table).
		ColumnExpr(`? (?)`, bun.Safe(req.Function), bun.Safe(req.Column))

	if req.From > 0 {
		query = query.Where("time >= to_timestamp(?)", req.From)
	}
	if req.To > 0 {
		query = query.Where("time < to_timestamp(?)", req.To)
	}

	var value string
	err := query.Scan(ctx, &value)
	return value, err
}

func (s Stats) TPS(ctx context.Context) (response storage.TPS, err error) {
	if err = s.db.DB().NewSelect().Table(storage.ViewBlockStatsByHour).
		ColumnExpr("max(tps) as high, min(tps) as low").
		Where("ts > date_trunc('hour', now()) - '1 week'::interval").
		Where("ts < date_trunc('hour', now())").
		Scan(ctx, &response.High, &response.Low); err != nil {
		return
	}

	if err = s.db.DB().NewSelect().Model((*storage.BlockStats)(nil)).
		ColumnExpr("sum(tx_count)/3600.0").
		Where("time > now() - '1 hour'::interval").
		Scan(ctx, &response.Current); err != nil {
		return
	}
	var prev float64
	if err = s.db.DB().NewSelect().Model((*storage.BlockStats)(nil)).
		ColumnExpr("sum(tx_count)/3600.0").
		Where("time > now() - '2 hour'::interval").
		Where("time <= now() - '1 hour'::interval").
		Scan(ctx, &prev); err != nil {
		return
	}

	switch {
	case prev == 0 && response.Current == 0:
		response.ChangeLastHourPct = 0
	case prev == 0 && response.Current > 0:
		response.ChangeLastHourPct = 1
	default:
		response.ChangeLastHourPct = (response.Current - prev) / prev
	}

	return
}

func (s Stats) Change24hBlockStats(ctx context.Context) (response storage.Change24hBlockStats, err error) {
	first := s.db.DB().NewSelect().
		Table(storage.ViewBlockStatsByHour).
		ColumnExpr(`
			sum(tx_count) as tx_count,
			sum(bytes_in_block) as bytes_in_block,
			sum(blobs_size) as blobs_size,
			sum(fee) as fee`).
		Where("ts > NOW() - '1 day':: interval")
	second := s.db.DB().NewSelect().
		Table(storage.ViewBlockStatsByHour).
		ColumnExpr(`
			sum(tx_count) as tx_count,
			sum(bytes_in_block) as bytes_in_block,
			sum(blobs_size) as blobs_size,
			sum(fee) as fee`).
		Where("ts <= NOW() - '1 day':: interval").
		Where("ts > NOW() - '2 days':: interval")

	err = s.db.DB().NewSelect().
		With("f", first).
		With("s", second).
		TableExpr("f, s").
		ColumnExpr(`
			case when s.tx_count > 0 then (f.tx_count - s.tx_count)/s.tx_count when f.tx_count > 0 then 1 else 0 end as tx_count_24h,
			case when s.bytes_in_block > 0 then (f.bytes_in_block - s.bytes_in_block)/s.bytes_in_block when f.bytes_in_block > 0 then 1 else 0 end as bytes_in_block_24h,
			case when s.blobs_size > 0 then (f.blobs_size - s.blobs_size)/s.blobs_size when f.blobs_size > 0 then 1 else 0 end as blobs_size_24h,
			case when s.fee > 0 then (f.fee - s.fee)/s.fee when f.fee > 0 then 1 else 0 end as fee_24h
		`).
		Scan(ctx, &response)
	return
}

func (s Stats) Series(ctx context.Context, timeframe storage.Timeframe, name string, req storage.SeriesRequest) (response []storage.SeriesItem, err error) {
	var view string
	switch timeframe {
	case storage.TimeframeHour:
		view = storage.ViewBlockStatsByHour
	case storage.TimeframeDay:
		view = storage.ViewBlockStatsByDay
	case storage.TimeframeWeek:
		view = storage.ViewBlockStatsByWeek
	case storage.TimeframeMonth:
		view = storage.ViewBlockStatsByMonth
	case storage.TimeframeYear:
		view = storage.ViewBlockStatsByYear
	default:
		return nil, errors.Errorf("unexpected timeframe %s", timeframe)
	}

	query := s.db.DB().NewSelect().Table(view)

	switch name {
	case storage.SeriesBlobsSize:
		query.ColumnExpr("ts, blobs_size as value")
	case storage.SeriesTPS:
		query.ColumnExpr("ts, tps as value, tps_max as max, tps_min as min")
	case storage.SeriesBPS:
		query.ColumnExpr("ts, bps as value, bps_max as max, bps_min as min")
	case storage.SeriesFee:
		query.ColumnExpr("ts, fee as value")
	case storage.SeriesSupplyChange:
		query.ColumnExpr("ts, supply_change as value")
	case storage.SeriesBlockTime:
		query.ColumnExpr("ts, block_time as value")
	case storage.SeriesTxCount:
		query.ColumnExpr("ts, tx_count as value")
	case storage.SeriesEventsCount:
		query.ColumnExpr("ts, events_count as value")
	case storage.SeriesGasPrice:
		query.ColumnExpr("ts, gas_price as value")
	case storage.SeriesGasEfficiency:
		query.ColumnExpr("ts, gas_efficiency as value")
	case storage.SeriesGasLimit:
		query.ColumnExpr("ts, gas_limit as value")
	case storage.SeriesGasUsed:
		query.ColumnExpr("ts, gas_used as value")
	case storage.SeriesBytesInBlock:
		query.ColumnExpr("ts, bytes_in_block as value")
	case storage.SeriesBlobsCount:
		query.ColumnExpr("ts, blobs_count as value")
	case storage.SeriesCommissions:
		query.ColumnExpr("ts, commissions as value")
	case storage.SeriesRewards:
		query.ColumnExpr("ts, rewards as value")
	default:
		return nil, errors.Errorf("unexpected series name: %s", name)
	}

	if !req.From.IsZero() {
		query = query.Where("ts >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("ts < ?", req.To)
	}

	err = query.Scan(ctx, &response)
	return
}

func (s Stats) NamespaceSeries(ctx context.Context, timeframe storage.Timeframe, name string, nsId uint64, req storage.SeriesRequest) (response []storage.SeriesItem, err error) {
	var view string
	switch timeframe {
	case storage.TimeframeHour:
		view = storage.ViewNamespaceStatsByHour
	case storage.TimeframeDay:
		view = storage.ViewNamespaceStatsByDay
	case storage.TimeframeWeek:
		view = storage.ViewNamespaceStatsByWeek
	case storage.TimeframeMonth:
		view = storage.ViewNamespaceStatsByMonth
	case storage.TimeframeYear:
		view = storage.ViewNamespaceStatsByYear
	default:
		return nil, errors.Errorf("unexpected timeframe %s", timeframe)
	}

	query := s.db.DB().NewSelect().Table(view).Where("namespace_id = ?", nsId)

	switch name {
	case storage.SeriesNsPfbCount:
		query.ColumnExpr("ts, pfb_count as value")
	case storage.SeriesNsSize:
		query.ColumnExpr("ts, size as value")
	default:
		return nil, errors.Errorf("unexpected series name: %s", name)
	}

	if !req.From.IsZero() {
		query = query.Where("ts >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("ts < ?", req.To)
	}

	err = query.Scan(ctx, &response)
	return
}

func (s Stats) CumulativeSeries(ctx context.Context, timeframe storage.Timeframe, name string, req storage.SeriesRequest) (response []storage.SeriesItem, err error) {
	query := s.db.DB().NewSelect()
	switch timeframe {
	case storage.TimeframeHour:
		query.Table(storage.ViewBlockStatsByHour)
	case storage.TimeframeDay:
		query.Table(storage.ViewBlockStatsByDay)
	case storage.TimeframeWeek:
		query.Table(storage.ViewBlockStatsByWeek)
	case storage.TimeframeMonth:
		query.Table(storage.ViewBlockStatsByMonth)
	case storage.TimeframeYear:
		query.Table(storage.ViewBlockStatsByYear)
	default:
		return nil, errors.Errorf("unexpected timeframe %s", timeframe)
	}

	switch name {
	case storage.SeriesBlobsSize:
		query.ColumnExpr("ts, sum(sum(blobs_size)) OVER(ORDER BY ts) as value")
	case storage.SeriesFee:
		query.ColumnExpr("ts, sum(sum(fee)) OVER(ORDER BY ts) as value")
	case storage.SeriesTxCount:
		query.ColumnExpr("ts, sum(sum(tx_count)) OVER(ORDER BY ts) as value")
	case storage.SeriesGasLimit:
		query.ColumnExpr("ts, sum(sum(gas_limit)) OVER(ORDER BY ts) as value")
	case storage.SeriesGasUsed:
		query.ColumnExpr("ts, sum(sum(gas_used)) OVER(ORDER BY ts) as value")
	case storage.SeriesBytesInBlock:
		query.ColumnExpr("ts, sum(sum(bytes_in_block)) OVER(ORDER BY ts) as value")
	case storage.SeriesBlobsCount:
		query.ColumnExpr("ts, sum(sum(blobs_count)) OVER(ORDER BY ts) as value")
	case storage.SeriesSupplyChange:
		query.ColumnExpr("ts, sum(sum(supply_change)) OVER(ORDER BY ts) as value")
	default:
		return nil, errors.Errorf("unexpected series name: %s", name)
	}

	withQuery := query.Group("ts")

	q := s.db.DB().
		NewSelect().
		With("q", withQuery).
		Table("q")

	if !req.From.IsZero() {
		q = q.Where("ts >= ?", req.From)
	}
	if !req.To.IsZero() {
		q = q.Where("ts < ?", req.To)
	}

	err = q.Scan(ctx, &response)
	return
}

func (s Stats) StakingSeries(ctx context.Context, timeframe storage.Timeframe, name string, validatorId uint64, req storage.SeriesRequest) (response []storage.SeriesItem, err error) {
	var view string
	switch timeframe {
	case storage.TimeframeHour:
		view = storage.ViewStakingByHour
	case storage.TimeframeDay:
		view = storage.ViewStakingByDay
	case storage.TimeframeMonth:
		view = storage.ViewStakingByMonth
	default:
		return nil, errors.Errorf("unexpected timeframe %s", timeframe)
	}

	query := s.db.DB().NewSelect().Table(view).Where("validator_id = ?", validatorId)

	switch name {
	case storage.SeriesRewards:
		query.ColumnExpr("time as ts, rewards as value")
	case storage.SeriesCommissions:
		query.ColumnExpr("time as ts, commissions as value")
	case storage.SeriesFlow:
		query.ColumnExpr("time as ts, flow as value")
	default:
		return nil, errors.Errorf("unexpected series name: %s", name)
	}

	if !req.From.IsZero() {
		query = query.Where("time >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("time < ?", req.To)
	}

	err = query.Scan(ctx, &response)
	return
}

type squareSize struct {
	SquareSize  int       `bun:"square_size"`
	Time        time.Time `bun:"ts"`
	CountBlocks string    `bun:"count_blocks"`
}

func (s Stats) SquareSize(ctx context.Context, from, to *time.Time) (result map[int][]storage.SeriesItem, err error) {
	query := s.db.DB().NewSelect().
		Table(storage.ViewSquareSize).
		OrderExpr("ts desc, square_size desc")

	switch {
	case from == nil && to == nil:
		query.
			Where("ts >= ?", time.Now().AddDate(0, -1, 0).UTC())

	case from != nil && to == nil:
		query.
			Where("ts >= ?", from.UTC()).
			Where("ts < ?", from.AddDate(0, 1, 0).UTC())

	case from == nil && to != nil:
		query.
			Where("ts >= ?", to.AddDate(0, -1, 0).UTC()).
			Where("ts < ?", to.UTC())

	case from != nil && to != nil:
		query.
			Where("ts >= ?", from.UTC()).
			Where("ts < ?", to.UTC())
	}

	var response []squareSize
	if err := query.Scan(ctx, &response); err != nil {
		return result, err
	}

	result = make(map[int][]storage.SeriesItem)
	for _, item := range response {
		seriesItem := storage.SeriesItem{
			Value: item.CountBlocks,
			Time:  item.Time,
		}
		if _, ok := result[item.SquareSize]; !ok {
			result[item.SquareSize] = make([]storage.SeriesItem, 0)
		}
		result[item.SquareSize] = append(result[item.SquareSize], seriesItem)
	}

	return
}

func (s Stats) RollupStats24h(ctx context.Context) (response []storage.RollupStats24h, err error) {
	inner := s.db.DB().NewSelect().
		Table(storage.ViewRollupStatsByHour).
		Column("namespace_id", "signer_id", "size", "fee", "blobs_count").
		Where("time > now() - '1 day'::interval")

	joined := s.db.DB().NewSelect().
		TableExpr("(?) as data", inner).
		ColumnExpr("rollup_id, sum(data.size) as size, sum(data.fee) as fee, sum(data.blobs_count) as blobs_count").
		Join("left join rollup_provider as rp on (rp.address_id = data.signer_id OR rp.address_id = 0) AND (rp.namespace_id = data.namespace_id OR rp.namespace_id = 0)").
		Group("rollup_id")

	err = s.db.DB().NewSelect().
		TableExpr("(?) as grouped", joined).
		ColumnExpr("grouped.*, r.name, r.logo").
		Join("left join rollup as r on r.id = grouped.rollup_id").
		Order("blobs_count desc").
		Scan(ctx, &response)
	return
}

func (s Stats) MessagesCount24h(ctx context.Context) (response []storage.CountItem, err error) {
	err = s.db.DB().NewSelect().
		Model((*storage.Message)(nil)).
		ColumnExpr("count(*) as value, type as name").
		Where("time > now() - '1 day'::interval").
		Group("type").
		Order("value desc").
		Scan(ctx, &response)
	return
}

func (s Stats) SizeGroups(ctx context.Context, timeFilter *time.Time) (groups []storage.SizeGroup, err error) {
	rangeQuery := s.db.DB().NewRaw(`SELECT *
      FROM ( VALUES 
		  (1, 1000, '<1Kb'),
		  (1001, 10000, '1-10Kb'),
		  (10001, 100000, '10-100Kb'),
		  (100001, 1000000, '100Kb-1Mb'),
		  (1000001, 100000000, '>1Mb') 
	  ) AS t(min_val, max_val, name)`)

	if timeFilter == nil {
		tf := time.Now().UTC().AddDate(0, 0, -1)
		timeFilter = &tf
	}

	err = s.db.DB().NewSelect().
		With("ranges", rangeQuery).
		Table("ranges").
		ColumnExpr("ranges.name as name, min(ranges.min_val) as min_val, count(blob_log.*), coalesce(sum(blob_log.size), 0) as size, coalesce(ceil(avg(blob_log.size)), 0) as avg_size").
		Join("left join blob_log on (blob_log.size between ranges.min_val and ranges.max_val) and time >= ?", timeFilter).
		Group("name").
		Order("min_val").
		Scan(ctx, &groups)

	return
}
