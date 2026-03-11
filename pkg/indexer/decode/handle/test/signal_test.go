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
	signal "github.com/celestiaorg/celestia-app/v7/x/signal/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestDecodeMsg_SuccessOnMsgSignalVersion(t *testing.T) {
	m := signal.NewMsgSignalVersion(
		"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x", 10,
	)
	block, now := testsuite.EmptyBlock()
	position := 7
	txId := uint64(1)

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, txId)

	msgExpected := createExpectations(
		block, now, m, position,
		storageTypes.MsgAddressTypeValidator,
		"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		[]byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
		storageTypes.MsgSignalVersion,
		58,
	)
	msgExpected.TxId = txId

	validator := storage.EmptyValidator()
	validator.Address = m.ValidatorAddress
	validator.Version = 10

	signal := &storage.SignalVersion{
		Height:    block.Height,
		Validator: &validator,
		Time:      block.Block.Time,
		Version:   m.Version,
		MsgId:     1,
		TxId:      txId,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)

	require.EqualValues(t, 1, decodeCtx.Upgrades.Len())
	upgrade, ok := decodeCtx.Upgrades.Get(10)
	require.True(t, ok)
	require.EqualValues(t, 1, upgrade.SignalsCount)
	require.EqualValues(t, block.Height, upgrade.Height)
	require.EqualValues(t, block.Block.Time, upgrade.Time)
	require.Len(t, decodeCtx.Signals, 1)
	require.Equal(t, signal, decodeCtx.Signals[0])
}

func TestDecodeMsg_SuccessOnMsgTryUpgrade(t *testing.T) {
	valAddress, err := types.AccAddressFromBech32("celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r")
	require.NoError(t, err)
	m := signal.NewMsgTryUpgrade(
		valAddress,
	)
	block, now := testsuite.EmptyBlock()
	position := 7

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	msgExpected := createExpectations(
		block, now, m, position,
		storageTypes.MsgAddressTypeSigner,
		"celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
		[]byte{0x50, 0xa1, 0xec, 0xc6, 0x67, 0x0c, 0x9a, 0x72, 0x1f, 0x26, 0x7e, 0x08, 0xcd, 0x7b, 0x2b, 0xbb, 0x22, 0xfd, 0xe6, 0xc8},
		storageTypes.MsgTryUpgrade,
		49,
	)

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.NotNil(t, decodeCtx.TryUpgrade)
}
