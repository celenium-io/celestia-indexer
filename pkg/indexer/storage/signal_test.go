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
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func makeUpgradesMap(upgrades ...*storage.Upgrade) *sdkSync.Map[uint64, *storage.Upgrade] {
	m := sdkSync.NewMap[uint64, *storage.Upgrade]()
	for _, u := range upgrades {
		m.Set(u.Version, u)
	}
	return m
}

// ---------------------------------------------------------------------------
// tryUpgrade
// ---------------------------------------------------------------------------

func TestTryUpgrade_NilUpgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	err := tryUpgrade(t.Context(), tx, repos.Repos(), nil, 3, "", nil)
	require.NoError(t, err)
}

func TestTryUpgrade_NoValidatorsSignaled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	// all validators signal version 0 (no version) or <= current version: nothing is tallied
	repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{
		{Id: 1, Version: 0},
		{Id: 2, Version: 3},
	}, nil)

	upgrade := &storage.Upgrade{Height: 100, Time: time.Now()}
	err := tryUpgrade(t.Context(), tx, repos.Repos(), upgrade, 3, "", nil)
	require.NoError(t, err)
}

func TestTryUpgrade_NoQuorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{
		{Id: 1, Version: 4},
		{Id: 2, Version: 3},
		{Id: 3, Version: 3},
	}, nil)
	// total power = 3; threshold = ceil(3 * 5/6) = 3; voted power 1 < 3
	repos.Validators.EXPECT().TotalVotingPower(gomock.Any()).Return(types.NumericFromInt64(3), nil)
	repos.SignalVersion.EXPECT().Tally(gomock.Any(), uint64(4)).Return(types.NumericFromInt64(1), nil)
	// SaveUpgrades must NOT be called

	upgrade := &storage.Upgrade{Height: 100, Time: time.Now()}
	err := tryUpgrade(t.Context(), tx, repos.Repos(), upgrade, 3, "", nil)
	require.NoError(t, err)
}

func TestTryUpgrade_WithQuorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{
		{Id: 1, Version: 4},
		{Id: 2, Version: 4},
		{Id: 3, Version: 4},
	}, nil)
	// total power = 6; threshold = 5; voted power 6 → quorum
	repos.Validators.EXPECT().TotalVotingPower(gomock.Any()).Return(types.NumericFromInt64(6), nil)
	repos.SignalVersion.EXPECT().Tally(gomock.Any(), uint64(4)).Return(types.NumericFromInt64(6), nil)
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.EqualValues(t, 4, upgrades[0].Version)
			require.Equal(t, "6", upgrades[0].VotingPower.String())
			require.Equal(t, "6", upgrades[0].VotedPower.String())
			require.Equal(t, types.UpgradeStatusWaitingUpgrade, upgrades[0].Status)
			// MsgTryUpgrade height plus the Mocha delay
			require.EqualValues(t, 100+66_462, upgrades[0].ExpectedHeight)
			return nil
		})
	tx.EXPECT().FixSignalsPower(gomock.Any(), uint64(4)).Return(nil)

	upgrade := &storage.Upgrade{Height: 100, EndHeight: 100, Time: time.Now()}
	err := tryUpgrade(t.Context(), tx, repos.Repos(), upgrade, 3, "mocha-5", nil)
	require.NoError(t, err)
}

func TestTryUpgrade_PicksMinimumQuorumVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{
		{Id: 1, Version: 5},
		{Id: 2, Version: 4},
		{Id: 3, Version: 4},
		{Id: 4, Version: 4},
	}, nil)
	repos.Validators.EXPECT().TotalVotingPower(gomock.Any()).Return(types.NumericFromInt64(12), nil)
	// versions are sorted: v4 has quorum, v5 is never tallied
	repos.SignalVersion.EXPECT().Tally(gomock.Any(), uint64(4)).Return(types.NumericFromInt64(12), nil)
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.EqualValues(t, 4, upgrades[0].Version)
			return nil
		})
	tx.EXPECT().FixSignalsPower(gomock.Any(), uint64(4)).Return(nil)

	upgrade := &storage.Upgrade{Height: 100, Time: time.Now()}
	err := tryUpgrade(t.Context(), tx, repos.Repos(), upgrade, 3, "", nil)
	require.NoError(t, err)
}

func TestTryUpgrade_ResolvesSigner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{{Id: 1, Version: 4}}, nil)
	repos.Validators.EXPECT().TotalVotingPower(gomock.Any()).Return(types.NumericFromInt64(1), nil)
	repos.SignalVersion.EXPECT().Tally(gomock.Any(), uint64(4)).Return(types.NumericFromInt64(1), nil)
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.EqualValues(t, 42, upgrades[0].SignerId)
			return nil
		})
	tx.EXPECT().FixSignalsPower(gomock.Any(), uint64(4)).Return(nil)

	upgrade := &storage.Upgrade{Height: 100, Signer: &storage.Address{Address: "signer"}}
	err := tryUpgrade(t.Context(), tx, repos.Repos(), upgrade, 3, "", map[string]uint64{"signer": 42})
	require.NoError(t, err)
}

func TestTryUpgrade_UnknownSigner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	upgrade := &storage.Upgrade{Height: 100, Signer: &storage.Address{Address: "signer"}}
	err := tryUpgrade(t.Context(), tx, repos.Repos(), upgrade, 3, "", map[string]uint64{})
	require.Error(t, err)
}

func TestSignalThreshold(t *testing.T) {
	tests := []struct {
		total int64
		want  string
	}{
		{0, "0"},
		{3, "3"},    // 2.5 → ceil
		{6, "5"},    // exact
		{7, "6"},    // 5.83 → ceil
		{100, "84"}, // 83.33 → ceil
	}
	for _, tt := range tests {
		got := signalThreshold(10, types.NumericFromInt64(tt.total))
		require.Equal(t, tt.want, got.String(), "total=%d", tt.total)
	}
}

