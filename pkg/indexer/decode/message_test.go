package decode

import (
	"testing"

	"cosmossdk.io/math"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
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
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
				Balance: decimal.Zero,
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
				Hash:    []byte("celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
				Balance: decimal.Zero,
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
		EvmAddress:        "0x10E0271ec47d55511a047516f2a7301801d55eaB",
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
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
				Balance: decimal.Zero,
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
		ValidatorDstAddress: "celestiavaloper1fg9l3xvfuu9wxremv2288777zawysg4r40gw7x",
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
				Hash:    []byte("celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorSrcAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorDstAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2288777zawysg4r40gw7x"),
				Balance: decimal.Zero,
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
		DelegatorAddress:  "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre77",
		ValidatorAddress:  "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw7x",
		Pubkey:            nil,
		Value:             cosmosTypes.Coin{},
		EvmAddress:        "",
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
				Hash:    []byte("celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre77"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw7x"),
				Balance: decimal.Zero,
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
		DelegatorAddress: "celestia1A2kqw44hdq5dwlcvsw8f2l49lkehtf9wp95kth",
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g77x",
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
				Hash:    []byte("celestia1A2kqw44hdq5dwlcvsw8f2l49lkehtf9wp95kth"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g77x"),
				Balance: decimal.Zero,
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
		DelegatorAddress: "celestia1A2kqw44hdq5dwlcvsw8f2l49lkehtf9wp99kth",
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g88x",
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
				Hash:    []byte("celestia1A2kqw44hdq5dwlcvsw8f2l49lkehtf9wp99kth"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g88x"),
				Balance: decimal.Zero,
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
		ValidatorAddr: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g11x",
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
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g11x"),
				Balance: decimal.Zero,
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
		FromAddress: "celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7",
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
				Hash:    []byte("celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeToAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"),
				Balance: decimal.Zero,
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
		FromAddress: "celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7",
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
				Hash:    []byte("celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeToAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"),
				Balance: decimal.Zero,
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
		FromAddress:    "celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7",
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
				Hash:    []byte("celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeToAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"),
				Balance: decimal.Zero,
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
		Signer:           "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkr777",
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
				Hash:    []byte("celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkr777"),
				Balance: decimal.Zero,
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
		Granter:   "celestia1l9qjhhnxc0t6tt93q8396gu0vttwlcc238gyvr",
		Grantee:   "celestia1vut644llcgwyvysmma6ww2xkmdytc8xspty5kx",
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
				Hash:    []byte("celestia1l9qjhhnxc0t6tt93q8396gu0vttwlcc238gyvr"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeGrantee,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1vut644llcgwyvysmma6ww2xkmdytc8xspty5kx"),
				Balance: decimal.Zero,
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
