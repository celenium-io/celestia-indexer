// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleHyperlaneRemoteTransfer(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/hyperlane.warp.v1.MsgRemoteTransfer" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processHyperlaneRemoteTransfer(ctx, c, msg)
}

func processHyperlaneRemoteTransfer(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	var transfer = &storage.HLTransfer{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		TxId:   msg.TxId,
	}

	for event := range c.MsgEvents("action") {
		switch event.Type {
		case types.EventTypeHyperlanecorev1EventDispatch:
			dispatchEvent, err := decode.NewHyperlaneDispatchEvent(event.Data)
			if err != nil {
				return errors.Wrap(err, "parse hyperlane dispatch event")
			}

			originMailboxId, err := util.DecodeHexAddress(dispatchEvent.OriginMailboxId)
			if err != nil {
				return errors.Wrap(err, "decode mailbox id")
			}

			transfer.Counterparty = dispatchEvent.Destination
			transfer.Version = dispatchEvent.Message.Version
			transfer.Nonce = dispatchEvent.Message.Nonce
			transfer.Body = dispatchEvent.Message.Body
			transfer.Type = types.HLTransferTypeSend

			transfer.Mailbox = &storage.HLMailbox{
				Mailbox:      originMailboxId.Bytes(),
				InternalId:   originMailboxId.GetInternalId(),
				SentMessages: 1,
			}

			ctx.AddHlTransfer(transfer)
		case types.EventTypeHyperlanewarpv1EventSendRemoteTransfer:
			sendEvent, err := decode.NewHyperlaneSendTransferEvent(event.Data)
			if err != nil {
				return errors.Wrap(err, "parse hyperlane send transfer event")
			}

			if err := makeHyperlaneTransferAddress(ctx, sendEvent.Sender, transfer, msg.Height); err != nil {
				return errors.Wrap(err, "makeHyperlaneTransferAddress")
			}
			if err := makeHyperlaneTransferAddress(ctx, sendEvent.Recipient, transfer, msg.Height); err != nil {
				return errors.Wrap(err, "makeHyperlaneTransferAddress")
			}

			transfer.Denom = sendEvent.Denom
			transfer.Amount = sendEvent.Amount
			tokenId, err := util.DecodeHexAddress(sendEvent.TokenId)
			if err != nil {
				return errors.Wrap(err, "decode token id")
			}
			transfer.Token = &storage.HLToken{
				TokenId:       tokenId.Bytes(),
				SentTransfers: 1,
				Sent:          sendEvent.Amount,
				Type:          types.HLTokenTypeCollateral,
			}
		case types.EventTypeHyperlanecorepostDispatchv1EventGasPayment:
			gasEvent, err := decode.NewHyperlaneGasPaymentEvent(event.Data)
			if err != nil {
				return errors.Wrap(err, "parse hyperlane gas payment event")
			}

			igpId, err := util.DecodeHexAddress(gasEvent.IgpId)
			if err != nil {
				return errors.Wrap(err, "decode igp id")
			}
			transfer.GasPayment = &storage.HLGasPayment{
				Height:    ctx.Block.Height,
				Time:      ctx.Block.Time,
				Amount:    gasEvent.Amount,
				GasAmount: gasEvent.GasAmount,
				Igp: &storage.HLIGP{
					IgpId: igpId.Bytes(),
				},
			}
		}
	}

	return nil
}
