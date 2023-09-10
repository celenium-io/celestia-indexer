package decode

import (
	"testing"

	"cosmossdk.io/math"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
	qgbTypes "github.com/celestiaorg/celestia-app/x/qgb/types"
	cosmosCodecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	cosmosFeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// MsgWithdrawValidatorCommission

func createMsgWithdrawValidatorCommission() cosmosTypes.Msg {
	m := cosmosDistributionTypes.MsgWithdrawValidatorCommission{
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgWithdrawValidatorCommission(t *testing.T) {
	m := createMsgWithdrawValidatorCommission()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgWithdrawValidatorCommission,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:    []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgWithdrawDelegatorReward

func createMsgWithdrawDelegatorReward() cosmosTypes.Msg {
	m := cosmosDistributionTypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgWithdrawDelegatorReward(t *testing.T) {
	m := createMsgWithdrawDelegatorReward()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgWithdrawDelegatorReward,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
				Hash:    []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:    []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgEditValidator

func createMsgEditValidator() cosmosTypes.Msg {
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

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgEditValidator,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:    []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgBeginRedelegate

func createMsgBeginRedelegate() cosmosTypes.Msg {
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

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgBeginRedelegate,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
				Hash:    []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorSrcAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:    []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorDstAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
				Hash:    []byte{0x56, 0x35, 0x87, 0x35, 0xf6, 0x7, 0xd1, 0x53, 0xc1, 0xb7, 0x94, 0xe6, 0x54, 0xbc, 0x47, 0x52, 0x1f, 0x83, 0xa9, 0x3b},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgCreateValidator

func createMsgCreateValidator() cosmosTypes.Msg {
	m := cosmosStakingTypes.MsgCreateValidator{
		Description:       cosmosStakingTypes.Description{},
		Commission:        cosmosStakingTypes.CommissionRates{},
		MinSelfDelegation: cosmosTypes.Int{}, // nolint
		DelegatorAddress:  "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorAddress:  "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		Pubkey:            nil,
		Value:             cosmosTypes.Coin{},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreateValidator(t *testing.T) {
	m := createMsgCreateValidator()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateValidator,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
				Hash:    []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:    []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgDelegate

func createMsgDelegate() cosmosTypes.Msg {

	msgDelegate := cosmosStakingTypes.MsgDelegate{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
		Amount: cosmosTypes.Coin{
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

	dm, err := Message(msgDelegate, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgDelegate,
		TxId:      0,
		Data:      structs.Map(msgDelegate),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:    []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
				Hash:    []byte{0x56, 0x35, 0x87, 0x35, 0xf6, 0x7, 0xd1, 0x53, 0xc1, 0xb7, 0x94, 0xe6, 0x54, 0xbc, 0x47, 0x52, 0x1f, 0x83, 0xa9, 0x3b},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgUndelegate

func createMsgUndelegate() cosmosTypes.Msg {
	m := cosmosStakingTypes.MsgUndelegate{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
		Amount: cosmosTypes.Coin{
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

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUndelegate,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:    []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
				Hash:    []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgUnjail

func createMsgUnjail() cosmosTypes.Msg {
	m := cosmosSlashingTypes.MsgUnjail{
		ValidatorAddr: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgUnjail(t *testing.T) {
	m := createMsgUnjail()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUnjail,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
				Hash:    []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgSend

func createMsgSend() cosmosTypes.Msg {
	m := cosmosBankTypes.MsgSend{
		FromAddress: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Amount: cosmosTypes.Coins{
			cosmosTypes.Coin{
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

	dm, err := Message(msgSend, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSend,
		TxId:      0,
		Data:      structs.Map(msgSend),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeFromAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:    []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeToAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
				Hash:    []byte{0x64, 0x18, 0x63, 0xcf, 0xc9, 0x3b, 0x6f, 0x4e, 0x37, 0x30, 0x0, 0x2e, 0x81, 0xcb, 0x8b, 0x27, 0x27, 0xe1, 0x0, 0x68},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgCreateVestingAccount

func createMsgCreateVestingAccount() cosmosTypes.Msg {
	m := cosmosVestingTypes.MsgCreateVestingAccount{
		FromAddress: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Amount: cosmosTypes.Coins{
			cosmosTypes.Coin{
				Denom:  "utia",
				Amount: math.NewInt(1000),
			},
		},
		EndTime: 0,
		Delayed: false,
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreateVestingAccount(t *testing.T) {
	m := createMsgCreateVestingAccount()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateVestingAccount,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeFromAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:    []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeToAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
				Hash:    []byte{0x64, 0x18, 0x63, 0xcf, 0xc9, 0x3b, 0x6f, 0x4e, 0x37, 0x30, 0x0, 0x2e, 0x81, 0xcb, 0x8b, 0x27, 0x27, 0xe1, 0x0, 0x68},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgCreatePeriodicVestingAccount

func createMsgCreatePeriodicVestingAccount() cosmosTypes.Msg {
	m := cosmosVestingTypes.MsgCreatePeriodicVestingAccount{
		FromAddress:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		ToAddress:      "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		StartTime:      0,
		VestingPeriods: nil,
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreatePeriodicVestingAccount(t *testing.T) {
	msgCreatePeriodicVestingAccount := createMsgCreatePeriodicVestingAccount()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(msgCreatePeriodicVestingAccount, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreatePeriodicVestingAccount,
		TxId:      0,
		Data:      structs.Map(msgCreatePeriodicVestingAccount),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeFromAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
				Hash:    []byte{0x94, 0x63, 0x42, 0xc7, 0x73, 0x2f, 0x6, 0x22, 0x37, 0xb6, 0x73, 0xf5, 0x26, 0xf2, 0x67, 0x4c, 0x7d, 0xb, 0xff, 0x98},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeToAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte{0x64, 0x18, 0x63, 0xcf, 0xc9, 0x3b, 0x6f, 0x4e, 0x37, 0x30, 0x0, 0x2e, 0x81, 0xcb, 0x8b, 0x27, 0x27, 0xe1, 0x0, 0x68},
				Address: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgPayForBlob

func createMsgPayForBlob() cosmosTypes.Msg {
	msgPayForBlob := appBlobTypes.MsgPayForBlobs{
		Signer:           "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
		Namespaces:       [][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22}},
		BlobSizes:        []uint32{1},
		ShareCommitments: [][]byte{{176, 28, 134, 119, 32, 117, 87, 107, 231, 67, 121, 255, 209, 106, 52, 99, 88, 183, 85, 36, 67, 137, 98, 199, 144, 159, 13, 178, 111, 190, 121, 36}},
		ShareVersions:    []uint32{0},
	}
	return &msgPayForBlob
}

func TestDecodeMsg_SuccessOnPayForBlob(t *testing.T) {
	msgPayForBlob := createMsgPayForBlob()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(msgPayForBlob, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:       0,
		Height:   blob.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgPayForBlobs,
		TxId:     0,
		Data:     structs.Map(msgPayForBlob),
		Namespace: []storage.Namespace{
			{
				Id:          0,
				FirstHeight: blob.Height,
				Version:     0,
				NamespaceID: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
				Size:        1,
				PfbCount:    1,
				Reserved:    false,
			},
		},
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeSigner,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1zefjxuq43xmjq9x4hhw23wkvvz6st5uhv40tys",
				Hash:    []byte{0x16, 0x53, 0x23, 0x70, 0x15, 0x89, 0xb7, 0x20, 0x14, 0xd5, 0xbd, 0xdc, 0xa8, 0xba, 0xcc, 0x60, 0xb5, 0x5, 0xd3, 0x97},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgGrantAllowance

func createMsgGrantAllowance() cosmosTypes.Msg {
	m := cosmosFeegrant.MsgGrantAllowance{
		Granter:   "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		Grantee:   "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
		Allowance: cosmosCodecTypes.UnsafePackAny(cosmosFeegrant.BasicAllowance{}),
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgGrantAllowance(t *testing.T) {
	m := createMsgGrantAllowance()
	blob, now := testsuite.EmptyBlock()
	position := 4

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgGrantAllowance,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeGranter,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
				Hash:    []byte{0x38, 0xf5, 0xc9, 0x8, 0x56, 0x46, 0xad, 0xc2, 0xc0, 0x71, 0x2c, 0xcc, 0x4a, 0x9e, 0xbe, 0x5, 0x41, 0x9e, 0xd2, 0xc8},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.TxAddressTypeGrantee,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
				Hash:    []byte{0x64, 0xd3, 0xfc, 0x6a, 0x2a, 0x52, 0x4e, 0x2f, 0x60, 0x3f, 0x51, 0xc7, 0xee, 0x4e, 0x8d, 0x35, 0xf7, 0x23, 0x22, 0xf8},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgRegisterEvmAddress

func createMsgRegisterEvmAddress() cosmosTypes.Msg {
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

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgRegisterEVMAddress,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestiavaloper1f5crra7r5m9kd6saw077u76x0n7dyjkkzk0qup",
				Hash:    []byte{0x4d, 0x30, 0x31, 0xf7, 0xc3, 0xa6, 0xcb, 0x66, 0xea, 0x1d, 0x73, 0xfd, 0xee, 0x7b, 0x46, 0x7c, 0xfc, 0xd2, 0x4a, 0xd6},
				Balance: storage.Balance{
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

func createMsgSetWithdrawAddress() cosmosTypes.Msg {
	m := cosmosDistributionTypes.MsgSetWithdrawAddress{
		DelegatorAddress: "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
		WithdrawAddress:  "celestia1nasjhf82cjuk3mxyhzw6ntpc66exzfe7qhl256",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgSetWithdrawAddress(t *testing.T) {
	m := createMsgSetWithdrawAddress()
	blob, now := testsuite.EmptyBlock()
	position := 4

	dm, err := Message(m, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgSetWithdrawAddress,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1u5pshtqpexjmuudrvq6q335qym2zggzhp7kq0p",
				Hash:    []byte{0xe5, 0x3, 0xb, 0xac, 0x1, 0xc9, 0xa5, 0xbe, 0x71, 0xa3, 0x60, 0x34, 0x8, 0xc6, 0x80, 0x26, 0xd4, 0x24, 0x20, 0x57},
				Balance: storage.Balance{
					Total: decimal.Zero,
				},
			},
		}, {
			Type: storageTypes.TxAddressTypeWithdraw,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Address: "celestia1nasjhf82cjuk3mxyhzw6ntpc66exzfe7qhl256",
				Hash:    []byte{0x9f, 0x61, 0x2b, 0xa4, 0xea, 0xc4, 0xb9, 0x68, 0xec, 0xc4, 0xb8, 0x9d, 0xa9, 0xac, 0x38, 0xd6, 0xb2, 0x61, 0x27, 0x3e},
				Balance: storage.Balance{
					Total: decimal.Zero,
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgUnknown

type UnknownMsgType struct{}

func (u *UnknownMsgType) Reset()                               {}
func (u *UnknownMsgType) String() string                       { return "unknown" }
func (u *UnknownMsgType) ProtoMessage()                        {}
func (u *UnknownMsgType) ValidateBasic() error                 { return nil }
func (u *UnknownMsgType) GetSigners() []cosmosTypes.AccAddress { return nil }

func createMsgUnknown() cosmosTypes.Msg {
	msgUnknown := UnknownMsgType{}
	return &msgUnknown
}

func TestDecodeMsg_MsgUnknown(t *testing.T) {
	msgUnknown := createMsgUnknown()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(msgUnknown, blob.Height, blob.Block.Time, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUnknown,
		TxId:      0,
		Data:      structs.Map(msgUnknown),
		Namespace: nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
}
