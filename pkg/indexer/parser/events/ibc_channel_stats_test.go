// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	decodeHandle "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/handle"
	transferTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	coreChannel "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

const (
	statsRelayer         = "celestia16p6lrlxf7f03c0ka8cv4sznr29rym27us7u4v6"
	statsCelestiaAddress = "celestia1nc44cmtgmp6cwch2ccfp6txdelu6qtz963ksrw"
	statsForeignAddress  = "osmo1vkdakqqg5htq5c3wy2kj2geq536q665xdexrtjuwqckpads2c2nsvhhcyv"
)

func statsPacket(t *testing.T, seq uint64, srcChannel, dstChannel, sender, receiver string) coreChannel.Packet {
	data, err := json.Marshal(transferTypes.FungibleTokenPacketData{
		Denom:    "utia",
		Amount:   "1000",
		Sender:   sender,
		Receiver: receiver,
	})
	require.NoError(t, err)
	return coreChannel.Packet{
		Sequence:           seq,
		SourcePort:         "transfer",
		SourceChannel:      srcChannel,
		DestinationPort:    "transfer",
		DestinationChannel: dstChannel,
		Data:               data,
	}
}

func packetMap(p coreChannel.Packet) types.PackedBytes {
	return types.PackedBytes{"Packet": map[string]any{
		"Sequence":           p.Sequence,
		"SourcePort":         p.SourcePort,
		"SourceChannel":      p.SourceChannel,
		"DestinationPort":    p.DestinationPort,
		"DestinationChannel": p.DestinationChannel,
	}}
}

// recvEvents: result is "success", "error" or "" for a redelivered packet (NOOP, no events)
func recvEvents(seq uint64, result string) []storage.Event {
	events := []storage.Event{
		{Type: "message", Data: map[string]string{"action": "/ibc.core.channel.v1.MsgRecvPacket"}},
	}
	switch result {
	case "success":
		events = append(events,
			storage.Event{Type: "recv_packet", Data: map[string]string{"packet_connection": "connection-2", "packet_dst_channel": "channel-2", "packet_sequence": strconv.FormatUint(seq, 10)}},
			storage.Event{Type: "message", Data: map[string]string{"module": "ibc_channel"}},
			storage.Event{Type: "coin_received", Data: map[string]string{"receiver": statsCelestiaAddress, "amount": "1000utia"}},
			storage.Event{Type: "fungible_token_packet", Data: map[string]string{"module": "transfer", "success": "true"}},
			storage.Event{Type: "write_acknowledgement", Data: map[string]string{"packet_connection": "connection-2"}},
			storage.Event{Type: "message", Data: map[string]string{"module": "ibc_channel"}},
		)
	case "error":
		events = append(events,
			storage.Event{Type: "recv_packet", Data: map[string]string{"packet_connection": "connection-2", "packet_dst_channel": "channel-2", "packet_sequence": strconv.FormatUint(seq, 10)}},
			storage.Event{Type: "message", Data: map[string]string{"module": "ibc_channel"}},
			storage.Event{Type: "fungible_token_packet", Data: map[string]string{"module": "transfer", "success": "false", "error": "insufficient funds"}},
			storage.Event{Type: "write_acknowledgement", Data: map[string]string{"packet_connection": "connection-2"}},
			storage.Event{Type: "message", Data: map[string]string{"module": "ibc_channel"}},
		)
	}
	return events
}

func newStatsContext() *context.Context {
	ctx := context.NewContext()
	ctx.Block = &storage.Block{Height: 100, Time: time.Now().UTC()}
	return ctx
}

type statsStep struct {
	seq    uint64
	result string
}

// scenarios share one channel within one block; each must end with exactly one confirmed transfer of seq=1
var statsScenarios = []struct {
	name  string
	steps []statsStep
}{
	{name: "failed transfer keeps other transfers' stats", steps: []statsStep{{1, "success"}, {2, "error"}}},
	{name: "redelivered packet is not counted", steps: []statsStep{{1, "success"}, {1, ""}}},
	{name: "redelivered then failed", steps: []statsStep{{1, "success"}, {1, ""}, {2, "error"}}},
}

