// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (module *Module) saveValidators(
	ctx context.Context,
	tx storage.Transaction,
	validators []*storage.Validator,
	jails *sync.Map[string, *storage.Jail],
) (int, error) {
	if jails.Len() > 0 {
		jailedVals := make([]*storage.Validator, 0)
		jailsArr := make([]storage.Jail, 0)

		err := jails.Range(func(address string, j *storage.Jail) (error, bool) {
			if id, ok := module.validatorsByConsAddress[address]; ok {
				j.ValidatorId = id
				j.Validator.Id = id
				jailedVals = append(jailedVals, j.Validator)
			} else {
				return errors.Errorf("unknown jailed validator: %s", address), false
			}

			jailsArr = append(jailsArr, *j)

			fraction := decimal.Zero
			switch j.Reason {
			case "double_sign":
				fraction = module.slashingForDoubleSign.Copy()
			case "missing_signature":
				fraction = module.slashingForDowntime.Copy()
			}
			if !fraction.IsPositive() {
				return nil, false
			}

			balanceUpdates, err := tx.UpdateSlashedDelegations(ctx, j.ValidatorId, fraction)
			if err != nil {
				return err, false
			}
			if err := tx.SaveBalances(ctx, balanceUpdates...); err != nil {
				return err, false
			}

			return nil, false
		})
		if err != nil {
			return 0, err
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

	for i := range validators {
		if validators[i].ConsAddress == "" {
			continue
		}
		module.validatorsByConsAddress[validators[i].ConsAddress] = validators[i].Id
		module.validatorsByAddress[validators[i].Address] = validators[i].Id
	}

	return count, nil
}
