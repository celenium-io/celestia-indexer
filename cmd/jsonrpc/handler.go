// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/base64"
	"errors"
	"slices"
	"sync"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/node"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
)

type BlobHandler struct {
	receiver   node.DalApi
	logs       storage.IBlobLog
	namespaces storage.INamespace
}

func NewBlobHandler(
	receiver node.DalApi,
	logs storage.IBlobLog,
	namespaces storage.INamespace,
) *BlobHandler {
	return &BlobHandler{
		receiver:   receiver,
		logs:       logs,
		namespaces: namespaces,
	}
}

func (h *BlobHandler) Get(ctx context.Context, height types.Level, namespace string, commitment string) (nodeTypes.Blob, error) {
	return h.receiver.Blob(ctx, height, namespace, commitment)
}

func (h *BlobHandler) GetAll(ctx context.Context, height types.Level, namespaces []string) ([]nodeTypes.Blob, error) {
	var (
		resultBlobs = make([][]nodeTypes.Blob, len(namespaces))
		resultErr   = make([]error, len(namespaces))
		wg          = new(sync.WaitGroup)
	)
	for i, namespace := range namespaces {
		wg.Add(1)
		go func(i int, namespace string) {
			defer wg.Done()

			hash, err := base64.StdEncoding.DecodeString(namespace)
			if err != nil {
				resultErr[i] = err
				return
			}
			ns, err := h.namespaces.ByNamespaceIdAndVersion(ctx, hash[1:], hash[0])
			if err != nil {
				resultErr[i] = err
				return
			}

			blobs, err := h.logs.ByNamespace(ctx, ns.Id, storage.BlobLogFilters{
				Limit:  100,
				Height: height,
			})
			if err != nil {
				resultErr[i] = err
				return
			}
			if len(blobs) > 0 {
				res := make([]nodeTypes.Blob, 0)
				for i := range blobs {
					b, err := h.receiver.Blob(ctx, height, namespace, blobs[i].Commitment)
					if err != nil {
						resultErr[i] = err
						return
					}
					res = append(res, b)
				}
				resultBlobs[i] = res
			}
		}(i, namespace)
	}
	wg.Wait()

	blobs := slices.Concat(resultBlobs...)
	err := errors.Join(resultErr...)
	return blobs, err
}
