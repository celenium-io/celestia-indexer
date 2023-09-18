package handle

import (
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func MsgGrant(level types.Level, m *authz.MsgGrant) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgGrant
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, level)
	return msgType, addresses, err
}
