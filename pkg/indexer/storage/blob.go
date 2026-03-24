// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

func saveBlobLogs(
	ctx context.Context,
	tx storage.Transaction,
	blobs []*storage.BlobLog,
	addrToId map[string]uint64,
) error {
	if len(blobs) == 0 {
		return nil
	}

	for i := range blobs {
		if err := processPayForBlob(addrToId, blobs[i]); err != nil {
			return err
		}
	}

	return tx.SaveBlobLogs(ctx, blobs...)
}

func processPayForBlob(addrToId map[string]uint64, blob *storage.BlobLog) error {
	if blob.Namespace == nil {
		return errors.New("nil namespace in pay for blob message")
	}
	if blob.Signer == nil {
		return errors.New("nil signer address in pay for blob message")
	}
	signerId, ok := addrToId[blob.Signer.Address]
	if !ok {
		return errors.Wrapf(errCantFindAddress, "signer for pay for blob message: %s", blob.Signer.Address)
	}
	blob.SignerId = signerId
	blob.NamespaceId = blob.Namespace.Id
	return nil
}
