// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	pg "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

type addedNamespace struct {
	bun.BaseModel `bun:"namespace"`
	*models.Namespace

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveNamespaces(ctx context.Context, namespaces ...*models.Namespace) (int64, error) {
	if len(namespaces) == 0 {
		return 0, nil
	}

	addedNamespaces := make([]addedNamespace, len(namespaces))
	for i := range namespaces {
		addedNamespaces[i].Namespace = namespaces[i]
	}

	_, err := tx.Tx().NewInsert().Model(&addedNamespaces).
		Column("version", "namespace_id", "pfb_count", "pff_count", "size", "fibre_size", "first_height", "last_height", "last_message_time", "blobs_count").
		On("CONFLICT ON CONSTRAINT namespace_id_version_idx DO UPDATE").
		Set("size = EXCLUDED.size + added_namespace.size").
		Set("fibre_size = EXCLUDED.fibre_size + added_namespace.fibre_size").
		Set("pfb_count = EXCLUDED.pfb_count + added_namespace.pfb_count").
		Set("pff_count = EXCLUDED.pff_count + added_namespace.pff_count").
		Set("last_height = EXCLUDED.last_height").
		Set("last_message_time = EXCLUDED.last_message_time").
		Set("blobs_count = EXCLUDED.blobs_count + added_namespace.blobs_count").
		Returning("xmax, id").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	for i := range addedNamespaces {
		if addedNamespaces[i].Xmax == 0 {
			count++
		}
	}

	return count, err
}

func (tx Transaction) SaveBlobLogs(ctx context.Context, logs ...*models.BlobLog) error {
	return pg.SaveBulkWithCopy(ctx, tx, logs, copyThreshold)
}

func (tx Transaction) RollbackNamespaceMessages(ctx context.Context, height types.Level) (msgs []models.NamespaceMessage, err error) {
	_, err = tx.Tx().NewDelete().Model(&msgs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackNamespaces(ctx context.Context, height types.Level) (ns []models.Namespace, err error) {
	_, err = tx.Tx().NewDelete().Model(&ns).Where("first_height = ?", height).Returning("*").Exec(ctx)
	return
}
