// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package binance

import (
	"context"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type OHLC struct {
	Time  time.Time
	Open  decimal.Decimal
	High  decimal.Decimal
	Low   decimal.Decimal
	Close decimal.Decimal
}

func (ohlc *OHLC) UnmarshalJSON(data []byte) error {
	var ts int64
	if err := json.Unmarshal(data, &[]any{
		&ts, &ohlc.Open, &ohlc.High, &ohlc.Low, &ohlc.Close,
	}); err != nil {
		return err
	}

	ohlc.Time = time.Unix(ts/1000, 0).UTC()
	return nil
}

type OHLCArgs struct {
	Start    int64
	End      int64
	Limit    int64
	TimeZone string
}

func (api API) OHLC(ctx context.Context, symbol, interval string, arguments *OHLCArgs) (candles []OHLC, err error) {
	args := map[string]string{
		"symbol":   symbol,
		"interval": interval,
	}
	if arguments != nil {
		if arguments.Start > 0 {
			args["startTime"] = strconv.FormatInt(arguments.Start, 10)
		}
		if arguments.End > 0 {
			args["endTime"] = strconv.FormatInt(arguments.End, 10)
		}
		if arguments.Limit > 0 {
			args["limit"] = strconv.FormatInt(arguments.Limit, 10)
		}
		if arguments.TimeZone != "" {
			args["timeZone"] = arguments.TimeZone
		}
	}
	err = api.get(ctx, "api/v3/klines", args, &candles)
	return
}
