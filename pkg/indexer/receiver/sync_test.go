// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"sort"
	"time"

	ic "github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"
	"github.com/pkg/errors"
	"go.uber.org/mock/gomock"
)

func (s *ModuleTestSuite) TestModule_SyncGracefullyStops() {
	s.InitApi(func() {
		s.api.EXPECT().
			Status(gomock.Any()).
			Return(nodeTypes.Status{}, errors.New("service is down")).
			MaxTimes(1)
	})

	receiverModule := s.createModule()

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	stopperModule := stopper.NewModule(cancelCtx)
	err := stopperModule.AttachTo(&receiverModule, StopOutput, stopper.InputName)
	s.Require().NoError(err)

	stopperCtx, stopperCtxCancel := context.WithCancel(context.Background())
	defer stopperCtxCancel()

	stopperModule.Start(stopperCtx)

	go receiverModule.sync(ctx)

	defer close(receiverModule.blocks)

	for range ctx.Done() {
		s.Require().ErrorIs(context.Canceled, ctx.Err())
		return
	}
}

func getResultBlock(level types.Level) types.ResultBlock {
	return types.ResultBlock{
		Block: &types.Block{
			Header: types.Header{
				Height: int64(level),
			},
		},
	}
}

func getResultBlockResults(level types.Level) types.ResultBlockResults {
	return types.ResultBlockResults{
		Height: level,
	}
}

func (s *ModuleTestSuite) TestModule_SyncReadsBlocks() {
	const blockCount = 5
	s.InitApi(func() {
		s.api.EXPECT().
			Status(gomock.Any()).
			Return(nodeTypes.Status{
				NodeInfo: nodeTypes.NodeInfo{
					ProtocolVersion: nodeTypes.ProtocolVersion{
						App: 4,
					},
				},
				SyncInfo: nodeTypes.SyncInfo{
					LatestBlockHash:   nil,
					LatestBlockHeight: 5,
				},
			}, nil).
			MaxTimes(1)

		levels := make([]types.Level, blockCount)
		bulkResult := make([]types.BlockData, blockCount)
		for i := types.Level(1); i <= blockCount; i++ {
			levels[i-1] = i
			bulkResult[i-1] = types.BlockData{
				ResultBlock:        getResultBlock(i),
				ResultBlockResults: getResultBlockResults(i),
			}
		}

		s.api.EXPECT().
			BlockBulkData(gomock.Any(), gomock.Any()).
			Return(bulkResult, nil).
			Times(1)
	})

	receiverModule := s.createModuleEmptyState(&ic.Indexer{
		Name:        cfgDefault.Name,
		BlockPeriod: cfgDefault.BlockPeriod,
	})

	ctx, cancelCtx := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer cancelCtx()

	go receiverModule.sync(ctx)

	defer close(receiverModule.blocks)

	syncedBlockData := make([]types.BlockData, blockCount)
	index := 0
	for b := range receiverModule.blocks {
		syncedBlockData[index] = b
		index++

		if index == 5 {
			break
		}
	}

	sort.Slice(syncedBlockData, func(i, j int) bool {
		return syncedBlockData[i].Height < syncedBlockData[j].Height
	})

	for i := types.Level(1); i <= blockCount; i++ {
		s.Require().EqualValues(i, syncedBlockData[i-1].Height)
	}
}
