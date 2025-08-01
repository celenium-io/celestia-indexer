// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// DenomMetadata -
type DenomMetadata struct {
	db *database.Bun
}

// NewDenomMetadata -
func NewDenomMetadata(db *database.Bun) *DenomMetadata {
	return &DenomMetadata{
		db: db,
	}
}

func (dm *DenomMetadata) All(ctx context.Context) (metadata []storage.DenomMetadata, err error) {
	err = dm.db.DB().NewSelect().Model(&metadata).Scan(ctx)
	return
}
