// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"strconv"

	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

const pathBlock = "block"

func (api *API) Block(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlock, error) {
	args := make(map[string]string)
	if level != 0 {
		args["height"] = strconv.FormatUint(uint64(level), 10)
	}

	var gbr types.Response[pkgTypes.ResultBlock]
	if err := api.get(ctx, pathBlock, args, &gbr); err != nil {
		return gbr.Result, errors.Wrap(err, "api.get")
	}

	if gbr.Error != nil {
		return gbr.Result, errors.Wrapf(types.ErrRequest, "request %d error: %s", gbr.Id, gbr.Error.Error())
	}

	return gbr.Result, nil
}

func (api *API) BlockData(ctx context.Context, level pkgTypes.Level) (pkgTypes.BlockData, error) {
	block := types.Response[pkgTypes.ResultBlock]{}
	results := types.Response[pkgTypes.ResultBlockResults]{}

	responses := []any{
		&block,
		&results,
	}

	levelString := level.String()
	requests := []types.Request{
		{
			Method:  pathBlock,
			JsonRpc: "2.0",
			Id:      -1,
			Params: []any{
				levelString,
			},
		}, {
			Method:  pathBlockResults,
			JsonRpc: "2.0",
			Id:      -1,
			Params: []any{
				levelString,
			},
		},
	}

	var blockData pkgTypes.BlockData

	if err := api.post(ctx, requests, &responses); err != nil {
		return blockData, errors.Wrap(err, "api.post")
	}

	if block.Error != nil {
		return blockData, errors.Wrapf(types.ErrRequest, "request error: %s", block.Error.Error())
	}

	if results.Error != nil {
		return blockData, errors.Wrapf(types.ErrRequest, "request error: %s", results.Error.Error())
	}

	blockData.ResultBlock = block.Result
	blockData.ResultBlockResults = results.Result
	return blockData, nil
}
