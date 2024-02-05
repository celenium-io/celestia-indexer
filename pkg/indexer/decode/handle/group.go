// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cosmos/cosmos-sdk/x/group"
)

// MsgCreateGroup is the Msg/CreateGroup request type.
func MsgCreateGroup(level types.Level, m *group.MsgCreateGroup) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateGroup
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, level)
	return msgType, addresses, err
}

// MsgUpdateGroupMembers is the Msg/UpdateGroupMembers request type.
func MsgUpdateGroupMembers(level types.Level, m *group.MsgUpdateGroupMembers) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateGroupMembers
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, level)
	return msgType, addresses, err
}

// MsgUpdateGroupAdmin is the Msg/UpdateGroupAdmin request type.
func MsgUpdateGroupAdmin(level types.Level, m *group.MsgUpdateGroupAdmin) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateGroupAdmin
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
		{t: storageTypes.MsgAddressTypeNewAdmin, address: m.NewAdmin},
	}, level)
	return msgType, addresses, err
}

// MsgUpdateGroupMetadata is the Msg/UpdateGroupMetadata request type.
func MsgUpdateGroupMetadata(level types.Level, m *group.MsgUpdateGroupMetadata) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateGroupMetadata
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, level)
	return msgType, addresses, err
}

// MsgCreateGroupPolicy is the Msg/CreateGroupPolicy request type.
func MsgCreateGroupPolicy(level types.Level, m *group.MsgCreateGroupPolicy) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateGroupPolicy
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, level)
	return msgType, addresses, err
}

// MsgUpdateGroupPolicyAdmin is the Msg/UpdateGroupPolicyAdmin request type.
func MsgUpdateGroupPolicyAdmin(level types.Level, m *group.MsgUpdateGroupPolicyAdmin) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateGroupPolicyAdmin
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, level)
	return msgType, addresses, err
}

// MsgCreateGroupWithPolicy is the Msg/CreateGroupWithPolicy request type.
func MsgCreateGroupWithPolicy(level types.Level, m *group.MsgCreateGroupWithPolicy) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgCreateGroupWithPolicy
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
	}, level)
	return msgType, addresses, err
}

// MsgUpdateGroupPolicyDecisionPolicy is the Msg/UpdateGroupPolicyDecisionPolicy request type.
func MsgUpdateGroupPolicyDecisionPolicy(level types.Level, m *group.MsgUpdateGroupPolicyDecisionPolicy) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateGroupPolicyDecisionPolicy
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
		{t: storageTypes.MsgAddressTypeGroupPolicyAddress, address: m.GroupPolicyAddress},
	}, level)
	return msgType, addresses, err
}

// MsgUpdateGroupPolicyMetadata is the Msg/UpdateGroupPolicyMetadata request type.
func MsgUpdateGroupPolicyMetadata(level types.Level, m *group.MsgUpdateGroupPolicyMetadata) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgUpdateGroupPolicyMetadata
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Admin},
		{t: storageTypes.MsgAddressTypeGroupPolicyAddress, address: m.GroupPolicyAddress},
	}, level)
	return msgType, addresses, err
}

// MsgSubmitProposal is the Msg/SubmitProposal request type.
func MsgSubmitProposalGroup(level types.Level, m *group.MsgSubmitProposal) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSubmitProposalGroup
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGroupPolicyAddress, address: m.GroupPolicyAddress},
		// Proposers - list of proposer addresses
	}, level)
	return msgType, addresses, err
}

// MsgWithdrawProposal is the Msg/WithdrawProposal request type.
func MsgWithdrawProposal(level types.Level, m *group.MsgWithdrawProposal) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgWithdrawProposal
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeAdmin, address: m.Address}, // address is the admin of the group policy or one of the proposer of the proposal.
	}, level)
	return msgType, addresses, err
}

// MsgVote is the Msg/Vote request type.
func MsgVoteGroup(level types.Level, m *group.MsgVote) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVoteGroup
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: m.Voter},
	}, level)
	return msgType, addresses, err
}

// MsgExec is the Msg/Exec request type.
func MsgExecGroup(level types.Level, m *group.MsgExec) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgExecGroup
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeExecutor, address: m.Executor},
	}, level)
	return msgType, addresses, err
}

// MsgLeaveGroup is the Msg/LeaveGroup request type.
func MsgLeaveGroup(level types.Level, m *group.MsgLeaveGroup) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgLeaveGroup
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeGroupMember, address: m.Address},
	}, level)
	return msgType, addresses, err
}
