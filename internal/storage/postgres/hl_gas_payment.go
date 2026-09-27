// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

type HLGasPayment struct {
	db bun.IDB
}

func NewHLGasPayment(db bun.IDB) storage.IHLGasPayment {
	return &HLGasPayment{db}
}

func (hl *HLGasPayment) List(ctx context.Context, limit, offset int) (payments []storage.HLGasPayment, err error) {
	query := hl.db.NewSelect().
		Model((*storage.HLGasPayment)(nil))

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx, &payments)

	return
}
