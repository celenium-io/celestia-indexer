package handle

import (
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func MsgWithdrawValidatorCommission(level types.Level, m *cosmosDistributionTypes.MsgWithdrawValidatorCommission) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawValidatorCommission
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

func MsgWithdrawDelegatorReward(level types.Level, m *cosmosDistributionTypes.MsgWithdrawDelegatorReward) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawDelegatorReward
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)

	return msgType, addresses, err
}

func MsgSetWithdrawalAddress(level types.Level, m *cosmosDistributionTypes.MsgSetWithdrawAddress) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSetWithdrawAddress
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.MsgAddressTypeWithdraw, address: m.WithdrawAddress},
	}, level)
	return msgType, addresses, err
}
