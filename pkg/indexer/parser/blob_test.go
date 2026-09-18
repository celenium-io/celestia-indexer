// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	blobTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// pfbBlob builds a blob log the way handle.MsgPayForBlobs does: the gas is
// already computed there, processBlobs only spreads the tx fee over it.
func pfbBlob(size int64, gas int64) *storage.BlobLog {
	return &storage.BlobLog{
		Size:        size,
		Source:      storageTypes.BlobSourcePfb,
		GasConsumed: decimal.NewFromInt(gas),
	}
}

// fibreBlob builds a blob log the way handle.MsgPayForFibre does: it already
// carries its own escrow payment and must not take a share of the tx fee.
func fibreBlob(size int64, gas int64) *storage.BlobLog {
	g := decimal.NewFromInt(gas)
	return &storage.BlobLog{
		Size:        size,
		Source:      storageTypes.BlobSourceFibre,
		GasConsumed: g,
		Fee:         storageTypes.NewNumeric(g),
	}
}

func attachedBlobs(n int) decode.DecodedTx {
	blobs := make([]*blobTypes.Blob, n)
	for i := range blobs {
		blobs[i] = &blobTypes.Blob{Data: []byte{0x01}}
	}
	return decode.DecodedTx{Blobs: blobs}
}

func blobTestTx() *storage.Tx {
	return &storage.Tx{
		Fee:     storageTypes.NumericFromInt64(100),
		GasUsed: 1000,
	}
}

func Test_processBlob(t *testing.T) {

	t.Run("one blob", func(t *testing.T) {
		blobs := []*storage.BlobLog{pfbBlob(1, 8)}
		tx := blobTestTx()

		err := processBlobs(blobs, attachedBlobs(1), tx)
		require.NoError(t, err)
		require.Equal(t, tx.Fee.String(), blobs[0].Fee.String())
	})

	t.Run("two equal blobs", func(t *testing.T) {
		blobs := []*storage.BlobLog{pfbBlob(1, 8), pfbBlob(1, 8)}
		tx := blobTestTx()

		err := processBlobs(blobs, attachedBlobs(2), tx)
		require.NoError(t, err)

		var totalFee storageTypes.Numeric
		for i := range blobs {
			totalFee = totalFee.Add(blobs[i].Fee)
		}
		require.Equal(t, tx.Fee.String(), totalFee.String())
		require.Equal(t, "50", blobs[0].Fee.String())
		require.Equal(t, "50", blobs[1].Fee.String())
	})

	t.Run("two different blobs", func(t *testing.T) {
		// 1 byte fits in one share, 1024 bytes need three: 8 and 24 gas.
		blobs := []*storage.BlobLog{pfbBlob(1, 8), pfbBlob(1024, 24)}
		tx := blobTestTx()

		err := processBlobs(blobs, attachedBlobs(2), tx)
		require.NoError(t, err)

		var totalFee storageTypes.Numeric
		for i := range blobs {
			totalFee = totalFee.Add(blobs[i].Fee)
		}
		require.Equal(t, tx.Fee.String(), totalFee.String())
		require.Equal(t, "49.2", blobs[0].Fee.String())
		require.Equal(t, "50.8", blobs[1].Fee.String())
	})

	t.Run("content type comes from the attached data", func(t *testing.T) {
		blobs := []*storage.BlobLog{pfbBlob(1, 8)}
		d := decode.DecodedTx{
			Blobs: []*blobTypes.Blob{{Data: []byte("<html><body>hi</body></html>")}},
		}

		require.NoError(t, processBlobs(blobs, d, blobTestTx()))
		require.Equal(t, "text/html; charset=utf-8", blobs[0].ContentType)
	})

	t.Run("fibre blob has no attached data and keeps its escrow fee", func(t *testing.T) {
		// A PFF transaction carries exactly one MsgPayForFibre and no blob tx,
		// so there is nothing to spread the tx fee over.
		blobs := []*storage.BlobLog{fibreBlob(262144, 695000)}

		require.NoError(t, processBlobs(blobs, decode.DecodedTx{}, blobTestTx()))
		require.Equal(t, "695000", blobs[0].Fee.String())
		require.Empty(t, blobs[0].ContentType)
	})

	t.Run("fibre blob is skipped when mixed with pfb blobs", func(t *testing.T) {
		// Defensive: even if a tx ever carried both, the fibre blob must not
		// consume an attached-data slot nor take a share of the tx fee.
		fibre := fibreBlob(262144, 695000)
		blobs := []*storage.BlobLog{fibre, pfbBlob(1, 8)}
		tx := blobTestTx()

		require.NoError(t, processBlobs(blobs, attachedBlobs(1), tx))
		require.Equal(t, tx.Fee.String(), blobs[1].Fee.String())
		require.Equal(t, "695000", fibre.Fee.String())
		require.Empty(t, fibre.ContentType)
	})

	t.Run("mismatch between pfb blobs and attached data is an error", func(t *testing.T) {
		blobs := []*storage.BlobLog{pfbBlob(1, 8)}

		err := processBlobs(blobs, attachedBlobs(2), blobTestTx())
		require.Error(t, err)
		require.Contains(t, err.Error(), "mismatch")
	})

	t.Run("nothing to do", func(t *testing.T) {
		require.NoError(t, processBlobs(nil, attachedBlobs(1), blobTestTx()))
		require.NoError(t, processBlobs([]*storage.BlobLog{pfbBlob(1, 8)}, decode.DecodedTx{}, blobTestTx()))
	})
}
