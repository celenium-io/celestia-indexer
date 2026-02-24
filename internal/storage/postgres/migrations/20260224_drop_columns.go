// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upDropZkismColumns, downDropZkismColumns)
}

func upDropZkismColumns(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."zk_ism" DROP COLUMN IF EXISTS "state_root";`); err != nil {
		return errors.Wrap(err, "drop state root")
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."zk_ism_update" DROP COLUMN IF EXISTS "new_state_root";`); err != nil {
		return errors.Wrap(err, "drop new state root")
	}
	return nil
}

func downDropZkismColumns(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."zk_ism" ADD COLUMN IF NOT EXISTS "state_root" BYTEA;`); err != nil {
		return errors.Wrap(err, "add state root")
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public."zk_ism_update" ADD COLUMN IF NOT EXISTS "new_state_root" BYTEA;`); err != nil {
		return errors.Wrap(err, "add new state root")
	}
	return nil
}
