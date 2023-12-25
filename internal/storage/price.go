// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

const (
	PriceTimeframeMinute = "1m"
	PriceTimeframeHour   = "1h"
	PriceTimeframeDay    = "1d"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IPrice interface {
	Save(ctx context.Context, price *Price) error
	Last(ctx context.Context) (Price, error)
	Get(ctx context.Context, timeframe string, start, end int64, limit int) ([]Price, error)
}

type Price struct {
	bun.BaseModel `bun:"price" comment:"Table with TIA price"`

	Time  time.Time       `bun:"time,pk"             comment:"Time of candles"`
	Open  decimal.Decimal `bun:"open,,type:numeric"  comment:"Open price"`
	High  decimal.Decimal `bun:"high,,type:numeric"  comment:"High price"`
	Low   decimal.Decimal `bun:"low,,type:numeric"   comment:"Low price"`
	Close decimal.Decimal `bun:"close,,type:numeric" comment:"Close price"`
}

// TableName -
func (Price) TableName() string {
	return "price"
}
