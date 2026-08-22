// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"time"

	ic "github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/node/types"
	"go.uber.org/mock/gomock"
)

func (s *ModuleTestSuite) runReceiveGenesis(cfg *ic.Indexer, initialHeight int64) *Module {
	s.InitApi(func() {
		s.api.EXPECT().
			Genesis(gomock.Any()).
			Return(nodeTypes.Genesis{InitialHeight: initialHeight}, nil).
			Times(1)
		s.cosmosApi.EXPECT().
			ModuleAccounts(gomock.Any()).
			Return(nil, nil).
			Times(1)
	})

	receiverModule := s.createModuleEmptyState(cfg)

	ctx, cancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- receiverModule.receiveGenesis(ctx)
	}()

	receiverModule.MustInput(GenesisDoneInput).Push(struct{}{})

	select {
	case err := <-done:
		s.Require().NoError(err)
	case <-time.After(5 * time.Second):
		s.FailNow("receiveGenesis did not return in time")
	}

	return receiverModule
}

// Default start_level + initial_height=1 must still start at block 1 (level=receivedLevel=0).
func (s *ModuleTestSuite) TestReceiveGenesis_DefaultStartLevel() {
	receiverModule := s.runReceiveGenesis(&ic.Indexer{Name: testIndexerName}, 1)

	s.Require().EqualValues(0, receiverModule.level)
	s.Require().EqualValues(0, receiverModule.receivedLevel)
	s.Require().Equal(receiverModule.level, receiverModule.receivedLevel,
		"level and receivedLevel must stay in sync or the sequencer will skip the first block")
}

// initial_height > 1 (re-genesis) must start fetching from initial_height, not block 1.
func (s *ModuleTestSuite) TestReceiveGenesis_InitialHeightAboveOne() {
	receiverModule := s.runReceiveGenesis(&ic.Indexer{Name: testIndexerName}, 1000)

	s.Require().EqualValues(999, receiverModule.level)
	s.Require().EqualValues(999, receiverModule.receivedLevel)
}

// An explicit start_level past initial_height (skip-ahead re-sync) must be preserved.
func (s *ModuleTestSuite) TestReceiveGenesis_ExplicitStartLevelOverride() {
	receiverModule := s.runReceiveGenesis(&ic.Indexer{Name: testIndexerName, StartLevel: 500_000}, 1)

	s.Require().EqualValues(500_000, receiverModule.level)
	s.Require().EqualValues(500_000, receiverModule.receivedLevel)
}

// A start_level below initial_height-1 must be clamped up: blocks below initial_height don't exist.
func (s *ModuleTestSuite) TestReceiveGenesis_StartLevelBelowInitialHeight() {
	receiverModule := s.runReceiveGenesis(&ic.Indexer{Name: testIndexerName, StartLevel: 500}, 1000)

	s.Require().EqualValues(999, receiverModule.level)
	s.Require().EqualValues(999, receiverModule.receivedLevel)
}
