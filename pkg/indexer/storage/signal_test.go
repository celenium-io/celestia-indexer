// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	indexerCfg "github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func makeUpgradesMap(upgrades ...*storage.Upgrade) *sync.Map[uint64, *storage.Upgrade] {
	m := sync.NewMap[uint64, *storage.Upgrade]()
	for _, u := range upgrades {
		m.Set(u.Version, u)
	}
	return m
}

func makeModule(ctrl *gomock.Controller, constants *mock.MockIConstant) (Module, *mock.MockTransaction) {
	tx := mock.NewMockTransaction(ctrl)
	constants.EXPECT().
		Get(gomock.Any(), types.ModuleNameStaking, "max_validators").
		Return(storage.Constant{Value: "100"}, nil).
		AnyTimes()

	m := NewModule(nil, constants, nil, nil, indexerCfg.Indexer{Name: testIndexerName})
	return m, tx
}

// ---------------------------------------------------------------------------
// tryUpgrade
// ---------------------------------------------------------------------------

func TestTryUpgrade_NilUpgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constants := mock.NewMockIConstant(ctrl)
	module, tx := makeModule(ctrl, constants)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := module.tryUpgrade(ctx, tx, nil, storage.State{Version: 3})
	require.NoError(t, err)
}

func TestTryUpgrade_NoValidatorsSignaled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constants := mock.NewMockIConstant(ctrl)
	module, tx := makeModule(ctrl, constants)

	// all validators signal version 0 (no version) or <= state.Version
	tx.EXPECT().BondedValidators(gomock.Any(), 100).Return([]storage.Validator{
		{Id: 1, Stake: types.NewNumeric(decimal.NewFromInt(1_000_000)), Version: 0},
		{Id: 2, Stake: types.NewNumeric(decimal.NewFromInt(1_000_000)), Version: 3},
	}, nil)
	// UpdateSignalsAfterUpgrade must NOT be called

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upgrade := &storage.Upgrade{Height: 100, Time: time.Now()}
	err := module.tryUpgrade(ctx, tx, upgrade, storage.State{Version: 3})
	require.NoError(t, err)
}

func TestTryUpgrade_NoQuorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constants := mock.NewMockIConstant(ctrl)
	module, tx := makeModule(ctrl, constants)

	// total stake = 3_000_000 → Shares = 3; threshold = 3 * 5/6 ≈ 2
	// voted for v4 = 1_000_000 → Shares = 1 < threshold → no quorum
	tx.EXPECT().BondedValidators(gomock.Any(), 100).Return([]storage.Validator{
		{Id: 1, Stake: types.NewNumeric(decimal.NewFromInt(1_000_000)), Version: 4},
		{Id: 2, Stake: types.NewNumeric(decimal.NewFromInt(1_000_000)), Version: 3},
		{Id: 3, Stake: types.NewNumeric(decimal.NewFromInt(1_000_000)), Version: 3},
	}, nil)
	tx.EXPECT().UpdateSignalsAfterUpgrade(gomock.Any(), uint64(4)).
		Return(types.NewNumeric(decimal.NewFromInt(1_000_000)), nil)
	// SaveUpgrades must NOT be called

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upgrade := &storage.Upgrade{Height: 100, Time: time.Now()}
	err := module.tryUpgrade(ctx, tx, upgrade, storage.State{Version: 3})
	require.NoError(t, err)
}

func TestTryUpgrade_WithQuorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constants := mock.NewMockIConstant(ctrl)
	module, tx := makeModule(ctrl, constants)

	// total stake = 6_000_000 → Shares = 6; threshold = 6 * 5/6 = 5
	// voted for v4 raw = 6_000_000 → Shares = 6 > 5 → quorum
	tx.EXPECT().BondedValidators(gomock.Any(), 100).Return([]storage.Validator{
		{Id: 1, Stake: types.NewNumeric(decimal.NewFromInt(2_000_000)), Version: 4},
		{Id: 2, Stake: types.NewNumeric(decimal.NewFromInt(2_000_000)), Version: 4},
		{Id: 3, Stake: types.NewNumeric(decimal.NewFromInt(2_000_000)), Version: 4},
	}, nil)
	tx.EXPECT().UpdateSignalsAfterUpgrade(gomock.Any(), uint64(4)).
		Return(types.NewNumeric(decimal.NewFromInt(6_000_000)), nil)
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.EqualValues(t, 4, upgrades[0].Version)
			require.Equal(t, types.UpgradeStatusWaitingUpgrade, upgrades[0].Status)
			return nil
		})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upgrade := &storage.Upgrade{Height: 100, Time: time.Now()}
	err := module.tryUpgrade(ctx, tx, upgrade, storage.State{Version: 3})
	require.NoError(t, err)
}

