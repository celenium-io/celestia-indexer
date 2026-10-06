// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddUpgradeStatusSkipped, downAddUpgradeStatusSkipped)
}

func upAddUpgradeStatusSkipped(ctx context.Context, db *bun.DB) error {
	// ALTER TYPE ... ADD VALUE cannot run inside a transaction block
	_, err := db.ExecContext(ctx, `ALTER TYPE upgrade_status ADD VALUE IF NOT EXISTS ?`,
		types.UpgradeStatusSkipped.String(),
	)
	return errors.Wrap(err, "add upgrade_status")
}

func downAddUpgradeStatusSkipped(ctx context.Context, db *bun.DB) error {
	// the enum value can be dropped only when no row uses it
	if err := restoreSkippedUpgrades(ctx, db); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, `DELETE FROM pg_enum
		WHERE enumlabel = ?
		AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'upgrade_status')`,
		types.UpgradeStatusSkipped.String(),
	)
	return errors.Wrap(err, "remove upgrade_status")
}

// restoreSkippedUpgrades puts back the status the upgrades had before they were skipped.
func restoreSkippedUpgrades(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `UPDATE upgrade
		SET status = CASE WHEN end_height > 0 THEN ?::upgrade_status ELSE ?::upgrade_status END
		WHERE status::text = ?`,
		types.UpgradeStatusWaitingUpgrade.String(),
		types.UpgradeStatusProcessing.String(),
		types.UpgradeStatusSkipped.String(),
	)
	return errors.Wrap(err, "restore skipped upgrades")
}
