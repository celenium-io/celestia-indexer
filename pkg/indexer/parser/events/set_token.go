// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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

func handleSetToken(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/hyperlane.warp.v1.MsgSetToken" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processSetToken(ctx, events, msg, idx)
}

func processSetToken(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false
	for !end {
		if events[*idx].Type == types.EventTypeHyperlanewarpv1EventSetToken {
			setToken, err := decode.NewSetToken(events[*idx].Data)
			if err != nil {
				return err
			}

			if setToken.NewOwner == "" {
				return nil
			}

			tokenId, err := util.DecodeHexAddress(setToken.TokenId)
			if err != nil {
				return errors.Wrap(err, "decode token id")
			}

			msg.HLToken = &storage.HLToken{
				TokenId: tokenId.Bytes(),
				Owner: &storage.Address{
					Address: setToken.NewOwner,
				},
			}
			end = true
		}

		action := decoder.StringFromMap(events[*idx].Data, "action")
		end = len(events)-1 == *idx || action != "" || end
		*idx += 1
	}

	toTheNextAction(events, idx)
	return nil
}
