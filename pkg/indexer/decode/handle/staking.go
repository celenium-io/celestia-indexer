// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

// MsgCreateValidator defines an SDK message for creating a new validator.
func MsgCreateValidator(level types.Level, status storageTypes.Status, m *cosmosStakingTypes.MsgCreateValidator) (storageTypes.MsgType, []storage.AddressWithType, *storage.Validator, error) {
	msgType := storageTypes.MsgCreateValidator
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, level)
	if status == storageTypes.StatusFailed {
		return msgType, addresses, nil, nil
	}

	validator := storage.Validator{
		Delegator:         m.DelegatorAddress,
		Address:           m.ValidatorAddress,
		Moniker:           m.Description.Moniker,
		Identity:          m.Description.Identity,
		Website:           m.Description.Website,
		Details:           m.Description.Details,
		Contacts:          m.Description.SecurityContact,
		Height:            level,
		Rate:              decimal.Zero,
		MaxRate:           decimal.Zero,
		MaxChangeRate:     decimal.Zero,
		MinSelfDelegation: decimal.Zero,
	}

	if !m.Commission.Rate.IsNil() {
		validator.Rate = decimal.RequireFromString(m.Commission.Rate.String())
	}

	if !m.Commission.MaxRate.IsNil() {
		validator.MaxRate = decimal.RequireFromString(m.Commission.MaxRate.String())
	}

	if !m.Commission.MaxChangeRate.IsNil() {
		validator.MaxChangeRate = decimal.RequireFromString(m.Commission.MaxChangeRate.String())
	}

	if !m.MinSelfDelegation.IsNil() {
		validator.MinSelfDelegation = decimal.RequireFromString(m.MinSelfDelegation.String())
	}

	return msgType, addresses, &validator, err
}

// MsgEditValidator defines a SDK message for editing an existing validator.
func MsgEditValidator(level types.Level, status storageTypes.Status, m *cosmosStakingTypes.MsgEditValidator) (storageTypes.MsgType, []storage.AddressWithType, *storage.Validator, error) {
	msgType := storageTypes.MsgEditValidator
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, level)
	if status == storageTypes.StatusFailed {
		return msgType, addresses, nil, nil
	}

	validator := storage.Validator{
		Address:           m.ValidatorAddress,
		Moniker:           m.Description.Moniker,
		Identity:          m.Description.Identity,
		Website:           m.Description.Website,
		Details:           m.Description.Details,
		Contacts:          m.Description.SecurityContact,
		Height:            level,
		Rate:              decimal.Zero,
		MinSelfDelegation: decimal.Zero,
	}

	if m.CommissionRate != nil && !m.CommissionRate.IsNil() {
		validator.Rate = decimal.RequireFromString(m.CommissionRate.String())
	}
	if m.MinSelfDelegation != nil && !m.MinSelfDelegation.IsNil() {
		validator.MinSelfDelegation = decimal.RequireFromString(m.MinSelfDelegation.String())
	}
	return msgType, addresses, &validator, err
}

// MsgDelegate defines a SDK message for performing a delegation of coins
// from a delegator to a validator.
func MsgDelegate(level types.Level, m *cosmosStakingTypes.MsgDelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgDelegate
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

// MsgBeginRedelegate defines an SDK message for performing a redelegation
// of coins from a delegator and source validator to a destination validator.
func MsgBeginRedelegate(level types.Level, m *cosmosStakingTypes.MsgBeginRedelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgBeginRedelegate
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorSrc, address: m.ValidatorSrcAddress},
		{t: storageTypes.MsgAddressTypeValidatorDst, address: m.ValidatorDstAddress},
	}, level)
	return msgType, addresses, err
}

// MsgUndelegate defines a SDK message for performing an undelegation from a
// delegate and a validator.
func MsgUndelegate(level types.Level, m *cosmosStakingTypes.MsgUndelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUndelegate
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

// MsgCancelUnbondingDelegation defines the SDK message for performing a cancel unbonding delegation for delegator
//
// Since: cosmos-sdk 0.46
func MsgCancelUnbondingDelegation(level types.Level, m *cosmosStakingTypes.MsgCancelUnbondingDelegation) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCancelUnbondingDelegation
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}
