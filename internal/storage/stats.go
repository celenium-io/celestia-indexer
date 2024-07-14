// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/stats"
	"github.com/pkg/errors"
)

type CountRequest struct {
	Table string
	From  uint64
	To    uint64
}

func (req CountRequest) Validate() error {
	if _, ok := stats.Tables[req.Table]; !ok {
		return errors.Wrapf(ErrValidation, "unknown table '%s' for stats computing", req.Table)
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
		return errors.Wrapf(ErrValidation, "unknown table '%s' for stats computing", req.Table)
	}

	column, ok := table.Columns[req.Column]
	if !ok {
		return errors.Wrapf(ErrValidation, "unknown column '%s' in table '%s' for stats computing", req.Column, req.Table)
	}

	if _, ok := column.Functions[req.Function]; !ok {
		return errors.Wrapf(ErrValidation, "unknown function '%s' for '%s'.'%s'", req.Function, req.Table, req.Column)
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

type TPS struct {
	Low               float64
	High              float64
	Current           float64
	ChangeLastHourPct float64
}

type Change24hBlockStats struct {
	TxCount      float64 `bun:"tx_count_24h"`
	Fee          float64 `bun:"fee_24h"`
	BlobsSize    float64 `bun:"blobs_size_24h"`
	BytesInBlock float64 `bun:"bytes_in_block_24h"`
}

type SeriesRequest struct {
	From time.Time
	To   time.Time
}

type DistributionItem struct {
	Name  int    `bun:"name"`
	Value string `bun:"value"`
}

type CountItem struct {
	Name  string `bun:"name"`
	Value int64  `bun:"value"`
}

func NewSeriesRequest(from, to int64) SeriesRequest {
	var seriesRequest SeriesRequest
	if from > 0 {
		seriesRequest.From = time.Unix(from, 0).UTC()
	}
	if to > 0 {
		seriesRequest.To = time.Unix(to, 0).UTC()
	}
	return seriesRequest
}

type RollupStats24h struct {
	RollupId   int64   `bun:"rollup_id"`
	Name       string  `bun:"name"`
	Logo       string  `bun:"logo"`
	Size       int64   `bun:"size"`
	Fee        float64 `bun:"fee"`
	BlobsCount int64   `bun:"blobs_count"`
}

type SeriesItem struct {
	Time  time.Time `bun:"ts"`
	Value string    `bun:"value"`
	Max   string    `bun:"max"`
	Min   string    `bun:"min"`
}

const (
	SeriesBlobsSize     = "blobs_size"
	SeriesBlobsCount    = "blobs_count"
	SeriesBlobsFee      = "blobs_fee"
	SeriesTPS           = "tps"
	SeriesBPS           = "bps"
	SeriesFee           = "fee"
	SeriesSupplyChange  = "supply_change"
	SeriesBlockTime     = "block_time"
	SeriesTxCount       = "tx_count"
	SeriesEventsCount   = "events_count"
	SeriesGasPrice      = "gas_price"
	SeriesGasUsed       = "gas_used"
	SeriesGasLimit      = "gas_limit"
	SeriesGasEfficiency = "gas_efficiency"
	SeriesNsPfbCount    = "pfb_count"
	SeriesNsSize        = "size"
	SeriesBytesInBlock  = "bytes_in_block"
	SeriesRewards       = "rewards"
	SeriesCommissions   = "commissions"
	SeriesFlow          = "flow"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IStats interface {
	Count(ctx context.Context, req CountRequest) (string, error)
	Summary(ctx context.Context, req SummaryRequest) (string, error)
	TPS(ctx context.Context) (TPS, error)
	Series(ctx context.Context, timeframe Timeframe, name string, req SeriesRequest) ([]SeriesItem, error)
	CumulativeSeries(ctx context.Context, timeframe Timeframe, name string, req SeriesRequest) ([]SeriesItem, error)
	NamespaceSeries(ctx context.Context, timeframe Timeframe, name string, nsId uint64, req SeriesRequest) (response []SeriesItem, err error)
	StakingSeries(ctx context.Context, timeframe Timeframe, name string, validatorId uint64, req SeriesRequest) (response []SeriesItem, err error)
	RollupStats24h(ctx context.Context) ([]RollupStats24h, error)
	SquareSize(ctx context.Context, from, to *time.Time) (map[int][]SeriesItem, error)
	Change24hBlockStats(ctx context.Context) (response Change24hBlockStats, err error)
	MessagesCount24h(ctx context.Context) ([]CountItem, error)
}
