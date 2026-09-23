// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/require"
)

var (
	testAddress     = "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"
	testHashAddress = []byte{0x96, 0xa, 0xa0, 0x36, 0x6b, 0x25, 0x4e, 0x1e, 0xa7, 0x9b, 0xda, 0x46, 0x7e, 0xb3, 0xaa, 0x5c, 0x97, 0xcb, 0xa5, 0xae}
	testIgpId       = []uint8{0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}
	testBlock       = storage.Block{
		Height: pkgTypes.Level(123456),
		Time:   time.Now(),
	}
)

func Test_parseCoinSpent(t *testing.T) {
	ibcDenom := "ibc/93E113CD8DF31891647AE56271673AEFA7E125A686AEC6CAD8D5106FE9600892"
	tests := []struct {
		name    string
		data    map[string]string
		height  pkgTypes.Level
		want    *storage.Address
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]string{
				"spender": testAddress,
				"amount":  "123utia",
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balances: []storage.Balance{
					{
						Currency:  currency.DefaultCurrency,
						Spendable: types.MustNumericFromString("-123"),
						Delegated: types.NumericZero(),
						Unbonding: types.NumericZero(),
					},
				},
			},
		}, {
			name: "test 2",
			data: map[string]string{
				"spender": testAddress,
				"amount":  "",
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balances:   []storage.Balance{},
			},
		}, {
			name: "test 3 multi-coin IBC and utia",
			data: map[string]string{
				"spender": testAddress,
				"amount":  "5000000" + ibcDenom + ",5000000utia",
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balances: []storage.Balance{
					{
						Currency:  ibcDenom,
						Spendable: types.MustNumericFromString("-5000000"),
						Delegated: types.NumericZero(),
						Unbonding: types.NumericZero(),
					},
					{
						Currency:  currency.DefaultCurrency,
						Spendable: types.MustNumericFromString("-5000000"),
						Delegated: types.NumericZero(),
						Unbonding: types.NumericZero(),
					},
				},
			},
		}, {
			name: "test 4 IBC only",
			data: map[string]string{
				"spender": testAddress,
				"amount":  "5000000" + ibcDenom,
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balances: []storage.Balance{
					{
						Currency:  ibcDenom,
						Spendable: types.MustNumericFromString("-5000000"),
						Delegated: types.NumericZero(),
						Unbonding: types.NumericZero(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			err := parseCoinSpent(ctx, tt.data, tt.height)
			require.True(t, (err == nil) != tt.wantErr)
			require.EqualValues(t, 1, ctx.Addresses.Len())
			for _, value := range ctx.Addresses.All() {
				require.Equal(t, tt.want, value)
			}
		})
	}
}

func Test_parseCoinReceived(t *testing.T) {
	ibcDenom := "ibc/93E113CD8DF31891647AE56271673AEFA7E125A686AEC6CAD8D5106FE9600892"
	tests := []struct {
		name    string
		data    map[string]string
		height  pkgTypes.Level
		want    *storage.Address
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]string{
				"receiver": testAddress,
				"amount":   "123utia",
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balances: []storage.Balance{
					{
						Currency:  currency.DefaultCurrency,
						Spendable: types.MustNumericFromString("123"),
						Delegated: types.NumericZero(),
						Unbonding: types.NumericZero(),
					},
				},
			},
		}, {
			name: "test 2",
			data: map[string]string{
				"receiver": testAddress,
				"amount":   "",
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balances:   []storage.Balance{},
			},
		}, {
			name: "test 3 multi-coin IBC and utia",
			data: map[string]string{
				"receiver": testAddress,
				"amount":   "5000000" + ibcDenom + ",5000000utia",
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balances: []storage.Balance{
					{
						Currency:  ibcDenom,
						Spendable: types.MustNumericFromString("5000000"),
						Delegated: types.NumericZero(),
						Unbonding: types.NumericZero(),
					},
					{
						Currency:  currency.DefaultCurrency,
						Spendable: types.MustNumericFromString("5000000"),
						Delegated: types.NumericZero(),
						Unbonding: types.NumericZero(),
					},
				},
			},
		}, {
			name: "test 4 IBC only",
			data: map[string]string{
				"receiver": testAddress,
				"amount":   "5000000" + ibcDenom,
			},
			height: pkgTypes.Level(58000),
			want: &storage.Address{
				Height:     pkgTypes.Level(58000),
				LastHeight: pkgTypes.Level(58000),
				Address:    testAddress,
				Hash:       testHashAddress,
				Balances: []storage.Balance{
					{
						Currency:  ibcDenom,
						Spendable: types.MustNumericFromString("5000000"),
						Delegated: types.NumericZero(),
						Unbonding: types.NumericZero(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()
			err := parseCoinReceived(ctx, tt.data, tt.height)
			require.True(t, (err == nil) != tt.wantErr)
			require.EqualValues(t, 1, ctx.Addresses.Len())
			for _, value := range ctx.Addresses.All() {
				require.Equal(t, tt.want, value)
			}
		})
	}
}

func Test_parseCreateIgp(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &testBlock

	tests := []struct {
		name    string
		data    map[string]string
		want    *storage.HLIGP
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]string{
				"denom":  "\"utia\"",
				"igp_id": "\"0x726f757465725f706f73745f6469737061746368000000040000000000000001\"",
				"owner":  testAddress,
			},
			want: &storage.HLIGP{
				Height: pkgTypes.Level(123456),
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
			for _, value := range ctx.Igps.All() {
				require.Equal(t, tt.want, value)
			}
		})
	}
}

