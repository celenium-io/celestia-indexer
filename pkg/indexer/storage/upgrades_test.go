// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	indexerCfg "github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var testValidators = []*storage.Validator{
	{
		Id:      1,
		Rate:    storageTypes.MustNumericFromString("0.150000000000000000"),
		MaxRate: storageTypes.MustNumericFromString("0.500000000000000000"),
	},
	{
		Id:      2,
		Rate:    storageTypes.MustNumericFromString("0.250000000000000000"),
		MaxRate: storageTypes.MustNumericFromString("0.700000000000000000"),
	},
}

func TestUpgradeV7(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validators := mock.NewMockIValidator(ctrl)
	validators.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(testValidators, nil).
		Times(1)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(nil, nil, validators, nil, indexerCfg.Indexer{Name: testIndexerName})
	dCtx := decodeContext.NewContext()

	err := module.upgradeV7(ctx, dCtx, 7)
	require.NoError(t, err)

	minCommissionRate := storageTypes.MustNumericFromString("0.200000000000000000")
	maxCommissionRate := storageTypes.MustNumericFromString("0.600000000000000000")

	for value := range dCtx.Validators.AllValues() {
		require.True(t, value.Rate.GreaterThanOrEqual(minCommissionRate))
		require.True(t, value.MaxRate.LessThanOrEqual(maxCommissionRate))
	}

	for value := range dCtx.Constants.AllValues() {
		if value.Name == "min_commission_rate" {
			require.Equal(t, "0.200000000000000000", value.Value)
		}
		if value.Name == "max_commission_rate" {
			require.Equal(t, "0.600000000000000000", value.Value)
		}
	}
}

// TestUpgrade_V8WithoutPriorV7 checks that upgrading to v8 when v7 was never applied
// also runs the v7 upgrade logic (validator commission adjustments + constants).
func TestUpgrade_V8WithoutPriorV7(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validators := mock.NewMockIValidator(ctrl)
	validators.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(testValidators, nil).
		Times(1)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(nil, nil, validators, nil, indexerCfg.Indexer{Name: testIndexerName})
	dCtx := decodeContext.NewContext()

	err := module.upgrade(ctx, dCtx, 6, 8)
	require.NoError(t, err)

	minCommissionRate := storageTypes.MustNumericFromString("0.200000000000000000")
	maxCommissionRate := storageTypes.MustNumericFromString("0.600000000000000000")

	for v := range dCtx.Validators.AllValues() {
		require.True(t, v.Rate.GreaterThanOrEqual(minCommissionRate))
		require.True(t, v.MaxRate.LessThanOrEqual(maxCommissionRate))
	}

	var foundMin, foundMax bool
	for c := range dCtx.Constants.AllValues() {
		if c.Name == "min_commission_rate" {
			require.Equal(t, "0.200000000000000000", c.Value)
			foundMin = true
		}
		if c.Name == "max_commission_rate" {
			require.Equal(t, "0.600000000000000000", c.Value)
			foundMax = true
		}
	}
	require.True(t, foundMin, "min_commission_rate constant must be set by v7 upgrade")
	require.True(t, foundMax, "max_commission_rate constant must be set by v7 upgrade")
}

// TestUpgrade_V8WithPriorV7 checks that when v7 was already applied (currentVersion=7),
// upgrading to v8 does NOT re-run the v7 logic.
func TestUpgrade_V8WithPriorV7(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validators := mock.NewMockIValidator(ctrl)
	// List must never be called — v7 logic must be skipped entirely.
	validators.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(nil, nil, validators, nil, indexerCfg.Indexer{Name: testIndexerName})
	dCtx := decodeContext.NewContext()

	err := module.upgrade(ctx, dCtx, 7, 8)
	require.NoError(t, err)

	require.Zero(t, dCtx.Validators.Len(), "no validators should be modified when v7 was already applied")
}

// TestUpgrade_NoOpWhenCurrentGTE checks the early-return guard.
func TestUpgrade_NoOpWhenCurrentGTE(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validators := mock.NewMockIValidator(ctrl)
	validators.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(nil, nil, validators, nil, indexerCfg.Indexer{Name: testIndexerName})
	dCtx := decodeContext.NewContext()

	err := module.upgrade(ctx, dCtx, 8, 8)
	require.NoError(t, err)

	err = module.upgrade(ctx, dCtx, 9, 8)
	require.NoError(t, err)
}

