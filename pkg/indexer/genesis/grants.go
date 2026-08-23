// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"encoding/json"

	"cosmossdk.io/x/feegrant"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/handle"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/pkg/errors"
)

func (module *Module) parseFeeGrants(
	ctx *decodeContext.Context, feeGrantsRaw []json.RawMessage,
) error {
	for i := range feeGrantsRaw {
		var feeGrant feegrant.MsgGrantAllowance
		if err := module.codec.UnmarshalJSON(feeGrantsRaw[i], &feeGrant); err != nil {
			return err
		}
		if _, err := handle.MsgGrantAllowance(ctx, storageTypes.StatusSuccess, 0, &feeGrant); err != nil {
			return err
		}
	}
	return nil
}

func (module *Module) parseAuthzGrants(
	ctx *decodeContext.Context, grantsRaw []json.RawMessage,
) error {
	for i := range grantsRaw {
		var grant authz.GrantAuthorization
		if err := module.codec.UnmarshalJSON(grantsRaw[i], &grant); err != nil {
			return err
		}
		if grant.Authorization == nil {
			return errors.Errorf("genesis authz grant has empty authorization: granter=%s grantee=%s",
				grant.Granter, grant.Grantee)
		}

		msg := authz.MsgGrant{
			Granter: grant.Granter,
			Grantee: grant.Grantee,
			Grant: authz.Grant{
				Authorization: grant.Authorization,
				Expiration:    grant.Expiration,
			},
		}
		if _, err := handle.MsgGrant(ctx, storageTypes.StatusSuccess, 0, &msg); err != nil {
			return err
		}
	}
	return nil
}
