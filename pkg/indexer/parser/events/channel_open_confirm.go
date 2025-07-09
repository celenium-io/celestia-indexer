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

func handleChannelOpenConfirm(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	action := decoder.StringFromMap(events[*idx].Data, "action")
	isValidMsg := action == "/ibc.core.channel.v1.MsgChannelOpenConfirm" || action == "/ibc.core.channel.v1.MsgChannelOpenAck"

	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processChannelOpenConfirm(ctx, events, msg, idx)
}

func processChannelOpenConfirm(_ *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if events[*idx].Type != storageTypes.EventTypeChannelOpenConfirm && events[*idx].Type != storageTypes.EventTypeChannelOpenAck {
		return errors.Errorf("invalid event type: %s", events[*idx].Type)
	}
	cc := decode.NewChannelChange(events[*idx].Data)

	msg.IbcChannel = &storage.IbcChannel{
		Id:                    cc.ChannelId,
		ConfirmationHeight:    msg.Height,
		ConfirmedAt:           msg.Time,
		PortId:                cc.PortId,
		ConnectionId:          cc.ConnectionId,
		CounterpartyPortId:    cc.CounterpartyPortId,
		CounterpartyChannelId: cc.CounterpartyChannelId,
		Status:                storageTypes.IbcChannelStatusOpened,
	}

	*idx += 2
	return nil
}
