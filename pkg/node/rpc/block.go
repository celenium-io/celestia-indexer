package rpc

import (
	"context"
	"strconv"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

const pathBlock = "block"

func (api *API) Block(ctx context.Context, level storage.Level) (types.ResultBlock, error) {
	args := make(map[string]string)
	if level != 0 {
		args["height"] = strconv.FormatUint(uint64(level), 10)
	}

	var gbr types.Response[types.ResultBlock]
	if err := api.get(ctx, pathBlock, args, &gbr); err != nil {
		return gbr.Result, errors.Wrap(err, "api.get")
	}

	if gbr.Error != nil {
		return gbr.Result, errors.Wrapf(types.ErrRequest, "request %d error: %s", gbr.Id, gbr.Error.Error())
	}

	return gbr.Result, nil
}
