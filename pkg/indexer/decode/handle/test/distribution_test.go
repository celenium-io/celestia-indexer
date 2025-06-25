// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

// MsgSetWithdrawAddress

func createMsgSetWithdrawAddress() types.Msg {
	m := cosmosDistributionTypes.MsgSetWithdrawAddress{
		DelegatorAddress: "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
		WithdrawAddress:  "celestia1nasjhf82cjuk3mxyhzw6ntpc66exzfe7qhl256",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgSetWithdrawAddress(t *testing.T) {
	m := createMsgSetWithdrawAddress()
	blob, now := testsuite.EmptyBlock()
	position := 4

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
				Hash:       []byte{0xe5, 0x3, 0xb, 0xac, 0x1, 0xc9, 0xa5, 0xbe, 0x71, 0xa3, 0x60, 0x34, 0x8, 0xc6, 0x80, 0x26, 0xd4, 0x24, 0x20, 0x57},
				Balance:    storage.EmptyBalance(),
			},
		}, {
			Type: storageTypes.MsgAddressTypeWithdraw,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1nasjhf82cjuk3mxyhzw6ntpc66exzfe7qhl256",
				Hash:       []byte{0x9f, 0x61, 0x2b, 0xa4, 0xea, 0xc4, 0xb9, 0x68, 0xec, 0xc4, 0xb8, 0x9d, 0xa9, 0xac, 0x38, 0xd6, 0xb2, 0x61, 0x27, 0x3e},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgSetWithdrawAddress,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      98,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgWithdrawDelegatorReward

func createMsgWithdrawDelegatorReward() types.Msg {
	m := cosmosDistributionTypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgWithdrawDelegatorReward(t *testing.T) {
	m := createMsgWithdrawDelegatorReward()
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
				Hash:       []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
				Balance:    storage.EmptyBalance(),
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:       []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgWithdrawDelegatorReward,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      105,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgWithdrawValidatorCommission

func createMsgWithdrawValidatorCommission() types.Msg {
	m := cosmosDistributionTypes.MsgWithdrawValidatorCommission{
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgWithdrawValidatorCommission(t *testing.T) {
	m := createMsgWithdrawValidatorCommission()
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:       []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgWithdrawValidatorCommission,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      56,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgFundCommunityPool

func createMsgFundCommunityPool() types.Msg {
	m := cosmosDistributionTypes.MsgFundCommunityPool{
		Amount:    nil,
		Depositor: "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgFundCommunityPool(t *testing.T) {
	m := createMsgFundCommunityPool()
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDepositor,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
				Hash:       []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgFundCommunityPool,
		TxId:      0,
		Size:      49,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
