// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"encoding/json"
	"testing"
	"time"

	"cosmossdk.io/x/feegrant"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/stretchr/testify/require"
)

const (
	basicAllowanceGrantJSON = `{
		"granter": "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		"grantee": "celestia1llx8v3hhw06cmtmv7q7uq4ztgyvd33nmk7mrfh",
		"allowance": {
			"@type": "/cosmos.feegrant.v1beta1.BasicAllowance",
			"spend_limit": [],
			"expiration": null
		}
	}`

	allowedMsgAllowanceGrantJSON = `{
		"granter": "celestia1gezy6mf5hy5nmq445djq3s3mwx7pl2xe7r3kp6",
		"grantee": "celestia1lamchvjp7c7d9asy8g9cxuycfa4qf4t6c3ufeg",
		"allowance": {
			"@type": "/cosmos.feegrant.v1beta1.AllowedMsgAllowance",
			"allowance": {
				"@type": "/cosmos.feegrant.v1beta1.BasicAllowance",
				"spend_limit": [],
				"expiration": null
			},
			"allowed_messages": [
				"/cosmos.bank.v1beta1.MsgSend",
				"/celestia.blob.v1.MsgPayForBlobs"
			]
		}
	}`
)

func testBlock() storage.Block {
	return storage.Block{
		Height: 1,
		Time:   time.Now().UTC(),
	}
}

func TestParseFeeGrants_BasicAllowance(t *testing.T) {
	data := newParsedData()
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(basicAllowanceGrantJSON)}
	block := testBlock()
	err := module.parseFeeGrants(raw, block, &data)
	require.NoError(t, err)

	require.Len(t, data.grants, 1)
	grant := data.grants[0]

	require.Equal(t, "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv", grant.Granter.Address)
	require.Equal(t, "celestia1llx8v3hhw06cmtmv7q7uq4ztgyvd33nmk7mrfh", grant.Grantee.Address)
	require.Equal(t, "fee", grant.Authorization)
	require.Equal(t, block.Height, grant.Height)
	require.False(t, grant.Revoked)
	require.Nil(t, grant.Expiration)

	require.NotNil(t, grant.Params)
	spendLimit, ok := grant.Params["SpendLimit"]
	require.True(t, ok)
	require.Empty(t, spendLimit)

	require.Contains(t, data.addresses, grant.Granter.Address)
	require.Contains(t, data.addresses, grant.Grantee.Address)

	granterAddr := data.addresses[grant.Granter.Address]
	require.Len(t, granterAddr.Balances, 1)
	require.True(t, granterAddr.Balances[0].Spendable.IsZero())
	require.Equal(t, currency.Utia, granterAddr.Balances[0].Currency)
}

func TestParseFeeGrants_AllowedMsgAllowance(t *testing.T) {
	data := newParsedData()
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(allowedMsgAllowanceGrantJSON)}
	err := module.parseFeeGrants(raw, testBlock(), &data)
	require.NoError(t, err)

	require.Len(t, data.grants, 1)
	grant := data.grants[0]

	require.Equal(t, "celestia1gezy6mf5hy5nmq445djq3s3mwx7pl2xe7r3kp6", grant.Granter.Address)
	require.Equal(t, "celestia1lamchvjp7c7d9asy8g9cxuycfa4qf4t6c3ufeg", grant.Grantee.Address)
	require.Equal(t, "fee", grant.Authorization)

	require.NotNil(t, grant.Params)

	allowedMessages, ok := grant.Params["AllowedMessages"].([]string)
	require.True(t, ok)
	require.ElementsMatch(t, []string{
		"/cosmos.bank.v1beta1.MsgSend",
		"/celestia.blob.v1.MsgPayForBlobs",
	}, allowedMessages)

	nestedAllowance, ok := grant.Params["Allowance"].(feegrant.BasicAllowance)
	require.True(t, ok)
	require.Empty(t, nestedAllowance.SpendLimit)
	require.Nil(t, nestedAllowance.Expiration)

	require.Nil(t, grant.Expiration)

	require.Contains(t, data.addresses, grant.Granter.Address)
	require.Contains(t, data.addresses, grant.Grantee.Address)
}

func TestParseFeeGrants_Multiple(t *testing.T) {
	data := newParsedData()
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{
		[]byte(basicAllowanceGrantJSON),
		[]byte(allowedMsgAllowanceGrantJSON),
	}
	err := module.parseFeeGrants(raw, testBlock(), &data)
	require.NoError(t, err)

	require.Len(t, data.grants, 2)
	require.Len(t, data.addresses, 4)
}

func TestParseFeeGrants_Empty(t *testing.T) {
	data := newParsedData()
	module := NewModule(postgres.Storage{}, config.Indexer{})

	err := module.parseFeeGrants(nil, testBlock(), &data)
	require.NoError(t, err)
	require.Empty(t, data.grants)
	require.Empty(t, data.addresses)
}

func TestParseFeeGrants_InvalidJSON(t *testing.T) {
	data := newParsedData()
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(`{"granter": "not valid`)}
	err := module.parseFeeGrants(raw, testBlock(), &data)
	require.Error(t, err)
}

func TestParseFeeGrants_ExistingAddressKeepsBalance(t *testing.T) {
	data := newParsedData()
	module := NewModule(postgres.Storage{}, config.Indexer{})

	const granterAddr = "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv"
	data.addresses[granterAddr] = &storage.Address{
		Address: granterAddr,
		Balances: []storage.Balance{
			{
				Currency:  currency.Utia,
				Spendable: storageTypes.NumericFromInt64(500),
			},
		},
	}

	raw := []json.RawMessage{[]byte(basicAllowanceGrantJSON)}
	err := module.parseFeeGrants(raw, testBlock(), &data)
	require.NoError(t, err)

	require.Equal(t, "500", data.addresses[granterAddr].Balances[0].Spendable.String())
}
