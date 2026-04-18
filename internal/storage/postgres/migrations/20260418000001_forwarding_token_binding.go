// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upForwardingTokenBinding, downForwardingTokenBinding)
}

func upForwardingTokenBinding(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		ALTER TABLE forwarding
			DROP COLUMN IF EXISTS success_count,
			DROP COLUMN IF EXISTS failed_count,
			DROP COLUMN IF EXISTS transfers,
			ADD COLUMN IF NOT EXISTS token_id   bigint  NOT NULL DEFAULT 0,
			ADD COLUMN IF NOT EXISTS message_id text    NOT NULL DEFAULT '',
			ADD COLUMN IF NOT EXISTS amount     numeric NOT NULL DEFAULT 0,
			ADD COLUMN IF NOT EXISTS denom      text    NOT NULL DEFAULT ''
	`)
	return err
}

func downForwardingTokenBinding(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		ALTER TABLE forwarding
			DROP COLUMN IF EXISTS token_id,
			DROP COLUMN IF EXISTS message_id,
			DROP COLUMN IF EXISTS amount,
			DROP COLUMN IF EXISTS denom,
			ADD COLUMN IF NOT EXISTS success_count bigint,
			ADD COLUMN IF NOT EXISTS failed_count  bigint,
			ADD COLUMN IF NOT EXISTS transfers      jsonb
	`)
	return err
}
