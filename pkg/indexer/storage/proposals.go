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
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
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
	votesMap *sdkSync.Map[decodeContext.VoteKey, []*storage.Vote],
	addrToId map[string]uint64,
	govSet *govValidators,
) (int64, error) {
	if votesMap.Len() > 0 {
		data := make([]*storage.Vote, 0, votesMap.Len())
		for votes := range votesMap.AllValues() {
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
				data = append(data, votes[i])
			}
		}

		votesCount, err := tx.SaveVotes(ctx, data...)
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

	filled, err := module.fillProposalsVotingPower(ctx, repos, height, proposals, govSet)
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

type weightedOption struct {
	option types.VoteOption
	weight types.Numeric
}

type voterOptions struct {
	voterId     uint64
	validatorId *uint64
	options     []weightedOption
}

// proposalTally mirrors x/gov Tally: each voter's power counts once toward the total,
// options get power*weight, and voters' delegations are deducted from their validators.
type proposalTally struct {
	results         map[types.VoteOption]types.Numeric
	total           types.Numeric
	validatorsPower map[uint64]types.Numeric
	validatorMinus  map[uint64]types.Numeric
	votedValidators map[uint64][]weightedOption
}

func newProposalTally(validatorsPower map[uint64]types.Numeric) *proposalTally {
	return &proposalTally{
		results:         make(map[types.VoteOption]types.Numeric),
		total:           types.NumericZero(),
		validatorsPower: validatorsPower,
		validatorMinus:  make(map[uint64]types.Numeric),
		votedValidators: make(map[uint64][]weightedOption),
	}
}

func (t *proposalTally) add(power types.Numeric, options []weightedOption) {
	t.total = t.total.Add(power)
	for _, o := range options {
		t.results[o.option] = t.result(o.option).Add(power.Mul(o.weight))
	}
}

func (t *proposalTally) result(option types.VoteOption) types.Numeric {
	if value, ok := t.results[option]; ok {
		return value
	}
	return types.NumericZero()
}

func (t *proposalTally) addVoters(ctx context.Context, delegations storage.IDelegation, voters []voterOptions) error {
	if len(voters) == 0 {
		return nil
	}

	ids := make([]uint64, len(voters))
	for i := range voters {
		ids[i] = voters[i].voterId
	}
	list, err := delegations.AddressDelegations(ctx, ids...)
	if err != nil {
		return errors.Wrap(err, "can't receive address delegations")
	}
	byAddress := make(map[uint64][]storage.Delegation, len(voters))
	for i := range list {
		byAddress[list[i].AddressId] = append(byAddress[list[i].AddressId], list[i])
	}

	for _, voter := range voters {
		if voter.validatorId != nil {
			t.votedValidators[*voter.validatorId] = voter.options
		}
		// self-delegation is not skipped: it is deducted from the validator below, as in x/gov
		for _, d := range byAddress[voter.voterId] {
			if _, bonded := t.validatorsPower[d.ValidatorId]; !bonded {
				continue
			}
			if minus, ok := t.validatorMinus[d.ValidatorId]; ok {
				t.validatorMinus[d.ValidatorId] = minus.Add(d.Amount)
			} else {
				t.validatorMinus[d.ValidatorId] = d.Amount
			}
			t.add(d.Amount, voter.options)
		}
	}
	return nil
}

func (t *proposalTally) addValidators() {
	for id, options := range t.votedValidators {
		power, bonded := t.validatorsPower[id]
		if !bonded {
			continue
		}
		if minus, ok := t.validatorMinus[id]; ok {
			power = power.Sub(minus)
		}
		t.add(power, options)
	}
}

// fill truncates the results like x/gov does after summing decimal powers.
func (t *proposalTally) fill(proposal *storage.Proposal) {
	proposal.VotingPower = t.total.Floor()
	proposal.YesVotingPower = t.result(types.VoteOptionYes).Floor()
	proposal.NoVotingPower = t.result(types.VoteOptionNo).Floor()
	proposal.NoWithVetoVotingPower = t.result(types.VoteOptionNoWithVeto).Floor()
	proposal.AbstainVotingPower = t.result(types.VoteOptionAbstain).Floor()
	proposal.Tallied = true
}

// proposalsTallyInterval is how often active proposals are recounted; finished ones are counted at once.
const proposalsTallyInterval = 600

func needsProposalsTally(height pkgTypes.Level, proposals *sdkSync.Map[uint64, *storage.Proposal]) bool {
	if height%proposalsTallyInterval == 0 {
		return true
	}
	for p := range proposals.AllValues() {
		if p.Finished() {
			return true
		}
	}
	return false
}

// govValidators is the validator set x/gov sees: gov's EndBlocker runs before staking's,
// so it is the set bonded before the block's validator_updates.
type govValidators struct {
	bonded map[uint64]struct{}
	// jailed in the block's BeginBlock: out of the tally, but their tokens are still bonded
	jailed map[uint64]struct{}
}

// takeGovValidators must run before the block's validator writes; nil means no tally in the block.
func takeGovValidators(
	ctx context.Context,
	start *blockStartValidators,
	height pkgTypes.Level,
	proposals *sdkSync.Map[uint64, *storage.Proposal],
) (*govValidators, error) {
	if !needsProposalsTally(height, proposals) {
		return nil, nil
	}
	bonded, err := start.get(ctx)
	if err != nil {
		return nil, err
	}
	set := &govValidators{
		bonded: make(map[uint64]struct{}, len(bonded)),
		jailed: make(map[uint64]struct{}),
	}
	for i := range bonded {
		set.bonded[bonded[i].Id] = struct{}{}
	}
	return set, nil
}

// markJailed needs validator ids resolved by saveValidators.
func (set *govValidators) markJailed(jails *sdkSync.Map[string, *storage.Jail]) {
	if set == nil {
		return
	}
	for jail := range jails.AllValues() {
		if _, ok := set.bonded[jail.ValidatorId]; ok {
			set.jailed[jail.ValidatorId] = struct{}{}
		}
	}
}

type govPower struct {
	// tallied validators' stake
	validators map[uint64]types.Numeric
	// TotalBondedTokens: the quorum denominator
	total types.Numeric
}

// power reads the stake after the block's transactions, as gov sees it.
func (set *govValidators) power(ctx context.Context, validators storage.IValidator) (govPower, error) {
	bonded, err := validators.BondedValidators(ctx)
	if err != nil {
		return govPower{}, errors.Wrap(err, "get bonded validators")
	}
	stakes := make(map[uint64]types.Numeric, len(set.bonded))
	for i := range bonded {
		// validators bonded by this block's validator_updates are not in the set
		if _, ok := set.bonded[bonded[i].Id]; ok {
			stakes[bonded[i].Id] = bonded[i].Stake
		}
	}
	// left the active set in this block: gov still counts them
	for id := range set.bonded {
		if _, ok := stakes[id]; ok {
			continue
		}
		validator, err := validators.GetByID(ctx, id)
		if err != nil {
			return govPower{}, errors.Wrapf(err, "get validator %d", id)
		}
		stakes[id] = validator.Stake
	}

	result := govPower{
		validators: make(map[uint64]types.Numeric, len(stakes)),
		total:      types.NumericZero(),
	}
	for id, stake := range stakes {
		result.total = result.total.Add(stake)
		if _, jailed := set.jailed[id]; !jailed {
			result.validators[id] = stake
		}
	}
	return result, nil
}

func (module *Module) fillProposalsVotingPower(
	ctx context.Context,
	repos storage.TxRepos,
	height pkgTypes.Level,
	proposals *sdkSync.Map[uint64, *storage.Proposal],
	govSet *govValidators,
) ([]*storage.Proposal, error) {
	// 1. Receive all active or just completed proposals

	// 1.1 Return if we don't have proposal updates and it's not certain block height (one block in hour)

	if !needsProposalsTally(height, proposals) {
		return proposals.Values(), nil
	}

	finished := make(map[uint64]*storage.Proposal)
	for value := range proposals.AllValues() {
		if value.Finished() {
			finished[value.Id] = value
		}
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

	if govSet == nil {
		return nil, errors.New("gov validator set was not taken at the block start")
	}
	power, err := govSet.power(ctx, repos.Validators)
	if err != nil {
		return nil, err
	}

	// 3. Compute voting results

	for _, proposal := range finished {
		if proposal.Finished() {
			proposal.TotalVotingPower = power.total

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

		if err := tallyProposal(ctx, repos, proposal, power.validators); err != nil {
			return nil, err
		}
	}

	return proposals.Values(), nil
}

const (
	votesPageSize   = 1000
	votersBatchSize = 100
)

type votersBatch []voterOptions

func (b votersBatch) Len() int {
	return len(b)
}

func (b votersBatch) Last() *voterOptions {
	if len(b) == 0 {
		return nil
	}
	return &b[len(b)-1]
}

// tallyProposal relies on votes ordered by voter, so options of one voter come in a row.
func tallyProposal(
	ctx context.Context,
	repos storage.TxRepos,
	proposal *storage.Proposal,
	validatorsPower map[uint64]types.Numeric,
) error {
	var (
		tally = newProposalTally(validatorsPower)
		batch = make(votersBatch, 0, votersBatchSize)
	)

	paginate := sdkSync.Paginate(
		ctx, votesPageSize,
		func(ctx context.Context, limit, offset int) ([]storage.Vote, error) {
			return repos.Votes.ListByProposal(ctx, proposal.Id, limit, offset)
		},
	)

	for vote, err := range paginate {
		if err != nil {
			return errors.Wrapf(err, "get proposal votes: proposal_id=%d", proposal.Id)
		}

		if n := batch.Len(); n == 0 || batch.Last().voterId != vote.VoterId {
			// voter changed, so every voter in the batch is complete
			if n >= votersBatchSize {
				if err := tally.addVoters(ctx, repos.Delegation, batch); err != nil {
					return errors.Wrapf(err, "proposal_id=%d", proposal.Id)
				}
				batch = batch[:0]
			}
			batch = append(batch, voterOptions{
				voterId:     vote.VoterId,
				validatorId: vote.ValidatorId,
			})
		}

		if last := batch.Last(); last != nil {
			last.options = append(last.options, weightedOption{option: vote.Option, weight: vote.Weight})
		}
	}

	if err := tally.addVoters(ctx, repos.Delegation, batch); err != nil {
		return errors.Wrapf(err, "proposal_id=%d", proposal.Id)
	}
	tally.addValidators()
	tally.fill(proposal)
	return nil
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
