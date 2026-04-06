// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	expectedValidators := map[string]*storage.Validator{
		"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x": {
			Address:           "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
			Moniker:           "newAgeValidator",
			Identity:          "UPort:1",
			Website:           "https://google.com",
			Contacts:          "tryme@gmail.com",
			Details:           "trust",
			Rate:              storageTypes.NumericZero(),
			MinSelfDelegation: storageTypes.NumericZero(),
			Stake:             storageTypes.NumericZero(),
			MessagesCount:     1,
		},
	}

	msgExpected := storage.Message{
		Id:         1,
		Height:     block.Height,
		Time:       now,
		Position:   0,
		Type:       storageTypes.MsgEditValidator,
		TxId:       0,
		Data:       mustMsgToMap(t, m),
		Size:       128,
		Namespace:  nil,
		Validators: []string{"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"},
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Len(t, dm.Msg.Validators, 1)
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
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgBeginRedelegate,
		TxId:      0,
		Size:      172,
		Data:      mustMsgToMap(t, m),
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
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
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	expectedValidators := map[string]*storage.Validator{
		"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x": {
			Delegator:         "celestia1fg9l3xvfuu9wxremv2229966zawysg4rss2hzq",
			Address:           "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
			ConsAddress:       "A8BEA00847066E6C765E7B064DD79265406D402B",
			Rate:              storageTypes.NumericZero(),
			MaxRate:           storageTypes.NumericZero(),
			MaxChangeRate:     storageTypes.NumericZero(),
			Stake:             storageTypes.NumericFromInt64(1),
			MinSelfDelegation: storageTypes.NumericFromInt64(1),
			Height:            block.Height,
			Jailed:            testsuite.Ptr(false),
			MessagesCount:     1,
			CreationTime:      block.Block.Time,
		},
	}

	data := mustMsgToMap(t, m)
	data["Pubkey"] = map[string]any{
		"key":  pk.PubKey().Bytes(),
		"type": "ed25519",
	}
	msgExpected := storage.Message{
		Id:         1,
		Height:     block.Height,
		Time:       now,
		Position:   0,
		Type:       storageTypes.MsgCreateValidator,
		TxId:       0,
		Data:       data,
		Size:       152,
		Namespace:  nil,
		Validators: []string{"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"},
	}
	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Len(t, dm.Msg.Validators, 1)
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
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgDelegate, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgDelegate,
		TxId:      0,
		Data:      mustMsgToMap(t, msgDelegate),
		Size:      119,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
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
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUndelegate,
		TxId:      0,
		Data:      mustMsgToMap(t, m),
		Size:      119,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
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
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCancelUnbondingDelegation,
		TxId:      0,
		Data:      mustMsgToMap(t, m),
		Size:      121,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
