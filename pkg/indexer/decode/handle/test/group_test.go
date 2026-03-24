// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/stretchr/testify/require"
)

func TestDecodeMsg_SuccessOnMsgCreateGroup(t *testing.T) {
	msg := &group.MsgCreateGroup{
		Admin: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateGroup,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgUpdateGroupMembers(t *testing.T) {
	msg := &group.MsgUpdateGroupMembers{
		Admin: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUpdateGroupMembers,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgUpdateGroupAdmin(t *testing.T) {
	msg := &group.MsgUpdateGroupAdmin{
		Admin:    "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		NewAdmin: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUpdateGroupAdmin,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      98,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgUpdateGroupMetadata(t *testing.T) {
	msg := &group.MsgUpdateGroupMetadata{
		Admin: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUpdateGroupMetadata,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgCreateGroupPolicy(t *testing.T) {
	msg := &group.MsgCreateGroupPolicy{
		Admin: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateGroupPolicy,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgUpdateGroupPolicyAdmin(t *testing.T) {
	msg := &group.MsgUpdateGroupPolicyAdmin{
		Admin: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUpdateGroupPolicyAdmin,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgCreateGroupWithPolicy(t *testing.T) {
	msg := &group.MsgCreateGroupWithPolicy{
		Admin: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateGroupWithPolicy,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgUpdateGroupPolicyDecisionPolicy(t *testing.T) {
	msg := &group.MsgUpdateGroupPolicyDecisionPolicy{
		Admin:              "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		GroupPolicyAddress: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUpdateGroupPolicyDecisionPolicy,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      98,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgUpdateGroupPolicyMetadata(t *testing.T) {
	msg := &group.MsgUpdateGroupPolicyMetadata{
		Admin:              "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
		GroupPolicyAddress: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUpdateGroupPolicyMetadata,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      98,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgSubmitProposalGroup(t *testing.T) {
	msg := &group.MsgSubmitProposal{
		GroupPolicyAddress: "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSubmitProposalGroup,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgWithdrawProposal(t *testing.T) {
	msg := &group.MsgWithdrawProposal{
		Address: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgWithdrawProposal,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgVoteGroup(t *testing.T) {
	msg := &group.MsgVote{
		Voter: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgVoteGroup,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgExecGroup(t *testing.T) {
	msg := &group.MsgExec{
		Executor: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgExecGroup,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

func TestDecodeMsg_SuccessOnMsgLeaveGroup(t *testing.T) {
	msg := &group.MsgLeaveGroup{
		Address: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgLeaveGroup,
		TxId:      0,
		Data:      mustMsgToMap(t, msg),
		Size:      49,
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}
