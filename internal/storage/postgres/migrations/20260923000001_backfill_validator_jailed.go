// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upBackfillValidatorJailed, downBackfillValidatorJailed)
}

// Rollback wrote a null jailed flag for every validator it reached through the staking
// logs, which is the whole active set on any block with rewards. Such a validator matches
// neither `jailed = true` nor `jailed = false`, and a null sorts first in ListByPower.
//
// The column stays nullable on purpose: SaveValidators sends a null in EXCLUDED to mean
// "keep the stored flag", and a not-null constraint rejects that row before the conflict
// is resolved, so the upsert would fail on every rewards event.
func upBackfillValidatorJailed(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `UPDATE validator SET jailed = false WHERE jailed IS NULL`); err != nil {
		return errors.Wrap(err, "backfill validator jailed")
	}

	_, err := db.ExecContext(ctx, `ALTER TABLE validator ALTER COLUMN jailed SET DEFAULT false`)
	return errors.Wrap(err, "set validator jailed default")
}

// Only the default is dropped. Which rows held a null is recorded nowhere, and false is the
// value they should have had, so the backfill is not undone.
func downBackfillValidatorJailed(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE validator ALTER COLUMN jailed DROP DEFAULT`)
	return errors.Wrap(err, "drop validator jailed default")
}
