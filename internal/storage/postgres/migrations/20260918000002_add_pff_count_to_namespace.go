// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddPffCountToNamespace, downAddPffCountToNamespace)
}

func upAddPffCountToNamespace(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `
		ALTER TABLE namespace
			ADD COLUMN IF NOT EXISTS pff_count  bigint,
			ADD COLUMN IF NOT EXISTS fibre_size bigint
	`); err != nil {
		return errors.Wrap(err, "add namespace fibre columns")
	}

	// The namespace upsert sums EXCLUDED values with the stored ones, so existing
	// rows must hold 0 rather than NULL.
	if _, err := db.ExecContext(ctx, `
		UPDATE namespace
		SET pff_count  = coalesce(pff_count, 0),
			fibre_size = coalesce(fibre_size, 0)
		WHERE pff_count IS NULL OR fibre_size IS NULL
	`); err != nil {
		return errors.Wrap(err, "backfill namespace fibre columns")
	}

	return nil
}

func downAddPffCountToNamespace(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		ALTER TABLE namespace
			DROP COLUMN IF EXISTS pff_count,
			DROP COLUMN IF EXISTS fibre_size
	`)
	return errors.Wrap(err, "drop namespace fibre columns")
}
