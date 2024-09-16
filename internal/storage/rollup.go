// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IRollup interface {
	sdk.Table[*Rollup]

	Leaderboard(ctx context.Context, sortField string, sort sdk.SortOrder, limit, offset int) ([]RollupWithStats, error)
	LeaderboardDay(ctx context.Context, sortField string, sort sdk.SortOrder, limit, offset int) ([]RollupWithDayStats, error)
	Namespaces(ctx context.Context, rollupId uint64, limit, offset int) (namespaceIds []uint64, err error)
	Providers(ctx context.Context, rollupId uint64) (providers []RollupProvider, err error)
	RollupsByNamespace(ctx context.Context, namespaceId uint64, limit, offset int) (rollups []Rollup, err error)
	ById(ctx context.Context, rollupId uint64) (RollupWithStats, error)
	Series(ctx context.Context, rollupId uint64, timeframe, column string, req SeriesRequest) (items []HistogramItem, err error)
	AllSeries(ctx context.Context) ([]RollupHistogramItem, error)
	Count(ctx context.Context) (int64, error)
	Distribution(ctx context.Context, rollupId uint64, series, groupBy string) (items []DistributionItem, err error)
	BySlug(ctx context.Context, slug string) (RollupWithStats, error)
}

// Rollup -
type Rollup struct {
	bun.BaseModel `bun:"rollup" comment:"Table with rollups."`

	Id             uint64               `bun:"id,pk,autoincrement"           comment:"Unique internal identity"`
	Name           string               `bun:"name"                          comment:"Rollup's name"`
	Description    string               `bun:"description"                   comment:"Rollup's description"`
	Website        string               `bun:"website"                       comment:"Website"`
	GitHub         string               `bun:"github"                        comment:"Github repository"`
	Twitter        string               `bun:"twitter"                       comment:"Twitter account"`
	Logo           string               `bun:"logo"                          comment:"Link to rollup logo"`
	Slug           string               `bun:"slug,unique:rollup_slug"       comment:"Rollup slug"`
	BridgeContract string               `bun:"bridge_contract"               comment:"Link to bridge contract"`
	L2Beat         string               `bun:"l2_beat"                       comment:"Link to L2 Beat"`
	Explorer       string               `bun:"explorer"                      comment:"Link to chain explorer"`
	Stack          string               `bun:"stack"                         comment:"Underlaying stack"`
	Compression    string               `bun:"compression"                   comment:"Compression"`
	Provider       string               `bun:"provider"                      comment:"RaaS provider"`
	Type           types.RollupType     `bun:"type,type:rollup_type"         comment:"Type of rollup: settled or sovereign"`
	Category       types.RollupCategory `bun:"category,type:rollup_category" comment:"Category of rollup"`
	VM             string               `bun:"vm"                            comment:"Virtual machine"`
	Links          []string             `bun:"links,array"                   comment:"Other links to rollup related sites"`

	Providers []*RollupProvider `bun:"rel:has-many,join:id=rollup_id"`
}

// TableName -
func (Rollup) TableName() string {
	return "rollup"
}

func (r Rollup) IsEmpty() bool {
	return r.Description == "" &&
		r.GitHub == "" &&
		r.Name == "" &&
		r.Twitter == "" &&
		r.Website == "" &&
		r.Logo == "" &&
		r.L2Beat == "" &&
		r.BridgeContract == "" &&
		r.Explorer == "" &&
		r.Stack == "" &&
		r.Links == nil &&
		r.Category == "" &&
		r.Compression == "" &&
		r.Provider == "" &&
		r.Type == "" &&
		r.VM == ""
}

type RollupWithStats struct {
	Rollup
	RollupStats
}

type RollupStats struct {
	Size            int64           `bun:"size"`
	BlobsCount      int64           `bun:"blobs_count"`
	LastActionTime  time.Time       `bun:"last_time"`
	FirstActionTime time.Time       `bun:"first_time"`
	Fee             decimal.Decimal `bun:"fee"`
	SizePct         float64         `bun:"size_pct"`
	FeePct          float64         `bun:"fee_pct"`
	BlobsCountPct   float64         `bun:"blobs_count_pct"`
}

type RollupWithDayStats struct {
	Rollup
	RolluDayStats
}

type RolluDayStats struct {
	AvgSize        float64         `bun:"avg_size"`
	BlobsCount     int64           `bun:"blobs_count"`
	TotalSize      int64           `bun:"total_size"`
	Throghput      int64           `bun:"throughput"`
	TotalFee       decimal.Decimal `bun:"total_fee"`
	NamespaceCount int64           `bun:"namespace_count"`
	PfbCount       int64           `bun:"pfb_count"`
	MBPrice        decimal.Decimal `bun:"mb_price"`
}

type RollupHistogramItem struct {
	Fee        string    `bun:"fee"`
	BlobsCount int64     `bun:"blobs_count"`
	Size       int64     `bun:"size"`
	Name       string    `bun:"name"`
	Logo       string    `bun:"logo"`
	Time       time.Time `bun:"time"`
}
