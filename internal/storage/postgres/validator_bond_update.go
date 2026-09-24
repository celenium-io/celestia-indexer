// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/go-lib/database"
)

// ValidatorBondUpdate -
type ValidatorBondUpdate struct {
	db *database.Bun
}

// NewValidatorBondUpdate -
func NewValidatorBondUpdate(db *database.Bun) *ValidatorBondUpdate {
	return &ValidatorBondUpdate{db}
}

// ListByValidator -
func (bu *ValidatorBondUpdate) ListByValidator(
	ctx context.Context, validatorId uint64, filters storage.FilterBondUpdatesListByValidator,
) (updates []storage.ValidatorBondUpdate, err error) {
	query := bu.db.DB().NewSelect().Model(&updates).
		Where("validator_id = ?", validatorId)

	query = limitScope(query, filters.Limit)
	query = timeAndIdSort(query, filters.Sort)
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err = query.Scan(ctx)
	return
}
