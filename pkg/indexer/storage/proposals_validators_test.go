// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newGovValidators(ids ...uint64) *govValidators {
	set := &govValidators{
		bonded: make(map[uint64]struct{}, len(ids)),
		jailed: make(map[uint64]struct{}),
	}
	for _, id := range ids {
		set.bonded[id] = struct{}{}
	}
	return set
}

func TestTakeGovValidators(t *testing.T) {
	bonded := []storage.Validator{{Id: 1}, {Id: 2}}

	t.Run("no tally in the block", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		validators := mock.NewMockIValidator(ctrl)

		set, err := takeGovValidators(t.Context(), newBlockStartValidators(validators), 601, proposalsMap(&storage.Proposal{
			Id: 1, Status: types.ProposalStatusActive,
		}))
		require.NoError(t, err)
		require.Nil(t, set)
	})

	t.Run("proposal finished in the block", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		validators := mock.NewMockIValidator(ctrl)
		validators.EXPECT().BondedValidators(gomock.Any()).Return(bonded, nil)

		set, err := takeGovValidators(t.Context(), newBlockStartValidators(validators), 601, proposalsMap(&storage.Proposal{
			Id: 1, Status: types.ProposalStatusRejected,
		}))
		require.NoError(t, err)
		require.Equal(t, newGovValidators(1, 2), set)
	})

	t.Run("recount height", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		validators := mock.NewMockIValidator(ctrl)
		validators.EXPECT().BondedValidators(gomock.Any()).Return(bonded, nil)

		set, err := takeGovValidators(t.Context(), newBlockStartValidators(validators), pkgTypes.Level(proposalsTallyInterval*3), proposalsMap())
		require.NoError(t, err)
		require.Equal(t, newGovValidators(1, 2), set)
	})
}

func TestBlockStartValidators_SharedBySignalAndGov(t *testing.T) {
	ctrl := gomock.NewController(t)
	repos := mock.NewTxRepos(ctrl)

	power := types.NumericFromInt64(10)
	repos.Upgrades.EXPECT().PendingVersions(gomock.Any(), uint64(3)).Return([]uint64{4}, nil)
	repos.Validators.EXPECT().BondedValidators(gomock.Any()).
		Return([]storage.Validator{{Id: 1, Power: &power}}, nil).
		Times(1)

	dCtx := context.NewContext()
	dCtx.Block = &storage.Block{Height: pkgTypes.Level(proposalsTallyInterval), VersionApp: 3}

	start := newBlockStartValidators(repos.Validators)
	round, err := prepareSignalRound(t.Context(), repos.Repos(), dCtx, 3, start)
	require.NoError(t, err)
	require.NotNil(t, round.snapshot)

	set, err := takeGovValidators(t.Context(), start, dCtx.Block.Height, dCtx.Proposals)
	require.NoError(t, err)
	require.Equal(t, newGovValidators(1), set)
}

func TestGovValidators_MarkJailed(t *testing.T) {
	jails := sdkSync.NewMap[string, *storage.Jail]()
	jails.Set("cons1", &storage.Jail{ValidatorId: 2})
	jails.Set("cons2", &storage.Jail{ValidatorId: 7}) // not bonded before the block

	set := newGovValidators(1, 2, 3)
	set.markJailed(jails)
	require.Equal(t, map[uint64]struct{}{2: {}}, set.jailed)
	require.Len(t, set.bonded, 3)

	// no tally in the block
	var empty *govValidators
	require.NotPanics(t, func() { empty.markJailed(jails) })
}

func TestGovValidators_Power(t *testing.T) {
	ctrl := gomock.NewController(t)
	validators := mock.NewMockIValidator(ctrl)

	// after the block's validator_updates: 4 joined the active set, 3 left it, 2 was jailed
	validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{
		{Id: 1, Stake: types.NumericFromInt64(100)},
		{Id: 4, Stake: types.NumericFromInt64(400)},
	}, nil)
	validators.EXPECT().GetByID(gomock.Any(), uint64(2)).Return(&storage.Validator{
		Id: 2, Stake: types.NumericFromInt64(20),
	}, nil)
	validators.EXPECT().GetByID(gomock.Any(), uint64(3)).Return(&storage.Validator{
		Id: 3, Stake: types.NumericFromInt64(30),
	}, nil)

	set := newGovValidators(1, 2, 3)
	set.jailed[2] = struct{}{}

	power, err := set.power(t.Context(), validators)
	require.NoError(t, err)

	got := make(map[uint64]string, len(power.validators))
	for id, stake := range power.validators {
		got[id] = stake.String()
	}
	// jailed validators are out of the tally, but their tokens are still in the bonded pool
	require.Equal(t, map[uint64]string{1: "100", 3: "30"}, got)
	require.Equal(t, "150", power.total.String())
}

