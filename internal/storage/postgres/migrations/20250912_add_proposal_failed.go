// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddProposalFailed, downAddProposalFailed)
}

func upAddProposalFailed(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TYPE proposal_status ADD VALUE IF NOT EXISTS ? AFTER ?`, types.ProposalStatusFailed.String(), types.ProposalStatusRejected.String()); err != nil {
		return err
	}
	return nil
}
func downAddProposalFailed(ctx context.Context, db *bun.DB) error {
	return nil
}
