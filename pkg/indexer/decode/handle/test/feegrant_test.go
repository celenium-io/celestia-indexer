// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
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

	expiration := time.Date(2123, 10, 13, 18, 47, 3, 0, time.UTC)
	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgGrantAllowance,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      233,
		Namespace: nil,
		Addresses: addressesExpected,
		Grants: []storage.Grant{
			{
				Height:        blob.Height,
				Time:          blob.Block.Time,
				Authorization: "fee",
				Granter: &storage.Address{
					Address: "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
				},
				Grantee: &storage.Address{
					Address: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
				},
				Expiration: &expiration,
				Params: map[string]any{
					"Allowance": feegrant.BasicAllowance{
						SpendLimit: nil,
						Expiration: &expiration,
					},
					"AllowedMessages": []string{
						"/cosmos.authz.v1beta1.MsgExec",
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
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
		Type:      storageTypes.MsgRevokeAllowance,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      98,
		Namespace: nil,
		Addresses: addressesExpected,
		Grants: []storage.Grant{
			{
				RevokeHeight:  &blob.Height,
				Revoked:       true,
				Authorization: "fee",
				Granter: &storage.Address{
					Address: "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
				},
				Grantee: &storage.Address{
					Address: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
