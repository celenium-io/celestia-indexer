// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

type IbcChannel struct {
	*database.Bun
}

func NewIbcChannel(conn *database.Bun) *IbcChannel {
	return &IbcChannel{conn}
}

func (c *IbcChannel) ById(ctx context.Context, id string) (channel storage.IbcChannel, err error) {
	query := c.DB().NewSelect().
		Model((*storage.IbcChannel)(nil)).
		Where("id = ?", id).
		Limit(1)

	err = c.DB().NewSelect().
		TableExpr("(?) as ibc_channel", query).
		ColumnExpr("ibc_channel.*").
		ColumnExpr("create_tx.hash as create_tx__hash").
		ColumnExpr("confirmation_tx.hash as confirmation_tx__hash").
		ColumnExpr("address.address as creator__address").
		ColumnExpr("celestial.id as creator__celestials__id, celestial.image_url as creator__celestials__image_url").
		ColumnExpr("ibc_client.chain_id as client__chain_id, ibc_client.type as client__type, ibc_client.connection_count as client__connection_count").
		Join("left join tx as create_tx on create_tx_id = create_tx.id").
		Join("left join tx as confirmation_tx on confirmation_tx_id = confirmation_tx.id").
		Join("left join address on address.id = creator_id").
		Join("left join celestial on celestial.address_id = creator_id and celestial.status = 'PRIMARY'").
		Join("left join ibc_client on ibc_client.id = ibc_channel.client_id").
		Scan(ctx, &channel)
	return
}

func (c *IbcChannel) List(ctx context.Context, fltrs storage.ListChannelFilters) (channels []storage.IbcChannel, err error) {
	query := c.DB().NewSelect().
		Model((*storage.IbcChannel)(nil))

	if fltrs.Offset > 0 {
		query.Offset(fltrs.Offset)
	}

	query = limitScope(query, fltrs.Limit)
	query = sortScope(query, "height", fltrs.Sort)

	if fltrs.ClientId != "" {
		query = query.Where("client_id = ?", fltrs.ClientId)
	}
	if fltrs.ConnectionId != "" {
		query = query.Where("connection_id = ?", fltrs.ConnectionId)
	}

	if fltrs.Status != "" {
		query = query.Where("status = ?", fltrs.Status)
	}

	err = c.DB().NewSelect().
		TableExpr("(?) as ibc_channel", query).
		ColumnExpr("ibc_channel.*").
		ColumnExpr("create_tx.hash as create_tx__hash").
		ColumnExpr("confirmation_tx.hash as confirmation_tx__hash").
		ColumnExpr("address.address as creator__address").
		ColumnExpr("celestial.id as creator__celestials__id, celestial.image_url as creator__celestials__image_url").
		ColumnExpr("ibc_client.chain_id as client__chain_id, ibc_client.type as client__type, ibc_client.connection_count as client__connection_count").
		Join("left join tx as create_tx on create_tx_id = create_tx.id").
		Join("left join tx as confirmation_tx on confirmation_tx_id = confirmation_tx.id").
		Join("left join address on address.id = creator_id").
		Join("left join celestial on celestial.address_id = creator_id and celestial.status = 'PRIMARY'").
		Join("left join ibc_client on ibc_client.id = ibc_channel.client_id").
		Scan(ctx, &channels)
	return
}

func (c *IbcChannel) StatsByChain(ctx context.Context, limit, offset int) (stats []storage.ChainStats, err error) {
	query := c.DB().NewSelect().
		Model((*storage.IbcChannel)(nil)).
		ColumnExpr("ibc_client.chain_id as chain_id").
		ColumnExpr("sum(received) as received").
		ColumnExpr("sum(sent) as sent").
		ColumnExpr("sum(received) + sum(sent) as flow").
		Join("left join ibc_client on ibc_client.id = ibc_channel.client_id").
		Group("chain_id")

	q := c.DB().NewSelect().
		With("stats", query).
		Table("stats").
		OrderExpr("flow desc").
		Offset(offset)

	q = limitScope(q, limit)

	err = q.Scan(ctx, &stats)
	return
}
