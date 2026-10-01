// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testVotingPeriod          = 7 * 24 * time.Hour
	testExpeditedVotingPeriod = 24 * time.Hour
	testMinDeposit            = "10000000000"
	testExpeditedMinDeposit   = "50000000000"
)

var govConstants = map[string]string{
	"min_deposit":             testMinDeposit,
	"expedited_min_deposit":   testExpeditedMinDeposit,
	"voting_period":           "604800000000000",
	"expedited_voting_period": "86400000000000",
	"max_deposit_period":      "604800000000000",
	"quorum":                  "0.334000000000000000",
	"threshold":               "0.500000000000000000",
	"expedited_threshold":     "0.667000000000000000",
	"veto_threshold":          "0.334000000000000000",
}

func expectGovConstants(repos *mock.TxRepos) {
	repos.Constants.EXPECT().
		Get(gomock.Any(), types.ModuleNameGov, gomock.Any()).
		DoAndReturn(func(_ context.Context, module types.ModuleName, name string) (storage.Constant, error) {
			value, ok := govConstants[name]
			if !ok {
				return storage.Constant{}, sql.ErrNoRows
			}
			return storage.Constant{Module: module, Name: name, Value: value}, nil
		}).
		AnyTimes()
}

// expectStoredProposals serves proposals already saved in the DB; the rest are not found
func expectStoredProposals(repos *mock.TxRepos, stored ...storage.Proposal) {
	byId := make(map[uint64]storage.Proposal)
	for i := range stored {
		byId[stored[i].Id] = stored[i]
	}
	repos.Proposals.EXPECT().
		GetByID(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, id uint64) (*storage.Proposal, error) {
			p, ok := byId[id]
			if !ok {
				return nil, sql.ErrNoRows
			}
			return &p, nil
		}).
		AnyTimes()
	repos.Proposals.EXPECT().
		IsNoRows(gomock.Any()).
		DoAndReturn(func(err error) bool { return errors.Is(err, sql.ErrNoRows) }).
		AnyTimes()
}

func proposalsMap(items ...*storage.Proposal) *sdkSync.Map[uint64, *storage.Proposal] {
	m := sdkSync.NewMap[uint64, *storage.Proposal]()
	for i := range items {
		m.Set(items[i].Id, items[i])
	}
	return m
}

func noVotes() *sdkSync.Map[decodeContext.VoteKey, []*storage.Vote] {
	return sdkSync.NewMap[decodeContext.VoteKey, []*storage.Vote]()
}

func TestFillExpeditedProposal(t *testing.T) {
	activation := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

	t.Run("regular proposal activated by deposit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)
		expectGovConstants(repos)
		expectStoredProposals(repos, storage.Proposal{
			Id:        1,
			Status:    types.ProposalStatusInactive,
			Expedited: testsuite.Ptr(false),
		})

		proposal := &storage.Proposal{Id: 1, Status: types.ProposalStatusActive, ActivationTime: &activation}
		err := fillExpeditedProposal(t.Context(), repos.Proposals, repos.Constants, proposalsMap(proposal))
		require.NoError(t, err)

		require.Equal(t, testMinDeposit, proposal.MinDeposit)
		// the regular end time is set later in saveProposals
		require.Nil(t, proposal.EndTime)
		require.Nil(t, proposal.Expedited)
	})

	t.Run("expedited proposal activated by deposit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)
		expectGovConstants(repos)
		expectStoredProposals(repos, storage.Proposal{
			Id:        2,
			Status:    types.ProposalStatusInactive,
			Expedited: testsuite.Ptr(true),
		})

		proposal := &storage.Proposal{Id: 2, Status: types.ProposalStatusActive, ActivationTime: &activation}
		err := fillExpeditedProposal(t.Context(), repos.Proposals, repos.Constants, proposalsMap(proposal))
		require.NoError(t, err)

		require.Equal(t, testExpeditedMinDeposit, proposal.MinDeposit)
		require.NotNil(t, proposal.EndTime)
		require.Equal(t, activation.Add(testExpeditedVotingPeriod), *proposal.EndTime)
	})

	// the initial deposit reaches expedited_min_deposit: the proposal is not in the DB yet
	t.Run("expedited proposal submitted and activated in one block", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)
		expectGovConstants(repos)
		expectStoredProposals(repos)

		proposal := &storage.Proposal{
			Id:             3,
			Status:         types.ProposalStatusActive,
			ActivationTime: &activation,
			CreatedAt:      activation,
			Expedited:      testsuite.Ptr(true),
		}
		err := fillExpeditedProposal(t.Context(), repos.Proposals, repos.Constants, proposalsMap(proposal))
		require.NoError(t, err)

		require.Equal(t, testExpeditedMinDeposit, proposal.MinDeposit)
		require.NotNil(t, proposal.EndTime)
		require.Equal(t, activation.Add(testExpeditedVotingPeriod), *proposal.EndTime)
		require.True(t, proposal.IsExpedited())
	})

	t.Run("regular proposal submitted and activated in one block", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)
		expectGovConstants(repos)
		expectStoredProposals(repos)

		proposal := &storage.Proposal{
			Id:             4,
			Status:         types.ProposalStatusActive,
			ActivationTime: &activation,
			CreatedAt:      activation,
			Expedited:      testsuite.Ptr(false),
		}
		err := fillExpeditedProposal(t.Context(), repos.Proposals, repos.Constants, proposalsMap(proposal))
		require.NoError(t, err)

		require.Equal(t, testMinDeposit, proposal.MinDeposit)
		require.Nil(t, proposal.EndTime)
	})

	// x/gov: VotingEndTime = VotingStartTime + voting_period, deposit untouched
	t.Run("expedited proposal rejected", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)
		expectGovConstants(repos)
		expeditedEnd := activation.Add(testExpeditedVotingPeriod)
		expectStoredProposals(repos, storage.Proposal{
			Id:             5,
			Status:         types.ProposalStatusActive,
			Expedited:      testsuite.Ptr(true),
			ActivationTime: &activation,
			EndTime:        &expeditedEnd,
			MinDeposit:     testExpeditedMinDeposit,
		})

		proposal := &storage.Proposal{Id: 5, Status: types.ProposalStatusActive, ExpeditedProposalRejected: true}
		err := fillExpeditedProposal(t.Context(), repos.Proposals, repos.Constants, proposalsMap(proposal))
		require.NoError(t, err)

		require.NotNil(t, proposal.EndTime)
		require.Equal(t, activation.Add(testVotingPeriod), *proposal.EndTime)
		require.NotNil(t, proposal.Expedited)
		require.False(t, *proposal.Expedited)
		// the proposal was activated as expedited, so its min deposit stays
		require.Empty(t, proposal.MinDeposit)
		require.Nil(t, proposal.ActivationTime)
	})

	t.Run("votes and finished proposals are skipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)

		vote := &storage.Proposal{Id: 6}
		finished := &storage.Proposal{Id: 7, Status: types.ProposalStatusApplied}
		err := fillExpeditedProposal(t.Context(), repos.Proposals, repos.Constants, proposalsMap(vote, finished))
		require.NoError(t, err)

		require.Nil(t, vote.EndTime)
		require.Empty(t, vote.MinDeposit)
		require.Nil(t, finished.EndTime)
		require.Empty(t, finished.MinDeposit)
	})
}

