// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"io"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// BlobLog -
type BlobLog struct {
	*postgres.Table[*storage.BlobLog]

	export *Export
}

// NewBlobLog -
func NewBlobLog(db *database.Bun, export *Export) *BlobLog {
	return &BlobLog{
		Table:  postgres.NewTable[*storage.BlobLog](db),
		export: export,
	}
}

func (bl *BlobLog) ByNamespace(ctx context.Context, nsId uint64, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	blobsQuery := bl.DB().NewSelect().Model((*storage.BlobLog)(nil)).
		Where("blob_log.namespace_id = ?", nsId)

	blobsQuery = blobLogFilters(blobsQuery, fltrs)

	err = bl.DB().NewSelect().
		ColumnExpr("blob_log.*").
		ColumnExpr("signer.address as signer__address").
		ColumnExpr("tx.id as tx__id, tx.height as tx__height, tx.time as tx__time, tx.position as tx__position, tx.gas_wanted as tx__gas_wanted, tx.gas_used as tx__gas_used, tx.timeout_height as tx__timeout_height, tx.events_count as tx__events_count, tx.messages_count as tx__messages_count, tx.fee as tx__fee, tx.status as tx__status, tx.error as tx__error, tx.codespace as tx__codespace, tx.hash as tx__hash, tx.memo as tx__memo, tx.message_types as tx__message_types").
		TableExpr("(?) as blob_log", blobsQuery).
		Join("left join address as signer on signer.id = blob_log.signer_id").
		Join("left join tx on tx.id = blob_log.tx_id").
		Scan(ctx, &logs)
	return
}

func (bl *BlobLog) ByProviders(ctx context.Context, providers []storage.RollupProvider, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	if len(providers) == 0 {
		return nil, nil
	}

	blobQuery := bl.DB().NewSelect().
		Model((*storage.BlobLog)(nil))

	for i := range providers {
		blobQuery = blobQuery.WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			sq.Where("blob_log.signer_id = ?", providers[i].AddressId)
			if providers[i].NamespaceId > 0 {
				sq.Where("blob_log.namespace_id = ?", providers[i].NamespaceId)
			}
			return sq
		})
	}

	blobQuery = blobLogFilters(blobQuery, fltrs)

	err = bl.DB().NewSelect().
		ColumnExpr("blob_log.*").
		ColumnExpr("signer.address as signer__address").
		ColumnExpr("ns.id as namespace__id, ns.size as namespace__size, ns.blobs_count as namespace__blobs_count, ns.version as namespace__version, ns.namespace_id as namespace__namespace_id, ns.reserved as namespace__reserved, ns.pfb_count as namespace__pfb_count, ns.last_height as namespace__last_height, ns.last_message_time as namespace__last_message_time").
		ColumnExpr("tx.id as tx__id, tx.height as tx__height, tx.time as tx__time, tx.position as tx__position, tx.gas_wanted as tx__gas_wanted, tx.gas_used as tx__gas_used, tx.timeout_height as tx__timeout_height, tx.events_count as tx__events_count, tx.messages_count as tx__messages_count, tx.fee as tx__fee, tx.status as tx__status, tx.error as tx__error, tx.codespace as tx__codespace, tx.hash as tx__hash, tx.memo as tx__memo, tx.message_types as tx__message_types").
		TableExpr("(?) as blob_log", blobQuery).
		Join("left join address as signer on signer.id = blob_log.signer_id").
		Join("left join namespace as ns on ns.id = blob_log.namespace_id").
		Join("left join tx on tx.id = blob_log.tx_id").
		Scan(ctx, &logs)
	return
}

const (
	maxExportPeriodInMonth = 1
)

func (bl *BlobLog) ExportByProviders(ctx context.Context, providers []storage.RollupProvider, from, to time.Time, stream io.Writer) (err error) {
	if len(providers) == 0 {
		return nil
	}

	blobQuery := bl.DB().NewSelect().
		Model((*storage.BlobLog)(nil))

	switch {
	case from.IsZero() && to.IsZero():
		blobQuery = blobQuery.
			Where("time >= ?", time.Now().AddDate(0, -maxExportPeriodInMonth, 0).UTC())

	case !from.IsZero() && to.IsZero():
		blobQuery = blobQuery.
			Where("time >= ?", from.UTC()).
			Where("time < ?", from.AddDate(0, maxExportPeriodInMonth, 0).UTC())

	case from.IsZero() && !to.IsZero():
		blobQuery = blobQuery.
			Where("time < ?", to.UTC()).
			Where("time >= ?", to.AddDate(0, -maxExportPeriodInMonth, 0).UTC())

	case !from.IsZero() && !to.IsZero():
		if to.Sub(from) > time.Hour*24*30 {
			blobQuery = blobQuery.
				Where("time >= ?", from.UTC()).
				Where("time < ?", from.AddDate(0, maxExportPeriodInMonth, 0).UTC())
		} else {
			blobQuery = blobQuery.
				Where("time >= ?", from.UTC()).
				Where("time < ?", to.UTC())
		}

	}

	for i := range providers {
		blobQuery = blobQuery.WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			sq.Where("blob_log.signer_id = ?", providers[i].AddressId)
			if providers[i].NamespaceId > 0 {
				sq.Where("blob_log.namespace_id = ?", providers[i].NamespaceId)
			}
			return sq
		})
	}

	query := bl.DB().NewSelect().
		ColumnExpr("blob_log.time, blob_log.height, blob_log.size, blob_log.commitment, blob_log.content_type").
		ColumnExpr("signer.address as signer").
		ColumnExpr("ns.version as namespace_version, ns.namespace_id as namespace_namespace_id").
		ColumnExpr("tx.hash as tx_hash").
		TableExpr("(?) as blob_log", blobQuery).
		Join("left join address as signer on signer.id = blob_log.signer_id").
		Join("left join namespace as ns on ns.id = blob_log.namespace_id").
		Join("left join tx on tx.id = blob_log.tx_id").
		Order("blob_log.time desc").
		String()

	err = bl.export.ToCsv(ctx, stream, query)
	return
}

