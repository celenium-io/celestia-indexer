// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddEvent, downAddEvent)
}

func upAddEvent(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeIbccallbackerrorIcs27Packet.String(), types.EventTypeSignalVersion.String()); err != nil {
		return err
	}
	return nil
}
func downAddEvent(ctx context.Context, db *bun.DB) error {
	return nil
}
