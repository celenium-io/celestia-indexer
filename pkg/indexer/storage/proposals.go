// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"iter"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

func appliedProposals(it iter.Seq[*storage.Proposal]) iter.Seq[*storage.Proposal] {
	return func(yield func(*storage.Proposal) bool) {
		for p := range it {
			if p.Status == types.ProposalStatusApplied {
				if !yield(p) {
					return
				}
			}
		}
	}
}

func (module *Module) saveProposals(
	ctx context.Context,
	tx storage.Transaction,
	repos storage.TxRepos,
	height pkgTypes.Level,
	proposals *sdkSync.Map[uint64, *storage.Proposal],
	votes []*storage.Vote,
	addrToId map[string]uint64,
) (int64, error) {
	if len(votes) > 0 {
		for i := range votes {
			if votes[i].Voter != nil {
				voterId, ok := addrToId[votes[i].Voter.Address]
				if !ok {
					return 0, errors.Wrap(errCantFindAddress, votes[i].Voter.Address)
				}
				votes[i].VoterId = voterId

				if validatorId, ok := module.validatorsByDelegator[votes[i].Voter.Address]; ok {
					votes[i].ValidatorId = &validatorId
				}
			} else {
				return 0, errors.Errorf("nil voter address")
			}
		}

		votesCount, err := tx.SaveVotes(ctx, votes...)
		if err != nil {
			return 0, errors.Wrap(err, "save votes")
		}

		for id, vc := range votesCount {
			proposal, ok := proposals.Get(id)
			if !ok {
				return 0, errors.Errorf("unknown proposal id during votes count computing: %d", id)
			}

			proposal.VotesCount += vc.VotesCount

			proposal.Abstain += vc.Abstain
			proposal.No += vc.No
			proposal.NoWithVeto += vc.NoWithVeto
			proposal.Yes += vc.Yes

			proposal.AbstainAddress += vc.AbstainAddress
			proposal.NoAddress += vc.NoAddress
			proposal.NoWithVetoAddress += vc.NoWithVetoAddress
			proposal.YesAddress += vc.YesAddress

			proposal.AbstainValidators += vc.AbstainValidators
			proposal.NoValidators += vc.NoValidators
			proposal.NoWithVetoValidators += vc.NoWithVetoValidators
			proposal.YesValidators += vc.YesValidators
		}
	}

	if err := fillExpeditedProposal(ctx, repos.Proposals, repos.Constants, proposals); err != nil {
		return 0, errors.Wrap(err, "fill expedited proposal")
	}

	filled, err := module.fillProposalsVotingPower(ctx, repos, height, proposals)
	if err != nil {
		return 0, errors.Wrap(err, "compute proposal shares")
	}

	for i := range filled {
		if filled[i].Proposer != nil {
			proposerId, ok := addrToId[filled[i].Proposer.Address]
			if !ok {
				return 0, errors.Wrap(errCantFindAddress, filled[i].Proposer.Address)
			}
			filled[i].ProposerId = proposerId
		}

		if !filled[i].CreatedAt.IsZero() {
			duration, err := module.getConstantDuration(ctx, repos, types.ModuleNameGov, "max_deposit_period")
			if err != nil {
				return 0, errors.Wrap(err, "getConstantDuration")
			}
			filled[i].DepositTime = filled[i].CreatedAt.Add(duration)
		}

		if filled[i].ActivationTime != nil && filled[i].EndTime == nil {
			duration, err := module.getConstantDuration(ctx, repos, types.ModuleNameGov, "voting_period")
			if err != nil {
				return 0, errors.Wrap(err, "getConstantDuration")
			}
			endTime := filled[i].ActivationTime.Add(duration)
			filled[i].EndTime = &endTime
		}
	}

	return tx.SaveProposals(ctx, filled...)
}

func (module *Module) getConstantDuration(
	ctx context.Context, repos storage.TxRepos, moduleName types.ModuleName, name string,
) (time.Duration, error) {
	constant, err := repos.Constants.Get(ctx, moduleName, name)
	if err != nil {
		return 0, errors.Wrapf(err, "can't find %s constant", name)
	}
	intValue, err := strconv.ParseInt(constant.Value, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "can't parse %s value", name)
	}
	return time.Duration(intValue), nil
}

