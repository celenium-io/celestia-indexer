// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleCreateClient(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/ibc.core.client.v1.MsgCreateClient" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processCreateClient(ctx, c, msg)
}

func processCreateClient(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	event, ok := c.Peek()
	if !ok {
		return errors.New("not enough events for create client")
	}
	cc, err := decode.NewUpdateClient(event.Data)
	if err != nil {
		return errors.Wrap(err, "parsing CreateClient event")
	}

	state, err := decoder.ClientStateFromMap(msg.Data, "ClientState")
	if err != nil {
		return errors.Wrap(err, "receiving ClientState from message")
	}

	signer := msg.Data.GetStringOrDefault("Signer")

	ibcClient := &storage.IbcClient{
		Height:                msg.Height,
		Type:                  cc.Type,
		CreatedAt:             msg.Time,
		UpdatedAt:             msg.Time,
		Id:                    cc.Id,
		TrustingPeriod:        state.TrustingPeriod,
		UnbondingPeriod:       state.UnbondingPeriod,
		MaxClockDrift:         state.MaxClockDrift,
		LatestRevisionHeight:  state.LatestHeight.RevisionHeight,
		LatestRevisionNumber:  state.LatestHeight.RevisionNumber,
		FrozenRevisionHeight:  state.FrozenHeight.RevisionHeight,
		FrozenRevisionNumber:  state.FrozenHeight.RevisionNumber,
		TrustLevelDenominator: state.TrustLevel.Denominator,
		TrustLevelNumerator:   state.TrustLevel.Numerator,
		ConnectionCount:       0,
		Creator: &storage.Address{
			Address: signer,
		},
		TxId: msg.TxId,
	}
	ctx.AddIbcClient(ibcClient)
	c.Skip(2)
	return nil
}
