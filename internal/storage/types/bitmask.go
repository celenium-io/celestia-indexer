// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

type Bits uint64

func (b *Bits) Set(flag Bits)     { *b |= flag }
func (b *Bits) Clear(flag Bits)   { *b &^= flag }
func (b Bits) Has(flag Bits) bool { return b&flag != 0 }
func (b Bits) CountBits() int {
	var count int
	for b != 0 {
		count += int(b & 1)
		b >>= 1
	}
	return count
}
func (b Bits) Empty() bool { return b == 0 }
