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
	Migrations.MustRegister(upSkipUnfinishedUpgrades, downSkipUnfinishedUpgrades)
}

// Separate from the enum migration: a new enum value is usable only after it is committed.
func upSkipUnfinishedUpgrades(ctx context.Context, db *bun.DB) error {
	// versions below the last applied one can no longer be reached
	_, err := db.ExecContext(ctx, `UPDATE upgrade
		SET status = ?
		WHERE status != ?
		AND version < (SELECT max(version) FROM upgrade WHERE status = ?)`,
		types.UpgradeStatusSkipped.String(),
		types.UpgradeStatusApplied.String(),
		types.UpgradeStatusApplied.String(),
	)
	return errors.Wrap(err, "skip unfinished upgrades")
}

func downSkipUnfinishedUpgrades(ctx context.Context, db *bun.DB) error {
	return restoreSkippedUpgrades(ctx, db)
}
