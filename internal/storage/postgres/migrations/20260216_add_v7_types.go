// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddV7Types, downAddV7Types)
}

func upAddV7Types(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgForward.String(), types.MsgModuleQuerySafe.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgCreateInterchainSecurityModule.String(), types.MsgForward.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgUpdateInterchainSecurityModule.String(), types.MsgCreateInterchainSecurityModule.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.MsgSubmitMessages.String(), types.MsgUpdateInterchainSecurityModule.String()); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `ALTER TABLE block
		ALTER COLUMN message_types
		TYPE bit varying(115)
		USING '0000'::bit varying(4) || message_types;`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE tx
		ALTER COLUMN message_types
		TYPE bit varying(115)
		USING '0000'::bit varying(4) || message_types;`); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeCelestiaforwardingv1EventTokenForwarded.String(), types.EventTypeIbccallbackerrorIcs27Packet.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeCelestiaforwardingv1EventForwardingComplete.String(), types.EventTypeCelestiaforwardingv1EventTokenForwarded.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeCelestiazkismv1EventCreateInterchainSecurityModule.String(), types.EventTypeCelestiaforwardingv1EventForwardingComplete.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeCelestiazkismv1EventUpdateInterchainSecurityModule.String(), types.EventTypeCelestiazkismv1EventCreateInterchainSecurityModule.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeCelestiazkismv1EventSubmitMessages.String(), types.EventTypeCelestiazkismv1EventUpdateInterchainSecurityModule.String()); err != nil {
		return err
	}
	return nil
}
func downAddV7Types(ctx context.Context, db *bun.DB) error {
	return nil
}
