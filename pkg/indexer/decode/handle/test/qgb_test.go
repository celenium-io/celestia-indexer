package handle_test

import (
	qgbTypes "github.com/celestiaorg/celestia-app/x/qgb/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/decode"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MsgRegisterEvmAddress

func createMsgRegisterEvmAddress() types.Msg {
	m := qgbTypes.MsgRegisterEVMAddress{
		ValidatorAddress: "celestiavaloper1f5crra7r5m9kd6saw077u76x0n7dyjkkzk0qup",
		EvmAddress:       "0xfDC46fBDd8AF50d9Bf7536Bf44ce8560E423352c",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgRegisterEvmAddress(t *testing.T) {
	m := createMsgRegisterEvmAddress()
	blob, now := testsuite.EmptyBlock()
	position := 4

	dm, err := decode.Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1f5crra7r5m9kd6saw077u76x0n7dyjkkzk0qup",
				Hash:       []byte{0x4d, 0x30, 0x31, 0xf7, 0xc3, 0xa6, 0xcb, 0x66, 0xea, 0x1d, 0x73, 0xfd, 0xee, 0x7b, 0x46, 0x7c, 0xfc, 0xd2, 0x4a, 0xd6},
				Balance: storage.Balance{
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgRegisterEVMAddress,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
