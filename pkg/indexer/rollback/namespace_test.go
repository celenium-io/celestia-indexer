// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rollback

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newNamespaceSize(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		want    namespaceSize
		wantErr bool
	}{
		{
			name: "test 1",
			data: map[string]any{
				"namespaces":        []any{"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAACt6d6t6d4="},
				"blob_sizes":        []any{12},
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			want: namespaceSize{
				"0000000000000000000000000000000000000000000000ade9deade9de": 12,
			},
			wantErr: false,
		}, {
			name: "test 2",
			data: map[string]any{
				"namespaces":        []any{},
				"blob_sizes":        []any{12},
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		}, {
			name: "test 3",
			data: map[string]any{
				"namespaces":        []any{"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAACt6d6t6d4="},
				"blob_sizes":        []any{},
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		}, {
			name: "test 4",
			data: map[string]any{
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		}, {
			name: "test 5",
			data: map[string]any{
				"namespaces":        []any{12},
				"blob_sizes":        []any{"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAACt6d6t6d4="},
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		}, {
			name: "test 6",
			data: map[string]any{
				"blob_sizes":        []any{12},
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		}, {
			name: "test 7",
			data: map[string]any{
				"namespaces":        []any{"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAACt6d6t6d4="},
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		}, {
			name: "test 8",
			data: map[string]any{
				"namespaces":        []any{"invalid string"},
				"blob_sizes":        []any{12},
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		}, {
			name: "test 9",
			data: map[string]any{
				"namespaces":        []any{"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAACt6d6t6d4="},
				"blob_sizes":        []any{"12"},
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		}, {
			name: "test 10",
			data: map[string]any{
				"namespaces":        "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAACt6d6t6d4=",
				"blob_sizes":        []any{"12"},
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		}, {
			name: "test 11",
			data: map[string]any{
				"namespaces":        []any{"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAACt6d6t6d4="},
				"blob_sizes":        12,
				"share_commitments": []any{"0CsLX630cjij9DR6nqoWfQcCH2pCQSoSuq63dTkd4Bw="},
				"share_versions":    []any{0},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newNamespaceSize(tt.data)
			require.Equal(t, tt.wantErr, err != nil, "want error")
			require.Equal(t, tt.want, got, "want namespace sizes")
		})
	}
}
