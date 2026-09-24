// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

func (module *Module) saveValidators(
	ctx context.Context,
	tx storage.Transaction,
	validators []*storage.Validator,
	jails *sdkSync.Map[string, *storage.Jail],
) (int, error) {
	if jails.Len() > 0 {
		jailedVals := make([]*storage.Validator, 0)
		jailsArr := make([]storage.Jail, 0)

		for address, j := range jails.All() {
			if id, ok := module.validatorsByConsAddress[address]; ok {
				j.ValidatorId = id
				j.Validator.Id = id
				jailedVals = append(jailedVals, j.Validator)
			} else {
				return 0, errors.Errorf("unknown jailed validator: %s", address)
			}

			jailsArr = append(jailsArr, *j)

			if j.Burned.IsZero() {
				continue
			}

			balanceUpdates, err := tx.UpdateSlashedDelegations(ctx, j.ValidatorId, j.Burned)
			if err != nil {
				return 0, err
			}
			if err := tx.SaveBalances(ctx, balanceUpdates...); err != nil {
				return 0, err
			}
		}

		if err := tx.Jail(ctx, jailedVals...); err != nil {
			return 0, err
		}

		if err := tx.SaveJails(ctx, jailsArr...); err != nil {
			return 0, err
		}
	}

	if len(validators) == 0 {
		return 0, nil
	}

	count, err := tx.SaveValidators(ctx, validators...)
	if err != nil {
		return 0, errors.Wrap(err, "saving validators")
	}

	if count == 0 {
		return 0, nil
	}

	module.fillValidatorsCache(validators)

	return count, nil
}

func (module *Module) fillValidatorsCache(validators []*storage.Validator) {
	for i := range validators {
		if validators[i].ConsAddress != "" {
			module.validatorsByConsAddress[validators[i].ConsAddress] = validators[i].Id
		}
		if validators[i].Address != "" {
			module.validatorsByAddress[validators[i].Address] = validators[i].Id
		}
		if validators[i].Delegator != "" {
			module.validatorsByDelegator[validators[i].Delegator] = validators[i].Id
		}
	}
}

// processValidatorBondUpdates resolves the operator address of every bond update and merges its power into dCtx.Validators.
func (module *Module) processValidatorBondUpdates(ctx *decodeContext.Context) error {
	if ctx.ValidatorUpdates.Len() == 0 {
		return nil
	}

	for consAddress, update := range ctx.ValidatorUpdates.All() {
		if update.Validator == nil {
			continue
		}
		validator := *update.Validator

		// a validator created in this block is not in the cache yet
		for v := range ctx.Validators.AllValues() {
			if v.ConsAddress == consAddress {
				validator.Address = v.Address
				break
			}
		}

		if validator.Address == "" {
			id, ok := module.validatorsByConsAddress[consAddress]
			if !ok {
				return errors.Errorf("unknown validator in bond update: %s", consAddress)
			}
			for address, validatorId := range module.validatorsByAddress {
				if validatorId == id {
					validator.Address = address
					break
				}
			}
			if validator.Address == "" {
				return errors.Errorf("can't find operator address of validator %d", id)
			}
			validator.Id = id
		}

		ctx.AddValidator(validator)
	}
	return nil
}

// saveValidatorBondUpdates must run after saveValidators, which caches ids of validators created in this block.
func (module *Module) saveValidatorBondUpdates(
	ctx context.Context,
	tx storage.Transaction,
	updates *sdkSync.Map[string, *storage.ValidatorBondUpdate],
) error {
	if updates.Len() == 0 {
		return nil
	}

	data := make([]*storage.ValidatorBondUpdate, 0, updates.Len())
	for consAddress, update := range updates.All() {
		id, ok := module.validatorsByConsAddress[consAddress]
		if !ok {
			return errors.Errorf("unknown validator in bond update: %s", consAddress)
		}
		update.ValidatorId = id
		data = append(data, update)
	}
	return tx.SaveBondUpdates(ctx, data...)
}
