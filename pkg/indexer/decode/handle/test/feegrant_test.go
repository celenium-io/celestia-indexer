// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"encoding/hex"
	"testing"

	"cosmossdk.io/x/feegrant"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"
)

// MsgGrantAllowance

func createMsgGrantAllowance() types.Msg {
	val, _ := hex.DecodeString("0a330a272f636f736d6f732e6665656772616e742e763162657461312e4261736963416c6c6f77616e63651208120608a7a5858a12121d2f636f736d6f732e617574687a2e763162657461312e4d736745786563")
	m := feegrant.MsgGrantAllowance{
		Granter: "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		Grantee: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
		Allowance: &codecTypes.Any{
			TypeUrl: "/cosmos.feegrant.v1beta1.AllowedMsgAllowance",
			Value:   val,
		},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgGrantAllowance(t *testing.T) {
	m := createMsgGrantAllowance()
	block, now := testsuite.EmptyBlock()
	position := 4

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgGrantAllowance,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      233,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

// MsgRevokeAllowance

func createMsgRevokeAllowance() types.Msg {
	m := feegrant.MsgRevokeAllowance{
		Granter: "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		Grantee: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgRevokeAllowance(t *testing.T) {
	m := createMsgRevokeAllowance()
	block, now := testsuite.EmptyBlock()
	position := 4

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgRevokeAllowance,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      98,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
