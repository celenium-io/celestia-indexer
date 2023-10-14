// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

// IBCTransfer defines a msg to transfer fungible tokens (i.e., Coins) between
// ICS20 enabled chains. See ICS Spec here:
// https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#data-structures
func IBCTransfer(level types.Level, m *ibcTypes.MsgTransfer) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.IBCTransfer
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSender, address: m.Sender},
		// {t: storageTypes.MsgAddressTypeReceiver,
		// address: m.Receiver}, // TODO: is it data to do IBC Transfer on cosmos network?
	}, level)
	return msgType, addresses, err
}
