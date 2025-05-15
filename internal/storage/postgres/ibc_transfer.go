// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
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
		Join("left join tx on tx_id = tx.id").
		Join("left join address as receiver on receiver.id = receiver_id").
		Join("left join celestial as cel_receiver on cel_receiver.address_id = receiver_id and cel_receiver.status = 'PRIMARY'").
		Join("left join address as sender on sender.id = sender_id").
		Join("left join celestial as cel_sender on cel_sender.address_id = sender_id and cel_sender.status = 'PRIMARY'").
		Scan(ctx, &transfers)
	return
}
