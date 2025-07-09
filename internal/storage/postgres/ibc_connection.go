// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

type IbcConnection struct {
	*database.Bun
}

func NewIbcConnection(conn *database.Bun) *IbcConnection {
	return &IbcConnection{conn}
}

func (c *IbcConnection) ById(ctx context.Context, id string) (conn storage.IbcConnection, err error) {
	query := c.DB().NewSelect().
		Model((*storage.IbcConnection)(nil)).
		Where("connection_id = ?", id).
		Limit(1)

	err = c.DB().NewSelect().
		TableExpr("(?) as ibc_connection", query).
		ColumnExpr("ibc_connection.*").
		ColumnExpr("create_tx.hash as create_tx__hash").
		ColumnExpr("connect_tx.hash as connection_tx__hash").
		ColumnExpr("ibc_client.chain_id as client__chain_id, ibc_client.type as client__type, ibc_client.connection_count as client__connection_count").
		Join("left join tx as create_tx on create_tx_id = create_tx.id").
		Join("left join tx as connect_tx on connection_tx_id = connect_tx.id").
		Join("left join ibc_client on client_id = ibc_client.id").
		Scan(ctx, &conn)
	return
}

func (c *IbcConnection) List(ctx context.Context, fltrs storage.ListConnectionFilters) (conns []storage.IbcConnection, err error) {
	query := c.DB().NewSelect().
		Model(&conns)

	if fltrs.Offset > 0 {
		query.Offset(fltrs.Offset)
	}

	query = limitScope(query, fltrs.Limit)
	query = sortScope(query, "height", fltrs.Sort)

	if fltrs.ClientId != "" {
		query = query.Where("client_id = ?", fltrs.ClientId)
	}

	err = c.DB().NewSelect().
		TableExpr("(?) as ibc_connection", query).
		ColumnExpr("ibc_connection.*").
		ColumnExpr("create_tx.hash as create_tx__hash").
		ColumnExpr("connect_tx.hash as connection_tx__hash").
		ColumnExpr("ibc_client.chain_id as client__chain_id, ibc_client.type as client__type, ibc_client.connection_count as client__connection_count").
		Join("left join tx as create_tx on create_tx_id = create_tx.id").
		Join("left join tx as connect_tx on connection_tx_id = connect_tx.id").
		Join("left join ibc_client on client_id = ibc_client.id").
		Scan(ctx, &conns)
	return
}
