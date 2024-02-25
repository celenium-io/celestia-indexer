// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"encoding/base64"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celestiaorg/celestia-app/pkg/namespace"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
	"github.com/pkg/errors"
)

// MsgPayForBlobs pays for the inclusion of a blob in the block.
func MsgPayForBlobs(ctx *context.Context, status storageTypes.Status, m *appBlobTypes.MsgPayForBlobs) (storageTypes.MsgType, []storage.AddressWithType, []storage.Namespace, []*storage.BlobLog, int64, error) {
	var blobsSize int64
	uniqueNs := make(map[string]*storage.Namespace)
	blobLogs := make([]*storage.BlobLog, 0)

	for nsI, ns := range m.Namespaces {
		if len(m.BlobSizes) < nsI {
			return storageTypes.MsgUnknown, nil, nil, nil, 0, errors.Errorf(
				"blob sizes length=%d is less then namespaces index=%d", len(m.BlobSizes), nsI)
		}
		if len(m.ShareCommitments) < nsI {
			return storageTypes.MsgUnknown, nil, nil, nil, 0, errors.Errorf(
				"share commitment sizes length=%d is less then namespaces index=%d", len(m.ShareCommitments), nsI)
		}

		appNS := namespace.Namespace{Version: ns[0], ID: ns[1:]}
		size := int64(m.BlobSizes[nsI])
		blobsSize += size
		namespace := storage.Namespace{
			FirstHeight:     ctx.Block.Height,
			Version:         appNS.Version,
			NamespaceID:     appNS.ID,
			PfbCount:        1,
			Reserved:        appNS.IsReserved(),
			LastHeight:      ctx.Block.Height,
			LastMessageTime: ctx.Block.Time,
		}

		if status == storageTypes.StatusSuccess {
			namespace.BlobsCount = 1
			namespace.Size = size

			blobLog := &storage.BlobLog{
				Commitment: base64.StdEncoding.EncodeToString(m.ShareCommitments[nsI]),
				Size:       size,
				Namespace:  &namespace,
				Height:     ctx.Block.Height,
				Time:       ctx.Block.Time,
				Signer: &storage.Address{
					Address: m.Signer,
				},
			}
			blobLogs = append(blobLogs, blobLog)
		}

		if n, ok := uniqueNs[namespace.String()]; ok {
			n.Size += size
			n.BlobsCount += namespace.BlobsCount
		} else {
			uniqueNs[namespace.String()] = namespace.Copy()
		}
	}

	namespaces := make([]storage.Namespace, 0, len(uniqueNs))
	for _, namespace := range uniqueNs {
		namespaces = append(namespaces, *namespace)
	}

	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)

	return storageTypes.MsgPayForBlobs, addresses, namespaces, blobLogs, blobsSize, err
}
