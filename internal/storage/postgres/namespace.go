// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Namespace -
type Namespace struct {
	*postgres.Table[*storage.Namespace]
}

// NewNamespace -
func NewNamespace(db *database.Bun) *Namespace {
	return &Namespace{
		Table: postgres.NewTable[*storage.Namespace](db),
	}
}

// ByNamespaceId -
func (n *Namespace) ByNamespaceId(ctx context.Context, namespaceId []byte) (namespace []storage.Namespace, err error) {
	err = n.DB().NewSelect().Model(&namespace).
		Where("namespace_id = ?", namespaceId).
		Scan(ctx)
	return
}

// ByNamespaceIdAndVersion -
func (n *Namespace) ByNamespaceIdAndVersion(ctx context.Context, namespaceId []byte, version byte) (namespace storage.Namespace, err error) {
	err = n.DB().NewSelect().Model(&namespace).
		Where("namespace_id = ?", namespaceId).
		Where("version = ?", version).
		Scan(ctx)
	return
}

// Messages -
func (n *Namespace) Messages(ctx context.Context, id uint64, limit, offset int) (msgs []storage.NamespaceMessage, err error) {
	subQuery := n.DB().NewSelect().Model(&msgs).
		Where("namespace_id = ?", id).
		Order("time desc")

	subQuery = limitScope(subQuery, limit)
	if offset > 0 {
		subQuery = subQuery.Offset(offset)
	}

	query := n.DB().NewSelect().
		TableExpr("(?) as msgs", subQuery).
		ColumnExpr("msgs.*").
		ColumnExpr("namespace.id as namespace__id, namespace.first_height as namespace__first_height, namespace.last_height as namespace__last_height, namespace.version as namespace__version, namespace.namespace_id as namespace__namespace_id, namespace.size as namespace__size, namespace.pfb_count as namespace__pfb_count, namespace.reserved as namespace__reserved, namespace.last_message_time as namespace__last_message_time").
		ColumnExpr("message.id as message__id, message.height as message__height, message.time as message__time, message.position as message__position, message.type as message__type, message.tx_id as message__tx_id, message.data as message__data").
		ColumnExpr("tx.id as tx__id, tx.height as tx__height, tx.time as tx__time, tx.position as tx__position, tx.gas_wanted as tx__gas_wanted, tx.gas_used as tx__gas_used, tx.timeout_height as tx__timeout_height, tx.events_count as tx__events_count, tx.messages_count as tx__messages_count, tx.fee as tx__fee, tx.status as tx__status, tx.error as tx__error, tx.codespace as tx__codespace, tx.hash as tx__hash, tx.memo as tx__memo, tx.message_types as tx__message_types").
		Join("LEFT JOIN namespace ON namespace.id = msgs.namespace_id").
		Join("LEFT JOIN message ON message.id = msgs.msg_id").
		Join("LEFT JOIN tx ON tx.id = msgs.tx_id").
		Order("msgs.time desc")
	err = query.Scan(ctx, &msgs)
	return
}

// MessagesByHeight -
func (n *Namespace) MessagesByHeight(ctx context.Context, height pkgTypes.Level, limit, offset int) (msgs []storage.NamespaceMessage, err error) {
	query := n.DB().NewSelect().Model(&msgs).
		Where("namespace_message.height = ?", height).
		Relation("Namespace").
		Relation("Message").
		Relation("Tx")
	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}

func (n *Namespace) CountMessagesByHeight(ctx context.Context, height pkgTypes.Level) (int, error) {
	return n.DB().NewSelect().Model((*storage.NamespaceMessage)(nil)).
		Where("namespace_message.height = ?", height).
		Count(ctx)
}

func (n *Namespace) ListWithSort(ctx context.Context, sortField string, sort sdk.SortOrder, limit, offset int) (ns []storage.Namespace, err error) {
	var field string
	switch sortField {
	case timeColumn:
		field = "last_message_time"
	case pfbCountColumn:
		field = pfbCountColumn
	case sizeColumn:
		field = sizeColumn
	default:
		field = "id"
	}

	if offset < 0 {
		offset = 0
	}

	query := n.DB().NewSelect().Model(&ns)
	limitScope(query, limit)
	sortScope(query, field, sort)

	err = query.Offset(offset).Scan(ctx)
	return
}

func (n *Namespace) MessagesByTxId(ctx context.Context, txId uint64, limit, offset int) (msgs []storage.NamespaceMessage, err error) {
	query := n.DB().NewSelect().Model(&msgs).
		Where("namespace_message.tx_id = ?", txId).
		Relation("Namespace").
		Relation("Message").
		Relation("Tx")
	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}

func (n *Namespace) CountMessagesByTxId(ctx context.Context, txId uint64) (int, error) {
	return n.DB().NewSelect().Model((*storage.NamespaceMessage)(nil)).
		Where("namespace_message.tx_id = ?", txId).
		Count(ctx)
}

func (n *Namespace) GetByIds(ctx context.Context, ids ...uint64) (ns []storage.Namespace, err error) {
	if len(ids) == 0 {
		return nil, nil
	}

	err = n.DB().NewSelect().Model(&ns).Where("id IN (?)", bun.In(ids)).Scan(ctx)
	return
}
