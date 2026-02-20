// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type ZkISM struct {
	*database.Bun
}

func NewZkISM(conn *database.Bun) *ZkISM {
	return &ZkISM{conn}
}

func (z *ZkISM) List(ctx context.Context, filter storage.ZkISMFilter) (items []storage.ZkISM, err error) {
	if filter.Sort == "" {
		filter.Sort = sdk.SortOrderDesc
	}

	query := z.DB().NewSelect().
		Model((*storage.ZkISM)(nil))

	if filter.CreatorId != nil {
		query = query.Where("creator_id = ?", *filter.CreatorId)
	}
	if filter.TxId != nil {
		query = query.Where("tx_id = ?", *filter.TxId)
	}

	query = query.OrderExpr("time ?0, id ?0", bun.Safe(filter.Sort))
	query = limitScope(query, filter.Limit)
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err = z.DB().NewSelect().
		TableExpr("(?) as zk_ism", query).
		ColumnExpr("zk_ism.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as creator__address").
		ColumnExpr("celestial.id as creator__celestials__id, celestial.image_url as creator__celestials__image_url").
		Join("left join tx on zk_ism.tx_id = tx.id").
		Join("left join address on address.id = zk_ism.creator_id").
		Join("left join celestial on celestial.address_id = zk_ism.creator_id and celestial.status = 'PRIMARY'").
		OrderExpr("zk_ism.time ?0, zk_ism.id ?0", bun.Safe(filter.Sort)).
		Scan(ctx, &items)
	return
}

func (z *ZkISM) ById(ctx context.Context, id uint64) (item storage.ZkISM, err error) {
	query := z.DB().NewSelect().
		Model(&item).
		Where("zk_ism.id = ?", id).
		Limit(1)

	err = z.DB().NewSelect().
		TableExpr("(?) as zk_ism", query).
		ColumnExpr("zk_ism.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as creator__address").
		ColumnExpr("celestial.id as creator__celestials__id, celestial.image_url as creator__celestials__image_url").
		Join("left join tx on zk_ism.tx_id = tx.id").
		Join("left join address on address.id = zk_ism.creator_id").
		Join("left join celestial on celestial.address_id = zk_ism.creator_id and celestial.status = 'PRIMARY'").
		Scan(ctx, &item)
	return
}

func (z *ZkISM) Updates(ctx context.Context, id uint64, filter storage.ZkISMUpdatesFilter) (items []storage.ZkISMUpdate, err error) {
	if filter.Sort == "" {
		filter.Sort = sdk.SortOrderDesc
	}

	query := z.DB().NewSelect().
		Model((*storage.ZkISMUpdate)(nil)).
		Where("zk_ism_id = ?", id)

	if filter.SignerId != nil {
		query = query.Where("signer_id = ?", *filter.SignerId)
	}
	if filter.TxId != nil {
		query = query.Where("tx_id = ?", *filter.TxId)
	}
	if !filter.From.IsZero() {
		query = query.Where("time >= ?", filter.From)
	}
	if !filter.To.IsZero() {
		query = query.Where("time < ?", filter.To)
	}

	query = query.OrderExpr("time ?0, id ?0", bun.Safe(filter.Sort))
	query = limitScope(query, filter.Limit)
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err = z.DB().NewSelect().
		TableExpr("(?) as u", query).
		ColumnExpr("u.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as signer__address").
		ColumnExpr("celestial.id as signer__celestials__id, celestial.image_url as signer__celestials__image_url").
		Join("left join tx on u.tx_id = tx.id").
		Join("left join address on address.id = u.signer_id").
		Join("left join celestial on celestial.address_id = u.signer_id and celestial.status = 'PRIMARY'").
		OrderExpr("u.time ?0, u.id ?0", bun.Safe(filter.Sort)).
		Scan(ctx, &items)
	return
}

func (z *ZkISM) Messages(ctx context.Context, id uint64, filter storage.ZkISMUpdatesFilter) (items []storage.ZkISMMessage, err error) {
	if filter.Sort == "" {
		filter.Sort = sdk.SortOrderDesc
	}

	query := z.DB().NewSelect().
		Model((*storage.ZkISMMessage)(nil)).
		Where("zk_ism_id = ?", id)

	if filter.SignerId != nil {
		query = query.Where("signer_id = ?", *filter.SignerId)
	}
	if filter.TxId != nil {
		query = query.Where("tx_id = ?", *filter.TxId)
	}
	if !filter.From.IsZero() {
		query = query.Where("time >= ?", filter.From)
	}
	if !filter.To.IsZero() {
		query = query.Where("time < ?", filter.To)
	}

	query = query.OrderExpr("time ?0, id ?0", bun.Safe(filter.Sort))
	query = limitScope(query, filter.Limit)
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err = z.DB().NewSelect().
		TableExpr("(?) as m", query).
		ColumnExpr("m.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("address.address as signer__address").
		ColumnExpr("celestial.id as signer__celestials__id, celestial.image_url as signer__celestials__image_url").
		Join("left join tx on m.tx_id = tx.id").
		Join("left join address on address.id = m.signer_id").
		Join("left join celestial on celestial.address_id = m.signer_id and celestial.status = 'PRIMARY'").
		OrderExpr("m.time ?0, m.id ?0", bun.Safe(filter.Sort)).
		Scan(ctx, &items)
	return
}
