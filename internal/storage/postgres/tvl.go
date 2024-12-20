// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// Tvl -
type Tvl struct {
	db *database.Bun
}

// NewTvl -
func NewTvl(db *database.Bun) *Tvl {
	return &Tvl{
		db: db,
	}
}

func (t *Tvl) Save(ctx context.Context, tvl *storage.Tvl) error {
	if tvl == nil {
		return nil
	}
	_, err := t.db.DB().NewInsert().Model(tvl).Exec(ctx)
	return err
}

func (t *Tvl) SaveBulk(ctx context.Context, tvls ...*storage.Tvl) error {
	if len(tvls) == 0 {
		return nil
	}
	_, err := t.db.DB().NewInsert().Model(&tvls).Exec(ctx)
	return err
}
