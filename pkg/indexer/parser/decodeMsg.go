package parser

import (
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

type decodedMsg struct {
	msg       storage.Message
	blobsSize uint64
	addresses []storage.AddressWithType
}

func decodeMsg(b types.BlockData, msg cosmosTypes.Msg, position int) (d decodedMsg, err error) {
	d.msg.Height = b.Height
	d.msg.Time = b.Block.Time
	d.msg.Position = uint64(position)
	d.msg.Data = structs.Map(msg)

	switch typedMsg := msg.(type) {
	case *cosmosDistributionTypes.MsgWithdrawValidatorCommission:
		d.msg.Type, d.addresses, err = handleMsgWithdrawValidatorCommission(b.Height, typedMsg)
	case *cosmosDistributionTypes.MsgWithdrawDelegatorReward:
		d.msg.Type, d.addresses, err = handleMsgWithdrawDelegatorReward(b.Height, typedMsg)
	case *cosmosStakingTypes.MsgEditValidator:
		d.msg.Type, d.addresses, err = handleMsgEditValidator(b.Height, typedMsg)
	case *cosmosStakingTypes.MsgBeginRedelegate:
		d.msg.Type, d.addresses, err = handleMsgBeginRedelegate(b.Height, typedMsg)
	case *cosmosStakingTypes.MsgCreateValidator:
		d.msg.Type, d.addresses, err = handleMsgCreateValidator(b.Height, typedMsg)
	case *cosmosStakingTypes.MsgDelegate:
		d.msg.Type, d.addresses, err = handleMsgDelegate(b.Height, typedMsg)
	case *cosmosStakingTypes.MsgUndelegate:
		d.msg.Type, d.addresses, err = handleMsgUndelegate(b.Height, typedMsg)
	case *cosmosSlashingTypes.MsgUnjail:
		d.msg.Type, d.addresses, err = handleMsgUnjail(b.Height, typedMsg)
	case *cosmosBankTypes.MsgSend:
		d.msg.Type, d.addresses, err = handleMsgSend(b.Height, typedMsg)
	case *cosmosVestingTypes.MsgCreateVestingAccount:
		d.msg.Type, d.addresses, err = handleMsgCreateVestingAccount(b.Height, typedMsg)
	case *cosmosVestingTypes.MsgCreatePeriodicVestingAccount:
		d.msg.Type, d.addresses, err = handleMsgCreatePeriodicVestingAccount(b.Height, typedMsg)
	case *appBlobTypes.MsgPayForBlobs:
		d.msg.Type, d.addresses, d.msg.Namespace, d.blobsSize, err = handleMsgPayForBlobs(b.Height, typedMsg)
	case *cosmosFeegrant.MsgGrantAllowance:
		d.msg.Type, d.addresses, err = handleMsgGrantAllowance(b.Height, typedMsg)
	default:
		d.msg.Type = storageTypes.MsgUnknown
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
