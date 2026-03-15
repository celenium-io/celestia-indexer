// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rpc

// Custom jsoniter decoder for tmTypes.Tx (= []byte).
//
// Why this exists:
//
// The default jsoniter base64Codec decodes a []byte JSON field in three steps:
//   1. ReadString()     → readStringSlowPath() → string(str)   [alloc: base64 chars as string]
//   2. []byte(s)        in base64.DecodeString                  [alloc: copy string→bytes again]
//   3. dbuf             in base64.DecodeString                  [alloc: decoded bytes]
//
// For a 1 MB blob tx (≈1.33 MB base64) that is three allocations ≈ 4 MB.
//
// ReadStringAsSlice() eliminates copies 1 and 2:
//   - Fast path: string fits in the iterator's read buffer → zero-copy slice, no allocation.
//   - Slow path: string spans buffer boundaries → one []byte allocation (same as readStringSlowPath),
//     but skips the wasteful string(str) round-trip and the []byte(s) re-conversion.
//
// base64.Decode(dst, src []byte) then writes directly into a pre-sized destination,
// eliminating copy 3's internal string→[]byte conversion.
//
// No escape handling is needed: the base64 alphabet {A-Z, a-z, 0-9, +, /, =}
// contains no JSON special characters, so ReadStringAsSlice's "no escape" fast
// path is always valid for transaction data.

import (
	"encoding/base64"
	"unsafe"

	tmTypes "github.com/cometbft/cometbft/types"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	jsoniter.RegisterTypeDecoder("types.Tx", &txJSONDecoder{})
}

type txJSONDecoder struct{}

func (d *txJSONDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.ReadNil() {
		*(*tmTypes.Tx)(ptr) = nil
		return
	}

	// ReadStringAsSlice reads the JSON string without building an intermediate Go string:
	//   - If the entire base64 value fits in iter's current buffer window, it returns
	//     a zero-copy slice directly into that buffer. The caller must not hold the
	//     slice past the next iterator call — base64.Decode consumes it immediately.
	//   - Otherwise it falls back to a single []byte allocation, skipping the
	//     string([]byte) and []byte(string) round-trips of the default codec.
	b64 := iter.ReadStringAsSlice()
	if iter.Error != nil {
		return
	}

	// Pre-allocate exactly the space needed for the decoded result.
	// DecodedLen may over-estimate by up to 2 bytes (padding); decoded[:n] trims that.
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(b64)))
	n, err := base64.StdEncoding.Decode(decoded, b64)
	if err != nil {
		iter.ReportError("Tx.Decode", err.Error())
		return
	}

	*(*tmTypes.Tx)(ptr) = decoded[:n]
}
