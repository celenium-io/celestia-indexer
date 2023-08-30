package rpc

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
)

type GenesisResult struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		types.Genesis `json:"genesis"`
	} `json:"result"`
}

func (api *API) GetGenesis(ctx context.Context) (types.Genesis, error) {
	path := "genesis"

	var gr GenesisResult
	if err := api.get(ctx, path, nil, &gr); err != nil {
		api.log.Err(err).Msg("genesis block request")
		return types.Genesis{}, err
	}

	return gr.Result.Genesis, nil
}
