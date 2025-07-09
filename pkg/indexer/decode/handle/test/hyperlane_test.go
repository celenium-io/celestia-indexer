// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	hyperlaneCore "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	hyperlaneWarp "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeMsg_SuccessOnMsgCreateMailbox(t *testing.T) {
	msg := &hyperlaneCore.MsgCreateMailbox{
		Owner:       "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		LocalDomain: 123,
		DefaultIsm:  util.CreateMockHexAddress("test", 123),
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateMailbox,
		TxId:      0,
		Size:      119,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgProcessMessage(t *testing.T) {
	msg := &hyperlaneCore.MsgProcessMessage{
		Relayer:   "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		MailboxId: util.CreateMockHexAddress("test", 123),
		Metadata:  "metadata",
		Message:   "message",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgProcessMessage,
		TxId:      0,
		Size:      136,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: []storage.AddressWithType{
			{
				Type: storageTypes.MsgAddressTypeRelayer,
				Address: storage.Address{
					Id:         0,
					Height:     block.Height,
					LastHeight: block.Height,
					Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
					Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
					Balance:    storage.EmptyBalance(),
				},
			},
		},
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, msgExpected.Addresses, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgSetMailbox(t *testing.T) {
	msg := &hyperlaneCore.MsgSetMailbox{
		Owner:     "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		MailboxId: util.CreateMockHexAddress("test", 123),
		NewOwner:  "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		}, {
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSetMailbox,
		TxId:      0,
		Size:      166,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgSetMailboxWithoutNewOwner(t *testing.T) {
	msg := &hyperlaneCore.MsgSetMailbox{
		Owner:     "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		MailboxId: util.CreateMockHexAddress("test", 123),
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSetMailbox,
		TxId:      0,
		Size:      117,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgCreateCollateralToken(t *testing.T) {
	msg := &hyperlaneWarp.MsgCreateCollateralToken{
		Owner:         "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		OriginMailbox: util.CreateMockHexAddress("test", 123),
		OriginDenom:   "DENOM",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateCollateralToken,
		TxId:      0,
		Size:      124,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgCreateSyntheticToken(t *testing.T) {
	msg := &hyperlaneWarp.MsgCreateSyntheticToken{
		Owner:         "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		OriginMailbox: util.CreateMockHexAddress("test", 123),
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateSyntheticToken,
		TxId:      0,
		Size:      117,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgSetToken(t *testing.T) {
	msg := &hyperlaneWarp.MsgSetToken{
		Owner:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		TokenId:  util.CreateMockHexAddress("test", 123),
		NewOwner: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		}, {
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSetToken,
		TxId:      0,
		Size:      166,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgEnrollRemoteRouter(t *testing.T) {
	msg := &hyperlaneWarp.MsgEnrollRemoteRouter{
		Owner:   "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		TokenId: util.CreateMockHexAddress("test", 123),
		RemoteRouter: &hyperlaneWarp.RemoteRouter{
			ReceiverDomain:   23,
			ReceiverContract: "contract",
			Gas:              math.NewInt(123456),
		},
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgEnrollRemoteRouter,
		TxId:      0,
		Size:      139,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgUnrollRemoteRouter(t *testing.T) {
	msg := &hyperlaneWarp.MsgUnrollRemoteRouter{
		Owner:          "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		TokenId:        util.CreateMockHexAddress("test", 123),
		ReceiverDomain: 23,
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeOwner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUnrollRemoteRouter,
		TxId:      0,
		Size:      119,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgRemoteTransfer(t *testing.T) {
	msg := &hyperlaneWarp.MsgRemoteTransfer{
		Sender:             "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		TokenId:            util.CreateMockHexAddress("test", 123),
		Amount:             math.NewInt(123456),
		GasLimit:           math.NewInt(456789),
		Recipient:          util.CreateMockHexAddress("Recipient", 234),
		DestinationDomain:  23,
		MaxFee:             types.NewCoin("DENOM", math.NewInt(13)),
		CustomHookMetadata: "test",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSender,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgRemoteTransfer,
		TxId:      0,
		Size:      222,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
