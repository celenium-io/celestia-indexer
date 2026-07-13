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
	uc, err := decode.NewUpdateClient(event.Data)
	if err != nil {
		return errors.Wrap(err, "parse update client event")
	}

	header, err := decoder.HeaderFromMap(msg.Data, "Header")
	if err != nil {
		return errors.Wrap(err, "receiving Header from message")
	}

	ibcClient := &storage.IbcClient{
		Id:                   uc.Id,
		UpdatedAt:            msg.Time,
		ChainId:              header.Header.ChainID,
		LatestRevisionHeight: uc.ConsensusHeight,
		LatestRevisionNumber: uc.Revision,
	}
	ctx.AddIbcClient(ibcClient)

	c.Skip(2)
	return nil
}
