package receiver

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/rollback"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
	"time"
)

type blockConciseData struct {
	level int64
	hash  types.Hex
}

const (
	asc int = 1 + iota
	desc
	random
)

var hashOf1000Block, _ = types.HexFromString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")

func createBlocks(order int, data ...blockConciseData) []types.BlockData {
	res := make([]types.BlockData, len(data))

	prevBlockHash := hashOf1000Block

	for i, d := range data {
		res[i].Height = types.Level(d.level)
		res[i].BlockID.Hash = d.hash
		res[i].Block = &types.Block{
			Header: types.Header{
				Height: d.level,
				LastBlockID: types.BlockId{
					Hash: prevBlockHash,
				},
			},
		}
		prevBlockHash = d.hash
	}

	if order == asc {
		return res
	}

	if order == desc {
		sort.Slice(res, func(i, j int) bool {
			return res[i].Height > res[j].Height
		})
	}

	if order == random {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(res), func(i, j int) { res[i], res[j] = res[j], res[i] })
	}

	return res
}

func Test_createBlock(t *testing.T) {
	tests := []struct {
		name       string
		order      int
		blocksData []blockConciseData
		want       []blockConciseData // use blockConciseData for brevity
		wantRandom bool
	}{
		{
			name:       "asc order",
			order:      asc,
			blocksData: blocksData,
			want:       blocksData,
		},
		{
			name:       "desc order",
			order:      desc,
			blocksData: blocksData,
			want: []blockConciseData{
				{level: 5, hash: []byte{0x05}},
				{level: 4, hash: []byte{0x04}},
				{level: 3, hash: []byte{0x03}},
				{level: 2, hash: []byte{0x02}},
				{level: 1, hash: []byte{0x01}},
			},
		},
		{
			name:       "random order",
			order:      random,
			blocksData: blocksData,
			want:       nil,
			wantRandom: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := createBlocks(tt.order, tt.blocksData...)

			assert.Len(t, blocks, len(tt.blocksData))
			if tt.order == random {
				return
			}

			for i, b := range blocks {
				assert.Equal(t, types.Level(tt.want[i].level), b.Height)
				assert.Equal(t, tt.want[i].level, b.Block.Height)
				assert.Equal(t, tt.want[i].hash, b.BlockID.Hash)
			}
		})
	}
}

var blocksData = []blockConciseData{
	{level: 1, hash: []byte{0x01}},
	{level: 2, hash: []byte{0x02}},
	{level: 3, hash: []byte{0x03}},
	{level: 4, hash: []byte{0x04}},
	{level: 5, hash: []byte{0x05}},
}

func (s *ModuleTestSuite) TestModule_SequencerOnEmptyState() {
	s.InitApi(nil)

	receiverModule := s.createModuleEmptyState(nil)

	blocksReaderModule := modules.New("ordered-blocks-reader")
	const orderedBlocksChannel = "ordered-blocks"
	blocksReaderModule.CreateInput(orderedBlocksChannel)
	err := blocksReaderModule.AttachTo(&receiverModule, BlocksOutput, orderedBlocksChannel)
	s.Require().NoError(err)

	tests := []struct {
		name   string
		order  int
		blocks []types.BlockData
		want   []blockConciseData
	}{
		{
			name:   "asc order",
			blocks: createBlocks(asc, blocksData...),
			want:   blocksData,
		},
		{
			name:   "desc order",
			blocks: createBlocks(desc, blocksData...),
			want:   blocksData,
		},
		{
			name:   "random order",
			blocks: createBlocks(random, blocksData...),
			want:   blocksData,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelCtx()

			receiverModule.setLevel(0, nil)
			go receiverModule.sequencer(ctx)

			for _, b := range tt.blocks {
				receiverModule.blocks <- b
			}

			index := 0
		out:
			for {
				select {
				case <-ctx.Done():
					s.T().Error("stop by cancelled context")
					return
				case ob := <-blocksReaderModule.MustInput(orderedBlocksChannel).Listen():
					orderedBlock := ob.(types.BlockData)
					s.Require().EqualValues(blocksData[index].level, orderedBlock.Height)
					s.Require().EqualValues(blocksData[index].level, orderedBlock.Block.Height)
					s.Require().EqualValues(blocksData[index].hash, orderedBlock.BlockID.Hash)
					index++

					if index == 5 {
						break out
					}
				}
			}

			receiverLevel, receiverHash := receiverModule.Level()
			s.Require().EqualValues(types.Level(5), receiverLevel)
			s.Require().EqualValues([]byte{0x05}, receiverHash)
		})
	}
}

