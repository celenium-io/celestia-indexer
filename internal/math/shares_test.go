package math_test

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/math"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestShares(t *testing.T) {
	tests := []struct {
		name  string
		stake decimal.Decimal
		want  decimal.Decimal
	}{
		{
			name:  "test 1: one",
			stake: decimal.RequireFromString("1000000"),
			want:  decimal.RequireFromString("1"),
		}, {
			name:  "test 2: zero",
			stake: decimal.RequireFromString("100000"),
			want:  decimal.RequireFromString("0"),
		}, {
			name:  "test 3: one thousand",
			stake: decimal.RequireFromString("1000999999"),
			want:  decimal.RequireFromString("1000"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := math.Shares(tt.stake)
			require.Equal(t, tt.want.String(), got.String())
		})
	}
}
