// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"testing"

	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	squareTx "github.com/celestiaorg/go-square/v3/tx"
	"github.com/cometbft/cometbft/crypto/tmhash"
	blobProto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmTypes "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

// buildIndexWrapperTx encodes an IndexWrapper proto message around innerTx.
//
//	message IndexWrapper {
//	  bytes           tx            = 1;
//	  repeated uint32 share_indexes = 2; // packed
//	  string          type_id       = 3;
//	}
func buildIndexWrapperTx(t testing.TB, innerTx []byte, shareIndexes []uint32) tmTypes.Tx {
	t.Helper()

	var b []byte
	b = protowire.AppendTag(b, 1, protowire.BytesType)
	b = protowire.AppendBytes(b, innerTx)

	if len(shareIndexes) > 0 {
		var packed []byte
		for _, idx := range shareIndexes {
			packed = protowire.AppendVarint(packed, uint64(idx))
		}
		b = protowire.AppendTag(b, 2, protowire.BytesType)
		b = protowire.AppendBytes(b, packed)
	}

	b = protowire.AppendTag(b, 3, protowire.BytesType)
	b = protowire.AppendBytes(b, []byte(squareTx.ProtoIndexWrapperTypeID))

	return tmTypes.Tx(b)
}

// minimalDeliverTx is a zero-value result used when the test only cares about
// the decoded hash and not the block-results fields.
var minimalDeliverTx = nodeTypes.ResponseDeliverTx{}

// mustDecodeB64 is a test helper that base64-decodes a string or fails the test.
func mustDecodeB64(t testing.TB, s string) []byte {
	t.Helper()
	b, err := base64.StdEncoding.DecodeString(s)
	require.NoError(t, err)
	return b
}

// ── plain cosmos tx ──────────────────────────────────────────────────────────

// TestDecodeTx_Hash_PlainTx verifies that for a plain (non-wrapped) cosmos tx
// the hash equals tmhash.Sum(rawBytes) — no stripping or transformation.
func TestDecodeTx_Hash_PlainTx(t *testing.T) {
	// MsgVote — a representative plain cosmos tx with no BlobTx/IndexWrapper envelope.
	rawTx := mustDecodeB64(t, "ClEKTwoWL2Nvc21vcy5nb3YudjEuTXNnVm90ZRI1CAQSL2NlbGVzdGlhMTJ6czdlM244cGpkOHk4ZXgwY3l2NjdldGh2MzBtZWtncXU2NjVyGAESaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAyJq13zdVvBc4sHiqsxdwmZuhu/+7jp5qybynJAUP4VeEgQKAggBGEQSFAoOCgR1dGlhEgY1MDAwMDAQ85EFGkB0CjjkpeDX/bfNeifAKWUWMSf5l7l8DqsDosnuQK3XMjiTlXN4AthomxLpSDqS/i7fsV7cLnaKV2trwJR5FvTc")

	block, _ := testsuite.CreateBlockWithTxs(minimalDeliverTx, rawTx, 1)
	dTx, err := Tx(block, 0)
	require.NoError(t, err)

	wantHash := tmhash.Sum(rawTx)
	require.Equal(t, wantHash, dTx.Hash)
}

// ── BlobTx ───────────────────────────────────────────────────────────────────

