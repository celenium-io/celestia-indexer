// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/fatih/structs"
)

func MsgPruneExpiredGrants(ctx *context.Context, m *authz.MsgPruneExpiredGrants) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgPruneExpiredGrants
	addresses, err := createAddresses(ctx, addressesData{
		{t: types.MsgAddressTypeExecutor, address: m.Pruner},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

// MsgGrant is a request type for Grant method. It declares authorization to the grantee
// on behalf of the granter with the provided expiration time.
func MsgGrant(ctx *context.Context, status types.Status, m *authz.MsgGrant) (types.MsgType, []storage.AddressWithType, []storage.Grant, error) {
	msgType := types.MsgGrant
	addresses, err := createAddresses(ctx, addressesData{
		{t: types.MsgAddressTypeGranter, address: m.Granter},
		{t: types.MsgAddressTypeGrantee, address: m.Grantee},
	}, ctx.Block.Height)
	if err != nil || status == types.StatusFailed {
		return msgType, addresses, nil, err
	}
	grants, err := parseGrants(m, ctx.Block.Time, ctx.Block.Height)
	return msgType, addresses, grants, err
}

// MsgExec attempts to execute the provided messages using
// authorizations granted to the grantee. Each message should have only
// one signer corresponding to the granter of the authorization.
func MsgExec(ctx *context.Context, status types.Status, m *authz.MsgExec) (types.MsgType, []storage.AddressWithType, []string, error) {
	msgType := types.MsgExec
	addresses, err := createAddresses(ctx, addressesData{
		{t: types.MsgAddressTypeGrantee, address: m.Grantee},
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
func MsgRevoke(ctx *context.Context, status types.Status, m *authz.MsgRevoke) (types.MsgType, []storage.AddressWithType, []storage.Grant, error) {
	msgType := types.MsgRevoke
	addresses, err := createAddresses(ctx, addressesData{
		{t: types.MsgAddressTypeGranter, address: m.Granter},
		{t: types.MsgAddressTypeGrantee, address: m.Grantee},
	}, ctx.Block.Height)
	if err != nil || status == types.StatusFailed {
		return msgType, addresses, nil, err
	}

	grant := storage.Grant{
		Granter: &storage.Address{
			Address: m.Granter,
		},
		Grantee: &storage.Address{
			Address: m.Grantee,
		},
		Revoked:       true,
		Authorization: m.MsgTypeUrl,
		RevokeHeight:  &ctx.Block.Height,
	}

	return msgType, addresses, []storage.Grant{grant}, nil
}

func parseGrants(msg *authz.MsgGrant, t time.Time, height pkgTypes.Level) ([]storage.Grant, error) {
	if msg == nil {
		return nil, nil
	}

	switch msg.Grant.Authorization.TypeUrl {
	case "/cosmos.authz.v1beta1.GenericAuthorization":
		var typ authz.GenericAuthorization
		if err := typ.Unmarshal(msg.Grant.Authorization.Value); err != nil {
			return nil, err
		}
		return []storage.Grant{
			{
				Params:        structs.Map(typ),
				Authorization: typ.MsgTypeURL(),
				Granter: &storage.Address{
					Address: msg.Granter,
				},
				Grantee: &storage.Address{
					Address: msg.Grantee,
				},
				Height:     height,
				Expiration: msg.Grant.Expiration,
				Time:       t,
			},
		}, nil
	case "/cosmos.bank.v1beta1.SendAuthorization":
		var typ bankTypes.SendAuthorization
		if err := typ.Unmarshal(msg.Grant.Authorization.Value); err != nil {
			return nil, err
		}
		return []storage.Grant{
			{
				Params:        structs.Map(typ),
				Authorization: "/cosmos.bank.v1beta1.MsgSend",
				Granter: &storage.Address{
					Address: msg.Granter,
				},
				Grantee: &storage.Address{
					Address: msg.Grantee,
				},
				Height:     height,
				Expiration: msg.Grant.Expiration,
				Time:       t,
			},
		}, nil
	case "/cosmos.staking.v1beta1.StakeAuthorization":
		var typ stakingTypes.StakeAuthorization
		if err := typ.Unmarshal(msg.Grant.Authorization.Value); err != nil {
			return nil, err
		}
		switch typ.AuthorizationType {
		case stakingTypes.AuthorizationType_AUTHORIZATION_TYPE_DELEGATE:
			return []storage.Grant{
				{
					Params:        structs.Map(typ),
					Authorization: "/cosmos.staking.v1beta1.MsgDelegate",
					Granter: &storage.Address{
						Address: msg.Granter,
					},
					Grantee: &storage.Address{
						Address: msg.Grantee,
					},
					Height:     height,
					Expiration: msg.Grant.Expiration,
					Time:       t,
				},
			}, nil
		case stakingTypes.AuthorizationType_AUTHORIZATION_TYPE_REDELEGATE:
			return []storage.Grant{
				{
					Params:        structs.Map(typ),
					Authorization: "/cosmos.staking.v1beta1.MsgRedelegate",
					Granter: &storage.Address{
						Address: msg.Granter,
					},
					Grantee: &storage.Address{
						Address: msg.Grantee,
					},
					Height:     height,
					Expiration: msg.Grant.Expiration,
					Time:       t,
				},
			}, nil
		case stakingTypes.AuthorizationType_AUTHORIZATION_TYPE_UNDELEGATE:
			return []storage.Grant{
				{
					Params:        structs.Map(typ),
					Authorization: "/cosmos.staking.v1beta1.MsgUndelegate",
					Granter: &storage.Address{
						Address: msg.Granter,
					},
					Grantee: &storage.Address{
						Address: msg.Grantee,
					},
					Height:     height,
					Expiration: msg.Grant.Expiration,
					Time:       t,
				},
			}, nil
		case stakingTypes.AuthorizationType_AUTHORIZATION_TYPE_UNSPECIFIED:
			return []storage.Grant{
				{
					Params:        structs.Map(typ),
					Authorization: "/cosmos.staking.v1beta1.MsgDelegate",
					Granter: &storage.Address{
						Address: msg.Granter,
					},
					Grantee: &storage.Address{
						Address: msg.Grantee,
					},
					Height:     height,
					Expiration: msg.Grant.Expiration,
					Time:       t,
				}, {
					Params:        structs.Map(typ),
					Authorization: "/cosmos.staking.v1beta1.MsgRedelegate",
					Granter: &storage.Address{
						Address: msg.Granter,
					},
					Grantee: &storage.Address{
						Address: msg.Grantee,
					},
					Height:     height,
					Expiration: msg.Grant.Expiration,
					Time:       t,
				}, {
					Params:        structs.Map(typ),
					Authorization: "/cosmos.staking.v1beta1.MsgUndelegate",
					Granter: &storage.Address{
						Address: msg.Granter,
					},
					Grantee: &storage.Address{
						Address: msg.Grantee,
					},
					Height:     height,
					Expiration: msg.Grant.Expiration,
					Time:       t,
				},
			}, nil
		}
	}
	return nil, nil
}
