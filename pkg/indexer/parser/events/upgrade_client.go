// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/pkg/errors"
)

func handleUpgradeClient(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/ibc.core.client.v1.MsgUpgradeClient" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processUpgradeClient(ctx, c, msg)
}

func processUpgradeClient(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	event, ok := c.Peek()
	if !ok {
		return errors.New("not enough events for upgrade client")
	}
	if event.Type != types.EventTypeUpgradeClient {
		return errors.Errorf("unexpected event %s for MsgUpgradeClient", event.Type)
	}

	// consensus_height is the upgraded client's latest height
	uc, err := decode.NewUpdateClient(event.Data)
	if err != nil {
		return errors.Wrap(err, "parse upgrade client event")
	}

	client := &storage.IbcClient{
		Id:                   uc.Id,
		UpdatedAt:            msg.Time,
		LatestRevisionHeight: uc.ConsensusHeight,
		LatestRevisionNumber: uc.Revision,
	}
	// ibc-go takes chain id and unbonding period from the upgraded state, the rest is kept
	if state, ok := msg.Data["ClientState"].(tendermint.ClientState); ok {
		client.ChainId = state.ChainId
		client.UnbondingPeriod = state.UnbondingPeriod
	}
	ctx.AddIbcClient(client)

	c.Skip(2)
	return nil
}
