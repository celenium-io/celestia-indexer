// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"
)

// MsgSend

func createMsgSend() types.Msg {
	amount, _ := math.NewIntFromString("1000")
	m := cosmosBankTypes.MsgSend{
		FromAddress: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Amount: types.Coins{
			types.Coin{
				Denom:  "utia",
				Amount: amount,
			},
		},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgSend(t *testing.T) {
	msgSend := createMsgSend()
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgSend, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSend,
		TxId:      0,
		Data:      structs.Map(msgSend),
		Size:      112,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

// MsgMultiSend

func createMsgMultiSend() types.Msg {
	m := cosmosBankTypes.MsgMultiSend{
		Inputs: []cosmosBankTypes.Input{
			{Address: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6", Coins: make(types.Coins, 0)},
			{Address: "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts", Coins: make(types.Coins, 0)},
		},
		Outputs: []cosmosBankTypes.Output{
			{Address: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l", Coins: make(types.Coins, 0)},
		},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgMultiSend(t *testing.T) {
	msgMultiSend := createMsgMultiSend()
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgMultiSend, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgMultiSend,
		TxId:      0,
		Data:      structs.Map(msgMultiSend),
		Size:      153,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgSetSendEnabled(t *testing.T) {
	msgMultiSend := &cosmosBankTypes.MsgSetSendEnabled{
		Authority: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		SendEnabled: []*cosmosBankTypes.SendEnabled{
			{
				Enabled: true,
				Denom:   "utia",
			},
		},
		UseDefaultFor: []string{
			"utia",
		},
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgMultiSend, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSetSendEnabled,
		TxId:      0,
		Data:      structs.Map(msgMultiSend),
		Size:      65,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
