// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"

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
