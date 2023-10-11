// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func DecimalFromMap(m map[string]any, key string) decimal.Decimal {
	str := StringFromMap(m, key)
	if str == "" {
		return decimal.Zero
	}
	str = strings.TrimRight(str, letters)
	dec, err := decimal.NewFromString(str)
	if err != nil {
		return decimal.Zero
	}
	return dec
}

func Amount(m map[string]any) decimal.Decimal {
	return DecimalFromMap(m, "amount")
}

func StringFromMap(m map[string]any, key string) string {
	val, ok := m[key]
	if !ok {
		return ""
	}
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

func BalanceFromMap(m map[string]any, key string) (*types.Coin, error) {
	str := StringFromMap(m, key)
	if str == "" {
		return nil, nil
	}
	coin, err := types.ParseCoinNormalized(str)
	if err != nil {
		return nil, err
	}
	return &coin, nil
}
