// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (module *Module) upgrade(ctx context.Context, decodeContext *decodeContext.Context, currentVersion, targetVersion uint64) error {
	if currentVersion >= targetVersion {
		return nil
	}

	switch targetVersion {
	case 1, 2, 3, 4, 5:
		// No upgrade logic needed for these versions
	case 6:
		// CIP-037 Reduce the validator unbonding period from 21 days to 14 days and 1 hour to improve capital efficiency while maintaining network security (https://cips.celestia.org/cip-037.html)
		decodeContext.AddConstant(types.ModuleNameStaking, "unbonding_time", "1213200000000000")

		// CIP-041:Reduce inflation to 2.5% and increase minimum validator commission to 10% to improve TIA’s suitability for financial applications (https://cips.celestia.org/cip-041.html)
		decodeContext.AddConstant(types.ModuleNameStaking, "min_commission_rate", "0.100000000000000000")

	case 7:
		if err := module.upgradeV7(ctx, decodeContext, targetVersion); err != nil {
			return errors.Wrap(err, "failed to upgrade to version 7")
		}
	default:
		return errors.Errorf("unsupported upgrade version: %d", targetVersion)
	}

	return nil
}

func (module *Module) upgradeV7(ctx context.Context, decodeContext *decodeContext.Context, targetVersion uint64) error {
	if targetVersion != 7 {
		return errors.Errorf("unsupported upgrade version: %d", targetVersion)
	}

	// CIP-044: Increase maximum validator commission to 60% and minimum commission to 20% (https://cips.celestia.org/cip-044.html)
	decodeContext.AddConstant(types.ModuleNameStaking, "min_commission_rate", "0.200000000000000000")
	decodeContext.AddConstant(types.ModuleNameStaking, "max_commission_rate", "0.600000000000000000")

	minCommissionRate := decimal.RequireFromString("0.200000000000000000")
	maxCommissionRate := decimal.RequireFromString("0.600000000000000000")
	limit := uint64(100)
	offset := uint64(0)
	end := false
	for !end {
		validators, err := module.validators.List(ctx, limit, offset, storage.SortOrderAsc)
		if err != nil {
			return errors.Wrap(err, "failed to list validators")
		}
		if len(validators) == 0 {
			break
		}

		for i := range validators {
			if validators[i].Rate.LessThan(minCommissionRate) {
				validators[i].Rate = minCommissionRate
			}
			if validators[i].MaxRate.GreaterThan(maxCommissionRate) {
				validators[i].MaxRate = maxCommissionRate
			}
			decodeContext.AddValidator(*validators[i])
		}
		offset += limit
		end = len(validators) < int(limit)
	}
	return nil
}
