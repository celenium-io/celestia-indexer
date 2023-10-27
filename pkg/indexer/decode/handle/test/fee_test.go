package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	fee "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeMsg_SuccessOnMsgRegisterPayee(t *testing.T) {
	msg := &fee.MsgRegisterPayee{
		Relayer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		Payee:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(msg, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeRelayer,
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
			Type: storageTypes.MsgAddressTypePayee,
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
		Type:      storageTypes.MsgRegisterPayee,
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

func TestDecodeMsg_SuccessOnMsgRegisterCounterpartyPayee(t *testing.T) {
	msg := &fee.MsgRegisterCounterpartyPayee{
		Relayer:           "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		CounterpartyPayee: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(msg, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeRelayer,
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
			Type: storageTypes.MsgAddressTypePayee,
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
		Type:      storageTypes.MsgRegisterCounterpartyPayee,
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

func TestDecodeMsg_SuccessOnMsgPayPacketFee(t *testing.T) {
	msg := &fee.MsgPayPacketFee{
		Signer: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(msg, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

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
		Type:      storageTypes.MsgPayPacketFee,
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

func TestDecodeMsg_SuccessOnMsgPayPacketFeeAsync(t *testing.T) {
	msg := &fee.MsgPayPacketFeeAsync{}
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(msg, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgPayPacketFeeAsync,
		TxId:      0,
		Data:      structs.Map(msg),
		Namespace: nil,
		Addresses: []storage.AddressWithType{},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, []storage.AddressWithType(nil), dm.Addresses)
}
