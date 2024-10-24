// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"net/http"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celestiaorg/celestia-app/v3/pkg/appconsts"
	appshares "github.com/celestiaorg/go-square/v2/share"
	"github.com/shopspring/decimal"
)

var gasPerBlobByte = decimal.NewFromInt(int64(appconsts.DefaultGasPerBlobByte))

func processBlob(blobs []*storage.BlobLog, d decode.DecodedTx, t *storage.Tx) {
	if len(blobs) == 0 || len(d.Blobs) != len(blobs) {
		return
	}
	t.BlobsCount += len(blobs)

	var (
		gasConsumedOnBlobs = decimal.Zero.Copy()
		gasConsumedPerBlob = make([]decimal.Decimal, len(blobs))
	)
	for i := range blobs {
		blobs[i].ContentType = http.DetectContentType(d.Blobs[i].Data)
		//nolint:gosec
		sharesUsed := appshares.SparseSharesNeeded(uint32(blobs[i].Size))
		gas := decimal.NewFromInt(int64(sharesUsed)).Mul(gasPerBlobByte)
		gasConsumedOnBlobs = gasConsumedOnBlobs.Add(gas)
		gasConsumedPerBlob[i] = gas

	}

	gasUsed := decimal.NewFromInt(t.GasUsed)

	// fix_gas_per_blob = (gas_used - consumed_gas_on_blobs) / blobs_count
	fix := gasUsed.Copy().
		Sub(gasConsumedOnBlobs).
		Div(decimal.NewFromInt(int64(len(blobs))))

	for i := range gasConsumedPerBlob {
		// share_in_gas = (gas_consumed_on_blob + fix_gas_per_blob) / gas_used
		share := gasConsumedPerBlob[i].Add(fix).Div(gasUsed)
		blobs[i].Fee = t.Fee.Copy().Mul(share)
	}
}
