package parser

import (
	"cosmossdk.io/math"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
	cosmosCodecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosFeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createMsgPayForBlob() cosmosTypes.Msg {

	msgPayForBlob := appBlobTypes.MsgPayForBlobs{
		Signer:           "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		Namespaces:       [][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22}},
		BlobSizes:        []uint32{1},
		ShareCommitments: [][]byte{{176, 28, 134, 119, 32, 117, 87, 107, 231, 67, 121, 255, 209, 106, 52, 99, 88, 183, 85, 36, 67, 137, 98, 199, 144, 159, 13, 178, 111, 190, 121, 36}},
		ShareVersions:    []uint32{0},
	}

	return &msgPayForBlob
}

func TestDecodeMsg_SuccessOnPayForBlob(t *testing.T) {
	msgPayForBlob := createMsgPayForBlob()
	blob, now := createEmptyBlock()
	position := 0

	msg, blobSize, err := decodeMsg(blob, msgPayForBlob, position)

	msgExpected := storage.Message{
		Id:       0,
		Height:   blob.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgTypePayForBlobs,
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

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), blobSize)
	assert.Equal(t, msgExpected, msg)
}

func createMsgDelegate() cosmosTypes.Msg {

	msgDelegate := cosmosStakingTypes.MsgDelegate{
		DelegatorAddress: "celestia1h2kqw44hdq5dwlcvsw8f2l49lkehtf9wp95kth",
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		Amount: cosmosTypes.Coin{
			Denom:  "utia",
			Amount: math.NewInt(1000),
		},
	}

	return &msgDelegate
}

func TestDecodeMsg_SuccessOnMsgDelegate(t *testing.T) {
	msgDelegate := createMsgDelegate()
	blob, now := createEmptyBlock()
	position := 0

	msg, blobSize, err := decodeMsg(blob, msgDelegate, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeDelegate,
		TxId:      0,
		Data:      structs.Map(msgDelegate),
		Namespace: nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), blobSize)
	assert.Equal(t, msgExpected, msg)
}

func createMsgSend() cosmosTypes.Msg {

	msgDelegate := cosmosBankTypes.MsgSend{
		FromAddress: "celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Amount: cosmosTypes.Coins{
			cosmosTypes.Coin{
				Denom:  "utia",
				Amount: math.NewInt(1000),
			},
		},
	}

	return &msgDelegate
}

func TestDecodeMsg_SuccessOnMsgSend(t *testing.T) {
	msgSend := createMsgSend()
	blob, now := createEmptyBlock()
	position := 0

	msg, blobSize, err := decodeMsg(blob, msgSend, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeSend,
		TxId:      0,
		Data:      structs.Map(msgSend),
		Namespace: nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), blobSize)
	assert.Equal(t, msgExpected, msg)
}

func createMsgGrantAllowance() cosmosTypes.Msg {

	msgDelegate := cosmosFeegrant.MsgGrantAllowance{
		Granter:   "celestia1l9qjhhnxc0t6tt93q8396gu0vttwlcc233gyvr",
		Grantee:   "celestia1vut644llcgwyvysmma6ww2xkmdytc8xspty8kx",
		Allowance: cosmosCodecTypes.UnsafePackAny(cosmosFeegrant.BasicAllowance{}),
	}

	return &msgDelegate
}

func TestDecodeMsg_SuccessOnMsgGrantAllowance(t *testing.T) {
	msgGrantAllowance := createMsgGrantAllowance()
	blob, now := createEmptyBlock()
	position := 4

	msg, blobSize, err := decodeMsg(blob, msgGrantAllowance, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgTypeGrantAllowance,
		TxId:      0,
		Data:      structs.Map(msgGrantAllowance),
		Namespace: nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), blobSize)
	assert.Equal(t, msgExpected, msg)
}

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
	blob, now := createEmptyBlock()
	position := 0

	msg, blobSize, err := decodeMsg(blob, msgUnknown, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeUnknown,
		TxId:      0,
		Data:      structs.Map(msgUnknown),
		Namespace: nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), blobSize)
	assert.Equal(t, msgExpected, msg)
}
