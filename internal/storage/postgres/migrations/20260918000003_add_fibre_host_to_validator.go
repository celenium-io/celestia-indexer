// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddFibreHostToValidator, downAddFibreHostToValidator)
}

func upAddFibreHostToValidator(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		ALTER TABLE validator
			ADD COLUMN IF NOT EXISTS fibre_host        text,
			ADD COLUMN IF NOT EXISTS fibre_host_height bigint
	`)
	return errors.Wrap(err, "add validator fibre host columns")
}

func downAddFibreHostToValidator(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		ALTER TABLE validator
			DROP COLUMN IF EXISTS fibre_host,
			DROP COLUMN IF EXISTS fibre_host_height
	`)
	return errors.Wrap(err, "drop validator fibre host columns")
}
