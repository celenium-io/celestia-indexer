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
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
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

	genericAuthzGrantJSON = `{
		"granter": "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		"grantee": "celestia1llx8v3hhw06cmtmv7q7uq4ztgyvd33nmk7mrfh",
		"authorization": {
			"@type": "/cosmos.authz.v1beta1.GenericAuthorization",
			"msg": "/cosmos.gov.v1.MsgVote"
		},
		"expiration": null
	}`

	sendAuthzGrantJSON = `{
		"granter": "celestia1gezy6mf5hy5nmq445djq3s3mwx7pl2xe7r3kp6",
		"grantee": "celestia1lamchvjp7c7d9asy8g9cxuycfa4qf4t6c3ufeg",
		"authorization": {
			"@type": "/cosmos.bank.v1beta1.SendAuthorization",
			"spend_limit": [],
			"allow_list": []
		},
		"expiration": null
	}`

	expiredSendAuthzGrantJSON = `{
		"granter": "celestia1gezy6mf5hy5nmq445djq3s3mwx7pl2xe7r3kp6",
		"grantee": "celestia1lamchvjp7c7d9asy8g9cxuycfa4qf4t6c3ufeg",
		"authorization": {
			"@type": "/cosmos.bank.v1beta1.SendAuthorization",
			"spend_limit": [],
			"allow_list": []
		},
		"expiration": "2000-01-01T00:00:00Z"
	}`

	nilAuthorizationGrantJSON = `{
		"granter": "celestia1gezy6mf5hy5nmq445djq3s3mwx7pl2xe7r3kp6",
		"grantee": "celestia1lamchvjp7c7d9asy8g9cxuycfa4qf4t6c3ufeg",
		"expiration": null
	}`

	malformedAddressAuthzGrantJSON = `{
		"granter": "not-a-valid-address",
		"grantee": "celestia1lamchvjp7c7d9asy8g9cxuycfa4qf4t6c3ufeg",
		"authorization": {
			"@type": "/cosmos.bank.v1beta1.SendAuthorization",
			"spend_limit": [],
			"allow_list": []
		},
		"expiration": null
	}`

	unrecognizedTypeAuthzGrantJSON = `{
		"granter": "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		"grantee": "celestia1llx8v3hhw06cmtmv7q7uq4ztgyvd33nmk7mrfh",
		"authorization": {
			"@type": "/ibc.applications.transfer.v1.TransferAuthorization",
			"allocations": []
		},
		"expiration": null
	}`
)

func stakeAuthzGrantJSON(authorizationType string) string {
	return `{
		"granter": "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		"grantee": "celestia1llx8v3hhw06cmtmv7q7uq4ztgyvd33nmk7mrfh",
		"authorization": {
			"@type": "/cosmos.staking.v1beta1.StakeAuthorization",
			"authorization_type": "` + authorizationType + `"
		},
		"expiration": null
	}`
}

func testBlock() storage.Block {
	return storage.Block{
		Height: 1,
		Time:   time.Now().UTC(),
	}
}

func TestParseFeeGrants_BasicAllowance(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(basicAllowanceGrantJSON)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseFeeGrants(ctx, raw)
	require.NoError(t, err)

	grants := ctx.Grants.Values()
	require.Len(t, grants, 1)
	grant := grants[0]

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

	_, ok = ctx.Addresses.Get(grant.Granter.Address)
	require.True(t, ok)
	_, ok = ctx.Addresses.Get(grant.Grantee.Address)
	require.True(t, ok)

	granterAddr, ok := ctx.Addresses.Get(grant.Granter.Address)
	require.True(t, ok)
	require.Len(t, granterAddr.Balances, 1)
	require.True(t, granterAddr.Balances[0].Spendable.IsZero())
	require.Equal(t, currency.Utia, granterAddr.Balances[0].Currency)
}

func TestParseFeeGrants_AllowedMsgAllowance(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(allowedMsgAllowanceGrantJSON)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseFeeGrants(ctx, raw)
	require.NoError(t, err)

	grants := ctx.Grants.Values()
	require.Len(t, grants, 1)
	grant := grants[0]

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

	_, ok = ctx.Addresses.Get(grant.Granter.Address)
	require.True(t, ok)
	_, ok = ctx.Addresses.Get(grant.Grantee.Address)
	require.True(t, ok)
}

func TestParseFeeGrants_Multiple(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{
		[]byte(basicAllowanceGrantJSON),
		[]byte(allowedMsgAllowanceGrantJSON),
	}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseFeeGrants(ctx, raw)
	require.NoError(t, err)

	require.Len(t, ctx.Grants.Values(), 2)
	require.Equal(t, 4, ctx.Addresses.Len())
}

