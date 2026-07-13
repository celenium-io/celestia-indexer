// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func (module *Module) saveConstantUpdates(
	ctx context.Context,
	tx storage.Transaction,
	consts *sdkSync.Map[string, *storage.Constant],
) error {

	newConstants := make([]storage.Constant, 0)

	for _, value := range consts.All() {
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
					return errors.Wrap(err, "slash_fraction_double_sign")
				}
				module.slashingForDoubleSign = val
			}
		case "slash_fraction_downtime":
			if value.Value != module.slashingForDowntime.String() {
				val, err := decimal.NewFromString(value.Value)
				if err != nil {
					return errors.Wrap(err, "slash_fraction_downtime")
				}
				module.slashingForDowntime = val
			}
		}
		newConstants = append(newConstants, *value)
	}

	if len(newConstants) == 0 {
		return nil
	}

	return tx.SaveConstants(ctx, newConstants...)
}
