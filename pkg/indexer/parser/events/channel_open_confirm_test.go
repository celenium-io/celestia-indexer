// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func Test_handleChannelOpenConfirm(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *context.Context
		events []storage.Event
		msg    []*storage.Message
		idx    *int
	}{
		{
			name: "test 1",
			ctx:  context.NewContext(),
			events: []storage.Event{
				{
					Height: 1163656,
					Type:   "message",
					Data: map[string]string{
						"action": "/ibc.core.channel.v1.MsgChannelOpenAck",
					},
				}, {
					Height: 1163656,
					Type:   "channel_open_ack",
					Data: map[string]string{
						"channel_id":              "channel-32",
						"connection_id":           "connection-55",
						"counterparty_channel_id": "channel-300",
						"counterparty_port_id":    "transfer",
						"port_id":                 "transfer",
					},
				}, {
					Height: 1163656,
					Type:   "message",
					Data: map[string]string{
						"module": "ibc_channel",
					},
				},
				{
					Height: 1163656,
					Type:   "message",
					Data: map[string]string{
						"action": "/ibc.core.channel.v1.MsgChannelOpenAck",
					},
				}, {
					Height: 1163656,
					Type:   "channel_open_ack",
					Data: map[string]string{
						"channel_id":              "channel-31",
						"connection_id":           "connection-55",
						"counterparty_channel_id": "channel-301",
						"counterparty_port_id":    "transfer",
						"port_id":                 "transfer",
					},
				}, {
					Height: 1163656,
					Type:   "message",
					Data: map[string]string{
						"module": "ibc_channel",
					},
				},
			},
			msg: []*storage.Message{
				{
					Type:   types.MsgChannelOpenAck,
					Height: 1163656,
				}, {
					Type:   types.MsgChannelOpenAck,
					Height: 1163656,
				},
			},
			idx: testsuite.Ptr(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.msg {
				err := handleChannelOpenConfirm(tt.ctx, tt.events, tt.msg[i], tt.idx)
				require.NoError(t, err)

				_ = tt.ctx.IbcChannels.Range(func(_ string, value *storage.IbcChannel) (error, bool) {
					require.NotEmpty(t, value.ConnectionId)
					require.NotEmpty(t, value.Id)
					require.NotEmpty(t, value.CounterpartyChannelId)
					require.NotEmpty(t, value.CounterpartyPortId)
					require.NotEmpty(t, value.PortId)
					return nil, false
				})
			}
		})
	}
}
