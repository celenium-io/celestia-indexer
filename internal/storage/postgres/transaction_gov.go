// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func (tx Transaction) SaveUpgrades(ctx context.Context, upgrades ...*models.Upgrade) error {
	if len(upgrades) == 0 {
		return nil
	}

	for i := range upgrades {
		if upgrades[i].Status == "" {
			upgrades[i].Status = storageTypes.UpgradeStatusProcessing
		}

		query := tx.Tx().NewInsert().Model(upgrades[i]).
			Column("version", "height", "time", "end_height", "end_time", "applied_at_level", "applied_at", "signer_id", "msg_id", "tx_id", "voting_power", "voted_power", "signals_count", "status").
			On("CONFLICT (version) DO UPDATE")

		if upgrades[i].EndHeight > 0 {
			query = query.Set("end_height = EXCLUDED.end_height")
		}
		if !upgrades[i].EndTime.IsZero() {
			query = query.Set("end_time = EXCLUDED.end_time")
		}
		if upgrades[i].AppliedAtLevel > 0 {
			query = query.Set("applied_at_level = EXCLUDED.applied_at_level")
		}
		if !upgrades[i].AppliedAt.IsZero() {
			query = query.Set("applied_at = EXCLUDED.applied_at")
		}
		if upgrades[i].SignerId > 0 {
			query = query.Set("signer_id = EXCLUDED.signer_id")
		}
		if upgrades[i].MsgId > 0 {
			query = query.Set("msg_id = EXCLUDED.msg_id")
		}
		if upgrades[i].TxId > 0 {
			query = query.Set("tx_id = EXCLUDED.tx_id")
		}
		if !upgrades[i].VotingPower.IsZero() {
			query = query.Set("voting_power = EXCLUDED.voting_power")
		}
		if !upgrades[i].VotedPower.IsZero() {
			query = query.Set("voted_power = EXCLUDED.voted_power")
		}
		if upgrades[i].SignalsCount > 0 {
			query = query.Set("signals_count = EXCLUDED.signals_count + upgrade.signals_count")
		}
		if upgrades[i].Status != storageTypes.UpgradeStatusProcessing {
			query = query.Set("status = EXCLUDED.status")
		}

		if _, err := query.Exec(ctx); err != nil {
			return errors.Wrapf(err, "save upgrade %d", upgrades[i].Version)
		}
	}

	return nil
}

