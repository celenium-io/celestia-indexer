// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"

	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestParseTxs_EmptyTxsResults(t *testing.T) {
	block := types.BlockData{
		ResultBlockResults: types.ResultBlockResults{
			TxsResults: make([]*types.ResponseDeliverTx, 0),
		},
	}

	resultTxs, err := parseTxs(block)

	assert.NoError(t, err)
	assert.Empty(t, resultTxs)
}

func TestParseTxs_SuccessTx(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{},
		Log:       "[]",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    nil,
		Codespace: "celestia-explorer",
	}
	block, now := testsuite.CreateTestBlock(txRes, 3)

	resultTxs, err := parseTxs(block)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 3)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusSuccess, f.Status)
	assert.Equal(t, "", f.Error)
	assert.Equal(t, int64(12000), f.GasWanted)
	assert.Equal(t, int64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}

func TestParseTxs_FailedTx(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code:      1,
		Data:      []byte{},
		Log:       "something wierd happened",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    nil,
		Codespace: "celestia-explorer",
	}
	block, now := testsuite.CreateTestBlock(txRes, 1)

	resultTxs, err := parseTxs(block)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusFailed, f.Status)
	assert.Equal(t, "something wierd happened", f.Error)
	assert.Equal(t, int64(12000), f.GasWanted)
	assert.Equal(t, int64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}

func TestParseTxs_FailedTxWithNonstandardErrorCode(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code:      300,
		Data:      []byte{},
		Log:       "something unusual happened",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    nil,
		Codespace: "celestia-explorer",
	}
	block, now := testsuite.CreateTestBlock(txRes, 1)

	resultTxs, err := parseTxs(block)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusFailed, f.Status)
	assert.Equal(t, "something unusual happened", f.Error)
	assert.Equal(t, int64(12000), f.GasWanted)
	assert.Equal(t, int64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}
