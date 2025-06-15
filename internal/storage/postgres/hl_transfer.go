// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/uptrace/bun"
)

// HLTransfer -
type HLTransfer struct {
	*database.Bun
}

// NewHLTransfer -
func NewHLTransfer(db *database.Bun) *HLTransfer {
	return &HLTransfer{
		Bun: db,
	}
}

func (t *HLTransfer) List(ctx context.Context, filters storage.ListHyperlaneTransfers) (transfers []storage.HLTransfer, err error) {
	query := t.DB().NewSelect().
		Model((*storage.HLTransfer)(nil))

	query = limitScope(query, filters.Limit)
	query = sortScope(query, "id", filters.Sort)

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	if filters.MailboxId > 0 {
		query = query.Where("mailbox_id = ?", filters.MailboxId)
	}
	if filters.AddressId > 0 {
		query = query.Where("address_id = ?", filters.AddressId)
	}
	if filters.RelayerId > 0 {
		query = query.Where("relayer_id = ?", filters.RelayerId)
	}
	if filters.TokenId > 0 {
		query = query.Where("token_id = ?", filters.TokenId)
	}
	if filters.Domain > 0 {
		query = query.Where("counterparty = ?", filters.Domain)
	}
	if len(filters.Type) > 0 {
		query = query.Where("type IN (?)", bun.In(filters.Type))
	}

	err = t.DB().NewSelect().
		TableExpr("(?) as transfer", query).
		ColumnExpr("transfer.*").
		ColumnExpr("hl_mailbox.mailbox as mailbox__mailbox").
		ColumnExpr("hl_token.token_id as token__token_id").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as address__address").
		ColumnExpr("celestial.id as address__celestials__id, celestial.image_url as address__celestials__image_url").
		ColumnExpr("relayer.address as relayer__address").
		ColumnExpr("relayer_celestials.id as relayer__celestials__id, relayer_celestials.image_url as relayer__celestials__image_url").
		Join("left join hl_mailbox on mailbox_id = hl_mailbox.id").
		Join("left join hl_token on transfer.token_id = hl_token.id").
		Join("left join tx on transfer.tx_id = tx.id").
		Join("left join address on address.id = transfer.address_id").
		Join("left join celestial on celestial.address_id = transfer.address_id and celestial.status = 'PRIMARY'").
		Join("left join address as relayer on relayer.id = transfer.relayer_id").
		Join("left join celestial as relayer_celestials on relayer_celestials.address_id = transfer.relayer_id and relayer_celestials.status = 'PRIMARY'").
		Scan(ctx, &transfers)
	return
}
