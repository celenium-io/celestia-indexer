// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
)

func genesisValidators(stakes map[string]int64) *decodeContext.Context {
	ctx := decodeContext.NewContext()
	for address, stake := range stakes {
		v := storage.EmptyValidator()
		v.Address = address
		v.Stake = types.NumericFromInt64(stake)
		ctx.AddValidator(v)
	}
	return ctx
}

func powerOf(t *testing.T, ctx *decodeContext.Context, address string) *types.Numeric {
	t.Helper()
	v, ok := ctx.Validators.Get(address)
	require.True(t, ok)
	return v.Power
}

func requirePower(t *testing.T, ctx *decodeContext.Context, address, expected string) {
	t.Helper()
	power := powerOf(t, ctx, address)
	require.NotNil(t, power, address)
	require.Equal(t, expected, power.String(), address)
}

// Genesis usually has fewer gentxs than max_validators.
func TestParseValidatorsPower_FewerThanMax(t *testing.T) {
	ctx := genesisValidators(map[string]int64{"a": 5_000_000, "b": 2_999_999})

	require.NotPanics(t, func() { parseValidatorsPower(ctx, 100) })

	requirePower(t, ctx, "a", "5")
	requirePower(t, ctx, "b", "2")
}

// Only the top max_validators by stake are bonded.
func TestParseValidatorsPower_TopByStake(t *testing.T) {
	ctx := genesisValidators(map[string]int64{
		"a": 1_000_000,
		"b": 9_000_000,
		"c": 3_000_000,
		"d": 7_000_000,
		"e": 5_000_000,
	})

	parseValidatorsPower(ctx, 2)

	requirePower(t, ctx, "b", "9")
	requirePower(t, ctx, "d", "7")
	for _, address := range []string{"a", "c", "e"} {
		power := powerOf(t, ctx, address)
		require.True(t, power == nil || power.IsZero(), address)
	}
}

func TestParseValidatorsPower_NoValidators(t *testing.T) {
	ctx := decodeContext.NewContext()
	require.NotPanics(t, func() { parseValidatorsPower(ctx, 100) })
}
