// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

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

func (s Stats) TxCountForLast24h(ctx context.Context) (response []storage.TxCountForLast24hItem, err error) {
	err = s.db.DB().NewSelect().Table(storage.ViewBlockStatsByHour).
		Column("ts", "tps", "tx_count").
		Where("ts = date_trunc('hour', now()) - '1 day'::interval").
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

	err = query.Limit(200).Scan(ctx, &response)
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

	err = query.Limit(100).Scan(ctx, &response)
	return
}

func (s Stats) CumulativeSeries(ctx context.Context, timeframe storage.Timeframe, name string, req storage.SeriesRequest) (response []storage.SeriesItem, err error) {
	query := s.db.DB().NewSelect()
	switch timeframe {
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

	err = query.Limit(100).Scan(ctx, &response)
	return
}
