// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	qgbTypes "github.com/celestiaorg/celestia-app/x/qgb/types"
)

// MsgRegisterEVMAddress registers an evm address to a validator.
func MsgRegisterEVMAddress(level types.Level, m *qgbTypes.MsgRegisterEVMAddress) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRegisterEVMAddress
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
		// TODO: think about EVM addresses
	}, level)
	return msgType, addresses, err
}