func TestParseFeeGrants_Empty(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseFeeGrants(ctx, nil)
	require.NoError(t, err)
	require.Empty(t, ctx.Grants.Values())
	require.Zero(t, ctx.Addresses.Len())
}

func TestParseFeeGrants_InvalidJSON(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(`{"granter": "not valid`)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseFeeGrants(ctx, raw)
	require.Error(t, err)
}

func TestParseFeeGrants_ExistingAddressKeepsBalance(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	const granterAddr = "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv"
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	ctx.Addresses.Set(granterAddr, &storage.Address{
		Address: granterAddr,
		Balances: []storage.Balance{
			{
				Currency:  currency.Utia,
				Spendable: storageTypes.NumericFromInt64(500),
			},
		},
	})

	raw := []json.RawMessage{[]byte(basicAllowanceGrantJSON)}
	err := module.parseFeeGrants(ctx, raw)
	require.NoError(t, err)

	granterAddress, ok := ctx.Addresses.Get(granterAddr)
	require.True(t, ok)
	require.Equal(t, "500", granterAddress.Balances[0].Spendable.String())
}

func TestParseAuthzGrants_GenericAuthorization(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(genericAuthzGrantJSON)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, raw)
	require.NoError(t, err)

	grants := ctx.Grants.Values()
	require.Len(t, grants, 1)
	grant := grants[0]

	require.Equal(t, "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv", grant.Granter.Address)
	require.Equal(t, "celestia1llx8v3hhw06cmtmv7q7uq4ztgyvd33nmk7mrfh", grant.Grantee.Address)
	require.Equal(t, "/cosmos.gov.v1.MsgVote", grant.Authorization)
	require.NotNil(t, grant.Params)
	require.False(t, grant.Revoked)
	require.Nil(t, grant.Expiration)

	_, ok := ctx.Addresses.Get(grant.Granter.Address)
	require.True(t, ok)
	_, ok = ctx.Addresses.Get(grant.Grantee.Address)
	require.True(t, ok)
}

func TestParseAuthzGrants_SendAuthorization(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(sendAuthzGrantJSON)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, raw)
	require.NoError(t, err)

	grants := ctx.Grants.Values()
	require.Len(t, grants, 1)
	grant := grants[0]

	require.Equal(t, "celestia1gezy6mf5hy5nmq445djq3s3mwx7pl2xe7r3kp6", grant.Granter.Address)
	require.Equal(t, "celestia1lamchvjp7c7d9asy8g9cxuycfa4qf4t6c3ufeg", grant.Grantee.Address)
	require.Equal(t, "/cosmos.bank.v1beta1.MsgSend", grant.Authorization)
	require.NotNil(t, grant.Params)

	_, ok := ctx.Addresses.Get(grant.Granter.Address)
	require.True(t, ok)
	_, ok = ctx.Addresses.Get(grant.Grantee.Address)
	require.True(t, ok)
}

func TestParseAuthzGrants_StakeAuthorization(t *testing.T) {
	cases := []struct {
		name              string
		authorizationType string
		wantAuthorization []string
	}{
		{"Delegate", "AUTHORIZATION_TYPE_DELEGATE", []string{"/cosmos.staking.v1beta1.MsgDelegate"}},
		{"Redelegate", "AUTHORIZATION_TYPE_REDELEGATE", []string{"/cosmos.staking.v1beta1.MsgRedelegate"}},
		{"Undelegate", "AUTHORIZATION_TYPE_UNDELEGATE", []string{"/cosmos.staking.v1beta1.MsgUndelegate"}},
		{"CancelUnbondingDelegation", "AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION", []string{"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation"}},
		{"Unspecified", "AUTHORIZATION_TYPE_UNSPECIFIED", []string{
			"/cosmos.staking.v1beta1.MsgDelegate",
			"/cosmos.staking.v1beta1.MsgRedelegate",
			"/cosmos.staking.v1beta1.MsgUndelegate",
			"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation",
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			module := NewModule(postgres.Storage{}, config.Indexer{})

			raw := []json.RawMessage{[]byte(stakeAuthzGrantJSON(tc.authorizationType))}
			block := testBlock()
			ctx := decodeContext.NewContext()
			ctx.Block = &block
			err := module.parseAuthzGrants(ctx, raw)
			require.NoError(t, err)

			grants := ctx.Grants.Values()
			require.Len(t, grants, len(tc.wantAuthorization))
			got := make([]string, len(grants))
			for i, g := range grants {
				got[i] = g.Authorization
			}
			require.ElementsMatch(t, tc.wantAuthorization, got)
		})
	}
}

