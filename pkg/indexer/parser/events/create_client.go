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

func handleCreateClient(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/ibc.core.client.v1.MsgCreateClient" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processCreateClient(ctx, events, msg, idx)
}

func processCreateClient(_ *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	cc, err := decode.NewUpdateClient(events[*idx].Data)
	if err != nil {
		return errors.Wrap(err, "parsing CreateClient event")
	}

	state, err := decoder.ClientStateFromMap(msg.Data, "ClientState")
	if err != nil {
		return errors.Wrap(err, "receiving ClientState from message")
	}

	signer := decoder.StringFromMap(msg.Data, "Signer")

	msg.IbcClient = &storage.IbcClient{
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
	}
	*idx += 2
	return nil
}
