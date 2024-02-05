// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package currency

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestStringTia(t *testing.T) {
	tests := []struct {
		name string
		val  decimal.Decimal
		want string
	}{
		{
			name: "test 1",
			val:  decimal.RequireFromString("0.123456789"),
			want: "0.123457",
		}, {
			name: "test 2",
			val:  decimal.RequireFromString("10000.123456789"),
			want: "10000.123457",
		}, {
			name: "test 3",
			val:  decimal.RequireFromString("10000"),
			want: "10000.000000",
		}, {
			name: "test 4",
			val:  decimal.RequireFromString("2"),
			want: "2.000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringTia(tt.val)
			require.Equal(t, tt.want, got)
		})
	}
}
