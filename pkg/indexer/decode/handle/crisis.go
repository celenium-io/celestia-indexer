package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// MsgVerifyInvariant represents a message to verify a particular invariance.
func MsgVerifyInvariant(level types.Level, m *crisisTypes.MsgVerifyInvariant) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVerifyInvariant
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Sender},
	}, level)
	return msgType, addresses, err
}
