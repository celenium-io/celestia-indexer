// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decoder

import (
	"strconv"
	"strings"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/cosmos/cosmos-sdk/types"
	channelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	tmTypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"
	"github.com/pkg/errors"
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

func AmountFromMap(m map[string]any, key string) decimal.Decimal {
	str := StringFromMap(m, key)
	if str == "" {
		return decimal.Zero
	}
	str = strings.TrimSuffix(str, currency.DefaultCurrency)
	return decimal.RequireFromString(str)
}

func TimeFromMap(m map[string]any, key string) (time.Time, error) {
	val, ok := m[key]
	if !ok {
		return time.Time{}, errors.Errorf("can't find key: %s", key)
	}
	str, ok := val.(string)
	if !ok {
		return time.Time{}, errors.Errorf("key '%s' is not a string", key)
	}
	return time.Parse(time.RFC3339, str)
}

func UnixNanoFromMap(m map[string]any, key string) (time.Time, error) {
	u, err := Int64FromMap(m, key)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "Int64FromMap")
	}
	if u == 0 {
		return time.Time{}, nil
	}
	return time.Unix(0, u).UTC(), nil
}

func Int64FromMap(m map[string]any, key string) (int64, error) {
	val, ok := m[key]
	if !ok {
		return 0, errors.Errorf("can't find key: %s", key)
	}
	str, ok := val.(string)
	if !ok {
		return 0, errors.Errorf("key '%s' is not a string", key)
	}
	return strconv.ParseInt(str, 10, 64)
}

func AuthMsgIndexFromMap(m map[string]any) (*int64, error) {
	val, ok := m["authz_msg_index"]
	if !ok {
		return nil, nil
	}
	str, ok := val.(string)
	if !ok {
		return nil, errors.New("key 'auth_msg_index' is not a string")
	}
	i, err := strconv.ParseInt(str, 10, 64)
	return &i, err
}

func Uint64FromMap(m map[string]any, key string) (uint64, error) {
	val, ok := m[key]
	if !ok {
		return 0, errors.Errorf("can't find key: %s", key)
	}
	str, ok := val.(string)
	if !ok {
		return 0, errors.Errorf("key '%s' is not a string", key)
	}
	return strconv.ParseUint(str, 10, 64)
}

func ClientStateFromMap(m map[string]any, key string) (*tmTypes.ClientState, error) {
	val, ok := m[key]
	if !ok {
		return nil, errors.Errorf("can't find key: %s", key)
	}
	cs, ok := val.(tmTypes.ClientState)
	if !ok {
		return nil, errors.Errorf("key '%s' is not a client state", key)
	}
	return &cs, nil
}

func HeaderFromMap(m map[string]any, key string) (*tmTypes.Header, error) {
	val, ok := m[key]
	if !ok {
		return nil, errors.Errorf("can't find key: %s", key)
	}
	header, ok := val.(tmTypes.Header)
	if !ok {
		return nil, errors.Errorf("key '%s' is not a header", key)
	}
	return &header, nil
}

func ChannelOrderingFromMap(m map[string]any, key string) (bool, error) {
	val, ok := m[key]
	if !ok {
		return false, errors.Errorf("can't find key: %s", key)
	}
	order, ok := val.(channelTypes.Order)
	if !ok {
		return false, errors.Errorf("key '%s' is not a Order", key)
	}
	return order == channelTypes.ORDERED, nil
}

func RevisionHeightFromMap(m map[string]any, key string) (uint64, uint64, error) {
	ch := StringFromMap(m, key)
	parts := strings.Split(ch, "-")
	if len(parts) != 2 {
		return 0, 0, errors.Errorf("invalid revision height: %s", ch)
	}
	revision, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, 0, errors.Wrap(err, "revision")
	}

	height, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, 0, errors.Wrap(err, "height")
	}

	return revision, height, nil
}
