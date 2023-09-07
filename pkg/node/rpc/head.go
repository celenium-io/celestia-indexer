package rpc

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func (api *API) Head(ctx context.Context) (types.ResultBlock, error) {
	return api.Block(ctx, 0)
}
