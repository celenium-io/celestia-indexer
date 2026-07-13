// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleChannelClose(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	action := decoder.StringFromMap(event.Data, "action")
	isValidMsg := action == "/ibc.core.channel.v1.MsgChannelCloseConfirm" || action == "/ibc.core.channel.v1.MsgChannelCloseInit"

	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processChannelClose(ctx, c, msg)
}

func processChannelClose(ctx *context.Context, c *Cursor, _ *storage.Message) error {
	event, ok := c.Peek()
	if !ok {
		return errors.New("not enough events for channel close")
	}
	if event.Type != storageTypes.EventTypeChannelCloseConfirm {
		return errors.Errorf("invalid event type: %s", event.Type)
	}
	cc := decode.NewChannelChange(event.Data)

	ibcChannel := &storage.IbcChannel{
		Id:     cc.ChannelId,
		Status: storageTypes.IbcChannelStatusClosed,
	}
	ctx.AddIbcChannel(ibcChannel)

	return nil
}
