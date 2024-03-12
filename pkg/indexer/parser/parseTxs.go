// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/parser/events"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func parseTxs(ctx *context.Context, b types.BlockData) ([]storage.Tx, error) {
	txs := make([]storage.Tx, len(b.TxsResults))

	for i := range b.TxsResults {
		if err := parseTx(ctx, b, i, b.TxsResults[i], &txs[i]); err != nil {
			return nil, err
		}
	}

	return txs, nil
}

func parseTx(ctx *context.Context, b types.BlockData, index int, txRes *types.ResponseDeliverTx, t *storage.Tx) error {
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

	for signer, signerBytes := range d.Signers {
		address := storage.Address{
			Address:    signer.String(),
			Height:     t.Height,
			LastHeight: t.Height,
			Hash:       signerBytes,
			Balance: storage.Balance{
				Currency:  currency.DefaultCurrency,
				Spendable: decimal.Zero,
				Delegated: decimal.Zero,
				Unbonding: decimal.Zero,
			},
		}
		t.Signers = append(t.Signers, address)
		if err := ctx.AddAddress(&address); err != nil {
			return err
		}
	}

	if txRes.IsFailed() {
		t.Status = storageTypes.StatusFailed
		t.Error = txRes.Log
	}

	t.Events, err = parseEvents(ctx, b, txRes.Events)
	if err != nil {
		return errors.Wrap(err, "parsing events")
	}

	var eventsIdx int

	// find first action
	for i := range t.Events {
		if t.Events[i].Type != storageTypes.EventTypeMessage {
			continue
		}
		if action := decoder.StringFromMap(t.Events[i].Data, "action"); action != "" {
			eventsIdx = i
			break
		}
	}

	for i := range d.Messages {
		dm, err := decode.Message(ctx, d.Messages[i], i, t.Status)
		if err != nil {
			return errors.Wrapf(err, "while parsing tx=%v on index=%d", t.Hash, t.Position)
		}

		processBlob(dm.Msg.BlobLogs, d, t)

		if txRes.IsFailed() {
			dm.Msg.Namespace = nil
			dm.BlobsSize = 0
		}

		t.Messages[i] = dm.Msg
		t.MessageTypes.SetByMsgType(dm.Msg.Type)
		t.BlobsSize += dm.BlobsSize

		if !txRes.IsFailed() {
			if err := events.Handle(ctx, t.Events, &t.Messages[i], &eventsIdx); err != nil {
				return err
			}
		}
	}

	ctx.Block.Stats.Fee = ctx.Block.Stats.Fee.Add(t.Fee)
	ctx.Block.MessageTypes.Set(t.MessageTypes.Bits)
	ctx.Block.Stats.BlobsSize += t.BlobsSize
	ctx.Block.Stats.GasLimit += t.GasWanted
	ctx.Block.Stats.GasUsed += t.GasUsed
	ctx.Block.Stats.BlobsCount += t.BlobsCount

	return nil
}
