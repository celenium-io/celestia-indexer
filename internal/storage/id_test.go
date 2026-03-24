// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/require"
)

func Test_idFromHeightAndPosition(t *testing.T) {
	tests := []struct {
		name     string
		height   types.Level
		position int64
		want     uint64
		wantErr  bool
	}{
		{
			name:     "test 1",
			height:   0,
			position: 0,
			want:     1,
		}, {
			name:     "test 2",
			height:   1,
			position: 1,
			want:     0x0000000001000001,
		}, {
			name:     "test 3",
			height:   16,
			position: 16,
			want:     0x0000000010000010,
		}, {
			// position+1 = maxPosition — the last 24-bit value that fits without overflow
			name:     "genesis max valid position",
			height:   0,
			position: maxPosition - 1,
			want:     maxPosition,
		}, {
			// highest position allowed in a regular (non-genesis) block
			name:     "non-genesis max position",
			height:   1,
			position: maxPosition,
			want:     0x0000000001ffffff,
		}, {
			// both fields at their upper bound: 0xffffffffff_ffffff
			name:     "max height and max position",
			height:   maxHeight,
			position: maxPosition,
			want:     0xffffffffffffffff,
		}, {
			// // position+1 = maxPosition+3 > maxPosition → error
			name:     "genesis position overflow",
			height:   0,
			position: maxPosition + 2,
			wantErr:  true,
		}, {
			// position = maxPosition+1 > maxPosition → error in the height>0 branch
			name:     "non-genesis position overflow",
			height:   1,
			position: maxPosition + 1,
			wantErr:  true,
		}, {
			name:     "height overflow",
			height:   maxHeight + 1,
			position: 0,
			wantErr:  true,
		},
		// position+1 = maxPosition+1 > maxPosition → error
		{
			name:     "genesis position overflow at boundary",
			height:   0,
			position: maxPosition,
			wantErr:  true,
		},
		{
			name:     "genesis position overflow at boundary+1",
			height:   0,
			position: maxPosition + 1,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := idFromHeightAndPosition(tt.height, tt.position)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, tt.want, id)
			}
		})
	}
}
