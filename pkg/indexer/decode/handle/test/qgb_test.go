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
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/legacy"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// MsgRegisterEvmAddress

func createMsgRegisterEvmAddress() types.Msg {
	m := legacy.MsgRegisterEVMAddress{
		ValidatorAddress: "celestiavaloper1f5crra7r5m9kd6saw077u76x0n7dyjkkzk0qup",
		EvmAddress:       "0xfDC46fBDd8AF50d9Bf7536Bf44ce8560E423352c",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgRegisterEvmAddress(t *testing.T) {
	m := createMsgRegisterEvmAddress()
	block, now := testsuite.EmptyBlock()
	position := 4

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
		Position:  4,
		Type:      storageTypes.MsgRegisterEVMAddress,
		TxId:      0,
		Data:      mustMsgToMap(t, m),
		Size:      100,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
