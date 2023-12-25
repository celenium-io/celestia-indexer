// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package binance

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestOHLC_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want OHLC
	}{
		{
			name: "test 1",
			data: []byte(`[
				1499040000000,
				"0.01634790",
				"0.80000000",
				"0.01575800",
				"0.01577100",
				"148976.11427815",
				1499644799999,
				"2434.19055334",
				308,
				"1756.87402397",
				"28.46694368",
				"0"
			]`),
			want: OHLC{
				Open:  decimal.RequireFromString("0.01634790"),
				High:  decimal.RequireFromString("0.80000000"),
				Low:   decimal.RequireFromString("0.01575800"),
				Close: decimal.RequireFromString("0.01577100"),
				Time:  time.Date(2017, 7, 3, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ohlc OHLC
			err := ohlc.UnmarshalJSON(tt.data)
			require.NoError(t, err)

			require.EqualValues(t, tt.want.Open.String(), ohlc.Open.String())
			require.EqualValues(t, tt.want.High.String(), ohlc.High.String())
			require.EqualValues(t, tt.want.Low.String(), ohlc.Low.String())
			require.EqualValues(t, tt.want.Close.String(), ohlc.Close.String())
			require.EqualValues(t, tt.want.Time.UTC().String(), ohlc.Time.String())
		})
	}
}
