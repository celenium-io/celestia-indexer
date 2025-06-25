// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// ApiKey -
type ApiKey struct {
	db *database.Bun
}

// NewApiKey -
func NewApiKey(db *database.Bun) *ApiKey {
	return &ApiKey{
		db: db,
	}
}

func (ak *ApiKey) Get(ctx context.Context, key string) (apikey storage.ApiKey, err error) {
	apikey.Key = key
	err = ak.db.DB().NewSelect().Model(&apikey).WherePK().Scan(ctx)
	return
}
