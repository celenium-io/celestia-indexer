// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

// HLToken -
type HLToken struct {
	db bun.IDB
}

// NewHLToken -
func NewHLToken(db bun.IDB) storage.IHLToken {
	return &HLToken{
		db: db,
	}
}

func (t *HLToken) ByHash(ctx context.Context, id []byte) (token storage.HLToken, err error) {
	query := t.db.NewSelect().
		Model((*storage.HLToken)(nil)).
		Where("token_id = ?", id).
		Limit(1)

	err = t.db.NewSelect().
		TableExpr("(?) as token", query).
		ColumnExpr("token.*").
		ColumnExpr("hl_mailbox.mailbox as mailbox__mailbox").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as owner__address").
		ColumnExpr("celestial.id as owner__celestials__id, celestial.image_url as owner__celestials__image_url").
		Join("left join hl_mailbox on mailbox_id = hl_mailbox.id").
		Join("left join tx on token.tx_id = tx.id").
		Join("left join address on address.id = token.owner_id").
		Join("left join celestial on celestial.address_id = token.owner_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &token)
	return
}

func (t *HLToken) List(ctx context.Context, filters storage.ListHyperlaneTokens) (tokens []storage.HLToken, err error) {
	query := t.db.NewSelect().
		Model((*storage.HLToken)(nil))

	query = limitScope(query, filters.Limit)
	query = sortScope(query, "id", filters.Sort)

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	if filters.MailboxId > 0 {
		query = query.Where("mailbox_id = ?", filters.MailboxId)
	}
	if filters.OwnerId > 0 {
		query = query.Where("owner_id = ?", filters.OwnerId)
	}
	if len(filters.Type) > 0 {
		query = query.Where("type IN ?", bun.Tuple(filters.Type))
	}

	err = t.db.NewSelect().
		TableExpr("(?) as token", query).
		ColumnExpr("token.*").
		ColumnExpr("hl_mailbox.mailbox as mailbox__mailbox").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as owner__address").
		ColumnExpr("celestial.id as owner__celestials__id, celestial.image_url as owner__celestials__image_url").
		Join("left join hl_mailbox on mailbox_id = hl_mailbox.id").
		Join("left join tx on token.tx_id = tx.id").
		Join("left join address on address.id = token.owner_id").
		Join("left join celestial on celestial.address_id = token.owner_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &tokens)
	return
}

func (t *HLToken) IdByTokenId(ctx context.Context, tokenId []byte) (id uint64, err error) {
	err = t.db.NewSelect().Model((*storage.HLToken)(nil)).
		Where("token_id = ?", tokenId).
		Column("id").
		Scan(ctx, &id)
	return
}
