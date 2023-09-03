package parser

import (
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createBlock(tx nodeTypes.ResponseDeliverTx, count int) (types.BlockData, time.Time) {
	now := time.Now()
	headerBlock := nodeTypes.Block{
		Header: nodeTypes.Header{
			Time: now,
		},
	}

	var txResults []*nodeTypes.ResponseDeliverTx
	for i := 0; i < count; i++ {
		txResults = append(txResults, &tx)
	}

	block := types.BlockData{
		ResultBlock: nodeTypes.ResultBlock{
			Block: &headerBlock,
		},
		ResultBlockResults: nodeTypes.ResultBlockResults{
			TxsResults: txResults,
		},
	}

	return block, now
}

func TestParseTxs_EmptyTxsResults(t *testing.T) {
	block := types.BlockData{
		ResultBlockResults: nodeTypes.ResultBlockResults{
			TxsResults: make([]*nodeTypes.ResponseDeliverTx, 0),
		},
	}

	resultTxs := parseTxs(block)

	assert.Empty(t, resultTxs)
}

func TestParseTxs_SuccessTx(t *testing.T) {
	txRes := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      []byte{},
		Log:       "[]",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    nil,
		Codespace: "celestia-explorer",
	}
	block, now := createBlock(txRes, 3)

	resultTxs := parseTxs(block)

	assert.Len(t, resultTxs, 3)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusSuccess, f.Status)
	assert.Equal(t, "", f.Error)
	assert.Equal(t, uint64(12000), f.GasWanted)
	assert.Equal(t, uint64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}

func TestParseTxs_FailedTx(t *testing.T) {
	txRes := nodeTypes.ResponseDeliverTx{
		Code:      1,
		Data:      []byte{},
		Log:       "something wierd happened",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    nil,
		Codespace: "celestia-explorer",
	}
	block, now := createBlock(txRes, 1)

	resultTxs := parseTxs(block)

	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusFailed, f.Status)
	assert.Equal(t, "something wierd happened", f.Error)
	assert.Equal(t, uint64(12000), f.GasWanted)
	assert.Equal(t, uint64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}

func TestParseTxs_FailedTxWithNonstandardErrorCode(t *testing.T) {
	txRes := nodeTypes.ResponseDeliverTx{
		Code:      300,
		Data:      []byte{},
		Log:       "something unusual happened",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    nil,
		Codespace: "celestia-explorer",
	}
	block, now := createBlock(txRes, 1)

	resultTxs := parseTxs(block)

	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusFailed, f.Status)
	assert.Equal(t, "something unusual happened", f.Error)
	assert.Equal(t, uint64(12000), f.GasWanted)
	assert.Equal(t, uint64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}
