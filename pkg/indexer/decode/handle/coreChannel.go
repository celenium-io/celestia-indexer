// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	coreChannel "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
)

// MsgChannelOpenInit defines an sdk.Msg to initialize a channel handshake. It
// is called by a relayer on Chain A.
func MsgChannelOpenInit(level types.Level, m *coreChannel.MsgChannelOpenInit) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelOpenInit
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgChannelOpenTry defines a msg sent by a Relayer to try to open a channel
// on Chain B. The version field within the Channel field has been deprecated. Its
// value will be ignored by core IBC.
func MsgChannelOpenTry(level types.Level, m *coreChannel.MsgChannelOpenTry) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelOpenTry
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgChannelOpenAck defines a msg sent by a Relayer to Chain A to acknowledge
// the change of channel state to TRYOPEN on Chain B.
func MsgChannelOpenAck(level types.Level, m *coreChannel.MsgChannelOpenAck) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelOpenAck
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgChannelOpenConfirm defines a msg sent by a Relayer to Chain B to
// acknowledge the change of channel state to OPEN on Chain A.
func MsgChannelOpenConfirm(level types.Level, m *coreChannel.MsgChannelOpenConfirm) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelOpenConfirm
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgChannelCloseInit defines a msg sent by a Relayer to Chain A
// to close a channel with Chain B.
func MsgChannelCloseInit(level types.Level, m *coreChannel.MsgChannelCloseInit) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelCloseInit
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgChannelCloseConfirm defines a msg sent by a Relayer to Chain B
// to acknowledge the change of channel state to CLOSED on Chain A.
func MsgChannelCloseConfirm(level types.Level, m *coreChannel.MsgChannelCloseConfirm) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgChannelCloseConfirm
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgRecvPacket receives an incoming IBC packet
func MsgRecvPacket(level types.Level, m *coreChannel.MsgRecvPacket) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRecvPacket
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgTimeout receives a timed-out packet
func MsgTimeout(level types.Level, m *coreChannel.MsgTimeout) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTimeout
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgTimeoutOnClose timed-out packet upon counterparty channel closure
func MsgTimeoutOnClose(level types.Level, m *coreChannel.MsgTimeoutOnClose) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTimeoutOnClose
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}

// MsgAcknowledgement receives incoming IBC acknowledgement
func MsgAcknowledgement(level types.Level, m *coreChannel.MsgAcknowledgement) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgAcknowledgement
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)
	return msgType, addresses, err
}