func Test_IbcChannelStats_RecvPacket(t *testing.T) {
	for _, tt := range statsScenarios {
		t.Run(tt.name, func(t *testing.T) {
			var events []storage.Event
			for _, st := range tt.steps {
				events = append(events, recvEvents(st.seq, st.result)...)
			}

			ctx := newStatsContext()
			c := NewCursor(events)
			for i, st := range tt.steps {
				packet := statsPacket(t, st.seq, "channel-99", "channel-2", statsForeignAddress, statsCelestiaAddress)
				data := packetMap(packet)
				msgType, err := decodeHandle.MsgRecvPacket(ctx, types.StatusSuccess, nil, data, 1, uint64(i), &coreChannel.MsgRecvPacket{Packet: packet, Signer: statsRelayer})
				require.NoError(t, err)
				require.NoError(t, handleRecvPacket(ctx, c, &storage.Message{Id: uint64(i), Type: msgType, Data: data}))
			}
			require.Empty(t, c.Remaining())

			require.Len(t, ctx.IbcTransfers, 1)
			require.EqualValues(t, 1, ctx.IbcTransfers[0].Sequence)

			require.Equal(t, 1, ctx.IbcChannels.Len())
			ch, ok := ctx.IbcChannels.Get("channel-2")
			require.True(t, ok)
			require.EqualValues(t, 1, ch.TransfersCount)
			require.Equal(t, "1000", ch.Received.String())
			require.True(t, ch.Sent.IsZero())
		})
	}
}

func Test_IbcChannelStats_Acknowledgement(t *testing.T) {
	for _, tt := range statsScenarios {
		t.Run(tt.name, func(t *testing.T) {
			var events []storage.Event
			for _, st := range tt.steps {
				if st.result == "" {
					events = append(events, storage.Event{Type: "message", Data: map[string]string{"action": "/ibc.core.channel.v1.MsgAcknowledgement"}})
					continue
				}
				events = append(events, ackTransferEvents(strconv.FormatUint(st.seq, 10), st.result)...)
			}

			ctx := newStatsContext()
			c := NewCursor(events)
			for i, st := range tt.steps {
				packet := statsPacket(t, st.seq, "channel-2", "channel-99", statsCelestiaAddress, statsForeignAddress)
				data := packetMap(packet)
				msgType, err := decodeHandle.MsgAcknowledgement(ctx, types.StatusSuccess, nil, data, 1, uint64(i), &coreChannel.MsgAcknowledgement{Packet: packet, Signer: statsRelayer})
				require.NoError(t, err)
				require.NoError(t, handleAcknowledgement(ctx, c, &storage.Message{Id: uint64(i), Type: msgType, Data: data}))
			}
			require.Empty(t, c.Remaining())

			require.Len(t, ctx.IbcTransfers, 1)
			require.EqualValues(t, 1, ctx.IbcTransfers[0].Sequence)

			require.Equal(t, 1, ctx.IbcChannels.Len())
			ch, ok := ctx.IbcChannels.Get("channel-2")
			require.True(t, ok)
			require.EqualValues(t, 1, ch.TransfersCount)
			require.Equal(t, "1000", ch.Sent.String())
			require.True(t, ch.Received.IsZero())
		})
	}
}

const statsEvmAddress = "0x9a3f5e2b1c4d7e8f0a1b2c3d4e5f6a7b8c9d0e1f"

