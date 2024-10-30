// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Message -
type Message struct {
	*postgres.Table[*storage.Message]
}

// NewMessage -
func NewMessage(db *database.Bun) *Message {
	return &Message{
		Table: postgres.NewTable[*storage.Message](db),
	}
}

// ByTxId -
func (m *Message) ByTxId(ctx context.Context, txId uint64, limit, offset int) (messages []storage.Message, err error) {
	query := m.DB().NewSelect().Model(&messages).
		Where("tx_id = ?", txId).
		Order("id asc")

	query = limitScope(query, limit)

	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}

func (m *Message) ListWithTx(ctx context.Context, filters storage.MessageListWithTxFilters) (msgs []storage.MessageWithTx, err error) {
	query := m.DB().NewSelect().Model(&msgs).Offset(filters.Offset)
	query = messagesFilter(query, filters)

	err = query.Relation("Tx").Scan(ctx)
	return
}

func (m *Message) ByAddress(ctx context.Context, addressId uint64, filters storage.AddressMsgsFilter) (msgs []storage.AddressMessageWithTx, err error) {
	query := m.DB().NewSelect().Model((*storage.MsgAddress)(nil)).
		Where("address_id = ?", addressId).
		Offset(filters.Offset)

	query = limitScope(query, filters.Limit)
	query = sortScope(query, "msg_id", filters.Sort)

	wrapQuery := m.DB().NewSelect().TableExpr("(?) as msg_address", query).
		ColumnExpr(`msg_address.address_id, msg_address.msg_id, msg_address.type, msg.id AS msg__id, msg.height AS msg__height, msg.time AS msg__time, msg.position AS msg__position, msg.type AS msg__type, msg.tx_id AS msg__tx_id, msg.size AS msg__size, msg.data AS msg__data`).
		ColumnExpr("tx.messages_count as tx__messages_count, tx.fee as tx__fee, tx.status as tx__status, tx.hash as tx__hash, tx.message_types as tx__message_types").
		Join("left join message as msg on msg_address.msg_id = msg.id").
		Join("left join tx on tx.id = msg.tx_id and tx.time = msg.time")
	if len(filters.MessageTypes) > 0 {
		wrapQuery = wrapQuery.Where("msg.type IN (?)", bun.In(filters.MessageTypes))
	}
	wrapQuery = sortScope(wrapQuery, "msg_id", filters.Sort)
	err = wrapQuery.Scan(ctx, &msgs)
	return
}
