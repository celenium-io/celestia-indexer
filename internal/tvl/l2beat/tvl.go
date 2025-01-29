// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package l2beat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type TVLResponse struct {
	Data    Data   `json:"data"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type Result struct {
	Data Data `json:"data"`
}

type Data struct {
	Usd   float64 `json:"usdValue"`
	Eth   float64 `json:"ethValue"`
	Chart Chart   `json:"chart"`
}

type Chart struct {
	Types []string `json:"types"`
	Data  []Item   `json:"data"`
}

type Item struct {
	Time      time.Time
	Native    decimal.Decimal
	Canonical decimal.Decimal
	External  decimal.Decimal
	EthPrice  decimal.Decimal
}

func (item *Item) UnmarshalJSON(data []byte) error {
	var items []float64
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	if len(items) != 5 {
		return errors.Errorf("invalid chart item: %s", string(data))
	}
	item.Time = time.Unix(int64(items[0]), 0).UTC()
	item.Native = decimal.NewFromFloat(items[1])
	item.Canonical = decimal.NewFromFloat(items[2])
	item.External = decimal.NewFromFloat(items[3])
	item.EthPrice = decimal.NewFromFloat(items[4])
	return nil
}

func (api API) TVL(ctx context.Context, rollupName string, timeframe TvlTimeframe) (result TVLResponse, err error) {
	args := make(map[string]string)
	args["range"] = timeframe.String()

	if err = api.get(ctx, fmt.Sprintf("scaling/tvl/%s", rollupName), args, &result); err != nil {
		return
	}

	if !result.Success {
		err = errors.Errorf("l2 beat error: %s", result.Error)
	}
	return
}
