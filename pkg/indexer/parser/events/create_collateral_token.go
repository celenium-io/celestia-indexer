// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
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
	"github.com/shopspring/decimal"
)

func handleCreateCollateralToken(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/hyperlane.warp.v1.MsgCreateCollateralToken" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processCreateCollateralToken(ctx, events, msg, idx)
}

func processCreateCollateralToken(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false
	for !end {
		if events[*idx].Type == types.EventTypeHyperlanewarpv1EventCreateCollateralToken {
			createToken, err := decode.NewCreateCollateralToken(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parsing create collateral token event")
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
				Type:   types.HLTokenTypeCollateral,
				Mailbox: &storage.HLMailbox{
					Height:     ctx.Block.Height,
					Time:       ctx.Block.Time,
					Mailbox:    originMailboxId.Bytes(),
					InternalId: originMailboxId.GetInternalId(),
				},
				TokenId:  tokenId.Bytes(),
				Received: decimal.Zero,
				Sent:     decimal.Zero,
			}

			if createToken.Owner != "" {
				token.Owner = &storage.Address{
					Address: createToken.Owner,
				}
			}

			msg.HLToken = token
			end = true
		}

		action := decoder.StringFromMap(events[*idx].Data, "action")
		end = len(events)-1 == *idx || action != "" || end
		*idx += 1
	}

	toTheNextAction(events, idx)
	return nil
}
