// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// VestingPeriod -
type VestingPeriod struct {
	*postgres.Table[*storage.VestingPeriod]
}

// NewVestingPeriod -
func NewVestingPeriod(db bun.IDB) storage.IVestingPeriod {
	return &VestingPeriod{
		Table: postgres.NewTable[*storage.VestingPeriod](db),
	}
}

func (v *VestingPeriod) ByVesting(ctx context.Context, id uint64, limit, offset int) (periods []storage.VestingPeriod, err error) {
	query := v.DB().NewSelect().
		Model(&periods).
		Where("vesting_account_id = ?", id).
		Order("id desc")

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}
