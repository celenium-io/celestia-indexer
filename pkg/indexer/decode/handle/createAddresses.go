// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
)

type addressData struct {
	t       storageTypes.MsgAddressType
	address string
}

type addressesData []addressData

func createAddresses(ctx *context.Context, data addressesData, level types.Level) ([]storage.AddressWithType, error) {
	addresses := make([]storage.AddressWithType, len(data))
	for i, d := range data {
		_, hash, err := types.Address(d.address).Decode()
		if err != nil {
			return nil, err
		}
		address := storage.Address{
			Hash:       hash,
			Height:     level,
			LastHeight: level,
			Address:    d.address,
			Balance:    storage.EmptyBalance(),
		}
		if err := ctx.AddAddress(&address); err != nil {
			return addresses, nil
		}

		addresses[i] = storage.AddressWithType{
			Type:    d.t,
			Address: address,
		}
	}
	return addresses, nil
}
