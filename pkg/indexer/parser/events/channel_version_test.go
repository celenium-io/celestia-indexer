// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func channelOpenEvents(action string, eventType types.EventType, data map[string]string) []storage.Event {
	return []storage.Event{
		{Type: types.EventTypeMessage, Data: map[string]string{"action": action}},
		{Type: eventType, Data: data},
		{Type: types.EventTypeMessage, Data: map[string]string{"module": "ibc_channel"}},
	}
}

func channelOpenMsg(msgType types.MsgType, msgVersion any) *storage.Message {
	settings := map[string]any{"Ordering": "ORDER_UNORDERED"}
	if msgVersion != nil {
		settings["Version"] = msgVersion
	}
	return &storage.Message{
		Type:   msgType,
		Height: 100,
		TxId:   1,
		Data: types.PackedBytes{
			"Channel": settings,
			"Signer":  "celestia1signer",
		},
	}
}

func Test_handleChannelOpenInit_Version(t *testing.T) {
	tests := []struct {
		name      string
		action    string
		eventType types.EventType
		eventVer  string
		msgType   types.MsgType
		msgVer    any
		want      string
	}{
		{
			name:      "try: event version wins over deprecated msg field",
			action:    "/ibc.core.channel.v1.MsgChannelOpenTry",
			eventType: types.EventTypeChannelOpenTry,
			eventVer:  `{"fee_version":"ics29-1","app_version":"ics20-1"}`,
			msgType:   types.MsgChannelOpenTry,
			msgVer:    "",
			want:      `{"fee_version":"ics29-1","app_version":"ics20-1"}`,
		}, {
			name:      "init: app-filled version from event",
			action:    "/ibc.core.channel.v1.MsgChannelOpenInit",
			eventType: types.EventTypeChannelOpenInit,
			eventVer:  "ics20-1",
			msgType:   types.MsgChannelOpenInit,
			msgVer:    "",
			want:      "ics20-1",
		}, {
			name:      "fallback to msg version without event attribute",
			action:    "/ibc.core.channel.v1.MsgChannelOpenInit",
			eventType: types.EventTypeChannelOpenInit,
			msgType:   types.MsgChannelOpenInit,
			msgVer:    "ics20-1",
			want:      "ics20-1",
		}, {
			name:      "no version anywhere",
			action:    "/ibc.core.channel.v1.MsgChannelOpenInit",
			eventType: types.EventTypeChannelOpenInit,
			msgType:   types.MsgChannelOpenInit,
			want:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]string{
				"channel_id":           "channel-1",
				"connection_id":        "connection-1",
				"counterparty_port_id": "transfer",
				"port_id":              "transfer",
			}
			if tt.eventVer != "" {
				data["version"] = tt.eventVer
			}
			ctx := context.NewContext()
			c := NewCursor(channelOpenEvents(tt.action, tt.eventType, data))

			err := handleChannelOpenInit(ctx, c, channelOpenMsg(tt.msgType, tt.msgVer))
			require.NoError(t, err)

			channel, ok := ctx.IbcChannels.Get("channel-1")
			require.True(t, ok)
			require.Equal(t, tt.want, channel.Version)
		})
	}
}

func Test_handleChannelOpenConfirm_Version(t *testing.T) {
	tests := []struct {
		name      string
		action    string
		eventType types.EventType
		msgType   types.MsgType
		msgData   types.PackedBytes
		want      string
	}{
		{
			name:      "ack takes counterparty version from msg",
			action:    "/ibc.core.channel.v1.MsgChannelOpenAck",
			eventType: types.EventTypeChannelOpenAck,
			msgType:   types.MsgChannelOpenAck,
			msgData:   types.PackedBytes{"CounterpartyVersion": "ics20-1"},
			want:      "ics20-1",
		}, {
			name:      "ack without counterparty version",
			action:    "/ibc.core.channel.v1.MsgChannelOpenAck",
			eventType: types.EventTypeChannelOpenAck,
			msgType:   types.MsgChannelOpenAck,
			want:      "",
		}, {
			name:      "confirm keeps version untouched",
			action:    "/ibc.core.channel.v1.MsgChannelOpenConfirm",
			eventType: types.EventTypeChannelOpenConfirm,
			msgType:   types.MsgChannelOpenConfirm,
			msgData:   types.PackedBytes{"CounterpartyVersion": "ics20-1"},
			want:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			c := NewCursor(channelOpenEvents(tt.action, tt.eventType, map[string]string{
				"channel_id":              "channel-1",
				"connection_id":           "connection-1",
				"counterparty_channel_id": "channel-2",
				"counterparty_port_id":    "transfer",
				"port_id":                 "transfer",
			}))
			msg := &storage.Message{Type: tt.msgType, Height: 100, TxId: 1, Data: tt.msgData}

			err := handleChannelOpenConfirm(ctx, c, msg)
			require.NoError(t, err)

			channel, ok := ctx.IbcChannels.Get("channel-1")
			require.True(t, ok)
			require.Equal(t, tt.want, channel.Version)
		})
	}
}
