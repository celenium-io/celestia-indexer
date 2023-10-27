package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// MsgSoftwareUpgrade is the Msg/SoftwareUpgrade request type.
func MsgSoftwareUpgrade(level types.Level, m *upgrade.MsgSoftwareUpgrade) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSoftwareUpgrade
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, level)
	return msgType, addresses, err
}

// MsgSoftwareUpgrade is the Msg/SoftwareUpgrade request type.
func MsgCancelUpgrade(level types.Level, m *upgrade.MsgCancelUpgrade) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCancelUpgrade
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, level)
	return msgType, addresses, err
}
