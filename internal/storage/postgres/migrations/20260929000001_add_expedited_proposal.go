// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddExpeditedProposal, downAddExpeditedProposal)
}

func upAddExpeditedProposal(ctx context.Context, db *bun.DB) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(ctx, `ALTER TABLE proposal ADD COLUMN expedited boolean NOT NULL DEFAULT false`); err != nil {
			return errors.Wrap(err, "add expedited")
		}
		return nil
	})
}

func downAddExpeditedProposal(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE proposal DROP COLUMN IF EXISTS expedited`)
	return errors.Wrap(err, "drop expedited")
}
