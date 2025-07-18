// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	paramsV1Beta "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stoewer/go-strcase"
)

func (module *Module) saveProposals(
	ctx context.Context,
	tx storage.Transaction,
	height pkgTypes.Level,
	proposals []*storage.Proposal,
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
					votes[i].ValidatorId = validatorId
				}

				for j := range proposals {
					if proposals[j].Id == votes[i].ProposalId {
						if votes[i].ValidatorId > 0 {
							switch votes[i].Option {
							case types.VoteOptionAbstain:
								proposals[j].AbstainValidators += 1
							case types.VoteOptionNo:
								proposals[j].NoValidators += 1
							case types.VoteOptionNoWithVeto:
								proposals[j].NoWithVetoValidators += 1
							case types.VoteOptionYes:
								proposals[j].YesValidators += 1
							}
						} else {
							switch votes[i].Option {
							case types.VoteOptionAbstain:
								proposals[j].AbstainAddress += 1
							case types.VoteOptionNo:
								proposals[j].NoAddress += 1
							case types.VoteOptionNoWithVeto:
								proposals[j].NoWithVetoAddress += 1
							case types.VoteOptionYes:
								proposals[j].YesAddress += 1
							}
						}
						break
					}
				}
			}
		}

		if err := tx.SaveVotes(ctx, votes...); err != nil {
			return 0, errors.Wrap(err, "save votes")
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

		if err := module.updateConstants(ctx, tx, filled[i]); err != nil {
			return 0, errors.Wrap(err, "update constants")
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

func (module *Module) fillProposalsVotingPower(ctx context.Context, tx storage.Transaction, height pkgTypes.Level, changedProposals []*storage.Proposal) ([]*storage.Proposal, error) {
	// 1. Receive all active or just completed proposals

	// 1.1 Return if we don't have proposal updates and it's not certain block height (one block in hour)

	proposals := make(map[uint64]*storage.Proposal)
	for i := range changedProposals {
		if !changedProposals[i].Finished() {
			continue
		}
		proposals[changedProposals[i].Id] = changedProposals[i]
	}

	if len(proposals) == 0 && height%600 > 0 {
		return changedProposals, nil
	}

	active, err := tx.ActiveProposals(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get active proposals")
	}

	for i := range active {
		if _, ok := proposals[active[i].Id]; !ok {
			// reset counters to avoid repeat of counting
			active[i].Abstain = 0
			active[i].AbstainAddress = 0
			active[i].AbstainValidators = 0
			active[i].No = 0
			active[i].NoAddress = 0
			active[i].NoValidators = 0
			active[i].NoWithVeto = 0
			active[i].NoWithVetoAddress = 0
			active[i].NoWithVetoValidators = 0
			active[i].Yes = 0
			active[i].YesAddress = 0
			active[i].YesValidators = 0
			active[i].VotesCount = 0
			active[i].Deposit = decimal.Zero

			proposals[active[i].Id] = &active[i]
		}
	}

	if len(active) == 0 && len(proposals) == 0 {
		return changedProposals, nil
	}

	// 2. Get all validators

	maxValsConsts, err := module.constants.Get(ctx, types.ModuleNameStaking, "max_validators")
	if err != nil {
		return nil, errors.Wrap(err, "get max validators value")
	}
	maxVals, err := strconv.Atoi(maxValsConsts.Value)
	if err != nil {
		return nil, errors.Wrap(err, "parse max validators value")
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

	totalVotingPower, err := module.validators.TotalVotingPower(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get total voting power")
	}

	for _, proposal := range proposals {
		proposal.VotingPower = decimal.Zero
		proposal.AbstainVotingPower = decimal.Zero
		proposal.NoVotingPower = decimal.Zero
		proposal.NoWithVetoVotingPower = decimal.Zero
		proposal.YesVotingPower = decimal.Zero
		proposal.TotalVotingPower = decimal.Zero

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
				delegations, err := tx.AddressDelegations(ctx, votes[i].VoterId)
				if err != nil {
					return nil, errors.Wrapf(err, "can't receive address delegations: %d", votes[i].VoterId)
				}

				for j := range delegations {
					if amount, ok := validatorMinus[delegations[j].ValidatorId]; ok {
						validatorMinus[delegations[j].ValidatorId] = amount.Add(delegations[j].Amount)
					} else {
						validatorMinus[delegations[j].ValidatorId] = delegations[j].Amount
					}
					proposal.VotingPower = proposal.VotingPower.Add(delegations[j].Amount)

					switch votes[i].Option {
					case types.VoteOptionAbstain:
						proposal.AbstainVotingPower = proposal.AbstainVotingPower.Add(delegations[j].Amount)
					case types.VoteOptionNo:
						proposal.NoVotingPower = proposal.NoVotingPower.Add(delegations[j].Amount)
					case types.VoteOptionNoWithVeto:
						proposal.NoWithVetoVotingPower = proposal.NoWithVetoVotingPower.Add(delegations[j].Amount)
					case types.VoteOptionYes:
						proposal.YesVotingPower = proposal.YesVotingPower.Add(delegations[j].Amount)
					}
				}

				if votes[i].ValidatorId > 0 {
					votedValidators[votes[i].ValidatorId] = votes[i].Option
				}
			}
		}

		for id, option := range votedValidators {
			if power, ok := validatorsPower[id]; ok {
				minus := validatorMinus[id]
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

	result := make([]*storage.Proposal, 0)
	for _, proposal := range proposals {
		result = append(result, proposal)
	}

	return result, nil
}

func (module *Module) updateConstants(ctx context.Context, tx storage.Transaction, proposal *storage.Proposal) error {
	changes, err := tx.Proposal(ctx, proposal.Id)
	if err != nil {
		if module.validators.IsNoRows(err) {
			return nil
		}
		return errors.Wrap(err, "receive proposal changes")
	}

	// save only constants from applied param change proposals
	if proposal.Status != types.ProposalStatusApplied || changes.Type != types.ProposalTypeParamChanged {
		return nil
	}

	var parsed []paramsV1Beta.ParamChange
	if err := json.Unmarshal(changes.Changes, &parsed); err != nil {
		return errors.Wrap(err, "parse proposal changes")
	}

	constants := make([]storage.Constant, 0)
	for i := range parsed {
		handledConstants, err := constantsHandle(parsed[i])
		if err != nil {
			return errors.Wrap(err, "handle proposal changes")
		}
		constants = append(constants, handledConstants...)
	}

	if len(constants) == 0 {
		return nil
	}

	return tx.SaveConstants(ctx, constants...)
}

func constantsHandle(change paramsV1Beta.ParamChange) ([]storage.Constant, error) {
	moduleName, err := types.ParseModuleName(change.GetSubspace())
	if err != nil {
		return nil, errors.Wrapf(err, "parsing module name in proposal changes: %s", change.GetSubspace())
	}
	key := change.GetKey()
	value := change.GetValue()

	constants := make([]storage.Constant, 0)

	switch moduleName {
	case types.ModuleNameConsensus, types.ModuleNameBaseapp:

		switch key {
		case "BlockParams":
			c, err := parseParamsToConstants(types.ModuleNameConsensus, "block_", value)
			if err != nil {
				return nil, errors.Wrap(err, "parse block params")
			}
			constants = append(constants, c...)
		case "EvidenceParams":
			c, err := parseParamsToConstants(types.ModuleNameConsensus, "evidence_", value)
			if err != nil {
				return nil, errors.Wrap(err, "parse evidence params")
			}
			constants = append(constants, c...)
		case "ValidatorParams":
			c, err := parseParamsToConstants(types.ModuleNameConsensus, "validator_", value)
			if err != nil {
				return nil, errors.Wrap(err, "parse validator params")
			}
			constants = append(constants, c...)
		}

	case types.ModuleNameGov:

		if key == "votingparams" {
			c, err := parseParamsToConstants(moduleName, "", value)
			if err != nil {
				return nil, errors.Wrap(err, "parse voting params")
			}
			constants = append(constants, c...)
		}

	case types.ModuleNameBlob:
		val := value
		if key == "GovMaxSquareSize" {
			val, err = strconv.Unquote(value)
			if err != nil {
				return nil, errors.Wrap(err, value)
			}
		}

		constants = append(constants, storage.Constant{
			Module: moduleName,
			Name:   strcase.SnakeCase(key),
			Value:  val,
		})
	default:

		constants = append(constants, storage.Constant{
			Module: moduleName,
			Name:   strcase.SnakeCase(key),
			Value:  value,
		})

	}

	return constants, nil
}

func parseParamsToConstants(moduleName types.ModuleName, keyPrefix, value string) ([]storage.Constant, error) {
	var params map[string]string
	if err := json.Unmarshal([]byte(value), &params); err != nil {
		return nil, errors.Wrap(err, "unmarshal params")
	}
	constants := make([]storage.Constant, 0)
	for k, v := range params {
		constants = append(constants, storage.Constant{
			Module: moduleName,
			Name:   keyPrefix + k,
			Value:  v,
		})
	}
	return constants, nil
}
