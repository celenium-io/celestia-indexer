// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
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
		dec := decimal.NewFromInt(100)
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
		dec := decimal.NewFromInt(100)
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
		dec := decimal.NewFromInt(100)
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
