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

func handleChannelClose(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	action := decoder.StringFromMap(events[*idx].Data, "action")
	isValidMsg := action == "/ibc.core.channel.v1.MsgChannelCloseConfirm" || action == "/ibc.core.channel.v1.MsgChannelCloseInit"

	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processChannelClose(ctx, events, msg, idx)
}

func processChannelClose(_ *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if events[*idx].Type != storageTypes.EventTypeChannelCloseConfirm {
		return errors.Errorf("invalid event type: %s", events[*idx].Type)
	}
	cc := decode.NewChannelChange(events[*idx].Data)

	msg.IbcChannel = &storage.IbcChannel{
		Id:     cc.ChannelId,
		Status: storageTypes.IbcChannelStatusClosed,
	}

	toTheNextAction(events, idx)
	return nil
}