func (bl *BlobLog) BySigner(ctx context.Context, signerId uint64, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	blobQuery := bl.DB().NewSelect().
		Model((*storage.BlobLog)(nil)).
		Where("signer_id = ?", signerId)

	blobQuery = blobLogFilters(blobQuery, fltrs)

	err = bl.DB().NewSelect().
		ColumnExpr("blob_log.*").
		ColumnExpr("ns.id as namespace__id, ns.size as namespace__size, ns.blobs_count as namespace__blobs_count, ns.version as namespace__version, ns.namespace_id as namespace__namespace_id, ns.reserved as namespace__reserved, ns.pfb_count as namespace__pfb_count, ns.last_height as namespace__last_height, ns.last_message_time as namespace__last_message_time").
		ColumnExpr("tx.id as tx__id, tx.height as tx__height, tx.time as tx__time, tx.position as tx__position, tx.gas_wanted as tx__gas_wanted, tx.gas_used as tx__gas_used, tx.timeout_height as tx__timeout_height, tx.events_count as tx__events_count, tx.messages_count as tx__messages_count, tx.fee as tx__fee, tx.status as tx__status, tx.error as tx__error, tx.codespace as tx__codespace, tx.hash as tx__hash, tx.memo as tx__memo, tx.message_types as tx__message_types").
		TableExpr("(?) as blob_log", blobQuery).
		Join("left join namespace as ns on ns.id = blob_log.namespace_id").
		Join("left join tx on tx.id = blob_log.tx_id").
		Scan(ctx, &logs)
	return
}

func (bl *BlobLog) ByTxId(ctx context.Context, txId uint64, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	blobLogQuery := bl.DB().NewSelect().
		Model((*storage.BlobLog)(nil)).
		Where("tx_id = ?", txId)

	blobLogQuery = blobLogFilters(blobLogQuery, fltrs)

	err = bl.DB().NewSelect().
		ColumnExpr("blob_log.*").
		ColumnExpr("signer.address as signer__address").
		ColumnExpr("ns.id as namespace__id, ns.size as namespace__size, ns.blobs_count as namespace__blobs_count, ns.version as namespace__version, ns.namespace_id as namespace__namespace_id, ns.reserved as namespace__reserved, ns.pfb_count as namespace__pfb_count, ns.last_height as namespace__last_height, ns.last_message_time as namespace__last_message_time").
		ColumnExpr("tx.id as tx__id, tx.height as tx__height, tx.time as tx__time, tx.position as tx__position, tx.gas_wanted as tx__gas_wanted, tx.gas_used as tx__gas_used, tx.timeout_height as tx__timeout_height, tx.events_count as tx__events_count, tx.messages_count as tx__messages_count, tx.fee as tx__fee, tx.status as tx__status, tx.error as tx__error, tx.codespace as tx__codespace, tx.hash as tx__hash, tx.memo as tx__memo, tx.message_types as tx__message_types").
		TableExpr("(?) as blob_log", blobLogQuery).
		Join("left join address as signer on signer.id = blob_log.signer_id").
		Join("left join namespace as ns on ns.id = blob_log.namespace_id").
		Join("left join tx on tx.id = blob_log.tx_id").
		Scan(ctx, &logs)
	return
}

func (bl *BlobLog) ByHeight(ctx context.Context, height types.Level, fltrs storage.BlobLogFilters) (logs []storage.BlobLog, err error) {
	blobLogQuery := bl.DB().NewSelect().
		Model((*storage.BlobLog)(nil)).
		Where("blob_log.height = ?", height)

	blobLogQuery = blobLogFilters(blobLogQuery, fltrs)

	err = bl.DB().NewSelect().
		ColumnExpr("blob_log.*").
		ColumnExpr("signer.address as signer__address").
		ColumnExpr("ns.id as namespace__id, ns.size as namespace__size, ns.blobs_count as namespace__blobs_count, ns.version as namespace__version, ns.namespace_id as namespace__namespace_id, ns.reserved as namespace__reserved, ns.pfb_count as namespace__pfb_count, ns.last_height as namespace__last_height, ns.last_message_time as namespace__last_message_time").
		ColumnExpr("tx.id as tx__id, tx.height as tx__height, tx.time as tx__time, tx.position as tx__position, tx.gas_wanted as tx__gas_wanted, tx.gas_used as tx__gas_used, tx.timeout_height as tx__timeout_height, tx.events_count as tx__events_count, tx.messages_count as tx__messages_count, tx.fee as tx__fee, tx.status as tx__status, tx.error as tx__error, tx.codespace as tx__codespace, tx.hash as tx__hash, tx.memo as tx__memo, tx.message_types as tx__message_types").
		TableExpr("(?) as blob_log", blobLogQuery).
		Join("left join address as signer on signer.id = blob_log.signer_id").
		Join("left join namespace as ns on ns.id = blob_log.namespace_id").
		Join("left join tx on tx.id = blob_log.tx_id").
		Scan(ctx, &logs)
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
