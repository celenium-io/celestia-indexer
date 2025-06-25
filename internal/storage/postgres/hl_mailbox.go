// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

type HLMailbox struct {
	*database.Bun
}

func NewHLMailbox(conn *database.Bun) *HLMailbox {
	return &HLMailbox{conn}
}

func (hl *HLMailbox) List(ctx context.Context, limit, offset int) (mailbox []storage.HLMailbox, err error) {
	query := hl.DB().NewSelect().
		Model((*storage.HLMailbox)(nil))

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	err = hl.DB().NewSelect().
		TableExpr("(?) as mailbox", query).
		ColumnExpr("mailbox.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as owner__address").
		ColumnExpr("celestial.id as owner__celestials__id, celestial.image_url as owner__celestials__image_url").
		Join("left join tx on mailbox.tx_id = tx.id").
		Join("left join address on address.id = mailbox.owner_id").
		Join("left join celestial on celestial.address_id = mailbox.owner_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &mailbox)
	return
}

func (hl *HLMailbox) ByHash(ctx context.Context, hash []byte) (mailbox storage.HLMailbox, err error) {
	query := hl.DB().NewSelect().
		Model(&mailbox).
		Where("mailbox = ?", hash).
		Limit(1)

	err = hl.DB().NewSelect().
		TableExpr("(?) as mailbox", query).
		ColumnExpr("mailbox.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as owner__address").
		ColumnExpr("celestial.id as owner__celestials__id, celestial.image_url as owner__celestials__image_url").
		Join("left join tx on mailbox.tx_id = tx.id").
		Join("left join address on address.id = mailbox.owner_id").
		Join("left join celestial on celestial.address_id = mailbox.owner_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &mailbox)
	return
}
