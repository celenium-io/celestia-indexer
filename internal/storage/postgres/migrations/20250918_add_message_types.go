// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upChangeMsgTypes, downChangeMsgTypes)
}

func upChangeMsgTypes(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgUpdateBlobParams.String(), types.MsgUpdateRoutingIsmOwner.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgPruneExpiredGrants.String(), types.MsgUpdateBlobParams.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgSetSendEnabled.String(), types.MsgPruneExpiredGrants.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgAuthorizeCircuitBreaker.String(), types.MsgSetSendEnabled.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgResetCircuitBreaker.String(), types.MsgAuthorizeCircuitBreaker.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgTripCircuitBreaker.String(), types.MsgResetCircuitBreaker.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgModuleQuerySafe.String(), types.MsgTripCircuitBreaker.String()); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `ALTER TABLE block
		ALTER COLUMN message_types
		TYPE bit varying(111)
		USING '000000000000'::bit varying(7) || message_types;`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE tx
		ALTER COLUMN message_types
		TYPE bit varying(111)
		USING '000000000000'::bit varying(7) || message_types;`); err != nil {
		return err
	}

	return nil
}
func downChangeMsgTypes(ctx context.Context, db *bun.DB) error {
	return nil
}
