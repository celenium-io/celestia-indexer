package handle

import (
	circuitTypes "cosmossdk.io/x/circuit/types"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

func MsgAuthorizeCircuitBreaker(ctx *context.Context, m *circuitTypes.MsgAuthorizeCircuitBreaker) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgAuthorizeCircuitBreaker
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

func MsgResetCircuitBreaker(ctx *context.Context, m *circuitTypes.MsgResetCircuitBreaker) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgResetCircuitBreaker
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

func MsgTripCircuitBreaker(ctx *context.Context, m *circuitTypes.MsgTripCircuitBreaker) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTripCircuitBreaker
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
