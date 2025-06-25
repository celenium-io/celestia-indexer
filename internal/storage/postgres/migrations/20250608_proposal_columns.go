// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upCreateProposalColumns, downProposalColumns)
	Migrations.MustRegister(upInitProposalColumns, downProposalColumns)
}

func upCreateProposalColumns(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		ALTER TABLE public.proposal ADD end_time timestamptz NULL;
		ALTER TABLE public.proposal ALTER COLUMN end_time SET STORAGE PLAIN;
		COMMENT ON COLUMN public.proposal.end_time IS 'Voting end time';
	`)
	if err != nil {
		return errors.Wrap(err, "create end_time")
	}
	_, err = db.ExecContext(ctx, `
		ALTER TABLE public.proposal ADD total_voting_power numeric DEFAULT 0 NULL;
		ALTER TABLE public.proposal ALTER COLUMN total_voting_power SET STORAGE MAIN;
		COMMENT ON COLUMN public.proposal.total_voting_power IS 'Total voting power in the network';

	`)
	if err != nil {
		return errors.Wrap(err, "create total_voting_power")
	}
	_, err = db.ExecContext(ctx, `
		ALTER TABLE public.proposal ADD quorum varchar NULL;
		ALTER TABLE public.proposal ALTER COLUMN quorum SET STORAGE EXTENDED;
		COMMENT ON COLUMN public.proposal.quorum IS 'The minimum percentage of voting power that needs to be cast on a proposal for the result to be valid';

	`)
	if err != nil {
		return errors.Wrap(err, "create quorum")
	}
	_, err = db.ExecContext(ctx, `
		ALTER TABLE public.proposal ADD veto_quorum varchar NULL;
		ALTER TABLE public.proposal ALTER COLUMN veto_quorum SET STORAGE EXTENDED;
		COMMENT ON COLUMN public.proposal.veto_quorum IS 'Minimum value of Veto votes to Total votes ratio for proposal to be vetoed';

	`)
	if err != nil {
		return errors.Wrap(err, "create veto_quorum")
	}
	_, err = db.ExecContext(ctx, `
		ALTER TABLE public.proposal ADD threshold varchar NULL;
		ALTER TABLE public.proposal ALTER COLUMN threshold SET STORAGE EXTENDED;
		COMMENT ON COLUMN public.proposal.threshold IS 'Minimum proportion of Yes votes for proposal to pass';

	`)
	if err != nil {
		return errors.Wrap(err, "create threshold")
	}
	_, err = db.ExecContext(ctx, `
		ALTER TABLE public.proposal ADD min_deposit varchar NULL;
		ALTER TABLE public.proposal ALTER COLUMN min_deposit SET STORAGE EXTENDED;
		COMMENT ON COLUMN public.proposal.min_deposit IS 'Minimum deposit for a proposal to enter voting period';

	`)

	return err
}

func upInitProposalColumns(ctx context.Context, db *bun.DB) error {
	var quorum storage.Constant
	err := db.NewSelect().
		Model(&quorum).
		Where("name = 'quorum'").
		Where("module = 'gov'").
		Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "receiving quorum")
	}

	_, err = db.NewUpdate().
		Model((*storage.Proposal)(nil)).
		Set("quorum = ?", quorum.Value).
		Where("id > 0").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "update quorum")
	}

	var threshold storage.Constant
	err = db.NewSelect().
		Model(&threshold).
		Where("name = 'threshold'").
		Where("module = 'gov'").
		Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "receiving threshold")
	}

	_, err = db.NewUpdate().
		Model((*storage.Proposal)(nil)).
		Set("threshold = ?", threshold.Value).
		Where("id > 0").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "update threshold")
	}

	var vetoQuorum storage.Constant
	err = db.NewSelect().
		Model(&vetoQuorum).
		Where("name = 'veto_threshold'").
		Where("module = 'gov'").
		Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "receiving veto_threshold")
	}

	_, err = db.NewUpdate().
		Model((*storage.Proposal)(nil)).
		Set("veto_quorum = ?", vetoQuorum.Value).
		Where("id > 0").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "update veto_quorum")
	}

	var minDeposit storage.Constant
	err = db.NewSelect().
		Model(&minDeposit).
		Where("name = 'min_deposit'").
		Where("module = 'gov'").
		Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "receiving min_deposit")
	}

	_, err = db.NewUpdate().
		Model((*storage.Proposal)(nil)).
		Set("min_deposit = ?", minDeposit.Value).
		Where("id > 0").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "update min_deposit")
	}

	var votingPeriod storage.Constant
	err = db.NewSelect().
		Model(&votingPeriod).
		Where("name = 'voting_period'").
		Where("module = 'gov'").
		Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "receiving voting_period")
	}

	period, err := strconv.ParseInt(votingPeriod.Value, 10, 64)
	if err != nil {
		return errors.Wrap(err, "parsing voting period")
	}
	durationVotingPeriod := time.Duration(period)

	var proposals []storage.Proposal
	err = db.NewSelect().
		Model(&proposals).
		Limit(1000).
		Where("status = 'applied'").
		WhereOr("status = 'rejected'").
		WhereOr("status = 'active'").
		Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "receiving proposals")
	}

	for i := range proposals {
		if proposals[i].ActivationTime != nil {
			_, err = db.NewUpdate().
				Model(&proposals[i]).
				Set("end_time = ?", proposals[i].ActivationTime.Add(durationVotingPeriod)).
				WherePK().
				Exec(ctx)
			if err != nil {
				return errors.Wrap(err, "update end_time")
			}
		}
	}

	return nil
}

func downProposalColumns(ctx context.Context, db *bun.DB) error {
	return nil
}
