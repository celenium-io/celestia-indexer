// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	encjson "encoding/json"
	"os"
	"testing"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/stretchr/testify/require"
)

const permanentLockedAccountAddress = "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6"

// addressesMap materializes ctx.Addresses (a sync map) into a plain map so it can be
// compared with require.Equal.
func addressesMap(ctx *decodeContext.Context) map[string]*storage.Address {
	got := make(map[string]*storage.Address, ctx.Addresses.Len())
	for key, addr := range ctx.Addresses.All() {
		got[key] = addr
	}
	return got
}

func TestParseAccounts(t *testing.T) {
	f, err := os.Open("../../../test/json/genesis.json")
	require.NoError(t, err)
	defer f.Close()

	var g types.Genesis
	err = json.ConfigFastest.NewDecoder(f).Decode(&g)
	require.NoError(t, err)

	module := NewModule(postgres.Storage{}, config.Indexer{})

	ctx := decodeContext.NewContext()
	block := storage.Block{
		Height: 1,
		Time:   time.Now(),
	}
	ctx.Block = &block

	module.parseDenomMetadata(ctx, g.AppState.Bank.DenomMetadata)
	denom := getBaseCurrency(ctx.DenomMetadata)

	err = module.parseAccounts(ctx, denom, g.AppState.Auth.Accounts)
	require.NoError(t, err)

	want := map[string]*storage.Address{
		"celestia1qqqpkhsnpyvtzx4knu53zsdfn7l88czztlp8tt": {
			Address:    "celestia1qqqpkhsnpyvtzx4knu53zsdfn7l88czztlp8tt",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0x0, 0x0, 0x1b, 0x5e, 0x13, 0x9, 0x18, 0xb1, 0x1a, 0xb6, 0x9f, 0x29, 0x11, 0x41, 0xa9, 0x9f, 0xbe, 0x73, 0xe0, 0x42},
			Balances:   []storage.Balance{storage.EmptyBalance()},
		},
		"celestia1qsfn7xq3spe6g3cvth7p6ld4ea8y0t262udez6": {
			Address:    "celestia1qsfn7xq3spe6g3cvth7p6ld4ea8y0t262udez6",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0x4, 0x13, 0x3f, 0x18, 0x11, 0x80, 0x73, 0xa4, 0x47, 0xc, 0x5d, 0xfc, 0x1d, 0x7d, 0xb5, 0xcf, 0x4e, 0x47, 0xad, 0x5a},
			Balances:   []storage.Balance{storage.EmptyBalance()},
		},
		"celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6": {
			Address:    "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0x4f, 0xea, 0x76, 0x42, 0x7b, 0x83, 0x45, 0x86, 0x1e, 0x80, 0xa3, 0x54, 0xa, 0x8a, 0x9d, 0x93, 0x6f, 0xd3, 0x93, 0x91},
			Balances:   []storage.Balance{storage.EmptyBalance()},
			Name:       "bonded_tokens_pool",
		},
		"celestia10n95tmwqtc5ua47m9vu52p7xwcf6gcdtjj9rfh": {
			Address:    "celestia10n95tmwqtc5ua47m9vu52p7xwcf6gcdtjj9rfh",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0x7c, 0xcb, 0x45, 0xed, 0xc0, 0x5e, 0x29, 0xce, 0xd7, 0xdb, 0x2b, 0x39, 0x45, 0x7, 0xc6, 0x76, 0x13, 0xa4, 0x61, 0xab},
			Balances:   []storage.Balance{storage.EmptyBalance()},
		},
		"celestia1e6mspkfqg9ud33m4ek3je0glzrlc9f0px9h40k": {
			Address:    "celestia1e6mspkfqg9ud33m4ek3je0glzrlc9f0px9h40k",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0xce, 0xb7, 0x0, 0xd9, 0x20, 0x41, 0x78, 0xd8, 0xc7, 0x75, 0xcd, 0xa3, 0x2c, 0xbd, 0x1f, 0x10, 0xff, 0x82, 0xa5, 0xe1},
			Balances:   []storage.Balance{storage.EmptyBalance()},
		},
		permanentLockedAccountAddress: {
			Address:    permanentLockedAccountAddress,
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
			Balances:   []storage.Balance{storage.EmptyBalance()},
		},
	}
	require.Equal(t, want, addressesMap(ctx))
	require.Len(t, ctx.VestingAccounts, 4)

	var permanent *storage.VestingAccount
	for _, v := range ctx.VestingAccounts {
		if v.Address.Address == permanentLockedAccountAddress {
			permanent = v
			break
		}
	}
	require.NotNil(t, permanent, "expected a vesting entry for the permanent locked account")
	require.Equal(t, storageTypes.VestingTypePermanent, permanent.Type)
	require.Nil(t, permanent.EndTime)
	require.Nil(t, permanent.StartTime)
	require.Empty(t, permanent.VestingPeriods)
}

