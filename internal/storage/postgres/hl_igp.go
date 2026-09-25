// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

type HLIGP struct {
	db bun.IDB
}

func NewHLIGP(db bun.IDB) storage.IHLIGP {
	return &HLIGP{db}
}

func (hl *HLIGP) List(ctx context.Context, limit, offset int) (igp []storage.HLIGP, err error) {
	query := hl.db.NewSelect().
		Model(&igp)

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Relation("Configs").
		Scan(ctx)
	return
}

func (hl *HLIGP) ByHash(ctx context.Context, hash []byte) (igp storage.HLIGP, err error) {
	query := hl.db.NewSelect().
		Model(&igp).
		Where("igp_id = ?", hash).
		Limit(1)

	err = query.Relation("Configs").
		Scan(ctx)
	return
}

// IdByHash resolves an IGP id without loading its configs.
func (hl *HLIGP) IdByHash(ctx context.Context, hash []byte) (id uint64, err error) {
	err = hl.db.NewSelect().
		Model((*storage.HLIGP)(nil)).
		Column("id").
		Where("igp_id = ?", hash).
		Limit(1).
		Scan(ctx, &id)
	return
}