// finishing proposal from events carries only id and status; the rules come from the stored proposal
func TestSaveProposals_FinishedThreshold(t *testing.T) {
	activation := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		expedited     bool
		wantThreshold string
	}{
		{name: "expedited passed in expedited period", expedited: true, wantThreshold: "0.667000000000000000"},
		{name: "regular or downgraded", expedited: false, wantThreshold: "0.500000000000000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repos := mock.NewTxRepos(ctrl)
			tx := mock.NewMockTransaction(ctrl)
			expectGovConstants(repos)

			stored := storage.Proposal{
				Id:             1,
				Status:         types.ProposalStatusActive,
				Expedited:      testsuite.Ptr(tt.expedited),
				ActivationTime: &activation,
			}
			expectStoredProposals(repos, stored)
			repos.Proposals.EXPECT().Active(gomock.Any()).Return([]storage.Proposal{stored}, nil).AnyTimes()
			repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return(nil, nil).AnyTimes()
			repos.Votes.EXPECT().ListByProposal(gomock.Any(), uint64(1), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

			var saved []*storage.Proposal
			tx.EXPECT().
				SaveProposals(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, proposals ...*storage.Proposal) (int64, error) {
					saved = proposals
					return 0, nil
				}).
				Times(1)

			module := NewModule(nil, nil, config.Indexer{})
			finished := &storage.Proposal{Id: 1, Status: types.ProposalStatusApplied}
			_, err := module.saveProposals(t.Context(), tx, repos.Repos(), pkgTypes.Level(101), proposalsMap(finished), noVotes(), nil, newGovValidators())
			require.NoError(t, err)

			require.Len(t, saved, 1)
			require.Equal(t, tt.wantThreshold, saved[0].Threshold)
			require.Equal(t, "0.334000000000000000", saved[0].Quorum)
			require.Equal(t, "0.334000000000000000", saved[0].VetoQuorum)
		})
	}
}

