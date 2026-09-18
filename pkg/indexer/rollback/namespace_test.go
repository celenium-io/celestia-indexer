// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// storedNamespaces are the rows rollbackNamespaces reads before applying a diff.
func storedNamespaces() map[uint64]storage.Namespace {
	return map[uint64]storage.Namespace{
		1: {
			Id: 1, Version: 0, NamespaceID: []byte{0xaa},
			PfbCount: 5, Size: 500,
			PffCount: 2, FibreSize: 524288,
			BlobsCount: 7,
		},
		2: {
			Id: 2, Version: 0, NamespaceID: []byte{0xbb},
			PfbCount: 1, Size: 100, BlobsCount: 1,
		},
	}
}

func Test_rollbackNamespaces(t *testing.T) {
	lastTime := time.Date(2026, 9, 18, 3, 10, 57, 0, time.UTC)

	tests := []struct {
		name      string
		nsMsgs    []storage.NamespaceMessage
		deletedNs []storage.Namespace
		want      map[uint64]storage.Namespace
	}{
		{
			name: "pay for blobs",
			nsMsgs: []storage.NamespaceMessage{
				{NamespaceId: 1, MsgId: 10, Size: 120, Source: types.BlobSourcePfb},
			},
			want: map[uint64]storage.Namespace{
				1: {PfbCount: 4, Size: 380, PffCount: 2, FibreSize: 524288, BlobsCount: 7},
			},
		}, {
			// A fibre blob moves pff_count, fibre_size and blobs_count, and must
			// leave the PFB counters alone.
			name: "pay for fibre",
			nsMsgs: []storage.NamespaceMessage{
				{NamespaceId: 1, MsgId: 10, Size: 262144, Source: types.BlobSourceFibre},
			},
			want: map[uint64]storage.Namespace{
				1: {PfbCount: 5, Size: 500, PffCount: 1, FibreSize: 262144, BlobsCount: 6},
			},
		}, {
			// Rows written before the source column existed read as pfb, which
			// is right: fibre did not exist before app v10.
			name: "empty source falls back to pfb",
			nsMsgs: []storage.NamespaceMessage{
				{NamespaceId: 1, MsgId: 10, Size: 120},
			},
			want: map[uint64]storage.Namespace{
				1: {PfbCount: 4, Size: 380, PffCount: 2, FibreSize: 524288, BlobsCount: 7},
			},
		}, {
			name: "several messages into one namespace accumulate",
			nsMsgs: []storage.NamespaceMessage{
				{NamespaceId: 1, MsgId: 10, Size: 120, Source: types.BlobSourcePfb},
				{NamespaceId: 1, MsgId: 11, Size: 262144, Source: types.BlobSourceFibre},
				{NamespaceId: 1, MsgId: 12, Size: 80, Source: types.BlobSourcePfb},
			},
			want: map[uint64]storage.Namespace{
				1: {PfbCount: 3, Size: 300, PffCount: 1, FibreSize: 262144, BlobsCount: 6},
			},
		}, {
			name: "several namespaces stay apart",
			nsMsgs: []storage.NamespaceMessage{
				{NamespaceId: 1, MsgId: 10, Size: 120, Source: types.BlobSourcePfb},
				{NamespaceId: 2, MsgId: 10, Size: 40, Source: types.BlobSourcePfb},
			},
			want: map[uint64]storage.Namespace{
				1: {PfbCount: 4, Size: 380, PffCount: 2, FibreSize: 524288, BlobsCount: 7},
				2: {PfbCount: 0, Size: 60, BlobsCount: 1},
			},
		}, {
			// A namespace deleted wholesale by the rollback needs no diff.
			name: "deleted namespaces are skipped",
			nsMsgs: []storage.NamespaceMessage{
				{NamespaceId: 1, MsgId: 10, Size: 120, Source: types.BlobSourcePfb},
				{NamespaceId: 2, MsgId: 10, Size: 40, Source: types.BlobSourcePfb},
			},
			deletedNs: []storage.Namespace{{Id: 2}},
			want: map[uint64]storage.Namespace{
				1: {PfbCount: 4, Size: 380, PffCount: 2, FibreSize: 524288, BlobsCount: 7},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			stored := storedNamespaces()
			tx := mock.NewMockTransaction(ctrl)
			tx.EXPECT().
				Namespace(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, id uint64) (storage.Namespace, error) {
					return stored[id], nil
				}).
				AnyTimes()
			tx.EXPECT().
				LastNamespaceMessage(gomock.Any(), gomock.Any()).
				Return(storage.NamespaceMessage{Height: 100, Time: lastTime}, nil).
				AnyTimes()

			var saved []*storage.Namespace
			tx.EXPECT().
				SaveNamespaces(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, namespaces ...*storage.Namespace) (int64, error) {
					saved = namespaces
					return int64(len(namespaces)), nil
				})

			module := new(Module)
			require.NoError(t, module.rollbackNamespaces(t.Context(), tx, tt.nsMsgs, tt.deletedNs))

			require.Len(t, saved, len(tt.want))
			for _, ns := range saved {
				want, ok := tt.want[ns.Id]
				require.Truef(t, ok, "namespace %d was not expected to change", ns.Id)

				require.EqualValuesf(t, want.PfbCount, ns.PfbCount, "ns %d pfb_count", ns.Id)
				require.EqualValuesf(t, want.Size, ns.Size, "ns %d size", ns.Id)
				require.EqualValuesf(t, want.PffCount, ns.PffCount, "ns %d pff_count", ns.Id)
				require.EqualValuesf(t, want.FibreSize, ns.FibreSize, "ns %d fibre_size", ns.Id)
				require.EqualValuesf(t, want.BlobsCount, ns.BlobsCount, "ns %d blobs_count", ns.Id)

				// The last surviving message becomes the namespace's new tip.
				require.EqualValuesf(t, 100, ns.LastHeight, "ns %d last_height", ns.Id)
				require.Equalf(t, lastTime, ns.LastMessageTime, "ns %d last_message_time", ns.Id)
			}
		})
	}

	t.Run("no namespace messages is a no-op", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// No call is expected on the transaction at all.
		tx := mock.NewMockTransaction(ctrl)
		module := new(Module)
		require.NoError(t, module.rollbackNamespaces(t.Context(), tx, nil, nil))
	})
}
