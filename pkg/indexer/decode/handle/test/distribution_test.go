// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
	"github.com/stretchr/testify/require"
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
	block, now := testsuite.EmptyBlock()
	position := 4

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
		Position:  4,
		Type:      storageTypes.MsgSetWithdrawAddress,
		TxId:      0,
		Data:      mustMsgToMap(t, m),
		Size:      98,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
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
		Type:      storageTypes.MsgWithdrawDelegatorReward,
		TxId:      0,
		Data:      mustMsgToMap(t, m),
		Size:      105,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
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
		Type:      storageTypes.MsgWithdrawValidatorCommission,
		TxId:      0,
		Data:      mustMsgToMap(t, m),
		Size:      56,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
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
		Type:      storageTypes.MsgFundCommunityPool,
		TxId:      0,
		Size:      49,
		Data:      mustMsgToMap(t, m),
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
