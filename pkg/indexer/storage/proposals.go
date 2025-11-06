// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"maps"
	"slices"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (module *Module) saveProposals(
	ctx context.Context,
	tx storage.Transaction,
	height pkgTypes.Level,
	proposals *sync.Map[uint64, *storage.Proposal],
	votes []*storage.Vote,
	addrToId map[string]uint64,
) (int64, error) {
	if len(votes) > 0 {
		for i := range votes {
			if votes[i].Voter != nil {
				voterId, ok := addrToId[votes[i].Voter.Address]
				if !ok {
					return 0, errors.Errorf("unknown voter address: %s", votes[i].Voter.Address)
				}
				votes[i].VoterId = voterId

				if validatorId, ok := module.validatorsByDelegator[votes[i].Voter.Address]; ok {
					votes[i].ValidatorId = &validatorId
				}
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

	filled, err := module.fillProposalsVotingPower(ctx, tx, height, proposals)
	if err != nil {
		return 0, errors.Wrap(err, "compute proposal shares")
	}

	for i := range filled {
		if filled[i].Proposer != nil {
			proposerId, ok := addrToId[filled[i].Proposer.Address]
			if !ok {
				return 0, errors.Errorf("unknown proposer address for proposal: %s", filled[i].Proposer.Address)
			}
			filled[i].ProposerId = proposerId
		}

		if !filled[i].CreatedAt.IsZero() {
			duration, err := module.getConstantDuration(ctx, types.ModuleNameGov, "max_deposit_period")
			if err != nil {
				return 0, errors.Wrap(err, "getConstantDuration")
			}
			filled[i].DepositTime = filled[i].CreatedAt.Add(duration)
		}

		if filled[i].ActivationTime != nil && filled[i].EndTime == nil {
			duration, err := module.getConstantDuration(ctx, types.ModuleNameGov, "voting_period")
			if err != nil {
				return 0, errors.Wrap(err, "getConstantDuration")
			}
			endTime := filled[i].ActivationTime.Add(duration)
			filled[i].EndTime = &endTime
		}
	}

	return tx.SaveProposals(ctx, filled...)
}

func (module *Module) getConstantDuration(ctx context.Context, moduleName types.ModuleName, name string) (time.Duration, error) {
	constant, err := module.constants.Get(ctx, moduleName, name)
	if err != nil {
		return 0, errors.Wrapf(err, "can't find %s constant", name)
	}
	intValue, err := strconv.ParseInt(constant.Value, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "can't parse %s value", name)
	}
	return time.Duration(intValue), nil
}

func (module *Module) fillProposalsVotingPower(ctx context.Context, tx storage.Transaction, height pkgTypes.Level, proposals *sync.Map[uint64, *storage.Proposal]) ([]*storage.Proposal, error) {
	// 1. Receive all active or just completed proposals

	// 1.1 Return if we don't have proposal updates and it's not certain block height (one block in hour)

	finished := make(map[uint64]*storage.Proposal)
	if err := proposals.Range(func(key uint64, value *storage.Proposal) (error, bool) {
		if value.Finished() {
			finished[value.Id] = value
		}
		return nil, false
	}); err != nil {
		return nil, errors.Wrap(err, "proposals range")
	}

	if len(finished) == 0 && height%600 > 0 {
		return proposals.Values(), nil
	}

	active, err := tx.ActiveProposals(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get active proposals")
	}

	for i := range active {
		if _, ok := finished[active[i].Id]; !ok {
			finished[active[i].Id] = &active[i]
		}
	}

	if len(finished) == 0 {
		return proposals.Values(), nil
	}

	// 2. Get all validators

	maxVals, err := getMaxValidatorsCount(ctx, module.constants)
	if err != nil {
		return nil, errors.Wrapf(err, "receiving max validators count")
	}

	validators, err := tx.BondedValidators(ctx, maxVals)
	if err != nil {
		return nil, errors.Wrap(err, "get validators")
	}
	validatorsPower := make(map[uint64]decimal.Decimal)
	for i := range validators {
		validatorsPower[validators[i].Id] = validators[i].Stake
	}

	// 3. Compute voting results

	const limit = 1000

	totalVotingPower, err := module.validators.TotalVotingPower(ctx, maxVals)
	if err != nil {
		return nil, errors.Wrap(err, "get total voting power")
	}

	for _, proposal := range finished {
		validatorMinus := make(map[uint64]decimal.Decimal)
		votedValidators := make(map[uint64]types.VoteOption)

		if proposal.Finished() {
			proposal.TotalVotingPower = totalVotingPower

			quorum, err := module.constants.Get(ctx, types.ModuleNameGov, "quorum")
			if err != nil {
				return nil, errors.Wrapf(err, "can't find quorum constant")
			}
			proposal.Quorum = quorum.Value

			minDeposit, err := module.constants.Get(ctx, types.ModuleNameGov, "min_deposit")
			if err != nil {
				return nil, errors.Wrapf(err, "can't find min_deposit constant")
			}
			proposal.MinDeposit = minDeposit.Value

			threshold, err := module.constants.Get(ctx, types.ModuleNameGov, "threshold")
			if err != nil {
				return nil, errors.Wrapf(err, "can't find threshold constant")
			}
			proposal.Threshold = threshold.Value

			veto, err := module.constants.Get(ctx, types.ModuleNameGov, "veto_threshold")
			if err != nil {
				return nil, errors.Wrapf(err, "can't find veto_threshold constant")
			}
			proposal.VetoQuorum = veto.Value
		}

		var offset int
		var end bool

		for !end {
			votes, err := tx.ProposalVotes(ctx, proposal.Id, limit, offset)
			if err != nil {
				return nil, errors.Wrapf(err, "get proposal votes: proposal_id=%d", proposal.Id)
			}
			offset += limit
			end = len(votes) < limit

			for i := range votes {
				if votes[i].ValidatorId != nil {
					votedValidators[*votes[i].ValidatorId] = votes[i].Option
				}

				delegations, err := tx.AddressDelegations(ctx, votes[i].VoterId)
				if err != nil {
					return nil, errors.Wrapf(err, "can't receive address delegations: %d", votes[i].VoterId)
				}

				for j := range delegations {
					if votes[i].ValidatorId != nil && delegations[j].ValidatorId == *votes[i].ValidatorId && delegations[j].AddressId == votes[i].VoterId {
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

					switch votes[i].Option {
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
		}

		for id, option := range votedValidators {
			minus, ok := validatorMinus[id]
			if !ok {
				minus = decimal.Zero
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

	return slices.Collect(maps.Values(finished)), nil
}
