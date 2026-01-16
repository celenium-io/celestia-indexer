// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"strconv"

	"github.com/celenium-io/celestia-indexer/internal/pool"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

var (
	requestsPool = pool.New(func() []types.Request {
		return make([]types.Request, 0, 20)
	})
	responsesPool = pool.New(func() []any {
		return make([]any, 0, 20)
	})
)

const pathBlock = "block"

func (api *API) Block(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlock, error) {
	args := make(map[string]string)
	if level != 0 {
		args["height"] = strconv.FormatInt(int64(level), 10)
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

func (api *API) BlockBulkData(ctx context.Context, levels ...pkgTypes.Level) ([]pkgTypes.BlockData, error) {
	if len(levels) == 0 {
		return nil, nil
	}

	// Get slices from pools
	responses := responsesPool.Get()
	requests := requestsPool.Get()

	// Ensure proper capacity
	neededSize := len(levels) * 2
	if cap(responses) < neededSize {
		responses = make([]any, 0, neededSize)
	}
	if cap(requests) < neededSize {
		requests = make([]types.Request, 0, neededSize)
	}

	// Reset and resize to needed length
	responses = responses[:neededSize]
	requests = requests[:neededSize]

	// Defer cleanup and return to pool
	defer func() {
		// Clear references to prevent memory leaks
		for i := range responses {
			responses[i] = nil
		}
		responses = responses[:0]
		responsesPool.Put(responses)

		// Clear request data
		requests = requests[:0]
		requestsPool.Put(requests)
	}()

	for i := range levels {
		responses[i*2] = &types.Response[pkgTypes.ResultBlock]{}
		responses[i*2+1] = &types.Response[pkgTypes.ResultBlockResults]{}

		levelString := levels[i].String()
		requests[i*2] = types.Request{
			Method:  pathBlock,
			JsonRpc: "2.0",
			Id:      -1,
			Params: []any{
				levelString,
			},
		}
		requests[i*2+1] = types.Request{
			Method:  pathBlockResults,
			JsonRpc: "2.0",
			Id:      -1,
			Params: []any{
				levelString,
			},
		}
	}

	if err := api.post(ctx, requests, &responses); err != nil {
		return nil, errors.Wrap(err, "api.post")
	}

	var blockData = make([]pkgTypes.BlockData, len(levels))

	for i := range responses {
		switch typ := responses[i].(type) {
		case *types.Response[pkgTypes.ResultBlock]:
			if typ.Error != nil {
				return nil, errors.Wrapf(types.ErrRequest, "request error: %s", typ.Error.Error())
			}
			blockData[i/2].ResultBlock = typ.Result
		case *types.Response[pkgTypes.ResultBlockResults]:
			if typ.Error != nil {
				return nil, errors.Wrapf(types.ErrRequest, "request error: %s", typ.Error.Error())
			}
			blockData[i/2].ResultBlockResults = typ.Result
		}
	}

	return blockData, nil
}

func (api *API) BlockDataGet(ctx context.Context, level pkgTypes.Level) (pkgTypes.BlockData, error) {
	var blockData pkgTypes.BlockData

	block, err := api.Block(ctx, level)
	if err != nil {
		return blockData, errors.Wrapf(types.ErrRequest, "request error: %s", err.Error())
	}

	results, err := api.BlockResults(ctx, level)
	if err != nil {
		return blockData, errors.Wrapf(types.ErrRequest, "request error: %s", err.Error())
	}

	blockData.ResultBlock = block
	blockData.ResultBlockResults = results
	return blockData, nil
}
