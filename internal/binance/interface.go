// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package binance

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IApi interface {
	OHLC(ctx context.Context, symbol, interval string, arguments *OHLCArgs) (candles []OHLC, err error)
}
