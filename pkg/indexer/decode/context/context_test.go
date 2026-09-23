// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package context

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
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

func Test_AddIbcClient_Merge(t *testing.T) {
	t0 := time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC)
	t1 := t0.Add(time.Second)

	ctx := NewContext()
	ctx.AddIbcClient(&storage.IbcClient{
		Id:                   "07-tendermint-1",
		Type:                 "07-tendermint",
		CreatedAt:            t0,
		UpdatedAt:            t0,
		TrustingPeriod:       time.Hour,
		LatestRevisionHeight: 100,
		LatestRevisionNumber: 1,
		Creator:              &storage.Address{Address: "celestia1signer"},
	})
	ctx.AddIbcClient(&storage.IbcClient{
		Id:                   "07-tendermint-1",
		UpdatedAt:            t1,
		ChainId:              "osmosis-1",
		LatestRevisionHeight: 150,
		LatestRevisionNumber: 1,
	})
	ctx.AddIbcClient(&storage.IbcClient{
		Id:                   "07-tendermint-1",
		UpdatedAt:            t1,
		FrozenRevisionHeight: 1,
	})
	ctx.AddIbcClient(&storage.IbcClient{
		Id:              "07-tendermint-1",
		ConnectionCount: 2,
	})

	require.Equal(t, 1, ctx.IbcClients.Len())
	client, ok := ctx.IbcClients.Get("07-tendermint-1")
	require.True(t, ok)
	require.Equal(t, "07-tendermint", client.Type)
	require.Equal(t, t0, client.CreatedAt)
	require.Equal(t, t1, client.UpdatedAt)
	require.Equal(t, time.Hour, client.TrustingPeriod)
	require.Equal(t, "osmosis-1", client.ChainId)
	require.EqualValues(t, 150, client.LatestRevisionHeight)
	require.EqualValues(t, 1, client.LatestRevisionNumber)
	require.EqualValues(t, 1, client.FrozenRevisionHeight)
	require.EqualValues(t, 2, client.ConnectionCount)
	require.Equal(t, "celestia1signer", client.Creator.Address)
}

// Chain keeps max(LatestHeight): a later update to a lower height must not move it back.
func Test_AddIbcClient_LatestHeightIsMax(t *testing.T) {
	ctx := NewContext()
	ctx.AddIbcClient(&storage.IbcClient{Id: "07-tendermint-1", LatestRevisionHeight: 150, LatestRevisionNumber: 1})
	ctx.AddIbcClient(&storage.IbcClient{Id: "07-tendermint-1", LatestRevisionHeight: 120, LatestRevisionNumber: 1})

	client, ok := ctx.IbcClients.Get("07-tendermint-1")
	require.True(t, ok)
	require.EqualValues(t, 150, client.LatestRevisionHeight)
}

func Test_AddIbcClient_NewRevisionResetsHeight(t *testing.T) {
	ctx := NewContext()
	ctx.AddIbcClient(&storage.IbcClient{Id: "07-tendermint-1", LatestRevisionHeight: 1000, LatestRevisionNumber: 1})
	ctx.AddIbcClient(&storage.IbcClient{Id: "07-tendermint-1", LatestRevisionHeight: 5, LatestRevisionNumber: 2})
	ctx.AddIbcClient(&storage.IbcClient{Id: "07-tendermint-1", LatestRevisionHeight: 2000, LatestRevisionNumber: 1})

	client, ok := ctx.IbcClients.Get("07-tendermint-1")
	require.True(t, ok)
	require.EqualValues(t, 2, client.LatestRevisionNumber)
	require.EqualValues(t, 5, client.LatestRevisionHeight)
}

