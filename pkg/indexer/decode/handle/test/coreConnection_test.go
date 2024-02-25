// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	coreConnection "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDecodeMsg_SuccessOnMsgConnectionOpenInit(t *testing.T) {
	msg := &coreConnection.MsgConnectionOpenInit{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance: storage.Balance{
					Id:        0,
					Spendable: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgConnectionOpenInit,
		TxId:      0,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgConnectionOpenTry(t *testing.T) {
	msg := &coreConnection.MsgConnectionOpenTry{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance: storage.Balance{
					Id:        0,
					Spendable: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgConnectionOpenTry,
		TxId:      0,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgConnectionOpenAck(t *testing.T) {
	msg := &coreConnection.MsgConnectionOpenAck{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance: storage.Balance{
					Id:        0,
					Spendable: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgConnectionOpenAck,
		TxId:      0,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func TestDecodeMsg_SuccessOnMsgConnectionOpenConfirm(t *testing.T) {
	msg := &coreConnection.MsgConnectionOpenConfirm{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeSigner,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance: storage.Balance{
					Id:        0,
					Spendable: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgConnectionOpenConfirm,
		TxId:      0,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
