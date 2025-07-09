// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package blob

import (
	"testing"

	blobTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"
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
			want:       "AAE=/100/Ag==",
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

func TestBase64ToUrl(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			s:    "zvwwM2k3fmfU8t6i1Mprs34+VQUIn2bdvz6IO2thcAU=",
			want: "zvwwM2k3fmfU8t6i1Mprs34-VQUIn2bdvz6IO2thcAU=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64ToUrl(tt.s)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
