// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

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
	query := n.DB().NewSelect().Model(&msgs).
		Where("namespace_message.namespace_id = ?", id).
		Order("namespace_message.time desc").
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

// MessagesByHeight -
func (n *Namespace) MessagesByHeight(ctx context.Context, height pkgTypes.Level, limit, offset int) (msgs []storage.NamespaceMessage, err error) {
	query := n.DB().NewSelect().Model(&msgs).
		Where("namespace_message.height = ?", height).
		Order("namespace_message.time desc").
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
	case "time":
		field = "last_message_time"
	case "pfb_count":
		field = "pfb_count"
	case "size":
		field = "size"
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
