// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

type HLIGP struct {
	*database.Bun
}

func NewHLIGP(conn *database.Bun) *HLIGP {
	return &HLIGP{conn}
}

func (hl *HLIGP) List(ctx context.Context, limit, offset int) (igp []storage.HLIGP, err error) {
	query := hl.DB().NewSelect().
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
	query := hl.DB().NewSelect().
		Model(&igp).
		Where("igp_id = ?", hash).
		Limit(1)

	err = query.Relation("Configs").
		Scan(ctx)
	return
}
