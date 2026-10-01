// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upTotalVotingPowerUtia, downTotalVotingPowerUtia)
}

// total_voting_power was consensus power (utia / 10^6), while the voting power is in utia.
// Migrations may replay on a database indexed by the new code, so only rows still in consensus
// power are converted: they are ~10^6 times below the current bonded stake, not just ~1000.
const totalVotingPowerUnitsCondition = `
	total_voting_power > 0 AND total_voting_power * 1000 < (
		SELECT coalesce(sum(stake), 0) FROM validator WHERE coalesce(power, 0) > 0
	)`

func upTotalVotingPowerUtia(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `UPDATE proposal
		SET total_voting_power = total_voting_power * 1000000
		WHERE `+totalVotingPowerUnitsCondition)
	return errors.Wrap(err, "convert total_voting_power to utia")
}

func downTotalVotingPowerUtia(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `UPDATE proposal
		SET total_voting_power = floor(total_voting_power / 1000000)
		WHERE total_voting_power > 0 AND NOT (`+totalVotingPowerUnitsCondition+`)`)
	return errors.Wrap(err, "convert total_voting_power to consensus power")
}
