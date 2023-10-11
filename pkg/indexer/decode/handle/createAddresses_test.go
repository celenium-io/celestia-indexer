// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAddresses_SingleAddress(t *testing.T) {
	data := addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts"},
	}
	level := types.Level(235236)

	addresses, err := createAddresses(data, level)
	assert.NoError(t, err)
	assert.NotEmpty(t, addresses)
	assert.Len(t, addresses, 1)

	addr := addresses[0]
	expectedAddr := storage.AddressWithType{
		Type: storageTypes.MsgAddressTypeVoter,
		Address: storage.Address{
			Hash:       []byte{8, 204, 180, 93, 112, 144, 218, 230, 174, 203, 58, 172, 76, 199, 190, 39, 45, 188, 116, 154},
			Height:     types.Level(235236),
			LastHeight: types.Level(235236),
			Address:    "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
			Balance: storage.Balance{
				Id:    0,
				Total: decimal.Zero,
			},
		},
	}
	assert.Equal(t, expectedAddr, addr)
}

func TestCreateAddresses_ListOfAddresses(t *testing.T) {
	data := addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37"},
		{t: storageTypes.MsgAddressTypeValidator, address: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym"},
	}
	level := types.Level(235236)

	addresses, err := createAddresses(data, level)
	assert.NoError(t, err)
	assert.NotEmpty(t, addresses)
	assert.Len(t, addresses, 2)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegator,
			Address: storage.Address{
				Id:         0,
				Height:     types.Level(235236),
				LastHeight: types.Level(235236),
				Address:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:       []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     types.Level(235236),
				LastHeight: types.Level(235236),
				Address:    "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
				Hash:       []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}
	assert.Equal(t, addressesExpected, addresses)
}

func TestCreateAddresses_ErrorOnDecodingAddress(t *testing.T) {
	data := addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: "NO_WAY_celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts"},
	}
	level := types.Level(235236)

	addresses, err := createAddresses(data, level)
	assert.Error(t, err, "decoding bech32 failed: string not all lowercase or all uppercase")
	assert.Empty(t, addresses)
}
