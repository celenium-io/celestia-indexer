// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestParseClientRecoveries(t *testing.T) {
	tests := []struct {
		name    string
		changes string
		want    []clientRecovery
		wantErr bool
	}{
		{
			name: "empty",
		}, {
			name:    "one",
			changes: `[{"SubjectClientId":"07-tendermint-1","SubstituteClientId":"07-tendermint-5"}]`,
			want:    []clientRecovery{{"07-tendermint-1", "07-tendermint-5"}},
		}, {
			name:    "many",
			changes: `[{"SubjectClientId":"07-tendermint-1","SubstituteClientId":"07-tendermint-5"},{"SubjectClientId":"06-solomachine-2","SubstituteClientId":"06-solomachine-3"}]`,
			want: []clientRecovery{
				{"07-tendermint-1", "07-tendermint-5"},
				{"06-solomachine-2", "06-solomachine-3"},
			},
		}, {
			name:    "null",
			changes: `null`,
			want:    []clientRecovery(nil),
		}, {
			name:    "entries without subject are dropped",
			changes: `[{"SubstituteClientId":"07-tendermint-5"},{"SubjectClientId":"07-tendermint-1","SubstituteClientId":"07-tendermint-6"}]`,
			want:    []clientRecovery{{"07-tendermint-1", "07-tendermint-6"}},
		}, {
			// mixed v1 proposal: changes hold param changes even when the type is client_update
			name:    "param changes",
			changes: `[{"subspace":"staking","key":"MaxValidators","value":"100"}]`,
			want:    []clientRecovery{},
		}, {
			name:    "object is not supported",
			changes: `{"SubjectClientId":"07-tendermint-1","SubstituteClientId":"07-tendermint-5"}`,
			wantErr: true,
		}, {
			name:    "invalid",
			changes: `[{"SubjectClientId":`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseClientRecoveries([]byte(tt.changes))
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestRecoverIbcClients(t *testing.T) {
	blockTime := time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC)

	newModule := func(ctrl *gomock.Controller) Module {
		return NewModule(nil, mock.NewMockIConstant(ctrl), mock.NewMockIValidator(ctrl), nil, config.Indexer{})
	}

	t.Run("no recoveries", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		tx := mock.NewMockTransaction(ctrl)
		module := newModule(ctrl)

		err := module.recoverIbcClients(t.Context(), tx, nil, sdkSync.NewMap[uint64, *storage.Proposal](), blockTime)
		require.NoError(t, err)
	})

	t.Run("substitutes from applied proposals", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		tx := mock.NewMockTransaction(ctrl)
		module := newModule(ctrl)

		proposals := sdkSync.NewMap[uint64, *storage.Proposal]()
		proposals.Set(1, &storage.Proposal{Id: 1, Status: types.ProposalStatusApplied})
		proposals.Set(2, &storage.Proposal{Id: 2, Status: types.ProposalStatusApplied})
		proposals.Set(3, &storage.Proposal{Id: 3, Status: types.ProposalStatusApplied})
		// not executed in this block: never read
		proposals.Set(4, &storage.Proposal{Id: 4, Status: types.ProposalStatusActive})

		tx.EXPECT().Proposal(gomock.Any(), uint64(1)).Return(storage.Proposal{
			Id:      1,
			Type:    types.ProposalTypeClientUpdate,
			Changes: []byte(`[{"SubjectClientId":"07-tendermint-1","SubstituteClientId":"07-tendermint-5"}]`),
		}, nil)
		tx.EXPECT().Proposal(gomock.Any(), uint64(2)).Return(storage.Proposal{
			Id:      2,
			Type:    types.ProposalTypeClientUpdate,
			Changes: []byte(`[{"SubjectClientId":"07-tendermint-2","SubstituteClientId":"07-tendermint-6"}]`),
		}, nil)
		tx.EXPECT().Proposal(gomock.Any(), uint64(3)).Return(storage.Proposal{
			Id:      3,
			Type:    types.ProposalTypeParamChanged,
			Changes: []byte(`[{"subspace":"staking","key":"MaxValidators","value":"100"}]`),
		}, nil)

		tx.EXPECT().RecoverIbcClient(gomock.Any(), "07-tendermint-1", "07-tendermint-5", blockTime).Return(nil)
		tx.EXPECT().RecoverIbcClient(gomock.Any(), "07-tendermint-2", "07-tendermint-6", blockTime).Return(nil)
		// no proposal knows the substitute: only unfrozen
		tx.EXPECT().RecoverIbcClient(gomock.Any(), "07-tendermint-9", "", blockTime).Return(nil)

		err := module.recoverIbcClients(t.Context(), tx,
			[]string{"07-tendermint-1", "07-tendermint-2", "07-tendermint-9"},
			proposals, blockTime)
		require.NoError(t, err)
	})

	t.Run("unparsable changes don't halt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		tx := mock.NewMockTransaction(ctrl)
		module := newModule(ctrl)

		proposals := sdkSync.NewMap[uint64, *storage.Proposal]()
		proposals.Set(1, &storage.Proposal{Id: 1, Status: types.ProposalStatusApplied})
		tx.EXPECT().Proposal(gomock.Any(), uint64(1)).Return(storage.Proposal{
			Id:      1,
			Type:    types.ProposalTypeClientUpdate,
			Changes: []byte(`[{"SubjectClientId":`),
		}, nil)
		tx.EXPECT().RecoverIbcClient(gomock.Any(), "07-tendermint-1", "", blockTime).Return(nil)

		err := module.recoverIbcClients(t.Context(), tx, []string{"07-tendermint-1"}, proposals, blockTime)
		require.NoError(t, err)
	})

	t.Run("proposal read error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		tx := mock.NewMockTransaction(ctrl)
		module := newModule(ctrl)

		proposals := sdkSync.NewMap[uint64, *storage.Proposal]()
		proposals.Set(1, &storage.Proposal{Id: 1, Status: types.ProposalStatusApplied})
		tx.EXPECT().Proposal(gomock.Any(), uint64(1)).Return(storage.Proposal{}, errors.New("db error"))

		err := module.recoverIbcClients(t.Context(), tx, []string{"07-tendermint-1"}, proposals, blockTime)
		require.Error(t, err)
	})
}
