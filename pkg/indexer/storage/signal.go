// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"slices"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/celestia-app/v10/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/v10/x/signal"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

// powerSnapshot is the consensus power x/signal reads during the block, left by the previous block's EndBlock.
// It is taken before saveValidators writes this block's EndBlock changes.
type powerSnapshot struct {
	powers map[uint64]types.Numeric
	total  types.Numeric
}

func takePowerSnapshot(ctx context.Context, validators storage.IValidator) (powerSnapshot, error) {
	bonded, err := validators.BondedValidators(ctx)
	if err != nil {
		return powerSnapshot{}, errors.Wrap(err, "get bonded validators")
	}
	snapshot := powerSnapshot{
		powers: make(map[uint64]types.Numeric, len(bonded)),
		total:  types.NumericZero(),
	}
	for i := range bonded {
		if bonded[i].Power == nil || !bonded[i].Power.IsPositive() {
			continue
		}
		snapshot.powers[bonded[i].Id] = *bonded[i].Power
		snapshot.total = snapshot.total.Add(*bonded[i].Power)
	}
	return snapshot, nil
}

// signalTally mirrors x/signal TallyVotingPower: latest signals of bonded validators weighted by the snapshot.
type signalTally struct {
	latest    []storage.SignalVersion
	snapshot  powerSnapshot
	byVersion map[uint64]types.Numeric
}

func loadTally(ctx context.Context, signals storage.ISignalVersion, snapshot powerSnapshot) (signalTally, error) {
	latest, err := signals.Latest(ctx)
	if err != nil {
		return signalTally{}, errors.Wrap(err, "get latest signals")
	}
	return newSignalTally(latest, snapshot), nil
}

func newSignalTally(latest []storage.SignalVersion, snapshot powerSnapshot) signalTally {
	tally := signalTally{
		latest:    latest,
		snapshot:  snapshot,
		byVersion: make(map[uint64]types.Numeric),
	}
	for i := range latest {
		power, ok := snapshot.powers[latest[i].ValidatorId]
		if !ok {
			continue
		}
		tally.byVersion[latest[i].Version] = tally.voted(latest[i].Version).Add(power)
	}
	return tally
}

func (t signalTally) voted(version uint64) types.Numeric {
	if power, ok := t.byVersion[version]; ok {
		return power
	}
	return types.NumericZero()
}

// versions returns the signalled versions above the current one in ascending order.
func (t signalTally) versions(currentVersion uint64) []uint64 {
	versions := make([]uint64, 0, len(t.byVersion))
	for version := range t.byVersion {
		if version > currentVersion {
			versions = append(versions, version)
		}
	}
	slices.Sort(versions)
	return versions
}

// counted returns the power each signal for version contributed, by signal id.
func (t signalTally) counted(version uint64) map[uint64]types.Numeric {
	powers := make(map[uint64]types.Numeric)
	for i := range t.latest {
		if t.latest[i].Version != version {
			continue
		}
		if power, ok := t.snapshot.powers[t.latest[i].ValidatorId]; ok {
			powers[t.latest[i].Id] = power
		}
	}
	return powers
}

// processSignalModule runs after saveValidators, so validators created in the block are known,
// and weighs signals by the snapshot taken before it.
func (module *Module) processSignalModule(
	ctx context.Context,
	tx storage.GovTx,
	repos storage.TxRepos,
	decodeContext *decodeContext.Context,
	stateVersion uint64,
	addrToId map[string]uint64,
	snapshot powerSnapshot,
) error {
	currentVersion := decodeContext.Block.VersionApp

	// closes the round if MsgTryUpgrade was not seen (e.g. sent via authz); a no-op otherwise
	if stateVersion < currentVersion {
		tally, err := loadTally(ctx, repos.SignalVersion, snapshot)
		if err != nil {
			return err
		}
		if err := tx.FixSignalsPower(ctx, tally.counted(currentVersion)); err != nil {
			return errors.Wrap(err, "fix signals power")
		}
	}

	// before signals: the tally only counts signals since the last applied upgrade
	if err := module.setUpgradeApplied(ctx, tx, stateVersion, decodeContext.Block); err != nil {
		return errors.Wrap(err, "set upgrade applied")
	}

	if err := module.saveSignals(ctx, tx, decodeContext.Signals); err != nil {
		return errors.Wrap(err, "save signals")
	}

	if err := saveUpgradeSignals(ctx, tx, decodeContext.Upgrades, currentVersion); err != nil {
		return err
	}

	pending, err := repos.Upgrades.PendingVersions(ctx, currentVersion)
	if err != nil {
		return errors.Wrap(err, "get pending upgrades")
	}
	if len(pending) == 0 && decodeContext.TryUpgrade == nil {
		return nil
	}

	tally, err := loadTally(ctx, repos.SignalVersion, snapshot)
	if err != nil {
		return err
	}

	if err := recountUpgrades(ctx, tx, pending, tally, currentVersion); err != nil {
		return errors.Wrap(err, "recount upgrades")
	}

	if err := tryUpgrade(
		ctx, tx, decodeContext.TryUpgrade, tally, currentVersion, decodeContext.Block.ChainId, addrToId,
	); err != nil {
		return errors.Wrap(err, "tryUpgrade")
	}
	return nil
}

