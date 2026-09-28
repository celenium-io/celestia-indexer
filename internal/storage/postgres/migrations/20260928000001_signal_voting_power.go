// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upSignalVotingPower, downSignalVotingPower)
}

// calculatedUpgrades lists every closed voting round: the version that won it,
// the height the tally was taken at and the height the round opened at.
const calculatedUpgrades = `
	WITH calc AS (
		SELECT u.version,
			COALESCE(NULLIF(u.end_height, 0), u.applied_at_level) AS close_height,
			COALESCE((
				SELECT MAX(p.applied_at_level) FROM upgrade p
				WHERE p.status = 'applied' AND p.version <> u.version
					AND p.applied_at_level <= COALESCE(NULLIF(u.end_height, 0), u.applied_at_level)
			), 0) AS open_height
		FROM upgrade u
		WHERE u.status = 'applied' OR u.end_height > 0
	)`

// signal_version.voting_power held the validator stake (utia), rewritten on every tally.
// Now it is null while the round is open, the counted power once it is closed, or 0 if not counted.
func upSignalVotingPower(ctx context.Context, db *bun.DB) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// the stake of counted signals was last rewritten close to the tally
		if _, err := tx.ExecContext(ctx, calculatedUpgrades+`
			UPDATE signal_version s SET voting_power = CASE
				WHEN NOT EXISTS (SELECT 1 FROM calc WHERE calc.close_height >= s.height) THEN NULL
				WHEN EXISTS (
					SELECT 1 FROM calc
					WHERE calc.version = s.version
						AND s.height BETWEEN calc.open_height AND calc.close_height
						AND NOT EXISTS (
							SELECT 1 FROM signal_version later
							WHERE later.validator_id = s.validator_id
								AND later.height > s.height
								AND later.height <= calc.close_height
						)
				) THEN trunc(s.voting_power / 1000000)
				ELSE 0
			END
		`); err != nil {
			return errors.Wrap(err, "convert signal voting power")
		}
		return nil
	})
}

// down is lossy: the stake of open and not counted signals is gone, so they get the current stake.
func downSignalVotingPower(ctx context.Context, db *bun.DB) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			UPDATE signal_version SET voting_power = CASE
				WHEN signal_version.voting_power > 0 THEN signal_version.voting_power * 1000000
				ELSE validator.stake
			END
			FROM validator
			WHERE validator.id = signal_version.validator_id
		`); err != nil {
			return errors.Wrap(err, "restore signal stake")
		}
		return nil
	})
}
