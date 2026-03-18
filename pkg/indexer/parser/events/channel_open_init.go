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

func handleChannelOpenInit(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	action := decoder.StringFromMap(events[*idx].Data, "action")
	isValidMsg := action == "/ibc.core.channel.v1.MsgChannelOpenInit" || action == "/ibc.core.channel.v1.MsgChannelOpenTry"

	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processChannelOpenInit(ctx, events, msg, idx)
}

func processChannelOpenInit(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if events[*idx].Type != storageTypes.EventTypeChannelOpenInit && events[*idx].Type != storageTypes.EventTypeChannelOpenTry {
		return errors.Errorf("invalid event type: %s", events[*idx].Type)
	}
	cc := decode.NewChannelChange(events[*idx].Data)

	channelSettings, ok := msg.Data["Channel"]
	if !ok {
		return errors.Errorf("can't receive channel settings from %s", msg.Type)
	}
	settings, ok := channelSettings.(storageTypes.PackedBytes)
	if !ok {
		return errors.Errorf("can't cast channel settings to map: %##v", channelSettings)
	}
	version, err := settings.GetString("Version")
	if err != nil {
		return errors.Wrap(err, "get string")
	}
	ordering, err := decoder.ChannelOrderingFromMap(settings, "Ordering")
	if err != nil {
		return errors.Wrap(err, "parse ordering")
	}

	signer, err := msg.Data.GetString("Signer")
	if err != nil {
		return errors.Wrap(err, "get string")
	}

	channel := &storage.IbcChannel{
		Id:                    cc.ChannelId,
		Height:                msg.Height,
		CreatedAt:             msg.Time,
		PortId:                cc.PortId,
		ConnectionId:          cc.ConnectionId,
		CounterpartyPortId:    cc.CounterpartyPortId,
		CounterpartyChannelId: cc.CounterpartyChannelId,
		Status:                storageTypes.IbcChannelStatusInitialization,
		Version:               version,
		Ordering:              ordering,
		Creator: &storage.Address{
			Address: signer,
		},
		CreateTxId: msg.TxId,
	}
	ctx.AddIbcChannel(channel)

	conn := &storage.IbcConnection{
		ConnectionId:  cc.ConnectionId,
		ChannelsCount: 1,
	}
	ctx.AddIbcConnection(conn)
	return nil
}
