package handle

import (
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

func MsgEditValidator(level types.Level, status storageTypes.Status, m *cosmosStakingTypes.MsgEditValidator) (storageTypes.MsgType, []storage.AddressWithType, *storage.Validator, error) {
	msgType := storageTypes.MsgEditValidator
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
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
		Height:            uint64(level),
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

func MsgBeginRedelegate(level types.Level, m *cosmosStakingTypes.MsgBeginRedelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgBeginRedelegate
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorSrcAddress, address: m.ValidatorSrcAddress},
		{t: storageTypes.MsgAddressTypeValidatorDstAddress, address: m.ValidatorDstAddress},
	}, level)
	return msgType, addresses, err
}

func MsgCreateValidator(level types.Level, status storageTypes.Status, m *cosmosStakingTypes.MsgCreateValidator) (storageTypes.MsgType, []storage.AddressWithType, *storage.Validator, error) {
	msgType := storageTypes.MsgCreateValidator
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
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
		Height:            uint64(level),
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

func MsgDelegate(level types.Level, m *cosmosStakingTypes.MsgDelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgDelegate
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

func MsgUndelegate(level types.Level, m *cosmosStakingTypes.MsgUndelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUndelegate
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}
