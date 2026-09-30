// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRollbackProposals(t *testing.T) {
	activation := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	regularEnd := activation.Add(7 * 24 * time.Hour)

	t.Run("expedited proposal rejected", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		tx := mock.NewMockTransaction(ctrl)
		constants := mock.NewMockIConstant(ctrl)
		proposals := mock.NewMockIProposal(ctrl)

		// constants keep durations in nanoseconds
		constants.EXPECT().
			Get(gomock.Any(), types.ModuleNameGov, "expedited_voting_period").
			Return(storage.Constant{Module: types.ModuleNameGov, Name: "expedited_voting_period", Value: "86400000000000"}, nil).
			AnyTimes()
		proposals.EXPECT().
			GetByID(gomock.Any(), uint64(5)).
			Return(&storage.Proposal{
				Id:             5,
				Status:         types.ProposalStatusActive,
				Expedited:      new(false),
				ActivationTime: &activation,
				EndTime:        &regularEnd,
			}, nil).
			AnyTimes()

		var saved []*storage.Proposal
		tx.EXPECT().
			SaveProposals(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, items ...*storage.Proposal) (int64, error) {
				saved = items
				return 0, nil
			}).
			Times(1)

		err := rollbackProposals(t.Context(), tx, constants, proposals, []storage.Event{{
			Type: types.EventTypeActiveProposal,
			Data: map[string]string{
				"proposal_id":     "5",
				"proposal_result": "expedited_proposal_rejected",
				"proposal_log":    "expedited proposal converted to regular",
			},
		}})
		require.NoError(t, err)

		require.Len(t, saved, 1)
		require.EqualValues(t, 5, saved[0].Id)
		require.True(t, saved[0].IsExpedited())
		require.NotNil(t, saved[0].EndTime)
		require.Equal(t, activation.Add(24*time.Hour), *saved[0].EndTime)
		// nothing else of the proposal is touched
		require.Empty(t, saved[0].Status)
		require.Nil(t, saved[0].ActivationTime)
	})

	t.Run("other events are ignored", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		tx := mock.NewMockTransaction(ctrl)
		constants := mock.NewMockIConstant(ctrl)
		proposals := mock.NewMockIProposal(ctrl)

		err := rollbackProposals(t.Context(), tx, constants, proposals, []storage.Event{
			{Type: types.EventTypeActiveProposal, Data: map[string]string{"proposal_id": "1", "proposal_result": "proposal_passed"}},
			{Type: types.EventTypeInactiveProposal, Data: map[string]string{"proposal_id": "2", "proposal_result": "proposal_dropped"}},
			{Type: types.EventTypeProposalVote, Data: map[string]string{"proposal_id": "3"}},
		})
		require.NoError(t, err)
	})
}
