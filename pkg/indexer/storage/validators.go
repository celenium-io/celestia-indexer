// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

func (module *Module) saveValidators(
	ctx context.Context,
	tx storage.Transaction,
	validators []*storage.Validator,
	jailed *sync.Map[string, struct{}],
	jails []storage.Jail,
) (int, error) {
	if jailed.Len() > 0 {
		jailedIds := make([]uint64, 0)
		err := jailed.Range(func(address string, _ struct{}) (error, bool) {
			if id, ok := module.validatorsByConsAddress[address]; ok {
				jailedIds = append(jailedIds, id)
				return nil, false
			}

			return errors.Errorf("unknown jailed validator: %s", address), false
		})
		if err != nil {
			return 0, err
		}

		if err := tx.Jail(ctx, jailedIds...); err != nil {
			return 0, err
		}
	}

	if len(jails) > 0 {
		for i := range jails {
			if id, ok := module.validatorsByConsAddress[jails[i].Validator.ConsAddress]; ok {
				jails[i].ValidatorId = id
			}
		}

		if err := tx.SaveJails(ctx, jails...); err != nil {
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
