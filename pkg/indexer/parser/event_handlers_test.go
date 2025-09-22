// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

var (
	testAddress     = "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"
	testHashAddress = []byte{0x96, 0xa, 0xa0, 0x36, 0x6b, 0x25, 0x4e, 0x1e, 0xa7, 0x9b, 0xda, 0x46, 0x7e, 0xb3, 0xaa, 0x5c, 0x97, 0xcb, 0xa5, 0xae}
	testIgpId       = []uint8{0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}
	testBlock       = storage.Block{
		Height: pkgTypes.Level(1488),
		Time:   time.Now(),
	}
)

func Test_parseCoinSpent(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		height  pkgTypes.Level
		want    *storage.Address
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]any{
				"spender": testAddress,
				"amount":  "123utia",
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balance: storage.Balance{
					Currency:  currency.DefaultCurrency,
					Spendable: decimal.RequireFromString("-123"),
					Delegated: decimal.Zero,
					Unbonding: decimal.Zero,
				},
			},
		}, {
			name: "test 2",
			data: map[string]any{
				"spender": testAddress,
				"amount":  nil,
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balance:    storage.EmptyBalance(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			err := parseCoinSpent(ctx, tt.data, tt.height)
			require.True(t, (err == nil) != tt.wantErr)
			require.EqualValues(t, 1, ctx.Addresses.Len())
			_ = ctx.Addresses.Range(func(_ string, value *storage.Address) (error, bool) {
				require.Equal(t, tt.want, value)
				return nil, false
			})
		})
	}
}

func Test_parseCreateIgp(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &testBlock

	tests := []struct {
		name    string
		data    map[string]any
		want    *storage.HLIGP
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]any{
				"denom":  "\"utia\"",
				"igp_id": "\"0x726f757465725f706f73745f6469737061746368000000040000000000000001\"",
				"owner":  testAddress,
			},
			want: &storage.HLIGP{
				Height: pkgTypes.Level(1488),
				Time:   testBlock.Time,
				Owner: &storage.Address{
					Address: testAddress,
				},
				IgpId: testIgpId,
				Denom: "utia",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseCreateIgp(ctx, tt.data)
			require.True(t, (err == nil) != tt.wantErr)
			require.EqualValues(t, 1, ctx.Igps.Len())
			_ = ctx.Igps.Range(func(_ string, value *storage.HLIGP) (error, bool) {
				require.Equal(t, tt.want, value)
				return nil, false
			})
		})
	}
}

func Test_parseSetDestinationGasConfig(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &testBlock

	tests := []struct {
		name    string
		data    map[string]any
		want    storage.HLIGPConfig
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]any{
				"gas_overhead":        "\"200000\"",
				"gas_price":           "\"1\"",
				"igp_id":              "\"0x726f757465725f706f73745f6469737061746368000000040000000000000001\"",
				"owner":               testAddress,
				"remote_domain":       "84532",
				"token_exchange_rate": "\"10000000000\"",
			},
			want: storage.HLIGPConfig{
				Height:            pkgTypes.Level(1488),
				Time:              testBlock.Time,
				IgpId:             testIgpId,
				GasPrice:          decimal.RequireFromString("1"),
				GasOverhead:       decimal.RequireFromString("200000"),
				RemoteDomain:      84532,
				TokenExchangeRate: "10000000000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseSetDestinationGasConfig(ctx, tt.data)
			require.True(t, (err == nil) != tt.wantErr)
			require.EqualValues(t, 1, len(ctx.IgpConfigs))
			for _, v := range ctx.IgpConfigs {
				require.Equal(t, tt.want, v)
			}
		})
	}
}

func Test_parseSetIgp(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &testBlock

	tests := []struct {
		name    string
		data    map[string]any
		want    *storage.HLIGP
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]any{
				"igp_id":             "\"0x726f757465725f706f73745f6469737061746368000000040000000000000001\"",
				"owner":              testAddress,
				"new_owner":          "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w61",
				"renounce_ownership": "false",
			},
			want: &storage.HLIGP{
				Height: pkgTypes.Level(1488),
				Time:   testBlock.Time,
				Owner: &storage.Address{
					Address: "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w61",
				},
				IgpId: testIgpId,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseSetIgp(ctx, tt.data)
			require.True(t, (err == nil) != tt.wantErr)
			require.EqualValues(t, 1, ctx.Igps.Len())
			_ = ctx.Igps.Range(func(_ string, value *storage.HLIGP) (error, bool) {
				require.Equal(t, tt.want, value)
				return nil, false
			})
		})
	}
}
