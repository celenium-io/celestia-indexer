// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func parseTxs(b types.BlockData) ([]storage.Tx, error) {
	txs := make([]storage.Tx, len(b.TxsResults))

	for i := range b.TxsResults {
		if err := parseTx(b, i, b.TxsResults[i], &txs[i]); err != nil {
			return nil, err
		}
	}

	return txs, nil
}

func parseTx(b types.BlockData, index int, txRes *types.ResponseDeliverTx, t *storage.Tx) error {
	d, err := decode.Tx(b, index)
	if err != nil {
		return errors.Wrapf(err, "while parsing Tx on index %d", index)
	}

	t.Height = b.Height
	t.Time = b.Block.Time
	t.Position = int64(index)
	t.GasWanted = txRes.GasWanted
	t.GasUsed = txRes.GasUsed
	t.TimeoutHeight = d.TimeoutHeight
	t.EventsCount = int64(len(txRes.Events))
	t.MessagesCount = int64(len(d.Messages))
	t.Fee = d.Fee
	t.Status = storageTypes.StatusSuccess
	t.Codespace = txRes.Codespace
	t.Hash = b.Block.Txs[index].Hash()
	t.Memo = d.Memo
	t.MessageTypes = storageTypes.NewMsgTypeBitMask()
	t.Messages = make([]storage.Message, len(d.Messages))
	t.Events = nil
	t.Signers = make([]storage.Address, 0)
	t.BlobsSize = 0
	t.BytesSize = int64(len(txRes.Data))

	for signer := range d.Signers {
		_, hash, err := types.Address(signer).Decode()
		if err != nil {
			return errors.Wrapf(err, "decode signer: %s", signer)
		}

		t.Signers = append(t.Signers, storage.Address{
			Address:    signer,
			Height:     t.Height,
			LastHeight: t.Height,
			Hash:       hash,
			Balance: storage.Balance{
				Total: decimal.Zero,
			},
		})
	}

	if txRes.IsFailed() {
		t.Status = storageTypes.StatusFailed
		t.Error = txRes.Log
	}

	t.Events = parseEvents(b, txRes.Events)
	for i := range d.Messages {
		dm, err := decode.Message(d.Messages[i], b.Height, b.Block.Time, i, t.Status)
		if err != nil {
			return errors.Wrapf(err, "while parsing tx=%v on index=%d", t.Hash, t.Position)
		}

		if txRes.IsFailed() {
			dm.Msg.Namespace = nil
			dm.BlobsSize = 0
		}

		t.Messages[i] = dm.Msg
		t.MessageTypes.SetByMsgType(dm.Msg.Type)
		t.BlobsSize += dm.BlobsSize
	}

	return nil
}
