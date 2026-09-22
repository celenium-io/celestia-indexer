// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIbcClient_LatestRevisionCompare(t *testing.T) {
	tests := []struct {
		name   string
		client IbcClient
		second IbcClient
		want   int
	}{
		{
			name:   "equal",
			client: IbcClient{LatestRevisionNumber: 1, LatestRevisionHeight: 100},
			second: IbcClient{LatestRevisionNumber: 1, LatestRevisionHeight: 100},
			want:   0,
		}, {
			name:   "both zero",
			client: IbcClient{},
			second: IbcClient{},
			want:   0,
		}, {
			name:   "same revision, higher height",
			client: IbcClient{LatestRevisionNumber: 1, LatestRevisionHeight: 200},
			second: IbcClient{LatestRevisionNumber: 1, LatestRevisionHeight: 100},
			want:   1,
		}, {
			name:   "same revision, lower height",
			client: IbcClient{LatestRevisionNumber: 1, LatestRevisionHeight: 100},
			second: IbcClient{LatestRevisionNumber: 1, LatestRevisionHeight: 200},
			want:   -1,
		}, {
			// chain upgrade resets heights: revision dominates
			name:   "higher revision, lower height",
			client: IbcClient{LatestRevisionNumber: 2, LatestRevisionHeight: 5},
			second: IbcClient{LatestRevisionNumber: 1, LatestRevisionHeight: 1000},
			want:   1,
		}, {
			name:   "lower revision, higher height",
			client: IbcClient{LatestRevisionNumber: 1, LatestRevisionHeight: 1000},
			second: IbcClient{LatestRevisionNumber: 2, LatestRevisionHeight: 5},
			want:   -1,
		}, {
			// solomachine: revision is always 0, height is the sequence
			name:   "solomachine sequence",
			client: IbcClient{LatestRevisionHeight: 8},
			second: IbcClient{LatestRevisionHeight: 7},
			want:   1,
		}, {
			// events without heights (e.g. invalid misbehaviour) must never win
			name:   "empty vs set",
			client: IbcClient{},
			second: IbcClient{LatestRevisionNumber: 1, LatestRevisionHeight: 100},
			want:   -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.client.LatestRevisionCompare(&tt.second))
			// antisymmetry: swapping arguments flips the sign
			require.Equal(t, -tt.want, tt.second.LatestRevisionCompare(&tt.client))
		})
	}
}
