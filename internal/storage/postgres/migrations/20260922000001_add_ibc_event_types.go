// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddIbcEventTypes, downAddIbcEventTypes)
}

// ibc-go v8 events (core, transfer, ICA host) that were previously stored as 'unknown'
var ibcEventTypes = []types.EventType{
	types.EventTypeClientMisbehaviour,
	types.EventTypeUpgradeClient,
	types.EventTypeRecoverClient,
	types.EventTypeScheduleIbcSoftwareUpgrade,
	types.EventTypeUpgradeChain,
	types.EventTypeChannelCloseInit,
	types.EventTypeChannelClose,
	types.EventTypeChannelUpgradeInit,
	types.EventTypeChannelUpgradeTry,
	types.EventTypeChannelUpgradeAck,
	types.EventTypeChannelUpgradeConfirm,
	types.EventTypeChannelUpgradeOpen,
	types.EventTypeChannelUpgradeTimeout,
	types.EventTypeChannelUpgradeCancelled,
	types.EventTypeChannelUpgradeError,
	types.EventTypeChannelFlushComplete,
	types.EventTypeChannelClosed,
	types.EventTypeDenominationTrace,
	types.EventTypeIbccallbackerrorDenominationTrace,
}

func upAddIbcEventTypes(ctx context.Context, db *bun.DB) error {
	// ALTER TYPE ... ADD VALUE can't be batched, so one idempotent statement per value
	for _, t := range ibcEventTypes {
		if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ?`, t.String()); err != nil {
			return err
		}
	}
	return nil
}

func downAddIbcEventTypes(ctx context.Context, db *bun.DB) error {
	// restore pre-migration 'unknown' type so rows survive dropping the enum value
	for _, t := range ibcEventTypes {
		if _, err := db.ExecContext(ctx, `UPDATE event SET type = ? WHERE type = ?`, types.EventTypeUnknown.String(), t.String()); err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, `DELETE FROM pg_enum
			WHERE enumlabel = ? AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'event_type')`,
			t.String(),
		); err != nil {
			return err
		}
	}
	return nil
}
