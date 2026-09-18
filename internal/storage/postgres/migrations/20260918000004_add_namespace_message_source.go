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
	Migrations.MustRegister(upAddNamespaceMessageSource, downAddNamespaceMessageSource)
}

func upAddNamespaceMessageSource(ctx context.Context, db *bun.DB) error {
	// blob_source is created by the blob_log migration, which runs first.
	// Existing rows predate app v10, so reading them as 'pfb' is correct, and
	// the default keeps the hypertable from being rewritten.
	if _, err := db.ExecContext(ctx,
		`ALTER TABLE namespace_message ADD COLUMN IF NOT EXISTS source blob_source NOT NULL DEFAULT ?`,
		types.BlobSourcePfb.String(),
	); err != nil {
		return errors.Wrap(err, "add namespace_message.source column")
	}

	return nil
}

func downAddNamespaceMessageSource(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE namespace_message DROP COLUMN IF EXISTS source`)
	return errors.Wrap(err, "drop namespace_message.source column")
}
