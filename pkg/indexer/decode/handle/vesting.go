// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// MsgCreateVestingAccount defines a message that enables creating a vesting
// account.
func MsgCreateVestingAccount(ctx *context.Context, m *cosmosVestingTypes.MsgCreateVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateVestingAccount
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgCreatePermanentLockedAccount defines a message that enables creating a permanent
// locked account.
//
// Since: cosmos-sdk 0.46
func MsgCreatePermanentLockedAccount(ctx *context.Context, m *cosmosVestingTypes.MsgCreatePermanentLockedAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreatePermanentLockedAccount
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgCreateVestingAccount defines a message that enables creating a vesting
// account.
//
// Since: cosmos-sdk 0.46
func MsgCreatePeriodicVestingAccount(ctx *context.Context, m *cosmosVestingTypes.MsgCreatePeriodicVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreatePeriodicVestingAccount
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