// TestDecodeTx_Hash_BlobTx verifies that for a BlobTx-wrapped transaction the
// hash equals tmhash.Sum(blobTx.Tx) — i.e. the inner cosmos-sdk tx bytes —
// and NOT tmhash.Sum(rawBytes) (hash of the full BlobTx envelope).
func TestDecodeTx_Hash_BlobTx(t *testing.T) {
	// Real MsgPayForBlobs transaction — a BlobTx with ProtoBlobTxTypeID.
	rawTx := mustDecodeB64(t, "CoUCCqABCp0BCiAvY2VsZXN0aWEuYmxvYi52MS5Nc2dQYXlGb3JCbG9icxJ5Ci9jZWxlc3RpYTFya3k5MDg2dDM0MG03cm1rY3R1ajRzcHh3djJnYzYydmx3eDU5dhIdAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ2Vyb0EaAqwFIiA92mk96XJQMA82kZz4lDP5Fbj4U7ss8LisNXzMW00q0kIBABIeCgkSBAoCCAEYjQUSEQoLCgR1dGlhEgMxODUQ+dAFGkCYFhvYyED7gTt9JbqSSJSFsQfgBcFU/H6n35PgNgZvWUp9EDMknrBVwRNwdHX00Ald9brD/Ir34FDdJAfc8p/tEs0FChwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAENlcm9BEqwFeyJnbG9iYWxTZXF1ZW5jZU51bWJlciI6MTAzOSwiYmxvY2tSYW5nZUNvdmVyZWQiOnsiYmxvY2tTdGFydCI6MzczNjgwMCwiYmxvY2tFbmQiOjM3NDA0MDB9LCJ0aW1lc3RhbXAiOjE3NDcxMzkwNDAsInJvbGx1cFNlcXVlbmNlcyI6W3sicm9sbHVwSWQiOjEsImJhdGNoZXMiOlt7InRpbWVzdGFtcCI6MTc0NzEzOTA0MCwibnVtYmVyIjoxMDM5fV19LHsicm9sbHVwSWQiOjIsImJhdGNoZXMiOlt7InRpbWVzdGFtcCI6MTc0NzEzOTA0MCwidHJhbnNhY3Rpb25zIjpbeyJ0eElkIjoiMDM5YTRmYWIwNmQ0Nzg1OTRmYWUxMmQ5NGYxN2I1ZTlmNDUyMTdiZDUzNDc3YTc5Y2I1MWY2YTQzZDY0ZDQzMSIsInJhd1RyYW5zYWN0aW9uIjoie2Zyb206c29tZXRlc3RhZGRyZXNzLCB0bzogc29tZW9uZWVsc2UsIGFtb3VudDogMC4wMDUsIHRpbWU6MTc0NzEzNTQwNyB9IiwiYmxvY2tIZWlnaHQiOjM3Mzc5MDMsInJvbGx1cElkIjoyfV0sIm51bWJlciI6MTAzOX1dfSx7InJvbGx1cElkIjozLCJiYXRjaGVzIjpbeyJ0aW1lc3RhbXAiOjE3NDcxMzkwNDAsIm51bWJlciI6MTAzOX1dfSx7InJvbGx1cElkIjo0LCJiYXRjaGVzIjpbeyJ0aW1lc3RhbXAiOjE3NDcxMzkwNDAsIm51bWJlciI6MTAzOX1dfSx7InJvbGx1cElkIjo1LCJiYXRjaGVzIjpbeyJ0aW1lc3RhbXAiOjE3NDcxMzkwNDAsIm51bWJlciI6MTAzOX1dfV19GgRCTE9C")

	// Derive expected hash using the generated proto unmarshaler (reference impl),
	// not our shallow parser — this is what makes the test non-circular.
	var bTx blobProto.BlobTx
	require.NoError(t, bTx.Unmarshal(rawTx))
	require.NotEmpty(t, bTx.Tx, "reference unmarshal must produce inner tx bytes")

	wantHash := tmhash.Sum(bTx.Tx)

	block, _ := testsuite.CreateBlockWithTxs(minimalDeliverTx, rawTx, 1)
	dTx, err := Tx(block, 0)
	require.NoError(t, err)

	require.Equal(t, wantHash, dTx.Hash,
		"BlobTx: hash must be tmhash.Sum(innerTx), not tmhash.Sum(fullEnvelope)")

	// Explicit: hash of the full envelope must differ from hash of the inner tx.
	require.False(t, bytes.Equal(tmhash.Sum(rawTx), dTx.Hash),
		"BlobTx: hashing the full envelope bytes must not equal the inner-tx hash")
}

// ── IndexWrapper ─────────────────────────────────────────────────────────────

