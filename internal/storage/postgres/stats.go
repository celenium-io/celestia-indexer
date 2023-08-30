package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
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

func (s Stats) HistogramCount(ctx context.Context, req storage.HistogramCountRequest) (response []storage.HistogramItem, err error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	query := s.db.DB().NewSelect().Table(req.Table).
		ColumnExpr(`COUNT(*) as value`).
		Group("bucket").
		Order("bucket desc")

	query, err = timeframeScope(query, req.Timeframe)
	if err != nil {
		return
	}

	if req.From > 0 {
		query = query.Where("time >= to_timestamp(?)", req.From)
	}
	if req.To > 0 {
		query = query.Where("time < to_timestamp(?)", req.To)
	}

	err = query.Scan(ctx, &response)
	return
}

func (s Stats) Histogram(ctx context.Context, req storage.HistogramRequest) (response []storage.HistogramItem, err error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	query := s.db.DB().NewSelect().Table(req.Table).
		ColumnExpr(`? (?) as value`, bun.Safe(req.Function), bun.Safe(req.Column)).
		Group("bucket").
		Order("bucket desc")

	query, err = timeframeScope(query, req.Timeframe)
	if err != nil {
		return
	}

	if req.From > 0 {
		query = query.Where("time >= to_timestamp(?)", req.From)
	}
	if req.To > 0 {
		query = query.Where("time < to_timestamp(?)", req.To)
	}

	err = query.Scan(ctx, &response)
	return
}
