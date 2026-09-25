// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

type IbcClient struct {
	db bun.IDB
}

func NewIbcClient(db bun.IDB) storage.IIbcClient {
	return &IbcClient{db}
}

func (c *IbcClient) ById(ctx context.Context, id string) (client storage.IbcClient, err error) {
	query := c.db.NewSelect().
		Model((*storage.IbcClient)(nil)).
		Where("id = ?", id).
		Limit(1)

	err = c.db.NewSelect().
		TableExpr("(?) as ibc_client", query).
		ColumnExpr("ibc_client.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as creator__address").
		ColumnExpr("celestial.id as creator__celestials__id, celestial.image_url as creator__celestials__image_url").
		Join("left join tx on tx_id = tx.id").
		Join("left join address on address.id = creator_id").
		Join("left join celestial on celestial.address_id = creator_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &client)
	return
}

func (c *IbcClient) List(ctx context.Context, filters storage.ListIbcClientsFilters) (clients []storage.IbcClient, err error) {
	query := c.db.NewSelect().
		Model(&clients)

	if filters.Offset > 0 {
		query.Offset(filters.Offset)
	}

	query = limitScope(query, filters.Limit)
	query = sortScope(query, "height", filters.Sort)

	if filters.CreatorId != nil {
		query = query.Where("creator_id = ?", *filters.CreatorId)
	}

	if filters.ChainId != "" {
		query = query.Where("chain_id = ?", filters.ChainId)
	}

	err = c.db.NewSelect().
		TableExpr("(?) as ibc_client", query).
		ColumnExpr("ibc_client.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("celestial.id as creator__celestials__id, celestial.image_url as creator__celestials__image_url").
		ColumnExpr("address.address as creator__address").
		Join("left join tx on tx_id = tx.id").
		Join("left join address on address.id = creator_id").
		Join("left join celestial on celestial.address_id = creator_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &clients)
	return
}

func (c *IbcClient) ByChainId(ctx context.Context, chainId string) (res []string, err error) {
	err = c.db.NewSelect().
		Column("id").
		Model((*storage.IbcClient)(nil)).
		Where("chain_id = ?", chainId).
		Scan(ctx, &res)
	return
}
