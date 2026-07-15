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

func handleConnectionOpenInit(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	action := decoder.StringFromMap(event.Data, "action")
	isValidMsg := action == "/ibc.core.connection.v1.MsgConnectionOpenInit" || action == "/ibc.core.connection.v1.MsgConnectionOpenTry"

	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processConnectionOpenInit(ctx, c, msg)
}

func processConnectionOpenInit(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	for _, event := range c.Remaining() {
		if event.Type != storageTypes.EventTypeConnectionOpenInit && event.Type != storageTypes.EventTypeConnectionOpenTry {
			continue
		}
		cc := decode.NewConnectionOpen(event.Data)

		conn := &storage.IbcConnection{
			Height:                   msg.Height,
			CreatedAt:                msg.Time,
			ClientId:                 cc.ClientId,
			ConnectionId:             cc.ConnectionId,
			CounterpartyClientId:     cc.CounterpartyClientId,
			CounterpartyConnectionId: cc.CounterpartyConnectionId,
			ChannelsCount:            0,
			CreateTxId:               msg.TxId,
		}
		ctx.AddIbcConnection(conn)
		break
	}

	c.Skip(2)
	return nil
}
