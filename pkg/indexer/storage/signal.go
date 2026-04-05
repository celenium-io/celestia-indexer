// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"slices"

	"github.com/celenium-io/celestia-indexer/internal/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

var signalsThreshold = types.NumericFromFloat64(5.0 / 6.0)

func (module *Module) saveSignals(
	ctx context.Context,
	tx storage.Transaction,
	signals []*storage.SignalVersion,
	upgrades *sync.Map[uint64, *storage.Upgrade],
	state storage.State,
) error {
	if len(signals) == 0 {
		return nil
	}

	votingPower, _, err := module.totalVotingPower(ctx, tx)
	if err != nil {
		return errors.Wrapf(err, "receiving total voting power")
	}
	votingPower = math.SharesNumeric(votingPower)

	for i := range signals {
		if signals[i].Validator == nil {
			return errors.Errorf("validator is nil in signal version")
		}
		validatorId, ok := module.validatorsByAddress[signals[i].Validator.Address]
		if !ok {
			return errors.Wrap(errCantFindAddress, signals[i].Validator.Address)
		}

		validator, err := tx.Validator(ctx, validatorId)
		if err != nil {
			return errors.Wrapf(err, "get validator by id: %d", validatorId)
		}

		signals[i].VotingPower = validator.Stake
		signals[i].ValidatorId = validatorId
		signals[i].Validator = &validator
	}

	if err := tx.SaveSignals(ctx, signals...); err != nil {
		return errors.Wrap(err, "saving signal version")
	}

	if err := saveUpgrades(ctx, tx, upgrades, state, votingPower); err != nil {
		return errors.Wrap(err, "save upgrades")
	}

	return nil
}

func (module *Module) tryUpgrade(
	ctx context.Context,
	tx storage.Transaction,
	upgrade *storage.Upgrade,
	state storage.State,
) error {
	if upgrade == nil {
		return nil
	}

	votingPower, validators, err := module.totalVotingPower(ctx, tx)
	if err != nil {
		return errors.Wrapf(err, "receiving total voting power")
	}
	votingPower = math.SharesNumeric(votingPower)
	threshold := votingPower.Mul(signalsThreshold)

	seen := make(map[uint64]struct{})
	var versions []uint64
	for i := range validators {
		v := validators[i].Version
		if v == 0 || v <= state.Version {
			continue
		}
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			versions = append(versions, v)
		}
	}
	slices.Sort(versions)

	for i := range versions {
		voted, err := tx.UpdateSignalsAfterUpgrade(ctx, versions[i])
		if err != nil {
			return errors.Wrapf(err, "update signals for version %d", versions[i])
		}
		votedShares := math.SharesNumeric(voted)
		if votedShares.GreaterThan(threshold) {
			upgrade.Version = versions[i]
			upgrade.VotingPower = votingPower
			upgrade.VotedPower = votedShares
			upgrade.Status = types.UpgradeStatusWaitingUpgrade
			return tx.SaveUpgrades(ctx, upgrade)
		}
	}

	return nil
}

func saveUpgrades(
	ctx context.Context,
	tx storage.Transaction,
	upgrades *sync.Map[uint64, *storage.Upgrade],
	state storage.State,
	votingPower types.Numeric,
) error {
	if upgrades.Len() == 0 {
		return nil
	}

	threshold := votingPower.Mul(signalsThreshold)

	var toSave []*storage.Upgrade
	err := upgrades.Range(func(version uint64, upgrade *storage.Upgrade) (error, bool) {
		if state.Version > 0 && state.Version >= version {
			return nil, false
		}

		voted, err := tx.UpdateSignalsAfterUpgrade(ctx, version)
		if err != nil {
			return errors.Wrapf(err, "update signals for version %d", version), true
		}

		upgrade.VotingPower = votingPower
		upgrade.VotedPower = math.SharesNumeric(voted)
		if upgrade.VotedPower.GreaterThan(threshold) {
			upgrade.Status = types.UpgradeStatusWaitingUpgrade
		}
		toSave = append(toSave, upgrade)

		return nil, false
	})
	if err != nil {
		return err
	}
	if len(toSave) == 0 {
		return nil
	}

	return tx.SaveUpgrades(ctx, toSave...)
}

func (module *Module) totalVotingPower(ctx context.Context, tx storage.Transaction) (types.Numeric, []storage.Validator, error) {
	maxVals, err := module.constants.Get(ctx, types.ModuleNameStaking, "max_validators")
	if err != nil {
		return types.Numeric{}, nil, errors.Wrap(err, "get max validators value")
	}
	validators, err := tx.BondedValidators(ctx, maxVals.MustInt())
	if err != nil {
		return types.Numeric{}, nil, errors.Wrap(err, "get validators")
	}

	var power types.Numeric
	for i := range validators {
		power = power.Add(validators[i].Stake)
	}
	return power, validators, nil
}
