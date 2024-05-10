// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package blob

import (
	"testing"

	"github.com/stretchr/testify/require"
	blobTypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestBlob_String(t *testing.T) {
	tests := []struct {
		name       string
		blob       *blobTypes.Blob
		commitment []byte
		height     uint64
		want       string
	}{
		{
			name: "test 1",
			blob: &blobTypes.Blob{
				Data:        []byte{0x01},
				NamespaceId: []byte{0x1},
			},
			commitment: []byte{0x02},
			height:     100,
			want:       "AQ==/100/Ag==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blob := Blob{
				Blob:       tt.blob,
				Commitment: tt.commitment,
				Height:     tt.height,
			}
			require.Equal(t, tt.want, blob.String())
		})
	}
}
