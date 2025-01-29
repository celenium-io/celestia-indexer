// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ITvl interface {
	LastSyncTime(ctx context.Context) (time.Time, error)

	Save(ctx context.Context, rollupTvl *Tvl) error
	SaveBulk(ctx context.Context, tvls ...*Tvl) error
}

// Tvl -
type Tvl struct {
	bun.BaseModel `bun:"tvl" comment:"Table with rollup TVL."`

	Id       uint64          `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Value    decimal.Decimal `bun:"value,type:numeric"        comment:"Value of TVL"`
	Time     time.Time       `bun:"time,pk,notnull"           comment:"The time of TVL"`
	RollupId uint64          `bun:"rollup_id"                 comment:"Rollup id"`

	Rollup *Rollup `bun:"rel:belongs-to,join:rollup_id=id"`
}

// TableName -
func (Tvl) TableName() string {
	return "tvl"
}
