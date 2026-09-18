// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package context

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test_AddSupply(t *testing.T) {
	tests := []struct {
		name string
		data map[string]string
		want decimal.Decimal
	}{
		{
			name: "valid amount",
			data: map[string]string{
				"amount": "1000000000000000000utia",
			},
			want: decimal.NewFromInt(1000000000000000000),
		}, {
			name: "valid amount but no utia",
			data: map[string]string{
				"amount": "1000000000000000000test",
			},
			want: decimal.NewFromInt(0),
		},
		{
			name: "invalid amount",
			data: map[string]string{
				"amount": "invalid_amount",
			},
			want: decimal.Zero,
		},
		{
			name: "amount without currency",
			data: map[string]string{
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
					SupplyChange: storageTypes.NumericZero(),
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
		data map[string]string
		want decimal.Decimal
	}{
		{
			name: "valid amount",
			data: map[string]string{
				"amount": "1000000000000000000utia",
			},
			want: decimal.NewFromInt(-1000000000000000000),
		}, {
			name: "valid amount but no utia",
			data: map[string]string{
				"amount": "1000000000000000000test",
			},
			want: decimal.NewFromInt(0),
		},
		{
			name: "invalid amount",
			data: map[string]string{
				"amount": "invalid_amount",
			},
			want: decimal.Zero,
		},
		{
			name: "amount without currency",
			data: map[string]string{
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
					SupplyChange: storageTypes.NumericZero(),
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
		Balances: []storage.Balance{
			{
				Currency:  currency.DefaultCurrency,
				Spendable: storageTypes.NumericFromInt64(100),
				Delegated: storageTypes.NumericFromInt64(50),
				Unbonding: storageTypes.NumericFromInt64(20),
			},
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
		Balances: []storage.Balance{
			{
				Currency:  currency.DefaultCurrency,
				Spendable: storageTypes.NumericFromInt64(100),
				Delegated: storageTypes.NumericFromInt64(50),
				Unbonding: storageTypes.NumericFromInt64(20),
			},
		},
	}

	err := ctx.AddAddress(address)
	require.NoError(t, err)

	addressUpdate := &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Balances: []storage.Balance{
			{
				Currency:  currency.DefaultCurrency,
				Spendable: storageTypes.NumericFromInt64(50),
				Delegated: storageTypes.NumericFromInt64(25),
				Unbonding: storageTypes.NumericFromInt64(10),
			},
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
		Balances: []storage.Balance{
			{
				Currency:  currency.DefaultCurrency,
				Spendable: storageTypes.NumericFromInt64(150),
				Delegated: storageTypes.NumericFromInt64(75),
				Unbonding: storageTypes.NumericFromInt64(30),
			},
		},
	}, addr)
}

func Test_AddAddress_ExistingWithInvalidCurrency(t *testing.T) {
	ctx := NewContext()

	address := &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Balances: []storage.Balance{
			{
				Currency:  currency.DefaultCurrency,
				Spendable: storageTypes.NumericFromInt64(100),
				Delegated: storageTypes.NumericFromInt64(50),
				Unbonding: storageTypes.NumericFromInt64(20),
			},
		},
	}

	err := ctx.AddAddress(address)
	require.NoError(t, err)

	addressUpdate := &storage.Address{
		Address:    "celestia1ydgj7csawc0k4f7qguy6zd5vs7q5cqx5cepy5e",
		Height:     1,
		LastHeight: 1,
		Balances: []storage.Balance{
			{
				Currency:  "invalid_currency",
				Spendable: storageTypes.NumericFromInt64(50),
				Delegated: storageTypes.NumericFromInt64(25),
				Unbonding: storageTypes.NumericFromInt64(10),
			},
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
		Balances: []storage.Balance{
			{
				Currency:  currency.DefaultCurrency,
				Spendable: storageTypes.NumericFromInt64(100),
				Delegated: storageTypes.NumericFromInt64(50),
				Unbonding: storageTypes.NumericFromInt64(20),
			},
			{
				Currency:  "invalid_currency",
				Spendable: storageTypes.NumericFromInt64(50),
				Delegated: storageTypes.NumericFromInt64(25),
				Unbonding: storageTypes.NumericFromInt64(10),
			},
		},
	}, addr)
}

// Test_AddNamespace_Merge covers the accumulator every namespace counter relies
// on: several blobs to one namespace inside one block must add up, and fibre
// traffic must be counted apart from PFB traffic.
func Test_AddNamespace_Merge(t *testing.T) {
	nsID := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22}

	newNs := func() *storage.Namespace {
		return &storage.Namespace{Version: 0, NamespaceID: nsID}
	}

	tests := []struct {
		name string
		add  []*storage.Namespace
		want storage.Namespace
	}{
		{
			name: "two pfb blobs",
			add: []*storage.Namespace{
				{Version: 0, NamespaceID: nsID, PfbCount: 1, Size: 10, BlobsCount: 1},
				{Version: 0, NamespaceID: nsID, PfbCount: 1, Size: 20, BlobsCount: 1},
			},
			want: storage.Namespace{Version: 0, NamespaceID: nsID, PfbCount: 2, Size: 30, BlobsCount: 2},
		}, {
			name: "two fibre blobs",
			add: []*storage.Namespace{
				{Version: 0, NamespaceID: nsID, PffCount: 1, FibreSize: 262144, BlobsCount: 1},
				{Version: 0, NamespaceID: nsID, PffCount: 1, FibreSize: 262144, BlobsCount: 1},
			},
			want: storage.Namespace{Version: 0, NamespaceID: nsID, PffCount: 2, FibreSize: 524288, BlobsCount: 2},
		}, {
			// The order matters: a PFB seen first must not swallow the fibre
			// counters of a PFF that lands later in the same block.
			name: "pfb then fibre",
			add: []*storage.Namespace{
				{Version: 0, NamespaceID: nsID, PfbCount: 1, Size: 10, BlobsCount: 1},
				{Version: 0, NamespaceID: nsID, PffCount: 1, FibreSize: 262144, BlobsCount: 1},
			},
			want: storage.Namespace{
				Version: 0, NamespaceID: nsID,
				PfbCount: 1, Size: 10,
				PffCount: 1, FibreSize: 262144,
				BlobsCount: 2,
			},
		}, {
			name: "fibre then pfb",
			add: []*storage.Namespace{
				{Version: 0, NamespaceID: nsID, PffCount: 1, FibreSize: 262144, BlobsCount: 1},
				{Version: 0, NamespaceID: nsID, PfbCount: 1, Size: 10, BlobsCount: 1},
			},
			want: storage.Namespace{
				Version: 0, NamespaceID: nsID,
				PfbCount: 1, Size: 10,
				PffCount: 1, FibreSize: 262144,
				BlobsCount: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()

			var last *storage.Namespace
			for _, ns := range tt.add {
				last = ctx.AddNamespace(ns)
			}

			require.Len(t, ctx.Namespaces.Values(), 1)
			// AddNamespace returns the shared instance, so blob logs of both
			// messages point at the same accumulated row.
			require.Same(t, ctx.Namespaces.Values()[0], last)
			require.Equal(t, tt.want, *last)
		})
	}

	t.Run("different namespaces stay apart", func(t *testing.T) {
		ctx := NewContext()
		other := newNs()
		other.NamespaceID = append([]byte{1}, nsID[1:]...)

		ctx.AddNamespace(newNs())
		ctx.AddNamespace(other)
		require.Len(t, ctx.Namespaces.Values(), 2)
	})
}
