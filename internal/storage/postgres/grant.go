// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Grant -
type Grant struct {
	*postgres.Table[*storage.Grant]
}

// NewGrant -
func NewGrant(db *database.Bun) *Grant {
	return &Grant{
		Table: postgres.NewTable[*storage.Grant](db),
	}
}

func (g *Grant) ByGrantee(ctx context.Context, id uint64, limit, offset int) (grants []storage.Grant, err error) {
	query := g.DB().NewSelect().
		Model((*storage.Grant)(nil)).
		Where("grantee_id = ?", id).
		Order("id desc")

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = g.DB().NewSelect().
		TableExpr("(?) as g", query).
		ColumnExpr("g.*").
		ColumnExpr("address.address as granter__address").
		ColumnExpr("celestial.id as granter__celestials__id, celestial.image_url as granter__celestials__image_url").
		Join("left join address on address.id = g.granter_id").
		Join("left join celestial on celestial.address_id = g.granter_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &grants)
	return
}

func (g *Grant) ByGranter(ctx context.Context, id uint64, limit, offset int) (grants []storage.Grant, err error) {
	query := g.DB().NewSelect().
		Model((*storage.Grant)(nil)).
		Where("granter_id = ?", id).
		Order("id desc")

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = g.DB().NewSelect().
		TableExpr("(?) as g", query).
		ColumnExpr("g.*").
		ColumnExpr("address.address as grantee__address").
		ColumnExpr("celestial.id as grantee__celestials__id, celestial.image_url as grantee__celestials__image_url").
		Join("left join address on address.id = g.grantee_id").
		Join("left join celestial on celestial.address_id = g.grantee_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &grants)
	return
}
