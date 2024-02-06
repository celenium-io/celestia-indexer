// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
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
			sq.Where("blob_log.signer_id = ?", providers[i].AddressId)
			if providers[i].NamespaceId > 0 {
				sq.Where("blob_log.namespace_id = ?", providers[i].NamespaceId)
			}
			return sq
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

func (bl *BlobLog) ByTxId(ctx context.Context, txId uint64, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	query := bl.DB().NewSelect().Model(&logs).
		Relation("Namespace").
		Relation("Tx").
		Relation("Signer").
		Where("tx_id = ?", txId)

	query = blobLogFilters(query, fltrs)

	err = query.Scan(ctx)
	return
}

func (bl *BlobLog) ByHeight(ctx context.Context, height types.Level, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	query := bl.DB().NewSelect().Model(&logs).
		Relation("Namespace").
		Relation("Tx").
		Relation("Signer").
		Where("blob_log.height = ?", height)

	query = blobLogFilters(query, fltrs)

	err = query.Scan(ctx)
	return
}

func (bl *BlobLog) CountByTxId(ctx context.Context, txId uint64) (int, error) {
	return bl.DB().NewSelect().Model((*storage.BlobLog)(nil)).
		Where("tx_id = ?", txId).
		Count(ctx)
}

func (bl *BlobLog) CountByHeight(ctx context.Context, height types.Level) (int, error) {
	return bl.DB().NewSelect().Model((*storage.BlobLog)(nil)).
		Where("height = ?", height).
		Count(ctx)
}
