// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package context

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test_AddSupply(t *testing.T) {
	tests := []struct {
		name string
		data map[string]any
		want decimal.Decimal
	}{
		{
			name: "valid amount",
			data: map[string]any{
				"amount": "1000000000000000000utia",
			},
			want: decimal.NewFromInt(1000000000000000000),
		}, {
			name: "valid amount but no utia",
			data: map[string]any{
				"amount": "1000000000000000000test",
			},
			want: decimal.NewFromInt(0),
		},
		{
			name: "invalid amount",
			data: map[string]any{
				"amount": "invalid_amount",
			},
			want: decimal.Zero,
		},
		{
			name: "amount without currency",
			data: map[string]any{
				"amount": "123456",
			},
			want: decimal.RequireFromString("123456"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			ctx.Block = &storage.Block{
				Stats: storage.BlockStats{
					SupplyChange: decimal.Zero,
				},
			}

			ctx.AddSupply(tt.data)
			require.Equal(t, tt.want.String(), ctx.Block.Stats.SupplyChange.String())
		})
	}
}

func Test_SubSupply(t *testing.T) {
	tests := []struct {
		name string
		data map[string]any
		want decimal.Decimal
	}{
		{
			name: "valid amount",
			data: map[string]any{
				"amount": "1000000000000000000utia",
			},
			want: decimal.NewFromInt(-1000000000000000000),
		}, {
			name: "valid amount but no utia",
			data: map[string]any{
				"amount": "1000000000000000000test",
			},
			want: decimal.NewFromInt(0),
		},
		{
			name: "invalid amount",
			data: map[string]any{
				"amount": "invalid_amount",
			},
			want: decimal.Zero,
		},
		{
			name: "amount without currency",
			data: map[string]any{
				"amount": "123456",
			},
			want: decimal.RequireFromString("-123456"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			ctx.Block = &storage.Block{
				Stats: storage.BlockStats{
					SupplyChange: decimal.Zero,
				},
			}

			ctx.SubSupply(tt.data)
			require.Equal(t, tt.want.String(), ctx.Block.Stats.SupplyChange.String())
		})
	}
}

func Test_AddAddress_New(t *testing.T) {
	ctx := NewContext()

	address := &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.NewFromInt(100),
			Delegated: decimal.NewFromInt(50),
			Unbonding: decimal.NewFromInt(20),
		},
	}

	err := ctx.AddAddress(address)
	require.NoError(t, err)

	addr, ok := ctx.Addresses.Get("celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e")
	require.True(t, ok)
	require.Equal(t, 1, ctx.Addresses.Len())
	require.Equal(t, address, addr)
}

func Test_AddAddress_Existing(t *testing.T) {
	ctx := NewContext()

	address := &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.NewFromInt(100),
			Delegated: decimal.NewFromInt(50),
			Unbonding: decimal.NewFromInt(20),
		},
	}

	err := ctx.AddAddress(address)
	require.NoError(t, err)

	addressUpdate := &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.NewFromInt(50),
			Delegated: decimal.NewFromInt(25),
			Unbonding: decimal.NewFromInt(10),
		},
	}

	err = ctx.AddAddress(addressUpdate)
	require.NoError(t, err)

	addr, ok := ctx.Addresses.Get("celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e")
	require.True(t, ok)
	require.Equal(t, 1, ctx.Addresses.Len())
	require.Equal(t, &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Hash:       []byte{0x23, 0x51, 0x2f, 0x62, 0x1d, 0x76, 0x1f, 0x6a, 0xa7, 0xc0, 0x47, 0x09, 0xa1, 0x36, 0x8c, 0x87, 0x81, 0x4c, 0x00, 0xd4},
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.NewFromInt(150),
			Delegated: decimal.NewFromInt(75),
			Unbonding: decimal.NewFromInt(30),
		},
	}, addr)
}

func Test_AddAddress_ExistingWithInvalidCurrency(t *testing.T) {
	ctx := NewContext()

	address := &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.NewFromInt(100),
			Delegated: decimal.NewFromInt(50),
			Unbonding: decimal.NewFromInt(20),
		},
	}

	err := ctx.AddAddress(address)
	require.NoError(t, err)

	addressUpdate := &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Balance: storage.Balance{
			Currency:  "invalid_currency",
			Spendable: decimal.NewFromInt(50),
			Delegated: decimal.NewFromInt(25),
			Unbonding: decimal.NewFromInt(10),
		},
	}

	err = ctx.AddAddress(addressUpdate)
	require.NoError(t, err)

	addr, ok := ctx.Addresses.Get("celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e")
	require.True(t, ok)
	require.Equal(t, 1, ctx.Addresses.Len())
	require.Equal(t, &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Hash:       []byte{0x23, 0x51, 0x2f, 0x62, 0x1d, 0x76, 0x1f, 0x6a, 0xa7, 0xc0, 0x47, 0x09, 0xa1, 0x36, 0x8c, 0x87, 0x81, 0x4c, 0x00, 0xd4},
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.NewFromInt(100),
			Delegated: decimal.NewFromInt(50),
			Unbonding: decimal.NewFromInt(20),
		},
	}, addr)
}
