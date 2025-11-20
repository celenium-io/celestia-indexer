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
	if _, err := db.ExecContext(
		ctx,
		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = ?) THEN
				CREATE TYPE ? AS ENUM (?);
			END IF;
		END$$`,
		"upgrade_status",
		bun.Safe("upgrade_status"),
		bun.In(types.UpgradeStatusValues()),
	); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `ALTER TABLE public."upgrade" ADD COLUMN IF NOT EXISTS "status" upgrade_status DEFAULT 'processing'`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."upgrade" ALTER COLUMN "status" SET STORAGE PLAIN`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public."upgrade"."status" IS 'Upgrade status'`); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `ALTER TABLE public."upgrade" ADD COLUMN IF NOT EXISTS "applied_at" timestamptz`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."upgrade" ALTER COLUMN "applied_at" SET STORAGE PLAIN`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public."upgrade"."applied_at" IS 'The time when upgrade was applied'`); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `ALTER TABLE public."upgrade" ADD COLUMN IF NOT EXISTS "applied_at_level" int8 DEFAULT 0`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."upgrade" ALTER COLUMN "applied_at_level" SET STORAGE PLAIN`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public."upgrade"."applied_at_level" IS 'The level when upgrade was applied'`); err != nil {
		return err
	}

	var state storage.State
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

	var blocks []storage.Block
	err = db.NewSelect().
		Model(&blocks).
		ColumnExpr("min(height) as height, min(time) as time, version_app").
		Where("version_app IS NOT NULL").
		Group("version_app").
		Scan(ctx)
	if err != nil {
		return err
	}

	for _, block := range blocks {
		_, err = db.NewUpdate().
			Model((*storage.Upgrade)(nil)).
			Set("applied_at = ?", block.Time).
			Set("applied_at_level = ?", block.Height).
			Where("version = ?", block.VersionApp).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
func downAddUpgradeStatusColumns(ctx context.Context, db *bun.DB) error {
	return nil
}