func TestTryUpgrade_PicksMinimumQuorumVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constants := mock.NewMockIConstant(ctrl)
	module, tx := makeModule(ctrl, constants)

	// two versions both have quorum; should pick v4 (minimum), not v5
	// total stake = 12_000_000 → Shares = 12; threshold = 12 * 5/6 = 10
	// voted raw = 12_000_000 → Shares = 12 > 10 → quorum for both
	tx.EXPECT().BondedValidators(gomock.Any(), 100).Return([]storage.Validator{
		{Id: 1, Stake: types.NewNumeric(decimal.NewFromInt(3_000_000)), Version: 4},
		{Id: 2, Stake: types.NewNumeric(decimal.NewFromInt(3_000_000)), Version: 4},
		{Id: 3, Stake: types.NewNumeric(decimal.NewFromInt(3_000_000)), Version: 4},
		{Id: 4, Stake: types.NewNumeric(decimal.NewFromInt(3_000_000)), Version: 5},
	}, nil)
	// versions are sorted, v4 is checked first and has quorum → SaveUpgrades called, v5 skipped
	tx.EXPECT().UpdateSignalsAfterUpgrade(gomock.Any(), uint64(4)).
		Return(types.NewNumeric(decimal.NewFromInt(12_000_000)), nil).MaxTimes(1)
	tx.EXPECT().UpdateSignalsAfterUpgrade(gomock.Any(), uint64(5)).
		Return(types.NewNumeric(decimal.NewFromInt(12_000_000)), nil).MaxTimes(1)
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.EqualValues(t, 4, upgrades[0].Version)
			return nil
		})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	upgrade := &storage.Upgrade{Height: 100, Time: time.Now()}
	err := module.tryUpgrade(ctx, tx, upgrade, storage.State{Version: 3})
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// saveUpgrades
// ---------------------------------------------------------------------------

func TestSaveUpgrades_EmptyMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	// nothing should be called
	err := saveUpgrades(context.Background(), tx, sync.NewMap[uint64, *storage.Upgrade](),
		storage.State{Version: 3}, types.NewNumeric(decimal.NewFromInt(3_000_000)))
	require.NoError(t, err)
}

func TestSaveUpgrades_SkipsAlreadyAppliedVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	// version 3 <= state.Version 3 → skipped; UpdateSignalsAfterUpgrade must NOT be called

	upgrades := makeUpgradesMap(&storage.Upgrade{Version: 3})

	err := saveUpgrades(context.Background(), tx, upgrades,
		storage.State{Version: 3}, types.NewNumeric(decimal.NewFromInt(3_000_000)))
	require.NoError(t, err)
}

func TestSaveUpgrades_NoQuorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	// total Shares = 3; threshold = 2; voted raw=1_000_000 → Shares=1 < 2
	tx.EXPECT().UpdateSignalsAfterUpgrade(gomock.Any(), uint64(4)).
		Return(types.NewNumeric(decimal.NewFromInt(1_000_000)), nil)
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.NotEqual(t, types.UpgradeStatusWaitingUpgrade, upgrades[0].Status)
			return nil
		})

	upgrades := makeUpgradesMap(&storage.Upgrade{Version: 4})

	err := saveUpgrades(context.Background(), tx, upgrades,
		storage.State{Version: 3}, types.NewNumeric(decimal.NewFromInt(3)))
	require.NoError(t, err)
}

