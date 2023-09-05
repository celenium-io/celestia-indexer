package parser

import (
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
	tmTypes "github.com/tendermint/tendermint/types"
	"testing"
	"time"
)

var txMsgBeginRedelegate = []byte{10, 252, 1, 10, 225, 1, 10, 42, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 66, 101, 103, 105, 110, 82, 101, 100, 101, 108, 101, 103, 97, 116, 101, 18, 178, 1, 10, 47, 99, 101, 108, 101, 115, 116, 105, 97, 49, 100, 97, 118, 122, 52, 48, 107, 97, 116, 57, 51, 116, 52, 57, 108, 106, 114, 107, 109, 107, 108, 53, 117, 113, 104, 113, 113, 52, 53, 101, 48, 116, 101, 100, 103, 102, 56, 97, 18, 54, 99, 101, 108, 101, 115, 116, 105, 97, 118, 97, 108, 111, 112, 101, 114, 49, 114, 102, 108, 117, 116, 107, 51, 101, 117, 119, 56, 100, 99, 119, 97, 101, 104, 120, 119, 117, 103, 99, 109, 57, 112, 101, 119, 107, 100, 110, 53, 54, 120, 106, 108, 104, 50, 54, 26, 54, 99, 101, 108, 101, 115, 116, 105, 97, 118, 97, 108, 111, 112, 101, 114, 49, 100, 97, 118, 122, 52, 48, 107, 97, 116, 57, 51, 116, 52, 57, 108, 106, 114, 107, 109, 107, 108, 53, 117, 113, 104, 113, 113, 52, 53, 101, 48, 116, 117, 106, 50, 115, 51, 109, 34, 15, 10, 4, 117, 116, 105, 97, 18, 7, 49, 48, 48, 48, 48, 48, 48, 18, 22, 116, 101, 115, 116, 32, 117, 105, 32, 114, 101, 100, 101, 108, 101, 103, 97, 116, 101, 32, 116, 120, 32, 18, 103, 10, 80, 10, 70, 10, 31, 47, 99, 111, 115, 109, 111, 115, 46, 99, 114, 121, 112, 116, 111, 46, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 2, 205, 82, 66, 173, 172, 164, 110, 151, 162, 183, 151, 111, 80, 96, 191, 38, 188, 141, 208, 175, 86, 52, 254, 146, 134, 204, 43, 40, 79, 127, 106, 1, 18, 4, 10, 2, 8, 127, 24, 39, 18, 19, 10, 13, 10, 4, 117, 116, 105, 97, 18, 5, 55, 50, 52, 51, 49, 16, 185, 215, 17, 26, 64, 98, 225, 18, 145, 187, 225, 213, 198, 229, 6, 6, 240, 177, 0, 28, 112, 160, 126, 193, 177, 221, 161, 96, 79, 5, 192, 224, 168, 253, 161, 12, 33, 9, 118, 215, 22, 219, 239, 73, 133, 79, 37, 218, 83, 238, 115, 44, 232, 16, 163, 242, 174, 100, 175, 162, 213, 142, 194, 58, 69, 84, 81, 3, 70}

func createEmptyBlock() (types.BlockData, time.Time) {
	return createBlock(nodeTypes.ResponseDeliverTx{}, 0)
}

func createBlockWithTxs(tx nodeTypes.ResponseDeliverTx, txData []byte, count int) (types.BlockData, time.Time) {
	now := time.Now()
	headerBlock := nodeTypes.Block{
		Header: nodeTypes.Header{
			Time: now,
		},
		Data: nodeTypes.Data{
			Txs: make(tmTypes.Txs, count),
		},
	}

	var txResults = make([]*nodeTypes.ResponseDeliverTx, count)
	for i := 0; i < count; i++ {
		txResults[i] = &tx
		headerBlock.Data.Txs[i] = txData
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

func createBlock(tx nodeTypes.ResponseDeliverTx, count int) (types.BlockData, time.Time) {
	now := time.Now()
	headerBlock := nodeTypes.Block{
		Header: nodeTypes.Header{
			Time: now,
		},
		Data: nodeTypes.Data{
			Txs: make(tmTypes.Txs, count),
		},
	}

	var txResults = make([]*nodeTypes.ResponseDeliverTx, count)
	for i := 0; i < count; i++ {
		txResults[i] = &tx
		headerBlock.Data.Txs[i] = txMsgBeginRedelegate
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

	resultTxs, err := parseTxs(block)

	assert.NoError(t, err)
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

	resultTxs, err := parseTxs(block)

	assert.NoError(t, err)
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

	resultTxs, err := parseTxs(block)

	assert.NoError(t, err)
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

	resultTxs, err := parseTxs(block)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusFailed, f.Status)
	assert.Equal(t, "something unusual happened", f.Error)
	assert.Equal(t, uint64(12000), f.GasWanted)
	assert.Equal(t, uint64(1000), f.GasUsed)
	assert.Equal(t, "celestia-explorer", f.Codespace)
}
