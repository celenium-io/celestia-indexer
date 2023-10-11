// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

// MsgUnjail defines the Msg/Unjail request type
func MsgUnjail(level types.Level, m *cosmosSlashingTypes.MsgUnjail) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUnjail
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddr},
	}, level)
	return msgType, addresses, err
}
