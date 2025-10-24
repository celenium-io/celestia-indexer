// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package math

import "github.com/shopspring/decimal"

var shareDivider = decimal.NewFromInt(1_000_000)

func Shares(stake decimal.Decimal) decimal.Decimal {
	return stake.Div(shareDivider).Floor()
}
