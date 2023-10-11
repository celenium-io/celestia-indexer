// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

// MsgSetWithdrawAddress sets the withdrawal address for
// a delegator (or validator self-delegation).
func MsgSetWithdrawAddress(level types.Level, m *cosmosDistributionTypes.MsgSetWithdrawAddress) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSetWithdrawAddress
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeWithdraw, address: m.WithdrawAddress},
	}, level)
	return msgType, addresses, err
}

// MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator
// from a single validator.
func MsgWithdrawDelegatorReward(level types.Level, m *cosmosDistributionTypes.MsgWithdrawDelegatorReward) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawDelegatorReward
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, level)

	return msgType, addresses, err
}

// MsgWithdrawValidatorCommission withdraws the full commission to the validator
// address.
func MsgWithdrawValidatorCommission(level types.Level, m *cosmosDistributionTypes.MsgWithdrawValidatorCommission) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawValidatorCommission
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

// MsgFundCommunityPool allows an account to directly
// fund the community pool.
func MsgFundCommunityPool(level types.Level, m *cosmosDistributionTypes.MsgFundCommunityPool) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgFundCommunityPool
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDepositor, address: m.Depositor},
	}, level)
	return msgType, addresses, err
}