type addedProposal struct {
	bun.BaseModel `bun:"proposal"`
	*models.Proposal

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveProposals(ctx context.Context, proposals ...*models.Proposal) (int64, error) {
	if len(proposals) == 0 {
		return 0, nil
	}

	var count int64
	for i := range proposals {
		if proposals[i].Type == "" {
			proposals[i].Type = storageTypes.ProposalTypeText
		}
		if proposals[i].Status == "" {
			proposals[i].Status = storageTypes.ProposalStatusInactive
		}

		add := addedProposal{
			Proposal: proposals[i],
		}

		query := tx.Tx().NewInsert().
			Column("id", "proposer_id", "height", "created_at", "deposit_time", "activation_time", "status", "type", "title", "description", "deposit", "metadata", "changes", "yes", "no", "no_with_veto", "abstain").
			Column("yes_vals", "no_vals", "no_with_veto_vals", "abstain_vals", "yes_addrs", "no_addrs", "no_with_veto_addrs", "abstain_addrs", "votes_count", "voting_power", "yes_voting_power", "no_voting_power", "no_with_veto_voting_power", "abstain_voting_power").
			Column("total_voting_power", "quorum", "veto_quorum", "threshold", "min_deposit", "end_time", "error").
			Model(&add).
			On("CONFLICT (id) DO UPDATE").
			Set("votes_count = added_proposal.votes_count + EXCLUDED.votes_count").
			Set("yes = added_proposal.yes + EXCLUDED.yes").
			Set("no = added_proposal.no + EXCLUDED.no").
			Set("no_with_veto = added_proposal.no_with_veto + EXCLUDED.no_with_veto").
			Set("abstain = added_proposal.abstain + EXCLUDED.abstain").
			Set("yes_vals = added_proposal.yes_vals + EXCLUDED.yes_vals").
			Set("no_vals = added_proposal.no_vals + EXCLUDED.no_vals").
			Set("no_with_veto_vals = added_proposal.no_with_veto_vals + EXCLUDED.no_with_veto_vals").
			Set("abstain_vals = added_proposal.abstain_vals + EXCLUDED.abstain_vals").
			Set("yes_addrs = added_proposal.yes_addrs + EXCLUDED.yes_addrs").
			Set("no_addrs = added_proposal.no_addrs + EXCLUDED.no_addrs").
			Set("no_with_veto_addrs = added_proposal.no_with_veto_addrs + EXCLUDED.no_with_veto_addrs").
			Set("abstain_addrs = added_proposal.abstain_addrs + EXCLUDED.abstain_addrs")

		if proposals[i].Deposit.IsPositive() {
			query.Set("deposit = added_proposal.deposit + EXCLUDED.deposit")
		}

		if !proposals[i].EmptyStatus() {
			query.Set("status = EXCLUDED.status")
		}

		if proposals[i].ActivationTime != nil {
			query.Set("activation_time = EXCLUDED.activation_time")
		}

		if proposals[i].EndTime != nil {
			query.Set("end_time = EXCLUDED.end_time")
		}

		if proposals[i].VotingPower.IsPositive() {
			query.Set("voting_power = EXCLUDED.voting_power")
		}
		if proposals[i].TotalVotingPower.IsPositive() {
			query.Set("total_voting_power = EXCLUDED.total_voting_power")
		}

		if proposals[i].Quorum != "" {
			query.Set("quorum = EXCLUDED.quorum")
		}
		if proposals[i].VetoQuorum != "" {
			query.Set("veto_quorum = EXCLUDED.veto_quorum")
		}
		if proposals[i].Threshold != "" {
			query.Set("threshold = EXCLUDED.threshold")
		}
		if proposals[i].MinDeposit != "" {
			query.Set("min_deposit = EXCLUDED.min_deposit")
		}
		if proposals[i].Error != "" {
			query.Set("error = EXCLUDED.error")
		}

		if proposals[i].YesVotingPower.IsPositive() {
			query.Set("yes_voting_power = EXCLUDED.yes_voting_power")
		}
		if proposals[i].NoVotingPower.IsPositive() {
			query.Set("no_voting_power = EXCLUDED.no_voting_power")
		}
		if proposals[i].NoWithVetoVotingPower.IsPositive() {
			query.Set("no_with_veto_voting_power = EXCLUDED.no_with_veto_voting_power")
		}
		if proposals[i].AbstainVotingPower.IsPositive() {
			query.Set("abstain_voting_power = EXCLUDED.abstain_voting_power")
		}

		if _, err := query.Returning("xmax, id").Exec(ctx); err != nil {
			return 0, err
		}

		if add.Xmax == 0 {
			count++
		}
	}

	return count, nil
}

var one = storageTypes.NumericFromInt64(1)

func (tx Transaction) SaveVotes(ctx context.Context, votes ...*models.Vote) (map[uint64]*models.VotesCount, error) {
	if len(votes) == 0 {
		return nil, nil
	}

	var votesCount = make(map[uint64]*models.VotesCount)
	for i := range votes {
		var existsVotes []models.Vote
		query := tx.Tx().NewSelect().
			Model(&existsVotes).
			Where("proposal_id = ?", votes[i].ProposalId)

		if votes[i].VoterId > 0 {
			query.Where("voter_id = ?", votes[i].VoterId)
		}
		if votes[i].ValidatorId != nil {
			query.Where("validator_id = ?", *votes[i].ValidatorId)
		}
		if err := query.Scan(ctx); err != nil {
			return nil, errors.Wrap(err, "receive existing votes")
		}

		if len(existsVotes) > 0 {
			ids := make([]uint64, len(existsVotes))
			totalWeight := votes[i].Weight.Copy()
			for j := range existsVotes {
				totalWeight = totalWeight.Add(existsVotes[j].Weight)
				ids[j] = existsVotes[j].Id
			}
			if totalWeight.GreaterThan(one) {
				if _, err := tx.Tx().NewDelete().Model((*models.Vote)(nil)).Where("id IN ?", bun.Tuple(ids)).Exec(ctx); err != nil {
					return nil, errors.Wrap(err, "remove existing votes")
				}

				for _, vote := range existsVotes {
					if vc, ok := votesCount[vote.ProposalId]; ok {
						vc.Update(-1, vote)
					} else {
						var vc models.VotesCount
						vc.Update(-1, vote)
						votesCount[vote.ProposalId] = &vc
					}
				}
			}
		}

		if vc, ok := votesCount[votes[i].ProposalId]; ok {
			vc.Update(1, *votes[i])
		} else {
			var vc models.VotesCount
			vc.Update(1, *votes[i])
			votesCount[votes[i].ProposalId] = &vc
		}
	}

	_, err := tx.Tx().NewInsert().Model(&votes).Exec(ctx)
	return votesCount, err
}

func (tx Transaction) Proposal(ctx context.Context, id uint64) (proposal models.Proposal, err error) {
	err = tx.Tx().NewSelect().Model(&proposal).
		Where("id = ?", id).
		Column("id", "changes", "type").
		Scan(ctx)
	return
}

func (tx Transaction) UpdateSignalsAfterUpgrade(ctx context.Context, version uint64) (storageTypes.Numeric, error) {
	_, err := tx.Tx().NewUpdate().Table("signal_version", "validator").
		SetColumn("voting_power", "validator.stake").
		Where("signal_version.version = ?", version).
		Where("validator.id = validator_id").
		Exec(ctx)
	if err != nil {
		return storageTypes.NumericZero(), err
	}

	var sum storageTypes.Numeric
	err = tx.Tx().NewSelect().
		TableExpr("(?) AS latest",
			tx.Tx().NewSelect().
				Table("signal_version").
				ColumnExpr("DISTINCT ON (validator_id) voting_power").
				Where("version = ?", version).
				OrderExpr("validator_id, height DESC"),
		).
		ColumnExpr("COALESCE(SUM(voting_power), 0)").
		Scan(ctx, &sum)
	return sum, err
}
