// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type IbcTransfer struct {
	*database.Bun
}

func NewIbcTransfer(conn *database.Bun) *IbcTransfer {
	return &IbcTransfer{conn}
}

func (c *IbcTransfer) List(ctx context.Context, fltrs storage.ListIbcTransferFilters) (transfers []storage.IbcTransfer, err error) {
	query := c.DB().NewSelect().
		Model((*storage.IbcTransfer)(nil))

	if fltrs.Offset > 0 {
		query.Offset(fltrs.Offset)
	}

	query = limitScope(query, fltrs.Limit)
	query = query.OrderExpr("time ?0, id ?0", bun.Safe(fltrs.Sort))

	if fltrs.ChannelId != "" {
		query = query.Where("channel_id = ?", fltrs.ChannelId)
	}

	if fltrs.AddressId != nil {
		query = query.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("receiver_id = ?", *fltrs.AddressId).WhereOr("sender_id = ?", *fltrs.AddressId)
		})
	}
	if fltrs.ReceiverId != nil {
		query = query.Where("receiver_id = ?", *fltrs.ReceiverId)
	}
	if fltrs.SenderId != nil {
		query = query.Where("sender_id = ?", *fltrs.SenderId)
	}

	err = c.DB().NewSelect().
		TableExpr("(?) as ibc_transfer", query).
		ColumnExpr("ibc_transfer.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("receiver.address as receiver__address").
		ColumnExpr("cel_receiver.id as receiver__celestials__id, cel_receiver.image_url as receiver__celestials__image_url").
		ColumnExpr("sender.address as sender__address").
		ColumnExpr("cel_sender.id as sender__celestials__id, cel_sender.image_url as sender__celestials__image_url").
		ColumnExpr("ibc_client.chain_id as connection__client__chain_id").
		Join("left join tx on tx_id = tx.id").
		Join("left join address as receiver on receiver.id = receiver_id").
		Join("left join celestial as cel_receiver on cel_receiver.address_id = receiver_id and cel_receiver.status = 'PRIMARY'").
		Join("left join address as sender on sender.id = sender_id").
		Join("left join celestial as cel_sender on cel_sender.address_id = sender_id and cel_sender.status = 'PRIMARY'").
		Join("left join ibc_connection on ibc_connection.connection_id = ibc_transfer.connection_id").
		Join("left join ibc_client on ibc_connection.client_id = ibc_client.id").
		Scan(ctx, &transfers)
	return
}

func (c *IbcTransfer) Series(ctx context.Context, channelId string, timeframe storage.Timeframe, column string, req storage.SeriesRequest) (items []storage.HistogramItem, err error) {
	query := c.DB().NewSelect().
		Order("time desc").
		Where("channel_id = ?", channelId)

	switch timeframe {
	case storage.TimeframeHour:
		query = query.Table(storage.ViewIbcTransfersByHour)
	case storage.TimeframeDay:
		query = query.Table(storage.ViewIbcTransfersByDay)
	case storage.TimeframeMonth:
		query = query.Table(storage.ViewIbcTransfersByMonth)
	default:
		return nil, errors.Errorf("invalid timeframe: %s", timeframe)
	}

	switch column {
	case "count":
		query = query.ColumnExpr("count as value, time as bucket")
	case "amount":
		query = query.ColumnExpr("amount as value, time as bucket")
	default:
		return nil, errors.Errorf("invalid column: %s", column)
	}

	if !req.From.IsZero() {
		query = query.Where("time >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("time < ?", req.To)
	}

	err = query.Scan(ctx, &items)
	return
}
