// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"strconv"

	"github.com/celenium-io/celestia-indexer/internal/pool"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	jxpkg "github.com/go-faster/jx"

	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

var requestsPool = pool.New(func() []types.Request {
	return make([]types.Request, 0, 20)
})

const pathBlock = "block"

func (api *API) Block(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlock, error) {
	args := make(map[string]string)
	if level != 0 {
		args["height"] = strconv.FormatInt(int64(level), 10)
	}

	var result pkgTypes.ResultBlock
	err := api.getStream(ctx, pathBlock, args, func(d *jxpkg.Decoder) error {
		return jxResponse(d, func(d *jxpkg.Decoder) error {
			var err error
			result, err = jxResultBlock(d)
			return err
		})
	})
	return result, errors.Wrap(err, "Block")
}

func (api *API) BlockBulkData(ctx context.Context, levels ...pkgTypes.Level) ([]pkgTypes.BlockData, error) {
	if len(levels) == 0 {
		return nil, nil
	}

	requests := requestsPool.Get()
	neededSize := len(levels) * 2
	if cap(requests) < neededSize {
		requests = make([]types.Request, 0, neededSize)
	}
	requests = requests[:neededSize]
	defer func() {
		requests = requests[:0]
		requestsPool.Put(requests)
	}()

	for i := range levels {
		levelString := levels[i].String()
		requests[i*2] = types.Request{Method: pathBlock, JsonRpc: "2.0", Id: -1, Params: []any{levelString}}
		requests[i*2+1] = types.Request{Method: pathBlockResults, JsonRpc: "2.0", Id: -1, Params: []any{levelString}}
	}

	result := make([]pkgTypes.BlockData, 0, len(levels))
	err := api.postStream(ctx, requests, func(d *jxpkg.Decoder) error {
		return jxBatchResponse(d, func(bd pkgTypes.BlockData) error {
			result = append(result, bd)
			return nil
		})
	})
	return result, errors.Wrap(err, "BlockBulkData")
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

func (api *API) BlockBulkDataStream(
	ctx context.Context,
	fn func(pkgTypes.BlockData) error,
	levels ...pkgTypes.Level,
) error {
	if len(levels) == 0 {
		return nil
	}
	// Get slices from pools
	requests := requestsPool.Get()

	// Ensure proper capacity
	neededSize := len(levels) * 2
	if cap(requests) < neededSize {
		requests = make([]types.Request, 0, neededSize)
	}
	// Reset and resize to needed length
	requests = requests[:neededSize]

	// Defer cleanup and return to pool
	defer func() {
		// Clear request data
		requests = requests[:0]
		requestsPool.Put(requests)
	}()

	for i := range levels {
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

	return api.postStream(ctx, requests, func(d *jxpkg.Decoder) error {
		return jxBatchResponse(d, fn)
	})
}
