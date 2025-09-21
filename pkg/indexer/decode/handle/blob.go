// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"encoding/base64"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	appBlobTypes "github.com/celestiaorg/celestia-app/v6/x/blob/types"
	nsPackage "github.com/celestiaorg/go-square/v3/share"
	"github.com/pkg/errors"
)

// MsgPayForBlobs pays for the inclusion of a blob in the block.
func MsgPayForBlobs(ctx *context.Context, status storageTypes.Status, m *appBlobTypes.MsgPayForBlobs) (storageTypes.MsgType, []storage.AddressWithType, []storage.Namespace, []*storage.BlobLog, int64, error) {
	var blobsSize int64
	uniqueNs := make(map[string]*storage.Namespace)
	blobLogs := make([]*storage.BlobLog, 0)

	for idx, ns := range m.Namespaces {
		if len(m.BlobSizes) < idx {
			return storageTypes.MsgUnknown, nil, nil, nil, 0, errors.Errorf(
				"blob sizes length=%d is less then namespaces index=%d", len(m.BlobSizes), idx)
		}
		if len(m.ShareCommitments) < idx {
			return storageTypes.MsgUnknown, nil, nil, nil, 0, errors.Errorf(
				"share commitment length=%d is less then namespaces index=%d", len(m.ShareCommitments), idx)
		}
		if len(m.ShareVersions) < idx {
			return storageTypes.MsgUnknown, nil, nil, nil, 0, errors.Errorf(
				"share versions length=%d is less then namespaces index=%d", len(m.ShareVersions), idx)
		}

		appNS, err := nsPackage.NewNamespaceFromBytes(ns)
		if err != nil {
			return storageTypes.MsgUnknown, nil, nil, nil, 0, errors.Wrap(err, "NewNamespaceFromBytes")
		}
		size := int64(m.BlobSizes[idx])
		blobsSize += size
		namespace := storage.Namespace{
			FirstHeight:     ctx.Block.Height,
			Version:         appNS.Version(),
			NamespaceID:     appNS.ID(),
			PfbCount:        1,
			Reserved:        appNS.IsReserved(),
			LastHeight:      ctx.Block.Height,
			LastMessageTime: ctx.Block.Time,
		}

		if status == storageTypes.StatusSuccess {
			namespace.BlobsCount = 1
			namespace.Size = size

			blobLog := &storage.BlobLog{
				Commitment: base64.StdEncoding.EncodeToString(m.ShareCommitments[idx]),
				Size:       size,
				Namespace:  &namespace,
				Height:     ctx.Block.Height,
				Time:       ctx.Block.Time,
				Signer: &storage.Address{
					Address: m.Signer,
				},
				ShareVersion: int(m.ShareVersions[idx]),
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

	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height)

	return storageTypes.MsgPayForBlobs, addresses, namespaces, blobLogs, blobsSize, err
}

// MsgUpdateBlobParams defines the sdk.Msg type to update the client parameters.
func MsgUpdateBlobParams(ctx *context.Context, status storageTypes.Status, m *appBlobTypes.MsgUpdateBlobParams) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateBlobParams
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height)
	return msgType, addresses, err
}
