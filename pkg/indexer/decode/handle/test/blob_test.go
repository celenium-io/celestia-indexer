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
	appBlobTypes "github.com/celestiaorg/celestia-app/v8/x/blob/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgPayForBlob, position, storageTypes.StatusSuccess, 0)

	blobLogs := []*storage.BlobLog{
		{
			Height:     block.Height,
			Time:       now,
			Size:       1,
			Commitment: "sByGdyB1V2vnQ3n/0Wo0Y1i3VSRDiWLHkJ8Nsm++eSQ=",
			Namespace: &storage.Namespace{
				FirstHeight:     block.Height,
				Version:         0,
				NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
				Size:            1,
				PfbCount:        1,
				Reserved:        false,
				LastHeight:      block.Height,
				LastMessageTime: block.Block.Time,
				BlobsCount:      1,
			},
			Signer: &storage.Address{
				Address: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
			},
			MsgId: 1,
		},
	}
	msgExpected := storage.Message{
		Id:       1,
		Height:   block.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgPayForBlobs,
		TxId:     0,
		Data:     mustMsgToMap(t, msgPayForBlob),
		Size:     120,
	}

	require.NoError(t, err)
	require.Equal(t, int64(1), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, blobLogs, dm.BlobLogs)
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
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgPayForBlob, position, storageTypes.StatusSuccess, 0)

	blobLogs := []*storage.BlobLog{
		{
			Height:     block.Height,
			Time:       now,
			Size:       1,
			Commitment: "sByGdyB1V2vnQ3n/0Wo0Y1i3VSRDiWLHkJ8Nsm++eSQ=",
			Namespace: &storage.Namespace{
				FirstHeight:     block.Height,
				Version:         0,
				NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
				Size:            6,
				PfbCount:        3,
				Reserved:        false,
				LastHeight:      block.Height,
				LastMessageTime: block.Block.Time,
				BlobsCount:      3,
			},
			Signer: &storage.Address{
				Address: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
			},
			MsgId: 1,
		}, {
			Height:     block.Height,
			Time:       now,
			Size:       2,
			Commitment: "sByGdyB1V2vnQ3n/0Wo0Y1i3VSRDiWLHkJ8Nsm++eSU=",
			Namespace: &storage.Namespace{
				FirstHeight:     block.Height,
				Version:         0,
				NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
				Size:            6,
				PfbCount:        3,
				Reserved:        false,
				LastHeight:      block.Height,
				LastMessageTime: block.Block.Time,
				BlobsCount:      3,
			},
			Signer: &storage.Address{
				Address: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
			},
			MsgId: 1,
		}, {
			Height:     block.Height,
			Time:       now,
			Size:       3,
			Commitment: "sByGdyB1V2vnQ3n/0Wo0Y1i3VSRDiWLHkJ8Nsm++eSY=",
			Namespace: &storage.Namespace{
				FirstHeight:     block.Height,
				Version:         0,
				NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
				Size:            6,
				PfbCount:        3,
				Reserved:        false,
				LastHeight:      block.Height,
				LastMessageTime: block.Block.Time,
				BlobsCount:      3,
			},
			Signer: &storage.Address{
				Address: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
			},
			MsgId: 1,
		},
	}
	msgExpected := storage.Message{
		Id:       1,
		Height:   block.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgPayForBlobs,
		TxId:     0,
		Data:     mustMsgToMap(t, msgPayForBlob),
		Size:     254,
	}

	require.NoError(t, err)
	require.Equal(t, int64(6), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Equal(t, blobLogs, dm.BlobLogs)
}

func TestDecodeMsg_FailedOnPayForBlob(t *testing.T) {
	msgPayForBlob := createMsgPayForBlob()
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msgPayForBlob, position, storageTypes.StatusFailed, 0)

	msgExpected := storage.Message{
		Id:       1,
		Height:   block.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgPayForBlobs,
		TxId:     0,
		Data:     mustMsgToMap(t, msgPayForBlob),
		Size:     120,
	}

	require.NoError(t, err)
	require.Equal(t, int64(1), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.Len(t, dm.BlobLogs, 0)
}
