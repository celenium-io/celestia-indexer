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
	signal "github.com/celestiaorg/celestia-app/v6/x/signal/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMsgSignalVersion() signal.MsgSignalVersion {
	m := signal.NewMsgSignalVersion(
		"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x", 10,
	)

	return *m
}

func TestDecodeMsg_SuccessOnMsgSignalVersion(t *testing.T) {
	m := createMsgSignalVersion()
	blob, now := testsuite.EmptyBlock()
	position := 7

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, &m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, &m, position,
		storageTypes.MsgAddressTypeValidator,
		"celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		[]byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
		storageTypes.MsgSignalVersion,
		58,
	)

	validator := storage.EmptyValidator()
	validator.Address = m.ValidatorAddress
	validator.Version = 10
	msgExpected.SignalVersion = &storage.SignalVersion{
		Height:    blob.Height,
		Validator: &validator,
		Time:      blob.Block.Time,
		Version:   m.Version,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, addressesExpected, dm.Addresses)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgTryUpgrade(t *testing.T) {
	valAddress, err := types.AccAddressFromBech32("celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r")
	require.NoError(t, err)
	m := signal.NewMsgTryUpgrade(
		valAddress,
	)
	blob, now := testsuite.EmptyBlock()
	position := 7

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeSigner,
		"celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r",
		[]byte{0x50, 0xa1, 0xec, 0xc6, 0x67, 0x0c, 0x9a, 0x72, 0x1f, 0x26, 0x7e, 0x08, 0xcd, 0x7b, 0x2b, 0xbb, 0x22, 0xfd, 0xe6, 0xc8},
		storageTypes.MsgTryUpgrade,
		49,
	)

	msgExpected.Upgrade = &storage.Upgrade{
		Height: blob.Height,
		Signer: &storage.Address{
			Address: m.Signer,
		},
		Time: blob.Block.Time,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
