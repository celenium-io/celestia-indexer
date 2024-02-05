// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

const (
	countOfStoringSignsInLevels = 1_000
)

func (module *Module) saveBlockSignatures(
	ctx context.Context,
	tx storage.Transaction,
	signs []storage.BlockSignature,
	height types.Level,
) error {
	retentionLevel := height - countOfStoringSignsInLevels
	if retentionLevel > 0 {
		if err := tx.RetentionBlockSignatures(ctx, retentionLevel); err != nil {
			return err
		}
	}

	if len(signs) == 0 {
		return nil
	}

	if len(module.validators) == 0 {
		validators, err := tx.Validators(ctx)
		if err != nil {
			return err
		}
		module.validators = make(map[string]uint64)
		for i := range validators {
			module.validators[validators[i].ConsAddress] = validators[i].Id
		}
	}

	for i := range signs {
		if signs[i].Validator == nil {
			return errors.New("nil validator of block signature")
		}

		if id, ok := module.validators[signs[i].Validator.ConsAddress]; ok {
			signs[i].ValidatorId = id
		} else {
			return errors.Errorf("unknown validator: %s", signs[i].Validator.ConsAddress)
		}
	}

	return tx.SaveBlockSignatures(ctx, signs...)
}
