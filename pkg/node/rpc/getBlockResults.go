// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"strconv"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

const pathBlockResults = "block_results"

func (api *API) BlockResults(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlockResults, error) {
	args := make(map[string]string)
	if level != 0 {
		args["height"] = strconv.FormatUint(uint64(level), 10)
	}

	var gbr types.Response[pkgTypes.ResultBlockResults]
	if err := api.get(ctx, pathBlockResults, args, &gbr); err != nil {
		return gbr.Result, errors.Wrap(err, "api.get")
	}

	if gbr.Error != nil {
		return gbr.Result, errors.Wrapf(types.ErrRequest, "request %d error: %s", gbr.Id, gbr.Error.Error())
	}

	return gbr.Result, nil
}
