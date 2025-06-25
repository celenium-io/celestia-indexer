// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (module *Module) rollbackMessages(ctx context.Context, tx storage.Transaction, height types.Level) (int64, error) {
	msgs, err := tx.RollbackMessages(ctx, height)
	if err != nil {
		return 0, err
	}

	if len(msgs) == 0 {
		return 0, nil
	}

	ids := make([]uint64, len(msgs))
	for i := range msgs {
		ids[i] = msgs[i].Id
	}

	if err := tx.RollbackMessageAddresses(ctx, ids); err != nil {
		return 0, err
	}

	nsMsgs, err := tx.RollbackNamespaceMessages(ctx, height)
	if err != nil {
		return 0, err
	}
	ns, err := tx.RollbackNamespaces(ctx, height)
	if err != nil {
		return 0, err
	}

	if err := module.rollbackNamespaces(ctx, tx, nsMsgs, ns, msgs); err != nil {
		return 0, errors.Wrap(err, "namespace rollback")
	}

	return int64(len(ns)), nil
}
