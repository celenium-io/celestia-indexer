// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	fibreTypes "github.com/celestiaorg/celestia-app/v10/x/fibre/types"
	sdkStorage "github.com/dipdup-net/indexer-sdk/pkg/storage"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

func (module *Module) upgrade(ctx context.Context, decodeContext *decodeContext.Context, currentVersion, targetVersion uint64) error {
	if currentVersion >= targetVersion {
		return nil
	}

	for version := currentVersion + 1; version <= targetVersion; version++ {
		switch version {
		case 1, 2, 3, 4, 5, 8, 9:
			// No upgrade logic needed for these versions
		case 6:
			// CIP-037 Reduce the validator unbonding period from 21 days to 14 days and 1 hour to improve capital efficiency while maintaining network security (https://cips.celestia.org/cip-037.html)
			decodeContext.AddConstant(types.ModuleNameStaking, "unbonding_time", "1213200000000000")

			// CIP-041:Reduce inflation to 2.5% and increase minimum validator commission to 10% to improve TIA’s suitability for financial applications (https://cips.celestia.org/cip-041.html)
			decodeContext.AddConstant(types.ModuleNameStaking, "min_commission_rate", "0.100000000000000000")

		case 7:
			if err := module.upgradeV7(ctx, decodeContext, version); err != nil {
				return errors.Wrap(err, "failed to upgrade to version 7")
			}
		case 10:
			if err := module.seedFibreParams(ctx, decodeContext); err != nil {
				return errors.Wrap(err, "failed to seed fibre params")
			}
		default:
			return errors.Errorf("unsupported upgrade version: %d", version)
		}
	}

	return nil
}

// seedFibreParams records the fibre params activated by v10. The upgrade writes
// nothing on chain and emits no event -- the keeper just starts answering with
// DefaultParams -- so the app defaults are the only source. A chain launched at
// v10 already got them from genesis, so existing values are never overwritten.
func (module *Module) seedFibreParams(ctx context.Context, decodeContext *decodeContext.Context) error {
	existing, err := module.constants.ByModule(ctx, types.ModuleNameFibre)
	if err != nil {
		return errors.Wrap(err, "get fibre constants")
	}
	if len(existing) > 0 {
		return nil
	}

	decodeContext.AddFibreParams(fibreTypes.DefaultParams())
	return nil
}

func (module *Module) upgradeV7(ctx context.Context, decodeContext *decodeContext.Context, targetVersion uint64) error {
	if targetVersion != 7 {
		return errors.Errorf("unsupported upgrade version: %d", targetVersion)
	}

	// CIP-044: Increase maximum validator commission to 60% and minimum commission to 20% (https://cips.celestia.org/cip-044.html)
	decodeContext.AddConstant(types.ModuleNameStaking, "min_commission_rate", "0.200000000000000000")
	decodeContext.AddConstant(types.ModuleNameStaking, "max_commission_rate", "0.600000000000000000")

	minCommissionRate := types.MustNumericFromString("0.200000000000000000")
	maxCommissionRate := types.MustNumericFromString("0.600000000000000000")

	paginate := sdkSync.Paginate(
		ctx, 100,
		func(ctx context.Context, limit, offset int) ([]*storage.Validator, error) {
			return module.validators.List(ctx, uint64(limit), uint64(offset), sdkStorage.SortOrderAsc)
		},
	)

	for validator, err := range paginate {
		if err != nil {
			return errors.Wrap(err, "list validators in upgrade v7")
		}
		validator.Rate = getMax(validator.Rate, minCommissionRate)
		validator.MaxRate = getMax(validator.MaxRate, maxCommissionRate)

		// only the rates go into the context: SaveValidators adds up stake, rewards,
		// commissions and messages_count, so the listed values would be counted twice
		newValidator := storage.EmptyValidator()
		newValidator.Address = validator.Address
		newValidator.Rate = validator.Rate
		newValidator.MaxRate = validator.MaxRate

		decodeContext.AddValidator(newValidator)
	}
	return nil
}

func getMax(a, b types.Numeric) types.Numeric {
	if a.GreaterThan(b) {
		return a
	}
	return b
}
