// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	minfeeTypes "github.com/celestiaorg/celestia-app/v7/x/minfee/types"
	fee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	"github.com/stretchr/testify/require"
)

func TestDecodeMsg_SuccessOnMsgRegisterPayee(t *testing.T) {
	msg := &fee.MsgRegisterPayee{
		Relayer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Payee:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
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
		Type:      storageTypes.MsgRegisterPayee,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      98,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgRegisterCounterpartyPayee(t *testing.T) {
	msg := &fee.MsgRegisterCounterpartyPayee{
		Relayer:           "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		CounterpartyPayee: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
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
		Type:      storageTypes.MsgRegisterCounterpartyPayee,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      98,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgPayPacketFee(t *testing.T) {
	msg := &fee.MsgPayPacketFee{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
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
		Type:      storageTypes.MsgPayPacketFee,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      51,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgPayPacketFeeAsync(t *testing.T) {
	msg := &fee.MsgPayPacketFeeAsync{}
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
		Type:      storageTypes.MsgPayPacketFeeAsync,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      6,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgUpdateMinfeeParams(t *testing.T) {
	msg := &minfeeTypes.MsgUpdateMinfeeParams{
		Authority: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Params:    minfeeTypes.DefaultParams(),
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
		Type:      storageTypes.MsgUpdateMinfeeParams,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      66,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
