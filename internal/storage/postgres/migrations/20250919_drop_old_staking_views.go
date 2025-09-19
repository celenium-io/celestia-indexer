// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upDropOldStakingViews, downDropOldStakingViews)
}

func upDropOldStakingViews(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `DROP MATERIALIZED VIEW IF EXISTS staking_by_hour CASCADE`); err != nil {
		return err
	}
	return nil
}
func downDropOldStakingViews(ctx context.Context, db *bun.DB) error {
	return nil
}
