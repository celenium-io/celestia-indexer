// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Forwarding -
type Forwarding struct {
	*postgres.Table[*storage.Forwarding]
}

// NewForwarding -
func NewForwarding(db *database.Bun) *Forwarding {
	return &Forwarding{
		Table: postgres.NewTable[*storage.Forwarding](db),
	}
}

func (f *Forwarding) ById(ctx context.Context, id uint64) (forwarding storage.Forwarding, prevTime time.Time, err error) {
	subQuery := f.DB().NewSelect().
		Model(&forwarding).
		Where("id = ?", id).
		WhereOr("id = ?", id-1).
		Order("id desc").
		Limit(2)

	var fwds []storage.Forwarding
	err = f.DB().NewSelect().
		TableExpr("(?) as forwarding", subQuery).
		ColumnExpr("forwarding.*").
		ColumnExpr("address.id as address__id, address.address as address__address").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join address on address.id = forwarding.address_id").
		Join("left join tx on tx.id = forwarding.tx_id").
		Scan(ctx, &fwds)
	if err != nil {
		return forwarding, prevTime, err
	}

	if len(fwds) == 0 {
		return forwarding, prevTime, sql.ErrNoRows
	}

	if len(fwds) > 1 {
		prevTime = fwds[1].Time
	}

	forwarding = fwds[0]
	return
}

func (f *Forwarding) Filter(ctx context.Context, filters storage.ForwardingFilter) (forwardings []storage.Forwarding, err error) {
	subQuery := f.DB().NewSelect().Model((*storage.Forwarding)(nil))

	if filters.Height != nil {
		subQuery = subQuery.Where("height = ?", *filters.Height)
	}
	if filters.AddressId != nil {
		subQuery = subQuery.Where("address_id = ?", *filters.AddressId)
	}
	if filters.TxId != nil {
		subQuery = subQuery.Where("tx_id = ?", *filters.TxId)
	}

	if filters.Sort == "" {
		filters.Sort = sdk.SortOrderAsc
	}

	if !filters.From.IsZero() {
		subQuery = subQuery.Where("time >= ?", filters.From)
	}
	if !filters.To.IsZero() {
		subQuery = subQuery.Where("time < ?", filters.To)
	}

	subQuery = subQuery.OrderExpr("time ?0, id ?0", bun.Safe(filters.Sort))
	subQuery = limitScope(subQuery, filters.Limit)
	if filters.Offset > 0 {
		subQuery = subQuery.Offset(filters.Offset)
	}

	err = f.DB().NewSelect().
		TableExpr("(?) as forwarding", subQuery).
		ColumnExpr("forwarding.*").
		ColumnExpr("address.id as address__id, address.address as address__address").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join address on address.id = forwarding.address_id").
		Join("left join tx on tx.id = forwarding.tx_id").
		OrderExpr("forwarding.time ?0, forwarding.id ?0", bun.Safe(filters.Sort)).
		Scan(ctx, &forwardings)

	return
}

func (f *Forwarding) Inputs(ctx context.Context, addressId uint64, from, to time.Time) (inputs []storage.ForwardingInput, err error) {
	messagesQuery := f.DB().NewSelect().
		Model((*storage.MsgAddress)(nil)).
		Where("address_id = ?", addressId).
		Where("type = ?", types.MsgAddressTypeToAddress)

	transfersQuery := f.DB().NewSelect().
		Model((*storage.HLTransfer)(nil)).
		Where("address_id = ?", addressId).
		Where("type = ?", types.HLTransferTypeReceive)
	if !from.IsZero() {
		transfersQuery = transfersQuery.Where("time > ?", from)
	}
	if !to.IsZero() {
		transfersQuery = transfersQuery.Where("time < ?", to)
	}

	subQuery := f.DB().NewSelect().
		Table("address_messages").
		ColumnExpr("message.height as height, message.time as time, tx.hash as hash, 'utia' as denom, '0' as amount, NULL as src, NULL as counterparty, message.data as data").
		Join("left join message on message.id = msg_id").
		Join("left join tx on tx.id = message.tx_id").
		Where("message.type = ?", types.MsgSend)

	if !from.IsZero() {
		subQuery = subQuery.Where("message.time > ?", from)
	}
	if !to.IsZero() {
		subQuery = subQuery.Where("message.time < ?", to)
	}

	subQuery = subQuery.UnionAll(
		f.DB().NewSelect().
			Table("address_transfers").
			ColumnExpr("tx.height, tx.time, tx.hash as hash, denom, amount, counterparty_address as src, counterparty, NULL as data").
			Join("left join tx on tx.id = address_transfers.tx_id"),
	)

	query := f.DB().NewSelect().
		With("address_messages", messagesQuery).
		With("address_transfers", transfersQuery).
		TableExpr("(?) as inputs", subQuery).
		Order("time desc")

	if err = query.Scan(ctx, &inputs); err != nil {
		return nil, err
	}

	for i := range inputs {
		if inputs[i].Data != nil {
			if src, ok := inputs[i].Data["FromAddress"]; ok {
				if srcStr, ok := src.(string); ok {
					inputs[i].From = srcStr
				}
			}

			if amount, ok := inputs[i].Data["Amount"]; ok {
				if arr, ok := amount.([]any); ok && len(arr) > 0 {
					if m, ok := arr[0].(map[string]any); ok {
						if amountStr, ok := m["Amount"].(string); ok {
							inputs[i].Amount = amountStr
						}
					}
				}
			}
		}
	}
	return
}
