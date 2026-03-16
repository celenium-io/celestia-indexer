// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package testutil

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func MustMsgToMap(t *testing.T, msg cosmosTypes.Msg) storageTypes.PackedBytes {
	t.Helper()
	m, err := decode.MsgToMap(msg)
	require.NoError(t, err)
	return m
}

// CreateExpectations builds the expected storage.Message for a decoded message test.
// In each test a fresh context is used, so the msgCounter always starts at 0:
// SetId(0) with height=0 gives id=1.
func CreateExpectations(
	t *testing.T,
	block *nodeTypes.BlockData,
	now time.Time,
	m cosmosTypes.Msg,
	position int,
	txType storageTypes.MsgType,
	size int,
) storage.Message {
	t.Helper()
	return storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  int64(position),
		Type:      txType,
		TxId:      0,
		Data:      MustMsgToMap(t, m),
		Size:      size,
		Namespace: nil,
	}
}
