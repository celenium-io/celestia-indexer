// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package lama

import (
	"context"
	"net/url"
)

type TVLResponse struct {
	Date int64   `json:"date"`
	TVL  float64 `json:"tvl"`
}

func (api API) TVL(ctx context.Context, rollupName string) (result []TVLResponse, err error) {
	path, err := url.JoinPath("v2/historicalChainTvl", rollupName)
	if err != nil {
		return nil, err
	}
	err = api.get(ctx, path, nil, &result)
	return
}
