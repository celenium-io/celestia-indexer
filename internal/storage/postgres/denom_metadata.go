// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
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
