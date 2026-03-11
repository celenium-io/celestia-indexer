// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestCreateAddresses_SingleAddress(t *testing.T) {
	data := addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts"},
	}
	level := types.Level(235236)
	ctx := context.NewContext()

	const msgId uint64 = 1
	err := createAddresses(ctx, data, level, msgId)
	require.NoError(t, err)

	require.Equal(t, 1, ctx.Addresses.Len())
	require.Equal(t, 1, ctx.AddressMessages.Len())
}

func TestCreateAddresses_ListOfAddresses(t *testing.T) {
	data := addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37"},
		{t: storageTypes.MsgAddressTypeValidator, address: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym"},
	}
	level := types.Level(235236)
	ctx := context.NewContext()

	const msgId uint64 = 2
	err := createAddresses(ctx, data, level, msgId)
	require.NoError(t, err)

	require.Equal(t, 2, ctx.Addresses.Len())
	require.Equal(t, 2, ctx.AddressMessages.Len())
}

func TestCreateAddresses_ErrorOnDecodingAddress(t *testing.T) {
	data := addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: "NO_WAY_celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts"},
	}
	level := types.Level(235236)
	ctx := context.NewContext()

	const msgId uint64 = 3
	err := createAddresses(ctx, data, level, msgId)
	require.Error(t, err, "decoding bech32 failed: string not all lowercase or all uppercase")
	require.Equal(t, 0, ctx.Addresses.Len())
}
