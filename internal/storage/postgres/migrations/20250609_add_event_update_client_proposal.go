// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddEventUpdateClientProposal, downAddEventUpdateClientProposal)
}

func upAddEventUpdateClientProposal(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeUpdateClientProposal.String(), types.EventTypeChannelCloseConfirm.String())
	return err
}
func downAddEventUpdateClientProposal(ctx context.Context, db *bun.DB) error {
	return nil
}
