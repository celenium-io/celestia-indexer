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
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS 'MsgSignalVersion' AFTER 'MsgAcknowledgement'`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS 'MsgTryUpgrade' AFTER 'MsgSignalVersion'`); err != nil {
		return err
	}
	return nil
}

func downSignalMsgTypes(_ context.Context, _ *bun.DB) error {
	return nil
}
