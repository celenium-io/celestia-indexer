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
	"github.com/stretchr/testify/assert"
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
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgSend, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSender,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.IBCTransfer,
		TxId:      0,
		Data:      structs.Map(msgSend),
		Size:      105,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
