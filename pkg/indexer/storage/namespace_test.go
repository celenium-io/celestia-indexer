// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_saveNamespaces(t *testing.T) {
	tests := []struct {
		name       string
		namespaces map[string]*storage.Namespace
		want       int64
		wantErr    bool
	}{
		{
			name: "test 1",
			namespaces: map[string]*storage.Namespace{
				"000010203040506070809000102030405060708090001020304050607": {
					FirstHeight: 100,
					Version:     0,
					NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
					Size:        10,
					PfbCount:    1,
					Reserved:    false,
				},
			},
			want: 1,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := mock.NewMockTransaction(ctrl)

			tx.EXPECT().SaveNamespaces(gomock.Any(), gomock.Any()).
				MaxTimes(1).
				MinTimes(1).
				DoAndReturn(func(_ context.Context, ns ...*storage.Namespace) (int64, error) {
					require.Equal(t, len(tt.namespaces), len(ns))
					return int64(len(ns)), nil
				})

			got, err := saveNamespaces(context.Background(), tx, tt.namespaces)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}
