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
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
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
	blob, now := testsuite.EmptyBlock()
	position := 4
	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeGranter,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
				Hash:       []byte{0x38, 0xf5, 0xc9, 0x8, 0x56, 0x46, 0xad, 0xc2, 0xc0, 0x71, 0x2c, 0xcc, 0x4a, 0x9e, 0xbe, 0x5, 0x41, 0x9e, 0xd2, 0xc8},
				Balance:    storage.EmptyBalance(),
			},
		},
		{
			Type: storageTypes.MsgAddressTypeGrantee,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
				Hash:       []byte{0x64, 0xd3, 0xfc, 0x6a, 0x2a, 0x52, 0x4e, 0x2f, 0x60, 0x3f, 0x51, 0xc7, 0xee, 0x4e, 0x8d, 0x35, 0xf7, 0x23, 0x22, 0xf8},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgGrant,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      146,
		Namespace: nil,
		Addresses: addressesExpected,
		Grants: []storage.Grant{
			{
				Height: blob.Height,
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
				Time: blob.Block.Time,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
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
	blob, now := testsuite.EmptyBlock()
	position := 4

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeGrantee,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
				Hash:       []byte{0x64, 0xd3, 0xfc, 0x6a, 0x2a, 0x52, 0x4e, 0x2f, 0x60, 0x3f, 0x51, 0xc7, 0xee, 0x4e, 0x8d, 0x35, 0xf7, 0x23, 0x22, 0xf8},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:           0,
		Height:       blob.Height,
		Time:         now,
		Position:     4,
		Type:         storageTypes.MsgExec,
		TxId:         0,
		Data:         structs.Map(m),
		Namespace:    nil,
		Size:         49,
		Addresses:    addressesExpected,
		InternalMsgs: make([]string, 0),
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
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
	blob, now := testsuite.EmptyBlock()
	position := 4

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeGranter,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
				Hash:       []byte{0x38, 0xf5, 0xc9, 0x8, 0x56, 0x46, 0xad, 0xc2, 0xc0, 0x71, 0x2c, 0xcc, 0x4a, 0x9e, 0xbe, 0x5, 0x41, 0x9e, 0xd2, 0xc8},
				Balance:    storage.EmptyBalance(),
			},
		},
		{
			Type: storageTypes.MsgAddressTypeGrantee,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
				Hash:       []byte{0x64, 0xd3, 0xfc, 0x6a, 0x2a, 0x52, 0x4e, 0x2f, 0x60, 0x3f, 0x51, 0xc7, 0xee, 0x4e, 0x8d, 0x35, 0xf7, 0x23, 0x22, 0xf8},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgRevoke,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      108,
		Namespace: nil,
		Addresses: addressesExpected,
		Grants: []storage.Grant{
			{
				RevokeHeight: &blob.Height,
				Granter: &storage.Address{
					Address: "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
				},
				Grantee: &storage.Address{
					Address: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
				},
				Authorization: "msg_type",
				Revoked:       true,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
