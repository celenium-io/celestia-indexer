package rpc

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"strconv"
)

const path = "block"

type getBlockResult struct {
	Result types.ResultBlock `json:"result"`
}

func (api *API) GetBlock(ctx context.Context, level storage.Level) (types.ResultBlock, error) {
	args := make(map[string]string)
	if level != 0 {
		args["height"] = strconv.FormatUint(uint64(level), 10)
	}

	var gbr getBlockResult
	if err := api.get(ctx, path, args, &gbr); err != nil {
		api.log.Err(err).Msg("node get block request")
		return types.ResultBlock{}, err
	}

	return gbr.Result, nil
}