func TestSaveUpgrades_WithQuorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	// total Shares = 6; threshold = 5; voted raw=6_000_000 → Shares=6 > 5 → quorum
	tx.EXPECT().UpdateSignalsAfterUpgrade(gomock.Any(), uint64(4)).
		Return(types.NewNumeric(decimal.NewFromInt(6_000_000)), nil)
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.Equal(t, types.UpgradeStatusWaitingUpgrade, upgrades[0].Status)
			require.True(t, upgrades[0].VotedPower.Equal(types.NumericFromInt64(6)))
			return nil
		})

	upgrades := makeUpgradesMap(&storage.Upgrade{Version: 4})

	err := saveUpgrades(context.Background(), tx, upgrades,
		storage.State{Version: 3}, types.NewNumeric(decimal.NewFromInt(6)))
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// saveSignals
// ---------------------------------------------------------------------------

func TestSaveSignals_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constants := mock.NewMockIConstant(ctrl)
	module := NewModule(nil, constants, nil, nil, indexerCfg.Indexer{Name: testIndexerName})

	err := module.saveSignals(context.Background(), nil, nil,
		sync.NewMap[uint64, *storage.Upgrade](), storage.State{})
	require.NoError(t, err)
}

func TestSaveSignals_WithQuorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constants := mock.NewMockIConstant(ctrl)
	constants.EXPECT().
		Get(gomock.Any(), types.ModuleNameStaking, "max_validators").
		Return(storage.Constant{Value: "100"}, nil)

	tx := mock.NewMockTransaction(ctrl)

	// total bonded stake
	tx.EXPECT().BondedValidators(gomock.Any(), 100).Return([]storage.Validator{
		{Id: 1, Stake: types.NewNumeric(decimal.NewFromInt(2_000_000))},
		{Id: 2, Stake: types.NewNumeric(decimal.NewFromInt(1_000_000))},
	}, nil)

	// saveSignals resolves validator by address
	tx.EXPECT().Validator(gomock.Any(), uint64(1)).Return(storage.Validator{
		Id: 1, Stake: types.NewNumeric(decimal.NewFromInt(2_000_000)),
	}, nil)

	tx.EXPECT().SaveSignals(gomock.Any(), gomock.Any()).Return(nil)

	// total Shares = 3; threshold = 2; voted raw=6_000_000 → Shares=6 > 2 → quorum
	tx.EXPECT().UpdateSignalsAfterUpgrade(gomock.Any(), uint64(4)).
		Return(types.NewNumeric(decimal.NewFromInt(6_000_000)), nil)
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.Equal(t, types.UpgradeStatusWaitingUpgrade, upgrades[0].Status)
			return nil
		})

	module := NewModule(nil, constants, nil, nil, indexerCfg.Indexer{Name: testIndexerName})
	module.validatorsByAddress["val1address"] = 1

	signals := []*storage.SignalVersion{
		{
			Version:   4,
			Height:    100,
			Validator: &storage.Validator{Address: "val1address"},
		},
	}
	upgrades := makeUpgradesMap(&storage.Upgrade{Version: 4})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := module.saveSignals(ctx, tx, signals, upgrades, storage.State{Version: 3})
	require.NoError(t, err)
}

func TestSaveSignals_UnknownValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constants := mock.NewMockIConstant(ctrl)
	constants.EXPECT().
		Get(gomock.Any(), types.ModuleNameStaking, "max_validators").
		Return(storage.Constant{Value: "100"}, nil)

	tx := mock.NewMockTransaction(ctrl)
	tx.EXPECT().BondedValidators(gomock.Any(), 100).Return([]storage.Validator{}, nil)

	module := NewModule(nil, constants, nil, nil, indexerCfg.Indexer{Name: testIndexerName})
	// validatorsByAddress intentionally empty

	signals := []*storage.SignalVersion{
		{
			Version:   4,
			Height:    100,
			Validator: &storage.Validator{Address: "unknown_address"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := module.saveSignals(ctx, tx, signals,
		sync.NewMap[uint64, *storage.Upgrade](), storage.State{Version: 3})
	require.Error(t, err)
}
