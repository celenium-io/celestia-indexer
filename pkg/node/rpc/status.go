// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

const pathStatus = "status"

func (api *API) Status(ctx context.Context) (types.Status, error) {
	var sr types.Response[types.Status]
	if err := api.get(ctx, pathStatus, nil, &sr); err != nil {
		return sr.Result, errors.Wrap(err, "api.get")
	}

	if sr.Error != nil {
		return sr.Result, errors.Wrapf(types.ErrRequest, "status request %d error: %s", sr.Id, sr.Error.Error())
	}

	return sr.Result, nil
}
