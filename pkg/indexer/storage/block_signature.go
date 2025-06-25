// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
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
	if retentionLevel > 0 && height%10 == 0 { // make retention on every ten block
		if err := tx.RetentionBlockSignatures(ctx, retentionLevel); err != nil {
			return err
		}
	}

	if len(signs) == 0 {
		return nil
	}

	for i := range signs {
		if signs[i].Validator == nil {
			return errors.New("nil validator of block signature")
		}

		if id, ok := module.validatorsByConsAddress[signs[i].Validator.ConsAddress]; ok {
			signs[i].ValidatorId = id
		} else {
			return errors.Errorf("unknown validator: %s", signs[i].Validator.ConsAddress)
		}
	}

	return tx.SaveBlockSignatures(ctx, signs...)
}
