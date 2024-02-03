// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package currency

import "github.com/shopspring/decimal"

type Denom string

const (
	Utia Denom = "utia"
	Tia  Denom = "tia"
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
