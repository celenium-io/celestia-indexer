// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// MsgEditValidator

func createMsgEditValidator() types.Msg {
	m := cosmosStakingTypes.MsgEditValidator{
		Description: cosmosStakingTypes.Description{
			Moniker:         "newAgeValidator",
			Identity:        "UPort:1",
			Website:         "https://google.com",
			SecurityContact: "tryme@gmail.com",
			Details:         "trust",
		},
		ValidatorAddress:  "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		CommissionRate:    nil,
		MinSelfDelegation: nil,
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgEditValidator(t *testing.T) {
	m := createMsgEditValidator()
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

	expectedValidators := map[string]*storage.Validator{
		"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x": {
			Address:           "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
			Moniker:           "newAgeValidator",
			Identity:          "UPort:1",
			Website:           "https://google.com",
			Contacts:          "tryme@gmail.com",
			Details:           "trust",
			Rate:              decimal.Zero,
			MinSelfDelegation: decimal.Zero,
			Stake:             decimal.Zero,
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgEditValidator,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      128,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
	for key, value := range expectedValidators {
		val, ok := decodeCtx.Validators.Get(key)
		require.True(t, ok)
		require.Equal(t, value, val)
	}
}

// MsgBeginRedelegate

func createMsgBeginRedelegate() types.Msg {
	m := cosmosStakingTypes.MsgBeginRedelegate{
		DelegatorAddress:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorSrcAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		ValidatorDstAddress: "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
		Amount:              types.NewCoin(currency.Utia, math.OneInt()),
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgBeginRedelegate(t *testing.T) {
	m := createMsgBeginRedelegate()
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
			Type: storageTypes.MsgAddressTypeValidatorSrc,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:       []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance:    storage.EmptyBalance(),
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidatorDst,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
				Hash:       []byte{0x56, 0x35, 0x87, 0x35, 0xf6, 0x7, 0xd1, 0x53, 0xc1, 0xb7, 0x94, 0xe6, 0x54, 0xbc, 0x47, 0x52, 0x1f, 0x83, 0xa9, 0x3b},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgBeginRedelegate,
		TxId:      0,
		Size:      172,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
}

// MsgCreateValidator

var pk = ed25519.GenPrivKeyFromSecret([]byte{0, 1, 2, 3})

func createMsgCreateValidator() types.Msg {
	pkAny, _ := codectypes.NewAnyWithValue(pk.PubKey())
	m := cosmosStakingTypes.MsgCreateValidator{
		Description:       cosmosStakingTypes.Description{},
		Commission:        cosmosStakingTypes.CommissionRates{},
		MinSelfDelegation: math.NewInt(1),
		ValidatorAddress:  "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		Pubkey:            pkAny,
		Value:             types.NewCoin("utia", math.OneInt()),
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreateValidator(t *testing.T) {
	m := createMsgCreateValidator()
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
	addressesExpected[0].Balance.Delegated = decimal.RequireFromString("1")

	expectedValidators := map[string]*storage.Validator{
		"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x": {
			Delegator:         "celestia1fg9l3xvfuu9wxremv2229966zawysg4rss2hzq",
			Address:           "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
			ConsAddress:       "A8BEA00847066E6C765E7B064DD79265406D402B",
			Rate:              decimal.Zero,
			MaxRate:           decimal.Zero,
			MaxChangeRate:     decimal.Zero,
			Stake:             decimal.RequireFromString("1"),
			MinSelfDelegation: decimal.RequireFromString("1"),
			Height:            blob.Height,
			Jailed:            testsuite.Ptr(false),
		},
	}

	data := structs.Map(m)
	data["Pubkey"] = map[string]any{
		"key":  pk.PubKey().Bytes(),
		"type": "ed25519",
	}
	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateValidator,
		TxId:      0,
		Data:      data,
		Size:      152,
		Namespace: nil,
		Addresses: addressesExpected,
	}
	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
	for key, value := range expectedValidators {
		val, ok := decodeCtx.Validators.Get(key)
		require.True(t, ok)
		require.Equal(t, value, val)
	}
}

// MsgDelegate

func createMsgDelegate() types.Msg {
	amount, _ := math.NewIntFromString("1000")
	msgDelegate := cosmosStakingTypes.MsgDelegate{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
		Amount: types.Coin{
			Denom:  "utia",
			Amount: amount,
		},
	}

	return &msgDelegate
}

func TestDecodeMsg_SuccessOnMsgDelegate(t *testing.T) {
	msgDelegate := createMsgDelegate()
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgDelegate, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:       []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance:    storage.EmptyBalance(),
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
				Hash:       []byte{0x56, 0x35, 0x87, 0x35, 0xf6, 0x7, 0xd1, 0x53, 0xc1, 0xb7, 0x94, 0xe6, 0x54, 0xbc, 0x47, 0x52, 0x1f, 0x83, 0xa9, 0x3b},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgDelegate,
		TxId:      0,
		Data:      structs.Map(msgDelegate),
		Size:      119,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
}

// MsgUndelegate

func createMsgUndelegate() types.Msg {
	amount, _ := math.NewIntFromString("1001")
	m := cosmosStakingTypes.MsgUndelegate{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
		Amount: types.Coin{
			Denom:  "utia",
			Amount: amount,
		},
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgUndelegate(t *testing.T) {
	m := createMsgUndelegate()
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
				Address:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:       []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance:    storage.EmptyBalance(),
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
				Hash:       []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUndelegate,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      119,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
}

// MsgCancelUnbondingDelegation

func createMsgCancelUnbondingDelegation() types.Msg {
	amount, _ := math.NewIntFromString("1001")
	m := cosmosStakingTypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
		Amount: types.Coin{
			Denom:  "utia",
			Amount: amount,
		},
		CreationHeight: 100,
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgCancelUnbondingDelegation(t *testing.T) {
	m := createMsgCancelUnbondingDelegation()
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
				Address:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:       []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance:    storage.EmptyBalance(),
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
				Hash:       []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCancelUnbondingDelegation,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      121,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, addressesExpected, dm.Addresses)
}
