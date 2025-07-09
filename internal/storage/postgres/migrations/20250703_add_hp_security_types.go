// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddHyperlaneSecurityTypes, downAddHyperlaneScurityTypes)
}

func upAddHyperlaneSecurityTypes(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgCreateMessageIdMultisigIsm.String(), types.MsgCreateNoopHook.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgCreateMerkleRootMultisigIsm.String(), types.MsgCreateMessageIdMultisigIsm.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgCreateNoopIsm.String(), types.MsgCreateMerkleRootMultisigIsm.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgAnnounceValidator.String(), types.MsgCreateNoopIsm.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgCreateRoutingIsm.String(), types.MsgAnnounceValidator.String()); err != nil {
		return err
	}
	return nil
}
func downAddHyperlaneScurityTypes(ctx context.Context, db *bun.DB) error {
	return nil
}
