// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upRecountProposalVotes, downRecountProposalVotes)
}

// Counters were incremented per vote row, so weighted votes were counted once per option.
// votes_count is the number of voters; option counters are voters who chose the option.
const recountProposalVotesQuery = `
	UPDATE proposal SET
		votes_count        = c.votes_count,
		yes                = c.yes,
		no                 = c.no,
		no_with_veto       = c.no_with_veto,
		abstain            = c.abstain,
		yes_vals           = c.yes_vals,
		no_vals            = c.no_vals,
		no_with_veto_vals  = c.no_with_veto_vals,
		abstain_vals       = c.abstain_vals,
		yes_addrs          = c.yes_addrs,
		no_addrs           = c.no_addrs,
		no_with_veto_addrs = c.no_with_veto_addrs,
		abstain_addrs      = c.abstain_addrs
	FROM (
		SELECT
			p.id,
			count(DISTINCT v.voter_id) AS votes_count,
			count(v.id) FILTER (WHERE v.option = 'yes')          AS yes,
			count(v.id) FILTER (WHERE v.option = 'no')           AS no,
			count(v.id) FILTER (WHERE v.option = 'no_with_veto') AS no_with_veto,
			count(v.id) FILTER (WHERE v.option = 'abstain')      AS abstain,
			count(v.id) FILTER (WHERE v.option = 'yes'          AND v.validator_id IS NOT NULL) AS yes_vals,
			count(v.id) FILTER (WHERE v.option = 'no'           AND v.validator_id IS NOT NULL) AS no_vals,
			count(v.id) FILTER (WHERE v.option = 'no_with_veto' AND v.validator_id IS NOT NULL) AS no_with_veto_vals,
			count(v.id) FILTER (WHERE v.option = 'abstain'      AND v.validator_id IS NOT NULL) AS abstain_vals,
			count(v.id) FILTER (WHERE v.option = 'yes'          AND v.validator_id IS NULL) AS yes_addrs,
			count(v.id) FILTER (WHERE v.option = 'no'           AND v.validator_id IS NULL) AS no_addrs,
			count(v.id) FILTER (WHERE v.option = 'no_with_veto' AND v.validator_id IS NULL) AS no_with_veto_addrs,
			count(v.id) FILTER (WHERE v.option = 'abstain'      AND v.validator_id IS NULL) AS abstain_addrs
		FROM proposal p
		LEFT JOIN vote v ON v.proposal_id = p.id
		GROUP BY p.id
	) c
	WHERE proposal.id = c.id
`

func upRecountProposalVotes(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, recountProposalVotesQuery)
	return errors.Wrap(err, "recount proposal votes")
}

// The old counters were wrong for weighted votes, so there is nothing to restore.
func downRecountProposalVotes(_ context.Context, _ *bun.DB) error {
	return nil
}
