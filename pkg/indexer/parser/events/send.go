// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleSend(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	action := decoder.StringFromMap(event.Data, "action")
	isValidMsg := action == "/cosmos.bank.v1beta1.MsgSend"
	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processSend(ctx, c, msg)
}

func processSend(_ *context.Context, c *Cursor, _ *storage.Message) error {
	event, ok := c.Peek()
	if !ok {
		return errors.New("not enough events for send")
	}
	msgIdx := decoder.StringFromMap(event.Data, "msg_index")
	newFormat := msgIdx != ""

	if newFormat {
		c.Skip(4)
	} else {
		c.Skip(5)
	}

	return nil
}