func (module *Module) fillProposalsVotingPower(
	ctx context.Context,
	repos storage.TxRepos,
	height pkgTypes.Level,
	proposals *sdkSync.Map[uint64, *storage.Proposal],
) ([]*storage.Proposal, error) {
	// 1. Receive all active or just completed proposals

	// 1.1 Return if we don't have proposal updates and it's not certain block height (one block in hour)

	finished := make(map[uint64]*storage.Proposal)

	for value := range proposals.AllValues() {
		if value.Finished() {
			finished[value.Id] = value
		}
	}

	if len(finished) == 0 && height%600 > 0 {
		return proposals.Values(), nil
	}

	active, err := repos.Proposals.Active(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get active proposals")
	}

	for i := range active {
		if proposal, ok := finished[active[i].Id]; !ok {
			if p, ok := proposals.Get(active[i].Id); ok {
				finished[active[i].Id] = p
			} else {
				finished[active[i].Id] = &active[i]
				proposals.Set(active[i].Id, finished[active[i].Id])
			}
		} else {
			proposal.Expedited = active[i].Expedited
		}
	}

	if len(finished) == 0 {
		return proposals.Values(), nil
	}

	// 2. Get all validators

	validators, err := repos.Validators.BondedValidators(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get validators")
	}
	validatorsPower := make(map[uint64]types.Numeric)
	for i := range validators {
		validatorsPower[validators[i].Id] = validators[i].Stake
	}

	// 3. Compute voting results

	const limit = 1000

	totalVotingPower, err := repos.Validators.TotalVotingPower(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get total voting power")
	}

	for _, proposal := range finished {
		validatorMinus := make(map[uint64]types.Numeric)
		votedValidators := make(map[uint64]types.VoteOption)

		if proposal.Finished() {
			proposal.TotalVotingPower = totalVotingPower

			quorum, err := repos.Constants.Get(ctx, types.ModuleNameGov, "quorum")
			if err != nil {
				return nil, errors.Wrapf(err, "can't find quorum constant")
			}
			proposal.Quorum = quorum.Value

			var thresholdConstantName = "threshold"
			if proposal.IsExpedited() {
				thresholdConstantName = "expedited_threshold"
			}
			threshold, err := repos.Constants.Get(ctx, types.ModuleNameGov, thresholdConstantName)
			if err != nil {
				return nil, errors.Wrapf(err, "can't find threshold constant")
			}
			proposal.Threshold = threshold.Value

			veto, err := repos.Constants.Get(ctx, types.ModuleNameGov, "veto_threshold")
			if err != nil {
				return nil, errors.Wrapf(err, "can't find veto_threshold constant")
			}
			proposal.VetoQuorum = veto.Value
		}

		paginate := sdkSync.Paginate(
			ctx, limit,
			func(ctx context.Context, limit, offset int) ([]storage.Vote, error) {
				return repos.Votes.ListByProposal(ctx, proposal.Id, limit, offset)
			},
		)

		for vote, err := range paginate {
			if err != nil {
				return nil, errors.Wrapf(err, "get proposal votes: proposal_id=%d", proposal.Id)
			}

			if vote.ValidatorId != nil {
				votedValidators[*vote.ValidatorId] = vote.Option
			}

			delegations, err := repos.Delegation.AddressDelegations(ctx, vote.VoterId)
			if err != nil {
				return nil, errors.Wrapf(err, "can't receive address delegations: %d", vote.VoterId)
			}

			for j := range delegations {
				if vote.ValidatorId != nil && delegations[j].ValidatorId == *vote.ValidatorId && delegations[j].AddressId == vote.VoterId {
					// skip self delegation
					continue
				}

				shares := delegations[j].Amount
				if amount, ok := validatorMinus[delegations[j].ValidatorId]; ok {
					validatorMinus[delegations[j].ValidatorId] = amount.Add(shares)
				} else {
					validatorMinus[delegations[j].ValidatorId] = shares
				}
				proposal.VotingPower = proposal.VotingPower.Add(shares)

				switch vote.Option {
				case types.VoteOptionAbstain:
					proposal.AbstainVotingPower = proposal.AbstainVotingPower.Add(shares)
				case types.VoteOptionNo:
					proposal.NoVotingPower = proposal.NoVotingPower.Add(shares)
				case types.VoteOptionNoWithVeto:
					proposal.NoWithVetoVotingPower = proposal.NoWithVetoVotingPower.Add(shares)
				case types.VoteOptionYes:
					proposal.YesVotingPower = proposal.YesVotingPower.Add(shares)
				}
			}
		}

		for id, option := range votedValidators {
			minus, ok := validatorMinus[id]
			if !ok {
				minus = types.NumericZero()
			}
			if power, ok := validatorsPower[id]; ok {
				proposal.VotingPower = proposal.VotingPower.Add(power).Sub(minus)

				switch option {
				case types.VoteOptionAbstain:
					proposal.AbstainVotingPower = proposal.AbstainVotingPower.Add(power).Sub(minus)
				case types.VoteOptionNo:
					proposal.NoVotingPower = proposal.NoVotingPower.Add(power).Sub(minus)
				case types.VoteOptionNoWithVeto:
					proposal.NoWithVetoVotingPower = proposal.NoWithVetoVotingPower.Add(power).Sub(minus)
				case types.VoteOptionYes:
					proposal.YesVotingPower = proposal.YesVotingPower.Add(power).Sub(minus)
				}
			}
		}
	}

	return proposals.Values(), nil
}

