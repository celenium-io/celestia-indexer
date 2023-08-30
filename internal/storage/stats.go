package storage

import (
	"context"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/stats"
	"github.com/pkg/errors"
)

type CountRequest struct {
	Table string
	From  uint64
	To    uint64
}

func (req CountRequest) Validate() error {
	if _, ok := stats.Tables[req.Table]; !ok {
		return errors.Errorf("unknown table '%s' for stats computing", req.Table)
	}

	return nil
}

type SummaryRequest struct {
	CountRequest
	Column   string
	Function string
}

func (req SummaryRequest) Validate() error {
	table, ok := stats.Tables[req.Table]
	if !ok {
		return errors.Errorf("unknown table '%s' for stats computing", req.Table)
	}

	column, ok := table.Columns[req.Column]
	if !ok {
		return errors.Errorf("unknown column '%s' in table '%s' for stats computing", req.Column, req.Table)
	}

	if _, ok := column.Functions[req.Function]; !ok {
		return errors.Errorf("unknown function '%s' for '%s'.'%s'", req.Function, req.Table, req.Column)
	}

	return nil
}

type Timeframe string

const (
	TimeframeHour  Timeframe = "hour"
	TimeframeDay   Timeframe = "day"
	TimeframeWeek  Timeframe = "week"
	TimeframeMonth Timeframe = "month"
	TimeframeYear  Timeframe = "year"
)

type HistogramRequest struct {
	SummaryRequest
	Timeframe Timeframe
}

func (req HistogramRequest) Validate() error {
	if err := req.SummaryRequest.Validate(); err != nil {
		return err
	}
	return nil
}

type HistogramCountRequest struct {
	CountRequest
	Timeframe Timeframe
}

func (req HistogramCountRequest) Validate() error {
	if err := req.CountRequest.Validate(); err != nil {
		return err
	}
	return nil
}

type HistogramItem struct {
	Time  time.Time `bun:"bucket"`
	Value string    `bun:"value"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IStats interface {
	Count(ctx context.Context, req CountRequest) (string, error)
	Summary(ctx context.Context, req SummaryRequest) (string, error)
	HistogramCount(ctx context.Context, req HistogramCountRequest) ([]HistogramItem, error)
	Histogram(ctx context.Context, req HistogramRequest) ([]HistogramItem, error)
}
