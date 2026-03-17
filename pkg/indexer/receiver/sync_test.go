// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"sort"
	"time"

	ic "github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"
	"github.com/pkg/errors"
	"go.uber.org/mock/gomock"
)

func (s *ModuleTestSuite) TestModule_SyncGracefullyStops() {
	s.InitApi(func() {
		s.api.EXPECT().
			CurrentHead(gomock.Any()).
			Return(0, errors.New("service is down")).
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

// TestPassBlocks_RollbackCancelsInFlightGoroutines verifies that calling
// cancelReadBlocks() (as startRollback does) stops in-flight fetch goroutines
// even when passBlocks was invoked with a never-cancelled parent context
// (the live() code path). Without the fix, goroutines would use the parent
// context and ignore the cancellation, leaking blocks into r.blocks after
// clearChannel drains it.
func (s *ModuleTestSuite) TestPassBlocks_RollbackCancelsInFlightGoroutines() {
	const (
		fetchConcurrency = 3
		bulkSize         = 2
		head             = types.Level(6) // 3 batches: [1,2], [3,4], [5,6]
	)

	s.InitApi(nil)

	receiverModule := s.createModuleEmptyState(&ic.Indexer{
		RequestBulkSize:  bulkSize,
		FetchConcurrency: fetchConcurrency,
	})
	// Use an unbuffered channel so r.blocks <- block in the callback always
	// blocks (no reader). This ensures the select in fetchBatch only has one
	// ready case after cancellation: <-ctx.Done(). Without this, a buffered
	// channel with free capacity would make the select nondeterministic.
	receiverModule.blocks = make(chan *types.BlockData)

	// fetchStarted is signalled by each goroutine once it is inside BlockBulkDataStream.
	fetchStarted := make(chan struct{}, fetchConcurrency)
	// releaseAfterCancel is closed after cancelReadBlocks() to unblock goroutines.
	releaseAfterCancel := make(chan struct{})

	s.api.EXPECT().
		BlockBulkDataStream(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f func(types.BlockData) error, lvls ...types.Level) error {
			fetchStarted <- struct{}{}
			select {
			case <-releaseAfterCancel:
			case <-ctx.Done():
				return ctx.Err()
			}
			// Attempt to push a block after cancellation.
			// With the fix, the select inside fetchBatch chooses <-ctx.Done()
			// and f returns an error without writing to r.blocks.
			_ = f(types.BlockData{ResultBlock: getResultBlock(lvls[0])})
			return nil
		}).
		Times(fetchConcurrency)

	// Simulate live() mode: passBlocks receives a never-cancelled parent ctx.
	done := make(chan struct{})
	go func() {
		defer close(done)
		receiverModule.passBlocks(context.Background(), head)
	}()

	// Wait until all goroutines are inside BlockBulkDataStream —
	// at this point cancelReadBlocks is registered and the goroutines are past
	// their own ctx.Done() check, so timing is deterministic.
	for i := 0; i < fetchConcurrency; i++ {
		select {
		case <-fetchStarted:
		case <-time.After(5 * time.Second):
			s.FailNow("fetch goroutines did not start in time")
		}
	}

	// Simulate startRollback: cancel in-flight fetches and drain the channel.
	receiverModule.cancelReadBlocks()
	clearChannel(receiverModule.blocks)

	// Now release the goroutines so they try to call f() and push a block.
	close(releaseAfterCancel)

	// passBlocks must return within a bounded time after cancellation.
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		s.FailNow("passBlocks did not return after cancelReadBlocks")
	}

	// The channel must be empty: no blocks should have been written after the
	// rollback drain.
	s.Require().Empty(receiverModule.blocks, "no blocks must appear in the channel after rollback cancellation")
}

func (s *ModuleTestSuite) TestModule_SyncReadsBlocks() {
	const blockCount = 5
	s.InitApi(func() {
		s.api.EXPECT().
			CurrentHead(gomock.Any()).
			Return(5, nil).
			Times(1)
	})

	receiverModule := s.createModuleEmptyState(&ic.Indexer{
		Name:            cfgDefault.Name,
		BlockPeriod:     cfgDefault.BlockPeriod,
		RequestBulkSize: 10,
	})

	levels := make([]types.Level, blockCount)
	bulkResult := make([]*types.BlockData, blockCount)
	for i := 0; i < blockCount; i++ {
		levels[i] = types.Level(i + 1)
		bulkResult[i] = &types.BlockData{
			ResultBlock:        getResultBlock(levels[i]),
			ResultBlockResults: getResultBlockResults(levels[i]),
		}
	}

	s.api.EXPECT().
		BlockBulkDataStream(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f func(types.BlockData) error, l ...types.Level) error {
			for i := range bulkResult {
				receiverModule.blocks <- bulkResult[i]
			}
			return nil
		}).
		Times(1)

	ctx, cancelCtx := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer cancelCtx()

	go receiverModule.sync(ctx)

	defer close(receiverModule.blocks)

	syncedBlockData := make([]*types.BlockData, blockCount)
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
