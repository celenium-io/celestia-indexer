// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/pkg/errors"
)

func handleUpdateClient(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/ibc.core.client.v1.MsgUpdateClient" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processUpdateClient(ctx, c, msg)
}

func processUpdateClient(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	event, ok := c.Peek()
	if !ok {
		return errors.New("not enough events for update client")
	}

	switch event.Type {
	case types.EventTypeClientMisbehaviour:
		if err := handleMisbehaviour(ctx, event, msg); err != nil {
			return errors.Wrap(err, "handle Misbehaviour")
		}
	case types.EventTypeUpdateClient:
		uc, err := decode.NewUpdateClient(event.Data)
		if err != nil {
			return errors.Wrap(err, "parse update client event")
		}
		switch header := msg.Data["Header"].(type) {
		case tendermint.Header:
			ctx.AddIbcClient(&storage.IbcClient{
				Id:                   uc.Id,
				UpdatedAt:            msg.Time,
				ChainId:              header.Header.ChainID,
				LatestRevisionHeight: uc.ConsensusHeight,
				LatestRevisionNumber: uc.Revision,
			})
		case solomachine.Header:
			ctx.AddIbcClient(&storage.IbcClient{
				Id:                   uc.Id,
				UpdatedAt:            msg.Time,
				LatestRevisionHeight: uc.ConsensusHeight,
				LatestRevisionNumber: uc.Revision,
			})
		default:
			// invalid Misbehaviour: update_client without heights, client is not frozen
			ctx.AddIbcClient(&storage.IbcClient{
				Id:        uc.Id,
				UpdatedAt: msg.Time,
			})
		}
	default:
		return errors.Errorf("unexpected event %s for MsgUpdateClient", event.Type)
	}

	c.Skip(2)

	return nil
}

func handleMisbehaviour(ctx *context.Context, event storage.Event, msg *storage.Message) error {
	uc, err := decode.NewClientMisbehaviour(event.Data)
	if err != nil {
		return errors.Wrap(err, "parse update client event")
	}
	ctx.AddIbcClient(&storage.IbcClient{
		Id:                   uc.Id,
		UpdatedAt:            msg.Time,
		FrozenRevisionHeight: 1,
	})
	return nil
}
