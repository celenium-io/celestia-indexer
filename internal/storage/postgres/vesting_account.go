// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// VestingAccount -
type VestingAccount struct {
	*postgres.Table[*storage.VestingAccount]
}

// NewVestingAccount -
func NewVestingAccount(db *database.Bun) *VestingAccount {
	return &VestingAccount{
		Table: postgres.NewTable[*storage.VestingAccount](db),
	}
}

func (v *VestingAccount) ByAddress(ctx context.Context, addressId uint64, limit, offset int, showEnded bool) (accs []storage.VestingAccount, err error) {
	query := v.DB().NewSelect().
		Model((*storage.VestingAccount)(nil)).
		Where("address_id = ?", addressId).
		Order("end_time desc")

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}
	if !showEnded {
		query = query.Where("end_time >= ?", time.Now().UTC())
	}

	err = v.DB().NewSelect().
		TableExpr("(?) as vesting_account", query).
		ColumnExpr("vesting_account.*").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join tx on tx.id = tx_id").
		Scan(ctx, &accs)
	return
}
