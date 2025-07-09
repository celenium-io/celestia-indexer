// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

type IbcClient struct {
	*database.Bun
}

func NewIbcClient(conn *database.Bun) *IbcClient {
	return &IbcClient{conn}
}

func (c *IbcClient) ById(ctx context.Context, id string) (client storage.IbcClient, err error) {
	query := c.DB().NewSelect().
		Model((*storage.IbcClient)(nil)).
		Where("id = ?", id).
		Limit(1)

	err = c.DB().NewSelect().
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

func (c *IbcClient) List(ctx context.Context, limit, offset int, sort sdk.SortOrder) (clients []storage.IbcClient, err error) {
	query := c.DB().NewSelect().
		Model(&clients)

	if offset > 0 {
		query.Offset(offset)
	}

	query = limitScope(query, limit)
	query = sortScope(query, "height", sort)

	err = c.DB().NewSelect().
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
