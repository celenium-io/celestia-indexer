// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFillProposalVotingPower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validators := mock.NewMockIValidator(ctrl)
	repos := storage.TxRepos{Validators: validators}

	module := NewModule(nil, nil, config.Indexer{})

	t.Run("not fill", func(t *testing.T) {

		proposals := sdkSync.NewMap[uint64, *storage.Proposal]()
		proposals.Set(1, &storage.Proposal{
			Id:     1,
			Status: types.ProposalStatusActive,
		})

		filled, err := module.fillProposalsVotingPower(t.Context(), repos, 1, proposals, nil)
		require.NoError(t, err)
		require.Len(t, filled, 1)
		require.False(t, filled[0].Tallied, "no tally outside the recount height")
	})

	t.Run("no active and finished", func(t *testing.T) {
		repos := mock.NewTxRepos(ctrl)
		repos.Proposals.
			EXPECT().
			Active(t.Context()).
			Return([]storage.Proposal{}, nil).
			Times(1)

		proposals := sdkSync.NewMap[uint64, *storage.Proposal]()
		proposals.Set(1, &storage.Proposal{
			Id:         1,
			Status:     types.ProposalStatusActive,
			Type:       types.ProposalTypeParamChanged,
			Abstain:    10,
			Yes:        100,
			No:         1,
			NoWithVeto: 1,
		})

		filled, err := module.fillProposalsVotingPower(t.Context(), repos.Repos(), 600, proposals, newGovValidators())
		require.NoError(t, err)
		require.Len(t, filled, 1)
	})

	t.Run("active and no finished", func(t *testing.T) {
		repos := mock.NewTxRepos(ctrl)

		repos.Proposals.
			EXPECT().
			Active(t.Context()).
			Return([]storage.Proposal{{
				Id:     1,
				Status: types.ProposalStatusActive,
				Type:   types.ProposalTypeParamChanged,
			}}, nil).
			Times(1)
		repos.Validators.
			EXPECT().
			BondedValidators(t.Context()).
			Return([]storage.Validator{{
				Id:    1,
				Stake: types.NumericFromInt64(100000000),
			}, {
				Id:    2,
				Stake: types.NumericFromInt64(200000000),
			}}, nil).
			Times(1)

		repos.Votes.
			EXPECT().
			ListByProposal(t.Context(), uint64(1), 1000, 0).
			Return([]storage.Vote{{
				VoterId: 1,
				Option:  types.VoteOptionNo,
				Weight:  types.NumericFromInt64(1),
			}, {
				VoterId: 2,
				Option:  types.VoteOptionYes,
				Weight:  types.NumericFromInt64(1),
			}, {
				ValidatorId: testsuite.Ptr(uint64(1)),
				VoterId:     3,
				Option:      types.VoteOptionAbstain,
				Weight:      types.NumericFromInt64(1),
			}}, nil).
			Times(1)

		repos.Delegation.
			EXPECT().
			AddressDelegations(t.Context(), uint64(1), uint64(2), uint64(3)).
			Return([]storage.Delegation{{
				ValidatorId: 1,
				AddressId:   1,
				Amount:      types.NumericFromInt64(50000000),
			}, {
				ValidatorId: 1,
				AddressId:   2,
				Amount:      types.NumericFromInt64(10000000),
			}, {
				ValidatorId: 1,
				AddressId:   3,
				Amount:      types.NumericFromInt64(10000000),
			}}, nil).
			Times(1)

		proposals := sdkSync.NewMap[uint64, *storage.Proposal]()
		proposals.Set(1, &storage.Proposal{
			Id:         1,
			Status:     types.ProposalStatusActive,
			Type:       types.ProposalTypeParamChanged,
			Abstain:    10,
			Yes:        100,
			No:         1,
			NoWithVeto: 1,
		})

		filled, err := module.fillProposalsVotingPower(t.Context(), repos.Repos(), 600, proposals, newGovValidators(1, 2))
		require.NoError(t, err)
		require.Len(t, filled, 1)
		require.True(t, filled[0].Tallied)

		require.Equal(t, "100000000", filled[0].VotingPower.String())
		require.Equal(t, "40000000", filled[0].AbstainVotingPower.String())
		require.Equal(t, "50000000", filled[0].NoVotingPower.String())
		require.Equal(t, "10000000", filled[0].YesVotingPower.String())
	})
}

func TestModule_getConstantDuration(t *testing.T) {
	t.Run("get constant", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		constants := mock.NewMockIConstant(ctrl)
		validators := mock.NewMockIValidator(ctrl)

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		constants.EXPECT().
			Get(gomock.Any(), types.ModuleNameGov, "voting_period").
			Return(storage.Constant{
				Module: types.ModuleNameGov,
				Name:   "voting_period",
				Value:  "86400000000000",
			}, nil).
			Times(1)

		module := NewModule(nil, nil, config.Indexer{})
		repos := storage.TxRepos{Constants: constants, Validators: validators}
		got, err := module.getConstantDuration(ctx, repos, types.ModuleNameGov, "voting_period")
		require.NoError(t, err)
		require.EqualValues(t, "24h0m0s", got.String())
	})
}