func Test_AddIbcChannelTransfer(t *testing.T) {
	ctx := NewContext()

	// outgoing from celestia: counted as sent
	ctx.AddIbcChannelTransfer(&storage.IbcTransfer{
		ChannelId: "channel-2",
		Amount:    storageTypes.NumericFromInt64(100),
		Sender:    &storage.Address{Address: "celestia1sender"},
	})
	// incoming to celestia: counted as received
	ctx.AddIbcChannelTransfer(&storage.IbcTransfer{
		ChannelId: "channel-2",
		Amount:    storageTypes.NumericFromInt64(30),
		Receiver:  &storage.Address{Address: "celestia1receiver"},
	})
	ctx.AddIbcChannelTransfer(&storage.IbcTransfer{
		ChannelId: "channel-3",
		Amount:    storageTypes.NumericFromInt64(7),
		Receiver:  &storage.Address{Address: "celestia1receiver"},
	})

	require.Equal(t, 2, ctx.IbcChannels.Len())

	ch, ok := ctx.IbcChannels.Get("channel-2")
	require.True(t, ok)
	require.EqualValues(t, 2, ch.TransfersCount)
	require.Equal(t, "100", ch.Sent.String())
	require.Equal(t, "30", ch.Received.String())
	require.Equal(t, storageTypes.IbcChannelStatusInitialization, ch.Status)

	ch, ok = ctx.IbcChannels.Get("channel-3")
	require.True(t, ok)
	require.EqualValues(t, 1, ch.TransfersCount)
	require.True(t, ch.Sent.IsZero())
	require.Equal(t, "7", ch.Received.String())
}

// Jail carries two numbers taken from the same `burned_coins` attribute: Burned is the
// magnitude written to the jail row, Validator.Stake is the delta the storage module adds
// to the stored stake. So a burn of 1000 arrives as Burned=1000 and Stake=-1000.
func jailOf(consAddress, reason string, burned int64) storage.Jail {
	jailed := true
	return storage.Jail{
		Reason: reason,
		Burned: storageTypes.NumericFromInt64(burned),
		Validator: &storage.Validator{
			ConsAddress: consAddress,
			Stake:       storageTypes.NumericFromInt64(-burned),
			Jailed:      &jailed,
		},
	}
}

func Test_AddJail_New(t *testing.T) {
	ctx := NewContext()
	ctx.AddJail(jailOf("A5B3", "double_sign", 1000))

	require.Equal(t, 1, ctx.Jails.Len())
	jail, ok := ctx.Jails.Get("A5B3")
	require.True(t, ok)
	require.Equal(t, "double_sign", jail.Reason)
	require.Equal(t, "1000", jail.Burned.String())
	require.Equal(t, "-1000", jail.Validator.Stake.String())
	require.NotNil(t, jail.Validator.Jailed)
	require.True(t, *jail.Validator.Jailed)
}

// Two slashes of the same validator in one block accumulate: the burned magnitudes are
// summed, the stake deltas are summed as well and stay negative.
func Test_AddJail_MergeTwoBurns(t *testing.T) {
	ctx := NewContext()
	ctx.AddJail(jailOf("A5B3", "double_sign", 1000))
	ctx.AddJail(jailOf("A5B3", "double_sign", 500))

	require.Equal(t, 1, ctx.Jails.Len())
	jail, ok := ctx.Jails.Get("A5B3")
	require.True(t, ok)
	require.Equal(t, "1500", jail.Burned.String())
	require.Equal(t, "-1500", jail.Validator.Stake.String())
}

// A downtime jail burns nothing on Celestia, so a later double sign in the same block must
// still bring its own burn and reason through.
func Test_AddJail_MergeZeroBurnThenSlash(t *testing.T) {
	ctx := NewContext()
	ctx.AddJail(jailOf("A5B3", "missing_signature", 0))
	ctx.AddJail(jailOf("A5B3", "double_sign", 1000))

	require.Equal(t, 1, ctx.Jails.Len())
	jail, ok := ctx.Jails.Get("A5B3")
	require.True(t, ok)
	require.Equal(t, "double_sign", jail.Reason)
	require.Equal(t, "1000", jail.Burned.String())
	require.Equal(t, "-1000", jail.Validator.Stake.String())
}

// A zero burn must leave what the previous slash accumulated untouched.
func Test_AddJail_MergeSlashThenZeroBurn(t *testing.T) {
	ctx := NewContext()
	ctx.AddJail(jailOf("A5B3", "double_sign", 1000))
	ctx.AddJail(jailOf("A5B3", "missing_signature", 0))

	jail, ok := ctx.Jails.Get("A5B3")
	require.True(t, ok)
	require.Equal(t, "1000", jail.Burned.String())
	require.Equal(t, "-1000", jail.Validator.Stake.String())
}

