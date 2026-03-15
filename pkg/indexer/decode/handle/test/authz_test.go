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
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/stretchr/testify/require"
)

// MsgGrant

func createMsgGrant() types.Msg {
	m := authz.MsgGrant{
		Granter: "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		Grantee: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
		Grant: authz.Grant{
			Authorization: &codecTypes.Any{
				TypeUrl: "/cosmos.authz.v1beta1.GenericAuthorization",
			},
			Expiration: nil,
		},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgGrant(t *testing.T) {
	m := createMsgGrant()
	block, now := testsuite.EmptyBlock()
	position := 4
	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	grants := []*storage.Grant{
		{
			Height: block.Height,
			Granter: &storage.Address{
				Address: "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
			},
			Grantee: &storage.Address{
				Address: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
			},
			Authorization: "",
			Params: map[string]any{
				"Msg": "",
			},
			Time: block.Block.Time,
		},
	}
	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgGrant,
		TxId:      0,
		Data:      mustMsgToMap(t, m),
		Size:      146,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.EqualValues(t, decodeCtx.Grants.Len(), 1)
	_ = decodeCtx.Grants.Range(func(_ string, value *storage.Grant) (error, bool) {
		require.Equal(t, value, grants[0])
		return nil, false
	})
}

// MsgExec

func createMsgExec() types.Msg {
	m := authz.MsgExec{
		Grantee: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
		Msgs:    make([]*codecTypes.Any, 0),
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgExec(t *testing.T) {
	m := createMsgExec()
	block, now := testsuite.EmptyBlock()
	position := 4

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:           1,
		Height:       block.Height,
		Time:         now,
		Position:     4,
		Type:         storageTypes.MsgExec,
		TxId:         0,
		Data:         mustMsgToMap(t, m),
		Namespace:    nil,
		Size:         49,
		InternalMsgs: make([]string, 0),
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

// MsgRevoke

func createMsgRevoke() types.Msg {
	m := authz.MsgRevoke{
		Granter:    "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		Grantee:    "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
		MsgTypeUrl: "msg_type",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgRevoke(t *testing.T) {
	m := createMsgRevoke()
	block, now := testsuite.EmptyBlock()
	position := 4

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)

	grants := []*storage.Grant{
		{
			RevokeHeight: &block.Height,
			Granter: &storage.Address{
				Address: "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
			},
			Grantee: &storage.Address{
				Address: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
			},
			Authorization: "msg_type",
			Revoked:       true,
		},
	}
	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgRevoke,
		TxId:      0,
		Data:      mustMsgToMap(t, m),
		Size:      108,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
	require.EqualValues(t, decodeCtx.Grants.Len(), 1)
	_ = decodeCtx.Grants.Range(func(_ string, value *storage.Grant) (error, bool) {
		require.Equal(t, value, grants[0])
		return nil, false
	})
}

// MsgPruneExpiredGrants

func TestDecodeMsg_SuccessOnMsgPruneExpiredGrants(t *testing.T) {
	m := &authz.MsgPruneExpiredGrants{
		Pruner: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
	}
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
		Type:      storageTypes.MsgPruneExpiredGrants,
		TxId:      0,
		Data:      mustMsgToMap(t, m),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
