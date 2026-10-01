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
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
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

// snapshotOf builds a snapshot from validator id → power
func snapshotOf(powers map[uint64]int64) powerSnapshot {
	snapshot := powerSnapshot{powers: make(map[uint64]types.Numeric), total: types.NumericZero()}
	for id, power := range powers {
		snapshot.powers[id] = types.NumericFromInt64(power)
		snapshot.total = snapshot.total.Add(types.NumericFromInt64(power))
	}
	return snapshot
}

func latestSignal(id, validatorId, version uint64) storage.SignalVersion {
	return storage.SignalVersion{Id: id, ValidatorId: validatorId, Version: version}
}

// ---------------------------------------------------------------------------
// power snapshot and tally
// ---------------------------------------------------------------------------

func TestTakePowerSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repos := mock.NewTxRepos(ctrl)
	three := types.NumericFromInt64(3)
	zero := types.NumericZero()
	repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{
		{Id: 1, Power: &three},
		{Id: 2},
		{Id: 3, Power: &zero},
	}, nil)

	snapshot, err := takePowerSnapshot(t.Context(), newBlockStartValidators(repos.Validators))
	require.NoError(t, err)
	require.Equal(t, "3", snapshot.total.String())
	require.Len(t, snapshot.powers, 1)
	require.Equal(t, "3", snapshot.powers[1].String())
}

func TestPrepareSignalRound(t *testing.T) {
	for _, tt := range []struct {
		name         string
		stateVersion uint64
		pending      []uint64
		upgrades     []*storage.Upgrade
		tryUpgrade   *storage.Upgrade
		wantSnapshot bool
	}{
		{name: "quiet block", stateVersion: 3},
		// a signal for the current version only withdraws a vote
		{name: "signal for current version", stateVersion: 3, upgrades: []*storage.Upgrade{{Version: 3}}},
		{name: "open round", stateVersion: 3, pending: []uint64{4}, wantSnapshot: true},
		{name: "first signal for a version", stateVersion: 3, upgrades: []*storage.Upgrade{{Version: 4}}, wantSnapshot: true},
		{name: "try upgrade", stateVersion: 3, tryUpgrade: &storage.Upgrade{}, wantSnapshot: true},
		{name: "applied upgrade", stateVersion: 2, wantSnapshot: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repos := mock.NewTxRepos(ctrl)
			repos.Upgrades.EXPECT().PendingVersions(gomock.Any(), uint64(3)).Return(tt.pending, nil)
			if tt.wantSnapshot {
				one := types.NumericFromInt64(1)
				repos.Validators.EXPECT().BondedValidators(gomock.Any()).Return([]storage.Validator{{Id: 1, Power: &one}}, nil)
			}

			dCtx := decodeContext.NewContext()
			dCtx.Block = &storage.Block{VersionApp: 3}
			dCtx.TryUpgrade = tt.tryUpgrade
			for _, upgrade := range tt.upgrades {
				dCtx.AddUpgrade(*upgrade)
			}

			round, err := prepareSignalRound(t.Context(), repos.Repos(), dCtx, tt.stateVersion, newBlockStartValidators(repos.Validators))
			require.NoError(t, err)
			require.Equal(t, tt.pending, round.pending)
			require.Equal(t, tt.wantSnapshot, round.snapshot != nil)
		})
	}
}

func TestSignalRoundTallyWithoutSnapshot(t *testing.T) {
	_, err := signalRound{}.tally(t.Context(), nil)
	require.Error(t, err)
}

func TestSignalTally(t *testing.T) {
	// validator 2 is unbonded, validator 4 has not signaled
	tally := newSignalTally([]storage.SignalVersion{
		latestSignal(10, 1, 4),
		latestSignal(11, 2, 4),
		latestSignal(12, 3, 5),
		latestSignal(13, 5, 3),
	}, snapshotOf(map[uint64]int64{1: 3, 3: 2, 4: 7, 5: 1}))

	require.Equal(t, "3", tally.voted(4).String())
	require.Equal(t, "2", tally.voted(5).String())
	require.Equal(t, "0", tally.voted(6).String())
	require.Equal(t, []uint64{4, 5}, tally.versions(3))
	require.Equal(t, []uint64{5}, tally.versions(4))
	require.Equal(t, map[uint64]types.Numeric{10: types.NumericFromInt64(3)}, tally.counted(4))
}

func TestLoadTally(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repos := mock.NewTxRepos(ctrl)
	repos.SignalVersion.EXPECT().Latest(gomock.Any()).Return([]storage.SignalVersion{latestSignal(10, 1, 4)}, nil)

	tally, err := loadTally(t.Context(), repos.SignalVersion, snapshotOf(map[uint64]int64{1: 3}))
	require.NoError(t, err)
	require.Equal(t, "3", tally.voted(4).String())
}

// ---------------------------------------------------------------------------
// tryUpgrade
// ---------------------------------------------------------------------------

func TestTryUpgrade_NilUpgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	err := tryUpgrade(t.Context(), tx, nil, signalTally{}, 3, "", nil)
	require.NoError(t, err)
}

