// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var signalsThreshold = decimal.NewFromFloat(5.0 / 6.0)

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

	votingPower, validators, err := module.totalVotingPower(ctx, tx)
	if err != nil {
		return errors.Wrapf(err, "receiving total voting power")
	}
	votingPower = math.Shares(votingPower)

	validatorsMap := make(map[uint64]storage.Validator, len(validators))
	for i := range validators {
		validatorsMap[validators[i].Id] = validators[i]
	}

	for i := range signals {
		if signals[i].Validator == nil {
			return errors.Errorf("validator is nil in signal version")
		}
		validatorId, ok := module.validatorsByAddress[signals[i].Validator.Address]
		if !ok {
			return errors.Wrap(errCantFindAddress, signals[i].Validator.Address)
		}

		validator, ok := validatorsMap[validatorId]
		if !ok {
			validator, err = tx.Validator(ctx, validatorId)
			if err != nil {
				return errors.Wrapf(err, "get validator by id: %d", validatorId)
			}
		}

		signals[i].VotingPower = validator.Stake
		signals[i].ValidatorId = validatorId
		signals[i].Validator = &validator
	}

	if err := tx.SaveSignals(ctx, signals...); err != nil {
		return errors.Wrap(err, "saving signal version")
	}

	if err := postProcessingSignal(ctx, tx, signals, upgrades, votingPower, validators); err != nil {
		return errors.Wrap(err, "postProcessingSignal")
	}

	if err := saveUpgrades(ctx, tx, upgrades, state, votingPower, validators); err != nil {
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
	votingPower = math.Shares(votingPower)

	pass, voted, err := recalculateSignalsForUpgrade(state.Version+1, state, votingPower, validators)
	if err != nil {
		return errors.Wrap(err, "recalculateSignalsForUpgrade")
	}

	upgrade.Version = state.Version + 1
	upgrade.VotingPower = votingPower
	upgrade.VotedPower = voted
	if pass {
		upgrade.Status = types.UpgradeStatusWaitingUpgrade
	}

	return tx.SaveUpgrades(ctx, upgrade)
}

func saveUpgrades(
	ctx context.Context,
	tx storage.Transaction,
	upgrades *sync.Map[uint64, *storage.Upgrade],
	state storage.State,
	votingPower decimal.Decimal,
	validators []storage.Validator,
) error {
	if upgrades.Len() == 0 {
		return nil
	}

	err := upgrades.Range(func(version uint64, upgrade *storage.Upgrade) (error, bool) {
		if state.Version > 0 && state.Version >= version {
			return nil, false
		}

		pass, voted, err := recalculateSignalsForUpgrade(version, state, votingPower, validators)
		if err != nil {
			return errors.Wrap(err, "recalculateSignalsForUpgrade"), true
		}

		upgrade.VotingPower = votingPower
		upgrade.VotedPower = voted
		if pass {
			upgrade.Status = types.UpgradeStatusWaitingUpgrade
		}

		return nil, false
	})

	if err != nil {
		return err
	}

	return tx.SaveUpgrades(ctx, upgrades.Values()...)
}

func recalculateSignalsForUpgrade(
	version uint64,
	state storage.State,
	votingPower decimal.Decimal,
	validators []storage.Validator,
) (bool, decimal.Decimal, error) {
	if version == 0 {
		return false, decimal.Zero, errors.New("recalculate signals for 0 upgrade")
	}

	if state.Version > 0 && state.Version >= version {
		return false, decimal.Zero, nil
	}

	var (
		pass      bool
		threshold = votingPower.Mul(signalsThreshold)
		voted     = decimal.Zero
	)

	for i := range validators {
		if validators[i].Version != version {
			continue
		}
		voted = voted.Add(math.Shares(validators[i].Stake))
		if voted.GreaterThan(threshold) {
			pass = true
		}
	}

	return pass, voted, nil
}

func postProcessingSignal(
	ctx context.Context,
	tx storage.Transaction,
	signals []*storage.SignalVersion,
	upgrades *sync.Map[uint64, *storage.Upgrade],
	votingPower decimal.Decimal,
	validators []storage.Validator,
) error {
	if len(signals) == 0 {
		return nil
	}

	versions := map[uint64]struct{}{}
	for i := range signals {
		versions[signals[i].Version] = struct{}{}
	}

	for version := range versions {
		var voted decimal.Decimal
		for i := range validators {
			if validators[i].Version != version {
				continue
			}
			voted = voted.Add(math.Shares(validators[i].Stake))
		}

		if val, ok := upgrades.Get(version); ok {
			val.VotedPower = voted
			val.VotingPower = votingPower
		} else {
			return errors.Errorf("found signal without upgrade version %d", version)
		}

		if err := tx.UpdateSignalsAfterUpgrade(ctx, version); err != nil {
			return errors.Wrapf(err, "update signals for version %d", version)
		}
	}

	return nil
}

func (module *Module) totalVotingPower(ctx context.Context, tx storage.Transaction) (decimal.Decimal, []storage.Validator, error) {
	maxVals, err := module.constants.Get(ctx, types.ModuleNameStaking, "max_validators")
	if err != nil {
		return decimal.Zero, nil, errors.Wrap(err, "get max validators value")
	}
	validators, err := tx.BondedValidators(ctx, maxVals.MustInt())
	if err != nil {
		return decimal.Zero, nil, errors.Wrap(err, "get validators")
	}

	power := decimal.Zero
	for i := range validators {
		power = power.Add(validators[i].Stake)
	}
	return power, validators, nil
}
