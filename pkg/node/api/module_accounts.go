// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package api

import (
	"context"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

func (api *API) ModuleAccounts(ctx context.Context) ([]types.Account, error) {
	var response types.Auth
	if err := api.get(ctx, "cosmos/auth/v1beta1/module_accounts", nil, &response); err != nil {
		return nil, errors.Wrap(err, "get")
	}
	return response.Accounts, nil
}
