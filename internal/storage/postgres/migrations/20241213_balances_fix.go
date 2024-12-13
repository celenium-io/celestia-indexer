// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upBalancesFix, downBalancesFix)
}

func upBalancesFix(ctx context.Context, db *bun.DB) error {
	var balances []*storage.Balance
	if err := db.NewSelect().Model(&balances).Where("delegated < 0").Scan(ctx); err != nil {
		return err
	}

	for i := range balances {
		var logs []*storage.StakingLog
		if err := db.NewSelect().Model(&logs).Where("address_id = ?", balances[i].Id).Scan(ctx); err != nil {
			return err
		}

		total := decimal.Zero
		for j := range logs {
			total = total.Add(logs[j].Change)
		}

		_, err := db.NewUpdate().
			Model(balances[i]).
			WherePK().
			Set("delegated = ?", total).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func downBalancesFix(_ context.Context, _ *bun.DB) error {
	return nil
}
