package parser

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func parseTxs(b types.BlockData) ([]storage.Tx, error) {
	txs := make([]storage.Tx, len(b.TxsResults))

	for i, txRes := range b.TxsResults {
		t, err := parseTx(b, i, txRes)
		if err != nil {
			return nil, err
		}

		txs[i] = t
	}

	return txs, nil
}

func parseTx(b types.BlockData, index int, txRes *nodeTypes.ResponseDeliverTx) (storage.Tx, error) {
	d, err := decodeTx(b, index)
	if err != nil {
		return storage.Tx{}, errors.Wrapf(err, "while parsing Tx on index %d", index)
	}

	t := storage.Tx{
		Height:        b.Height,
		Time:          b.Block.Time,
		Position:      uint64(index),
		GasWanted:     uint64(txRes.GasWanted),
		GasUsed:       uint64(txRes.GasUsed),
		TimeoutHeight: d.timeoutHeight,
		EventsCount:   uint64(len(txRes.Events)),
		MessagesCount: uint64(len(d.messages)),
		Fee:           d.fee,
		Status:        storageTypes.StatusSuccess,
		Codespace:     txRes.Codespace,
		Hash:          b.Block.Txs[index].Hash(),
		Memo:          d.memo,
		MessageTypes:  storageTypes.NewMsgTypeBitMask(),

		Messages:  make([]storage.Message, len(d.messages)),
		Events:    nil,
		Addresses: make([]storage.AddressWithType, 0),
		BlobsSize: 0,
	}

	if txRes.Code != 0 {
		t.Status = storageTypes.StatusFailed
		t.Error = txRes.Log
	}

	t.Events = parseEvents(b, txRes.Events)
	for position, sdkMsg := range d.messages {
		dm, err := decodeMsg(b, sdkMsg, position)
		if err != nil {
			return storage.Tx{}, errors.Wrapf(err, "while parsing tx=%v on index=%d", t.Hash, t.Position)
		}

		t.Messages[position] = dm.msg
		t.MessageTypes.SetBit(dm.msg.Type)
		t.BlobsSize += dm.blobsSize
		t.Addresses = append(t.Addresses, dm.addresses...)
	}

	return t, nil
}
