// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upForwarding, downForwarding)
}

func upForwarding(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."address" ADD is_forwarding BOOLEAN NOT NULL DEFAULT FALSE;`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."address" ALTER COLUMN "is_forwarding" SET STORAGE PLAIN`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public."address"."is_forwarding" IS 'Is the address used for forwarding'`); err != nil {
		return err
	}

	return nil
}
func downForwarding(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE public."address" DROP COLUMN IF EXISTS "is_forwarding";`)
	return err
}
