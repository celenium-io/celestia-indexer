// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func (tx Transaction) UpdateRollup(ctx context.Context, rollup *models.Rollup) error {
	if rollup == nil || (rollup.IsEmpty() && !rollup.Verified) {
		return nil
	}

	query := tx.Tx().NewUpdate().Model(rollup).WherePK()

	if rollup.Name != "" {
		query = query.Set("name = ?", rollup.Name)
	}
	if rollup.Slug != "" {
		query = query.Set("slug = ?", rollup.Slug)
	}
	if rollup.Description != "" {
		query = query.Set("description = ?", rollup.Description)
	}
	if rollup.Twitter != "" {
		query = query.Set("twitter = ?", rollup.Twitter)
	}
	if rollup.GitHub != "" {
		query = query.Set("github = ?", rollup.GitHub)
	}
	if rollup.Website != "" {
		query = query.Set("website = ?", rollup.Website)
	}
	if rollup.Logo != "" {
		query = query.Set("logo = ?", rollup.Logo)
	}
	if rollup.L2Beat != "" {
		query = query.Set("l2_beat = ?", rollup.L2Beat)
	}
	if rollup.Explorer != "" {
		query = query.Set("explorer = ?", rollup.Explorer)
	}
	if rollup.BridgeContract != "" {
		query = query.Set("bridge_contract = ?", rollup.BridgeContract)
	}
	if rollup.Stack != "" {
		query = query.Set("stack = ?", rollup.Stack)
	}
	if rollup.Links != nil {
		query = query.Set("links = ?", pgdialect.Array(rollup.Links))
	}
	if rollup.Type != "" {
		query = query.Set("type = ?", rollup.Type)
	}
	if rollup.Category != "" {
		query = query.Set("category = ?", rollup.Category)
	}
	if rollup.Tags != nil {
		query = query.Set("tags = ?", pgdialect.Array(rollup.Tags))
	}
	if rollup.Provider != "" {
		query = query.Set("provider = ?", rollup.Provider)
	}
	if rollup.Compression != "" {
		query = query.Set("compression = ?", rollup.Compression)
	}
	if rollup.VM != "" {
		query = query.Set("vm = ?", rollup.VM)
	}
	if rollup.DeFiLama != "" {
		query = query.Set("defi_lama = ?", rollup.DeFiLama)
	}
	if rollup.SettledOn != "" {
		query = query.Set("settled_on = ?", rollup.SettledOn)
	}
	if rollup.Color != "" {
		query = query.Set("color = ?", rollup.Color)
	}

	query = query.Set("verified = ?", rollup.Verified)

	_, err := query.Exec(ctx)
	return err
}

func (tx Transaction) DeleteProviders(ctx context.Context, rollupId uint64) error {
	if rollupId == 0 {
		return nil
	}
	_, err := tx.Tx().NewDelete().
		Model((*models.RollupProvider)(nil)).
		Where("rollup_id = ?", rollupId).
		Exec(ctx)
	return err
}

func (tx Transaction) DeleteRollup(ctx context.Context, rollupId uint64) error {
	if rollupId == 0 {
		return nil
	}
	_, err := tx.Tx().NewDelete().
		Model((*models.Rollup)(nil)).
		Where("id = ?", rollupId).
		Exec(ctx)
	return err
}