// TestUpgrade_V10SeedsFibreParams covers a v9 -> v10 upgrade: the chain writes
// nothing on chain for the new module, so the app defaults are seeded.
func TestUpgrade_V10SeedsFibreParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validators := mock.NewMockIValidator(ctrl)
	validators.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	constants := mock.NewMockIConstant(ctrl)
	constants.EXPECT().
		ByModule(gomock.Any(), storageTypes.ModuleNameFibre).
		Return(nil, nil).
		Times(1)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(nil, constants, validators, nil, indexerCfg.Indexer{Name: testIndexerName})
	dCtx := decodeContext.NewContext()

	err := module.upgrade(ctx, dCtx, 9, 10)
	require.NoError(t, err)

	got := make(map[string]string)
	for c := range dCtx.Constants.AllValues() {
		require.Equal(t, storageTypes.ModuleNameFibre, c.Module)
		got[c.Name] = c.Value
	}

	// Durations are in nanoseconds: 24h, 1h and 4h from the app defaults.
	require.Equal(t, map[string]string{
		"withdrawal_delay":              "86400000000000",
		"payment_promise_timeout":       "3600000000000",
		"payment_promise_height_window": "1000",
		"shard_retention":               "14400000000000",
		"full_stake_storage_budget":     "2199023255552",
	}, got)
}

// TestUpgrade_V10KeepsGenesisFibreParams covers a chain launched at v10: genesis
// already recorded the real params, so the defaults must not clobber them.
//
// currentVersion is 9 here, not 0: since the upgrade() fix (case 10 must chain
// through every earlier version), starting from 0 would also cascade through
// v6/v7 and add their staking constants + touch validators, which would
// muddy this test's one job of checking the fibre-params genesis guard. That
// cascading behavior is covered on its own by TestUpgrade_V8WithoutPriorV7.
func TestUpgrade_V10KeepsGenesisFibreParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validators := mock.NewMockIValidator(ctrl)
	validators.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	constants := mock.NewMockIConstant(ctrl)
	constants.EXPECT().
		ByModule(gomock.Any(), storageTypes.ModuleNameFibre).
		Return([]storage.Constant{{
			Module: storageTypes.ModuleNameFibre,
			Name:   "withdrawal_delay",
			Value:  "129600000000000",
		}}, nil).
		Times(1)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(nil, constants, validators, nil, indexerCfg.Indexer{Name: testIndexerName})
	dCtx := decodeContext.NewContext()

	err := module.upgrade(ctx, dCtx, 9, 10)
	require.NoError(t, err)
	require.Zero(t, dCtx.Constants.Len(), "genesis params must be left alone")
}

// TestUpgrade_V10ChainsThroughEarlierVersions covers a devnet/private-network
// jump straight from below v7 to v10 (e.g. an indexer started fresh after the
// chain already upgraded past v10): every intermediate upgrade must still run,
// not just the v10 one, matching the per-version loop in upgrade().
func TestUpgrade_V10ChainsThroughEarlierVersions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validators := mock.NewMockIValidator(ctrl)
	validators.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(testValidators, nil).
		Times(1)

	constants := mock.NewMockIConstant(ctrl)
	constants.EXPECT().
		ByModule(gomock.Any(), storageTypes.ModuleNameFibre).
		Return(nil, nil).
		Times(1)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(nil, constants, validators, nil, indexerCfg.Indexer{Name: testIndexerName})
	dCtx := decodeContext.NewContext()

	err := module.upgrade(ctx, dCtx, 6, 10)
	require.NoError(t, err)

	// v7's validator commission adjustment ran. testValidators share the
	// same (empty) Address, so AddValidator folds them into one entry.
	minCommissionRate := storageTypes.MustNumericFromString("0.200000000000000000")
	maxCommissionRate := storageTypes.MustNumericFromString("0.600000000000000000")
	require.NotZero(t, dCtx.Validators.Len())
	for v := range dCtx.Validators.AllValues() {
		require.True(t, v.Rate.GreaterThanOrEqual(minCommissionRate))
		require.True(t, v.MaxRate.LessThanOrEqual(maxCommissionRate))
	}

	// v7's staking constants and v10's fibre params both landed.
	got := make(map[string]string)
	for c := range dCtx.Constants.AllValues() {
		got[c.Name] = c.Value
	}
	require.Equal(t, "0.200000000000000000", got["min_commission_rate"])
	require.Equal(t, "0.600000000000000000", got["max_commission_rate"])
	require.Contains(t, got, "withdrawal_delay")
}
