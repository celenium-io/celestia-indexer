// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"database/sql"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	fibreTypes "github.com/celestiaorg/celestia-app/v10/x/fibre/types"
	sdkStorage "github.com/dipdup-net/indexer-sdk/pkg/storage"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

// LegacyDec form of appconsts.DefaultNetworkMinGasPrice
const defaultNetworkMinGasPrice = "0.000001000000000000"

func upgrade(
	ctx context.Context, repos storage.TxRepos, decodeContext *decodeContext.Context, currentVersion, targetVersion uint64,
) error {
	if currentVersion >= targetVersion {
		return nil
	}

	for version := currentVersion + 1; version <= targetVersion; version++ {
		switch version {
		case 1, 3, 5, 8, 9:
			// No upgrade logic needed for these versions
		case 2:
			// CIP-006: Minimum gas price enforcement (https://cips.celestia.org/cip-006.html)
			if err := createConstantIfNotExists(
				ctx, decodeContext, repos.Constants, types.ModuleNameMinfee, "network_min_gas_price", defaultNetworkMinGasPrice,
			); err != nil {
				return errors.Wrap(err, "failed to upgrade to version 2")
			}
		case 6:
			// CIP-037 Reduce the validator unbonding period from 21 days to 14 days and 1 hour to improve capital
			// efficiency while maintaining network security (https://cips.celestia.org/cip-037.html)
			decodeContext.AddConstant(types.ModuleNameStaking, "unbonding_time", "1213200000000000")

			// CIP-041:Reduce inflation to 2.5% and increase minimum validator commission to 10% to improve TIA’s
			// suitability for financial applications (https://cips.celestia.org/cip-041.html)
			decodeContext.AddConstant(types.ModuleNameStaking, "min_commission_rate", "0.100000000000000000")

		case 7:
			if err := upgradeV7(ctx, repos, decodeContext, version); err != nil {
				return errors.Wrap(err, "failed to upgrade to version 7")
			}
		case 10:
			if err := seedFibreParams(ctx, repos, decodeContext); err != nil {
				return errors.Wrap(err, "failed to seed fibre params")
			}
		case 4:
			if err := upgradeV4(ctx, decodeContext, repos.Constants); err != nil {
				return err
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
func seedFibreParams(
	ctx context.Context, repos storage.TxRepos, decodeContext *decodeContext.Context,
) error {
	existing, err := repos.Constants.ByModule(ctx, types.ModuleNameFibre)
	if err != nil {
		return errors.Wrap(err, "get fibre constants")
	}
	if len(existing) > 0 {
		return nil
	}

	decodeContext.AddFibreParams(fibreTypes.DefaultParams())
	return nil
}

func upgradeV7(
	ctx context.Context, repos storage.TxRepos, decodeContext *decodeContext.Context, targetVersion uint64,
) error {
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
			return repos.Validators.List(ctx, uint64(limit), uint64(offset), sdkStorage.SortOrderAsc)
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

func upgradeV4(
	ctx context.Context,
	decodeContext *decodeContext.Context,
	constants storage.IConstant,
) error {
	if err := createConstantIfNotExists(
		ctx, decodeContext, constants, types.ModuleNameGov, "expedited_voting_period", "86400000000000",
	); err != nil {
		return err
	}
	if err := createConstantIfNotExists(
		ctx, decodeContext, constants, types.ModuleNameGov, "expedited_threshold", "0.667000000000000000",
	); err != nil {
		return err
	}
	if err := createConstantIfNotExists(
		ctx, decodeContext, constants, types.ModuleNameGov, "expedited_min_deposit", "50000000000utia",
	); err != nil {
		return err
	}
	return nil
}

func createConstantIfNotExists(
	ctx context.Context,
	decodeContext *decodeContext.Context,
	constants storage.IConstant,
	module types.ModuleName,
	name, value string,
) error {
	if _, err := constants.Get(ctx, module, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			decodeContext.AddConstant(module, name, value)
		} else {
			return err
		}
	}
	return nil
}
