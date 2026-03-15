// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

// UnmarshalBlobTxShallow — shallow BlobTx parser that caps Blob.Data at 512 bytes.
//
// Problem: tmTypes.UnmarshalBlobTx → bTx.Unmarshal (generated protobuf) allocates
// and copies the FULL blob payload into Blob.Data.  A single 1 MiB blob costs ≈1 MiB
// of heap; with many blobs per block this is the dominant allocation site in the
// indexer.
//
// Root cause: the only consumer of Blob.Data in the entire pipeline is
// http.DetectContentType(blob.Data) in parser/blob.go, which reads at most the
// first 512 bytes (net/http's sniffLen constant).  Everything else — namespace,
// size, share version, fee — comes from the MsgPayForBlobs cosmos message, not
// from the raw blob bytes.
//
// Fix: parse the BlobTx wire format manually using protowire. For the Blob.data
// field (field 2) read only the first blobDataSniffLen bytes; advance iNdEx past
// the rest without copying.  For every other field the behaviour is identical to
// the generated Unmarshal.

import (
	"github.com/celestiaorg/go-square/v3/share"
	squareTx "github.com/celestiaorg/go-square/v3/tx"
	blobTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	tmTypes "github.com/cometbft/cometbft/types"
	"google.golang.org/protobuf/encoding/protowire"
)

// blobDataSniffLen mirrors net/http.sniffLen — the maximum number of bytes
// that http.DetectContentType will ever examine.
const blobDataSniffLen = 512

// UnmarshalBlobTxShallow parses a BlobTx from its wire-format bytes while
// capping each Blob.Data at blobDataSniffLen bytes.
//
// It replicates the validation logic of tmTypes.UnmarshalBlobTx (type_id check,
// non-empty blobs, correct namespace ID length) so it can be used as a drop-in
// replacement.
func UnmarshalBlobTxShallow(tx tmTypes.Tx) (bTx blobTypes.BlobTx, isBlob bool) {
	if err := parseBlobTxShallow([]byte(tx), &bTx); err != nil {
		return blobTypes.BlobTx{}, false
	}
	if bTx.TypeId != squareTx.ProtoBlobTxTypeID {
		return bTx, false
	}
	if len(bTx.Blobs) == 0 {
		return bTx, false
	}
	for _, b := range bTx.Blobs {
		if len(b.NamespaceId) != share.NamespaceIDSize {
			return bTx, false
		}
	}
	return bTx, true
}

// parseBlobTxShallow decodes the BlobTx proto message:
//
//	message BlobTx {
//	  bytes  tx      = 1;
//	  Blob   blobs   = 2; // repeated
//	  string type_id = 3;
//	}
func parseBlobTxShallow(b []byte, out *blobTypes.BlobTx) error {
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return protowire.ParseError(n)
		}
		b = b[n:]

		switch {
		case num == 1 && typ == protowire.BytesType: // tx bytes
			v, n := protowire.ConsumeBytes(b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			out.Tx = append(out.Tx[:0], v...)
			b = b[n:]

		case num == 2 && typ == protowire.BytesType: // repeated Blob message
			v, n := protowire.ConsumeBytes(b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			blob := new(blobTypes.Blob)
			if err := parseBlobShallow(v, blob); err != nil {
				return err
			}
			out.Blobs = append(out.Blobs, blob)
			b = b[n:]

		case num == 3 && typ == protowire.BytesType: // type_id string
			v, n := protowire.ConsumeBytes(b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			out.TypeId = string(v)
			b = b[n:]

		default:
			n := protowire.ConsumeFieldValue(num, typ, b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			b = b[n:]
		}
	}
	return nil
}

// parseBlobShallow decodes a Blob proto message:
//
//	message Blob {
//	  bytes  namespace_id      = 1;
//	  bytes  data              = 2;  // truncated to blobDataSniffLen
//	  uint32 share_version     = 3;
//	  uint32 namespace_version = 4;
//	}
//
// For field 2 (data), only the first blobDataSniffLen bytes are copied into
// out.Data.  protowire.ConsumeBytes returns a zero-copy subslice of the input,
// so the large blob payload is never allocated — the pointer simply skips past
// it via b = b[n:].
func parseBlobShallow(b []byte, out *blobTypes.Blob) error {
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return protowire.ParseError(n)
		}
		b = b[n:]

		switch {
		case num == 1 && typ == protowire.BytesType: // namespace_id
			v, n := protowire.ConsumeBytes(b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			out.NamespaceId = append(out.NamespaceId[:0], v...)
			b = b[n:]

		case num == 2 && typ == protowire.BytesType: // data — TRUNCATED
			v, n := protowire.ConsumeBytes(b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			// v is a zero-copy subslice of the input buffer; capping it to
			// blobDataSniffLen does not allocate.  The subsequent append copies
			// at most 512 bytes regardless of actual blob size.
			if len(v) > blobDataSniffLen {
				v = v[:blobDataSniffLen]
			}
			out.Data = append(out.Data[:0], v...)
			b = b[n:] // advances past the FULL field, not just v[:512]

		case num == 3 && typ == protowire.VarintType: // share_version
			v, n := protowire.ConsumeVarint(b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			out.ShareVersion = uint32(v)
			b = b[n:]

		case num == 4 && typ == protowire.VarintType: // namespace_version
			v, n := protowire.ConsumeVarint(b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			out.NamespaceVersion = uint32(v)
			b = b[n:]

		default:
			n := protowire.ConsumeFieldValue(num, typ, b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			b = b[n:]
		}
	}
	return nil
}