// TestDecodeTx_Hash_IndexWrapper verifies that for an IndexWrapper-wrapped
// transaction the hash equals tmhash.Sum(indexWrapper.Tx) — the inner tx bytes —
// and NOT tmhash.Sum(rawBytes) (hash of the full IndexWrapper envelope).
//
// The IndexWrapper is constructed synthetically by wrapping a known valid cosmos
// tx so that we control the expected inner bytes precisely.
func TestDecodeTx_Hash_IndexWrapper(t *testing.T) {
	// Use a real plain cosmos tx as the inner payload — MsgVote from the existing
	// test corpus. Its bytes are valid for txDecoder, so the full decode succeeds.
	innerTxBytes := mustDecodeB64(t, "ClEKTwoWL2Nvc21vcy5nb3YudjEuTXNnVm90ZRI1CAQSL2NlbGVzdGlhMTJ6czdlM244cGpkOHk4ZXgwY3l2NjdldGh2MzBtZWtncXU2NjVyGAESaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAyJq13zdVvBc4sHiqsxdwmZuhu/+7jp5qybynJAUP4VeEgQKAggBGEQSFAoOCgR1dGlhEgY1MDAwMDAQ85EFGkB0CjjkpeDX/bfNeifAKWUWMSf5l7l8DqsDosnuQK3XMjiTlXN4AthomxLpSDqS/i7fsV7cLnaKV2trwJR5FvTc")

	rawTx := buildIndexWrapperTx(t, innerTxBytes, []uint32{0, 1, 2})
	wantHash := tmhash.Sum(innerTxBytes)

	block, _ := testsuite.CreateBlockWithTxs(minimalDeliverTx, rawTx, 1)
	dTx, err := Tx(block, 0)
	require.NoError(t, err)

	require.Equal(t, wantHash, dTx.Hash,
		"IndexWrapper: hash must be tmhash.Sum(innerTx), not tmhash.Sum(fullEnvelope)")

	// Explicit: the full-envelope hash must differ — proves stripping occurred.
	require.False(t, bytes.Equal(tmhash.Sum(rawTx), dTx.Hash),
		"IndexWrapper: hashing the full envelope bytes must not equal the inner-tx hash")
}

// TestDecodeTx_Hash_IndexWrapper_NoShareIndexes checks IndexWrapper without the
// optional share_indexes field — the hash contract must hold regardless.
func TestDecodeTx_Hash_IndexWrapper_NoShareIndexes(t *testing.T) {
	innerTxBytes := mustDecodeB64(t, "ClEKTwoWL2Nvc21vcy5nb3YudjEuTXNnVm90ZRI1CAQSL2NlbGVzdGlhMTJ6czdlM244cGpkOHk4ZXgwY3l2NjdldGh2MzBtZWtncXU2NjVyGAESaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAyJq13zdVvBc4sHiqsxdwmZuhu/+7jp5qybynJAUP4VeEgQKAggBGEQSFAoOCgR1dGlhEgY1MDAwMDAQ85EFGkB0CjjkpeDX/bfNeifAKWUWMSf5l7l8DqsDosnuQK3XMjiTlXN4AthomxLpSDqS/i7fsV7cLnaKV2trwJR5FvTc")

	rawTx := buildIndexWrapperTx(t, innerTxBytes, nil)
	wantHash := tmhash.Sum(innerTxBytes)

	block, _ := testsuite.CreateBlockWithTxs(minimalDeliverTx, rawTx, 1)
	dTx, err := Tx(block, 0)
	require.NoError(t, err)

	require.Equal(t, wantHash, dTx.Hash)
}

