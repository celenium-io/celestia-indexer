// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_decodeName(t *testing.T) {
	tests := []struct {
		name string
		nsId string
		want string
	}{
		{
			name: "test 1",
			nsId: "6d656d6573",
			want: "memes",
		}, {
			name: "test 2",
			nsId: "0000000000000000000000000000006d656d6573",
			want: "memes",
		}, {
			name: "test 3",
			nsId: "0000000000000000000000000000000000000000e6edd3ffbef8c7d8",
			want: "e6edd3ffbef8c7d8",
		}, {
			name: "test 4",
			nsId: "00000000000000000000000000000000000000e6edd3ffbef8c700d8",
			want: "e6edd3ffbef8c700d8",
		}, {
			name: "test 5",
			nsId: "00000000000000000000000000000000006d656d657300",
			want: "6d656d657300",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := hex.DecodeString(tt.nsId)
			require.NoError(t, err)

			got := decodeName(decoded)
			require.Equal(t, tt.want, got)
		})
	}
}
