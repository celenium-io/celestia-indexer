// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/pkg/errors"
)

// Price -
type Price struct {
	db *database.Bun
}

// NewPrice -
func NewPrice(db *database.Bun) *Price {
	return &Price{
		db: db,
	}
}

func (p *Price) Get(ctx context.Context, timeframe string, start, end time.Time, limit int) (price []storage.Price, err error) {
	query := p.db.DB().NewSelect()
	switch timeframe {
	case storage.PriceTimeframeDay:
		query = query.Table("price_by_day")
	case storage.PriceTimeframeHour:
		query = query.Table("price_by_hour")
	case storage.PriceTimeframeMinute:
		query = query.Table("price")
	default:
		return nil, errors.Errorf("unknown price timeframe: %s", timeframe)
	}

	if !start.IsZero() {
		query = query.Where("time >= ?", start)
	}
	if !end.IsZero() {
		query = query.Where("time < ?", end)
	}
	limitScope(query, limit)
	err = query.Order("time desc").Scan(ctx, &price)
	return
}

func (p *Price) Save(ctx context.Context, price *storage.Price) error {
	if price == nil {
		return nil
	}
	_, err := p.db.DB().NewInsert().Model(price).Exec(ctx)
	return err
}

func (p *Price) Last(ctx context.Context) (price storage.Price, err error) {
	err = p.db.DB().NewSelect().Model(&price).Order("time desc").Limit(1).Scan(ctx)
	return
}
