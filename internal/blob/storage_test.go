// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package blob

import (
	"testing"

	blobTypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestBlob_String(t *testing.T) {
	type fields struct {
		Blob       *blobTypes.Blob
		Commitment []byte
		Height     uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blob := Blob{
				Blob:       tt.fields.Blob,
				Commitment: tt.fields.Commitment,
				Height:     tt.fields.Height,
			}
			if got := blob.String(); got != tt.want {
				t.Errorf("Blob.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
