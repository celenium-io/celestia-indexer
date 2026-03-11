// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	fwdTypes "github.com/celestiaorg/celestia-app/v7/x/forwarding/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"
)

func TestMsgForward(t *testing.T) {
	msg := &fwdTypes.MsgForward{
		Signer:        "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		DestDomain:    123,
		DestRecipient: "address",
		ForwardAddr:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		MaxIgpFee:     types.NewCoin("utia", math.NewInt(123)),
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   now,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgForward,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      122,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
