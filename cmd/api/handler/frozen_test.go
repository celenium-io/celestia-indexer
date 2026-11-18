// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFrozenState_CachesOnceActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	state := mock.NewMockIState(ctrl)
	f := NewFrozenState(state, testIndexerName)

	// before v11 the state is read on every call
	state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(storage.State{ChainId: "mocha-5", Version: 10}, nil).
		Times(2)
	for range 2 {
		chainId, err := f.activeChainId(t.Context())
		require.NoError(t, err)
		require.Empty(t, chainId)
	}

	// from v11 it is read once and cached
	state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(storage.State{ChainId: "mocha-5", Version: 11}, nil).
		Times(1)
	for range 2 {
		chainId, err := f.activeChainId(t.Context())
		require.NoError(t, err)
		require.Equal(t, "mocha-5", chainId)
	}
}

func TestFrozenState_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	state := mock.NewMockIState(ctrl)
	f := NewFrozenState(state, testIndexerName)

	state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(storage.State{}, errors.New("db is down")).
		Times(1)
	_, err := f.activeChainId(t.Context())
	require.Error(t, err)
}
