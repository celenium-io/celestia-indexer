// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// MsgSetWithdrawAddress sets the withdrawal address for
// a delegator (or validator self-delegation).
func MsgSetWithdrawAddress(ctx *context.Context, m *cosmosDistributionTypes.MsgSetWithdrawAddress) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSetWithdrawAddress
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeWithdraw, address: m.WithdrawAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator
// from a single validator.
func MsgWithdrawDelegatorReward(ctx *context.Context, m *cosmosDistributionTypes.MsgWithdrawDelegatorReward) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawDelegatorReward
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height)

	return msgType, addresses, err
}

// MsgWithdrawValidatorCommission withdraws the full commission to the validator
// address.
func MsgWithdrawValidatorCommission(ctx *context.Context, m *cosmosDistributionTypes.MsgWithdrawValidatorCommission) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawValidatorCommission
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgFundCommunityPool allows an account to directly
// fund the community pool.
func MsgFundCommunityPool(ctx *context.Context, m *cosmosDistributionTypes.MsgFundCommunityPool) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgFundCommunityPool
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDepositor, address: m.Depositor},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