func TestParseBalances_NewAddressMultiCoin(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	ctx := decodeContext.NewContext()
	block := storage.Block{Height: 1}
	ctx.Block = &block

	const addr = "celestia1qqqpkhsnpyvtzx4knu53zsdfn7l88czztlp8tt"
	balances := []types.Balances{
		{
			Address: addr,
			Coins: []types.Coins{
				{Denom: currency.Utia, Amount: "100"},
				{Denom: "ibc/AAA", Amount: "50"},
			},
		},
	}

	err := module.parseBalances(ctx, balances)
	require.NoError(t, err)

	got, ok := ctx.Addresses.Get(addr)
	require.True(t, ok)
	require.Len(t, got.Balances, 2)

	byCurrency := make(map[string]storageTypes.Numeric)
	for _, b := range got.Balances {
		byCurrency[b.Currency] = b.Spendable
	}
	require.Equal(t, storageTypes.NumericFromInt64(100), byCurrency[currency.Utia])
	require.Equal(t, storageTypes.NumericFromInt64(50), byCurrency["ibc/AAA"])
}

// Regression guard: parseBalances must mutate an existing address in place, not
// re-run it through ctx.AddAddress, which would double every balance it just set.
func TestParseBalances_ExistingAddressNewCurrency(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	ctx := decodeContext.NewContext()
	block := storage.Block{Height: 1}
	ctx.Block = &block

	const addr = "celestia1qqqpkhsnpyvtzx4knu53zsdfn7l88czztlp8tt"
	ctx.Addresses.Set(addr, &storage.Address{
		Address:    addr,
		Height:     1,
		LastHeight: 1,
		Balances:   []storage.Balance{storage.EmptyBalance()},
	})

	balances := []types.Balances{
		{
			Address: addr,
			Coins: []types.Coins{
				{Denom: currency.Utia, Amount: "100"},
				{Denom: "ibc/AAA", Amount: "50"},
			},
		},
	}

	err := module.parseBalances(ctx, balances)
	require.NoError(t, err)

	got, ok := ctx.Addresses.Get(addr)
	require.True(t, ok)
	require.Len(t, got.Balances, 2)

	byCurrency := make(map[string]storageTypes.Numeric)
	for _, b := range got.Balances {
		byCurrency[b.Currency] = b.Spendable
	}
	require.Equal(t, storageTypes.NumericFromInt64(100), byCurrency["utia"])
	require.Equal(t, storageTypes.NumericFromInt64(50), byCurrency["ibc/AAA"])
}

func loadGenesisFixture(t *testing.T) types.Genesis {
	f, err := os.Open("../../../test/json/genesis.json")
	require.NoError(t, err)
	defer f.Close()

	var g types.Genesis
	err = json.ConfigFastest.NewDecoder(f).Decode(&g)
	require.NoError(t, err)
	return g
}

// Exported genesis has validators/delegations in app_state.staking.*, not gen_txs — must fail, not index an empty set.
func TestParse_ExportedGenesisFailsFast(t *testing.T) {
	g := loadGenesisFixture(t)
	g.AppState.Staking.Exported = true

	module := NewModule(postgres.Storage{}, config.Indexer{})
	_, err := module.parse(types.GenesisOutput{Genesis: g})
	require.Error(t, err)
	require.Contains(t, err.Error(), "exported")
}

// gen_txs cleared: the fixture's gen_tx has an empty fee amount, which DecodeFee rejects unrelatedly.
func TestParse_NonExportedGenesisSucceeds(t *testing.T) {
	g := loadGenesisFixture(t)
	require.False(t, g.AppState.Staking.Exported, "fixture is expected to represent a fresh, non-exported genesis")
	g.AppState.Genutil.GenTxs = nil

	module := NewModule(postgres.Storage{}, config.Indexer{})
	_, err := module.parse(types.GenesisOutput{Genesis: g})
	require.NoError(t, err)
}

// Regression guard: an earlier version had parseFeeGrants/parseAuthzGrants each
// snapshot the shared ctx.Grants into data.grants, duplicating entries from the first call.
func TestParse_FeegrantAndAuthzGrants_NoDuplication(t *testing.T) {
	g := loadGenesisFixture(t)
	g.AppState.Genutil.GenTxs = nil
	g.AppState.Feegrant.FeeGrants = []encjson.RawMessage{encjson.RawMessage(basicAllowanceGrantJSON)}
	g.AppState.Authz.Authorization = []encjson.RawMessage{encjson.RawMessage(sendAuthzGrantJSON)}

	module := NewModule(postgres.Storage{}, config.Indexer{})
	dCtx, err := module.parse(types.GenesisOutput{Genesis: g})
	require.NoError(t, err)

	grants := dCtx.Grants.Values()
	require.Len(t, grants, 2)
	var haveFee, haveAuthz bool
	for _, grant := range grants {
		switch grant.Authorization {
		case "fee":
			haveFee = true
		case "/cosmos.bank.v1beta1.MsgSend":
			haveAuthz = true
		}
	}
	require.True(t, haveFee, "feegrant entry must be present exactly once")
	require.True(t, haveAuthz, "authz entry must be present exactly once")
}
