// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// MsgGrant is a request type for Grant method. It declares authorization to the grantee
// on behalf of the granter with the provided expiration time.
func MsgGrant(ctx *context.Context, m *authz.MsgGrant) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgGrant
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgExec attempts to execute the provided messages using
// authorizations granted to the grantee. Each message should have only
// one signer corresponding to the granter of the authorization.
func MsgExec(ctx *context.Context, status types.Status, m *authz.MsgExec) (storageTypes.MsgType, []storage.AddressWithType, []string, error) {
	msgType := storageTypes.MsgExec
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, ctx.Block.Height)

	// MsgExecute also has Msgs field, where also can be addresses.
	// Authorization Msg requests to execute. Each msg must implement Authorization interface
	// The x/authz will try to find a grant matching (msg.signers[0], grantee, MsgTypeURL(msg))
	// triple and validate it.

	if err != nil {
		return msgType, addresses, nil, err
	}

	if status == types.StatusFailed {
		return msgType, addresses, nil, nil
	}

	msgs := make([]string, len(m.Msgs))
	for i := range m.Msgs {
		msgs[i] = m.Msgs[i].TypeUrl
	}

	return msgType, addresses, msgs, nil
}

// MsgRevoke revokes any authorization with the provided sdk.Msg type on the
// granter's account with that has been granted to the grantee.
func MsgRevoke(ctx *context.Context, m *authz.MsgRevoke) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgRevoke
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
