// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

// MsgCreateVestingAccount defines a message that enables creating a vesting
// account.
func MsgCreateVestingAccount(level types.Level, m *cosmosVestingTypes.MsgCreateVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateVestingAccount
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}

// MsgCreatePermanentLockedAccount defines a message that enables creating a permanent
// locked account.
//
// Since: cosmos-sdk 0.46
func MsgCreatePermanentLockedAccount(level types.Level, m *cosmosVestingTypes.MsgCreatePermanentLockedAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreatePermanentLockedAccount
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}

// MsgCreateVestingAccount defines a message that enables creating a vesting
// account.
//
// Since: cosmos-sdk 0.46
func MsgCreatePeriodicVestingAccount(level types.Level, m *cosmosVestingTypes.MsgCreatePeriodicVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreatePeriodicVestingAccount
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}
