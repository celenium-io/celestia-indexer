// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

type HLMailbox struct {
	db bun.IDB
}

func NewHLMailbox(db bun.IDB) storage.IHLMailbox {
	return &HLMailbox{db}
}

func (hl *HLMailbox) List(ctx context.Context, limit, offset int) (mailbox []storage.HLMailbox, err error) {
	query := hl.db.NewSelect().
		Model((*storage.HLMailbox)(nil))

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	err = hl.db.NewSelect().
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
	query := hl.db.NewSelect().
		Model(&mailbox).
		Where("mailbox = ?", hash).
		Limit(1)

	err = hl.db.NewSelect().
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

func (hl *HLMailbox) ByInternalId(ctx context.Context, internalId uint64) (mailbox storage.HLMailbox, err error) {
	err = hl.db.NewSelect().Model(&mailbox).
		Where("internal_id = ?", internalId).
		Column("id").
		Scan(ctx)
	return
}
