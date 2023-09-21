package handle

import (
	qgbTypes "github.com/celestiaorg/celestia-app/x/qgb/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

// MsgRegisterEVMAddress registers an evm address to a validator.
func MsgRegisterEVMAddress(level types.Level, m *qgbTypes.MsgRegisterEVMAddress) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRegisterEVMAddress
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeValidator, address: m.ValidatorAddress},
		// TODO: think about EVM addresses
	}, level)
	return msgType, addresses, err
}
