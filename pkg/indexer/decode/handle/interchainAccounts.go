package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	interchainAccounts "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
)

// MsgRegisterInterchainAccount defines the payload for Msg/MsgRegisterInterchainAccount
func MsgRegisterInterchainAccount(level types.Level, m *interchainAccounts.MsgRegisterInterchainAccount) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRegisterInterchainAccount
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.Owner},
	}, level)
	return msgType, addresses, err
}

// MsgSendTx defines the payload for Msg/SendTx
func MsgSendTx(level types.Level, m *interchainAccounts.MsgSendTx) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSendTx
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeOwner, address: m.Owner},
	}, level)
	return msgType, addresses, err
}
