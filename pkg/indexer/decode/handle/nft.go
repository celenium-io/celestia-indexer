package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

// MsgSendNFT represents a message to send a nft from one account to another account.
func MsgSendNFT(level types.Level, m *nft.MsgSend) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSendNFT
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Sender},
		{t: storageTypes.MsgAddressTypeReceiver, address: m.Receiver},
	}, level)
	return msgType, addresses, err
}
