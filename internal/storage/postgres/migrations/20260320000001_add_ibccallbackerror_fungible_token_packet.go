// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddEventType, downAddEventType)
}

func upAddEventType(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx,
		`ALTER TYPE event_type ADD VALUE IF NOT EXISTS 'ibccallbackerror-fungible_token_packet'`,
	)
	return err
}

func downAddEventType(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx,
		`DELETE FROM pg_enum
		WHERE enumlabel = 'ibccallbackerror-fungible_token_packet'
		AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'event_type')`,
	)
	return err
}
