// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
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

		var state storage.State
		if err := tx.NewSelect().
			Model(&state).
			Limit(1).
			Scan(ctx); err != nil {
			return err
		}
		if state.Version < 4 {
			return nil
		}

		constants := []storage.Constant{
			{
				Module: types.ModuleNameGov,
				Name:   "expedited_voting_period",
				Value:  "86400000000000",
			}, {
				Module: types.ModuleNameGov,
				Name:   "expedited_threshold",
				Value:  "0.667000000000000000",
			}, {
				Module: types.ModuleNameGov,
				Name:   "expedited_min_deposit",
				Value:  "50000000000utia",
			},
		}

		_, err := tx.NewInsert().
			Model(&constants).
			On("CONFLICT DO NOTHING").
			Exec(ctx)
		return errors.Wrap(err, "insert constants")
	})
}

func downAddExpeditedProposal(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE proposal DROP COLUMN IF EXISTS expedited`); err != nil {
		return errors.Wrap(err, "drop expedited")
	}
	_, err := db.NewDelete().
		Model((*storage.Constant)(nil)).
		Where("module = ?", types.ModuleNameGov.String()).
		Where("name IN (?)", bun.List([]string{"expedited_voting_period", "expedited_threshold", "expedited_min_deposit"})).
		Exec(ctx)
	return errors.Wrap(err, "delete constants")
}
