// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func handleHyperlaneProcessMessage(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/hyperlane.core.v1.MsgProcessMessage" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processHyperlaneProcessMessage(ctx, events, msg, idx)
}

func processHyperlaneProcessMessage(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false

	var transfer = &storage.HLTransfer{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
	}

	for !end {
		switch events[*idx].Type {
		case types.EventTypeHyperlanecorev1EventProcess:
			processEvent, err := decode.NewHyperlaneProcessEvent(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parse hyperlane process event")
			}

			originMailboxId, err := util.DecodeHexAddress(processEvent.OriginMailboxId)
			if err != nil {
				return errors.Wrap(err, "decode mailbox id")
			}

			transfer.Mailbox = &storage.HLMailbox{
				Mailbox:          originMailboxId.Bytes(),
				InternalId:       originMailboxId.GetInternalId(),
				ReceivedMessages: 1,
			}

			transfer.Counterparty = processEvent.Origin
			transfer.Version = processEvent.Message.Version
			transfer.Nonce = processEvent.Message.Nonce
			transfer.Body = processEvent.Message.Body
			transfer.Type = types.HLTransferTypeReceive

			if metadata := decoder.StringFromMap(msg.Data, "Metadata"); metadata != "" {
				decodedMetadata, err := util.DecodeEthHex(metadata)
				if err != nil {
					return errors.Wrap(err, "decode process message metadata")
				}
				transfer.Metadata = decodedMetadata
			}

			if relayer := decoder.StringFromMap(msg.Data, "Relayer"); relayer != "" {
				transfer.Relayer = &storage.Address{
					Address: relayer,
				}
			}
			msg.HLTransfer = transfer
		case types.EventTypeHyperlanewarpv1EventReceiveRemoteTransfer:
			event, err := decode.NewHyperlaneReceiveTransferEvent(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parse hyperlane receive transfer event")
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
				TokenId:          tokenId.Bytes(),
				ReceiveTransfers: 1,
				Received:         event.Amount,
				Type:             types.HLTokenTypeCollateral,
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

	return nil
}

func makeHyperlaneTransferAddress(str string, transfer *storage.HLTransfer, height pkgTypes.Level) {
	if prefix, hash, err := pkgTypes.Address(str).Decode(); err == nil && prefix == pkgTypes.AddressPrefixCelestia {
		transfer.Address = &storage.Address{
			Address:    str,
			Hash:       hash,
			Height:     height,
			LastHeight: height,
		}
	} else {
		str = strings.TrimPrefix(str, "0x")
		str = strings.TrimLeft(str, "0")
		transfer.CounterpartyAddress = str
	}
}
