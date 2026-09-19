// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"net/http"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func processBlobs(blobs []*storage.BlobLog, d decode.DecodedTx, t *storage.Tx) error {
	if len(blobs) == 0 || len(d.Blobs) == 0 {
		return nil
	}

	pfbBlobs := make([]*storage.BlobLog, 0, len(blobs))
	for i := range blobs {
		if blobs[i].Source == storageTypes.BlobSourcePfb {
			pfbBlobs = append(pfbBlobs, blobs[i])
		}
	}

	if len(pfbBlobs) != len(d.Blobs) {
		return errors.Errorf("parsed PFB blobs and attached data count mismatch: %d != %d", len(pfbBlobs), len(d.Blobs))
	}

	var (
		gasConsumedOnBlobs = decimal.Zero.Copy()
		count              int64
	)
	for i := range pfbBlobs {
		if pfbBlobs[i].Source == storageTypes.BlobSourcePfb {
			pfbBlobs[i].ContentType = http.DetectContentType(d.Blobs[i].Data)
			gasConsumedOnBlobs = gasConsumedOnBlobs.Add(pfbBlobs[i].GasConsumed)
			count++
		}
	}

	gasUsed := decimal.NewFromInt(t.GasUsed)

	// fix_gas_per_blob = (gas_used - consumed_gas_on_blobs) / blobs_count
	fix := gasUsed.Copy().
		Sub(gasConsumedOnBlobs).
		Div(decimal.NewFromInt(count))

	for i := range pfbBlobs {
		if pfbBlobs[i].Source == storageTypes.BlobSourcePfb {
			// share_in_gas = (gas_consumed_on_blob + fix_gas_per_blob) / gas_used
			share := pfbBlobs[i].GasConsumed.Add(fix).Div(gasUsed)
			pfbBlobs[i].Fee = t.Fee.Copy().Mul(storageTypes.NewNumeric(share))
		}
	}

	return nil
}
