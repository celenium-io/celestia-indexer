// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
)

func (tx Transaction) SaveZkISMs(ctx context.Context, items ...*models.ZkISM) error {
	if len(items) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&items).
		On("CONFLICT (external_id) DO UPDATE").
		Set("state = EXCLUDED.state").
		Returning("id").
		Exec(ctx)
	return err
}
