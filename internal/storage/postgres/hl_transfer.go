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

func (t *HLTransfer) List(ctx context.Context, filters storage.ListHyperlaneTransferFilters) (transfers []storage.HLTransfer, err error) {
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
	if filters.TxId > 0 {
		query = query.Where("tx_id = ?", filters.TxId)
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

func (t *HLTransfer) ById(ctx context.Context, id uint64) (transfer storage.HLTransfer, err error) {
	query := t.DB().NewSelect().
		Model((*storage.HLTransfer)(nil)).
		Where("id = ?", id)

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
		Scan(ctx, &transfer)
	return
}

func (t *HLTransfer) Series(ctx context.Context, domainId uint64, timeframe storage.Timeframe, column string, req storage.SeriesRequest) (items []storage.HistogramItem, err error) {
	query := t.DB().NewSelect().
		Order("time desc").
		Where("counterparty = ?", domainId)

	switch timeframe {
	case storage.TimeframeHour:
		query = query.Table(storage.ViewHlTransfersByHour)
	case storage.TimeframeDay:
		query = query.Table(storage.ViewHlTransfersByDay)
	case storage.TimeframeMonth:
		query = query.Table(storage.ViewHlTransfersByMonth)
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

func (t *HLTransfer) StatsByDomain(ctx context.Context, limit, offset int) (stats []storage.DomainStats, err error) {
	query := t.DB().NewSelect().
		Table(storage.ViewHlTransfersByMonth).
		ColumnExpr("counterparty as domain_id, sum(count) as tx_count, sum(amount) as amount").
		Group("domain_id").
		Order("amount DESC").
		Offset(offset)

	query = limitScope(query, limit)

	err = query.Scan(ctx, &stats)
	return
}
