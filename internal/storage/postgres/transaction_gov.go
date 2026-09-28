// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
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
			Column("version", "height", "time", "end_height", "end_time", "applied_at_level", "applied_at",
				"expected_upgrade_height", "signer_id", "msg_id", "tx_id", "voting_power", "voted_power",
				"signals_count", "status").
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
		if upgrades[i].ExpectedHeight > 0 {
			// write-once: the MsgTryUpgrade forecast is kept when the upgrade is applied
			query = query.Set("expected_upgrade_height = COALESCE(NULLIF(upgrade.expected_upgrade_height, 0), EXCLUDED.expected_upgrade_height)")
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

// FixSignalsPower closes the open round: counted signals keep the power they contributed, the rest get 0.
// Updates go row by row: a round has at most one signal per validator.
func (tx Transaction) FixSignalsPower(ctx context.Context, powers map[uint64]storageTypes.Numeric) error {
	if len(powers) > 0 {
		signals := make([]models.SignalVersion, 0, len(powers))
		for id, power := range powers {
			signals = append(signals, models.SignalVersion{
				Id:          id,
				VotingPower: power,
			})
		}

		if _, err := tx.Tx().NewUpdate().
			With("_data", tx.Tx().NewValues(&signals).Column("id", "voting_power")).
			Model((*models.SignalVersion)(nil)).
			TableExpr("_data").
			Set("voting_power = _data.voting_power").
			Where("signal_version.id = _data.id").
			Where("signal_version.voting_power IS NULL").
			Exec(ctx); err != nil {
			return errors.Wrap(err, "fix counted signals")
		}
	}

	// the round is every signal still open; older rounds are closed already
	if _, err := tx.Tx().NewUpdate().
		Table("signal_version").
		Set("voting_power = 0").
		Where("voting_power IS NULL").
		Exec(ctx); err != nil {
		return errors.Wrap(err, "zero not counted signals")
	}
	return nil
}

// RollbackUpgrades reverts what a block wrote over existing upgrades and signals.
// It must run before the block's signals and upgrades are deleted by height.
func (tx Transaction) RollbackUpgrades(ctx context.Context, height pkgTypes.Level) error {
	// signals only counted for versions above the block's app version
	if _, err := tx.Tx().ExecContext(ctx, `
		UPDATE upgrade SET signals_count = upgrade.signals_count - c.cnt
		FROM (
			SELECT version, count(*) AS cnt FROM signal_version
			WHERE height = ? AND version > (SELECT version_app FROM block WHERE height = ?)
			GROUP BY version
		) AS c
		WHERE upgrade.version = c.version
	`, height, height); err != nil {
		return errors.Wrap(err, "rollback signals count")
	}

	// reopen the round closed by the block: by MsgTryUpgrade or, without it, by the applied upgrade;
	// plain queries, as TimescaleDB fails to plan subqueries in an UPDATE of a hypertable
	closed, err := tx.Tx().NewSelect().
		Model((*models.Upgrade)(nil)).
		WhereOr("end_height = ?", height).
		WhereOr("applied_at_level = ? AND end_height = 0", height).
		Exists(ctx)
	if err != nil {
		return errors.Wrap(err, "find closed round")
	}
	if closed {
		var roundStart pkgTypes.Level
		if err := tx.Tx().NewSelect().
			Model((*models.Upgrade)(nil)).
			ColumnExpr("COALESCE(MAX(applied_at_level), 0)").
			Where("status = ?", storageTypes.UpgradeStatusApplied).
			Where("applied_at_level < ?", height).
			Scan(ctx, &roundStart); err != nil {
			return errors.Wrap(err, "find round start")
		}
		if _, err := tx.Tx().NewUpdate().
			Model((*models.SignalVersion)(nil)).
			Set("voting_power = NULL").
			Where("height >= ?", roundStart).
			Where("height < ?", height).
			Exec(ctx); err != nil {
			return errors.Wrap(err, "reopen signals round")
		}
	}

	if _, err := tx.Tx().NewUpdate().
		Model((*models.Upgrade)(nil)).
		Set("end_height = 0").
		Set("end_time = ?", time.Time{}).
		Set("expected_upgrade_height = NULL").
		Set("signer_id = 0").
		Set("msg_id = 0").
		Set("tx_id = 0").
		Set("status = ?", storageTypes.UpgradeStatusProcessing).
		Where("end_height = ?", height).
		Exec(ctx); err != nil {
		return errors.Wrap(err, "rollback try upgrade")
	}

	// rows created by setUpgradeApplied alone have no signal height
	if _, err := tx.Tx().NewDelete().
		Model((*models.Upgrade)(nil)).
		Where("applied_at_level = ?", height).
		Where("height = 0").
		Exec(ctx); err != nil {
		return errors.Wrap(err, "delete applied upgrade")
	}
	if _, err := tx.Tx().NewUpdate().
		Model((*models.Upgrade)(nil)).
		Set("applied_at_level = 0").
		Set("applied_at = ?", time.Time{}).
		// without MsgTryUpgrade the height was taken from the applied block
		Set("expected_upgrade_height = CASE WHEN end_height > 0 THEN expected_upgrade_height ELSE NULL END").
		Set("status = CASE WHEN end_height > 0 THEN ?::upgrade_status ELSE ?::upgrade_status END",
			storageTypes.UpgradeStatusWaitingUpgrade, storageTypes.UpgradeStatusProcessing).
		Where("applied_at_level = ?", height).
		Exec(ctx); err != nil {
		return errors.Wrap(err, "rollback applied upgrade")
	}
	return nil
}

// UpdateUpgradeTally writes a recounted tally; unlike SaveUpgrades it stores zero power and lowers the status back.
func (tx Transaction) UpdateUpgradeTally(ctx context.Context, version uint64, votingPower, votedPower storageTypes.Numeric, status storageTypes.UpgradeStatus) error {
	_, err := tx.Tx().NewUpdate().
		Model((*models.Upgrade)(nil)).
		Set("voting_power = ?", votingPower).
		Set("voted_power = ?", votedPower).
		Set("status = ?", status).
		Where("version = ?", version).
		Where("(voting_power, voted_power, status) IS DISTINCT FROM (?, ?, ?)", votingPower, votedPower, status).
		Exec(ctx)
	return err
}
