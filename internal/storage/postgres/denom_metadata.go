// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

// DenomMetadata -
type DenomMetadata struct {
	db bun.IDB
}

// NewDenomMetadata -
func NewDenomMetadata(db bun.IDB) storage.IDenomMetadata {
	return &DenomMetadata{
		db: db,
	}
}

func (dm *DenomMetadata) All(ctx context.Context) (metadata []storage.DenomMetadata, err error) {
	err = dm.db.NewSelect().Model(&metadata).Scan(ctx)
	return
}
