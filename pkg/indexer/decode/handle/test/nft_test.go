// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	"cosmossdk.io/x/nft"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// MsgSend (NFT)

func createMsgSendNFT() types.Msg {
	m := nft.MsgSend{
		Sender:   "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Receiver: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgSendNFT(t *testing.T) {
	msg := createMsgSendNFT()
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
		Type:      storageTypes.MsgSendNFT,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      98,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
