// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	fee "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
)

// MsgRegisterPayee defines the request type for the RegisterPayee rpc
func MsgRegisterPayee(level types.Level, m *fee.MsgRegisterPayee) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRegisterPayee
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeRelayer, address: m.Relayer},
		{t: storageTypes.MsgAddressTypePayee, address: m.Payee},
	}, level)
	return msgType, addresses, err
}

// MsgRegisterCounterpartyPayee defines the request type for the RegisterCounterpartyPayee rpc
func MsgRegisterCounterpartyPayee(level types.Level, m *fee.MsgRegisterCounterpartyPayee) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRegisterCounterpartyPayee
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeRelayer, address: m.Relayer},
		{t: storageTypes.MsgAddressTypePayee, address: m.CounterpartyPayee}, // the counterparty payee address
	}, level)
	return msgType, addresses, err
}

// MsgPayPacketFee defines the request type for the PayPacketFee rpc
// This Msg can be used to pay for a packet at the next sequence send & should be combined with the Msg that will be
// paid for
func MsgPayPacketFee(level types.Level, m *fee.MsgPayPacketFee) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgPayPacketFee
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgPayPacketFeeAsync defines the request type for the PayPacketFeeAsync rpc
// This Msg can be used to pay for a packet at a specified sequence (instead of the next sequence sends)
func MsgPayPacketFeeAsync() (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgPayPacketFeeAsync
	return msgType, []storage.AddressWithType{}, nil
}
