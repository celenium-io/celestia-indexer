// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
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

	workersCtx, cancelWorkers := context.WithCancel(ctx)
	receiverModule.cancelWorkers = cancelWorkers
	receiverModule.pool.Start(workersCtx)

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
				SyncInfo: nodeTypes.SyncInfo{
					LatestBlockHash:   nil,
					LatestBlockHeight: 5,
				},
			}, nil).
			MaxTimes(1)

		for i := types.Level(1); i <= blockCount; i++ {
			s.api.EXPECT().
				BlockData(gomock.Any(), i).
				Return(types.BlockData{
					ResultBlock:        getResultBlock(i),
					ResultBlockResults: getResultBlockResults(i),
				}, nil).
				MaxTimes(1).
				MinTimes(1)
		}
	})

	receiverModule := s.createModuleEmptyState(&ic.Indexer{
		Name:         cfgDefault.Name,
		ThreadsCount: blockCount,
		BlockPeriod:  cfgDefault.BlockPeriod,
	})

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	workersCtx, cancelWorkers := context.WithCancel(ctx)
	receiverModule.cancelWorkers = cancelWorkers
	receiverModule.pool.Start(workersCtx)

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
