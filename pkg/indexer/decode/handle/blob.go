// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"encoding/base64"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/celestia-app/pkg/namespace"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
	"github.com/pkg/errors"
)

// MsgPayForBlobs pays for the inclusion of a blob in the block.
func MsgPayForBlobs(level types.Level, blockTime time.Time, m *appBlobTypes.MsgPayForBlobs) (storageTypes.MsgType, []storage.AddressWithType, []storage.Namespace, []*storage.BlobLog, int64, error) {
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
			FirstHeight:     level,
			Version:         appNS.Version,
			NamespaceID:     appNS.ID,
			Size:            size,
			PfbCount:        1,
			Reserved:        appNS.IsReserved(),
			LastHeight:      level,
			LastMessageTime: blockTime,
		}
		if n, ok := uniqueNs[namespace.String()]; ok {
			n.Size += size
		} else {
			uniqueNs[namespace.String()] = &namespace
		}

		blobLog := &storage.BlobLog{
			Commitment: base64.StdEncoding.EncodeToString(m.ShareCommitments[nsI]),
			Size:       size,
			Namespace:  &namespace,
			Height:     level,
			Time:       blockTime,
		}

		blobLogs = append(blobLogs, blobLog)
	}

	namespaces := make([]storage.Namespace, 0, len(uniqueNs))
	for _, namespace := range uniqueNs {
		namespaces = append(namespaces, *namespace)
	}

	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSigner, address: m.Signer},
	}, level)

	for i := range blobLogs {
		blobLogs[i].Signer = &storage.Address{
			Address: m.Signer,
		}
	}

	return storageTypes.MsgPayForBlobs, addresses, namespaces, blobLogs, blobsSize, err
}
