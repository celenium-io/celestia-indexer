package rpc

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"strconv"
)

const pathBlockResults = "block_results"

type getBlockResults struct {
	Result types.ResultBlockResults `json:"result"`
}

func (api *API) GetBlockResults(ctx context.Context, level storage.Level) (types.ResultBlockResults, error) {
	args := make(map[string]string)
	if level != 0 {
		args["height"] = strconv.FormatUint(uint64(level), 10)
	}

	var gbr getBlockResults
	if err := api.get(ctx, pathBlockResults, args, &gbr); err != nil {
		api.log.Err(err).Msg("node get block_results request")
		return types.ResultBlockResults{}, err
	}

	return gbr.Result, nil
}
