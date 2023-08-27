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
			b:    0,
			flag: 1,
			want: 1,
		}, {
			name: "test 2",
			b:    1,
			flag: 1,
			want: 1,
		}, {
			name: "test 3",
			b:    1,
			flag: 0,
			want: 1,
		}, {
			name: "test 4",
			b:    0,
			flag: 0,
			want: 0,
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
			b:    1,
			flag: 1,
			want: 0,
		}, {
			name: "test 2",
			b:    0,
			flag: 1,
			want: 0,
		}, {
			name: "test 3",
			b:    1,
			flag: 0,
			want: 1,
		}, {
			name: "test 3",
			b:    0,
			flag: 0,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Clear(tt.flag)
			require.EqualValues(t, tt.want, tt.b)
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
			b:    0,
			flag: 0,
			want: false,
		}, {
			name: "test 2",
			b:    1,
			flag: 0,
			want: false,
		}, {
			name: "test 3",
			b:    0,
			flag: 1,
			want: false,
		}, {
			name: "test 4",
			b:    1,
			flag: 1,
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
			b:    0,
			want: 0,
		}, {
			name: "test 2",
			b:    1,
			want: 1,
		}, {
			name: "test 3",
			b:    2,
			want: 1,
		}, {
			name: "test 3",
			b:    3,
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
			b:    0,
			want: true,
		}, {
			name: "test 2",
			b:    1,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.b.Empty())
		})
	}
}
