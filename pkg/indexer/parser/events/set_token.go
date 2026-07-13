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

func handleSetToken(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/hyperlane.warp.v1.MsgSetToken" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processSetToken(ctx, c, msg)
}

func processSetToken(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	for event := range c.MsgEvents("action") {
		if event.Type != types.EventTypeHyperlanewarpv1EventSetToken {
			continue
		}

		setToken, err := decode.NewSetToken(event.Data)
		if err != nil {
			return err
		}

		if setToken.NewOwner == "" {
			continue
		}

		tokenId, err := util.DecodeHexAddress(setToken.TokenId)
		if err != nil {
			return errors.Wrap(err, "decode token id")
		}

		token := &storage.HLToken{
			TokenId: tokenId.Bytes(),
			Owner: &storage.Address{
				Address:    setToken.NewOwner,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balances:   []storage.Balance{storage.EmptyBalance()},
			},
			Type: types.HLTokenTypeCollateral,
		}
		ctx.AddHlToken(token)
		if err := ctx.AddAddress(token.Owner); err != nil {
			return errors.Wrap(err, "add address")
		}
		break
	}

	return nil
}
