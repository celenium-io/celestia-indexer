// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddUpgradeStatusColumns, downAddUpgradeStatusColumns)
}

func upAddUpgradeStatusColumns(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."upgrade" ADD COLUMN IF NOT EXISTS "status" upgrade_status DEFAULT 'processing'`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."upgrade" ALTER COLUMN "status" SET STORAGE PLAIN`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public."upgrade"."status" IS 'Upgrade status'`); err != nil {
		return err
	}

	var state *storage.State
	err := db.NewSelect().
		Model(&state).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewUpdate().
		Model((*storage.Upgrade)(nil)).
		Set("status = ?", types.UpgradeStatusApplied).
		Where("version <= ?", state.Version).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewUpdate().
		Model((*storage.Upgrade)(nil)).
		Set("status = ?", types.UpgradeStatusWaitingUpgrade).
		Where("version > ?", state.Version).
		Where("end_height is not null").
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
func downAddUpgradeStatusColumns(ctx context.Context, db *bun.DB) error {
	return nil
}