func Test_parseSetDestinationGasConfig(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &testBlock

	tests := []struct {
		name    string
		data    map[string]string
		want    *storage.HLIGPConfig
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]string{
				"gas_overhead":        "\"200000\"",
				"gas_price":           "\"1\"",
				"igp_id":              "\"0x726f757465725f706f73745f6469737061746368000000040000000000000001\"",
				"owner":               testAddress,
				"remote_domain":       "84532",
				"token_exchange_rate": "\"10000000000\"",
			},
			want: &storage.HLIGPConfig{
				Height:            pkgTypes.Level(123456),
				Time:              testBlock.Time,
				GasPrice:          types.NumericFromInt64(1),
				GasOverhead:       types.NumericFromInt64(200000),
				RemoteDomain:      84532,
				TokenExchangeRate: "10000000000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseSetDestinationGasConfig(ctx, tt.data)
			require.True(t, (err == nil) != tt.wantErr)
			require.EqualValues(t, 1, ctx.IgpConfigs.Len())
			for _, value := range ctx.IgpConfigs.All() {
				require.Equal(t, tt.want, value)
			}
		})
	}
}

func Test_parseSetIgp(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &testBlock

	tests := []struct {
		name    string
		data    map[string]string
		want    *storage.HLIGP
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]string{
				"igp_id":             "\"0x726f757465725f706f73745f6469737061746368000000040000000000000001\"",
				"owner":              testAddress,
				"new_owner":          "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w61",
				"renounce_ownership": "false",
			},
			want: &storage.HLIGP{
				Height: pkgTypes.Level(123456),
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
			for _, value := range ctx.Igps.All() {
				require.Equal(t, tt.want, value)
			}
		})
	}
}

