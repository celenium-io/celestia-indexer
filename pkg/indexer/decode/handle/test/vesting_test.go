// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/stretchr/testify/require"
)

// MsgCreateVestingAccount

func createMsgCreateVestingAccount() types.Msg {
	amount, _ := math.NewIntFromString("1000")
	m := cosmosVestingTypes.MsgCreateVestingAccount{
		FromAddress: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Amount: types.Coins{
			types.Coin{
				Denom:  "utia",
				Amount: amount,
			},
		},
		EndTime: 0,
		Delayed: false,
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreateVestingAccount(t *testing.T) {
	m := createMsgCreateVestingAccount()
	block, now := testsuite.EmptyBlock()
	position := 0
	txId := uint64(1)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, txId)

	vestingAccount := &storage.VestingAccount{
		Height: block.Height,
		Time:   block.Block.Time,
		Address: &storage.Address{
			Address: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		},
		Amount: storageTypes.NumericFromInt64(1000),
		Type:   storageTypes.VestingTypeContinuous,
		TxId:   testsuite.Ptr(txId),
	}
	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateVestingAccount,
		TxId:      1,
		Data:      mustMsgToMap(t, m),
		Size:      112,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Len(t, decodeCtx.VestingAccounts, 1)
	require.Equal(t, vestingAccount, decodeCtx.VestingAccounts[0])
}

// MsgCreatePermanentLockedAccount

func createMsgCreatePermanentLockedAccount() types.Msg {
	m := cosmosVestingTypes.MsgCreatePermanentLockedAccount{
		FromAddress: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Amount:      make(types.Coins, 0),
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreatePermanentLockedAccount(t *testing.T) {
	msgCreatePeriodicVestingAccount := createMsgCreatePermanentLockedAccount()
	block, now := testsuite.EmptyBlock()
	position := 0
	txId := uint64(1)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgCreatePeriodicVestingAccount, position, storageTypes.StatusSuccess, txId)

	vestingAccount := &storage.VestingAccount{
		Height: block.Height,
		Time:   block.Block.Time,
		Address: &storage.Address{
			Address: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		},
		Amount: storageTypes.NumericFromInt64(0),
		Type:   storageTypes.VestingTypePermanent,
		TxId:   testsuite.Ptr(txId),
	}
	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreatePermanentLockedAccount,
		TxId:      1,
		Data:      mustMsgToMap(t, msgCreatePeriodicVestingAccount),
		Size:      98,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Len(t, decodeCtx.VestingAccounts, 1)
	require.Equal(t, vestingAccount, decodeCtx.VestingAccounts[0])
}

// MsgCreatePeriodicVestingAccount

func createMsgCreatePeriodicVestingAccount() types.Msg {
	m := cosmosVestingTypes.MsgCreatePeriodicVestingAccount{
		FromAddress: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		StartTime:   1710357710,
		VestingPeriods: []cosmosVestingTypes.Period{
			{
				Length: 1000,
				Amount: types.NewCoins(types.NewCoin("utia", math.OneInt())),
			},
		},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreatePeriodicVestingAccount(t *testing.T) {
	msgCreatePeriodicVestingAccount := createMsgCreatePeriodicVestingAccount()
	block, now := testsuite.EmptyBlock()
	position := 0
	txId := uint64(1)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgCreatePeriodicVestingAccount, position, storageTypes.StatusSuccess, txId)

	startTime := time.Date(2024, 03, 13, 19, 21, 50, 0, time.UTC)
	vestingAccount := &storage.VestingAccount{
		Height: block.Height,
		Time:   block.Block.Time,
		Address: &storage.Address{
			Address: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		},
		Amount:    storageTypes.NumericFromInt64(1),
		Type:      storageTypes.VestingTypePeriodic,
		StartTime: &startTime,
		TxId:      testsuite.Ptr(txId),
		VestingPeriods: []storage.VestingPeriod{
			{
				Height: block.Height,
				Amount: storageTypes.NumericFromInt64(1),
				Time:   time.Date(2024, 03, 13, 19, 38, 30, 0, time.UTC),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreatePeriodicVestingAccount,
		TxId:      1,
		Data:      mustMsgToMap(t, msgCreatePeriodicVestingAccount),
		Size:      120,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Len(t, decodeCtx.VestingAccounts, 1)
	require.Equal(t, vestingAccount, decodeCtx.VestingAccounts[0])
}
