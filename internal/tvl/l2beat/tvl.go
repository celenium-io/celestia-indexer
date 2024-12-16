package l2beat

import (
	"context"
)

type Data struct {
	JSON [][]interface{} `json:"json"`
}

type Result struct {
	Data Data `json:"data"`
}

type TVLResponse []struct {
	Result Result `json:"result"`
}

type TVLArgs struct {
	Batch int64    `json:"batch"`
	Input JsonData `json:"input"`
}

type JsonData struct {
	Filter                  Filter `json:"filter"`
	Range                   string `json:"range"`
	ExcludeAssociatedTokens bool   `json:"excludeAssociatedTokens"`
}

type Filter struct {
	Type       string   `json:"type"`
	ProjectIds []string `json:"projectIds"`
}

func (api API) TVL(ctx context.Context, arguments *TVLArgs) (result TVLResponse, err error) {
	// TODO: cast args to query params
	args := map[string]string{}
	err = api.get(ctx, "trpc/tvl.chart", args, &result)
	return
}
