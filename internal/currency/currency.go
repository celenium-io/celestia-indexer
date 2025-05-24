// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package currency

import "github.com/shopspring/decimal"

const (
	Utia string = "utia"
	Tia  string = "tia"
)

const (
	DefaultCurrency = "utia"
)

func StringTia(val decimal.Decimal) string {
	return val.StringFixed(6)
}

func StringUtia(val decimal.Decimal) string {
	return val.StringFixed(0)
}

var divider = decimal.NewFromInt(10).Pow(decimal.NewFromInt(-6))

func StringTiaFromUtia(val decimal.Decimal) string {
	return val.Mul(divider).StringFixed(6)
}
