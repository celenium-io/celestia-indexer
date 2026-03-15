// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"bytes"
	"testing"

	squareTx "github.com/celestiaorg/go-square/v3/tx"
	blobTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	tmTypes "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

// buildBlobTx encodes a BlobTx proto message containing blobs whose Data fields
// are filled with the given byte pattern, each of length dataLen.
func buildBlobTx(t testing.TB, nsID []byte, dataLen int, nBlobs int) tmTypes.Tx {
	t.Helper()

	blobData := bytes.Repeat([]byte{0xAB}, dataLen)

	// encode a single Blob message
	encodeBlob := func() []byte {
		var b []byte
		b = protowire.AppendTag(b, 1, protowire.BytesType)
		b = protowire.AppendBytes(b, nsID)
		b = protowire.AppendTag(b, 2, protowire.BytesType)
		b = protowire.AppendBytes(b, blobData)
		b = protowire.AppendTag(b, 3, protowire.VarintType)
		b = protowire.AppendVarint(b, 0) // share_version
		b = protowire.AppendTag(b, 4, protowire.VarintType)
		b = protowire.AppendVarint(b, 0) // namespace_version
		return b
	}

	fakeTx := []byte{0x01, 0x02, 0x03} // minimal cosmos tx bytes

	var blobTxBytes []byte
	blobTxBytes = protowire.AppendTag(blobTxBytes, 1, protowire.BytesType)
	blobTxBytes = protowire.AppendBytes(blobTxBytes, fakeTx)

	for range nBlobs {
		blobTxBytes = protowire.AppendTag(blobTxBytes, 2, protowire.BytesType)
		blobTxBytes = protowire.AppendBytes(blobTxBytes, encodeBlob())
	}

	blobTxBytes = protowire.AppendTag(blobTxBytes, 3, protowire.BytesType)
	blobTxBytes = protowire.AppendBytes(blobTxBytes, []byte(squareTx.ProtoBlobTxTypeID))

	return tmTypes.Tx(blobTxBytes)
}

// validNsID returns a NamespaceID of the required length (28 bytes for go-square/v3).
func validNsID() []byte {
	id := make([]byte, 28)
	id[27] = 0x01
	return id
}

// TestUnmarshalBlobTxShallow_DataTruncated verifies that Blob.Data is capped at
// blobDataSniffLen bytes even when the encoded blob is much larger.
func TestUnmarshalBlobTxShallow_DataTruncated(t *testing.T) {
	const bigBlob = 1 << 20 // 1 MiB
	tx := buildBlobTx(t, validNsID(), bigBlob, 1)

	bTx, isBlob := UnmarshalBlobTxShallow(tx)
	require.True(t, isBlob)
	require.Len(t, bTx.Blobs, 1)
	require.Equal(t, blobDataSniffLen, len(bTx.Blobs[0].Data),
		"Data must be truncated to blobDataSniffLen")
}

// TestUnmarshalBlobTxShallow_SmallBlob verifies that blobs smaller than
// blobDataSniffLen are stored in full.
func TestUnmarshalBlobTxShallow_SmallBlob(t *testing.T) {
	const smallBlob = 100
	tx := buildBlobTx(t, validNsID(), smallBlob, 1)

	bTx, isBlob := UnmarshalBlobTxShallow(tx)
	require.True(t, isBlob)
	require.Len(t, bTx.Blobs, 1)
	require.Equal(t, smallBlob, len(bTx.Blobs[0].Data),
		"Data shorter than sniff limit must be stored in full")
}

// TestUnmarshalBlobTxShallow_ContentPreserved checks that the first
// blobDataSniffLen bytes match the original blob prefix exactly.
func TestUnmarshalBlobTxShallow_ContentPreserved(t *testing.T) {
	const bigBlob = 2 << 20 // 2 MiB
	tx := buildBlobTx(t, validNsID(), bigBlob, 1)

	bTx, isBlob := UnmarshalBlobTxShallow(tx)
	require.True(t, isBlob)

	want := bytes.Repeat([]byte{0xAB}, blobDataSniffLen)
	require.Equal(t, want, bTx.Blobs[0].Data)
}

// TestUnmarshalBlobTxShallow_MultipleBlobs ensures each blob is independently
// truncated and other fields (TypeId, Tx, ShareVersion, NamespaceId) are intact.
func TestUnmarshalBlobTxShallow_MultipleBlobs(t *testing.T) {
	const (
		nBlobs  = 5
		bigBlob = 512 * 1024 // 512 KiB
	)
	nsID := validNsID()
	tx := buildBlobTx(t, nsID, bigBlob, nBlobs)

	bTx, isBlob := UnmarshalBlobTxShallow(tx)
	require.True(t, isBlob)
	require.Equal(t, squareTx.ProtoBlobTxTypeID, bTx.TypeId)
	require.Len(t, bTx.Blobs, nBlobs)

	for i, b := range bTx.Blobs {
		require.Equal(t, blobDataSniffLen, len(b.Data), "blob[%d].Data", i)
		require.Equal(t, nsID, b.NamespaceId, "blob[%d].NamespaceId", i)
		require.Equal(t, uint32(0), b.ShareVersion, "blob[%d].ShareVersion", i)
	}
}

// TestUnmarshalBlobTxShallow_ParityWithGenerated checks that fields other than
// Data are identical to what the generated bTx.Unmarshal would produce.
func TestUnmarshalBlobTxShallow_ParityWithGenerated(t *testing.T) {
	const dataLen = 200 // small enough that both parsers return the same Data
	tx := buildBlobTx(t, validNsID(), dataLen, 2)

	// shallow parser
	shallow, isBlob := UnmarshalBlobTxShallow(tx)
	require.True(t, isBlob)

	// generated parser
	var gen blobTypes.BlobTx
	require.NoError(t, gen.Unmarshal([]byte(tx)))

	require.Equal(t, gen.TypeId, shallow.TypeId)
	require.Equal(t, gen.Tx, shallow.Tx)
	require.Len(t, shallow.Blobs, len(gen.Blobs))
	for i := range gen.Blobs {
		require.Equal(t, gen.Blobs[i].NamespaceId, shallow.Blobs[i].NamespaceId)
		require.Equal(t, gen.Blobs[i].ShareVersion, shallow.Blobs[i].ShareVersion)
		require.Equal(t, gen.Blobs[i].NamespaceVersion, shallow.Blobs[i].NamespaceVersion)
		require.Equal(t, gen.Blobs[i].Data, shallow.Blobs[i].Data) // same when dataLen <= 512
	}
}
