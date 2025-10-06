// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func (module *Module) saveConstantUpdates(
	ctx context.Context,
	tx storage.Transaction,
	consts *sync.Map[string, *storage.Constant],
) error {

	newConstants := make([]storage.Constant, 0)
	err := consts.Range(func(key string, value *storage.Constant) (error, bool) {
		switch value.Name {
		case "evidence_max_age_num_blocks":
			if value.Value != module.maxAgeNumBlocks {
				module.maxAgeNumBlocks = value.Value
			}
		case "evidence_max_age_duration":
			if value.Value != module.maxAgeDuration {
				module.maxAgeDuration = value.Value
			}
		case "slash_fraction_double_sign":
			if value.Value != module.slashingForDoubleSign.String() {
				val, err := decimal.NewFromString(value.Value)
				if err != nil {
					return errors.Wrap(err, "slash_fraction_double_sign"), true
				}
				module.slashingForDoubleSign = val
			}
		case "slash_fraction_downtime":
			if value.Value != module.slashingForDowntime.String() {
				val, err := decimal.NewFromString(value.Value)
				if err != nil {
					return errors.Wrap(err, "slash_fraction_downtime"), true
				}
				module.slashingForDowntime = val
			}
		}
		newConstants = append(newConstants, *value)
		return nil, false
	})

	if err != nil {
		return errors.Wrap(err, "range")
	}

	if len(newConstants) == 0 {
		return nil
	}

	return tx.SaveConstants(ctx, newConstants...)
}
