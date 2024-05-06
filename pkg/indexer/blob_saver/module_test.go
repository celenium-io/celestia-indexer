// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package blobsaver

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/blob"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/proto/tendermint/types"
	"go.uber.org/mock/gomock"
)

func TestBlobSaverModule(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := blob.NewMockStorage(ctrl)

	module, err := NewModule("mock")
	require.NoError(t, err, "create module")
	module.storage = storage

	b1 := &types.Blob{
		NamespaceId:      []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10},
		Data:             []byte{0, 1, 2, 3},
		ShareVersion:     0,
		NamespaceVersion: 0,
	}
	b2 := &types.Blob{
		NamespaceId:      []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x67, 0x6d},
		Data:             []byte("0x676d"),
		ShareVersion:     0,
		NamespaceVersion: 0,
	}
	commitment, err := base64.StdEncoding.DecodeString("uwghsElFtoHNqQ3JrsDGj8uLW456izVbegVL/AunMOw=")
	require.NoError(t, err, "decode commitment")

	commitment2, err := base64.StdEncoding.DecodeString("tHdPeRQ5+xwJ7ik4Er99x8/HXyDEndFtBhRZ3I9f6/8=")
	require.NoError(t, err, "decode commitment2")

	storage.EXPECT().
		Head(ctx).
		Return(100, nil).
		Times(1)

	storage.EXPECT().
		SaveBulk(ctx, []blob.Blob{
			{
				Blob:       b1,
				Height:     101,
				Commitment: commitment,
			}, {
				Blob:       b2,
				Height:     101,
				Commitment: commitment2,
			},
		}).
		Return(nil).
		Times(1)

	storage.EXPECT().
		UpdateHead(ctx, uint64(101)).
		Return(nil).
		Times(1)

	module.Start(ctx)
	require.NoError(t, err, "init module")
	require.EqualValues(t, 100, module.head)

	input := module.MustInput(InputName)
	input.Push(&Msg{
		Height: 101,
		Blob:   b1,
	})
	input.Push(&Msg{
		Height: 101,
		Blob:   b2,
	})
	input.Push(&Msg{
		Height:   101,
		EndBlock: true,
	})

	var end bool
	for !end {
		end = module.blocks.Len() == 0
		time.Sleep(time.Millisecond)
	}

	cancel()
	err = module.Close()
	require.NoError(t, err, "closing module")
}
