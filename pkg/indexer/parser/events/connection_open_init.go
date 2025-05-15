// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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

func handleConnectionOpenInit(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	action := decoder.StringFromMap(events[*idx].Data, "action")
	isValidMsg := action == "/ibc.core.connection.v1.MsgConnectionOpenInit" || action == "/ibc.core.connection.v1.MsgConnectionOpenTry"

	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processConnectionOpenInit(ctx, events, msg, idx)
}

func processConnectionOpenInit(_ *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	for i := *idx; i < len(events); i++ {
		if events[i].Type != storageTypes.EventTypeConnectionOpenInit && events[i].Type != storageTypes.EventTypeConnectionOpenTry {
			continue
		}
		cc := decode.NewConnectionOpen(events[i].Data)

		msg.IbcConnection = &storage.IbcConnection{
			Height:                   msg.Height,
			CreatedAt:                msg.Time,
			ClientId:                 cc.ClientId,
			ConnectionId:             cc.ConnectionId,
			CounterpartyClientId:     cc.CounterpartyClientId,
			CounterpartyConnectionId: cc.CounterpartyConnectionId,
			ChannelsCount:            0,
		}
	}
	return nil
}