func (s *ModuleTestSuite) TestModule_SequencerOnNonEmptyState() {
	s.InitApi(nil)

	receiverModule := s.createModule()

	blocksReaderModule := modules.New("ordered-blocks-reader")
	const orderedBlocksChannel = "ordered-blocks"
	blocksReaderModule.CreateInput(orderedBlocksChannel)
	err := blocksReaderModule.AttachTo(&receiverModule, BlocksOutput, orderedBlocksChannel)
	s.Require().NoError(err)

	blocksData := []blockConciseData{
		{level: 1001, hash: []byte{0x10, 0x10, 0x10, 0x01}},
		{level: 1002, hash: []byte{0x10, 0x10, 0x10, 0x02}},
		{level: 1003, hash: []byte{0x10, 0x10, 0x10, 0x03}},
		{level: 1004, hash: []byte{0x10, 0x10, 0x10, 0x04}},
		{level: 1005, hash: []byte{0x10, 0x10, 0x10, 0x05}},
	}

	tests := []struct {
		name   string
		order  int
		blocks []types.BlockData
		want   []blockConciseData
	}{
		{
			name:   "asc order",
			blocks: createBlocks(asc, blocksData...),
			want:   blocksData,
		},
		{
			name:   "desc order",
			blocks: createBlocks(desc, blocksData...),
			want:   blocksData,
		},
		{
			name:   "random order",
			blocks: createBlocks(random, blocksData...),
			want:   blocksData,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelCtx()

			receiverModule.setLevel(1000, hashOf1000Block)
			go receiverModule.sequencer(ctx)

			for _, b := range tt.blocks {
				receiverModule.blocks <- b
			}

			index := 0
		out:
			for {
				select {
				case <-ctx.Done():
					s.T().Error("stop by cancelled context")
					return
				case ob := <-blocksReaderModule.MustInput(orderedBlocksChannel).Listen():
					orderedBlock := ob.(types.BlockData)
					s.Require().EqualValues(blocksData[index].level, orderedBlock.Height)
					s.Require().EqualValues(blocksData[index].level, orderedBlock.Block.Height)
					s.Require().EqualValues(blocksData[index].hash, orderedBlock.BlockID.Hash)
					index++

					if index == 5 {
						break out
					}
				}
			}

			receiverLevel, receiverHash := receiverModule.Level()
			s.Require().EqualValues(types.Level(1005), receiverLevel)
			s.Require().EqualValues([]byte{0x10, 0x10, 0x10, 0x05}, receiverHash)
		})
	}
}

func (s *ModuleTestSuite) TestModule_SequencerGracefullyStops() {
	s.InitApi(nil)

	receiverModule := s.createModuleEmptyState(nil)

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	stopperModule := stopper.NewModule(cancelCtx)
	err := stopperModule.AttachTo(&receiverModule, StopOutput, stopper.InputName)
	s.Require().NoError(err)

	stopperCtx, stopperCtxCancel := context.WithCancel(context.Background())
	defer stopperCtxCancel()

	stopperModule.Start(stopperCtx)
	go receiverModule.sequencer(ctx)

	close(receiverModule.blocks)

	for range ctx.Done() {
		s.Require().ErrorIs(context.Canceled, ctx.Err())
		return
	}
}

