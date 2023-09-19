package storage

import (
	"testing"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/stretchr/testify/require"
)

func Test_setNamespacesFromMessage(t *testing.T) {
	tests := []struct {
		name       string
		msg        storage.Message
		namespaces map[string]*storage.Namespace
	}{
		{
			name: "test 1",
			msg: storage.Message{
				Namespace: []storage.Namespace{
					{
						FirstHeight: 100,
						Version:     0,
						NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
						Size:        10,
						PfbCount:    1,
						Reserved:    false,
					},
				},
			},
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
		}, {
			name: "test 2",
			msg: storage.Message{
				Namespace: []storage.Namespace{
					{
						FirstHeight: 100,
						Version:     0,
						NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
						Size:        10,
						PfbCount:    1,
						Reserved:    false,
					}, {
						FirstHeight: 100,
						Version:     0,
						NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
						Size:        10,
						PfbCount:    1,
						Reserved:    false,
					},
				},
			},
			namespaces: map[string]*storage.Namespace{
				"000010203040506070809000102030405060708090001020304050607": {
					FirstHeight: 100,
					Version:     0,
					NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
					Size:        20,
					PfbCount:    2,
					Reserved:    false,
				},
			},
		}, {
			name: "test 3",
			msg: storage.Message{
				Namespace: []storage.Namespace{
					{
						FirstHeight: 100,
						Version:     0,
						NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
						Size:        10,
						PfbCount:    1,
						Reserved:    false,
					}, {
						FirstHeight: 100,
						Version:     0,
						NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
						Size:        10,
						PfbCount:    1,
						Reserved:    false,
					}, {
						FirstHeight: 100,
						Version:     0,
						NamespaceID: []byte{1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
						Size:        10,
						PfbCount:    1,
						Reserved:    false,
					},
				},
			},
			namespaces: map[string]*storage.Namespace{
				"000010203040506070809000102030405060708090001020304050607": {
					FirstHeight: 100,
					Version:     0,
					NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
					Size:        20,
					PfbCount:    2,
					Reserved:    false,
				},
				"001010203040506070809000102030405060708090001020304050607": {
					FirstHeight: 100,
					Version:     0,
					NamespaceID: []byte{1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
					Size:        10,
					PfbCount:    1,
					Reserved:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namespaces := make(map[string]*storage.Namespace)
			setNamespacesFromMessage(tt.msg, namespaces)
			require.Equal(t, tt.namespaces, namespaces)
		})
	}
}
