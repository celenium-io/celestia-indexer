// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddShareVersionColumns, downAddShareVersionColumns)
}

func upAddShareVersionColumns(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."blob_log" ADD "share_version" int4 DEFAULT 0`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."blob_log" ALTER COLUMN "share_version" SET STORAGE PLAIN`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public."blob_log"."share_version" IS 'Share version'`); err != nil {
		return err
	}

	return nil
}
func downAddShareVersionColumns(ctx context.Context, db *bun.DB) error {
	return nil
}
