// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	zkismTypes "github.com/celestiaorg/celestia-app/v7/x/zkism/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

func TestZkism_MsgCreateInterchainSecurityModule(t *testing.T) {
	msg := &zkismTypes.MsgCreateInterchainSecurityModule{
		Creator:             "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		State:               testsuite.RandomBytes(10),
		MerkleTreeAddress:   testsuite.RandomBytes(32),
		Groth16Vkey:         testsuite.RandomBytes(16),
		StateTransitionVkey: testsuite.RandomBytes(16),
		StateMembershipVkey: testsuite.RandomBytes(16),
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   now,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSender,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
				Hash:       []byte{0x64, 0x18, 0x63, 0xcf, 0xc9, 0x3b, 0x6f, 0x4e, 0x37, 0x30, 0x0, 0x2e, 0x81, 0xcb, 0x8b, 0x27, 0x27, 0xe1, 0x0, 0x68},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateInterchainSecurityModule,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      149,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestZkism_MsgUpdateInterchainSecurityModule(t *testing.T) {
	msg := &zkismTypes.MsgUpdateInterchainSecurityModule{
		Signer:       "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Proof:        testsuite.RandomBytes(32),
		PublicValues: testsuite.RandomBytes(16),
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   now,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
				Hash:       []byte{0x64, 0x18, 0x63, 0xcf, 0xc9, 0x3b, 0x6f, 0x4e, 0x37, 0x30, 0x0, 0x2e, 0x81, 0xcb, 0x8b, 0x27, 0x27, 0xe1, 0x0, 0x68},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUpdateInterchainSecurityModule,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      169,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestZkism_MsgSubmitMessages(t *testing.T) {
	msg := &zkismTypes.MsgSubmitMessages{
		Signer:       "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Proof:        testsuite.RandomBytes(32),
		PublicValues: testsuite.RandomBytes(16),
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   now,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     block.Height,
				LastHeight: block.Height,
				Address:    "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
				Hash:       []byte{0x64, 0x18, 0x63, 0xcf, 0xc9, 0x3b, 0x6f, 0x4e, 0x37, 0x30, 0x0, 0x2e, 0x81, 0xcb, 0x8b, 0x27, 0x27, 0xe1, 0x0, 0x68},
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSubmitMessages,
		TxId:      0,
		Data:      structs.Map(msg),
		Size:      169,
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
