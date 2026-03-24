// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/x/group"
)

// MsgCreateGroup is the Msg/CreateGroup request type.
func MsgCreateGroup(ctx *context.Context, msgId uint64, m *group.MsgCreateGroup) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateGroup
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUpdateGroupMembers is the Msg/UpdateGroupMembers request type.
func MsgUpdateGroupMembers(ctx *context.Context, msgId uint64, m *group.MsgUpdateGroupMembers) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateGroupMembers
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUpdateGroupAdmin is the Msg/UpdateGroupAdmin request type.
func MsgUpdateGroupAdmin(ctx *context.Context, msgId uint64, m *group.MsgUpdateGroupAdmin) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateGroupAdmin
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
		{t: storageTypes.MsgAddressTypeNewAdmin, address: m.NewAdmin},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUpdateGroupMetadata is the Msg/UpdateGroupMetadata request type.
func MsgUpdateGroupMetadata(ctx *context.Context, msgId uint64, m *group.MsgUpdateGroupMetadata) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateGroupMetadata
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateGroupPolicy is the Msg/CreateGroupPolicy request type.
func MsgCreateGroupPolicy(ctx *context.Context, msgId uint64, m *group.MsgCreateGroupPolicy) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateGroupPolicy
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUpdateGroupPolicyAdmin is the Msg/UpdateGroupPolicyAdmin request type.
func MsgUpdateGroupPolicyAdmin(ctx *context.Context, msgId uint64, m *group.MsgUpdateGroupPolicyAdmin) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateGroupPolicyAdmin
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgCreateGroupWithPolicy is the Msg/CreateGroupWithPolicy request type.
func MsgCreateGroupWithPolicy(ctx *context.Context, msgId uint64, m *group.MsgCreateGroupWithPolicy) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateGroupWithPolicy
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUpdateGroupPolicyDecisionPolicy is the Msg/UpdateGroupPolicyDecisionPolicy request type.
func MsgUpdateGroupPolicyDecisionPolicy(ctx *context.Context, msgId uint64, m *group.MsgUpdateGroupPolicyDecisionPolicy) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateGroupPolicyDecisionPolicy
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
		{t: storageTypes.MsgAddressTypeGroupPolicyAddress, address: m.GroupPolicyAddress},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgUpdateGroupPolicyMetadata is the Msg/UpdateGroupPolicyMetadata request type.
func MsgUpdateGroupPolicyMetadata(ctx *context.Context, msgId uint64, m *group.MsgUpdateGroupPolicyMetadata) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgUpdateGroupPolicyMetadata
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
		{t: storageTypes.MsgAddressTypeGroupPolicyAddress, address: m.GroupPolicyAddress},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgSubmitProposal is the Msg/SubmitProposal request type.
func MsgSubmitProposalGroup(ctx *context.Context, msgId uint64, m *group.MsgSubmitProposal) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSubmitProposalGroup
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeGroupPolicyAddress, address: m.GroupPolicyAddress},
		// Proposers - list of proposer addresses
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgWithdrawProposal is the Msg/WithdrawProposal request type.
func MsgWithdrawProposal(ctx *context.Context, msgId uint64, m *group.MsgWithdrawProposal) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgWithdrawProposal
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Address}, // address is the admin of the group policy or one of the proposer of the proposal.
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgVote is the Msg/Vote request type.
func MsgVoteGroup(ctx *context.Context, msgId uint64, m *group.MsgVote) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgVoteGroup
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: m.Voter},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgExec is the Msg/Exec request type.
func MsgExecGroup(ctx *context.Context, msgId uint64, m *group.MsgExec) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgExecGroup
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeExecutor, address: m.Executor},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

// MsgLeaveGroup is the Msg/LeaveGroup request type.
func MsgLeaveGroup(ctx *context.Context, msgId uint64, m *group.MsgLeaveGroup) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgLeaveGroup
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeGroupMember, address: m.Address},
	}, ctx.Block.Height, msgId)
	return msgType, err
}
