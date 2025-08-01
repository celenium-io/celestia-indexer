// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	dCtx "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	tmTypes "github.com/cometbft/cometbft/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createModules(t *testing.T) (modules.BaseModule, string, Module) {
	writerModule := modules.New("writer-module")
	outputName := "write"
	writerModule.CreateOutput(outputName)
	parserModule := NewModule(config.Indexer{})

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
		MessageTypes:       storageTypes.NewMsgTypeBitMask(),
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
		ProposerAddress:    types.Hex{0x0, 0x0, 0x1, 0x9}.String(),
		ChainId:            "celestia-explorer-test",
		Txs:                make([]storage.Tx, 0),
		Events:             make([]storage.Event, 0),
		BlockSignatures: []storage.BlockSignature{
			{
				Height: 999,
				Validator: &storage.Validator{
					ConsAddress: "960AA0366B254E1EA79BDA467EB3AA5C97CBA5AE",
				},
				Time: time.Time{},
			},
		},
		Stats: storage.BlockStats{
			Id:            0,
			Height:        100,
			Time:          time.Time{},
			TxCount:       0,
			EventsCount:   0,
			BlobsSize:     0,
			SupplyChange:  decimal.Zero,
			InflationRate: decimal.Zero,
			Fee:           decimal.Zero,
			Rewards:       decimal.Zero,
			Commissions:   decimal.Zero,
		},
	}
}

func getBlock() types.BlockData {
	return types.BlockData{
		ResultBlock: types.ResultBlock{
			BlockID: types.BlockId{
				Hash: types.Hex{0x0, 0x0, 0x0, 0x2},
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
					LastBlockID: types.BlockId{
						Hash: types.Hex{0x0, 0x0, 0x0, 0x1},
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
				LastCommit: &types.Commit{
					Height: 999,
					Round:  1,
					Signatures: []tmTypes.CommitSig{
						{
							BlockIDFlag:      tmTypes.BlockIDFlagCommit,
							ValidatorAddress: testHashAddress,
							Timestamp:        time.Time{},
							Signature:        testsuite.MustHexDecode("0011"),
						},
					},
				},
			},
		},
		ResultBlockResults: types.ResultBlockResults{
			Height:              100,
			TxsResults:          nil,
			FinalizeBlockEvents: nil,
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

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
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

			received, ok := msg.(*dCtx.Context)
			assert.Truef(t, ok, "invalid message type: %T", msg)

			expectedBlock := getExpectedBlock()
			assert.Equal(t, expectedBlock, *received.Block)
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

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*1)
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

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*1)
	defer cancel()

	parserModule.Start(ctx)

	block := getBlock()
	block.Block.Txs = tmTypes.Txs{
		// unfinished sequence of tx bytes
		{10, 171, 1, 10, 168, 1, 10, 35, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98},
	}
	block.Block.SquareSize = 1
	block.TxsResults = []*types.ResponseDeliverTx{
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

func getBlockByHeight(height uint64) (types.BlockData, error) {
	blockFile, err := os.Open(fmt.Sprintf("../../../test/json/block_%d.json", height))
	if err != nil {
		return types.BlockData{}, err
	}
	defer blockFile.Close()

	var block types.ResultBlock
	if err := json.NewDecoder(blockFile).Decode(&block); err != nil {
		return types.BlockData{}, err
	}

	blockResultsFile, err := os.Open(fmt.Sprintf("../../../test/json/results_%d.json", height))
	if err != nil {
		return types.BlockData{}, err
	}
	defer blockResultsFile.Close()

	var blockResults types.ResultBlockResults
	if err := json.NewDecoder(blockResultsFile).Decode(&blockResults); err != nil {
		return types.BlockData{}, err
	}

	return types.BlockData{
		ResultBlock:        block,
		ResultBlockResults: blockResults,
	}, nil
}

func TestModule_1768659(t *testing.T) {
	writerModule, writerOutputName, parserModule := createModules(t)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*1)
	defer cancel()

	parserModule.Start(ctx)

	block, err := getBlockByHeight(1768659)
	require.NoError(t, err)
	writerModule.MustOutput(writerOutputName).Push(block)

	<-ctx.Done()

	err = parserModule.Close()
	require.NoError(t, err)
}
