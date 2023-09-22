package handle_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/decode"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MsgEditValidator

func createMsgEditValidator() types.Msg {
	m := cosmosStakingTypes.MsgEditValidator{
		Description: cosmosStakingTypes.Description{
			Moniker:         "newAgeValidator",
			Identity:        "UPort:1",
			Website:         "https://google.com",
			SecurityContact: "tryme@gmail.com",
			Details:         "trust",
		},
		ValidatorAddress:  "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		CommissionRate:    nil,
		MinSelfDelegation: nil,
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgEditValidator(t *testing.T) {
	m := createMsgEditValidator()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:       []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
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
		Type:      storageTypes.MsgEditValidator,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
		Validator: &storage.Validator{
			Address:           "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
			Moniker:           "newAgeValidator",
			Identity:          "UPort:1",
			Website:           "https://google.com",
			Contacts:          "tryme@gmail.com",
			Details:           "trust",
			Rate:              decimal.Zero,
			MinSelfDelegation: decimal.Zero,
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgBeginRedelegate

func createMsgBeginRedelegate() types.Msg {
	m := cosmosStakingTypes.MsgBeginRedelegate{
		DelegatorAddress:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorSrcAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		ValidatorDstAddress: "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgBeginRedelegate(t *testing.T) {
	m := createMsgBeginRedelegate()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
				Hash:       []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidatorSrc,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:       []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidatorDst,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
				Hash:       []byte{0x56, 0x35, 0x87, 0x35, 0xf6, 0x7, 0xd1, 0x53, 0xc1, 0xb7, 0x94, 0xe6, 0x54, 0xbc, 0x47, 0x52, 0x1f, 0x83, 0xa9, 0x3b},
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
		Type:      storageTypes.MsgBeginRedelegate,
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

// MsgCreateValidator

func createMsgCreateValidator() types.Msg {
	m := cosmosStakingTypes.MsgCreateValidator{
		Description:       cosmosStakingTypes.Description{},
		Commission:        cosmosStakingTypes.CommissionRates{},
		MinSelfDelegation: types.NewInt(1),
		DelegatorAddress:  "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorAddress:  "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		Pubkey:            nil,
		Value:             types.Coin{},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreateValidator(t *testing.T) {
	m := createMsgCreateValidator()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
				Hash:       []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:       []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
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
		Type:      storageTypes.MsgCreateValidator,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
		Validator: &storage.Validator{
			Delegator:         "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
			Address:           "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
			Rate:              decimal.Zero,
			MaxRate:           decimal.Zero,
			MaxChangeRate:     decimal.Zero,
			MinSelfDelegation: decimal.RequireFromString("1"),
			Height:            uint64(blob.Height),
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgDelegate

func createMsgDelegate() types.Msg {

	msgDelegate := cosmosStakingTypes.MsgDelegate{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
		Amount: types.Coin{
			Denom:  "utia",
			Amount: math.NewInt(1000),
		},
	}

	return &msgDelegate
}

func TestDecodeMsg_SuccessOnMsgDelegate(t *testing.T) {
	msgDelegate := createMsgDelegate()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(msgDelegate, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:       []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
				Hash:       []byte{0x56, 0x35, 0x87, 0x35, 0xf6, 0x7, 0xd1, 0x53, 0xc1, 0xb7, 0x94, 0xe6, 0x54, 0xbc, 0x47, 0x52, 0x1f, 0x83, 0xa9, 0x3b},
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
		Type:      storageTypes.MsgDelegate,
		TxId:      0,
		Data:      structs.Map(msgDelegate),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgUndelegate

func createMsgUndelegate() types.Msg {
	m := cosmosStakingTypes.MsgUndelegate{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
		Amount: types.Coin{
			Denom:  "utia",
			Amount: math.NewInt(1001),
		},
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgUndelegate(t *testing.T) {
	m := createMsgUndelegate()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:       []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
				Hash:       []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
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
		Type:      storageTypes.MsgUndelegate,
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

// MsgCancelUnbondingDelegation

func createMsgCancelUnbondingDelegation() types.Msg {
	m := cosmosStakingTypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
		Amount: types.Coin{
			Denom:  "utia",
			Amount: math.NewInt(1001),
		},
		CreationHeight: 100,
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgCancelUnbondingDelegation(t *testing.T) {
	m := createMsgCancelUnbondingDelegation()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := decode.Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:       []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidator,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
				Hash:       []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
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
		Type:      storageTypes.MsgCancelUnbondingDelegation,
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