func Test_AddJail_DifferentValidators(t *testing.T) {
	ctx := NewContext()
	ctx.AddJail(jailOf("A5B3", "double_sign", 1000))
	ctx.AddJail(jailOf("1122", "double_sign", 500))

	require.Equal(t, 2, ctx.Jails.Len())
	first, ok := ctx.Jails.Get("A5B3")
	require.True(t, ok)
	require.Equal(t, "-1000", first.Validator.Stake.String())
	second, ok := ctx.Jails.Get("1122")
	require.True(t, ok)
	require.Equal(t, "-500", second.Validator.Stake.String())
}

func cancelOf(validator, delegator string, creationHeight pkgTypes.Level, amount int64) storage.Undelegation {
	return storage.Undelegation{
		Height:         2000,
		CreationHeight: creationHeight,
		Amount:         storageTypes.NumericFromInt64(amount),
		Validator:      &storage.Validator{Address: validator},
		Address:        &storage.Address{Address: delegator},
	}
}

func Test_AddCancelUndelegation_SameEntrySums(t *testing.T) {
	ctx := NewContext()
	require.NoError(t, ctx.AddCancelUndelegation(cancelOf("valoper1", "celestia1", 1000, 200)))
	require.NoError(t, ctx.AddCancelUndelegation(cancelOf("valoper1", "celestia1", 1000, 400)))

	values := ctx.CancelUnbonding.Values()
	require.Len(t, values, 1)
	require.Equal(t, "600", values[0].Amount.String())
	require.EqualValues(t, 1000, values[0].CreationHeight)
	require.Equal(t, "valoper1", values[0].Validator.Address)
	require.Equal(t, "celestia1", values[0].Address.Address)
}

// Validator is built by EmptyValidator with only the operator address set, so the key must use it
func Test_AddCancelUndelegation_DifferentValidators(t *testing.T) {
	ctx := NewContext()
	require.NoError(t, ctx.AddCancelUndelegation(cancelOf("valoper1", "celestia1", 1000, 200)))
	require.NoError(t, ctx.AddCancelUndelegation(cancelOf("valoper2", "celestia1", 1000, 400)))

	require.Equal(t, 2, ctx.CancelUnbonding.Len())
	amounts := make(map[string]string)
	for u := range ctx.CancelUnbonding.AllValues() {
		amounts[u.Validator.Address] = u.Amount.String()
	}
	require.Equal(t, map[string]string{"valoper1": "200", "valoper2": "400"}, amounts)
}

func Test_AddCancelUndelegation_DifferentKeys(t *testing.T) {
	ctx := NewContext()
	require.NoError(t, ctx.AddCancelUndelegation(cancelOf("valoper1", "celestia1", 1000, 200)))
	require.NoError(t, ctx.AddCancelUndelegation(cancelOf("valoper1", "celestia1", 1001, 300)))
	require.NoError(t, ctx.AddCancelUndelegation(cancelOf("valoper1", "celestia2", 1000, 400)))

	require.Equal(t, 3, ctx.CancelUnbonding.Len())
	for u := range ctx.CancelUnbonding.AllValues() {
		switch {
		case u.Address.Address == "celestia2":
			require.Equal(t, "400", u.Amount.String())
		case u.CreationHeight == 1001:
			require.Equal(t, "300", u.Amount.String())
		default:
			require.Equal(t, "200", u.Amount.String())
		}
	}
}

// Merging must not mutate the amount of the event that was added second
func Test_AddCancelUndelegation_DoesNotMutateInput(t *testing.T) {
	ctx := NewContext()
	first := cancelOf("valoper1", "celestia1", 1000, 200)
	second := cancelOf("valoper1", "celestia1", 1000, 400)
	require.NoError(t, ctx.AddCancelUndelegation(first))
	require.NoError(t, ctx.AddCancelUndelegation(second))

	require.Equal(t, "200", first.Amount.String())
	require.Equal(t, "400", second.Amount.String())
}

func Test_AddCancelUndelegation_NilPointers(t *testing.T) {
	ctx := NewContext()

	noValidator := cancelOf("valoper1", "celestia1", 1000, 200)
	noValidator.Validator = nil
	require.Error(t, ctx.AddCancelUndelegation(noValidator))

	noAddress := cancelOf("valoper1", "celestia1", 1000, 200)
	noAddress.Address = nil
	require.Error(t, ctx.AddCancelUndelegation(noAddress))

	require.Equal(t, 0, ctx.CancelUnbonding.Len())
}