func Test_parseCompleteUnbonding(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &testBlock

	tests := []struct {
		name    string
		data    map[string]string
		want    storage.StakingLog
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]string{
				"amount":    "35570000utia",
				"delegator": "celestia1nwm73xdjhdwfpw6uxc3pkcspw0kr5m06z067t3",
				"mode":      "EndBlock",
				"validator": "celestiavaloper140l6y2gp3gxvay6qtn70re7z2s0gn57zcvqd22",
			},
			want: storage.StakingLog{
				Height: pkgTypes.Level(123456),
				Time:   testBlock.Time,
				Address: &storage.Address{
					Address:    "celestia1nwm73xdjhdwfpw6uxc3pkcspw0kr5m06z067t3",
					Height:     pkgTypes.Level(123456),
					LastHeight: pkgTypes.Level(123456),
					Hash:       []byte{0x9b, 0xb7, 0xe8, 0x99, 0xb2, 0xbb, 0x5c, 0x90, 0xbb, 0x5c, 0x36, 0x22, 0x1b, 0x62, 0x01, 0x73, 0xec, 0x3a, 0x6d, 0xfa},
					Balances: []storage.Balance{
						{
							Currency:  currency.Utia,
							Unbonding: types.MustNumericFromString("-35570000"),
							Spendable: types.NumericZero(),
							Delegated: types.NumericZero(),
						},
					},
				},
				Validator: &storage.Validator{
					Address:           "celestiavaloper140l6y2gp3gxvay6qtn70re7z2s0gn57zcvqd22",
					Moniker:           storage.DoNotModify,
					Website:           storage.DoNotModify,
					Identity:          storage.DoNotModify,
					Contacts:          storage.DoNotModify,
					Details:           storage.DoNotModify,
					Rate:              types.NumericZero(),
					MaxRate:           types.NumericZero(),
					MaxChangeRate:     types.NumericZero(),
					MinSelfDelegation: types.NumericZero(),
					Stake:             types.NumericZero(),
					Rewards:           types.NumericZero(),
					Commissions:       types.NumericZero(),
				},
				Change: types.MustNumericFromString("-35570000"),
				Type:   types.StakingLogTypeUnbonded,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseCompleteUnbonding(ctx, tt.data)
			require.True(t, (err == nil) != tt.wantErr)
			require.Len(t, ctx.StakingLogs, 1)
			require.Equal(t, tt.want, ctx.StakingLogs[0])
			require.EqualValues(t, ctx.Addresses.Len(), 1)
			addr, ok := ctx.Addresses.Get("celestia1nwm73xdjhdwfpw6uxc3pkcspw0kr5m06z067t3")
			require.True(t, ok)
			require.Equal(t, tt.want.Address, addr)
			require.EqualValues(t, ctx.Validators.Len(), 1)
			val, ok := ctx.Validators.Get("celestiavaloper140l6y2gp3gxvay6qtn70re7z2s0gn57zcvqd22")
			require.True(t, ok)
			require.Equal(t, tt.want.Validator, val)
		})
	}
}

func Test_parseCompleteRedelegation(t *testing.T) {
	ctx := context.NewContext()
	ctx.Block = &testBlock

	tests := []struct {
		name    string
		data    map[string]string
		want    storage.StakingLog
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]string{
				"amount":                "265636688utia",
				"delegator":             "celestia1nwm73xdjhdwfpw6uxc3pkcspw0kr5m06z067t3",
				"destination_validator": "celestiavaloper1snun9qqk9eussvyhkqm03lz6f265ekhnnlw043",
				"mode":                  "EndBlock",
				"source_validator":      "celestiavaloper140l6y2gp3gxvay6qtn70re7z2s0gn57zcvqd22",
			},
			want: storage.StakingLog{
				Height: pkgTypes.Level(123456),
				Time:   testBlock.Time,
				Address: &storage.Address{
					Address:    "celestia1nwm73xdjhdwfpw6uxc3pkcspw0kr5m06z067t3",
					Height:     pkgTypes.Level(123456),
					LastHeight: pkgTypes.Level(123456),
					Hash:       []byte{0x9b, 0xb7, 0xe8, 0x99, 0xb2, 0xbb, 0x5c, 0x90, 0xbb, 0x5c, 0x36, 0x22, 0x1b, 0x62, 0x01, 0x73, 0xec, 0x3a, 0x6d, 0xfa},
					Balances:   []storage.Balance{},
				},
				Validator: &storage.Validator{
					Address:           "celestiavaloper140l6y2gp3gxvay6qtn70re7z2s0gn57zcvqd22",
					Moniker:           storage.DoNotModify,
					Website:           storage.DoNotModify,
					Identity:          storage.DoNotModify,
					Contacts:          storage.DoNotModify,
					Details:           storage.DoNotModify,
					Rate:              types.NumericZero(),
					MaxRate:           types.NumericZero(),
					MaxChangeRate:     types.NumericZero(),
					MinSelfDelegation: types.NumericZero(),
					Stake:             types.NumericZero(),
					Rewards:           types.NumericZero(),
					Commissions:       types.NumericZero(),
				},
				Change: types.MustNumericFromString("-265636688"),
				Type:   types.StakingLogTypeUnbonded,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseCompleteRedelegation(ctx, tt.data)
			require.True(t, (err == nil) != tt.wantErr)
			require.Len(t, ctx.StakingLogs, 1)
			require.Equal(t, tt.want, ctx.StakingLogs[0])
			require.EqualValues(t, ctx.Addresses.Len(), 1)
			addr, ok := ctx.Addresses.Get("celestia1nwm73xdjhdwfpw6uxc3pkcspw0kr5m06z067t3")
			require.True(t, ok)
			require.Equal(t, tt.want.Address, addr)
			require.EqualValues(t, ctx.Validators.Len(), 1)
			val, ok := ctx.Validators.Get("celestiavaloper140l6y2gp3gxvay6qtn70re7z2s0gn57zcvqd22")
			require.True(t, ok)
			require.Equal(t, tt.want.Validator, val)
		})
	}
}

