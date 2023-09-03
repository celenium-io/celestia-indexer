package parser

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

func parseTxs(b types.BlockData) []storage.Tx {
	txs := make([]storage.Tx, len(b.TxsResults))

	for i, txRes := range b.TxsResults {
		tx := parseTx(b, i, txRes)
		txs[i] = tx
	}

	return txs
}

func parseTx(b types.BlockData, index int, txRes *nodeTypes.ResponseDeliverTx) storage.Tx {
	tx := storage.Tx{
		Height:        b.Height,
		Time:          b.Block.Time,
		Position:      uint64(index),
		GasWanted:     uint64(txRes.GasWanted),
		GasUsed:       uint64(txRes.GasUsed),
		TimeoutHeight: 0, // TODO like nodes.guru
		EventsCount:   uint64(len(txRes.Events)),
		MessagesCount: 0,            // TODO
		Fee:           decimal.Zero, // TODO like nodes.guru
		Status:        storageTypes.StatusSuccess,
		Codespace:     txRes.Codespace,
		Hash:          make([]byte, 0),            // TODO like nodes.guru
		Memo:          "",                         // TODO like nodes.guru
		MessageTypes:  storageTypes.MsgTypeBits{}, // TODO

		Messages: nil, // make([]storage.Message, 0), // TODO
		Events:   nil,
	}

	if txRes.Code != 0 {
		tx.Status = storageTypes.StatusFailed
		tx.Error = txRes.Log
	}

	tx.Events = parseEvents(b, txRes.Events)

	return tx
}
