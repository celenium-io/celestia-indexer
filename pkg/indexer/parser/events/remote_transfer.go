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

func handleHyperlaneRemoteTransfer(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/hyperlane.warp.v1.MsgRemoteTransfer" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processHyperlaneRemoteTransfer(ctx, events, msg, idx)
}

func processHyperlaneRemoteTransfer(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false

	var transfer = &storage.HLTransfer{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
	}

	for !end {

		switch events[*idx].Type {
		case types.EventTypeHyperlanecorev1EventDispatch:
			dispatchEvent, err := decode.NewHyperlaneDispatchEvent(events[*idx].Data)
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

			msg.HLTransfer = transfer
			end = true
		case types.EventTypeHyperlanewarpv1EventSendRemoteTransfer:
			event, err := decode.NewHyperlaneSendTransferEvent(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parse hyperlane send transfer event")
			}

			makeHyperlaneTransferAddress(event.Sender, transfer, msg.Height)
			makeHyperlaneTransferAddress(event.Recipient, transfer, msg.Height)

			transfer.Denom = event.Denom
			transfer.Amount = event.Amount
			tokenId, err := util.DecodeHexAddress(event.TokenId)
			if err != nil {
				return errors.Wrap(err, "decode token id")
			}
			transfer.Token = &storage.HLToken{
				TokenId:       tokenId.Bytes(),
				SentTransfers: 1,
				Sent:          event.Amount,
				Type:          types.HLTokenTypeCollateral,
			}
		case types.EventTypeHyperlanecorepostDispatchv1EventGasPayment:
			event, err := decode.NewHyperlaneGasPaymentEvent(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parse hyperlane gas payment event")
			}

			igpId, err := util.DecodeHexAddress(event.IgpId)
			if err != nil {
				return errors.Wrap(err, "decode igp id")
			}
			transfer.GasPayment = &storage.HLGasPayment{
				Height:    ctx.Block.Height,
				Time:      ctx.Block.Time,
				Amount:    event.Amount,
				GasAmount: event.GasAmount,
				Igp: &storage.HLIGP{
					IgpId: igpId.Bytes(),
				},
			}
		}

		action := decoder.StringFromMap(events[*idx].Data, "action")
		end = len(events)-1 == *idx || action != "" || end
		*idx += 1
	}

	toTheNextAction(events, idx)
	return nil
}
