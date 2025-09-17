// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddMessagesCount, downMessagesCount)
}

func upAddMessagesCount(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."validator" ADD "messages_count" int8 NULL`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."validator" ALTER COLUMN "messages_count" SET STORAGE PLAIN`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public."validator"."messages_count" IS 'Count of validator messages'`); err != nil {
		return err
	}
	return nil
}
func downMessagesCount(ctx context.Context, db *bun.DB) error {
	return nil
}
