// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	blobTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test_processBlob(t *testing.T) {

	t.Run("one blob", func(t *testing.T) {
		blobs := []*storage.BlobLog{
			{
				Size: 1,
			},
		}
		d := decode.DecodedTx{
			Blobs: []*blobTypes.Blob{
				{
					Data: []byte{0x01},
				},
			},
		}
		tx := &storage.Tx{
			Fee:     decimal.RequireFromString("100"),
			GasUsed: 1000,
		}

		processBlob(blobs, d, tx)
		require.Equal(t, tx.Fee.String(), blobs[0].Fee.String())
	})

	t.Run("two equal blobs", func(t *testing.T) {
		blobs := []*storage.BlobLog{
			{
				Size: 1,
			}, {
				Size: 1,
			},
		}
		d := decode.DecodedTx{
			Blobs: []*blobTypes.Blob{
				{
					Data: []byte{0x01},
				}, {
					Data: []byte{0x01},
				},
			},
		}
		tx := &storage.Tx{
			Fee:     decimal.RequireFromString("100"),
			GasUsed: 1000,
		}

		processBlob(blobs, d, tx)

		totalFee := decimal.Zero
		for i := range blobs {
			totalFee = totalFee.Add(blobs[i].Fee)
		}
		require.Equal(t, tx.Fee.String(), totalFee.String())
		require.Equal(t, "50", blobs[0].Fee.String())
		require.Equal(t, "50", blobs[1].Fee.String())
	})

	t.Run("two different blobs", func(t *testing.T) {
		blobs := []*storage.BlobLog{
			{
				Size: 1,
			}, {
				Size: 1024,
			},
		}
		d := decode.DecodedTx{
			Blobs: []*blobTypes.Blob{
				{
					Data: []byte{0x01},
				}, {
					Data: []byte{0x01},
				},
			},
		}
		tx := &storage.Tx{
			Fee:     decimal.RequireFromString("100"),
			GasUsed: 1000,
		}

		processBlob(blobs, d, tx)

		totalFee := decimal.Zero
		for i := range blobs {
			totalFee = totalFee.Add(blobs[i].Fee)
		}
		require.Equal(t, tx.Fee.String(), totalFee.String())
		require.Equal(t, "49.2", blobs[0].Fee.String())
		require.Equal(t, "50.8", blobs[1].Fee.String())
	})
}
