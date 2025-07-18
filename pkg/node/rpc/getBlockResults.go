// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"strconv"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

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
