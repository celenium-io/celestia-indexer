// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

// MsgCreateValidator defines an SDK message for creating a new validator.
func MsgCreateValidator(ctx *context.Context, status storageTypes.Status, m *cosmosStakingTypes.MsgCreateValidator) (storageTypes.MsgType, []storage.AddressWithType, []string, error) {
	msgType := storageTypes.MsgCreateValidator
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height)
	if err != nil {
		return msgType, addresses, nil, err
	}

	validators := []string{m.ValidatorAddress}
	if status == storageTypes.StatusFailed {
		return msgType, addresses, nil, nil
	}

	var consAddress string
	if m.Pubkey != nil {
		pk, ok := m.Pubkey.GetCachedValue().(cryptotypes.PubKey)
		if ok {
			consAddress = pk.Address().String()
		} else {
			log.Warn().Msg("can't decode consensus address of validator")
		}
	}

	validatorAddress := types.Address(m.ValidatorAddress)
	_, b, err := validatorAddress.Decode()
	if err != nil {
		return msgType, addresses, nil, errors.Wrap(err, m.ValidatorAddress)
	}
	addr, err := types.NewAddressFromBytes(b)
	if err != nil {
		return msgType, addresses, nil, errors.Wrap(err, m.ValidatorAddress)
	}

	jailed := false
	validator := storage.Validator{
		Delegator:         addr.String(),
		Address:           m.ValidatorAddress,
		ConsAddress:       consAddress,
		Moniker:           m.Description.Moniker,
		Identity:          m.Description.Identity,
		Website:           m.Description.Website,
		Details:           m.Description.Details,
		Contacts:          m.Description.SecurityContact,
		Height:            ctx.Block.Height,
		Rate:              decimal.Zero,
		MaxRate:           decimal.Zero,
		MaxChangeRate:     decimal.Zero,
		MinSelfDelegation: decimal.Zero,
		Stake:             decimal.Zero,
		Jailed:            &jailed,
	}

	if !m.Value.IsNil() {
		amount := decimal.RequireFromString(m.Value.Amount.String())
		validator.Stake = amount

		address := storage.Address{
			Address: addr.String(),
			Balance: storage.Balance{
				Currency:  currency.DefaultCurrency,
				Spendable: decimal.Zero,
				Unbonding: decimal.Zero,
				Delegated: amount.Copy(),
			},
		}
		if err := ctx.AddAddress(&address); err != nil {
			return msgType, nil, validators, err
		}
		addresses[0].Balance = address.Balance

		ctx.AddDelegation(storage.Delegation{
			Address:   &address,
			Validator: &validator,
			Amount:    amount.Copy(),
		})

		ctx.AddStakingLog(storage.StakingLog{
			Time:      ctx.Block.Time,
			Height:    ctx.Block.Height,
			Address:   &address,
			Validator: &validator,
			Change:    amount.Copy(),
			Type:      storageTypes.StakingLogTypeDelegation,
		})
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

	ctx.AddValidator(validator)

	return msgType, addresses, validators, err
}

// MsgEditValidator defines a SDK message for editing an existing validator.
func MsgEditValidator(ctx *context.Context, status storageTypes.Status, m *cosmosStakingTypes.MsgEditValidator) (storageTypes.MsgType, []storage.AddressWithType, []string, error) {
	msgType := storageTypes.MsgEditValidator
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height)
	if err != nil {
		return msgType, addresses, nil, err
	}

	validators := []string{m.ValidatorAddress}
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
		Height:            ctx.Block.Height,
		Rate:              decimal.Zero,
		MinSelfDelegation: decimal.Zero,
		Stake:             decimal.Zero,
	}

	if m.CommissionRate != nil && !m.CommissionRate.IsNil() {
		validator.Rate = decimal.RequireFromString(m.CommissionRate.String())
	}
	if m.MinSelfDelegation != nil && !m.MinSelfDelegation.IsNil() {
		validator.MinSelfDelegation = decimal.RequireFromString(m.MinSelfDelegation.String())
	}
	ctx.AddValidator(validator)
	return msgType, addresses, validators, err
}

// MsgDelegate defines a SDK message for performing a delegation of coins
// from a delegator to a validator.
func MsgDelegate(ctx *context.Context, m *cosmosStakingTypes.MsgDelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgDelegate
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height)

	return msgType, addresses, err
}

// MsgBeginRedelegate defines an SDK message for performing a redelegation
// of coins from a delegator and source validator to a destination validator.
func MsgBeginRedelegate(ctx *context.Context, m *cosmosStakingTypes.MsgBeginRedelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgBeginRedelegate
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorSrc, address: m.ValidatorSrcAddress},
		{t: storageTypes.MsgAddressTypeValidatorDst, address: m.ValidatorDstAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgUndelegate defines a SDK message for performing an undelegation from a
// delegate and a validator.
func MsgUndelegate(ctx *context.Context, m *cosmosStakingTypes.MsgUndelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUndelegate
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgCancelUnbondingDelegation defines the SDK message for performing a cancel unbonding delegation for delegator
//
// Since: cosmos-sdk 0.46
func MsgCancelUnbondingDelegation(ctx *context.Context, m *cosmosStakingTypes.MsgCancelUnbondingDelegation) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCancelUnbondingDelegation
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeDelegator, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
