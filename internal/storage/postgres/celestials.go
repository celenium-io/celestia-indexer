// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"iter"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdkStorage "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

type Celestials struct {
	*database.Bun
}

func NewCelestials(db *database.Bun) *Celestials {
	return &Celestials{
		Bun: db,
	}
}

func (c *Celestials) ById(ctx context.Context, id string) (result storage.Celestial, err error) {
	err = c.DB().NewSelect().
		Model(&result).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	return
}

func (c *Celestials) ByAddressId(ctx context.Context, addressId uint64, limit, offset int) (result []storage.Celestial, err error) {
	query := c.DB().NewSelect().
		Model(&result).
		Where("address_id = ?", addressId).
		Offset(offset).
		OrderExpr("change_id desc")

	query = limitScope(query, limit)

	err = query.Scan(ctx)
	return
}

type CelestialState struct {
	db *database.Bun
}

func NewCelestialState(db *database.Bun) *CelestialState {
	return &CelestialState{
		db: db,
	}
}

func (cs *CelestialState) ByName(ctx context.Context, name string) (result storage.CelestialState, err error) {
	err = cs.db.DB().NewSelect().
		Model(&result).
		Where("name = ?", name).
		Limit(1).
		Scan(ctx)
	return
}

func (cs *CelestialState) Save(ctx context.Context, state *storage.CelestialState) error {
	_, err := cs.db.DB().NewInsert().
		Model(state).
		Exec(ctx)
	return err
}

type CelestialTransaction struct {
	sdkStorage.Transaction
}

func BeginCelestialTransaction(ctx context.Context, tx sdkStorage.Transactable) (CelestialTransaction, error) {
	t, err := tx.BeginTransaction(ctx)
	return CelestialTransaction{t}, err
}

func (tx CelestialTransaction) SaveCelestials(ctx context.Context, celestials iter.Seq[storage.Celestial]) error {
	for cel := range celestials {
		_, err := tx.Tx().NewInsert().
			Model(&cel).
			Column("id", "address_id", "image_url", "change_id").
			On("CONFLICT (id) DO UPDATE").
			Set("address_id = EXCLUDED.address_id").
			Set("image_url = EXCLUDED.image_url").
			Set("change_id = EXCLUDED.change_id").
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tx CelestialTransaction) UpdateState(ctx context.Context, state *storage.CelestialState) error {
	_, err := tx.Tx().NewUpdate().
		Model(state).
		Set("change_id = ?", state.ChangeId).
		WherePK().
		Exec(ctx)
	return err
}
