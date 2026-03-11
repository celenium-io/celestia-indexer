// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"encoding/base64"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	icaTypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transferTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	coreChannel "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"
)

func TestDecodeMsg_SuccessOnMsgChannelOpenInit(t *testing.T) {
	msg := &coreChannel.MsgChannelOpenInit{
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
		Type:      storageTypes.MsgChannelOpenInit,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      53,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgChannelOpenTry(t *testing.T) {
	msg := &coreChannel.MsgChannelOpenTry{
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
		Type:      storageTypes.MsgChannelOpenTry,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      55,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgChannelOpenAck(t *testing.T) {
	msg := &coreChannel.MsgChannelOpenAck{
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
		Type:      storageTypes.MsgChannelOpenAck,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      51,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgChannelOpenConfirm(t *testing.T) {
	msg := &coreChannel.MsgChannelOpenConfirm{
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
		Type:      storageTypes.MsgChannelOpenConfirm,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      51,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgChannelCloseInit(t *testing.T) {
	msg := &coreChannel.MsgChannelCloseInit{
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
		Type:      storageTypes.MsgChannelCloseInit,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgChannelCloseConfirm(t *testing.T) {
	msg := &coreChannel.MsgChannelCloseConfirm{
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
		Type:      storageTypes.MsgChannelCloseConfirm,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      51,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
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

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	data := structs.Map(msg)
	packet, ok := data["Packet"].(map[string]any)
	require.True(t, ok)
	packet["Data"] = map[string]any{
		"Memo": "",
		"Type": icaTypes.EXECUTE_TX,
		"Data": []cosmosTypes.Msg{
			&transferTypes.MsgTransfer{
				SourcePort:    "transfer",
				SourceChannel: "channel-8",
				Token: cosmosTypes.Coin{
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
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgRecvPacket,
		TxId:      0,
		Data:      data,
		Size:      426,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
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

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

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
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgRecvPacket,
		TxId:      0,
		Data:      data,
		Size:      544,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgRecvPacket_IcaHost3(t *testing.T) {
	raw, err := base64.StdEncoding.DecodeString("eyJkYXRhIjoiQ293Q0NoNHZZMjl6Ylc5ekxtRjFkR2g2TG5ZeFltVjBZVEV1VFhOblIzSmhiblFTNlFFS1EyTmxiR1Z6ZEdsaE1YRnhkSEkwY0hwb09YSmphbmR3ZW5FMU0zYzFNRFIwZWpjeU4yMDFlV2g0ZEhSNU1ubG9lV1ZoYlhwNWR6aHhPR0V6WkhOeWEyaDFibXdTUTJObGJHVnpkR2xoTVdwbWRtb3dPR2h6YUROM2RtdDBiVFZtTW1kc2QzaHpZM1p5TTNKbFkyRXdOSEF3WkRabWJHNXhOVzQyT1d4a2JqSnJaM05oZHpOd2QzY2FYUXBUQ2lvdlkyOXpiVzl6TG1GMWRHaDZMbll4WW1WMFlURXVSMlZ1WlhKcFkwRjFkR2h2Y21sNllYUnBiMjRTSlFvakwyTnZjMjF2Y3k1emRHRnJhVzVuTG5ZeFltVjBZVEV1VFhOblJHVnNaV2RoZEdVU0JnamUxcmJFRkFyQkFRb3lMMk52YzIxdmN5NWthWE4wY21saWRYUnBiMjR1ZGpGaVpYUmhNUzVOYzJkVFpYUlhhWFJvWkhKaGQwRmtaSEpsYzNNU2lnRUtRMk5sYkdWemRHbGhNWEZ4ZEhJMGNIcG9PWEpqYW5kd2VuRTFNM2MxTURSMGVqY3lOMjAxZVdoNGRIUjVNbmxvZVdWaGJYcDVkemh4T0dFelpITnlhMmgxYm13U1EyTmxiR1Z6ZEdsaE1YcDNibTA1YzNCaGNYWjZObVp5TTJoeFpUQnplSFY2ZEhSNFpIUXlkbWRuY1dZMmNUQnliSG8wZG5wNk1qSTRaRGN5YW5Gd2F6TmpNbVU9IiwibWVtbyI6IiIsInR5cGUiOiJUWVBFX0VYRUNVVEVfVFgifQ==")
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

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)

	data := structs.Map(msg)
	packet, ok := data["Packet"].(map[string]any)
	require.True(t, ok)

	expirationDate := time.Date(2144, 9, 3, 17, 48, 14, 0, time.UTC)
	packet["Data"] = map[string]any{
		"Memo": "",
		"Type": icaTypes.EXECUTE_TX,
		"Data": []cosmosTypes.Msg{
			&authz.MsgGrant{
				Granter: "celestia1qqtr4pzh9rcjwpzq53w504tz727m5yhxtty2yhyeamzyw8q8a3dsrkhunl",
				Grantee: "celestia1jfvj08hsh3wvktm5f2glwxscvr3reca04p0d6flnq5n69ldn2kgsaw3pww",
				Grant: authz.Grant{
					Expiration: &expirationDate,
				},
			},
			&distributionTypes.MsgSetWithdrawAddress{
				DelegatorAddress: "celestia1qqtr4pzh9rcjwpzq53w504tz727m5yhxtty2yhyeamzyw8q8a3dsrkhunl",
				WithdrawAddress:  "celestia1zwnm9spaqvz6fr3hqe0sxuzttxdt2vggqf6q0rlz4vzz228d72jqpk3c2e",
			},
		},
	}

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgRecvPacket,
		TxId:      0,
		Data:      data,
		Size:      836,
		Namespace: nil,
	}

	_, err = dm.Msg.Data.ToBytes()
	require.NoError(t, err)

	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)

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

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

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
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgRecvPacket,
		TxId:      0,
		Data:      data,
		Size:      253,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgTimeout(t *testing.T) {
	msg := &coreChannel.MsgTimeout{
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
		Type:      storageTypes.MsgTimeout,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      55,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgTimeoutOnClose(t *testing.T) {
	msg := &coreChannel.MsgTimeoutOnClose{
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
		Type:      storageTypes.MsgTimeoutOnClose,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      55,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
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
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

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
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgAcknowledgement,
		TxId:      0,
		Data:      data,
		Size:      231,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
