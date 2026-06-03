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

	err = dCtx.Validators.Range(func(key string, value *storage.Validator) (error, bool) {
		require.True(t, value.Rate.GreaterThanOrEqual(minCommissionRate))
		require.True(t, value.MaxRate.LessThanOrEqual(maxCommissionRate))
		return nil, true
	})
	require.NoError(t, err)

	err = dCtx.Constants.Range(func(_ string, value *storage.Constant) (error, bool) {
		if value.Name == "min_commission_rate" {
			require.Equal(t, "0.200000000000000000", value.Value)
		}
		if value.Name == "max_commission_rate" {
			require.Equal(t, "0.600000000000000000", value.Value)
		}
		return nil, true
	})
	require.NoError(t, err)
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

	err = dCtx.Validators.Range(func(_ string, v *storage.Validator) (error, bool) {
		require.True(t, v.Rate.GreaterThanOrEqual(minCommissionRate))
		require.True(t, v.MaxRate.LessThanOrEqual(maxCommissionRate))
		return nil, false
	})
	require.NoError(t, err)

	var foundMin, foundMax bool
	err = dCtx.Constants.Range(func(_ string, c *storage.Constant) (error, bool) {
		if c.Name == "min_commission_rate" {
			require.Equal(t, "0.200000000000000000", c.Value)
			foundMin = true
		}
		if c.Name == "max_commission_rate" {
			require.Equal(t, "0.600000000000000000", c.Value)
			foundMax = true
		}
		return nil, false
	})
	require.NoError(t, err)
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

	count := 0
	_ = dCtx.Validators.Range(func(_ string, _ *storage.Validator) (error, bool) {
		count++
		return nil, true
	})
	require.Zero(t, count, "no validators should be modified when v7 was already applied")
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
