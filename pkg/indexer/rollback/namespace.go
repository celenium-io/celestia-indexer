// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
)

func (module *Module) rollbackNamespaces(
	ctx context.Context,
	tx storage.Transaction,
	nsMsgs []storage.NamespaceMessage,
	deletedNs []storage.Namespace,
) error {
	if len(nsMsgs) == 0 {
		return nil
	}
	deleted := make(map[uint64]struct{}, len(deletedNs))
	for i := range deletedNs {
		deleted[deletedNs[i].Id] = struct{}{}
	}

	diffs := make(map[uint64]*storage.Namespace)
	for i := range nsMsgs {
		nsId := nsMsgs[i].NamespaceId
		if _, ok := deleted[nsId]; ok {
			continue
		}

		diff, ok := diffs[nsId]
		if !ok {
			ns, err := tx.Namespace(ctx, nsId)
			if err != nil {
				return err
			}
			// SaveNamespaces upserts by adding onto the stored row
			// (size = EXCLUDED.size + added_namespace.size), so diff must carry
			// only the negative delta this rollback removes, not the namespace's
			// post-rollback absolute totals. Keep everything but the counters,
			// which are needed only for the ON CONFLICT match (namespace_id, version).
			ns.PfbCount, ns.Size, ns.PffCount, ns.FibreSize, ns.BlobsCount = 0, 0, 0, 0, 0
			diff = &ns
			diffs[nsId] = diff
		}

		// The deleted namespace_message row carries both the size it added and
		// the source it came from, so the message payload is never parsed here.
		size := int64(nsMsgs[i].Size) //nolint:gosec // sizes are bounded by the max blob size
		diff.BlobsCount -= 1
		switch nsMsgs[i].Source {
		case storageTypes.BlobSourceFibre:
			diff.PffCount -= 1
			diff.FibreSize -= size
		default:
			diff.PfbCount -= 1
			diff.Size -= size
		}
	}

	namespaces := make([]*storage.Namespace, 0, len(diffs))
	for key := range diffs {
		last, err := tx.LastNamespaceMessage(ctx, diffs[key].Id)
		if err != nil {
			return errors.Wrap(err, "receiving last namespace message")
		}
		diffs[key].LastHeight = last.Height
		diffs[key].LastMessageTime = last.Time
		namespaces = append(namespaces, diffs[key])
	}

	_, err := tx.SaveNamespaces(ctx, namespaces...)
	return err
}
