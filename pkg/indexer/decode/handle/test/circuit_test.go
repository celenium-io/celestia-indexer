// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	circuitTypes "cosmossdk.io/x/circuit/types"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func TestDecodeMsg_SuccessOnMsgAuthorizeCircuitBreaker(t *testing.T) {
	msg := &circuitTypes.MsgAuthorizeCircuitBreaker{
		Granter: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Grantee: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Permissions: &circuitTypes.Permissions{
			Level: circuitTypes.Permissions_LEVEL_ALL_MSGS,
		},
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgAuthorizeCircuitBreaker,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      102,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgResetCircuitBreaker(t *testing.T) {
	msg := &circuitTypes.MsgResetCircuitBreaker{
		Authority: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgResetCircuitBreaker,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgTripCircuitBreaker(t *testing.T) {
	msg := &circuitTypes.MsgTripCircuitBreaker{
		Authority: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTripCircuitBreaker,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
