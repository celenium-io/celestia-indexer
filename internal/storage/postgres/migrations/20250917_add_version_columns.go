// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddVersionColumns, downAddVersionColumns)
}

func upAddVersionColumns(ctx context.Context, db *bun.DB) error {
	// Create validator's version column
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."validator" ADD "version" int8 NULL`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."validator" ALTER COLUMN "version" SET STORAGE PLAIN`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public."validator"."version" IS 'Signal version'`); err != nil {
		return err
	}

	// Set current validator versions
	versionsQuery := db.NewSelect().Model((*storage.SignalVersion)(nil)).
		ColumnExpr("max(version) as version, validator_id").
		Group("validator_id")

	if _, err := db.NewUpdate().
		With("versions", versionsQuery).
		Table("validator", "versions").
		SetColumn("version", "versions.version").
		Where("validator.id = validator_id").
		Exec(ctx); err != nil {
		return err
	}

	// Create version column at the state
	if _, err := db.ExecContext(ctx, `ALTER TABLE public.state ADD "version" int8 NULL`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public.state ALTER COLUMN "version" SET STORAGE PLAIN`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public.state."version" IS 'Version'`); err != nil {
		return err
	}

	// Set current version
	var currentVersion uint64
	if err := db.NewSelect().Model((*storage.Upgrade)(nil)).
		ColumnExpr("max(version)").
		Scan(ctx, &currentVersion); err != nil {
		return err
	}

	var state storage.State
	if err := db.NewSelect().Model(&state).Where("name = 'celestia_indexer'").Scan(ctx); err != nil {
		return err
	}

	if _, err := db.NewUpdate().Model(&state).WherePK().Set("version = ?", currentVersion).Exec(ctx); err != nil {
		return err
	}

	return nil
}
func downAddVersionColumns(ctx context.Context, db *bun.DB) error {
	return nil
}
