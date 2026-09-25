// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

type HLIGPConfig struct {
	db bun.IDB
}

func NewHLIGPConfig(db bun.IDB) storage.IHLIGPConfig {
	return &HLIGPConfig{db}
}

func (hl *HLIGPConfig) List(ctx context.Context, limit, offset int) (config []storage.HLIGPConfig, err error) {
	query := hl.db.NewSelect().
		Model((*storage.HLIGPConfig)(nil))

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx, &config)

	return
}
