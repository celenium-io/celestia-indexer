package handle

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func MsgVote(level types.Level, voterAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVote
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, level)
	return msgType, addresses, err
}

func MsgVoteWeighted(level types.Level, voterAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVoteWeighted
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, level)
	return msgType, addresses, err
}

func MsgSubmitProposal(level types.Level, proposerAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSubmitProposal
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeProposer, address: proposerAddress},
	}, level)
	return msgType, addresses, err
}
