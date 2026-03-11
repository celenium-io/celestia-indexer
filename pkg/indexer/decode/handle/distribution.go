// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// MsgSetWithdrawAddress sets the withdrawal address for
// a delegator (or validator self-delegation).
func MsgSetWithdrawAddress(ctx *context.Context, msgId uint64, m *cosmosDistributionTypes.MsgSetWithdrawAddress) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSetWithdrawAddress
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeWithdraw, address: m.WithdrawAddress},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator
// from a single validator.
func MsgWithdrawDelegatorReward(ctx *context.Context, msgId uint64, m *cosmosDistributionTypes.MsgWithdrawDelegatorReward) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgWithdrawDelegatorReward
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height, msgId)

	return msgType, err
}

// MsgWithdrawValidatorCommission withdraws the full commission to the validator
// address.
func MsgWithdrawValidatorCommission(ctx *context.Context, msgId uint64, m *cosmosDistributionTypes.MsgWithdrawValidatorCommission) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgWithdrawValidatorCommission
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgFundCommunityPool allows an account to directly
// fund the community pool.
func MsgFundCommunityPool(ctx *context.Context, msgId uint64, m *cosmosDistributionTypes.MsgFundCommunityPool) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgFundCommunityPool
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDepositor, address: m.Depositor},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

func MsgUpdateParamsDistr(ctx *context.Context, msgId uint64, m *cosmosDistributionTypes.MsgUpdateParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
