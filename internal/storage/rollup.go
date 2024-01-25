// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IRollup interface {
	storage.Table[*Rollup]

	Leaderboard(ctx context.Context, sortField string, sort sdk.SortOrder, limit, offset int) ([]RollupWithStats, error)
	Namespaces(ctx context.Context, rollupId uint64, limit, offset int) (namespaceIds []uint64, err error)
	Providers(ctx context.Context, rollupId uint64) (providers []RollupProvider, err error)
	Stats(ctx context.Context, rollupId uint64) (RollupStats, error)
	Series(ctx context.Context, rollupId uint64, timeframe, column string, req SeriesRequest) (items []HistogramItem, err error)
	Count(ctx context.Context) (int64, error)
}

// Rollup -
type Rollup struct {
	bun.BaseModel `bun:"rollup" comment:"Table with rollups."`

	Id          uint64 `bun:"id,pk,autoincrement" comment:"Unique internal identity"`
	Name        string `bun:"name"                comment:"Rollup's name"`
	Description string `bun:"description"         comment:"Rollup's description"`
	Website     string `bun:"website"             comment:"Website"`
	GitHub      string `bun:"github"              comment:"Github repository"`
	Twitter     string `bun:"twitter"             comment:"Twitter account"`
	Logo        string `bun:"logo"                comment:"Link to rollup logo"`

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
		r.Logo == ""
}

type RollupWithStats struct {
	Rollup
	RollupStats
}

type RollupStats struct {
	Size           int64     `bun:"size"`
	BlobsCount     int64     `bun:"blobs_count"`
	LastActionTime time.Time `bun:"last_time"`
}
