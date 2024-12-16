package lama

import (
	"context"
	"net/url"
)

type TVLResponse struct {
	Date int64   `json:"date"`
	TVL  float64 `json:"tvl"`
}

type TVLArgs struct {
	Chain string `json:"chain"`
}

func (api API) TVL(ctx context.Context, arguments *TVLArgs) (result []TVLResponse, err error) {
	path, err := url.JoinPath("v2/historicalChainTvl", arguments.Chain)
	if err != nil {
		return nil, err
	}
	err = api.get(ctx, path, nil, &result)
	return
}
