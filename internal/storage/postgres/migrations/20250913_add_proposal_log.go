// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddProposalLog, downAddProposalLog)
}

func upAddProposalLog(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE public.proposal ADD error varchar NULL`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE public.proposal ALTER COLUMN error SET STORAGE EXTENDED`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `COMMENT ON COLUMN public.proposal.error IS 'Proposal error'`); err != nil {
		return err
	}
	return nil
}
func downAddProposalLog(ctx context.Context, db *bun.DB) error {
	return nil
}
