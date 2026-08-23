// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddModuleMinfee, downAddModuleMinfee)
}

func upAddModuleMinfee(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		ALTER TYPE module_name ADD VALUE IF NOT EXISTS 'minfee'
	`)
	return err
}

func downAddModuleMinfee(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM pg_enum 
		WHERE enumlabel = 'minfee' AND enumtypid = (
			SELECT oid FROM pg_type WHERE typname = 'module_name'
		)
	`)
	return err
}