func TestFillProposalVotingPower_TotalVotingPower(t *testing.T) {
	ctrl := gomock.NewController(t)
	repos := mock.NewTxRepos(ctrl)

	repos.Proposals.EXPECT().Active(gomock.Any()).Return(nil, nil)
	repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{
		{Id: 1, Stake: types.NumericFromInt64(1_000_000)},
		{Id: 2, Stake: types.NumericFromInt64(3_000_000)},
	}, nil)
	for _, name := range []string{"quorum", "threshold", "veto_threshold"} {
		repos.Constants.EXPECT().Get(gomock.Any(), types.ModuleNameGov, name).Return(storage.Constant{Value: "0.5"}, nil)
	}
	repos.Votes.EXPECT().ListByProposal(gomock.Any(), uint64(1), votesPageSize, 0).Return(nil, nil)

	module := NewModule(nil, nil, config.Indexer{})
	filled, err := module.fillProposalsVotingPower(t.Context(), repos.Repos(), 601, proposalsMap(&storage.Proposal{
		Id: 1, Status: types.ProposalStatusRejected,
	}), newGovValidators(1, 2))
	require.NoError(t, err)
	require.Len(t, filled, 1)
	// utia, the same unit as the voting power
	require.Equal(t, "4000000", filled[0].TotalVotingPower.String())
}

func TestFillProposalVotingPower_RequiresGovValidators(t *testing.T) {
	ctrl := gomock.NewController(t)
	repos := mock.NewTxRepos(ctrl)
	repos.Proposals.EXPECT().Active(gomock.Any()).Return(nil, nil)

	module := NewModule(nil, nil, config.Indexer{})
	_, err := module.fillProposalsVotingPower(t.Context(), repos.Repos(), 601, proposalsMap(&storage.Proposal{
		Id: 1, Status: types.ProposalStatusApplied,
	}), nil)
	require.ErrorContains(t, err, "gov validator set was not taken")
}

// a canceled proposal keeps the tally at the cancel block, like the other final statuses
func TestFillProposalVotingPower_Cancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	repos := mock.NewTxRepos(ctrl)

	repos.Proposals.EXPECT().Active(gomock.Any()).Return(nil, nil)
	repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{
		{Id: 1, Stake: types.NumericFromInt64(1000)},
	}, nil)
	for _, name := range []string{"quorum", "threshold", "veto_threshold"} {
		repos.Constants.EXPECT().Get(gomock.Any(), types.ModuleNameGov, name).Return(storage.Constant{Value: "0.5"}, nil)
	}
	repos.Votes.EXPECT().ListByProposal(gomock.Any(), uint64(1), votesPageSize, 0).Return([]storage.Vote{
		{VoterId: 10, ValidatorId: testsuite.Ptr(uint64(1)), Option: types.VoteOptionYes, Weight: types.NumericFromInt64(1)},
	}, nil)
	repos.Delegation.EXPECT().AddressDelegations(gomock.Any(), uint64(10)).Return(nil, nil)

	module := NewModule(nil, nil, config.Indexer{})
	filled, err := module.fillProposalsVotingPower(t.Context(), repos.Repos(), 601, proposalsMap(&storage.Proposal{
		Id: 1, Status: types.ProposalStatusCancelled,
	}), newGovValidators(1))
	require.NoError(t, err)
	require.Len(t, filled, 1)

	p := filled[0]
	require.True(t, p.Tallied)
	require.Equal(t, "1000", p.YesVotingPower.String())
	require.Equal(t, "1000", p.TotalVotingPower.String())
	require.Equal(t, "0.5", p.Quorum)
	require.Equal(t, "0.5", p.Threshold)
	require.Equal(t, "0.5", p.VetoQuorum)
}
