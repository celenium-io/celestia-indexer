// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package l2beat

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestItem_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want Item
	}{
		{
			name: "test 1",
			data: []byte(`[1735516800,7731223048.1,5923475975.24,4919978970.45,3347.85]`),
			want: Item{
				Time:      time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC),
				Native:    decimal.RequireFromString("7731223048.1"),
				Canonical: decimal.RequireFromString("5923475975.24"),
				External:  decimal.RequireFromString("4919978970.45"),
				EthPrice:  decimal.RequireFromString("3347.85"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var item Item
			err := item.UnmarshalJSON(tt.data)
			require.NoError(t, err)
			require.Equal(t, tt.want, item)
		})
	}
}
