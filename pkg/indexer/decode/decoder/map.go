// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decoder

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	channelTypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	tmTypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func DecimalFromMap(m map[string]string, key string) decimal.Decimal {
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

func CoinFromMap(m map[string]string, key string) (cosmosTypes.Coin, error) {
	str := StringFromMap(m, key)
	if str == "" {
		return cosmosTypes.Coin{}, nil
	}
	return cosmosTypes.ParseCoinNormalized(str)
}

func Map(m map[string]any, key string) (map[string]any, error) {
	val, ok := m[key]
	if !ok {
		return nil, errors.Errorf("can't find key: %s", key)
	}
	result, ok := val.(map[string]any)
	if !ok {
		return nil, errors.Errorf("value of '%s' is not map: %##v", key, val)
	}
	return result, nil
}

func StringFromMap(m map[string]string, key string) string {
	val, ok := m[key]
	if !ok {
		return ""
	}
	return val
}

func BalanceFromMap(m map[string]string, key string) (*cosmosTypes.Coin, error) {
	str := StringFromMap(m, key)
	if str == "" {
		return nil, nil
	}
	coin, err := cosmosTypes.ParseCoinNormalized(str)
	if err != nil {
		return nil, err
	}
	return &coin, nil
}

func AmountFromMap(m map[string]string, key string) decimal.Decimal {
	str := StringFromMap(m, key)
	if str == "" {
		return decimal.Zero
	}
	str = strings.TrimSuffix(str, currency.DefaultCurrency)
	return decimal.RequireFromString(str)
}

func TimeFromMap(m map[string]string, key string) (time.Time, error) {
	val, ok := m[key]
	if !ok {
		return time.Time{}, errors.Errorf("can't find key: %s", key)
	}
	return time.Parse(time.RFC3339, val)
}

var (
	nsDivider = decimal.NewFromInt(10).Pow(decimal.NewFromInt(9))
)

func UnixNanoFromMap(m map[string]string, key string) time.Time {
	value := DecimalFromMap(m, key)
	if value.IsZero() {
		return time.Time{}
	}
	x := value.Div(nsDivider)
	return time.Unix(x.IntPart(), value.Mod(nsDivider).IntPart()).UTC()
}

func Int64FromMap(m map[string]string, key string) (int64, error) {
	val, ok := m[key]
	if !ok {
		return 0, errors.Errorf("can't find key: %s", key)
	}
	return strconv.ParseInt(val, 10, 64)
}

func AuthMsgIndexFromMap(m map[string]string) (*int64, error) {
	val, ok := m["authz_msg_index"]
	if !ok {
		return nil, nil
	}
	i, err := strconv.ParseInt(val, 10, 64)
	return &i, err
}

func Uint64FromMap(m map[string]string, key string) (uint64, error) {
	val, ok := m[key]
	if !ok {
		return 0, errors.Errorf("can't find key: %s", key)
	}
	return strconv.ParseUint(val, 10, 64)
}

func Uint64(m map[string]any, key string) (uint64, error) {
	val, ok := m[key]
	if !ok {
		return 0, errors.Errorf("can't find key: %s", key)
	}
	switch v := val.(type) {
	case uint64:
		return v, nil
	case float64:
		return uint64(v), nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	}
	return 0, errors.Errorf("key '%s' is not a uint64", key)
}

func BoolFromMap(m map[string]string, key string) (bool, error) {
	val, ok := m[key]
	if !ok {
		return false, errors.Errorf("can't find key: %s", key)
	}
	return strconv.ParseBool(val)
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
	switch v := val.(type) {
	case channelTypes.Order:
		return v == channelTypes.ORDERED, nil
	case string:
		order, ok := channelTypes.Order_value[v]
		if !ok {
			return false, errors.Errorf("key '%s' has unknown Order value: %s", key, v)
		}
		return channelTypes.Order(order) == channelTypes.ORDERED, nil
	default:
		return false, errors.Errorf("key '%s' is not a Order", key)
	}
}

func RevisionHeightFromMap(m map[string]string, key string) (uint64, uint64, error) {
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

func MessagesFromMap(m map[string]any, key string) ([]cosmosTypes.Msg, error) {
	val, ok := m[key]
	if !ok {
		return nil, errors.Errorf("can't find key: %s", key)
	}
	msgs, ok := val.([]cosmosTypes.Msg)
	if !ok {
		return nil, errors.Errorf("key '%s' is not a messages", key)
	}
	return msgs, nil
}

func HyperlaneMessageFromMap(m map[string]string, key string) (*util.HyperlaneMessage, error) {
	str, ok := m[key]
	if !ok {
		return nil, nil
	}
	if str == "" {
		return nil, nil
	}

	unquoted, err := parseUnquoteOptional(str)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unquote hyperlane message")
	}

	messageBytes, err := util.DecodeEthHex(unquoted)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode hyperlane message")
	}

	result, err := util.ParseHyperlaneMessage(messageBytes)
	return &result, err
}

func parseUnquoteOptional(s string) (string, error) {
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return strconv.Unquote(s)
	}
	return s, nil
}

func BytesFromMap(m map[string]string, key string) ([]byte, error) {
	val, ok := m[key]
	if !ok {
		return nil, errors.Errorf("can't find key: %s", key)
	}
	str, err := parseUnquoteOptional(val)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unquote string")
	}
	return hex.DecodeString(strings.TrimPrefix(str, "0x"))
}
