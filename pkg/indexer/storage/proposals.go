// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/json"
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
								proposals[i].AbstainValidators += 1
							case types.VoteOptionNo:
								proposals[i].NoValidators += 1
							case types.VoteOptionNoWithVeto:
								proposals[i].NoWithVetoValidators += 1
							case types.VoteOptionYes:
								proposals[i].YesValidators += 1
							}
						} else {
							switch votes[i].Option {
							case types.VoteOptionAbstain:
								proposals[i].AbstainAddress += 1
							case types.VoteOptionNo:
								proposals[i].NoAddress += 1
							case types.VoteOptionNoWithVeto:
								proposals[i].NoWithVetoAddress += 1
							case types.VoteOptionYes:
								proposals[i].YesAddress += 1
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
			constant, err := module.constants.Get(ctx, types.ModuleNameGov, "max_deposit_period")
			if err != nil {
				return 0, errors.Wrap(err, "can't find max_deposit_period constant")
			}
			maxDepositPeriod, err := time.ParseDuration(constant.Value)
			if err != nil {
				return 0, errors.Wrap(err, "can't parse max_deposit_period value")
			}
			filled[i].DepositTime = filled[i].CreatedAt.Add(maxDepositPeriod)
		}

		if err := module.updateConstants(ctx, tx, filled[i]); err != nil {
			return 0, errors.Wrap(err, "update constants")
		}
	}

	return tx.SaveProposals(ctx, filled...)
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

	validators, err := tx.Validators(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get validators")
	}
	validatorsPower := make(map[uint64]decimal.Decimal)
	for i := range validators {
		validatorsPower[validators[i].Id] = validators[i].Stake
	}

	// 3. Compute voting results

	const limit = 1000

	for _, proposal := range proposals {
		proposal.VotingPower = decimal.Zero
		proposal.AbstainVotingPower = decimal.Zero
		proposal.NoVotingPower = decimal.Zero
		proposal.NoWithVetoVotingPower = decimal.Zero
		proposal.YesVotingPower = decimal.Zero

		validatorMinus := make(map[uint64]decimal.Decimal)
		votedValidators := make(map[uint64]types.VoteOption)

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
				if votes[i].ValidatorId > 0 {
					votedValidators[votes[i].ValidatorId] = votes[i].Option
					continue
				}
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
			} else {
				return nil, errors.Errorf("unknown validator id: %d", id)
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
			c, err := parseParamsToConstants(moduleName, "block_", value)
			if err != nil {
				return nil, errors.Wrap(err, "parse block params")
			}
			constants = append(constants, c...)
		case "EvidenceParams":
			c, err := parseParamsToConstants(moduleName, "evidence_", value)
			if err != nil {
				return nil, errors.Wrap(err, "parse evidence params")
			}
			constants = append(constants, c...)
		case "ValidatorParams":
			c, err := parseParamsToConstants(moduleName, "validator_", value)
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
