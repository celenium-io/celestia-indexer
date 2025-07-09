// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	cosmosGovTypesV1Beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

func createExpectations(
	blob nodeTypes.BlockData,
	now time.Time,
	m types.Msg,
	position int,
	addrType storageTypes.MsgAddressType,
	address string,
	hash []byte,
	txType storageTypes.MsgType,
	size int,
) ([]storage.AddressWithType, storage.Message) {
	addressesExpected := []storage.AddressWithType{
		{
			Type: addrType,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    address,
				Hash:       hash,
				Balance:    storage.EmptyBalance(),
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  int64(position),
		Type:      txType,
		TxId:      0,
		Data:      structs.Map(m),
		Size:      size,
		Namespace: nil,
		Addresses: addressesExpected,
	}
	return addressesExpected, msgExpected
}

// v1.MsgSubmitProposal

func createMsgSubmitProposalV1() types.Msg {
	// Data from: ADDAF8EA30C75A7B3A069B1F9E24975CA6EA769CC42A850AD816432B4B0BE38F
	value := "CmUKLi9jb3Ntb3MucGFyYW1zLnYxYmV0YTEuUGFyYW1ldGVyQ2hhbmdlUHJvcG9zYWwSMwoFdGl0bGUSC2Rlc2NyaXB0aW9uGh0KB3N0YWtpbmcSDU1heFZhbGlkYXRvcnMaAzEwMxISCgR1dGlhEgoxMDAwMDAwMDAwGi9jZWxlc3RpYTEwZDA3eTI2NWdtbXV2dDR6MHc5YXc4ODBqbnNyNzAwanRnejR2Nw=="
	val, _ := base64.StdEncoding.DecodeString(value)
	m := cosmosGovTypesV1.MsgSubmitProposal{
		Messages: []*codecTypes.Any{
			{
				TypeUrl: "/cosmos.gov.v1beta1.MsgSubmitProposal",
				Value:   val,
			},
		},
		InitialDeposit: make([]types.Coin, 0),
		Proposer:       "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
		Metadata:       "",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgSubmitProposal_V1(t *testing.T) {
	m := createMsgSubmitProposalV1()
	blob, now := testsuite.EmptyBlock()
	position := 7

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeProposer,
		"celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
		[]byte{123, 95, 226, 43, 84, 70, 247, 198, 46, 162, 123, 139, 215, 28, 239, 148, 224, 63, 61, 242},
		storageTypes.MsgSubmitProposal,
		266,
	)

	msgExpected.Proposal = &storage.Proposal{
		Height:    blob.Height,
		CreatedAt: blob.Block.Time,
		Proposer: &storage.Address{
			Address: "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
		},
		Status:      storageTypes.ProposalStatusInactive,
		Type:        storageTypes.ProposalTypeText,
		Title:       "Proposal with messages",
		Description: "Proposal contains messages:\r\n1. /cosmos.gov.v1beta1.MsgSubmitProposal\r\n",
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// v1beta1.MsgSubmitProposal

func createMsgSubmitProposalV1Beta1() types.Msg {
	// Data from: 8137be4397d31f265c3c61f11b9ed79601ea6c407b5517ddd5b259b989c135ef
	value := "Chnwn5KOQ2VsZXN0aWEgQWlyZHJvcCDinIUgElhHZXQg8J+SjkNlbGVzdGlhIEFpcmRyb3Ag4pyFIHZpc2l0aW5nIHVybDogd3d3LlRlcnJhUHJvLmF0CgotIG1vcmUgaW5mbzogd3d3LlRlcnJhV2ViLmF0"
	val, _ := base64.StdEncoding.DecodeString(value)
	m := cosmosGovTypesV1Beta1.MsgSubmitProposal{
		Content: &codecTypes.Any{
			TypeUrl: "/cosmos.gov.v1beta1.TextProposal",
			Value:   val,
		},
		InitialDeposit: make(types.Coins, 0),
		Proposer:       "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgSubmitProposal_V1Beta1(t *testing.T) {
	m := createMsgSubmitProposalV1Beta1()
	blob, now := testsuite.EmptyBlock()
	position := 8

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeProposer,
		"celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
		[]byte{123, 95, 226, 43, 84, 70, 247, 198, 46, 162, 123, 139, 215, 28, 239, 148, 224, 63, 61, 242},
		storageTypes.MsgSubmitProposal,
		205,
	)
	msgExpected.Data["Content"] = cosmosGovTypesV1Beta1.TextProposal{
		Title:       "💎Celestia Airdrop ✅ ",
		Description: "Get 💎Celestia Airdrop ✅ visiting url: www.TerraPro.at\n\n- more info: www.TerraWeb.at",
	}
	msgExpected.Proposal = &storage.Proposal{
		Title:       "💎Celestia Airdrop ✅ ",
		Description: "Get 💎Celestia Airdrop ✅ visiting url: www.TerraPro.at\n\n- more info: www.TerraWeb.at",
		Height:      blob.Height,
		CreatedAt:   blob.Block.Time,
		Type:        storageTypes.ProposalTypeText,
		Status:      storageTypes.ProposalStatusInactive,
		Proposer: &storage.Address{
			Address: "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgExecLegacyContent

func createMsgExecLegacyContent() types.Msg {
	m := cosmosGovTypesV1.MsgExecLegacyContent{
		Authority: "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
		Content:   nil,
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgExecLegacyContent(t *testing.T) {
	m := createMsgExecLegacyContent()
	blob, now := testsuite.EmptyBlock()
	position := 9

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeAuthority,
		"celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7",
		[]byte{123, 95, 226, 43, 84, 70, 247, 198, 46, 162, 123, 139, 215, 28, 239, 148, 224, 63, 61, 242},
		storageTypes.MsgExecLegacyContent,
		49,
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// v1.MsgVote

func createMsgVoteV1() types.Msg {
	// Data from: 0A4BA0A30449C3269F313B5D974560F8D3A8179BE994054724898FF2D6866928
	m := cosmosGovTypesV1.MsgVote{
		ProposalId: 1,
		Voter:      "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		Option:     cosmosGovTypesV1.VoteOption_VOTE_OPTION_YES,
		Metadata:   "",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgVote_V1(t *testing.T) {
	m := createMsgVoteV1()
	blob, now := testsuite.EmptyBlock()
	position := 7

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeVoter,
		"celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		[]byte{8, 204, 180, 93, 112, 144, 218, 230, 174, 203, 58, 172, 76, 199, 190, 39, 45, 188, 116, 154},
		storageTypes.MsgVote,
		53,
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// v1beta1.MsgVote

func createMsgVoteV1Beta1() types.Msg {
	// Data from: 0A4BA0A30449C3269F313B5D974560F8D3A8179BE994054724898FF2D6866928
	m := cosmosGovTypesV1Beta1.MsgVote{
		ProposalId: 1,
		Voter:      "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		Option:     cosmosGovTypesV1Beta1.OptionYes,
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgVote_V1Beta1(t *testing.T) {
	m := createMsgVoteV1Beta1()
	blob, now := testsuite.EmptyBlock()
	position := 8

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeVoter,
		"celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		[]byte{8, 204, 180, 93, 112, 144, 218, 230, 174, 203, 58, 172, 76, 199, 190, 39, 45, 188, 116, 154},
		storageTypes.MsgVote,
		53,
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// v1.MsgVoteWeighted

func createMsgVoteWeightedV1() types.Msg {
	m := cosmosGovTypesV1.MsgVoteWeighted{
		ProposalId: 1,
		Voter:      "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		Options:    make([]*cosmosGovTypesV1.WeightedVoteOption, 0),
		Metadata:   "",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgVoteWeighted_V1(t *testing.T) {
	m := createMsgVoteWeightedV1()
	blob, now := testsuite.EmptyBlock()
	position := 7

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeVoter,
		"celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		[]byte{8, 204, 180, 93, 112, 144, 218, 230, 174, 203, 58, 172, 76, 199, 190, 39, 45, 188, 116, 154},
		storageTypes.MsgVoteWeighted,
		51,
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// v1beta1.MsgVoteWeighted

func createMsgVoteWeightedV1Beta1() types.Msg {
	m := cosmosGovTypesV1Beta1.MsgVoteWeighted{
		ProposalId: 1,
		Voter:      "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		Options:    make([]cosmosGovTypesV1Beta1.WeightedVoteOption, 0),
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgVoteWeighted_V1Beta1(t *testing.T) {
	m := createMsgVoteWeightedV1Beta1()
	blob, now := testsuite.EmptyBlock()
	position := 8

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeVoter,
		"celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		[]byte{8, 204, 180, 93, 112, 144, 218, 230, 174, 203, 58, 172, 76, 199, 190, 39, 45, 188, 116, 154},
		storageTypes.MsgVoteWeighted,
		51,
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// v1.MsgDeposit

func createMsgDepositV1() types.Msg {
	m := cosmosGovTypesV1.MsgDeposit{
		ProposalId: 1,
		Depositor:  "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		Amount:     make([]types.Coin, 0),
	}

	return &m
}

func TestDecodeMsg_SuccessMsgDeposit_V1(t *testing.T) {
	m := createMsgDepositV1()
	blob, now := testsuite.EmptyBlock()
	position := 7

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeDepositor,
		"celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		[]byte{8, 204, 180, 93, 112, 144, 218, 230, 174, 203, 58, 172, 76, 199, 190, 39, 45, 188, 116, 154},
		storageTypes.MsgDeposit,
		51,
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// v1beta1.MsgDeposit

func createMsgDepositV1Beta1() types.Msg {
	m := cosmosGovTypesV1Beta1.MsgDeposit{
		ProposalId: 1,
		Depositor:  "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		Amount:     make(types.Coins, 0),
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgDeposit_V1Beta1(t *testing.T) {
	m := createMsgDepositV1Beta1()
	blob, now := testsuite.EmptyBlock()
	position := 8

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: blob.Height,
		Time:   blob.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess)

	addressesExpected, msgExpected := createExpectations(
		blob, now, m, position,
		storageTypes.MsgAddressTypeDepositor,
		"celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		[]byte{8, 204, 180, 93, 112, 144, 218, 230, 174, 203, 58, 172, 76, 199, 190, 39, 45, 188, 116, 154},
		storageTypes.MsgDeposit,
		51,
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}
