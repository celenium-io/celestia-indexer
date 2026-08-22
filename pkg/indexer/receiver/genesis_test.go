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

// TestReceiveGenesis_DefaultStartLevel verifies that with the default (unset)
// start_level and the network's usual initial_height=1, the receiver still
// begins fetching from block 1 (level=0, receivedLevel=0), exactly as before
// genesis.InitialHeight was taken into account. A regression here would
// silently skip the chain's first block.
func (s *ModuleTestSuite) TestReceiveGenesis_DefaultStartLevel() {
	receiverModule := s.runReceiveGenesis(&ic.Indexer{Name: testIndexerName}, 1)

	s.Require().EqualValues(0, receiverModule.level)
	s.Require().EqualValues(0, receiverModule.receivedLevel)
	s.Require().Equal(receiverModule.level, receiverModule.receivedLevel,
		"level and receivedLevel must stay in sync or the sequencer will skip the first block")
}

// TestReceiveGenesis_InitialHeightAboveOne verifies that when a chain's
// genesis declares initial_height > 1 (e.g. a re-genesis that continues
// height numbering), the receiver starts fetching from initial_height
// instead of from block 1, which doesn't exist on such a chain.
func (s *ModuleTestSuite) TestReceiveGenesis_InitialHeightAboveOne() {
	receiverModule := s.runReceiveGenesis(&ic.Indexer{Name: testIndexerName}, 1000)

	s.Require().EqualValues(999, receiverModule.level)
	s.Require().EqualValues(999, receiverModule.receivedLevel)
}

// TestReceiveGenesis_ExplicitStartLevelOverride verifies that an operator's
// explicit start_level (used to skip ahead in a re-sync) is preserved when
// it is already past the chain's initial_height.
func (s *ModuleTestSuite) TestReceiveGenesis_ExplicitStartLevelOverride() {
	receiverModule := s.runReceiveGenesis(&ic.Indexer{Name: testIndexerName, StartLevel: 500_000}, 1)

	s.Require().EqualValues(500_000, receiverModule.level)
	s.Require().EqualValues(500_000, receiverModule.receivedLevel)
}

// TestReceiveGenesis_StartLevelBelowInitialHeight verifies that a start_level
// lower than the chain's initial_height-1 is clamped up, since blocks below
// initial_height don't exist on the node.
func (s *ModuleTestSuite) TestReceiveGenesis_StartLevelBelowInitialHeight() {
	receiverModule := s.runReceiveGenesis(&ic.Indexer{Name: testIndexerName, StartLevel: 500}, 1000)

	s.Require().EqualValues(999, receiverModule.level)
	s.Require().EqualValues(999, receiverModule.receivedLevel)
}
