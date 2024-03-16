// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// StakingLog -
type StakingLog struct {
	*postgres.Table[*storage.StakingLog]
}

// NewStakingLog -
func NewStakingLog(db *database.Bun) *StakingLog {
	return &StakingLog{
		Table: postgres.NewTable[*storage.StakingLog](db),
	}
}

func (d *StakingLog) ByValidator(ctx context.Context, id uint64, limit, offset int) (logs []storage.StakingLog, err error) {
	query := d.DB().NewSelect().
		Model(&logs).
		Where("validator_id = ?", id)

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx)
	return
}
