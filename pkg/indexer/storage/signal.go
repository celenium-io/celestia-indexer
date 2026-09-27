// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"slices"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/celestia-app/v10/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/v10/x/signal"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

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

// recountUpgrades refreshes the tally of every upgrade that can still change on each block:
// signals move votes between versions and stake changes move power.
func recountUpgrades(
	ctx context.Context,
	tx storage.GovTx,
	repos storage.TxRepos,
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

	versions, err := repos.Upgrades.PendingVersions(ctx, currentVersion)
	if err != nil {
		return errors.Wrap(err, "get pending upgrades")
	}
	if len(versions) == 0 {
		return nil
	}

	votingPower, err := repos.Validators.TotalVotingPower(ctx)
	if err != nil {
		return errors.Wrap(err, "receiving total voting power")
	}
	threshold := signalThreshold(currentVersion, votingPower)

	for _, version := range versions {
		voted, err := repos.SignalVersion.Tally(ctx, version)
		if err != nil {
			return errors.Wrapf(err, "tally version %d", version)
		}

		status := types.UpgradeStatusProcessing
		if voted.GreaterThanOrEqual(threshold) {
			status = types.UpgradeStatusWaitingUpgrade
		}
		if err := tx.UpdateUpgradeTally(ctx, version, votingPower, voted, status); err != nil {
			return errors.Wrapf(err, "update tally of version %d", version)
		}
	}
	return nil
}

func tryUpgrade(
	ctx context.Context,
	tx storage.GovTx,
	repos storage.TxRepos,
	upgrade *storage.Upgrade,
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

	validators, err := repos.Validators.BondedValidators(ctx)
	if err != nil {
		return errors.Wrap(err, "get bonded validators")
	}

	seen := make(map[uint64]struct{})
	var versions []uint64
	for i := range validators {
		v := validators[i].Version
		if v == 0 || v <= currentVersion {
			continue
		}
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			versions = append(versions, v)
		}
	}
	if len(versions) == 0 {
		return nil
	}
	slices.Sort(versions)

	votingPower, err := repos.Validators.TotalVotingPower(ctx)
	if err != nil {
		return errors.Wrap(err, "receiving total voting power")
	}
	threshold := signalThreshold(currentVersion, votingPower)

	for i := range versions {
		voted, err := repos.SignalVersion.Tally(ctx, versions[i])
		if err != nil {
			return errors.Wrapf(err, "tally version %d", versions[i])
		}
		if voted.GreaterThanOrEqual(threshold) {
			upgrade.Version = versions[i]
			upgrade.VotingPower = votingPower
			upgrade.VotedPower = voted
			upgrade.Status = types.UpgradeStatusWaitingUpgrade
			// x/signal schedules the upgrade at the MsgTryUpgrade height plus a per-chain delay
			upgrade.ExpectedHeight = upgrade.EndHeight + pkgTypes.Level(appconsts.GetUpgradeHeightDelay(chainId))
			if err := tx.SaveUpgrades(ctx, upgrade); err != nil {
				return errors.Wrap(err, "save upgrade")
			}
			// keep who contributed which power to the upgrade
			return tx.FixSignalsPower(ctx, versions[i])
		}
	}

	return nil
}

// signalThreshold mirrors x/signal Keeper.GetVotingPowerThreshold.
func signalThreshold(appVersion uint64, totalPower types.Numeric) types.Numeric {
	threshold := signal.Threshold(appVersion).MulInt64(totalPower.IntPart()).Ceil().TruncateInt()
	return types.NumericFromBigInt(threshold.BigInt(), 0)
}
