// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"strconv"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	jxpkg "github.com/go-faster/jx"

	"github.com/pkg/errors"
)

const pathBlockResults = "block_results"

func (api *API) BlockResults(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlockResults, error) {
	args := make(map[string]string)
	if level != 0 {
		args["height"] = strconv.FormatUint(uint64(level), 10)
	}

	var result pkgTypes.ResultBlockResults
	err := api.getStream(ctx, pathBlockResults, args, func(d *jxpkg.Decoder) error {
		return jxResponse(d, func(d *jxpkg.Decoder) error {
			var err error
			result, err = jxResultBlockResults(d)
			return err
		})
	})
	return result, errors.Wrap(err, "BlockResults")
}
