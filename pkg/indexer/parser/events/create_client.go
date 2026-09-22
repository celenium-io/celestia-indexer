// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
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

	clientStateData, ok := msg.Data["ClientState"]
	if !ok {
		return errors.Errorf("can't find 'ClientState' field in message data")
	}

	switch state := clientStateData.(type) {
	case tendermint.ClientState:
		ibcClient := newIbcClientForTendermint(state, msg, cc)
		ctx.AddIbcClient(ibcClient)
	case solomachine.ClientState:
		ibcClient := newIbcClientForSolomachine(state, msg, cc)
		ctx.AddIbcClient(ibcClient)
	}

	c.Skip(2)
	return nil
}

func newIbcClientForTendermint(state tendermint.ClientState, msg *storage.Message, cc decode.UpdateClient) *storage.IbcClient {
	ibcClient := newIbcClientForDefault(msg, cc)
	ibcClient.TrustingPeriod = state.TrustingPeriod
	ibcClient.UnbondingPeriod = state.UnbondingPeriod
	ibcClient.MaxClockDrift = state.MaxClockDrift
	ibcClient.LatestRevisionHeight = state.LatestHeight.RevisionHeight
	ibcClient.LatestRevisionNumber = state.LatestHeight.RevisionNumber
	ibcClient.FrozenRevisionHeight = state.FrozenHeight.RevisionHeight
	ibcClient.FrozenRevisionNumber = state.FrozenHeight.RevisionNumber
	ibcClient.TrustLevelDenominator = state.TrustLevel.Denominator
	ibcClient.TrustLevelNumerator = state.TrustLevel.Numerator
	ibcClient.ChainId = state.ChainId
	return ibcClient
}

func newIbcClientForSolomachine(state solomachine.ClientState, msg *storage.Message, cc decode.UpdateClient) *storage.IbcClient {
	ibcClient := newIbcClientForDefault(msg, cc)
	ibcClient.LatestRevisionHeight = state.Sequence
	return ibcClient
}

func newIbcClientForDefault(msg *storage.Message, cc decode.UpdateClient) *storage.IbcClient {
	signer := msg.Data.GetStringOrDefault("Signer")
	return &storage.IbcClient{
		Height:          msg.Height,
		Type:            cc.Type,
		CreatedAt:       msg.Time,
		UpdatedAt:       msg.Time,
		Id:              cc.Id,
		ConnectionCount: 0,
		Creator: &storage.Address{
			Address: signer,
		},
		TxId: msg.TxId,
	}
}
