// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decoder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDecimalFromMap(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]any
		key  string
		want string
	}{
		{
			name: "test 1",
			m: map[string]any{
				"amount": "123utia",
			},
			key:  "amount",
			want: "123",
		}, {
			name: "test 2",
			m: map[string]any{
				"amount": "123utia",
			},
			key:  "invalid",
			want: "0",
		}, {
			name: "test 3",
			m: map[string]any{
				"amount": "123uta",
			},
			key:  "amount",
			want: "123",
		}, {
			name: "test 4",
			m: map[string]any{
				"amount": 123,
			},
			key:  "amount",
			want: "0",
		}, {
			name: "test 5",
			m: map[string]any{
				"amount": "123test",
			},
			key:  "amount",
			want: "123",
		}, {
			name: "test 6",
			m: map[string]any{
				"amount": "1-23test",
			},
			key:  "amount",
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecimalFromMap(tt.m, tt.key)
			require.Equal(t, tt.want, got.String())
		})
	}
}

func TestUnixNanoFromMap(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]any
		key  string
		want time.Time
	}{
		{
			name: "test 1",
			m: map[string]any{
				"packet_timeout_timestamp": "9439823803807825920",
			},
			key:  "packet_timeout_timestamp",
			want: time.Date(2269, 02, 19, 05, 16, 43, 807825920, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnixNanoFromMap(tt.m, tt.key)
			require.Equal(t, tt.want, got)
		})
	}
}
