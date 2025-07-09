// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBits_Set(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		flag Bits
		want Bits
	}{
		{
			name: "test 1",
			b:    NewEmptyBits(),
			flag: NewBits(1),
			want: NewBits(1),
		}, {
			name: "test 2",
			b:    NewBits(1),
			flag: NewBits(1),
			want: NewBits(1),
		}, {
			name: "test 3",
			b:    NewBits(1),
			flag: NewEmptyBits(),
			want: NewBits(1),
		}, {
			name: "test 4",
			b:    NewEmptyBits(),
			flag: NewEmptyBits(),
			want: NewEmptyBits(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Set(tt.flag)
			require.EqualValues(t, tt.want, tt.b)
		})
	}
}

func TestBits_Clear(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		flag Bits
		want Bits
	}{
		{
			name: "test 1",
			b:    NewBits(1),
			flag: NewBits(1),
			want: NewEmptyBits(),
		}, {
			name: "test 2",
			b:    NewBits(0),
			flag: NewBits(1),
			want: NewBits(0),
		}, {
			name: "test 3",
			b:    NewBits(1),
			flag: NewBits(0),
			want: NewBits(1),
		}, {
			name: "test 3",
			b:    NewBits(0),
			flag: NewBits(0),
			want: NewBits(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Clear(tt.flag)
			require.EqualValues(t, tt.want.value.Uint64(), tt.b.value.Uint64())
		})
	}
}

func TestBits_Has(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		flag Bits
		want bool
	}{
		{
			name: "test 1",
			b:    NewBits(0),
			flag: NewBits(0),
			want: false,
		}, {
			name: "test 2",
			b:    NewBits(1),
			flag: NewBits(0),
			want: false,
		}, {
			name: "test 3",
			b:    NewBits(0),
			flag: NewBits(1),
			want: false,
		}, {
			name: "test 4",
			b:    NewBits(1),
			flag: NewBits(1),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Has(tt.flag)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestBits_CountBits(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want int
	}{
		{
			name: "test 1",
			b:    NewBits(0),
			want: 0,
		}, {
			name: "test 2",
			b:    NewBits(1),
			want: 1,
		}, {
			name: "test 3",
			b:    NewBits(2),
			want: 1,
		}, {
			name: "test 3",
			b:    NewBits(3),
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := tt.b
			got := tt.b.CountBits()
			require.Equal(t, tt.want, got)
			require.Equal(t, val, tt.b)
		})
	}
}

func TestBits_Empty(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want bool
	}{
		{
			name: "test 1",
			b:    NewBits(0),
			want: true,
		}, {
			name: "test 2",
			b:    NewBits(1),
			want: false,
		}, {
			name: "test 3",
			b:    Bits{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.b.Empty())
		})
	}
}

func TestBits_SetBit(t *testing.T) {
	tests := []struct {
		name   string
		value  Bits
		number int
		want   Bits
	}{
		{
			name:   "test 1",
			value:  NewBits(1),
			number: 2,
			want:   NewBits(5),
		}, {
			name:   "test 2",
			value:  NewBits(1),
			number: 1,
			want:   NewBits(3),
		}, {
			name:   "test 1",
			value:  NewBits(1),
			number: 0,
			want:   NewBits(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.value.SetBit(tt.number)
			require.EqualValues(t, tt.want, tt.value)
		})
	}
}
