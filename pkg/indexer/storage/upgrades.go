// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/pkg/errors"
)

func Upgrade(ctx *decodeContext.Context, currentVersion, targetVersion uint64) error {
	if currentVersion >= targetVersion {
		return nil
	}

	switch targetVersion {
	case 1, 2, 3, 4, 5:
		// No upgrade logic needed for these versions
	case 6:
		// CIP-037 Reduce the validator unbonding period from 21 days to 14 days and 1 hour to improve capital efficiency while maintaining network security (https://cips.celestia.org/cip-037.html)
		ctx.AddConstant(types.ModuleNameStaking, "unbonding_time", "1213200000000000")

		// CIP-041:Reduce inflation to 2.5% and increase minimum validator commission to 10% to improve TIA’s suitability for financial applications (https://cips.celestia.org/cip-041.html)
		ctx.AddConstant(types.ModuleNameStaking, "min_commission_rate", "0.100000000000000000")

	default:
		return errors.Errorf("unsupported upgrade version: %d", targetVersion)
	}

	return nil
}
