// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddFibreTypes, downAddFibreTypes)
}

// msg_type values added by app v10, in bit-mask order: each one must be appended
// after its predecessor so the enum order keeps matching MsgTypeValues().
var fibreMsgTypes = []types.MsgType{
	types.MsgDepositToEscrow,
	types.MsgRequestWithdrawal,
	types.MsgPayForFibre,
	types.MsgPaymentPromiseTimeout,
	types.MsgUpdateFibreParams,
	types.MsgSetFibreProviderInfo,
}

var fibreEventTypes = []types.EventType{
	types.EventTypeCelestiafibrev1EventDepositToEscrow,
	types.EventTypeCelestiafibrev1EventWithdrawFromEscrowRequest,
	types.EventTypeCelestiafibrev1EventWithdrawFromEscrowExecuted,
	types.EventTypeCelestiafibrev1EventPayForFibre,
	types.EventTypeCelestiafibrev1EventPaymentPromiseTimeout,
	types.EventTypeCelestiafibrev1EventUpdateFibreParams,
	types.EventTypeCelestiafibrev1EventProcessedPaymentPruned,
	types.EventTypeSetFibreProviderInfo,
}

func upAddFibreTypes(ctx context.Context, db *bun.DB) error {
	// ALTER TYPE ... ADD VALUE cannot run inside a transaction block, so every
	// statement goes on its own and is idempotent.
	prev := types.MsgSubmitMessages.String()
	for _, t := range fibreMsgTypes {
		if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`, t.String(), prev); err != nil {
			return err
		}
		prev = t.String()
	}

	for _, t := range fibreEventTypes {
		if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ?`, t.String()); err != nil {
			return err
		}
	}

	if _, err := db.ExecContext(ctx, `ALTER TYPE module_name ADD VALUE IF NOT EXISTS 'fibre'`); err != nil {
		return err
	}

	// New types take the high bits, so old masks are left-padded to keep every
	// existing bit in place. Padding must be a literal: there is no text -> bit cast.
	//
	// This step is not naturally idempotent like the ones above: on a fresh
	// deployment the table is created with message_types already at
	// MsgTypeBitsCount (the table doesn't exist yet, so initDatabaseWithMigrations
	// skips migrateDatabase entirely and no migration gets recorded as applied).
	// A later restart then replays every migration against that already-current
	// schema, and re-padding an already MsgTypeBitsCount-wide column would push
	// it past MsgTypeBitsCount and silently truncate the low bits on cast back
	// down, corrupting every existing mask. Guard on the column's current width.
	padding := strings.Repeat("0", len(fibreMsgTypes))
	for _, table := range []string{"block", "tx"} {
		var width int
		if err := db.QueryRowContext(ctx, `SELECT COALESCE(character_maximum_length, 0)
			FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = ? AND column_name = 'message_types'`,
			table,
		).Scan(&width); err != nil {
			return err
		}
		if width == types.MsgTypeBitsCount {
			continue
		}

		if _, err := db.ExecContext(ctx, fmt.Sprintf(`ALTER TABLE %s
			ALTER COLUMN message_types
			TYPE bit(%d)
			USING (B'%s' || message_types)::bit(%d)`,
			table, types.MsgTypeBitsCount, padding, types.MsgTypeBitsCount,
		)); err != nil {
			return err
		}
	}

	return nil
}

func downAddFibreTypes(ctx context.Context, db *bun.DB) error {
	// Drop the padding first, then the rows using the new types: an enum value
	// can only be deleted from pg_enum while nothing references it.
	padding := len(fibreMsgTypes)
	width := types.MsgTypeBitsCount - padding
	for _, table := range []string{"block", "tx"} {
		if _, err := db.ExecContext(ctx, fmt.Sprintf(`ALTER TABLE %s
			ALTER COLUMN message_types
			TYPE bit(%d)
			USING substring(message_types from %d for %d)::bit(%d)`,
			table, width, padding+1, width, width,
		)); err != nil {
			return err
		}
	}

	for _, t := range fibreMsgTypes {
		if _, err := db.ExecContext(ctx, `DELETE FROM message WHERE type = ?`, t.String()); err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, `DELETE FROM pg_enum
			WHERE enumlabel = ? AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'msg_type')`,
			t.String(),
		); err != nil {
			return err
		}
	}

	for _, t := range fibreEventTypes {
		if _, err := db.ExecContext(ctx, `DELETE FROM event WHERE type = ?`, t.String()); err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, `DELETE FROM pg_enum
			WHERE enumlabel = ? AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'event_type')`,
			t.String(),
		); err != nil {
			return err
		}
	}

	if _, err := db.ExecContext(ctx, `DELETE FROM constant WHERE module = 'fibre'`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `DELETE FROM pg_enum
		WHERE enumlabel = 'fibre' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'module_name')`,
	); err != nil {
		return err
	}

	return nil
}
