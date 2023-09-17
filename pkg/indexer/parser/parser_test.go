package parser

import (
	"context"
	"testing"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/bytes"
	tmTypes "github.com/tendermint/tendermint/types"
)

func createModules(t *testing.T) (modules.BaseModule, string, Module) {
	writerModule := modules.New("writer-module")
	outputName := "write"
	writerModule.CreateOutput(outputName)
	parserModule := NewModule()

	err := parserModule.AttachTo(&writerModule, outputName, InputName)
	assert.NoError(t, err)

	return writerModule, outputName, parserModule
}

func getExpectedBlock() storage.Block {
	return storage.Block{
		Id:                 0,
		Height:             100,
		Time:               time.Time{},
		VersionBlock:       1,
		VersionApp:         2,
		MessageTypes:       storageTypes.MsgTypeBits{},
		Hash:               types.Hex{0x0, 0x0, 0x0, 0x2},
		ParentHash:         types.Hex{0x0, 0x0, 0x0, 0x1},
		LastCommitHash:     types.Hex{0x0, 0x0, 0x1, 0x1},
		DataHash:           types.Hex{0x0, 0x0, 0x1, 0x2},
		ValidatorsHash:     types.Hex{0x0, 0x0, 0x1, 0x3},
		NextValidatorsHash: types.Hex{0x0, 0x0, 0x1, 0x4},
		ConsensusHash:      types.Hex{0x0, 0x0, 0x1, 0x5},
		AppHash:            types.Hex{0x0, 0x0, 0x1, 0x6},
		LastResultsHash:    types.Hex{0x0, 0x0, 0x1, 0x7},
		EvidenceHash:       types.Hex{0x0, 0x0, 0x1, 0x8},
		ProposerAddress:    types.Hex{0x0, 0x0, 0x1, 0x9},
		ChainId:            "celestia-explorer-test",
		Txs:                make([]storage.Tx, 0),
		Events:             make([]storage.Event, 0),
		Addresses:          make([]storage.Address, 0),
		Stats: storage.BlockStats{
			Id:          0,
			Height:      100,
			Time:        time.Time{},
			TxCount:     0,
			EventsCount: 0,
			BlobsSize:   0,
			// SupplyChange: decimal.Zero,
			// InflationRate: decimal.Zero,
			Fee: decimal.Zero,
		},
	}
}

func getBlock() types.BlockData {
	return types.BlockData{
		ResultBlock: types.ResultBlock{
			BlockID: tmTypes.BlockID{
				Hash: bytes.HexBytes{0x0, 0x0, 0x0, 0x2},
				PartSetHeader: tmTypes.PartSetHeader{
					Total: 0,
					Hash:  nil,
				},
			},
			Block: &types.Block{
				Header: types.Header{
					Version: types.Consensus{
						Block: 1,
						App:   2,
					},
					ChainID: "celestia-explorer-test",
					Height:  1000,
					Time:    time.Time{},
					LastBlockID: tmTypes.BlockID{
						Hash: bytes.HexBytes{0x0, 0x0, 0x0, 0x1},
						PartSetHeader: tmTypes.PartSetHeader{
							Total: 0,
							Hash:  nil,
						},
					},
					LastCommitHash:     types.Hex{0x0, 0x0, 0x1, 0x1},
					DataHash:           types.Hex{0x0, 0x0, 0x1, 0x2},
					ValidatorsHash:     types.Hex{0x0, 0x0, 0x1, 0x3},
					NextValidatorsHash: types.Hex{0x0, 0x0, 0x1, 0x4},
					ConsensusHash:      types.Hex{0x0, 0x0, 0x1, 0x5},
					AppHash:            types.Hex{0x0, 0x0, 0x1, 0x6},
					LastResultsHash:    types.Hex{0x0, 0x0, 0x1, 0x7},
					EvidenceHash:       types.Hex{0x0, 0x0, 0x1, 0x8},
					ProposerAddress:    types.Hex{0x0, 0x0, 0x1, 0x9},
				},
				Data: types.Data{
					Txs:        nil,
					SquareSize: 0,
				},
				LastCommit: nil,
			},
		},
		ResultBlockResults: types.ResultBlockResults{
			Height:           100,
			TxsResults:       nil,
			BeginBlockEvents: nil,
			EndBlockEvents:   nil,
			ValidatorUpdates: nil,
			ConsensusParamUpdates: &types.ConsensusParams{
				Block: &types.BlockParams{
					MaxBytes: 0,
					MaxGas:   0,
				},
				Evidence: &types.EvidenceParams{
					MaxAgeNumBlocks: 0,
					MaxAgeDuration:  0,
					MaxBytes:        0,
				},
				Validator: &types.ValidatorParams{
					PubKeyTypes: nil,
				},
				Version: &types.VersionParams{
					AppVersion: 0,
				},
			},
		},
	}
}

func TestParserModule_Success(t *testing.T) {
	writerModule, outputName, parserModule := createModules(t)

	readerModule := modules.New("reader-module")
	readerInputName := "read"
	readerModule.CreateInput(readerInputName)

	err := readerModule.AttachTo(&parserModule, OutputName, readerInputName)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	parserModule.Start(ctx)

	block := getBlock()
	writerModule.MustOutput(outputName).Push(block)

	for {
		select {
		case <-ctx.Done():
			t.Error("stop by cancelled context")
		case msg, ok := <-readerModule.MustInput(readerInputName).Listen():
			assert.True(t, ok, "received value should be delivered by successful send operation")

			parsedBlock, ok := msg.(storage.Block)
			assert.Truef(t, ok, "invalid message type: %T", msg)

			expectedBlock := getExpectedBlock()
			assert.Equal(t, expectedBlock, parsedBlock)
			return
		}
	}
}

func TestModule_OnClosedChannel(t *testing.T) {
	_, _, parserModule := createModules(t)

	stopperModule := modules.New("stopper-module")
	stopInputName := "stop-signal"
	stopperModule.CreateInput(stopInputName)

	err := stopperModule.AttachTo(&parserModule, StopOutput, stopInputName)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	parserModule.Start(ctx)

	err = parserModule.MustInput(InputName).Close()
	assert.NoError(t, err)

	for {
		select {
		case <-ctx.Done():
			t.Error("stop by cancelled context")
		case msg := <-stopperModule.MustInput(stopInputName).Listen():
			assert.Equal(t, struct{}{}, msg)
			return
		}
	}
}

func TestModule_OnParseError(t *testing.T) {
	writerModule, writerOutputName, parserModule := createModules(t)

	stopperModule := modules.New("stopper-module")
	stopInputName := "stop-signal"
	stopperModule.CreateInput(stopInputName)

	err := stopperModule.AttachTo(&parserModule, StopOutput, stopInputName)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	parserModule.Start(ctx)

	block := getBlock()
	block.Block.Data.Txs = tmTypes.Txs{
		// unfinished sequence of tx bytes
		{10, 171, 1, 10, 168, 1, 10, 35, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98},
	}
	block.Block.Data.SquareSize = 1
	block.ResultBlockResults.TxsResults = []*types.ResponseDeliverTx{
		{
			Code:      0,
			Data:      []byte{18, 45, 10, 43, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97},
			Log:       "",
			Info:      "",
			GasWanted: 20,
			GasUsed:   10,
			Events:    nil,
			Codespace: "",
		},
	}
	writerModule.MustOutput(writerOutputName).Push(block)

	for {
		select {
		case <-ctx.Done():
			t.Error("stop by cancelled context")
		case msg := <-stopperModule.MustInput(stopInputName).Listen():
			assert.Equal(t, struct{}{}, msg)
			return
		}
	}
}
