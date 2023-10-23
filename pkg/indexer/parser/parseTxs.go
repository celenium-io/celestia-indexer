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

	for i, txRes := range b.TxsResults {
		t, err := parseTx(b, i, txRes)
		if err != nil {
			return nil, err
		}

		txs[i] = t
	}

	return txs, nil
}

func parseTx(b types.BlockData, index int, txRes *types.ResponseDeliverTx) (storage.Tx, error) {
	d, err := decode.Tx(b, index)
	if err != nil {
		return storage.Tx{}, errors.Wrapf(err, "while parsing Tx on index %d", index)
	}

	t := storage.Tx{
		Height:        b.Height,
		Time:          b.Block.Time,
		Position:      int64(index),
		GasWanted:     txRes.GasWanted,
		GasUsed:       txRes.GasUsed,
		TimeoutHeight: d.TimeoutHeight,
		EventsCount:   int64(len(txRes.Events)),
		MessagesCount: int64(len(d.Messages)),
		Fee:           d.Fee,
		Status:        storageTypes.StatusSuccess,
		Codespace:     txRes.Codespace,
		Hash:          b.Block.Txs[index].Hash(),
		Memo:          d.Memo,
		MessageTypes:  storageTypes.NewMsgTypeBitMask(),

		Messages:  make([]storage.Message, len(d.Messages)),
		Events:    nil,
		Signers:   make([]storage.Address, 0),
		BlobsSize: 0,
	}

	for signer := range d.Signers {
		_, hash, err := types.Address(signer).Decode()
		if err != nil {
			return t, errors.Wrapf(err, "decode signer: %s", signer)
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
	for position, sdkMsg := range d.Messages {
		dm, err := decode.Message(sdkMsg, b.Height, b.Block.Time, position, t.Status)
		if err != nil {
			return storage.Tx{}, errors.Wrapf(err, "while parsing tx=%v on index=%d", t.Hash, t.Position)
		}

		if txRes.IsFailed() {
			dm.Msg.Namespace = nil
			dm.BlobsSize = 0
		}

		t.Messages[position] = dm.Msg
		t.MessageTypes.SetBit(dm.Msg.Type)
		t.BlobsSize += dm.BlobsSize
	}

	return t, nil
}
