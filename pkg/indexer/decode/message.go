package decode

import (
	"time"

	"github.com/celestiaorg/celestia-app/pkg/namespace"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
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
		d.Msg.Type, d.Addresses, err = handleMsgWithdrawValidatorCommission(height, typedMsg)
	case *cosmosDistributionTypes.MsgWithdrawDelegatorReward:
		d.Msg.Type, d.Addresses, err = handleMsgWithdrawDelegatorReward(height, typedMsg)
	case *cosmosStakingTypes.MsgEditValidator:
		d.Msg.Type, d.Addresses, err = handleMsgEditValidator(height, typedMsg)
	case *cosmosStakingTypes.MsgBeginRedelegate:
		d.Msg.Type, d.Addresses, err = handleMsgBeginRedelegate(height, typedMsg)
	case *cosmosStakingTypes.MsgCreateValidator:
		d.Msg.Type, d.Addresses, err = handleMsgCreateValidator(height, typedMsg)
	case *cosmosStakingTypes.MsgDelegate:
		d.Msg.Type, d.Addresses, err = handleMsgDelegate(height, typedMsg)
	case *cosmosStakingTypes.MsgUndelegate:
		d.Msg.Type, d.Addresses, err = handleMsgUndelegate(height, typedMsg)
	case *cosmosSlashingTypes.MsgUnjail:
		d.Msg.Type, d.Addresses, err = handleMsgUnjail(height, typedMsg)
	case *cosmosBankTypes.MsgSend:
		d.Msg.Type, d.Addresses, err = handleMsgSend(height, typedMsg)
	case *cosmosVestingTypes.MsgCreateVestingAccount:
		d.Msg.Type, d.Addresses, err = handleMsgCreateVestingAccount(height, typedMsg)
	case *cosmosVestingTypes.MsgCreatePeriodicVestingAccount:
		d.Msg.Type, d.Addresses, err = handleMsgCreatePeriodicVestingAccount(height, typedMsg)
	case *appBlobTypes.MsgPayForBlobs:
		d.Msg.Type, d.Addresses, d.Msg.Namespace, d.BlobsSize, err = handleMsgPayForBlobs(height, typedMsg)
	case *cosmosFeegrant.MsgGrantAllowance:
		d.Msg.Type, d.Addresses, err = handleMsgGrantAllowance(height, typedMsg)
	default:
		d.Msg.Type = storageTypes.MsgUnknown
	}

	if err != nil {
		err = errors.Wrapf(err, "while decoding msg(%T) on position=%d", msg, position)
	}

	return
}

type addressesData []struct {
	t       storageTypes.TxAddressType
	address string
}

func createAddresses(data addressesData, level types.Level) []storage.AddressWithType {
	addresses := make([]storage.AddressWithType, len(data))
	for i, d := range data {
		addresses[i] = storage.AddressWithType{
			Type: d.t,
			Address: storage.Address{
				Height:  level,
				Hash:    []byte(d.address),
				Balance: decimal.Zero,
			},
		}
	}
	return addresses
}

func handleMsgWithdrawValidatorCommission(level types.Level, m *cosmosDistributionTypes.MsgWithdrawValidatorCommission) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawValidatorCommission
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgWithdrawDelegatorReward(level types.Level, m *cosmosDistributionTypes.MsgWithdrawDelegatorReward) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawDelegatorReward
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)

	return msgType, addresses, nil
}

func handleMsgEditValidator(level types.Level, m *cosmosStakingTypes.MsgEditValidator) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgEditValidator
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgBeginRedelegate(level types.Level, m *cosmosStakingTypes.MsgBeginRedelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgBeginRedelegate
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorSrcAddress, address: m.ValidatorSrcAddress},
		{t: storageTypes.TxAddressTypeValidatorDstAddress, address: m.ValidatorDstAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgCreateValidator(level types.Level, m *cosmosStakingTypes.MsgCreateValidator) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateValidator
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgDelegate(level types.Level, m *cosmosStakingTypes.MsgDelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgDelegate
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgUndelegate(level types.Level, m *cosmosStakingTypes.MsgUndelegate) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUndelegate
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgUnjail(level types.Level, m *cosmosSlashingTypes.MsgUnjail) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUnjail
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddr},
	}, level)
	return msgType, addresses, nil
}

func handleMsgSend(level types.Level, m *cosmosBankTypes.MsgSend) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSend
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.TxAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgCreateVestingAccount(level types.Level, m *cosmosVestingTypes.MsgCreateVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateVestingAccount
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.TxAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgCreatePeriodicVestingAccount(level types.Level, m *cosmosVestingTypes.MsgCreatePeriodicVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreatePeriodicVestingAccount
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.TxAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, nil
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

	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeSigner, address: m.Signer},
	}, level)

	return storageTypes.MsgPayForBlobs, addresses, namespaces, blobsSize, nil
}

func handleMsgGrantAllowance(level types.Level, m *cosmosFeegrant.MsgGrantAllowance) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgGrantAllowance
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeGranter, address: m.Granter},
		{t: storageTypes.TxAddressTypeGrantee, address: m.Grantee},
	}, level)
	return msgType, addresses, nil
}
