// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	coreClient "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

// MsgCreateClient defines a message to create an IBC client
func MsgCreateClient(level types.Level, m *coreClient.MsgCreateClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateClient
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgUpdateClient defines a sdk.Msg to update an IBC client state using the given header
func MsgUpdateClient(level types.Level, m *coreClient.MsgUpdateClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateClient
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgUpgradeClient defines a sdk.Msg to upgrade an IBC client to a new client state
func MsgUpgradeClient(level types.Level, m *coreClient.MsgUpgradeClient) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpgradeClient
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgSubmitMisbehaviour defines a sdk.Msg type that submits Evidence for light client misbehavior
func MsgSubmitMisbehaviour(level types.Level, m *coreClient.MsgSubmitMisbehaviour) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSubmitMisbehaviour
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}