func (s *ModuleTestSuite) TestModule_SequencerCallsRollback() {
	s.InitApi(nil)

	receiverModule := s.createModule()

	blocksReaderModule := modules.New("ordered-blocks-reader")
	const orderedBlocksChannel = "ordered-blocks"
	blocksReaderModule.CreateInput(orderedBlocksChannel)
	err := blocksReaderModule.AttachTo(&receiverModule, BlocksOutput, orderedBlocksChannel)
	s.Require().NoError(err)

	rollbackModule := modules.New("rollback")
	rollbackModule.CreateInput(rollback.InputName)
	rollbackModule.CreateOutput(rollback.OutputName)
	err = receiverModule.AttachTo(&rollbackModule, rollback.OutputName, RollbackInput)
	s.Require().NoError(err)
	err = rollbackModule.AttachTo(&receiverModule, RollbackOutput, rollback.InputName)
	s.Require().NoError(err)

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-rollbackModule.MustInput(rollback.InputName).Listen():
				rollbackModule.MustOutput(rollback.OutputName).Push(storage.State{
					LastHeight: types.Level(4),
					LastHash:   []byte{0x04},
				})
			}
		}
	}()

	receiverModule.setLevel(0, nil)
	go receiverModule.sequencer(ctx)

	blocks := createBlocks(asc, blocksData...)

	blocks[4].Block.LastBlockID.Hash = []byte{0x09} // not equal to hash that is saved in the state of receiver.

	for _, b := range blocks {
		receiverModule.blocks <- b
	}

	index := 0
out:
	for {
		select {
		case <-ctx.Done():
			s.T().Error("stop by cancelled context")
			return
		case ob := <-blocksReaderModule.MustInput(orderedBlocksChannel).Listen():
			orderedBlock := ob.(types.BlockData)
			s.Require().EqualValues(blocksData[index].level, orderedBlock.Height)
			s.Require().EqualValues(blocksData[index].level, orderedBlock.Block.Height)
			s.Require().EqualValues(blocksData[index].hash, orderedBlock.BlockID.Hash)
			index++

			if index == 4 {
				break out
			}
		}
	}

	receiverLevel, receiverHash := receiverModule.Level()
	s.Require().EqualValues(types.Level(4), receiverLevel)
	s.Require().EqualValues([]byte{0x04}, receiverHash)
}

func (s *ModuleTestSuite) TestModule_SequencerCallsRollbackWithinPreSavedBlocks() {
	s.InitApi(nil)

	receiverModule := s.createModule()

	blocksReaderModule := modules.New("ordered-blocks-reader")
	const orderedBlocksChannel = "ordered-blocks"
	blocksReaderModule.CreateInput(orderedBlocksChannel)
	err := blocksReaderModule.AttachTo(&receiverModule, BlocksOutput, orderedBlocksChannel)
	s.Require().NoError(err)

	rollbackModule := modules.New("rollback")
	rollbackModule.CreateInput(rollback.InputName)
	rollbackModule.CreateOutput(rollback.OutputName)
	err = receiverModule.AttachTo(&rollbackModule, rollback.OutputName, RollbackInput)
	s.Require().NoError(err)
	err = rollbackModule.AttachTo(&receiverModule, RollbackOutput, rollback.InputName)
	s.Require().NoError(err)

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-rollbackModule.MustInput(rollback.InputName).Listen():
				rollbackModule.MustOutput(rollback.OutputName).Push(storage.State{
					LastHeight: types.Level(3),
					LastHash:   []byte{0x03},
				})
			}
		}
	}()

	receiverModule.setLevel(0, nil)
	go receiverModule.sequencer(ctx)

	blocks := createBlocks(desc, blocksData...)

	blocks[1].Block.LastBlockID.Hash = []byte{0x09} // not equal to hash that is saved in the state of receiver.

	for _, b := range blocks {
		receiverModule.blocks <- b
	}

	index := 0
out:
	for {
		select {
		case <-ctx.Done():
			s.T().Error("stop by cancelled context")
			return
		case ob := <-blocksReaderModule.MustInput(orderedBlocksChannel).Listen():
			orderedBlock := ob.(types.BlockData)
			s.Require().EqualValues(blocksData[index].level, orderedBlock.Height)
			s.Require().EqualValues(blocksData[index].level, orderedBlock.Block.Height)
			s.Require().EqualValues(blocksData[index].hash, orderedBlock.BlockID.Hash)
			index++

			if index == 3 {
				break out
			}
		}
	}

	receiverLevel, receiverHash := receiverModule.Level()
	s.Require().EqualValues(types.Level(3), receiverLevel)
	s.Require().EqualValues([]byte{0x03}, receiverHash)
}
