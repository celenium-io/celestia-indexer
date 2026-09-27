// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

type IbcConnection struct {
	db bun.IDB
}

func NewIbcConnection(db bun.IDB) storage.IIbcConnection {
	return &IbcConnection{db}
}

func (c *IbcConnection) ById(ctx context.Context, id string) (conn storage.IbcConnection, err error) {
	query := c.db.NewSelect().
		Model((*storage.IbcConnection)(nil)).
		Where("connection_id = ?", id).
		Limit(1)

	err = c.db.NewSelect().
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
	query := c.db.NewSelect().
		Model(&conns)

	if fltrs.Offset > 0 {
		query.Offset(fltrs.Offset)
	}

	query = limitScope(query, fltrs.Limit)
	query = sortScope(query, "height", fltrs.Sort)

	if fltrs.ClientId != "" {
		query = query.Where("client_id = ?", fltrs.ClientId)
	}

	err = c.db.NewSelect().
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

func (c *IbcConnection) IdsByClients(ctx context.Context, clientIds ...string) (res []string, err error) {
	if len(clientIds) == 0 {
		return res, nil
	}

	err = c.db.NewSelect().
		Column("connection_id").
		Model((*storage.IbcConnection)(nil)).
		Where("client_id IN ?", bun.Tuple(clientIds)).
		Scan(ctx, &res)
	return
}

func (c *IbcConnection) ByIds(ctx context.Context, ids ...string) (result []storage.ConnIdAndClientId, err error) {
	if len(ids) == 0 {
		return
	}

	err = c.db.NewSelect().
		Model((*storage.IbcConnection)(nil)).
		Column("connection_id", "client_id").
		Where("connection_id IN ?", bun.Tuple(ids)).
		Scan(ctx, &result)
	return
}
