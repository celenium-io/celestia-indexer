// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"
	"database/sql"

	"github.com/celestiaorg/celestia-app/v10/pkg/appconsts"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upUpgradeExpectedHeight, downUpgradeExpectedHeight)
}

func upUpgradeExpectedHeight(ctx context.Context, db *bun.DB) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(ctx, `ALTER TABLE upgrade ADD COLUMN IF NOT EXISTS expected_upgrade_height bigint`); err != nil {
			return errors.Wrap(err, "add expected_upgrade_height")
		}

		// the chain switches the version in the block after the scheduled height,
		// which is exact regardless of the delay constants of old app versions
		if _, err := tx.ExecContext(ctx, `
			UPDATE upgrade SET expected_upgrade_height = applied_at_level - 1
			WHERE status = 'applied' AND applied_at_level > 0
		`); err != nil {
			return errors.Wrap(err, "fill applied upgrades")
		}

		var chainId string
		err := tx.NewSelect().Table("state").Column("chain_id").Limit(1).Scan(ctx, &chainId)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		case err != nil:
			return errors.Wrap(err, "get chain id")
		}

		if _, err := tx.ExecContext(ctx, `
			UPDATE upgrade SET expected_upgrade_height = end_height + ?
			WHERE status != 'applied' AND end_height > 0
		`, appconsts.GetUpgradeHeightDelay(chainId)); err != nil {
			return errors.Wrap(err, "fill scheduled upgrades")
		}
		return nil
	})
}

func downUpgradeExpectedHeight(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE upgrade DROP COLUMN IF EXISTS expected_upgrade_height`)
	return errors.Wrap(err, "drop expected_upgrade_height")
}
