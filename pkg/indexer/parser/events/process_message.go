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

func handleHyperlaneProcessMessage(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/hyperlane.core.v1.MsgProcessMessage" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processHyperlaneProcessMessage(ctx, c, msg)
}

func processHyperlaneProcessMessage(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	var transfer = &storage.HLTransfer{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		TxId:   msg.TxId,
	}

	for event := range c.MsgEvents("action") {
		switch event.Type {
		case types.EventTypeHyperlanecorev1EventProcess:
			processEvent, err := decode.NewHyperlaneProcessEvent(event.Data)
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

			if metadata := msg.Data.GetStringOrDefault("Metadata"); metadata != "" {
				decodedMetadata, err := util.DecodeEthHex(metadata)
				if err != nil {
					return errors.Wrap(err, "decode process message metadata")
				}
				transfer.Metadata = decodedMetadata
			}

			if relayer := msg.Data.GetStringOrDefault("Relayer"); relayer != "" {
				transfer.Relayer = &storage.Address{
					Address:    relayer,
					Height:     msg.Height,
					LastHeight: msg.Height,
					Balances:   []storage.Balance{storage.EmptyBalance()},
				}
			}
			ctx.AddHlTransfer(transfer)
		case types.EventTypeHyperlanewarpv1EventReceiveRemoteTransfer:
			receiveEvent, err := decode.NewHyperlaneReceiveTransferEvent(event.Data)
			if err != nil {
				return errors.Wrap(err, "parse hyperlane receive transfer event")
			}

			if err := makeHyperlaneTransferAddress(ctx, receiveEvent.Sender, transfer, msg.Height); err != nil {
				return errors.Wrap(err, "makeHyperlaneTransferAddress")
			}
			if err := makeHyperlaneTransferAddress(ctx, receiveEvent.Recipient, transfer, msg.Height); err != nil {
				return errors.Wrap(err, "makeHyperlaneTransferAddress")
			}

			transfer.Denom = receiveEvent.Denom
			transfer.Amount = receiveEvent.Amount
			tokenId, err := util.DecodeHexAddress(receiveEvent.TokenId)
			if err != nil {
				return errors.Wrap(err, "decode token id")
			}
			transfer.Token = &storage.HLToken{
				TokenId:          tokenId.Bytes(),
				ReceiveTransfers: 1,
				Received:         receiveEvent.Amount,
				Type:             types.HLTokenTypeCollateral,
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

func makeHyperlaneTransferAddress(ctx *context.Context, str string, transfer *storage.HLTransfer, height pkgTypes.Level) error {
	if prefix, hash, err := pkgTypes.Address(str).Decode(); err == nil && prefix == pkgTypes.AddressPrefixCelestia {
		transfer.Address = &storage.Address{
			Address:    str,
			Hash:       hash,
			Height:     height,
			LastHeight: height,
		}
		return ctx.AddAddress(transfer.Address)
	}
	str = strings.TrimPrefix(str, "0x")
	str = strings.TrimLeft(str, "0")
	transfer.CounterpartyAddress = str
	return nil
}
