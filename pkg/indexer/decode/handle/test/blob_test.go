// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	appBlobTypes "github.com/celestiaorg/celestia-app/v4/x/blob/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

// MsgPayForBlob

func createMsgPayForBlob() types.Msg {
	msgPayForBlob := appBlobTypes.MsgPayForBlobs{
		Signer:           "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
		Namespaces:       [][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22}},
		BlobSizes:        []uint32{1},
		ShareCommitments: [][]byte{{176, 28, 134, 119, 32, 117, 87, 107, 231, 67, 121, 255, 209, 106, 52, 99, 88, 183, 85, 36, 67, 137, 98, 199, 144, 159, 13, 178, 111, 190, 121, 36}},
		ShareVersions:    []uint32{0},
	}
	return &msgPayForBlob
}

func TestDecodeMsg_SuccessOnPayForBlob(t *testing.T) {
	msgPayForBlob := createMsgPayForBlob()
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgPayForBlob, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
				Hash:       []byte{0x16, 0x53, 0x23, 0x70, 0x15, 0x89, 0xb7, 0x20, 0x14, 0xd5, 0xbd, 0xdc, 0xa8, 0xba, 0xcc, 0x60, 0xb5, 0x5, 0xd3, 0x97},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:       0,
		Height:   blob.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgPayForBlobs,
		TxId:     0,
		Data:     structs.Map(msgPayForBlob),
		Size:     120,
		Namespace: []storage.Namespace{
			{
				Id:              0,
				FirstHeight:     blob.Height,
				Version:         0,
				NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
				Size:            1,
				PfbCount:        1,
				Reserved:        false,
				LastHeight:      blob.Height,
				LastMessageTime: blob.Block.Time,
				BlobsCount:      1,
			},
		},
		Addresses: addressesExpected,
		BlobLogs: []*storage.BlobLog{
			{
				Height:     blob.Height,
				Time:       now,
				Size:       1,
				Commitment: "sByGdyB1V2vnQ3n/0Wo0Y1i3VSRDiWLHkJ8Nsm++eSQ=",
				Namespace: &storage.Namespace{
					Id:              0,
					FirstHeight:     blob.Height,
					Version:         0,
					NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
					Size:            1,
					PfbCount:        1,
					Reserved:        false,
					LastHeight:      blob.Height,
					LastMessageTime: blob.Block.Time,
					BlobsCount:      1,
				},
				Signer: &storage.Address{
					Address: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(1), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func createMultipleMsgPayForBlob() types.Msg {
	msgPayForBlob := appBlobTypes.MsgPayForBlobs{
		Signer: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
		Namespaces: [][]byte{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
		},
		BlobSizes: []uint32{1, 2, 3},
		ShareCommitments: [][]byte{
			{176, 28, 134, 119, 32, 117, 87, 107, 231, 67, 121, 255, 209, 106, 52, 99, 88, 183, 85, 36, 67, 137, 98, 199, 144, 159, 13, 178, 111, 190, 121, 36},
			{176, 28, 134, 119, 32, 117, 87, 107, 231, 67, 121, 255, 209, 106, 52, 99, 88, 183, 85, 36, 67, 137, 98, 199, 144, 159, 13, 178, 111, 190, 121, 37},
			{176, 28, 134, 119, 32, 117, 87, 107, 231, 67, 121, 255, 209, 106, 52, 99, 88, 183, 85, 36, 67, 137, 98, 199, 144, 159, 13, 178, 111, 190, 121, 38},
		},
		ShareVersions: []uint32{0, 0, 0},
	}
	return &msgPayForBlob
}

func TestDecodeMsg_ManyUpdatesInOnePayForBlob(t *testing.T) {
	msgPayForBlob := createMultipleMsgPayForBlob()
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgPayForBlob, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
				Hash:       []byte{0x16, 0x53, 0x23, 0x70, 0x15, 0x89, 0xb7, 0x20, 0x14, 0xd5, 0xbd, 0xdc, 0xa8, 0xba, 0xcc, 0x60, 0xb5, 0x5, 0xd3, 0x97},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:       0,
		Height:   blob.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgPayForBlobs,
		TxId:     0,
		Data:     structs.Map(msgPayForBlob),
		Size:     254,
		Namespace: []storage.Namespace{
			{
				Id:              0,
				FirstHeight:     blob.Height,
				Version:         0,
				NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
				Size:            6,
				PfbCount:        1,
				Reserved:        false,
				LastHeight:      blob.Height,
				LastMessageTime: blob.Block.Time,
				BlobsCount:      3,
			},
		},
		Addresses: addressesExpected,
		BlobLogs: []*storage.BlobLog{
			{
				Height:     blob.Height,
				Time:       now,
				Size:       1,
				Commitment: "sByGdyB1V2vnQ3n/0Wo0Y1i3VSRDiWLHkJ8Nsm++eSQ=",
				Namespace: &storage.Namespace{
					Id:              0,
					FirstHeight:     blob.Height,
					Version:         0,
					NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
					Size:            1,
					PfbCount:        1,
					Reserved:        false,
					LastHeight:      blob.Height,
					LastMessageTime: blob.Block.Time,
					BlobsCount:      1,
				},
				Signer: &storage.Address{
					Address: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
				},
			}, {
				Height:     blob.Height,
				Time:       now,
				Size:       2,
				Commitment: "sByGdyB1V2vnQ3n/0Wo0Y1i3VSRDiWLHkJ8Nsm++eSU=",
				Namespace: &storage.Namespace{
					Id:              0,
					FirstHeight:     blob.Height,
					Version:         0,
					NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
					Size:            2,
					PfbCount:        1,
					Reserved:        false,
					LastHeight:      blob.Height,
					LastMessageTime: blob.Block.Time,
					BlobsCount:      1,
				},
				Signer: &storage.Address{
					Address: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
				},
			}, {
				Height:     blob.Height,
				Time:       now,
				Size:       3,
				Commitment: "sByGdyB1V2vnQ3n/0Wo0Y1i3VSRDiWLHkJ8Nsm++eSY=",
				Namespace: &storage.Namespace{
					Id:              0,
					FirstHeight:     blob.Height,
					Version:         0,
					NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
					Size:            3,
					PfbCount:        1,
					Reserved:        false,
					LastHeight:      blob.Height,
					LastMessageTime: blob.Block.Time,
					BlobsCount:      1,
				},
				Signer: &storage.Address{
					Address: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(6), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_FailedOnPayForBlob(t *testing.T) {
	msgPayForBlob := createMsgPayForBlob()
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgPayForBlob, position, storageTypes.StatusFailed)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
				Hash:       []byte{0x16, 0x53, 0x23, 0x70, 0x15, 0x89, 0xb7, 0x20, 0x14, 0xd5, 0xbd, 0xdc, 0xa8, 0xba, 0xcc, 0x60, 0xb5, 0x5, 0xd3, 0x97},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:       0,
		Height:   blob.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgPayForBlobs,
		TxId:     0,
		Data:     structs.Map(msgPayForBlob),
		Size:     120,
		Namespace: []storage.Namespace{
			{
				Id:              0,
				FirstHeight:     blob.Height,
				Version:         0,
				NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
				Size:            0,
				PfbCount:        1,
				Reserved:        false,
				LastHeight:      blob.Height,
				LastMessageTime: blob.Block.Time,
				BlobsCount:      0,
			},
		},
		Addresses: addressesExpected,
		BlobLogs:  []*storage.BlobLog{},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(1), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
