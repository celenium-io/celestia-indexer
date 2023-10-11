// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package genesis

import (
	"os"
	"testing"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/goccy/go-json"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestParseAccounts(t *testing.T) {
	f, err := os.Open("../../../test/json/genesis.json")
	require.NoError(t, err)
	defer f.Close()

	var g types.Genesis
	err = json.NewDecoder(f).Decode(&g)
	require.NoError(t, err)

	data := newParsedData()

	module := NewModule(postgres.Storage{}, config.Indexer{})

	module.parseDenomMetadata(g.AppState.Bank.DenomMetadata, &data)

	err = module.parseAccounts(g.AppState.Auth.Accounts, 1, &data)
	require.NoError(t, err)

	want := map[string]*storage.Address{
		"celestia1qqqpkhsnpyvtzx4knu53zsdfn7l88czztlp8tt": {
			Address:    "celestia1qqqpkhsnpyvtzx4knu53zsdfn7l88czztlp8tt",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0x0, 0x0, 0x1b, 0x5e, 0x13, 0x9, 0x18, 0xb1, 0x1a, 0xb6, 0x9f, 0x29, 0x11, 0x41, 0xa9, 0x9f, 0xbe, 0x73, 0xe0, 0x42},
			Balance: storage.Balance{
				Id:       0,
				Total:    decimal.Zero,
				Currency: "utia",
			},
		},
		"celestia1qsfn7xq3spe6g3cvth7p6ld4ea8y0t262udez6": {
			Address:    "celestia1qsfn7xq3spe6g3cvth7p6ld4ea8y0t262udez6",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0x4, 0x13, 0x3f, 0x18, 0x11, 0x80, 0x73, 0xa4, 0x47, 0xc, 0x5d, 0xfc, 0x1d, 0x7d, 0xb5, 0xcf, 0x4e, 0x47, 0xad, 0x5a},
			Balance: storage.Balance{
				Id:       0,
				Total:    decimal.Zero,
				Currency: "utia",
			},
		},
		"celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6": {
			Address:    "celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0x4f, 0xea, 0x76, 0x42, 0x7b, 0x83, 0x45, 0x86, 0x1e, 0x80, 0xa3, 0x54, 0xa, 0x8a, 0x9d, 0x93, 0x6f, 0xd3, 0x93, 0x91},
			Balance: storage.Balance{
				Id:       0,
				Total:    decimal.Zero,
				Currency: "utia",
			},
		},
		"celestia10n95tmwqtc5ua47m9vu52p7xwcf6gcdtjj9rfh": {
			Address:    "celestia10n95tmwqtc5ua47m9vu52p7xwcf6gcdtjj9rfh",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0x7c, 0xcb, 0x45, 0xed, 0xc0, 0x5e, 0x29, 0xce, 0xd7, 0xdb, 0x2b, 0x39, 0x45, 0x7, 0xc6, 0x76, 0x13, 0xa4, 0x61, 0xab},
			Balance: storage.Balance{
				Id:       0,
				Total:    decimal.Zero,
				Currency: "utia",
			},
		},
		"celestia1e6mspkfqg9ud33m4ek3je0glzrlc9f0px9h40k": {
			Address:    "celestia1e6mspkfqg9ud33m4ek3je0glzrlc9f0px9h40k",
			Height:     1,
			LastHeight: 1,
			Hash:       []byte{0xce, 0xb7, 0x0, 0xd9, 0x20, 0x41, 0x78, 0xd8, 0xc7, 0x75, 0xcd, 0xa3, 0x2c, 0xbd, 0x1f, 0x10, 0xff, 0x82, 0xa5, 0xe1},
			Balance: storage.Balance{
				Id:       0,
				Total:    decimal.Zero,
				Currency: "utia",
			},
		},
	}
	require.Equal(t, want, data.addresses)
}
