// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"
)

type addedAddress struct {
	bun.BaseModel `bun:"address"`
	*models.Address

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveAddresses(ctx context.Context, addresses ...*models.Address) (int64, error) {
	if len(addresses) == 0 {
		return 0, nil
	}

	addr := make([]addedAddress, len(addresses))
	for i := range addresses {
		addr[i].Address = addresses[i]
	}

	_, err := tx.Tx().NewInsert().Model(&addr).
		Column("address", "height", "last_height", "hash", "name", "is_forwarding").
		On("CONFLICT ON CONSTRAINT address_idx DO UPDATE").
		Set("last_height = EXCLUDED.last_height").
		Set("is_forwarding = ?TableAlias.is_forwarding OR EXCLUDED.is_forwarding").
		Returning("xmax, id").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	for i := range addr {
		if addr[i].Xmax == 0 {
			count++
		}
	}

	return count, err
}

func (tx Transaction) SaveBalances(ctx context.Context, balances ...models.Balance) error {
	if len(balances) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&balances).
		Column("id", "currency", "spendable", "delegated", "unbonding").
		On("CONFLICT (id, currency) DO UPDATE").
		Set("spendable = EXCLUDED.spendable + balance.spendable").
		Set("delegated = EXCLUDED.delegated + balance.delegated").
		Set("unbonding = EXCLUDED.unbonding + balance.unbonding").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveVestingAccounts(ctx context.Context, accs ...*models.VestingAccount) error {
	if len(accs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&accs).Returning("id").Exec(ctx)
	return err
}

func (tx Transaction) SaveGrants(ctx context.Context, grants ...*models.Grant) error {
	if len(grants) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().
		Model(&grants).
		Column("height", "time", "granter_id", "grantee_id", "authorization", "expiration", "revoked", "revoke_height", "params").
		On("CONFLICT ON CONSTRAINT grant_key DO UPDATE").
		Set("revoked = EXCLUDED.revoked").
		Set("revoke_height = EXCLUDED.revoke_height").
		Exec(ctx)
	return err
}

func (tx Transaction) RollbackAddresses(ctx context.Context, height types.Level) (address []models.Address, err error) {
	_, err = tx.Tx().NewDelete().Model(&address).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackGrants(ctx context.Context, height types.Level) (err error) {
	if _, err = tx.Tx().NewDelete().
		Model((*models.Grant)(nil)).
		Where("height = ?", height).
		Exec(ctx); err != nil {
		return err
	}

	_, err = tx.Tx().NewUpdate().
		Model((*models.Grant)(nil)).
		Where("revoke_height = ?", height).
		Set("revoked = false").
		Set("revoke_height = null").
		Exec(ctx)
	return
}

func (tx Transaction) DeleteBalances(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := tx.Tx().NewDelete().
		Model((*models.Balance)(nil)).
		Where("id IN ?", bun.Tuple(ids)).
		Exec(ctx)
	return err
}
