// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestNewValidatorUptime(t *testing.T) {
	t.Run("validator with all levels", func(t *testing.T) {
		var (
			levels       = []types.Level{4, 3, 2, 1}
			currentLevel = types.Level(4)
			count        = types.Level(3)
		)
		uptime := NewValidatorUptime(levels, currentLevel, count)
		require.EqualValues(t, "1.0000", uptime.Uptime)
		require.Len(t, uptime.Blocks, 3)

		for i := range uptime.Blocks {
			require.True(t, uptime.Blocks[i].Signed)
		}
	})

	t.Run("validator with skipped levels", func(t *testing.T) {
		var (
			levels       = []types.Level{4, 1}
			currentLevel = types.Level(4)
			count        = types.Level(3)
		)
		uptime := NewValidatorUptime(levels, currentLevel, count)
		require.EqualValues(t, "0.3333", uptime.Uptime)
		require.Len(t, uptime.Blocks, 3)

		require.True(t, uptime.Blocks[0].Signed)
		require.False(t, uptime.Blocks[1].Signed)
		require.False(t, uptime.Blocks[2].Signed)
	})

	t.Run("current level less than requested count", func(t *testing.T) {
		var (
			levels       = []types.Level{4, 3, 2}
			currentLevel = types.Level(4)
			count        = types.Level(100)
		)
		uptime := NewValidatorUptime(levels, currentLevel, count)
		require.EqualValues(t, "0.7500", uptime.Uptime)
		require.Len(t, uptime.Blocks, 4)

		require.True(t, uptime.Blocks[0].Signed)
		require.True(t, uptime.Blocks[1].Signed)
		require.True(t, uptime.Blocks[2].Signed)
		require.False(t, uptime.Blocks[3].Signed)
	})
}

func TestNewValidator(t *testing.T) {
	t.Run("validator with nil jailed field", func(t *testing.T) {
		dec := storageTypes.NumericFromInt64(100)
		validator := storage.Validator{
			Jailed:            nil,
			Id:                1,
			Rate:              dec,
			MaxRate:           dec,
			MaxChangeRate:     dec,
			MinSelfDelegation: dec,
			Stake:             dec,
			Rewards:           dec,
			Commissions:       dec,
		}
		val := NewValidator(validator)
		require.False(t, val.Jailed)
	})

	t.Run("validator with false jailed field", func(t *testing.T) {
		dec := storageTypes.NumericFromInt64(100)
		validator := storage.Validator{
			Jailed:            testsuite.Ptr(false),
			Id:                1,
			Rate:              dec,
			MaxRate:           dec,
			MaxChangeRate:     dec,
			MinSelfDelegation: dec,
			Stake:             dec,
			Rewards:           dec,
			Commissions:       dec,
		}
		val := NewValidator(validator)
		require.False(t, val.Jailed)
	})

	t.Run("validator with true jailed field", func(t *testing.T) {
		dec := storageTypes.NumericFromInt64(100)
		validator := storage.Validator{
			Jailed:            testsuite.Ptr(true),
			Id:                1,
			Rate:              dec,
			MaxRate:           dec,
			MaxChangeRate:     dec,
			MinSelfDelegation: dec,
			Stake:             dec,
			Rewards:           dec,
			Commissions:       dec,
		}
		val := NewValidator(validator)
		require.True(t, val.Jailed)
	})
}

func TestNewValidatorStatus(t *testing.T) {
	power := func(v int64) *storageTypes.Numeric {
		n := storageTypes.NumericFromInt64(v)
		return &n
	}

	tests := []struct {
		name        string
		power       *storageTypes.Numeric
		jailed      *bool
		status      storageTypes.ValidatorStatus
		votingPower string
	}{
		{name: "bonded", power: power(10), jailed: testsuite.Ptr(false), status: storageTypes.ValidatorStatusActive, votingPower: "10"},
		{name: "bonded, unknown jailed flag", power: power(10), status: storageTypes.ValidatorStatusActive, votingPower: "10"},
		{name: "never bonded", jailed: testsuite.Ptr(false), status: storageTypes.ValidatorStatusNotActive},
		{name: "left active set", power: power(0), jailed: testsuite.Ptr(false), status: storageTypes.ValidatorStatusNotActive, votingPower: "0"},
		{name: "jailed", power: power(0), jailed: testsuite.Ptr(true), status: storageTypes.ValidatorStatusJailed, votingPower: "0"},
		{name: "jailed without power", jailed: testsuite.Ptr(true), status: storageTypes.ValidatorStatusJailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := storage.EmptyValidator()
			val.Id = 1
			val.Power = tt.power
			val.Jailed = tt.jailed
			val.BondUpdatesCount = 3

			resp := NewValidator(val)
			require.NotNil(t, resp)
			require.Equal(t, tt.status.String(), resp.Status)
			require.Equal(t, tt.votingPower, resp.VotingPower)
			require.EqualValues(t, 3, resp.BondUpdatesCount)
		})
	}
}

func TestNewBondUpdate(t *testing.T) {
	blockTime := time.Date(2023, 7, 4, 3, 10, 57, 0, time.UTC)
	power := storageTypes.NumericFromInt64(42)

	update := NewBondUpdate(storage.ValidatorBondUpdate{Id: 1, Height: 100, Time: blockTime, ValidatorId: 5, Power: &power})
	require.EqualValues(t, 100, update.Height)
	require.Equal(t, blockTime, update.Time)
	require.NotNil(t, update.Power)
	require.Equal(t, "42", *update.Power)

	zero := storageTypes.NumericZero()
	update = NewBondUpdate(storage.ValidatorBondUpdate{Height: 101, Power: &zero})
	require.NotNil(t, update.Power)
	require.Equal(t, "0", *update.Power)

	// no power is omitted from json instead of rendered as an empty string
	update = NewBondUpdate(storage.ValidatorBondUpdate{Height: 102})
	require.Nil(t, update.Power)
	raw, err := json.Marshal(update)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "power")
}
