// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddHyperlanePostTypes, downAddHyperlanePostTypes)
}

func upAddHyperlanePostTypes(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgCreateIgp.String(), types.MsgUpdateMinfeeParams.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgSetIgpOwner.String(), types.MsgCreateIgp.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgSetDestinationGasConfig.String(), types.MsgSetIgpOwner.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgPayForGas.String(), types.MsgSetDestinationGasConfig.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgClaim.String(), types.MsgPayForGas.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgCreateMerkleTreeHook.String(), types.MsgClaim.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgCreateNoopHook.String(), types.MsgCreateMerkleTreeHook.String()); err != nil {
		return err
	}
	return nil
}
func downAddHyperlanePostTypes(ctx context.Context, db *bun.DB) error {
	return nil
}
