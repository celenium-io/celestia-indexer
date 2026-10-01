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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type tallyResult struct {
	total, yes, no, veto, abstain string
}

func requireTally(t *testing.T, p *storage.Proposal, want tallyResult) {
	t.Helper()
	require.True(t, p.Tallied, "tallied")
	require.Equal(t, want.total, p.VotingPower.String(), "total")
	require.Equal(t, want.yes, p.YesVotingPower.String(), "yes")
	require.Equal(t, want.no, p.NoVotingPower.String(), "no")
	require.Equal(t, want.veto, p.NoWithVetoVotingPower.String(), "no with veto")
	require.Equal(t, want.abstain, p.AbstainVotingPower.String(), "abstain")
}

func weight(t *testing.T, s string) types.Numeric {
	t.Helper()
	w, err := types.NumericFromString(s)
	require.NoError(t, err)
	return w
}

func TestTallyProposal(t *testing.T) {
	validatorsPower := map[uint64]types.Numeric{
		1: types.NumericFromInt64(1000),
		2: types.NumericFromInt64(500),
	}

	t.Run("weighted delegator vote counts power once", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)

		repos.Votes.EXPECT().
			ListByProposal(gomock.Any(), uint64(1), votesPageSize, 0).
			Return([]storage.Vote{
				{VoterId: 10, Option: types.VoteOptionYes, Weight: weight(t, "0.7")},
				{VoterId: 10, Option: types.VoteOptionNo, Weight: weight(t, "0.3")},
			}, nil)
		repos.Delegation.EXPECT().
			AddressDelegations(gomock.Any(), uint64(10)).
			Return([]storage.Delegation{
				{AddressId: 10, ValidatorId: 1, Amount: types.NumericFromInt64(100)},
			}, nil)

		proposal := &storage.Proposal{Id: 1}
		require.NoError(t, tallyProposal(t.Context(), repos.Repos(), proposal, validatorsPower))
		requireTally(t, proposal, tallyResult{total: "100", yes: "70", no: "30", veto: "0", abstain: "0"})
	})

	t.Run("validator vote gets stake minus voters' delegations", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)

		repos.Votes.EXPECT().
			ListByProposal(gomock.Any(), uint64(1), votesPageSize, 0).
			Return([]storage.Vote{
				{VoterId: 10, Option: types.VoteOptionYes, Weight: weight(t, "0.7")},
				{VoterId: 10, Option: types.VoteOptionNo, Weight: weight(t, "0.3")},
				{VoterId: 20, ValidatorId: testsuite.Ptr(uint64(1)), Option: types.VoteOptionAbstain, Weight: weight(t, "1")},
			}, nil)
		repos.Delegation.EXPECT().
			AddressDelegations(gomock.Any(), uint64(10), uint64(20)).
			Return([]storage.Delegation{
				{AddressId: 10, ValidatorId: 1, Amount: types.NumericFromInt64(100)},
			}, nil)

		proposal := &storage.Proposal{Id: 1}
		require.NoError(t, tallyProposal(t.Context(), repos.Repos(), proposal, validatorsPower))
		requireTally(t, proposal, tallyResult{total: "1000", yes: "70", no: "30", veto: "0", abstain: "900"})
	})

	t.Run("weighted validator vote with self-delegation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)

		repos.Votes.EXPECT().
			ListByProposal(gomock.Any(), uint64(1), votesPageSize, 0).
			Return([]storage.Vote{
				{VoterId: 20, ValidatorId: testsuite.Ptr(uint64(1)), Option: types.VoteOptionYes, Weight: weight(t, "0.5")},
				{VoterId: 20, ValidatorId: testsuite.Ptr(uint64(1)), Option: types.VoteOptionNo, Weight: weight(t, "0.5")},
			}, nil)
		repos.Delegation.EXPECT().
			AddressDelegations(gomock.Any(), uint64(20)).
			Return([]storage.Delegation{
				{AddressId: 20, ValidatorId: 1, Amount: types.NumericFromInt64(200)},
			}, nil)

		proposal := &storage.Proposal{Id: 1}
		require.NoError(t, tallyProposal(t.Context(), repos.Repos(), proposal, validatorsPower))
		requireTally(t, proposal, tallyResult{total: "1000", yes: "500", no: "500", veto: "0", abstain: "0"})
	})

	t.Run("delegation to not bonded validator is ignored", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)

		repos.Votes.EXPECT().
			ListByProposal(gomock.Any(), uint64(1), votesPageSize, 0).
			Return([]storage.Vote{
				{VoterId: 10, Option: types.VoteOptionNoWithVeto, Weight: weight(t, "1")},
			}, nil)
		repos.Delegation.EXPECT().
			AddressDelegations(gomock.Any(), uint64(10)).
			Return([]storage.Delegation{
				{AddressId: 10, ValidatorId: 2, Amount: types.NumericFromInt64(50)},
				{AddressId: 10, ValidatorId: 3, Amount: types.NumericFromInt64(1000)},
			}, nil)

		proposal := &storage.Proposal{Id: 1}
		require.NoError(t, tallyProposal(t.Context(), repos.Repos(), proposal, validatorsPower))
		requireTally(t, proposal, tallyResult{total: "50", yes: "0", no: "0", veto: "50", abstain: "0"})
	})

	t.Run("results are truncated only at the end", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)

		third := weight(t, "0.333333333333333333")
		repos.Votes.EXPECT().
			ListByProposal(gomock.Any(), uint64(1), votesPageSize, 0).
			Return([]storage.Vote{
				{VoterId: 10, Option: types.VoteOptionYes, Weight: third},
				{VoterId: 10, Option: types.VoteOptionNo, Weight: third},
				{VoterId: 10, Option: types.VoteOptionAbstain, Weight: weight(t, "0.333333333333333334")},
				{VoterId: 11, Option: types.VoteOptionYes, Weight: third},
				{VoterId: 11, Option: types.VoteOptionNo, Weight: third},
				{VoterId: 11, Option: types.VoteOptionAbstain, Weight: weight(t, "0.333333333333333334")},
			}, nil)
		repos.Delegation.EXPECT().
			AddressDelegations(gomock.Any(), uint64(10), uint64(11)).
			Return([]storage.Delegation{
				{AddressId: 10, ValidatorId: 1, Amount: types.NumericFromInt64(100)},
				{AddressId: 11, ValidatorId: 1, Amount: types.NumericFromInt64(50)},
			}, nil)

		proposal := &storage.Proposal{Id: 1}
		require.NoError(t, tallyProposal(t.Context(), repos.Repos(), proposal, validatorsPower))
		// 33.33.. + 16.66.. = 49.99.. per option: truncating each voter's share first would give 49
		requireTally(t, proposal, tallyResult{total: "150", yes: "49", no: "49", veto: "0", abstain: "50"})
	})

	t.Run("voters are batched and a voter is not split between batches", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)

		votes := make([]storage.Vote, 0, votersBatchSize*2+2)
		for i := range votersBatchSize + 1 {
			id := uint64(i + 1)
			votes = append(votes,
				storage.Vote{VoterId: id, Option: types.VoteOptionYes, Weight: weight(t, "0.5")},
				storage.Vote{VoterId: id, Option: types.VoteOptionNo, Weight: weight(t, "0.5")},
			)
		}
		repos.Votes.EXPECT().
			ListByProposal(gomock.Any(), uint64(1), votesPageSize, 0).
			Return(votes, nil)

		delegationsOf := func(ids []uint64) []storage.Delegation {
			result := make([]storage.Delegation, len(ids))
			for i, id := range ids {
				result[i] = storage.Delegation{AddressId: id, ValidatorId: 1, Amount: types.NumericFromInt64(2)}
			}
			return result
		}
		var batches [][]uint64
		repos.Delegation.EXPECT().
			AddressDelegations(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, ids ...uint64) ([]storage.Delegation, error) {
				batches = append(batches, ids)
				return delegationsOf(ids), nil
			}).
			Times(2)

		proposal := &storage.Proposal{Id: 1}
		require.NoError(t, tallyProposal(t.Context(), repos.Repos(), proposal, validatorsPower))

		require.Len(t, batches, 2)
		require.Len(t, batches[0], votersBatchSize)
		require.Equal(t, []uint64{votersBatchSize + 1}, batches[1])
		requireTally(t, proposal, tallyResult{total: "202", yes: "101", no: "101", veto: "0", abstain: "0"})
	})

	t.Run("no votes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)

		repos.Votes.EXPECT().
			ListByProposal(gomock.Any(), uint64(1), votesPageSize, 0).
			Return([]storage.Vote{}, nil)

		proposal := &storage.Proposal{Id: 1}
		require.NoError(t, tallyProposal(t.Context(), repos.Repos(), proposal, validatorsPower))
		requireTally(t, proposal, tallyResult{total: "0", yes: "0", no: "0", veto: "0", abstain: "0"})
	})
}