func fillExpeditedProposal(
	ctx context.Context,
	repo storage.IProposal,
	constants storage.IConstant,
	proposals *sdkSync.Map[uint64, *storage.Proposal],
) error {
	for proposal := range proposals.AllValues() {
		if proposal.Status != types.ProposalStatusActive {
			continue
		}
		if proposal.ActivationTime == nil && !proposal.ExpeditedProposalRejected {
			continue
		}

		stored, err := repo.GetByID(ctx, proposal.Id)
		if err != nil && !repo.IsNoRows(err) {
			return errors.Wrapf(err, "get proposal by id: %d", proposal.Id)
		}

		if proposal.ExpeditedProposalRejected {
			if stored == nil {
				return errors.Errorf("unknown expedited proposal: %d", proposal.Id)
			}
			if err := fillExpeditedProposalEndTime(
				ctx, constants, "voting_period", stored.ActivationTime, proposal,
			); err != nil {
				return errors.Wrap(err, "can't fill expedited proposal")
			}
			proposal.Expedited = new(false)
			continue
		}

		expedited := stored != nil && stored.IsExpedited()
		if proposal.Expedited != nil {
			expedited = *proposal.Expedited
		}

		if err := fillProposalMinDeposit(ctx, constants, proposal, expedited); err != nil {
			return err
		}

		if expedited {
			if err := fillExpeditedProposalEndTime(
				ctx, constants, "expedited_voting_period", proposal.ActivationTime, proposal,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func computeEndTimeForExpeditionProposal(activationTime *time.Time, votingPeriod string) (*time.Time, error) {
	if activationTime == nil {
		return nil, errors.Errorf("nil activation time")
	}
	ns, err := strconv.ParseInt(votingPeriod, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "parse voting period: %s", votingPeriod)
	}
	endTime := activationTime.Add(time.Duration(ns))
	return &endTime, nil
}

func fillExpeditedProposalEndTime(
	ctx context.Context, constants storage.IConstant, name string, activationTime *time.Time, proposal *storage.Proposal,
) error {
	votingPeriod, err := constants.Get(ctx, types.ModuleNameGov, name)
	if err != nil {
		return errors.Wrapf(err, "get '%s'", name)
	}
	endTime, err := computeEndTimeForExpeditionProposal(activationTime, votingPeriod.Value)
	if err != nil {
		return errors.Wrapf(err, "compute end time for expedition proposal: %d", proposal.Id)
	}
	proposal.EndTime = endTime
	return nil
}

func fillProposalMinDeposit(ctx context.Context, constants storage.IConstant, proposal *storage.Proposal, expedited bool) error {
	name := "min_deposit"
	if expedited {
		name = "expedited_min_deposit"
	}
	minDeposit, err := constants.Get(ctx, types.ModuleNameGov, name)
	if err != nil {
		return errors.Wrapf(err, "get '%s'", name)
	}
	proposal.MinDeposit = minDeposit.Value
	return nil
}