func Test_parseSlash(t *testing.T) {
	// the event carries a bech32 consensus address, the context is keyed by its uppercase hex
	const (
		valcons    = "celestiavalcons15keu050z7sy3ydzk0zgpydzk0zg2hn00835r7k"
		valconsHex = "A5B3C7D1E2F40912345678901234567890ABCDEF"
	)

	t.Run("double sign", func(t *testing.T) {
		ctx := context.NewContext()
		ctx.Block = &testBlock

		err := parseSlash(ctx, map[string]string{
			"address":      valcons,
			"power":        "100",
			"reason":       "double_sign",
			"burned_coins": "1000",
		})
		require.NoError(t, err)
		require.EqualValues(t, 1, ctx.Jails.Len())

		jail, ok := ctx.Jails.Get(valconsHex)
		require.True(t, ok)
		require.Equal(t, testBlock.Height, jail.Height)
		require.Equal(t, testBlock.Time, jail.Time)
		require.Equal(t, "double_sign", jail.Reason)
		require.Equal(t, "1000", jail.Burned.String())
		require.Equal(t, valconsHex, jail.Validator.ConsAddress)
		// stake is a delta: the burn must arrive negated or the jail would raise the stake
		require.Equal(t, "-1000", jail.Validator.Stake.String())
		require.NotNil(t, jail.Validator.Jailed)
		require.True(t, *jail.Validator.Jailed)
	})

	// Celestia runs with slash_fraction_downtime = 0, so a downtime jail burns nothing
	t.Run("downtime burns nothing", func(t *testing.T) {
		ctx := context.NewContext()
		ctx.Block = &testBlock

		err := parseSlash(ctx, map[string]string{
			"address":      valcons,
			"power":        "100",
			"reason":       "missing_signature",
			"jailed":       valcons,
			"burned_coins": "0",
		})
		require.NoError(t, err)
		require.EqualValues(t, 1, ctx.Jails.Len())

		jail, ok := ctx.Jails.Get(valconsHex)
		require.True(t, ok)
		require.Equal(t, "missing_signature", jail.Reason)
		require.True(t, jail.Burned.IsZero())
		require.True(t, jail.Validator.Stake.IsZero())
		require.NotNil(t, jail.Validator.Jailed)
		require.True(t, *jail.Validator.Jailed)
	})

	// the second event of the evidence path carries only `jailed` and must be skipped
	t.Run("jail only event is ignored", func(t *testing.T) {
		ctx := context.NewContext()
		ctx.Block = &testBlock

		err := parseSlash(ctx, map[string]string{
			"jailed": valcons,
		})
		require.NoError(t, err)
		require.EqualValues(t, 0, ctx.Jails.Len())
	})

	t.Run("two slashes in one block accumulate", func(t *testing.T) {
		ctx := context.NewContext()
		ctx.Block = &testBlock

		for _, burned := range []string{"1000", "500"} {
			err := parseSlash(ctx, map[string]string{
				"address":      valcons,
				"power":        "100",
				"reason":       "double_sign",
				"burned_coins": burned,
			})
			require.NoError(t, err)
		}
		require.EqualValues(t, 1, ctx.Jails.Len())

		jail, ok := ctx.Jails.Get(valconsHex)
		require.True(t, ok)
		require.Equal(t, "1500", jail.Burned.String())
		require.Equal(t, "-1500", jail.Validator.Stake.String())
	})

	t.Run("invalid address", func(t *testing.T) {
		ctx := context.NewContext()
		ctx.Block = &testBlock

		err := parseSlash(ctx, map[string]string{
			"address":      "invalid",
			"reason":       "double_sign",
			"burned_coins": "1000",
		})
		require.Error(t, err)
		require.EqualValues(t, 0, ctx.Jails.Len())
	})
}
