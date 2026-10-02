// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddHyperlaneIsmEvents, downAddHyperlaneIsmEvents)
}

var hyperlaneIsmEvents = []types.EventType{
	types.EventTypeHyperlanecoreismv1EventRemoveRoutingIsmDomain,
	types.EventTypeHyperlanecoreismv1EventAnnounceStorageLocation,
	types.EventTypeHyperlanecoreismv1EventCreateMessageIdMultisigIsm,
	types.EventTypeHyperlanecoreismv1EventCreateMerkleRootMultisigIsm,
}

func upAddHyperlaneIsmEvents(ctx context.Context, db *bun.DB) error {
	// ALTER TYPE ... ADD VALUE cannot run inside a transaction block
	for _, event := range hyperlaneIsmEvents {
		if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ?`, event.String()); err != nil {
			return errors.Wrapf(err, "add event_type %s", event)
		}
	}
	return nil
}

func downAddHyperlaneIsmEvents(ctx context.Context, db *bun.DB) error {
	for _, event := range hyperlaneIsmEvents {
		// before the migration these events were stored as unknown
		if _, err := db.ExecContext(ctx, `UPDATE event SET type = ? WHERE type = ?`,
			types.EventTypeUnknown.String(), event.String(),
		); err != nil {
			return errors.Wrapf(err, "restore unknown %s", event)
		}
		if _, err := db.ExecContext(ctx, `DELETE FROM pg_enum
			WHERE enumlabel = ? AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'event_type')`,
			event.String(),
		); err != nil {
			return errors.Wrapf(err, "drop event_type %s", event)
		}
	}
	return nil
}
