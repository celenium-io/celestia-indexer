package decode

import (
	"github.com/rs/zerolog/log"
	"time"

	"github.com/celestiaorg/celestia-app/pkg/namespace"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
	qgbTypes "github.com/celestiaorg/celestia-app/x/qgb/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	cosmosFeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type DecodedMsg struct {
	Msg       storage.Message
	BlobsSize uint64
	Addresses []storage.AddressWithType
}

func Message(msg cosmosTypes.Msg, height types.Level, time time.Time, position int) (d DecodedMsg, err error) {
	d.Msg.Height = height
	d.Msg.Time = time
	d.Msg.Position = uint64(position)
	d.Msg.Data = structs.Map(msg)

	switch typedMsg := msg.(type) {
	case *cosmosDistributionTypes.MsgWithdrawValidatorCommission:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgWithdrawValidatorCommission(height, typedMsg)
	case *cosmosDistributionTypes.MsgWithdrawDelegatorReward:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgWithdrawDelegatorReward(height, typedMsg)
	case *cosmosStakingTypes.MsgEditValidator:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Validator, err = handleMsgEditValidator(height, typedMsg)
	case *cosmosStakingTypes.MsgBeginRedelegate:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgBeginRedelegate(height, typedMsg)
	case *cosmosStakingTypes.MsgCreateValidator:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Validator, err = handleMsgCreateValidator(height, typedMsg)
	case *cosmosStakingTypes.MsgDelegate:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgDelegate(height, typedMsg)
	case *cosmosStakingTypes.MsgUndelegate:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgUndelegate(height, typedMsg)
	case *cosmosSlashingTypes.MsgUnjail:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgUnjail(height, typedMsg)
	case *cosmosBankTypes.MsgSend:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgSend(height, typedMsg)
	case *cosmosVestingTypes.MsgCreateVestingAccount:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgCreateVestingAccount(height, typedMsg)
	case *cosmosVestingTypes.MsgCreatePeriodicVestingAccount:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgCreatePeriodicVestingAccount(height, typedMsg)
	case *appBlobTypes.MsgPayForBlobs:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Namespace, d.BlobsSize, err = handleMsgPayForBlobs(height, typedMsg)
	case *cosmosFeegrant.MsgGrantAllowance:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgGrantAllowance(height, typedMsg)
	case *qgbTypes.MsgRegisterEVMAddress:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgRegisterEVMAddress(height, typedMsg)
	case *cosmosDistributionTypes.MsgSetWithdrawAddress:
		d.Msg.Type, d.Msg.Addresses, err = handleMsgSetWithdrawalAddress(height, typedMsg)
	default:
		log.Err(errors.New("unknown message type")).Msgf("got type %T", msg)
		d.Msg.Type = storageTypes.MsgUnknown
	}

	if err != nil {
		err = errors.Wrapf(err, "while decoding msg(%T) on position=%d", msg, position)
	}

	d.Addresses = append(d.Addresses, d.Msg.Addresses...)
	return
}

type addressesData []struct {
	t       storageTypes.MsgAddressType
	address string
}

func createAddresses(data addressesData, level types.Level) ([]storage.AddressWithType, error) {
	addresses := make([]storage.AddressWithType, len(data))
	for i, d := range data {
		_, hash, err := types.Address(d.address).Decode()
		if err != nil {
			return nil, err
		}
		addresses[i] = storage.AddressWithType{
			Type: d.t,
			Address: storage.Address{
				Hash:    hash,
				Height:  level,
				Address: d.address,
				Balance: storage.Balance{
					Total: decimal.Zero,
				},
			},
		}
	}
	return addresses, nil
}

func handleMsgWithdrawValidatorCommission(level types.Level, m *cosmosDistributionTypes.MsgWithdrawValidatorCommission) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawValidatorCommission
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

func handleMsgWithdrawDelegatorReward(level types.Level, m *cosmosDistributionTypes.MsgWithdrawDelegatorReward) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawDelegatorReward
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)

	return msgType, addresses, err
}

func handleMsgEditValidator(level types.Level, m *cosmosStakingTypes.MsgEditValidator) (storageTypes.MsgType, []storage.AddressWithType, *storage.Validator, error) {
	msgType := storageTypes.MsgEditValidator
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
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

func handleMsgBeginRedelegate(level types.Level, m *cosmosStakingTypes.MsgBeginRedelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgBeginRedelegate
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorSrcAddress, address: m.ValidatorSrcAddress},
		{t: storageTypes.MsgAddressTypeValidatorDstAddress, address: m.ValidatorDstAddress},
	}, level)
	return msgType, addresses, err
}

func handleMsgCreateValidator(level types.Level, m *cosmosStakingTypes.MsgCreateValidator) (storageTypes.MsgType, []storage.AddressWithType, *storage.Validator, error) {
	msgType := storageTypes.MsgCreateValidator
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
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

func handleMsgDelegate(level types.Level, m *cosmosStakingTypes.MsgDelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgDelegate
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

func handleMsgUndelegate(level types.Level, m *cosmosStakingTypes.MsgUndelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUndelegate
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

func handleMsgUnjail(level types.Level, m *cosmosSlashingTypes.MsgUnjail) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUnjail
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddr},
	}, level)
	return msgType, addresses, err
}

func handleMsgSend(level types.Level, m *cosmosBankTypes.MsgSend) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSend
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}

func handleMsgCreateVestingAccount(level types.Level, m *cosmosVestingTypes.MsgCreateVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateVestingAccount
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}

func handleMsgCreatePeriodicVestingAccount(level types.Level, m *cosmosVestingTypes.MsgCreatePeriodicVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreatePeriodicVestingAccount
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}

func handleMsgPayForBlobs(level types.Level, m *appBlobTypes.MsgPayForBlobs) (storageTypes.MsgType, []storage.AddressWithType, []storage.Namespace, uint64, error) {
	var blobsSize uint64
	namespaces := make([]storage.Namespace, len(m.Namespaces))

	for nsI, ns := range m.Namespaces {
		if len(m.BlobSizes) < nsI {
			return storageTypes.MsgUnknown, nil, nil, 0, errors.Errorf(
				"blob sizes length=%d is less then namespaces index=%d", len(m.BlobSizes), nsI)
		}

		appNS := namespace.Namespace{Version: ns[0], ID: ns[1:]}
		size := uint64(m.BlobSizes[nsI])
		blobsSize += size
		namespaces[nsI] = storage.Namespace{
			FirstHeight: level,
			Version:     appNS.Version,
			NamespaceID: appNS.ID,
			Size:        size,
			PfbCount:    1,
			Reserved:    appNS.IsReserved(),
		}
	}

	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)

	return storageTypes.MsgPayForBlobs, addresses, namespaces, blobsSize, err
}

func handleMsgGrantAllowance(level types.Level, m *cosmosFeegrant.MsgGrantAllowance) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgGrantAllowance
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, level)
	return msgType, addresses, err
}

func handleMsgRegisterEVMAddress(level types.Level, m *qgbTypes.MsgRegisterEVMAddress) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRegisterEVMAddress
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
		// TODO: think about EVM addresses
	}, level)
	return msgType, addresses, err
}

func handleMsgSetWithdrawalAddress(level types.Level, m *cosmosDistributionTypes.MsgSetWithdrawAddress) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSetWithdrawAddress
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeWithdraw, address: m.WithdrawAddress},
	}, level)
	return msgType, addresses, err
}
