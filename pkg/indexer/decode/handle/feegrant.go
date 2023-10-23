// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// MsgGrantAllowance adds permission for Grantee to spend up to Allowance
// of fees from the account of Granter.
func MsgGrantAllowance(level types.Level, m *feegrant.MsgGrantAllowance) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgGrantAllowance
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, level)
	return msgType, addresses, err
}

// MsgRevokeAllowance removes any existing Allowance from Granter to Grantee.
func MsgRevokeAllowance(level types.Level, m *feegrant.MsgRevokeAllowance) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRevokeAllowance
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, level)
	return msgType, addresses, err
}
