// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"
)

// MsgUnjail

func createMsgUnjail() types.Msg {
	m := cosmosSlashingTypes.MsgUnjail{
		ValidatorAddr: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgUnjail(t *testing.T) {
	m := createMsgUnjail()
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUnjail,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      56,
		Namespace: nil,
	}
	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
