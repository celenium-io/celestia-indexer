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

func createAddresses(ctx *context.Context, data addressesData, level types.Level, msgId uint64) error {
	for i := range data {
		_, hash, err := types.Address(data[i].address).Decode()
		if err != nil {
			return err
		}
		address := storage.Address{
			Hash:       hash,
			Height:     level,
			LastHeight: level,
			Address:    data[i].address,
			Balance:    storage.EmptyBalance(),
		}
		if err := ctx.AddAddress(&address); err != nil {
			return nil
		}

		ctx.AddAddressMessage(&storage.MsgAddress{
			MsgId:   msgId,
			Type:    data[i].t,
			Address: &address,
		})
	}
	return nil
}