// ---------------------------------------------------------------------------
// recountUpgrades
// ---------------------------------------------------------------------------

func TestRecountUpgrades_SavesOnlyFutureVersions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	// a signal for the current version 3 withdraws a vote and creates no upgrade row
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.EqualValues(t, 4, upgrades[0].Version)
			return nil
		})
	repos.Upgrades.EXPECT().PendingVersions(gomock.Any(), uint64(3)).Return(nil, nil)

	upgrades := makeUpgradesMap(&storage.Upgrade{Version: 3}, &storage.Upgrade{Version: 4})
	err := recountUpgrades(t.Context(), tx, repos.Repos(), upgrades, 3)
	require.NoError(t, err)
}

func TestRecountUpgrades_RecountsVersionsNotInBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	// a validator switched from v4 to v5: v4 must be recounted although nobody signaled it in this block
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).Return(nil)
	repos.Upgrades.EXPECT().PendingVersions(gomock.Any(), uint64(3)).Return([]uint64{4, 5}, nil)
	// total power = 6; threshold = 5
	repos.Validators.EXPECT().TotalVotingPower(gomock.Any()).Return(types.NumericFromInt64(6), nil)
	repos.SignalVersion.EXPECT().Tally(gomock.Any(), uint64(4)).Return(types.NumericFromInt64(4), nil)
	repos.SignalVersion.EXPECT().Tally(gomock.Any(), uint64(5)).Return(types.NumericFromInt64(2), nil)
	tx.EXPECT().UpdateUpgradeTally(gomock.Any(), uint64(4), types.NumericFromInt64(6), types.NumericFromInt64(4), types.UpgradeStatusProcessing).Return(nil)
	tx.EXPECT().UpdateUpgradeTally(gomock.Any(), uint64(5), types.NumericFromInt64(6), types.NumericFromInt64(2), types.UpgradeStatusProcessing).Return(nil)

	upgrades := makeUpgradesMap(&storage.Upgrade{Version: 5})
	err := recountUpgrades(t.Context(), tx, repos.Repos(), upgrades, 3)
	require.NoError(t, err)
}

func TestRecountUpgrades_Quorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).Return(nil)
	repos.Upgrades.EXPECT().PendingVersions(gomock.Any(), uint64(3)).Return([]uint64{4}, nil)
	// total power = 6; threshold = 5; voted = 5 → quorum
	repos.Validators.EXPECT().TotalVotingPower(gomock.Any()).Return(types.NumericFromInt64(6), nil)
	repos.SignalVersion.EXPECT().Tally(gomock.Any(), uint64(4)).Return(types.NumericFromInt64(5), nil)
	tx.EXPECT().UpdateUpgradeTally(gomock.Any(), uint64(4), types.NumericFromInt64(6), types.NumericFromInt64(5), types.UpgradeStatusWaitingUpgrade).Return(nil)

	err := recountUpgrades(t.Context(), tx, repos.Repos(), makeUpgradesMap(&storage.Upgrade{Version: 4}), 3)
	require.NoError(t, err)
}

func TestRecountUpgrades_ZeroTallyAfterReset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	repos := mock.NewTxRepos(ctrl)

	// v4 was applied, v5 was voted before it: the chain wiped those votes, the row goes back to zero
	tx.EXPECT().SaveUpgrades(gomock.Any()).Return(nil)
	repos.Upgrades.EXPECT().PendingVersions(gomock.Any(), uint64(4)).Return([]uint64{5}, nil)
	repos.Validators.EXPECT().TotalVotingPower(gomock.Any()).Return(types.NumericFromInt64(6), nil)
	repos.SignalVersion.EXPECT().Tally(gomock.Any(), uint64(5)).Return(types.NumericZero(), nil)
	tx.EXPECT().UpdateUpgradeTally(gomock.Any(), uint64(5), types.NumericFromInt64(6), types.NumericZero(), types.UpgradeStatusProcessing).Return(nil)

	err := recountUpgrades(t.Context(), tx, repos.Repos(), sdkSync.NewMap[uint64, *storage.Upgrade](), 4)
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// saveSignals
// ---------------------------------------------------------------------------

func TestSaveSignals_Empty(t *testing.T) {
	module := NewModule(nil, nil, indexerCfg.Indexer{Name: testIndexerName})

	err := module.saveSignals(t.Context(), nil, nil)
	require.NoError(t, err)
}

func TestSaveSignals_PowerNotFixed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	tx.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(&[]*storage.SignalVersion{})).Return(nil)

	module := NewModule(nil, nil, indexerCfg.Indexer{Name: testIndexerName})
	module.validatorsByAddress["val1address"] = 1

	signals := []*storage.SignalVersion{
		{Version: 4, Height: 100, Validator: &storage.Validator{Address: "val1address"}},
	}

	err := module.saveSignals(t.Context(), tx, signals)
	require.NoError(t, err)
	require.EqualValues(t, 1, signals[0].ValidatorId)
	require.True(t, signals[0].VotingPower.IsZero())
}

func TestSaveSignals_UnknownValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)

	module := NewModule(nil, nil, indexerCfg.Indexer{Name: testIndexerName})
	// validatorsByAddress intentionally empty

	signals := []*storage.SignalVersion{
		{Version: 4, Height: 100, Validator: &storage.Validator{Address: "unknown_address"}},
	}

	err := module.saveSignals(t.Context(), tx, signals)
	require.Error(t, err)
}
