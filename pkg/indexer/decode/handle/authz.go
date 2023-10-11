// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

// MsgGrant is a request type for Grant method. It declares authorization to the grantee
// on behalf of the granter with the provided expiration time.
func MsgGrant(level types.Level, m *authz.MsgGrant) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgGrant
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, level)
	return msgType, addresses, err
}

// MsgExec attempts to execute the provided messages using
// authorizations granted to the grantee. Each message should have only
// one signer corresponding to the granter of the authorization.
func MsgExec(level types.Level, m *authz.MsgExec) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgExec
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, level)

	// MsgExecute also has Msgs field, where also can be addresses.
	// Authorization Msg requests to execute. Each msg must implement Authorization interface
	// The x/authz will try to find a grant matching (msg.signers[0], grantee, MsgTypeURL(msg))
	// triple and validate it.

	return msgType, addresses, err
}

// MsgRevoke revokes any authorization with the provided sdk.Msg type on the
// granter's account with that has been granted to the grantee.
func MsgRevoke(level types.Level, m *authz.MsgRevoke) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRevoke
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, level)
	return msgType, addresses, err
}
