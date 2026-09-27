// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

// ValidatorBondUpdate -
type ValidatorBondUpdate struct {
	db bun.IDB
}

// NewValidatorBondUpdate -
func NewValidatorBondUpdate(db bun.IDB) storage.IValidatorBondUpdate {
	return &ValidatorBondUpdate{db}
}

// ListByValidator -
func (bu *ValidatorBondUpdate) ListByValidator(
	ctx context.Context, validatorId uint64, filters storage.FilterBondUpdatesListByValidator,
) (updates []storage.ValidatorBondUpdate, err error) {
	query := bu.db.NewSelect().Model(&updates).
		Where("validator_id = ?", validatorId)

	query = limitScope(query, filters.Limit)
	query = timeAndIdSort(query, filters.Sort)
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err = query.Scan(ctx)
	return
}

// LastBondUpdate -
func (bu *ValidatorBondUpdate) LastBondUpdate(ctx context.Context, validatorId uint64) (update storage.ValidatorBondUpdate, err error) {
	err = bu.db.NewSelect().Model(&update).
		Where("validator_id = ?", validatorId).
		OrderExpr("time desc, id desc").
		Limit(1).
		Scan(ctx)
	return
}
