// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	dCtx "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/require"
)

// TestCursorGolden pins down the exact output of parseTxs on a real production
// block (height 8880520) whose first transaction carries nine messages —
// one MsgUpdateClient followed by eight consecutive MsgRecvPacket — sharing a
// single events cursor across message boundaries.
//
// This is the scenario the events-cursor refactor (idx *int -> Cursor) is
// riskiest for: eight same-type messages in a row exercise the "module"-keyed
// IBC dispatch (handle/ibcEventHandlers) repeatedly, back to back, with no
// other message type in between to reset any accidental shared state. If the
// refactor shifts the cursor by even one event, message count, positions, or
// the derived IBC channel state will diverge from the values captured here on
// the pre-refactor code, and this test will fail.
//
// The expected values below were captured from the current (pre-refactor)
// parser output, not derived from first principles — treat them as a
// snapshot to preserve, not as a spec to re-derive.
func TestCursorGolden(t *testing.T) {
	block, err := getBlockByHeight(8880520)
	require.NoError(t, err)

	decodeCtx := dCtx.NewContext()
	decodeCtx.Block = &storage.Block{
		Height:       8880520,
		MessageTypes: storageTypes.NewMsgTypeBitMask(),
	}

	p := NewModule(config.Indexer{})
	resultTxs, err := p.parseTxs(decodeCtx, &block)
	require.NoError(t, err)
	require.Len(t, resultTxs, 5)

	// tx0: MsgUpdateClient + 8x MsgRecvPacket sharing one cursor.
	tx0 := resultTxs[0]
	require.Equal(t, storageTypes.StatusSuccess, tx0.Status)
	require.Empty(t, tx0.Error)
	require.EqualValues(t, 4998079, tx0.GasUsed)
	require.EqualValues(t, 5999331, tx0.GasWanted)
	require.ElementsMatch(t, []storageTypes.MsgType{
		storageTypes.MsgUpdateClient,
		storageTypes.MsgRecvPacket,
	}, tx0.MessageTypes.Names())

	require.Equal(t, storageTypes.StatusSuccess, resultTxs[1].Status)
	require.ElementsMatch(t, []storageTypes.MsgType{storageTypes.IBCTransfer}, resultTxs[1].MessageTypes.Names())

	for i := 2; i <= 4; i++ {
		require.Equal(t, storageTypes.StatusSuccess, resultTxs[i].Status)
		require.ElementsMatch(t, []storageTypes.MsgType{storageTypes.MsgPayForBlobs}, resultTxs[i].MessageTypes.Names())
	}

	// The cursor must split tx0's shared event slice into exactly nine
	// distinct messages, in order, each with its own position — not merge
	// adjacent MsgRecvPacket messages together or drop/duplicate one.
	require.Len(t, decodeCtx.Messages, 13)

	wantTypes := []storageTypes.MsgType{
		storageTypes.MsgUpdateClient,
		storageTypes.MsgRecvPacket, storageTypes.MsgRecvPacket, storageTypes.MsgRecvPacket, storageTypes.MsgRecvPacket,
		storageTypes.MsgRecvPacket, storageTypes.MsgRecvPacket, storageTypes.MsgRecvPacket, storageTypes.MsgRecvPacket,
		storageTypes.IBCTransfer,
		storageTypes.MsgPayForBlobs, storageTypes.MsgPayForBlobs, storageTypes.MsgPayForBlobs,
	}
	wantPositions := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 0}
	for i, msg := range decodeCtx.Messages {
		require.Equalf(t, wantTypes[i], msg.Type, "message %d type", i)
		require.Equalf(t, wantPositions[i], msg.Position, "message %d position", i)
		require.Equalf(t, pkgTypes.Level(8880520), msg.Height, "message %d height", i)
		require.Emptyf(t, msg.InternalMsgs, "message %d must not have internal msgs", i)
	}

	// The eight MsgRecvPacket handlers converge on a single IBC channel
	// derived from the packets' shared destination channel.
	require.Equal(t, 1, decodeCtx.IbcChannels.Len())
	channel, ok := decodeCtx.IbcChannels.Get("channel-4")
	require.True(t, ok)
	require.Equal(t, storageTypes.IbcChannelStatusInitialization, channel.Status)
	require.Equal(t, "200176", channel.Received.String())
	require.True(t, channel.Sent.IsZero())
	require.EqualValues(t, 1, channel.TransfersCount)

	// Every in-flight transfer opened during processing was resolved
	// (matched to a fungible-token-packet event or explicitly discarded) by
	// the time parseTxs returns — none should be left dangling.
	require.Empty(t, decodeCtx.IbcTransfers)

	require.Equal(t, 71, decodeCtx.Addresses.Len())
	require.Len(t, decodeCtx.BlobLogs, 32)
	require.Len(t, decodeCtx.Events, 356)
	require.ElementsMatch(t, []storageTypes.MsgType{
		storageTypes.MsgPayForBlobs,
		storageTypes.IBCTransfer,
		storageTypes.MsgUpdateClient,
		storageTypes.MsgRecvPacket,
	}, decodeCtx.Block.MessageTypes.Names())
}
