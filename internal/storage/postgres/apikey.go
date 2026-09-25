// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

// ApiKey -
type ApiKey struct {
	db bun.IDB
}

// NewApiKey -
func NewApiKey(db bun.IDB) storage.IApiKey {
	return &ApiKey{
		db: db,
	}
}

func (ak *ApiKey) Get(ctx context.Context, key string) (apikey storage.ApiKey, err error) {
	apikey.Key = key
	err = ak.db.NewSelect().Model(&apikey).WherePK().Scan(ctx)
	return
}
