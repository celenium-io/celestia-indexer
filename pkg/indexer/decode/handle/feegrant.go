// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"cosmossdk.io/x/feegrant"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/fatih/structs"
)

// MsgGrantAllowance adds permission for Grantee to spend up to Allowance
// of fees from the account of Granter.
func MsgGrantAllowance(ctx *context.Context, status storageTypes.Status, msgId uint64, m *feegrant.MsgGrantAllowance) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgGrantAllowance
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	grant := storage.Grant{
		Granter: &storage.Address{
			Address: m.Granter,
		},
		Grantee: &storage.Address{
			Address: m.Grantee,
		},
		Authorization: "fee",
		Height:        ctx.Block.Height,
		Time:          ctx.Block.Time,
	}

	err = parseGrantFee(m, &grant)
	ctx.AddGrants(&grant)
	return msgType, err
}

// MsgRevokeAllowance removes any existing Allowance from Granter to Grantee.
func MsgRevokeAllowance(ctx *context.Context, status storageTypes.Status, msgId uint64, m *feegrant.MsgRevokeAllowance) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgRevokeAllowance
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeGranter, address: m.Granter},
		{t: storageTypes.MsgAddressTypeGrantee, address: m.Grantee},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	grant := storage.Grant{
		Granter: &storage.Address{
			Address: m.Granter,
		},
		Grantee: &storage.Address{
			Address: m.Grantee,
		},
		Revoked:       true,
		Authorization: "fee",
		RevokeHeight:  &ctx.Block.Height,
	}

	ctx.AddGrants(&grant)
	return msgType, nil
}

func parseGrantFee(m *feegrant.MsgGrantAllowance, g *storage.Grant) error {
	switch m.Allowance.TypeUrl {
	case "/cosmos.feegrant.v1beta1.BasicAllowance":
		var body feegrant.BasicAllowance
		if err := body.Unmarshal(m.Allowance.Value); err != nil {
			return err
		}
		g.Params = structs.Map(body)
		g.Expiration = body.Expiration
	case "/cosmos.feegrant.v1beta1.PeriodicAllowance":
		var body feegrant.PeriodicAllowance
		if err := body.Unmarshal(m.Allowance.Value); err != nil {
			return err
		}
		g.Params = structs.Map(body)
		g.Expiration = body.Basic.Expiration
	case "/cosmos.feegrant.v1beta1.AllowedMsgAllowance":
		var body feegrant.AllowedMsgAllowance
		if err := body.Unmarshal(m.Allowance.Value); err != nil {
			return err
		}

		var basic feegrant.BasicAllowance
		if err := basic.Unmarshal(body.Allowance.Value); err != nil {
			return err
		}
		g.Params = structs.Map(body)
		g.Params["Allowance"] = basic
		g.Expiration = basic.Expiration
	}
	return nil
}
