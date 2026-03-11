// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upDropAutoincrement, downDropAutoincrement)
}

func upDropAutoincrement(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."tx" ALTER COLUMN "id" DROP DEFAULT;`); err != nil {
		return errors.Wrap(err, "drop default")
	}
	if _, err := db.ExecContext(ctx, `DROP sequence tx_id_seq;`); err != nil {
		return errors.Wrap(err, "drop sequence tx_id_seq")
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."message" ALTER COLUMN "id" DROP DEFAULT;`); err != nil {
		return errors.Wrap(err, "drop default")
	}
	if _, err := db.ExecContext(ctx, `DROP sequence message_id_seq;`); err != nil {
		return errors.Wrap(err, "drop sequence message_id_seq")
	}
	return nil
}

func downDropAutoincrement(_ context.Context, _ *bun.DB) error {
	return nil
}