func Test_IbcTransfer_NonBech32Address(t *testing.T) {
	t.Run("recv from non-bech32 sender", func(t *testing.T) {
		ctx := newStatsContext()
		c := NewCursor(recvEvents(1, "success"))

		packet := statsPacket(t, 1, "channel-99", "channel-2", statsEvmAddress, statsCelestiaAddress)
		data := packetMap(packet)
		msgType, err := decodeHandle.MsgRecvPacket(ctx, types.StatusSuccess, nil, data, 1, 1, &coreChannel.MsgRecvPacket{Packet: packet, Signer: statsRelayer})
		require.NoError(t, err)
		require.NoError(t, handleRecvPacket(ctx, c, &storage.Message{Id: 1, Type: msgType, Data: data}))

		require.Len(t, ctx.IbcTransfers, 1)
		transfer := ctx.IbcTransfers[0]
		require.Nil(t, transfer.Sender)
		require.NotNil(t, transfer.SenderAddress)
		require.Equal(t, statsEvmAddress, *transfer.SenderAddress)
		require.NotNil(t, transfer.Receiver)
		require.Equal(t, statsCelestiaAddress, transfer.Receiver.Address)
		require.Nil(t, transfer.ReceiverAddress)
		require.Equal(t, "connection-2", transfer.ConnectionId)

		ch, ok := ctx.IbcChannels.Get("channel-2")
		require.True(t, ok)
		require.EqualValues(t, 1, ch.TransfersCount)
		require.Equal(t, "1000", ch.Received.String())
	})

	t.Run("ack of transfer to non-bech32 receiver", func(t *testing.T) {
		ctx := newStatsContext()
		c := NewCursor(ackTransferEvents("1", "success"))

		packet := statsPacket(t, 1, "channel-2", "channel-99", statsCelestiaAddress, statsEvmAddress)
		data := packetMap(packet)
		msgType, err := decodeHandle.MsgAcknowledgement(ctx, types.StatusSuccess, nil, data, 1, 1, &coreChannel.MsgAcknowledgement{Packet: packet, Signer: statsRelayer})
		require.NoError(t, err)
		require.NoError(t, handleAcknowledgement(ctx, c, &storage.Message{Id: 1, Type: msgType, Data: data}))

		require.Len(t, ctx.IbcTransfers, 1)
		transfer := ctx.IbcTransfers[0]
		require.Nil(t, transfer.Receiver)
		require.NotNil(t, transfer.ReceiverAddress)
		require.Equal(t, statsEvmAddress, *transfer.ReceiverAddress)
		require.NotNil(t, transfer.Sender)
		require.Equal(t, statsCelestiaAddress, transfer.Sender.Address)

		ch, ok := ctx.IbcChannels.Get("channel-2")
		require.True(t, ok)
		require.EqualValues(t, 1, ch.TransfersCount)
		require.Equal(t, "1000", ch.Sent.String())
	})
}

func Test_IbcTransfer_BoundToMessage(t *testing.T) {
	t.Run("failed non-bech32 transfer doesn't remove the previous one", func(t *testing.T) {
		events := append(recvEvents(1, "success"), recvEvents(2, "error")...)
		ctx := newStatsContext()
		c := NewCursor(events)

		for i, sender := range []string{statsForeignAddress, statsEvmAddress} {
			id := uint64(i + 1)
			packet := statsPacket(t, id, "channel-99", "channel-2", sender, statsCelestiaAddress)
			data := packetMap(packet)
			msgType, err := decodeHandle.MsgRecvPacket(ctx, types.StatusSuccess, nil, data, 1, id, &coreChannel.MsgRecvPacket{Packet: packet, Signer: statsRelayer})
			require.NoError(t, err)
			require.NoError(t, handleRecvPacket(ctx, c, &storage.Message{Id: id, Type: msgType, Data: data}))
		}

		require.Len(t, ctx.IbcTransfers, 1)
		require.EqualValues(t, 1, ctx.IbcTransfers[0].Sequence)
		ch, ok := ctx.IbcChannels.Get("channel-2")
		require.True(t, ok)
		require.EqualValues(t, 1, ch.TransfersCount)
	})

	t.Run("redelivered ICA packet doesn't remove the previous transfer", func(t *testing.T) {
		// the icahost packet creates no transfer; its NOOP redelivery emits no events
		events := append(recvEvents(1, "success"), recvEvents(7, "")...)
		ctx := newStatsContext()
		c := NewCursor(events)

		packet := statsPacket(t, 1, "channel-99", "channel-2", statsForeignAddress, statsCelestiaAddress)
		data := packetMap(packet)
		msgType, err := decodeHandle.MsgRecvPacket(ctx, types.StatusSuccess, nil, data, 1, 1, &coreChannel.MsgRecvPacket{Packet: packet, Signer: statsRelayer})
		require.NoError(t, err)
		require.NoError(t, handleRecvPacket(ctx, c, &storage.Message{Id: 1, Type: msgType, Data: data}))

		icaMsg := &storage.Message{
			Id:   2,
			Type: types.MsgRecvPacket,
			Data: types.PackedBytes{"Packet": map[string]any{"DestinationPort": "icahost", "DestinationChannel": "channel-5"}},
		}
		require.NoError(t, handleRecvPacket(ctx, c, icaMsg))
		require.Empty(t, c.Remaining())

		require.Len(t, ctx.IbcTransfers, 1)
		require.EqualValues(t, 1, ctx.IbcTransfers[0].Sequence)
		ch, ok := ctx.IbcChannels.Get("channel-2")
		require.True(t, ok)
		require.EqualValues(t, 1, ch.TransfersCount)
	})
}
