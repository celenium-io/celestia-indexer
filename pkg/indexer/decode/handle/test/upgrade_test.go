// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeMsg_SuccessOnMsgSoftwareUpgrade(t *testing.T) {
	msg := &upgrade.MsgSoftwareUpgrade{
		Authority: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(msg, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeAuthority,
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
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSoftwareUpgrade,
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

func TestDecodeMsg_SuccessOnMsgCancelUpgrade(t *testing.T) {
	msg := &upgrade.MsgCancelUpgrade{
		Authority: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(msg, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeAuthority,
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
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCancelUpgrade,
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
