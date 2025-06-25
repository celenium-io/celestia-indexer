// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func saveNamespaces(
	ctx context.Context,
	tx storage.Transaction,
	namespaces map[string]*storage.Namespace,
) (int64, error) {
	if len(namespaces) == 0 {
		return 0, nil
	}

	data := make([]*storage.Namespace, 0, len(namespaces))
	for key := range namespaces {
		data = append(data, namespaces[key])
	}

	totalNamespaces, err := tx.SaveNamespaces(ctx, data...)
	if err != nil {
		return 0, err
	}

	return totalNamespaces, nil
}