func (module *Module) saveSignals(
	ctx context.Context,
	tx storage.GovTx,
	signals []*storage.SignalVersion,
) error {
	if len(signals) == 0 {
		return nil
	}

	for i := range signals {
		if signals[i].Validator == nil {
			return errors.Errorf("validator is nil in signal version")
		}
		validatorId, ok := module.validatorsByAddress[signals[i].Validator.Address]
		if !ok {
			return errors.Wrap(errCantFindAddress, signals[i].Validator.Address)
		}

		// stored as null until an upgrade closes the round; the API shows the current power meanwhile
		signals[i].VotingPower = types.NumericZero()
		signals[i].ValidatorId = validatorId
	}

	if err := storage.Insert(ctx, tx, signals...); err != nil {
		return errors.Wrap(err, "saving signal version")
	}
	return nil
}

// saveUpgradeSignals creates upgrade rows and counts signals; a signal for the current version only withdraws a vote.
func saveUpgradeSignals(
	ctx context.Context,
	tx storage.GovTx,
	upgrades *sdkSync.Map[uint64, *storage.Upgrade],
	currentVersion uint64,
) error {
	var toSave []*storage.Upgrade
	for version, upgrade := range upgrades.All() {
		if version > currentVersion {
			toSave = append(toSave, upgrade)
		}
	}
	if err := tx.SaveUpgrades(ctx, toSave...); err != nil {
		return errors.Wrap(err, "save upgrades")
	}
	return nil
}

// recountUpgrades refreshes the tally of every upgrade that can still change on each block:
// signals move votes between versions and stake changes move power.
func recountUpgrades(
	ctx context.Context,
	tx storage.GovTx,
	versions []uint64,
	tally signalTally,
	currentVersion uint64,
) error {
	threshold := signalThreshold(currentVersion, tally.snapshot.total)
	for _, version := range versions {
		voted := tally.voted(version)
		status := types.UpgradeStatusProcessing
		if voted.GreaterThanOrEqual(threshold) {
			status = types.UpgradeStatusWaitingUpgrade
		}
		if err := tx.UpdateUpgradeTally(ctx, version, tally.snapshot.total, voted, status); err != nil {
			return errors.Wrapf(err, "update tally of version %d", version)
		}
	}
	return nil
}

func tryUpgrade(
	ctx context.Context,
	tx storage.GovTx,
	upgrade *storage.Upgrade,
	tally signalTally,
	currentVersion uint64,
	chainId string,
	addrToId map[string]uint64,
) error {
	if upgrade == nil {
		return nil
	}
	if upgrade.Signer != nil {
		signerId, ok := addrToId[upgrade.Signer.Address]
		if !ok {
			return errors.Wrapf(errCantFindAddress, "try upgrade signer: %s", upgrade.Signer.Address)
		}
		upgrade.SignerId = signerId
	}

	threshold := signalThreshold(currentVersion, tally.snapshot.total)
	for _, version := range tally.versions(currentVersion) {
		voted := tally.voted(version)
		if !voted.GreaterThanOrEqual(threshold) {
			continue
		}
		upgrade.Version = version
		upgrade.VotingPower = tally.snapshot.total
		upgrade.VotedPower = voted
		upgrade.Status = types.UpgradeStatusWaitingUpgrade
		// x/signal schedules the upgrade at the MsgTryUpgrade height plus a per-chain delay
		upgrade.ExpectedHeight = upgrade.EndHeight + pkgTypes.Level(appconsts.GetUpgradeHeightDelay(chainId))
		if err := tx.SaveUpgrades(ctx, upgrade); err != nil {
			return errors.Wrap(err, "save upgrade")
		}
		// keep who contributed which power to the upgrade
		return tx.FixSignalsPower(ctx, tally.counted(version))
	}
	return nil
}

// signalThreshold mirrors x/signal Keeper.GetVotingPowerThreshold.
func signalThreshold(appVersion uint64, totalPower types.Numeric) types.Numeric {
	threshold := signal.Threshold(appVersion).MulInt64(totalPower.IntPart()).Ceil().TruncateInt()
	return types.NumericFromBigInt(threshold.BigInt(), 0)
}
