// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Jail -
type Jail struct {
	*postgres.Table[*storage.Jail]
}

// NewJail -
func NewJail(db *database.Bun) *Jail {
	return &Jail{
		Table: postgres.NewTable[*storage.Jail](db),
	}
}

func (j *Jail) ByValidator(ctx context.Context, id uint64, limit, offset int) (jails []storage.Jail, err error) {
	query := j.DB().NewSelect().Model(&jails).
		Where("validator_id = ?", id).
		Order("time desc")

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}
