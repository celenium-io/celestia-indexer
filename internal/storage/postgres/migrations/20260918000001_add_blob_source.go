// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddBlobSource, downAddBlobSource)
}

func upAddBlobSource(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx,
		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'blob_source') THEN
				CREATE TYPE blob_source AS ENUM ?;
			END IF;
		END$$;`,
		bun.Tuple(types.BlobSourceValues()),
	); err != nil {
		return errors.Wrap(err, "create blob_source type")
	}

	// Default makes existing rows read as 'pfb' without rewriting the hypertable.
	if _, err := db.ExecContext(ctx,
		`ALTER TABLE blob_log ADD COLUMN IF NOT EXISTS source blob_source NOT NULL DEFAULT ?`,
		types.BlobSourcePfb.String(),
	); err != nil {
		return errors.Wrap(err, "add blob_log.source column")
	}

	return nil
}

func downAddBlobSource(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE blob_log DROP COLUMN IF EXISTS source`); err != nil {
		return errors.Wrap(err, "drop blob_log.source column")
	}

	if _, err := db.ExecContext(ctx, `DROP TYPE IF EXISTS blob_source`); err != nil {
		return errors.Wrap(err, "drop blob_source type")
	}

	return nil
}
