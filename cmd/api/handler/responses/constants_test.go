// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_roundCounstant(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want string
	}{
		{
			name: "test 1",
			val:  "0.000000000000000000",
			want: "0",
		}, {
			name: "test 2",
			val:  "0.334000000000000000",
			want: "0.334",
		}, {
			name: "test 3",
			val:  "utia",
			want: "utia",
		}, {
			name: "test 4",
			val:  "604800s",
			want: "604800s",
		}, {
			name: "test 5",
			val:  "10000",
			want: "10000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roundCounstant(tt.val)
			require.Equal(t, tt.want, got)
		})
	}
}
