// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"strconv"

	"github.com/celenium-io/celestia-indexer/internal/pool"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/goccy/go-json"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

var (
	requestsPool = pool.New(func() []types.Request {
		return make([]types.Request, 0, 20)
	})
	blockResponsesPool = pool.New(func() []*types.Response[pkgTypes.ResultBlock] {
		return make([]*types.Response[pkgTypes.ResultBlock], 0, 20)
	})
	resultsResponsesPool = pool.New(func() []*types.Response[pkgTypes.ResultBlockResults] {
		return make([]*types.Response[pkgTypes.ResultBlockResults], 0, 20)
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

func (api *API) BlockBulkData(ctx context.Context, levels ...pkgTypes.Level) ([]pkgTypes.BlockData, error) {
	if len(levels) == 0 {
		return nil, nil
	}

	// Get slices from pools
	blockResponses := blockResponsesPool.Get()
	resultResponses := resultsResponsesPool.Get()
	requests := requestsPool.Get()

	// Ensure proper capacity
	requestsSize := len(levels) * 2
	if cap(blockResponses) < len(levels) {
		blockResponses = make([]*types.Response[pkgTypes.ResultBlock], 0, len(levels))
	}
	if cap(resultResponses) < len(levels) {
		resultResponses = make([]*types.Response[pkgTypes.ResultBlockResults], 0, len(levels))
	}
	if cap(requests) < requestsSize {
		requests = make([]types.Request, 0, requestsSize)
	}

	// Reset and resize to needed length
	blockResponses = blockResponses[:len(levels)]
	resultResponses = resultResponses[:len(levels)]
	requests = requests[:requestsSize]

	// Defer cleanup and return to pool
	defer func() {
		// Clear references to prevent memory leaks
		for i := range blockResponses {
			blockResponses[i] = nil
		}
		blockResponses = blockResponses[:0]
		blockResponsesPool.Put(blockResponses)

		for i := range resultResponses {
			resultResponses[i] = nil
		}
		resultResponses = resultResponses[:0]
		resultsResponsesPool.Put(resultResponses)

		// Clear request data
		requests = requests[:0]
		requestsPool.Put(requests)
	}()

	for i := range levels {
		blockResponses[i] = new(types.Response[pkgTypes.ResultBlock])
		resultResponses[i] = new(types.Response[pkgTypes.ResultBlockResults])

		levelString := levels[i].String()
		requests[i*2] = types.Request{
			Method:  pathBlock,
			JsonRpc: "2.0",
			Id:      int64(i) * 2,
			Params: []any{
				levelString,
			},
		}

		requests[i*2+1] = types.Request{
			Method:  pathBlockResults,
			JsonRpc: "2.0",
			Id:      int64(i)*2 + 1,
			Params: []any{
				levelString,
			},
		}
	}

	err := api.postStream(ctx, requests, func(dec *json.Decoder) error {
		if _, err := dec.Token(); err != nil {
			return err
		}
		for i := range levels {
			if err := dec.Decode(blockResponses[i]); err != nil {
				return err
			}
			if err := dec.Decode(resultResponses[i]); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "api.postStream")
	}

	blockData := make([]pkgTypes.BlockData, len(levels))
	for i := range levels {
		if blockResponses[i].Error != nil {
			return nil, errors.Wrapf(types.ErrRequest, "block error: %s", blockResponses[i].Error.Error())
		}
		if resultResponses[i].Error != nil {
			return nil, errors.Wrapf(types.ErrRequest, "results error: %s", resultResponses[i].Error.Error())
		}
		blockData[i].ResultBlock = blockResponses[i].Result
		blockData[i].ResultBlockResults = resultResponses[i].Result
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
