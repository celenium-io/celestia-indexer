// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"math/big"
	"math/bits"

	"github.com/pkg/errors"
)

var (
	zero = big.NewInt(0)
)

type Bits struct {
	value *big.Int
}

func NewEmptyBits() Bits {
	return Bits{
		big.NewInt(0),
	}
}

func NewBits(value int) Bits {
	return Bits{
		big.NewInt(int64(value)),
	}
}

func NewBitsWithPosition(position int) Bits {
	b := big.NewInt(0)
	b = b.SetBit(b, position, 1)
	return Bits{
		b,
	}
}

func NewBitsFromString(value string) (Bits, error) {
	b := big.NewInt(0)
	b, ok := b.SetString(value, 2)
	if !ok {
		return Bits{}, errors.Errorf("invalid mask value: %s", value)
	}
	return Bits{b}, nil
}

func (b *Bits) Set(flag Bits) { b.value = b.value.Or(b.value, flag.value) }
func (b *Bits) SetBit(position int) {
	b.value = b.value.SetBit(b.value, position, 1)
}

func (b *Bits) Clear(flag Bits) { b.value = b.value.AndNot(b.value, flag.value) }
func (b Bits) Has(flag Bits) bool {
	and := b.value.And(b.value, flag.value)
	return and.Cmp(zero) != 0
}
func (b Bits) HasBit(position int) bool {
	return b.value.Bit(position) != 0
}
func (b Bits) CountBits() int {
	var count int
	for _, x := range b.value.Bits() {
		count += bits.OnesCount(uint(x))
	}
	return count
}
func (b Bits) Empty() bool { return b.value == nil || b.value.Cmp(zero) == 0 }

func (b Bits) String() string {
	if b.value == nil {
		return "nil"
	}
	return b.value.String()
}
