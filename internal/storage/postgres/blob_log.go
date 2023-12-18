// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
)

// BlobLog -
type BlobLog struct {
	*postgres.Table[*storage.BlobLog]
}

// NewBlobLog -
func NewBlobLog(db *database.Bun) *BlobLog {
	return &BlobLog{
		Table: postgres.NewTable[*storage.BlobLog](db),
	}
}

func (bl *BlobLog) ByNamespace(ctx context.Context, nsId uint64, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	query := bl.DB().NewSelect().Model(&logs).
		Where("blob_log.namespace_id = ?", nsId).
		Relation("Signer")

	if fltrs.Offset > 0 {
		query.Offset(fltrs.Offset)
	}

	query = limitScope(query, fltrs.Limit)

	switch fltrs.SortBy {
	case sizeColumn, "time":
		query = sortScope(query, fltrs.SortBy, fltrs.Sort)
	case "":
		query = sortScope(query, "id", fltrs.Sort)
	default:
		return nil, errors.Errorf("invalid sort by parameter: %s", fltrs.SortBy)
	}

	err = query.Scan(ctx)
	return
}
