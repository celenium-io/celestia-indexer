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
	"github.com/cosmos/cosmos-sdk/types"
	ibcTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcCoreClientTypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"
)

// MsgTransfer

func createIBCMsgTransfer() types.Msg {
	m := ibcTypes.MsgTransfer{
		SourcePort:       "",
		SourceChannel:    "",
		Token:            types.Coin{},
		Sender:           "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Receiver:         "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		TimeoutHeight:    ibcCoreClientTypes.Height{},
		TimeoutTimestamp: 0,
		Memo:             "",
	}

	return &m
}

func TestDecodeMsg_SuccessOnIBCMsgTransfer(t *testing.T) {
	msgSend := createIBCMsgTransfer()
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
		Type:      storageTypes.IBCTransfer,
		TxId:      0,
		Data:      structs.Map(msgSend),
		Size:      105,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
