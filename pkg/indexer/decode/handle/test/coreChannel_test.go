// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"encoding/base64"
	"testing"

	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icaTypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibcTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	transferTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	coreChannel "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeMsg_SuccessOnMsgChannelOpenInit(t *testing.T) {
	msg := &coreChannel.MsgChannelOpenInit{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
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
		Type:      storageTypes.MsgChannelOpenInit,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      53,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgChannelOpenTry(t *testing.T) {
	msg := &coreChannel.MsgChannelOpenTry{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
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
		Type:      storageTypes.MsgChannelOpenTry,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      55,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgChannelOpenAck(t *testing.T) {
	msg := &coreChannel.MsgChannelOpenAck{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
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
		Type:      storageTypes.MsgChannelOpenAck,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      51,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgChannelOpenConfirm(t *testing.T) {
	msg := &coreChannel.MsgChannelOpenConfirm{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
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
		Type:      storageTypes.MsgChannelOpenConfirm,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      51,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgChannelCloseInit(t *testing.T) {
	msg := &coreChannel.MsgChannelCloseInit{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
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
		Type:      storageTypes.MsgChannelCloseInit,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      49,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgChannelCloseConfirm(t *testing.T) {
	msg := &coreChannel.MsgChannelCloseConfirm{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
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
		Type:      storageTypes.MsgChannelCloseConfirm,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      51,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgRecvPacket_IcaHost(t *testing.T) {
	raw, err := base64.StdEncoding.DecodeString("eyJkYXRhIjoiQ3VjQkNpa3ZhV0pqTG1Gd2NHeHBZMkYwYVc5dWN5NTBjbUZ1YzJabGNpNTJNUzVOYzJkVWNtRnVjMlpsY2hLNUFRb0lkSEpoYm5ObVpYSVNDV05vWVc1dVpXd3RPQm9QQ2dSMWRHbGhFZ2N5TlRnek1USXdJa05qWld4bGMzUnBZVEV6Y1dVNVpuaGpaRFl6ZVcwMVozUTBabU15TXpWMVoyUjJPWHA2YW1WcWRYZHJlVGRuYkhGeE9IaDBaR00yTm5JNVp6WnpialIyWm5JMktrSnVaWFYwY205dU1YRTNjR04wT1RONGJUaHhkVzV0Ym1wNWRUSm1aV1Y2YW5wdWFIbGxaSE5tT1RZNGRUY3ljRFp0WVRKdFkyYzNOVFYwTUhOMWMyd3laM1k0L1BPcm1QVEZwcVFZIiwibWVtbyI6IiIsInR5cGUiOiJUWVBFX0VYRUNVVEVfVFgifQ==")
	require.NoError(t, err)
	msg := &coreChannel.MsgRecvPacket{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Packet: coreChannel.Packet{
			Data:            raw,
			DestinationPort: "icahost",
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
			Type: storageTypes.MsgAddressTypeSigner,
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

	data := structs.Map(msg)
	packet, ok := data["Packet"].(map[string]any)
	require.True(t, ok)
	packet["Data"] = map[string]any{
		"Memo": "",
		"Type": icaTypes.EXECUTE_TX,
		"Data": []cosmosTypes.Msg{
			&ibcTypes.MsgTransfer{
				SourcePort:    "transfer",
				SourceChannel: "channel-8",
				Token: types.Coin{
					Denom:  "utia",
					Amount: math.NewInt(2583120),
				},
				Sender:           "celestia13qe9fxcd63ym5gt4fc235ugdv9zzjejuwky7glqq8xtdc66r9g6sn4vfr6",
				Receiver:         "neutron1q7pct93xm8qunmnjyu2feezjznhyedsf968u72p6ma2mcg755t0susl2gv",
				Memo:             "",
				TimeoutTimestamp: 1749817983012370940,
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgRecvPacket,
		TxId:      0,
		Data:      data,
		Size:      426,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgRecvPacket_IcaHost2(t *testing.T) {
	raw, err := base64.StdEncoding.DecodeString("eyJkYXRhIjoiZXlKdFpYTnpZV2RsY3lJNlczc2lRSFI1Y0dVaU9pSXZZMjl6Ylc5ekxtSmhibXN1ZGpGaVpYUmhNUzVOYzJkVFpXNWtJaXdpWm5KdmJWOWhaR1J5WlhOeklqb2lZMjl6Ylc5ek1UVmpZM05vYUcxd01HZHplREk1Y1hCeGNUWm5OSHB0YkhSdWJuWm5iWGwxT1hWbGRXRmthRGw1TW01ak5YcHFNSE42YkhNMVozUmtaSG9pTENKMGIxOWhaR1J5WlhOeklqb2lZMjl6Ylc5ek1UQm9PWE4wWXpWMk5tNTBaMlY1WjJZMWVHWTVORFZ1YW5GeE5XZ3pNbkkxTTNWeGRYWjNJaXdpWVcxdmRXNTBJanBiZXlKa1pXNXZiU0k2SW5OMFlXdGxJaXdpWVcxdmRXNTBJam9pTVRBd01DSjlYWDFkZlE9PSIsIm1lbW8iOiJtZW1vIiwidHlwZSI6IlRZUEVfRVhFQ1VURV9UWCJ9")
	require.NoError(t, err)
	msg := &coreChannel.MsgRecvPacket{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Packet: coreChannel.Packet{
			Data:               raw,
			DestinationPort:    "icahost",
			DestinationChannel: "channel-1",
			Sequence:           1,
			SourceChannel:      "channel-4310",
			SourcePort:         "icacontroller-cosmos1epqzuh6myrwrp4zr8zjamcye4nvkkg9xd8ywak",
			TimeoutTimestamp:   1725050827576431600,
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
			Type: storageTypes.MsgAddressTypeSigner,
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

	data := structs.Map(msg)
	packet, ok := data["Packet"].(map[string]any)
	require.True(t, ok)
	packet["Data"] = map[string]any{
		"Memo": "memo",
		"Type": icaTypes.EXECUTE_TX,
		"Data": []cosmosTypes.Msg{
			&bankTypes.MsgSend{
				FromAddress: "cosmos15ccshhmp0gsx29qpqq6g4zmltnnvgmyu9ueuadh9y2nc5zj0szls5gtddz",
				ToAddress:   "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
				Amount:      cosmosTypes.NewCoins(cosmosTypes.NewCoin("stake", math.NewInt(1000))),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgRecvPacket,
		TxId:      0,
		Data:      data,
		Size:      544,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgRecvPacket_Transfer(t *testing.T) {
	msg := &coreChannel.MsgRecvPacket{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Packet: coreChannel.Packet{
			Data:               []byte(`{"amount":"2000000","denom":"transfer/channel-6994/utia","receiver":"celestia19863f6vse7qc8jegpmg8wzagdy7n0h6fwkzw3k","sender":"osmo19863f6vse7qc8jegpmg8wzagdy7n0h6fh8qwaf"}`),
			DestinationPort:    "transfer",
			DestinationChannel: "channel-0",
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
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		}, {
			Type: storageTypes.MsgAddressTypeReceiver,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia19863f6vse7qc8jegpmg8wzagdy7n0h6fwkzw3k",
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	data := structs.Map(msg)
	packet, ok := data["Packet"].(map[string]any)
	require.True(t, ok)
	packet["Data"] = transferTypes.FungibleTokenPacketData{
		Amount:   "2000000",
		Denom:    "transfer/channel-6994/utia",
		Receiver: "celestia19863f6vse7qc8jegpmg8wzagdy7n0h6fwkzw3k",
		Sender:   "osmo19863f6vse7qc8jegpmg8wzagdy7n0h6fh8qwaf",
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgRecvPacket,
		TxId:      0,
		Data:      data,
		Size:      253,
		Namespace: nil,
		Addresses: addressesExpected,
		IbcTransfer: &storage.IbcTransfer{
			Time:   now,
			Height: block.Height,
			Receiver: &storage.Address{
				Address: "celestia19863f6vse7qc8jegpmg8wzagdy7n0h6fwkzw3k",
				Balance: storage.EmptyBalance(),
			},
			SenderAddress: testsuite.Ptr("osmo19863f6vse7qc8jegpmg8wzagdy7n0h6fh8qwaf"),
			Amount:        decimal.RequireFromString("2000000"),
			Denom:         "utia",
			ChannelId:     "channel-0",
			Port:          "transfer",
		},
		IbcChannel: &storage.IbcChannel{
			Id:             "channel-0",
			TransfersCount: 1,
			Received:       decimal.RequireFromString("2000000"),
			Status:         storageTypes.IbcChannelStatusInitialization,
		},
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgTimeout(t *testing.T) {
	msg := &coreChannel.MsgTimeout{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
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
		Type:      storageTypes.MsgTimeout,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      55,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgTimeoutOnClose(t *testing.T) {
	msg := &coreChannel.MsgTimeoutOnClose{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
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
		Type:      storageTypes.MsgTimeoutOnClose,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      55,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgAcknowledgement(t *testing.T) {
	msg := &coreChannel.MsgAcknowledgement{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Packet: coreChannel.Packet{
			SourcePort:    "transfer",
			SourceChannel: "channel-0",
			Data:          []byte(`{"amount":"1000000","denom":"utia","receiver":"osmo1gutppfxgmwcrm4ws796ma467reu4cj8qg0fxgn","sender":"celestia1gutppfxgmwcrm4ws796ma467reu4cj8q37txyv"}`),
		},
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance:    storage.EmptyBalance(),
			},
		}, {
			Type: storageTypes.MsgAddressTypeSender,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1gutppfxgmwcrm4ws796ma467reu4cj8q37txyv",
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	data := structs.Map(msg)
	packet, ok := data["Packet"].(map[string]any)
	require.True(t, ok)
	packet["Data"] = transferTypes.FungibleTokenPacketData{
		Amount:   "1000000",
		Denom:    "utia",
		Receiver: "osmo1gutppfxgmwcrm4ws796ma467reu4cj8qg0fxgn",
		Sender:   "celestia1gutppfxgmwcrm4ws796ma467reu4cj8q37txyv",
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgAcknowledgement,
		TxId:      0,
		Data:      data,
		Size:      231,
		Namespace: nil,
		Addresses: addressesExpected,
		IbcChannel: &storage.IbcChannel{
			Id:             "channel-0",
			TransfersCount: 1,
			Sent:           decimal.RequireFromString("1000000"),
			Status:         storageTypes.IbcChannelStatusInitialization,
		},
		IbcTransfer: &storage.IbcTransfer{
			Height:          blob.Height,
			Time:            now,
			Port:            "transfer",
			ChannelId:       "channel-0",
			ReceiverAddress: testsuite.Ptr("osmo1gutppfxgmwcrm4ws796ma467reu4cj8qg0fxgn"),
			Sender: &storage.Address{
				Address: "celestia1gutppfxgmwcrm4ws796ma467reu4cj8q37txyv",
				Balance: storage.EmptyBalance(),
			},
			Amount: decimal.RequireFromString("1000000"),
			Denom:  "utia",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
