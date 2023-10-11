// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"

	testsuite "github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/stretchr/testify/require"
)

func TestNamespace_Hash(t *testing.T) {
	tests := []struct {
		name string
		ns   Namespace
		want string
	}{
		{
			name: "test 1",
			ns: Namespace{
				NamespaceID: testsuite.MustHexDecode("0000000000000000000000000000000000006e5bbbdfcea081b366b1"),
				Version:     0,
			},
			want: "AAAAAAAAAAAAAAAAAAAAAAAAAG5bu9/OoIGzZrE=",
		}, {
			name: "test 2",
			ns: Namespace{
				NamespaceID: testsuite.MustHexDecode("000000000000000000000000000000000000a476c00deb8796b16999"),
				Version:     0,
			},
			want: "AAAAAAAAAAAAAAAAAAAAAAAAAKR2wA3rh5axaZk=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ns.Hash()
			require.Equal(t, tt.want, got)
		})
	}
}
