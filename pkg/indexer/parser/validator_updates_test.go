// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	dCtx "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/require"
)

// Consensus addresses are sha256(pubkey)[:20] of the keys below.
const (
	testValPubKey1   = "AQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eHyA="
	testValConsAddr1 = "AE216C2EF5247A3782C135EFA279A3E4CDC61094"
	testValPubKey2   = "AgMEBQYHCAkKCwwNDg8QERITFBUWFxgZGhscHR4fICE="
	testValConsAddr2 = "6F103F3D9BA4C7E4D49642FB221098B83BCF07AC"
)

func validatorUpdate(t *testing.T, pubKey string, power *string) types.ValidatorUpdate {
	t.Helper()
	key, err := base64.StdEncoding.DecodeString(pubKey)
	require.NoError(t, err)

	var update types.ValidatorUpdate
	update.PubKey.Sum.Type = "tendermint.crypto.PublicKey_Ed25519"
	update.PubKey.Sum.Value.Ed25519 = key
	update.Power = power
	return update
}

func ptr[T any](v T) *T { return &v }

// parse sets the block before validator updates are parsed.
func newBlockContext() *dCtx.Context {
	ctx := dCtx.NewContext()
	ctx.Block = &storage.Block{Height: 100, Time: time.Unix(1_700_000_000, 0).UTC()}
	return ctx
}

func updatesByConsAddress(ctx *dCtx.Context) map[string]*storage.Validator {
	result := make(map[string]*storage.Validator)
	for _, v := range ctx.ValidatorUpdates.Values() {
		result[v.Validator.ConsAddress] = v.Validator
	}
	return result
}

func TestParseValidatorUpdates_Empty(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	require.NoError(t, module.parseValidatorUpdates(ctx, nil))
	require.Equal(t, 0, ctx.ValidatorUpdates.Len())
}

func TestParseValidatorUpdates_Single(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	err := module.parseValidatorUpdates(ctx, []types.ValidatorUpdate{
		validatorUpdate(t, testValPubKey1, ptr("1234")),
	})
	require.NoError(t, err)

	vals := updatesByConsAddress(ctx)
	require.Len(t, vals, 1)
	val, ok := vals[testValConsAddr1]
	require.True(t, ok)
	require.NotNil(t, val.Power)
	require.Equal(t, "1234", val.Power.String())
	require.Nil(t, val.Jailed)
	require.Equal(t, storage.DoNotModify, val.Moniker)
	require.True(t, val.Stake.IsZero())
	require.Empty(t, val.Address)

	// id and operator address are resolved later by the storage module
	require.Equal(t, 0, ctx.Validators.Len())
}

// Power 0 is omitted by cmtjson and means the validator left the active set.
func TestParseValidatorUpdates_MissingPowerIsZero(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	err := module.parseValidatorUpdates(ctx, []types.ValidatorUpdate{
		validatorUpdate(t, testValPubKey1, nil),
	})
	require.NoError(t, err)

	val, ok := updatesByConsAddress(ctx)[testValConsAddr1]
	require.True(t, ok)
	require.NotNil(t, val.Power)
	require.True(t, val.Power.IsZero())
}

func TestParseValidatorUpdates_Multiple(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	err := module.parseValidatorUpdates(ctx, []types.ValidatorUpdate{
		validatorUpdate(t, testValPubKey1, ptr("10")),
		validatorUpdate(t, testValPubKey2, nil),
	})
	require.NoError(t, err)

	vals := updatesByConsAddress(ctx)
	require.Len(t, vals, 2)
	require.Equal(t, "10", vals[testValConsAddr1].Power.String())
	require.True(t, vals[testValConsAddr2].Power.IsZero())
}

// A repeated key in one block keeps the last power.
func TestParseValidatorUpdates_SameKeyTwice(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	err := module.parseValidatorUpdates(ctx, []types.ValidatorUpdate{
		validatorUpdate(t, testValPubKey1, ptr("10")),
		validatorUpdate(t, testValPubKey1, nil),
	})
	require.NoError(t, err)

	vals := updatesByConsAddress(ctx)
	require.Len(t, vals, 1)
	require.True(t, vals[testValConsAddr1].Power.IsZero())
}

func TestParseValidatorUpdates_InvalidPower(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	err := module.parseValidatorUpdates(ctx, []types.ValidatorUpdate{
		validatorUpdate(t, testValPubKey1, ptr("abc")),
	})
	require.Error(t, err)
}

// Unsupported key types are skipped with a warning.
func TestParseValidatorUpdates_NonEd25519KeySkipped(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	var update types.ValidatorUpdate
	update.PubKey.Sum.Type = "tendermint.crypto.PublicKey_Secp256K1"
	update.Power = ptr("1")

	err := module.parseValidatorUpdates(ctx, []types.ValidatorUpdate{
		update,
		validatorUpdate(t, testValPubKey1, ptr("2")),
	})
	require.NoError(t, err)

	vals := updatesByConsAddress(ctx)
	require.Len(t, vals, 1)
	require.Equal(t, "2", vals[testValConsAddr1].Power.String())
}

func TestParseValidatorUpdates_EmptyEd25519KeySkipped(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	update := validatorUpdate(t, testValPubKey1, ptr("1"))
	update.PubKey.Sum.Value.Ed25519 = nil

	require.NoError(t, module.parseValidatorUpdates(ctx, []types.ValidatorUpdate{update}))
	require.Equal(t, 0, ctx.ValidatorUpdates.Len())
}

// A key of the wrong size is skipped instead of panicking inside ed25519.PubKey.Address.
func TestParseValidatorUpdates_WrongSizeEd25519KeySkipped(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	update := validatorUpdate(t, testValPubKey1, ptr("1"))
	update.PubKey.Sum.Value.Ed25519 = update.PubKey.Sum.Value.Ed25519[:31]

	require.NotPanics(t, func() {
		err := module.parseValidatorUpdates(ctx, []types.ValidatorUpdate{update})
		require.NoError(t, err)
	})
	require.Equal(t, 0, ctx.ValidatorUpdates.Len())
}

// Updates carry block height and time and never create a validator with an empty operator address.
func TestParseValidatorUpdates_BlockFieldsAndNoValidators(t *testing.T) {
	module := NewModule(config.Indexer{})
	ctx := newBlockContext()

	err := module.parseValidatorUpdates(ctx, []types.ValidatorUpdate{
		validatorUpdate(t, testValPubKey1, ptr("10")),
		validatorUpdate(t, testValPubKey2, nil),
	})
	require.NoError(t, err)
	require.Equal(t, 0, ctx.Validators.Len())

	for _, update := range ctx.ValidatorUpdates.Values() {
		require.EqualValues(t, 100, update.Height)
		require.Equal(t, ctx.Block.Time, update.Time)
		require.NotNil(t, update.Power)
		require.Equal(t, update.Power.String(), update.Validator.Power.String())
		require.EqualValues(t, 1, update.Validator.BondUpdatesCount)
		require.Zero(t, update.ValidatorId)
	}
}
