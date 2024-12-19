// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package l2beat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type Data struct {
	Json [][]interface{} `json:"json"`
}

type Result struct {
	Data Data `json:"data"`
}

type TVLResponse []struct {
	Result Result `json:"result"`
}

type RequestData map[string]struct {
	Json JsonData `json:"json"`
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

func (api API) TVL(ctx context.Context, rollupName string, timeframe storage.TvlTimeframe) (result TVLResponse, err error) {
	data := RequestData{
		"0": {
			JsonData{
				Filter: Filter{
					Type:       "projects",
					ProjectIds: []string{rollupName},
				},
				Range:                   string(timeframe),
				ExcludeAssociatedTokens: false,
			},
		},
	}

	var dataString, err1 = json.Marshal(data)
	if err1 != nil {
		return nil, fmt.Errorf("serialization error: %v", err)
	}

	args := make(map[string]string)
	args["batch"] = "1"
	args["input"] = string(dataString)

	err = api.get(ctx, "trpc/tvl.chart", args, &result)
	return
}