// the downgrade must reach SaveProposals even when the block recounts voting power
func TestSaveProposals_ExpeditedRejectedIsSaved(t *testing.T) {
	activation := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	expeditedEnd := activation.Add(testExpeditedVotingPeriod)

	for _, height := range []pkgTypes.Level{601, 600} {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)
		tx := mock.NewMockTransaction(ctrl)
		expectGovConstants(repos)

		stored := storage.Proposal{
			Id:             5,
			Status:         types.ProposalStatusActive,
			Expedited:      testsuite.Ptr(true),
			ActivationTime: &activation,
			EndTime:        &expeditedEnd,
			MinDeposit:     testExpeditedMinDeposit,
		}
		expectStoredProposals(repos, stored)
		// Active does not select expedited and the counters
		repos.Proposals.EXPECT().Active(gomock.Any()).Return([]storage.Proposal{{
			Id:             5,
			Status:         types.ProposalStatusActive,
			ActivationTime: &activation,
			EndTime:        &expeditedEnd,
			MinDeposit:     testExpeditedMinDeposit,
		}}, nil).AnyTimes()
		repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return(nil, nil).AnyTimes()
		repos.Votes.EXPECT().ListByProposal(gomock.Any(), uint64(5), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

		var saved []*storage.Proposal
		tx.EXPECT().
			SaveProposals(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, proposals ...*storage.Proposal) (int64, error) {
				saved = proposals
				return 0, nil
			}).
			Times(1)

		module := NewModule(nil, nil, config.Indexer{})
		rejected := &storage.Proposal{Id: 5, Status: types.ProposalStatusActive, ExpeditedProposalRejected: true}
		_, err := module.saveProposals(t.Context(), tx, repos.Repos(), height, proposalsMap(rejected), noVotes(), nil, newGovValidators())
		require.NoError(t, err, "height %d", height)

		require.Len(t, saved, 1, "height %d", height)
		require.NotNil(t, saved[0].Expedited, "height %d", height)
		require.False(t, *saved[0].Expedited, "height %d", height)
		require.NotNil(t, saved[0].EndTime, "height %d", height)
		require.Equal(t, activation.Add(testVotingPeriod), *saved[0].EndTime, "height %d", height)
	}
}

// constants change when the proposal passes, not when it is submitted
func TestSaveProposalConstantUpdates(t *testing.T) {
	collect := func(saved *[]storage.Constant) func(context.Context, ...storage.Constant) error {
		return func(_ context.Context, constants ...storage.Constant) error {
			*saved = append(*saved, constants...)
			return nil
		}
	}

	t.Run("applied param change", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)
		tx := mock.NewMockTransaction(ctrl)

		repos.Proposals.EXPECT().
			GetByID(gomock.Any(), uint64(1)).
			Return(&storage.Proposal{
				Id:      1,
				Type:    types.ProposalTypeParamChanged,
				Changes: []byte(`[{"subspace":"gov","key":"expedited_threshold","value":"0.700000000000000000"},{"subspace":"staking","key":"max_validators","value":"110"}]`),
			}, nil).
			Times(1)
		repos.Proposals.EXPECT().
			GetByID(gomock.Any(), uint64(3)).
			Return(&storage.Proposal{Id: 3, Type: types.ProposalTypeText}, nil).
			Times(1)

		var saved []storage.Constant
		tx.EXPECT().SaveConstants(gomock.Any(), gomock.Any()).DoAndReturn(collect(&saved)).AnyTimes()

		proposals := proposalsMap(
			&storage.Proposal{Id: 1, Status: types.ProposalStatusApplied},
			// rejected and failed proposals change nothing
			&storage.Proposal{Id: 2, Status: types.ProposalStatusRejected},
			&storage.Proposal{Id: 3, Status: types.ProposalStatusApplied},
			&storage.Proposal{Id: 4, Status: types.ProposalStatusFailed},
		)
		require.NoError(t, saveProposalConstantUpdates(t.Context(), tx, repos.Proposals, proposals))

		got := make(map[string]string, len(saved))
		for _, c := range saved {
			got[c.Module.String()+"_"+c.Name] = c.Value
		}
		require.Equal(t, map[string]string{
			"gov_expedited_threshold": "0.700000000000000000",
			"staking_max_validators":  "110",
		}, got)
	})

	t.Run("nothing applied", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repos := mock.NewTxRepos(ctrl)
		tx := mock.NewMockTransaction(ctrl)

		var saved []storage.Constant
		tx.EXPECT().SaveConstants(gomock.Any(), gomock.Any()).DoAndReturn(collect(&saved)).AnyTimes()

		proposals := proposalsMap(&storage.Proposal{Id: 1, Status: types.ProposalStatusActive})
		require.NoError(t, saveProposalConstantUpdates(t.Context(), tx, repos.Proposals, proposals))
		require.Empty(t, saved)
	})
}

func TestSaveConstantUpdates(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)

	// empty: no write at all
	require.NoError(t, saveConstantUpdates(t.Context(), tx, sdkSync.NewMap[string, *storage.Constant]()))

	var saved []storage.Constant
	tx.EXPECT().
		SaveConstants(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, constants ...storage.Constant) error {
			saved = constants
			return nil
		}).
		Times(1)

	constants := sdkSync.NewMap[string, *storage.Constant]()
	constants.Set("gov_expedited_threshold", &storage.Constant{Module: types.ModuleNameGov, Name: "expedited_threshold", Value: "0.667000000000000000"})
	require.NoError(t, saveConstantUpdates(t.Context(), tx, constants))
	require.Len(t, saved, 1)
	require.Equal(t, "0.667000000000000000", saved[0].Value)
}