func TestTryUpgrade_NoSignals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	// only a withdrawn vote for the current version: nothing to upgrade to
	tally := newSignalTally([]storage.SignalVersion{latestSignal(10, 1, 3)}, snapshotOf(map[uint64]int64{1: 1}))

	upgrade := &storage.Upgrade{Height: 100, Time: time.Now()}
	err := tryUpgrade(t.Context(), tx, upgrade, tally, 3, "", nil)
	require.NoError(t, err)
}

func TestTryUpgrade_NoQuorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	// total power = 3; threshold = ceil(3 * 5/6) = 3; voted power 1 < 3
	tally := newSignalTally([]storage.SignalVersion{latestSignal(10, 1, 4)}, snapshotOf(map[uint64]int64{1: 1, 2: 1, 3: 1}))
	// SaveUpgrades must NOT be called

	upgrade := &storage.Upgrade{Height: 100, Time: time.Now()}
	err := tryUpgrade(t.Context(), tx, upgrade, tally, 3, "", nil)
	require.NoError(t, err)
}

func TestTryUpgrade_WithQuorum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	// total power = 6; threshold = 5; voted power 5 → quorum; validator 3 did not signal
	tally := newSignalTally([]storage.SignalVersion{
		latestSignal(10, 1, 4),
		latestSignal(11, 2, 4),
	}, snapshotOf(map[uint64]int64{1: 2, 2: 3, 3: 1}))

	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.EqualValues(t, 4, upgrades[0].Version)
			require.Equal(t, "6", upgrades[0].VotingPower.String())
			require.Equal(t, "5", upgrades[0].VotedPower.String())
			require.Equal(t, types.UpgradeStatusWaitingUpgrade, upgrades[0].Status)
			// MsgTryUpgrade height plus the Mocha delay
			require.EqualValues(t, 100+66_462, upgrades[0].ExpectedHeight)
			return nil
		})
	tx.EXPECT().FixSignalsPower(gomock.Any(), map[uint64]types.Numeric{
		10: types.NumericFromInt64(2),
		11: types.NumericFromInt64(3),
	}).Return(nil)

	upgrade := &storage.Upgrade{Height: 100, EndHeight: 100, Time: time.Now()}
	err := tryUpgrade(t.Context(), tx, upgrade, tally, 3, "mocha-5", nil)
	require.NoError(t, err)
}

func TestTryUpgrade_ResolvesSigner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	tally := newSignalTally([]storage.SignalVersion{latestSignal(10, 1, 4)}, snapshotOf(map[uint64]int64{1: 1}))
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.EqualValues(t, 42, upgrades[0].SignerId)
			return nil
		})
	tx.EXPECT().FixSignalsPower(gomock.Any(), gomock.Any()).Return(nil)

	upgrade := &storage.Upgrade{Height: 100, Signer: &storage.Address{Address: "signer"}}
	err := tryUpgrade(t.Context(), tx, upgrade, tally, 3, "", map[string]uint64{"signer": 42})
	require.NoError(t, err)
}

func TestTryUpgrade_UnknownSigner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)

	upgrade := &storage.Upgrade{Height: 100, Signer: &storage.Address{Address: "signer"}}
	err := tryUpgrade(t.Context(), tx, upgrade, signalTally{}, 3, "", map[string]uint64{})
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
// saveUpgradeSignals and recountUpgrades
// ---------------------------------------------------------------------------

func TestSaveUpgradeSignals_OnlyFutureVersions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	// a signal for the current version 3 withdraws a vote and creates no upgrade row
	tx.EXPECT().SaveUpgrades(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, upgrades ...*storage.Upgrade) error {
			require.Len(t, upgrades, 1)
			require.EqualValues(t, 4, upgrades[0].Version)
			return nil
		})

	upgrades := makeUpgradesMap(&storage.Upgrade{Version: 3}, &storage.Upgrade{Version: 4})
	err := saveUpgradeSignals(t.Context(), tx, upgrades, 3)
	require.NoError(t, err)
}

func TestRecountUpgrades(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	// total power = 6; threshold = 5
	tally := newSignalTally([]storage.SignalVersion{
		latestSignal(10, 1, 4),
		latestSignal(11, 2, 5),
	}, snapshotOf(map[uint64]int64{1: 5, 2: 1}))

	total := types.NumericFromInt64(6)
	tx.EXPECT().UpdateUpgradeTally(gomock.Any(), uint64(4), total, types.NumericFromInt64(5), types.UpgradeStatusWaitingUpgrade).Return(nil)
	tx.EXPECT().UpdateUpgradeTally(gomock.Any(), uint64(5), total, types.NumericFromInt64(1), types.UpgradeStatusProcessing).Return(nil)
	// votes wiped by an applied upgrade: the row goes back to zero
	tx.EXPECT().UpdateUpgradeTally(gomock.Any(), uint64(6), total, types.NumericZero(), types.UpgradeStatusProcessing).Return(nil)

	err := recountUpgrades(t.Context(), tx, []uint64{4, 5, 6}, tally, 3)
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
	require.ErrorIs(t, err, errCantFindAddress)
}
