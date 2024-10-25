// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upSignalMsgTypes, downSignalMsgTypes)
}

func upSignalMsgTypes(ctx context.Context, db *bun.DB) error {
	var maxLen int
	if err := db.NewSelect().
		TableExpr("information_schema.columns").
		Column("character_maximum_length").
		Where("table_name = 'tx'").
		Where("column_name = 'message_types'").
		Scan(ctx, &maxLen); err != nil {
		return err
	}

	if maxLen == 76 {
		return nil
	}

	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS 'MsgSignalVersion' AFTER 'MsgAcknowledgement'`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS 'MsgTryUpgrade' AFTER 'MsgSignalVersion'`); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `ALTER TABLE block ADD COLUMN IF NOT EXISTS message_types_1 bit(76)`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE tx ADD COLUMN IF NOT EXISTS message_types_1 bit(76)`); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `UPDATE block SET message_types_1 = ('00' || message_types::text)::bit(76)`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `UPDATE tx SET message_types_1 = ('00' || message_types::text)::bit(76)`); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `ALTER TABLE block DROP COLUMN IF EXISTS message_types CASCADE`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE tx DROP COLUMN IF EXISTS message_types CASCADE`); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `ALTER TABLE block RENAME COLUMN message_types_1 TO message_types`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE tx RENAME COLUMN message_types_1 TO message_types`); err != nil {
		return err
	}
	return nil
}

func downSignalMsgTypes(_ context.Context, _ *bun.DB) error {
	return nil
}