func TestParseAuthzGrants_Expiration(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(expiredSendAuthzGrantJSON)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, raw)
	require.NoError(t, err)

	grants := ctx.Grants.Values()
	require.Len(t, grants, 1)
	grant := grants[0]

	require.NotNil(t, grant.Expiration)
	require.True(t, grant.Expiration.Before(time.Now()), "expiration is in the past and must not be filtered out")
	require.False(t, grant.Revoked)
}

func TestParseAuthzGrants_NilAuthorization(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(nilAuthorizationGrantJSON)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, raw)
	require.Error(t, err)
}

// An unrecognized authorization type falls back to a generic grant (unknown_type_url,
// no decoded Params) rather than being skipped or erroring out.
func TestParseAuthzGrants_UnrecognizedType(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(unrecognizedTypeAuthzGrantJSON)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, raw)
	require.NoError(t, err)

	grants := ctx.Grants.Values()
	require.Len(t, grants, 1)
	grant := grants[0]
	require.Equal(t, "unknown_type_url", grant.Authorization)
	require.Nil(t, grant.Params)
	require.Equal(t, "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv", grant.Granter.Address)
	require.Equal(t, "celestia1llx8v3hhw06cmtmv7q7uq4ztgyvd33nmk7mrfh", grant.Grantee.Address)
}

// Repeating the same grant JSON must not produce two grant rows: ctx.Grants dedups by
// grant.String(), same key both times.
func TestParseAuthzGrants_DuplicateKnownType(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(genericAuthzGrantJSON), []byte(genericAuthzGrantJSON)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, raw)
	require.NoError(t, err)
	require.Len(t, ctx.Grants.Values(), 1)
}

func TestParseAuthzGrants_Multiple(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{
		[]byte(genericAuthzGrantJSON),
		[]byte(sendAuthzGrantJSON),
	}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, raw)
	require.NoError(t, err)

	require.Len(t, ctx.Grants.Values(), 2)
	require.Equal(t, 4, ctx.Addresses.Len())
}

func TestParseAuthzGrants_Empty(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, nil)
	require.NoError(t, err)
	require.Empty(t, ctx.Grants.Values())
	require.Zero(t, ctx.Addresses.Len())
}

func TestParseAuthzGrants_InvalidJSON(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(`{"granter": "not valid`)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, raw)
	require.Error(t, err)
}

func TestParseAuthzGrants_MalformedAddress(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	raw := []json.RawMessage{[]byte(malformedAddressAuthzGrantJSON)}
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	err := module.parseAuthzGrants(ctx, raw)
	require.Error(t, err)
}

func TestParseAuthzGrants_ExistingAddressKeepsBalance(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	const granterAddr = "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv"
	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block
	ctx.Addresses.Set(granterAddr, &storage.Address{
		Address: granterAddr,
		Balances: []storage.Balance{
			{
				Currency:  currency.Utia,
				Spendable: storageTypes.NumericFromInt64(500),
			},
		},
	})

	raw := []json.RawMessage{[]byte(genericAuthzGrantJSON)}
	err := module.parseAuthzGrants(ctx, raw)
	require.NoError(t, err)

	granterAddress, ok := ctx.Addresses.Get(granterAddr)
	require.True(t, ok)
	require.Equal(t, "500", granterAddress.Balances[0].Spendable.String())
}

// Regression guard: feegrant and authz grants share ctx.Grants; neither parse call
// may clobber entries the other already wrote.
func TestParseGrants_MixedFeegrantAndAuthz_NoClobber(t *testing.T) {
	module := NewModule(postgres.Storage{}, config.Indexer{})

	block := testBlock()
	ctx := decodeContext.NewContext()
	ctx.Block = &block

	feeRaw := []json.RawMessage{[]byte(basicAllowanceGrantJSON)}
	err := module.parseFeeGrants(ctx, feeRaw)
	require.NoError(t, err)

	authzRaw := []json.RawMessage{[]byte(sendAuthzGrantJSON)}
	err = module.parseAuthzGrants(ctx, authzRaw)
	require.NoError(t, err)

	grants := ctx.Grants.Values()
	require.Len(t, grants, 2)
	var haveFee, haveAuthz bool
	for _, g := range grants {
		switch g.Authorization {
		case "fee":
			haveFee = true
		case "/cosmos.bank.v1beta1.MsgSend":
			haveAuthz = true
		}
	}
	require.True(t, haveFee, "feegrant entry must be present")
	require.True(t, haveAuthz, "authz entry must be present")
}
