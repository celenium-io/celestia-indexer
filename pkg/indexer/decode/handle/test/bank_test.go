// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/decode"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MsgSend

func createMsgSend() types.Msg {
	m := cosmosBankTypes.MsgSend{
		FromAddress: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Amount: types.Coins{
			types.Coin{
				Denom:  "utia",
				Amount: math.NewInt(1000),
			},
		},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgSend(t *testing.T) {
	msgSend := createMsgSend()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(msgSend, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeFromAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeToAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
				Hash:       []byte{0x64, 0x18, 0x63, 0xcf, 0xc9, 0x3b, 0x6f, 0x4e, 0x37, 0x30, 0x0, 0x2e, 0x81, 0xcb, 0x8b, 0x27, 0x27, 0xe1, 0x0, 0x68},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSend,
		TxId:      0,
		Data:      structs.Map(msgSend),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgMultiSend

func createMsgMultiSend() types.Msg {
	m := cosmosBankTypes.MsgMultiSend{
		Inputs: []cosmosBankTypes.Input{
			{Address: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6", Coins: make(types.Coins, 0)},
			{Address: "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts", Coins: make(types.Coins, 0)},
		},
		Outputs: []cosmosBankTypes.Output{
			{Address: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l", Coins: make(types.Coins, 0)},
		},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgMultiSend(t *testing.T) {
	msgMultiSend := createMsgMultiSend()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(msgMultiSend, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeInput,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:       []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeInput,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
				Hash:       []byte{8, 204, 180, 93, 112, 144, 218, 230, 174, 203, 58, 172, 76, 199, 190, 39, 45, 188, 116, 154},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeOutput,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
				Hash:       []byte{0x64, 0x18, 0x63, 0xcf, 0xc9, 0x3b, 0x6f, 0x4e, 0x37, 0x30, 0x0, 0x2e, 0x81, 0xcb, 0x8b, 0x27, 0x27, 0xe1, 0x0, 0x68},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgMultiSend,
		TxId:      0,
		Data:      structs.Map(msgMultiSend),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