func TestDecodeTx_Hash_Table(t *testing.T) {
	tests := []struct {
		name     string
		rawTx    string
		wantHash string
	}{
		{
			name:     "e0ff6261629b88f1802c3a96207b2b665bc445371d6744e73805bf59e6e5579d mainnet",
			rawTx:    "Cv8CCp8BCpwBCiAvY2VsZXN0aWEuYmxvYi52MS5Nc2dQYXlGb3JCbG9icxJ4Ci9jZWxlc3RpYTFuZWZ1a2h5dzhyNnBsYWdsMHdodTU0dXRyZDhrdnlndnVoNDAyaxIdAAAAAAAAAAAAAAAAAAAAAAAAAHNvbGF4eS1zb3YaAQ0iIFtL6ieX+4oFNR7qM5qDNJsoSc8BZAQjZVxHcnA/tdpZQgEBEpgBClIKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDgK7ap458cN7nQS/wShv8o4t/YTEFUuRHwUNYVjP5Uu0SBAoCCAEYlIATEkIKCwoEdXRpYRIDMjM4EMHPAxovY2VsZXN0aWExbmVmdWtoeXc4cjZwbGFnbDB3aHU1NHV0cmQ4a3Z5Z3Z1aDQwMmsaQMqr5s/VmxDM+BA/69rPB6Z6xd5HCrr1jThIyCMCUkhGDxcmM/0kUVhobfjmMDqdELYlKJXN5CbVhrrLmzzHrDoSRQocAAAAAAAAAAAAAAAAAAAAAAAAc29sYXh5LXNvdhINRLYEAAAAAAAAAAAAARgBKhSeU8tcjjj0H/Ufe6/KV4sbT2YRDBoEQkxPQg==",
			wantHash: "e0ff6261629b88f1802c3a96207b2b665bc445371d6744e73805bf59e6e5579d",
		}, {
			name:     "a54696917a20d22a842475bc9e34b98cbda683286e342a01b58c9f04737cf3ff mainnet",
			rawTx:    "CpgICpUICikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchLnBwoIdHJhbnNmZXISCWNoYW5uZWwtMhoPCgR1dGlhEgcxNjEwMDAwIi9jZWxlc3RpYTFhZXZoNm5mNHk3cTV5dGg0anl1NzgwdHBwcTQ3cHNyaHAwa2hodSo/b3NtbzEwYTNrNGh2azM3Y2M0aG54Y3R3NHA5NWZoc2NkMno2aDJybXgwYXVrYzZybTh1OXFxeDlzbWZzaDd1MgA4mOaznI/Sws4YQsAGeyJ3YXNtIjp7ImNvbnRyYWN0Ijoib3NtbzEwYTNrNGh2azM3Y2M0aG54Y3R3NHA5NWZoc2NkMno2aDJybXgwYXVrYzZybTh1OXFxeDlzbWZzaDd1IiwibXNnIjp7InN3YXBfYW5kX2FjdGlvbiI6eyJ1c2VyX3N3YXAiOnsic3dhcF9leGFjdF9hc3NldF9pbiI6eyJzd2FwX3ZlbnVlX25hbWUiOiJvc21vc2lzLXBvb2xtYW5hZ2VyIiwib3BlcmF0aW9ucyI6W3sicG9vbCI6IjEyNDgiLCJkZW5vbV9pbiI6ImliYy9ENzlFN0Q4M0FCMzk5QkZGRjkzNDMzRTU0RkFBNDgwQzE5MTI0OEZDNTU2OTI0QTJBODM1MUFFMjYzOEIzODc3IiwiZGVub21fb3V0IjoidW9zbW8ifSx7InBvb2wiOiI5IiwiZGVub21faW4iOiJ1b3NtbyIsImRlbm9tX291dCI6ImliYy9FNjkzMUY3ODA1N0Y3Q0M1REEwRkQ2Q0VGODJGRjM5MzczQTZFMDQ1MkJGMUZENzY5MTBCOTMyOTJDRjM1NkMxIn1dfX0sIm1pbl9hc3NldCI6eyJuYXRpdmUiOnsiZGVub20iOiJpYmMvRTY5MzFGNzgwNTdGN0NDNURBMEZENkNFRjgyRkYzOTM3M0E2RTA0NTJCRjFGRDc2OTEwQjkzMjkyQ0YzNTZDMSIsImFtb3VudCI6IjcwOTMwNDE3OCJ9fSwidGltZW91dF90aW1lc3RhbXAiOjE3NzM1ODU0NDU5MzE1MzA4MzAsInBvc3Rfc3dhcF9hY3Rpb24iOnsiaWJjX3RyYW5zZmVyIjp7ImliY19pbmZvIjp7InNvdXJjZV9jaGFubmVsIjoiY2hhbm5lbC01IiwicmVjZWl2ZXIiOiJjcm8xYWV2aDZuZjR5N3E1eXRoNGp5dTc4MHRwcHE0N3BzcmhnNzA3M3EiLCJtZW1vIjoiIiwicmVjb3Zlcl9hZGRyZXNzIjoib3NtbzFhZXZoNm5mNHk3cTV5dGg0anl1NzgwdHBwcTQ3cHNyaGM3NWhtciJ9fX0sImFmZmlsaWF0ZXMiOltdfX19fRJnClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDXtBYcJDqVYk0iaaf39uoQZmF0N9iaygiiJUUM42gXxESBAoCCH8Y7gESEgoMCgR1dGlhEgQxMzQwELuWCBpAWuL7ros6mZmtodW7+oFrb5EjIGCqCWuksqlTxeQ9/o9WxcfAi4hq6b5nNzCa613/aTvJWnUp03+cTyFgryBa9Q==",
			wantHash: "a54696917a20d22a842475bc9e34b98cbda683286e342a01b58c9f04737cf3ff",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rawTx := mustDecodeB64(t, tc.rawTx)
			block, _ := testsuite.CreateBlockWithTxs(minimalDeliverTx, rawTx, 1)
			dTx, err := Tx(block, 0)
			require.NoError(t, err)
			hash := hex.EncodeToString(dTx.Hash)
			require.Equal(t, tc.wantHash, hash)
		})
	}
}
