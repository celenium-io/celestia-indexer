// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
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
		Relation("Signer").
		Relation("Tx")

	query = blobLogFilters(query, fltrs)

	err = query.Scan(ctx)
	return
}

func (bl *BlobLog) ByProviders(ctx context.Context, providers []storage.RollupProvider, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	if len(providers) == 0 {
		return nil, nil
	}

	query := bl.DB().NewSelect().Model(&logs).
		Relation("Signer").
		Relation("Namespace").
		Relation("Tx")

	for i := range providers {
		query = query.WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.
				Where("blob_log.namespace_id = ?", providers[i].NamespaceId).
				Where("blob_log.signer_id = ?", providers[i].AddressId)
		})
	}

	query = blobLogFilters(query, fltrs)

	err = query.Scan(ctx)
	return
}

func (bl *BlobLog) BySigner(ctx context.Context, signerId uint64, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	query := bl.DB().NewSelect().Model(&logs).
		Relation("Namespace").
		Relation("Tx").
		Where("signer_id = ?", signerId)

	query = blobLogFilters(query, fltrs)

	err = query.Scan(ctx)
	return
}
