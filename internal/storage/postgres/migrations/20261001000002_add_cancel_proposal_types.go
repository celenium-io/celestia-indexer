// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"
	"fmt"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddCancelProposalTypes, downAddCancelProposalTypes)
}

var messageTypesTables = []string{"block", "tx"}

// Mask width after this migration. Hardcoded: MsgTypeBitsCount grows with later migrations.
const cancelProposalWidth = 122

func upAddCancelProposalTypes(ctx context.Context, db *bun.DB) error {
	// ALTER TYPE ... ADD VALUE cannot run inside a transaction block, so every
	// statement goes on its own and is idempotent. msg_type keeps the bit-mask order.
	if _, err := db.ExecContext(ctx, `ALTER TYPE msg_type ADD VALUE IF NOT EXISTS ? AFTER ?`,
		types.MsgCancelProposal.String(), types.MsgSetFibreProviderInfo.String(),
	); err != nil {
		return errors.Wrap(err, "add msg_type")
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ?`,
		types.EventTypeCancelProposal.String(),
	); err != nil {
		return errors.Wrap(err, "add event_type")
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE proposal_status ADD VALUE IF NOT EXISTS ?`,
		types.ProposalStatusCancelled.String(),
	); err != nil {
		return errors.Wrap(err, "add proposal_status")
	}

	// The new type takes the high bit, so old masks are left-padded with one zero.
	// A fresh database is created at the current width, so skip already widened columns.
	for _, table := range messageTypesTables {
		width, err := messageTypesWidth(ctx, db, table)
		if err != nil {
			return err
		}
		if width >= cancelProposalWidth {
			continue
		}
		if _, err := db.ExecContext(ctx, fmt.Sprintf(`ALTER TABLE %s
			ALTER COLUMN message_types
			TYPE bit(%d)
			USING (B'0' || message_types)::bit(%d)`,
			table, cancelProposalWidth, cancelProposalWidth,
		)); err != nil {
			return errors.Wrapf(err, "widen %s.message_types", table)
		}
	}
	return nil
}

func downAddCancelProposalTypes(ctx context.Context, db *bun.DB) error {
	width := cancelProposalWidth - 1
	for _, table := range messageTypesTables {
		current, err := messageTypesWidth(ctx, db, table)
		if err != nil {
			return err
		}
		if current == width {
			continue
		}
		if current != cancelProposalWidth {
			return errors.Errorf("%s.message_types is bit(%d), expected bit(%d)", table, current, cancelProposalWidth)
		}
		if _, err := db.ExecContext(ctx, fmt.Sprintf(`ALTER TABLE %s
			ALTER COLUMN message_types
			TYPE bit(%d)
			USING substring(message_types from 2 for %d)::bit(%d)`,
			table, width, width, width,
		)); err != nil {
			return errors.Wrapf(err, "narrow %s.message_types", table)
		}
	}

	// rows keep the values they had before the migration, so the enum values can be dropped
	if _, err := db.ExecContext(ctx, `UPDATE message SET type = ? WHERE type = ?`,
		types.MsgUnknown.String(), types.MsgCancelProposal.String(),
	); err != nil {
		return errors.Wrap(err, "restore unknown messages")
	}
	if _, err := db.ExecContext(ctx, `UPDATE event SET type = ? WHERE type = ?`,
		types.EventTypeUnknown.String(), types.EventTypeCancelProposal.String(),
	); err != nil {
		return errors.Wrap(err, "restore unknown events")
	}
	// before the migration a canceled proposal stayed in its last status
	if _, err := db.ExecContext(ctx, `UPDATE proposal
		SET status = CASE WHEN activation_time IS NULL THEN ?::proposal_status ELSE ?::proposal_status END
		WHERE status = ?`,
		types.ProposalStatusInactive.String(), types.ProposalStatusActive.String(), types.ProposalStatusCancelled.String(),
	); err != nil {
		return errors.Wrap(err, "restore proposal status")
	}

	for typ, value := range map[string]string{
		"msg_type":        types.MsgCancelProposal.String(),
		"event_type":      types.EventTypeCancelProposal.String(),
		"proposal_status": types.ProposalStatusCancelled.String(),
	} {
		if _, err := db.ExecContext(ctx, `DELETE FROM pg_enum
			WHERE enumlabel = ? AND enumtypid = (SELECT oid FROM pg_type WHERE typname = ?)`,
			value, typ,
		); err != nil {
			return errors.Wrapf(err, "drop %s value", typ)
		}
	}
	return nil
}

func messageTypesWidth(ctx context.Context, db *bun.DB, table string) (int, error) {
	var width int
	err := db.QueryRowContext(ctx, `SELECT COALESCE(character_maximum_length, 0)
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = ? AND column_name = 'message_types'`,
		table,
	).Scan(&width)
	return width, errors.Wrapf(err, "%s.message_types width", table)
}
