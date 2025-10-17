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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFillProposalVotingPower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constants := mock.NewMockIConstant(ctrl)
	validators := mock.NewMockIValidator(ctrl)

	module := NewModule(nil, constants, validators, nil, config.Indexer{})

	t.Run("not fill", func(t *testing.T) {
		tx := mock.NewMockTransaction(ctrl)

		filled, err := module.fillProposalsVotingPower(t.Context(), tx, 1, []*storage.Proposal{{
			Status: types.ProposalStatusActive,
		}})
		require.NoError(t, err)
		require.Len(t, filled, 1)
	})

	t.Run("no active and finished", func(t *testing.T) {
		tx := mock.NewMockTransaction(ctrl)

		tx.EXPECT().
			ActiveProposals(t.Context()).
			Return([]storage.Proposal{}, nil).
			Times(1)

		filled, err := module.fillProposalsVotingPower(t.Context(), tx, 600, []*storage.Proposal{{
			Id:         1,
			Status:     types.ProposalStatusActive,
			Type:       types.ProposalTypeParamChanged,
			Abstain:    10,
			Yes:        100,
			No:         1,
			NoWithVeto: 1,
		}})
		require.NoError(t, err)
		require.Len(t, filled, 1)
	})

	t.Run("active and no finished", func(t *testing.T) {
		tx := mock.NewMockTransaction(ctrl)

		constants.EXPECT().
			Get(gomock.Any(), types.ModuleNameStaking, "max_validators").
			Return(storage.Constant{
				Name:   "max_validators",
				Module: types.ModuleNameStaking,
				Value:  "100",
			}, nil).
			Times(1)

		tx.EXPECT().
			ActiveProposals(t.Context()).
			Return([]storage.Proposal{{
				Id:         1,
				Status:     types.ProposalStatusActive,
				Type:       types.ProposalTypeParamChanged,
				Abstain:    10,
				Yes:        100,
				No:         1,
				NoWithVeto: 1,
			}}, nil).
			Times(1)

		validators.EXPECT().
			TotalVotingPower(gomock.Any(), 100).
			Return(decimal.RequireFromString("10000"), nil).
			Times(1)

		tx.EXPECT().
			BondedValidators(t.Context(), 100).
			Return([]storage.Validator{{
				Id:    1,
				Stake: decimal.RequireFromString("100000000"),
			}, {
				Id:    2,
				Stake: decimal.RequireFromString("200000000"),
			}}, nil).
			Times(1)

		tx.EXPECT().
			ProposalVotes(t.Context(), uint64(1), 1000, 0).
			Return([]storage.Vote{{
				ValidatorId: testsuite.Ptr(uint64(1)),
				VoterId:     3,
				Option:      types.VoteOptionAbstain,
			}, {
				VoterId: 1,
				Option:  types.VoteOptionNo,
			}, {
				VoterId: 2,
				Option:  types.VoteOptionYes,
			}}, nil).
			Times(1)

		tx.EXPECT().
			AddressDelegations(t.Context(), uint64(1)).
			Return([]storage.Delegation{{
				ValidatorId: 1,
				AddressId:   1,
				Amount:      decimal.RequireFromString("50000000"),
			}}, nil).
			Times(1)

		tx.EXPECT().
			AddressDelegations(t.Context(), uint64(2)).
			Return([]storage.Delegation{{
				ValidatorId: 1,
				AddressId:   1,
				Amount:      decimal.RequireFromString("10000000"),
			}}, nil).
			Times(1)

		tx.EXPECT().
			AddressDelegations(t.Context(), uint64(3)).
			Return([]storage.Delegation{{
				ValidatorId: 1,
				AddressId:   1,
				Amount:      decimal.RequireFromString("10000000"),
			}}, nil).
			Times(1)

		filled, err := module.fillProposalsVotingPower(t.Context(), tx, 600, []*storage.Proposal{{
			Status: types.ProposalStatusActive,
		}})
		require.NoError(t, err)
		require.Len(t, filled, 1)

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

		module := NewModule(nil, constants, validators, nil, config.Indexer{})
		got, err := module.getConstantDuration(ctx, types.ModuleNameGov, "voting_period")
		require.NoError(t, err)
		require.EqualValues(t, "24h0m0s", got.String())
	})
}
