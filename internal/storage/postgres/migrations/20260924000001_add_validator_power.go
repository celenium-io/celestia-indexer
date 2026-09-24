// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddValidatorPower, downAddValidatorPower)
}

// The validator_bond_update table is new, so initDatabase creates it after migrations.
func upAddValidatorPower(ctx context.Context, db *bun.DB) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// power stays nullable: SaveValidators sends null for validators absent from validator_updates
		if _, err := tx.ExecContext(ctx, `
			ALTER TABLE validator
				ADD COLUMN IF NOT EXISTS power              numeric,
				ADD COLUMN IF NOT EXISTS bond_updates_count bigint NOT NULL DEFAULT 0
		`); err != nil {
			return errors.Wrap(err, "add validator power columns")
		}

		if _, err := tx.ExecContext(ctx, `UPDATE validator SET messages_count = 0 WHERE messages_count IS NULL`); err != nil {
			return errors.Wrap(err, "fill null messages_count")
		}
		if _, err := tx.ExecContext(ctx, `
			ALTER TABLE validator
				ALTER COLUMN messages_count SET DEFAULT 0,
				ALTER COLUMN messages_count SET NOT NULL
		`); err != nil {
			return errors.Wrap(err, "set messages_count default and not null")
		}

		// validator_updates of past blocks were not indexed, so the bonded set is approximated
		// as the top max_validators unjailed by stake, power = tokens / 10^6 truncated.
		if _, err := tx.ExecContext(ctx, `
			WITH bonded AS (
				SELECT id FROM validator
				WHERE jailed = false
				ORDER BY stake DESC
				LIMIT COALESCE((SELECT value::int FROM constant WHERE module = 'staking' AND name = 'max_validators'), 100)
			)
			UPDATE validator SET power = CASE
				WHEN validator.id IN (SELECT id FROM bonded) THEN trunc(validator.stake / 1000000)
				ELSE 0
			END
		`); err != nil {
			return errors.Wrap(err, "backfill validator power")
		}
		return nil
	})
}

func downAddValidatorPower(ctx context.Context, db *bun.DB) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS validator_bond_update`); err != nil {
			return errors.Wrap(err, "drop validator_bond_update")
		}
		if _, err := tx.ExecContext(ctx, `
			ALTER TABLE validator
				ALTER COLUMN messages_count DROP NOT NULL,
				ALTER COLUMN messages_count DROP DEFAULT,
				DROP COLUMN IF EXISTS power,
				DROP COLUMN IF EXISTS bond_updates_count
		`); err != nil {
			return errors.Wrap(err, "drop validator power columns")
		}
		return nil
	})
}
