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

func handleForward(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/celestia.forwarding.v1.MsgForward" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processForward(ctx, c, msg)
}

func processForward(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	var forwarding = storage.Forwarding{
		Height: msg.Height,
		Time:   msg.Time,
		TxId:   msg.TxId,
	}

	for {
		event, ok := c.Peek()
		if !ok {
			break
		}

		switch event.Type {
		case types.EventTypeCelestiaforwardingv1EventTokenForwarded:
			forwarded, err := decode.NewEventTokenForwarded(event.Data)
			if err != nil {
				return errors.Wrap(err, "decoding token forwarded event")
			}

			// Pre-v8 format lacks token_id — drain remaining events for this
			// message without creating a forwarding entity.
			if forwarded.TokenId == "" {
				c.Next()
				c.SkipToNext("action")
				return nil
			}

			forwarding.Amount, err = types.NumericFromString(forwarded.Amount)
			if err != nil {
				return errors.Wrap(err, "parsing amount as numeric")
			}
			forwarding.Denom = forwarded.Denom
			forwarding.MessageId = forwarded.MessageId

			tokenId, err := util.DecodeHexAddress(forwarded.TokenId)
			if err != nil {
				return errors.Wrap(err, "decode token id")
			}

			forwarding.Token = &storage.HLToken{
				TokenId: tokenId.Bytes(),
			}
			forwarding.Address = &storage.Address{
				Address:      forwarded.ForwardAddress,
				IsForwarding: true,
				Height:       msg.Height,
				LastHeight:   msg.Height,
				Balances:     []storage.Balance{storage.EmptyBalance()},
			}
			if err = ctx.AddAddress(forwarding.Address); err != nil {
				return errors.Wrap(err, "add forwarding address")
			}
			c.Next()

		case types.EventTypeHyperlanewarpv1EventSendRemoteTransfer:
			transferEvent, err := decode.NewHyperlaneSendTransferEvent(event.Data)
			if err != nil {
				return errors.Wrap(err, "parse hyperlane send transfer event")
			}

			recipient, err := util.DecodeHexAddress(transferEvent.Recipient)
			if err != nil {
				return errors.Wrap(err, "decode recipient address")
			}
			forwarding.DestDomain = transferEvent.DestinationDomain
			forwarding.DestRecipient = recipient.Bytes()
			c.Next()

		default:
			if action := decoder.StringFromMap(event.Data, "action"); action != "" {
				ctx.AddForwarding(&forwarding)
				return nil
			}
			c.Next()
		}
	}

	if forwarding.Token == nil {
		// if token is absent in events, we can't process forwarding
		return nil
	}

	ctx.AddForwarding(&forwarding)
	return nil
}
