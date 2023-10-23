// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
)

type addressData struct {
	t       storageTypes.MsgAddressType
	address string
}

type addressesData []addressData

func createAddresses(data addressesData, level types.Level) ([]storage.AddressWithType, error) {
	addresses := make([]storage.AddressWithType, len(data))
	for i, d := range data {
		_, hash, err := types.Address(d.address).Decode()
		if err != nil {
			return nil, err
		}
		addresses[i] = storage.AddressWithType{
			Type: d.t,
			Address: storage.Address{
				Hash:       hash,
				Height:     level,
				LastHeight: level,
				Address:    d.address,
				Balance:    storage.EmptyBalance(),
			},
		}
	}
	return addresses, nil
}
