// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"encoding/base64"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	appBlobTypes "github.com/celestiaorg/celestia-app/v7/x/blob/types"
	nsPackage "github.com/celestiaorg/go-square/v3/share"
	"github.com/pkg/errors"
)

// MsgPayForBlobs pays for the inclusion of a blob in the block.
func MsgPayForBlobs(ctx *context.Context, status storageTypes.Status, msgId, txId uint64, m *appBlobTypes.MsgPayForBlobs) (storageTypes.MsgType, []*storage.BlobLog, int64, error) {
	var blobsSize int64
	blobLogs := make([]*storage.BlobLog, 0)

	for idx, ns := range m.Namespaces {
		if len(m.BlobSizes) < idx {
			return storageTypes.MsgUnknown, nil, 0, errors.Errorf(
				"blob sizes length=%d is less then namespaces index=%d", len(m.BlobSizes), idx)
		}
		if len(m.ShareCommitments) < idx {
			return storageTypes.MsgUnknown, nil, 0, errors.Errorf(
				"share commitment length=%d is less then namespaces index=%d", len(m.ShareCommitments), idx)
		}
		if len(m.ShareVersions) < idx {
			return storageTypes.MsgUnknown, nil, 0, errors.Errorf(
				"share versions length=%d is less then namespaces index=%d", len(m.ShareVersions), idx)
		}

		appNS, err := nsPackage.NewNamespaceFromBytes(ns)
		if err != nil {
			return storageTypes.MsgUnknown, nil, 0, errors.Wrap(err, "NewNamespaceFromBytes")
		}
		size := int64(m.BlobSizes[idx])
		blobsSize += size
		namespace := &storage.Namespace{
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

			ns := ctx.AddNamespace(namespace)

			blobLog := &storage.BlobLog{
				Commitment: base64.StdEncoding.EncodeToString(m.ShareCommitments[idx]),
				Size:       size,
				Namespace:  ns,
				Height:     ctx.Block.Height,
				Time:       ctx.Block.Time,
				Signer: &storage.Address{
					Address: m.Signer,
				},
				ShareVersion: int(m.ShareVersions[idx]),
				MsgId:        msgId,
				TxId:         txId,
			}
			blobLogs = append(blobLogs, blobLog)

			ctx.AddNamespaceMessage(&storage.NamespaceMessage{
				MsgId:     msgId,
				TxId:      txId,
				Height:    ctx.Block.Height,
				Time:      ctx.Block.Time,
				Namespace: ns,
				Size:      uint64(size),
			})
		}
	}

	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, ctx.Block.Height, msgId)

	return storageTypes.MsgPayForBlobs, blobLogs, blobsSize, err
}

// MsgUpdateBlobParams defines the sdk.Msg type to update the client parameters.
func MsgUpdateBlobParams(ctx *context.Context, status storageTypes.Status, msgId uint64, m *appBlobTypes.MsgUpdateBlobParams) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateBlobParams
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
