// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

func saveVestings(
	ctx context.Context,
	tx storage.Transaction,
	accounts []*storage.VestingAccount,
	addrToId map[string]uint64,
) error {
	if len(accounts) == 0 {
		return nil
	}

	for i := range accounts {
		if accounts[i].Address == nil {
			return errors.Errorf("nil address in vesting account")
		}
		addrId, ok := addrToId[accounts[i].Address.Address]
		if !ok {
			return errors.Wrap(errCantFindAddress, accounts[i].Address.Address)
		}
		accounts[i].AddressId = addrId
	}

	if err := tx.SaveVestingAccounts(ctx, accounts...); err != nil {
		return errors.Wrap(err, "saving vesting accounts")
	}

	vestingPeriods := make([]storage.VestingPeriod, 0)
	for i := range accounts {
		for j := range accounts[i].VestingPeriods {
			accounts[i].VestingPeriods[j].VestingAccountId = accounts[i].Id
		}
		vestingPeriods = append(vestingPeriods, accounts[i].VestingPeriods...)
	}

	if err := tx.SaveVestingPeriods(ctx, vestingPeriods...); err != nil {
		return errors.Wrap(err, "saving vesting periods")
	}

	return nil
}
