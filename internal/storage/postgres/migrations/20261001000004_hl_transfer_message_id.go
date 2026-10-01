// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upHlTransferMessageId, downHlTransferMessageId)
}

// Existing rows keep NULL: message ids are filled by reindexing.
func upHlTransferMessageId(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE hl_transfer ADD COLUMN IF NOT EXISTS message_id bytea`)
	return errors.Wrap(err, "add hl_transfer.message_id")
}

func downHlTransferMessageId(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE hl_transfer DROP COLUMN IF EXISTS message_id`)
	return errors.Wrap(err, "drop hl_transfer.message_id")
}
