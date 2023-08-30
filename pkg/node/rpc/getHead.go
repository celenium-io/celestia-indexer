package rpc

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
)

func (api *API) GetHead(ctx context.Context) (types.ResultBlock, error) {
	return api.GetBlock(ctx, 0)
}
