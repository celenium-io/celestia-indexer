// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// SignalVersion -
type SignalVersion struct {
	*database.Bun
}

// NewSignalVersion -
func NewSignalVersion(db *database.Bun) *SignalVersion {
	return &SignalVersion{
		Bun: db,
	}
}

func (t *SignalVersion) List(ctx context.Context, filters storage.ListSignalsFilter) (signals []storage.SignalVersion, err error) {
	query := t.DB().NewSelect().
		Model((*storage.SignalVersion)(nil))

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	if !filters.From.IsZero() {
		query = query.Where("time >= ?", filters.From)
	}
	if !filters.To.IsZero() {
		query = query.Where("time < ?", filters.To)
	}

	query = limitScope(query, filters.Limit)
	query = sortScope(query, "height", filters.Sort)

	if filters.TxId != nil {
		query = query.Where("tx_id = ?", *filters.TxId)
	}
	if filters.ValidatorId > 0 {
		query = query.Where("validator_id = ?", filters.ValidatorId)
	}
	if filters.Version > 0 {
		query = query.Where("version = ?", filters.Version)
	}

	err = t.DB().NewSelect().
		TableExpr("(?) as signal_version", query).
		ColumnExpr("signal_version.*").
		ColumnExpr("validator.address as validator__address").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join validator as validator on validator.id = validator_id").
		Join("left join tx on tx_id = tx.id").
		Scan(ctx, &signals)

	return
}
