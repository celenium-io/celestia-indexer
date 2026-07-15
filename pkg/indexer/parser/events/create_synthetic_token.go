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

func handleCreateSyntheticToken(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/hyperlane.warp.v1.MsgCreateSyntheticToken" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processCreateSyntheticToken(ctx, c, msg)
}

func processCreateSyntheticToken(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	for event := range c.MsgEvents("action") {
		if event.Type != types.EventTypeHyperlanewarpv1EventCreateSyntheticToken {
			continue
		}

		createToken, err := decode.NewCreateSyntheticToken(event.Data)
		if err != nil {
			return errors.Wrap(err, "parsing create synthetic token event")
		}

		originMailboxId, err := util.DecodeHexAddress(createToken.MailboxId)
		if err != nil {
			return errors.Wrap(err, "decode mailbox id")
		}

		tokenId, err := util.DecodeHexAddress(createToken.TokenId)
		if err != nil {
			return errors.Wrap(err, "decode token id")
		}

		token := &storage.HLToken{
			Height: ctx.Block.Height,
			Time:   ctx.Block.Time,
			Denom:  createToken.Denom,
			Type:   types.HLTokenTypeSynthetic,
			Mailbox: &storage.HLMailbox{
				Height:     ctx.Block.Height,
				Time:       ctx.Block.Time,
				Mailbox:    originMailboxId.Bytes(),
				InternalId: originMailboxId.GetInternalId(),
			},
			TokenId:  tokenId.Bytes(),
			Sent:     types.NumericZero(),
			Received: types.NumericZero(),
			TxId:     msg.TxId,
		}

		if createToken.Owner != "" {
			token.Owner = &storage.Address{
				Address:    createToken.Owner,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balances:   []storage.Balance{storage.EmptyBalance()},
			}
			if err := ctx.AddAddress(token.Owner); err != nil {
				return errors.Wrap(err, "add address")
			}
		}

		ctx.AddHlToken(token)
		break
	}
	return nil
}
