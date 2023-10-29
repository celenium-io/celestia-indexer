// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgTypeBits_Names(t *testing.T) {
	tests := []struct {
		name    string
		msgType int
		want    []MsgType
	}{
		{
			name:    string(MsgUnknown),
			msgType: MsgTypeBitsUnknown,
			want:    []MsgType{MsgUnknown},
		}, {
			name:    string(MsgSetWithdrawAddress),
			msgType: MsgTypeBitsSetWithdrawAddress,
			want:    []MsgType{MsgSetWithdrawAddress},
		}, {
			name:    string(MsgWithdrawDelegatorReward),
			msgType: MsgTypeBitsWithdrawDelegatorReward,
			want:    []MsgType{MsgWithdrawDelegatorReward},
		}, {
			name:    string(MsgWithdrawValidatorCommission),
			msgType: MsgTypeBitsWithdrawValidatorCommission,
			want:    []MsgType{MsgWithdrawValidatorCommission},
		}, {
			name:    string(MsgFundCommunityPool),
			msgType: MsgTypeBitsFundCommunityPool,
			want:    []MsgType{MsgFundCommunityPool},
		}, {
			name:    string(MsgCreateValidator),
			msgType: MsgTypeBitsCreateValidator,
			want:    []MsgType{MsgCreateValidator},
		}, {
			name:    string(MsgEditValidator),
			msgType: MsgTypeBitsEditValidator,
			want:    []MsgType{MsgEditValidator},
		}, {
			name:    string(MsgDelegate),
			msgType: MsgTypeBitsDelegate,
			want:    []MsgType{MsgDelegate},
		}, {
			name:    string(MsgBeginRedelegate),
			msgType: MsgTypeBitsBeginRedelegate,
			want:    []MsgType{MsgBeginRedelegate},
		}, {
			name:    string(MsgUndelegate),
			msgType: MsgTypeBitsUndelegate,
			want:    []MsgType{MsgUndelegate},
		}, {
			name:    string(MsgCancelUnbondingDelegation),
			msgType: MsgTypeBitsCancelUnbondingDelegation,
			want:    []MsgType{MsgCancelUnbondingDelegation},
		}, {
			name:    string(MsgUnjail),
			msgType: MsgTypeBitsUnjail,
			want:    []MsgType{MsgUnjail},
		}, {
			name:    string(MsgSend),
			msgType: MsgTypeBitsSend,
			want:    []MsgType{MsgSend},
		}, {
			name:    string(MsgMultiSend),
			msgType: MsgTypeBitsMultiSend,
			want:    []MsgType{MsgMultiSend},
		}, {
			name:    string(MsgCreateVestingAccount),
			msgType: MsgTypeBitsCreateVestingAccount,
			want:    []MsgType{MsgCreateVestingAccount},
		}, {
			name:    string(MsgCreatePermanentLockedAccount),
			msgType: MsgTypeBitsCreatePermanentLockedAccount,
			want:    []MsgType{MsgCreatePermanentLockedAccount},
		}, {
			name:    string(MsgCreatePeriodicVestingAccount),
			msgType: MsgTypeBitsCreatePeriodicVestingAccount,
			want:    []MsgType{MsgCreatePeriodicVestingAccount},
		}, {
			name:    string(MsgPayForBlobs),
			msgType: MsgTypeBitsPayForBlobs,
			want:    []MsgType{MsgPayForBlobs},
		}, {
			name:    string(MsgGrant),
			msgType: MsgTypeBitsGrant,
			want:    []MsgType{MsgGrant},
		}, {
			name:    string(MsgExec),
			msgType: MsgTypeBitsExec,
			want:    []MsgType{MsgExec},
		}, {
			name:    string(MsgRevoke),
			msgType: MsgTypeBitsRevoke,
			want:    []MsgType{MsgRevoke},
		}, {
			name:    string(MsgGrantAllowance),
			msgType: MsgTypeBitsGrantAllowance,
			want:    []MsgType{MsgGrantAllowance},
		}, {
			name:    string(MsgRevokeAllowance),
			msgType: MsgTypeBitsRevokeAllowance,
			want:    []MsgType{MsgRevokeAllowance},
		}, {
			name:    string(MsgRegisterEVMAddress),
			msgType: MsgTypeBitsRegisterEVMAddress,
			want:    []MsgType{MsgRegisterEVMAddress},
		}, {
			name:    string(MsgSubmitProposal),
			msgType: MsgTypeBitsSubmitProposal,
			want:    []MsgType{MsgSubmitProposal},
		}, {
			name:    string(MsgExecLegacyContent),
			msgType: MsgTypeBitsExecLegacyContent,
			want:    []MsgType{MsgExecLegacyContent},
		}, {
			name:    string(MsgVote),
			msgType: MsgTypeBitsVote,
			want:    []MsgType{MsgVote},
		}, {
			name:    string(MsgVoteWeighted),
			msgType: MsgTypeBitsVoteWeighted,
			want:    []MsgType{MsgVoteWeighted},
		}, {
			name:    string(MsgDeposit),
			msgType: MsgTypeBitsDeposit,
			want:    []MsgType{MsgDeposit},
		}, {
			name:    string(IBCTransfer),
			msgType: MsgTypeBitsIBCTransfer,
			want:    []MsgType{IBCTransfer},
		}, {
			name:    string(MsgVerifyInvariant),
			msgType: MsgTypeBitsVerifyInvariant,
			want:    []MsgType{MsgVerifyInvariant},
		}, {
			name:    string(MsgSubmitEvidence),
			msgType: MsgTypeBitsSubmitEvidence,
			want:    []MsgType{MsgSubmitEvidence},
		}, {
			name:    string(MsgCreateGroup),
			msgType: MsgTypeBitsCreateGroup,
			want:    []MsgType{MsgCreateGroup},
		}, {
			name:    string(MsgUpdateGroupMembers),
			msgType: MsgTypeBitsUpdateGroupMembers,
			want:    []MsgType{MsgUpdateGroupMembers},
		}, {
			name:    string(MsgUpdateGroupAdmin),
			msgType: MsgTypeBitsUpdateGroupAdmin,
			want:    []MsgType{MsgUpdateGroupAdmin},
		}, {
			name:    string(MsgUpdateGroupMetadata),
			msgType: MsgTypeBitsUpdateGroupMetadata,
			want:    []MsgType{MsgUpdateGroupMetadata},
		}, {
			name:    string(MsgCreateGroupPolicy),
			msgType: MsgTypeBitsCreateGroupPolicy,
			want:    []MsgType{MsgCreateGroupPolicy},
		}, {
			name:    string(MsgUpdateGroupPolicyAdmin),
			msgType: MsgTypeBitsUpdateGroupPolicyAdmin,
			want:    []MsgType{MsgUpdateGroupPolicyAdmin},
		}, {
			name:    string(MsgCreateGroupWithPolicy),
			msgType: MsgTypeBitsCreateGroupWithPolicy,
			want:    []MsgType{MsgCreateGroupWithPolicy},
		}, {
			name:    string(MsgUpdateGroupPolicyDecisionPolicy),
			msgType: MsgTypeBitsUpdateGroupPolicyDecisionPolicy,
			want:    []MsgType{MsgUpdateGroupPolicyDecisionPolicy},
		}, {
			name:    string(MsgUpdateGroupPolicyMetadata),
			msgType: MsgTypeBitsUpdateGroupPolicyMetadata,
			want:    []MsgType{MsgUpdateGroupPolicyMetadata},
		}, {
			name:    string(MsgSubmitProposalGroup),
			msgType: MsgTypeBitsSubmitProposalGroup,
			want:    []MsgType{MsgSubmitProposalGroup},
		}, {
			name:    string(MsgWithdrawProposal),
			msgType: MsgTypeBitsWithdrawProposal,
			want:    []MsgType{MsgWithdrawProposal},
		}, {
			name:    string(MsgVoteGroup),
			msgType: MsgTypeBitsVoteGroup,
			want:    []MsgType{MsgVoteGroup},
		}, {
			name:    string(MsgExecGroup),
			msgType: MsgTypeBitsExecGroup,
			want:    []MsgType{MsgExecGroup},
		}, {
			name:    string(MsgLeaveGroup),
			msgType: MsgTypeBitsLeaveGroup,
			want:    []MsgType{MsgLeaveGroup},
		}, {
			name:    string(MsgSoftwareUpgrade),
			msgType: MsgTypeBitsSoftwareUpgrade,
			want:    []MsgType{MsgSoftwareUpgrade},
		}, {
			name:    string(MsgCancelUpgrade),
			msgType: MsgTypeBitsCancelUpgrade,
			want:    []MsgType{MsgCancelUpgrade},
		}, {
			name:    string(MsgRegisterInterchainAccount),
			msgType: MsgTypeBitsRegisterInterchainAccount,
			want:    []MsgType{MsgRegisterInterchainAccount},
		}, {
			name:    string(MsgSendTx),
			msgType: MsgTypeBitsSendTx,
			want:    []MsgType{MsgSendTx},
		}, {
			name:    string(MsgRegisterPayee),
			msgType: MsgTypeBitsRegisterPayee,
			want:    []MsgType{MsgRegisterPayee},
		}, {
			name:    string(MsgRegisterCounterpartyPayee),
			msgType: MsgTypeBitsRegisterCounterpartyPayee,
			want:    []MsgType{MsgRegisterCounterpartyPayee},
		}, {
			name:    string(MsgPayPacketFee),
			msgType: MsgTypeBitsPayPacketFee,
			want:    []MsgType{MsgPayPacketFee},
		}, {
			name:    string(MsgPayPacketFeeAsync),
			msgType: MsgTypeBitsPayPacketFeeAsync,
			want:    []MsgType{MsgPayPacketFeeAsync},
		}, {
			name:    string(MsgTransfer),
			msgType: MsgTypeBitsTransfer,
			want:    []MsgType{MsgTransfer},
		}, {
			name:    string(MsgCreateClient),
			msgType: MsgTypeBitsCreateClient,
			want:    []MsgType{MsgCreateClient},
		}, {
			name:    string(MsgUpdateClient),
			msgType: MsgTypeBitsUpdateClient,
			want:    []MsgType{MsgUpdateClient},
		}, {
			name:    string(MsgUpgradeClient),
			msgType: MsgTypeBitsUpgradeClient,
			want:    []MsgType{MsgUpgradeClient},
		}, {
			name:    string(MsgSubmitMisbehaviour),
			msgType: MsgTypeBitsSubmitMisbehaviour,
			want:    []MsgType{MsgSubmitMisbehaviour},
		}, {
			name:    string(MsgConnectionOpenInit),
			msgType: MsgTypeBitsConnectionOpenInit,
			want:    []MsgType{MsgConnectionOpenInit},
		}, {
			name:    string(MsgConnectionOpenTry),
			msgType: MsgTypeBitsConnectionOpenTry,
			want:    []MsgType{MsgConnectionOpenTry},
		}, {
			name:    string(MsgConnectionOpenAck),
			msgType: MsgTypeBitsConnectionOpenAck,
			want:    []MsgType{MsgConnectionOpenAck},
		}, {
			name:    string(MsgConnectionOpenConfirm),
			msgType: MsgTypeBitsConnectionOpenConfirm,
			want:    []MsgType{MsgConnectionOpenConfirm},
		},
		{
			name:    string(MsgChannelOpenInit),
			msgType: MsgTypeBitsChannelOpenInit,
			want:    []MsgType{MsgChannelOpenInit},
		}, {
			name:    string(MsgChannelOpenTry),
			msgType: MsgTypeBitsChannelOpenTry,
			want:    []MsgType{MsgChannelOpenTry},
		}, {
			name:    string(MsgChannelOpenAck),
			msgType: MsgTypeBitsChannelOpenAck,
			want:    []MsgType{MsgChannelOpenAck},
		}, {
			name:    string(MsgChannelOpenConfirm),
			msgType: MsgTypeBitsChannelOpenConfirm,
			want:    []MsgType{MsgChannelOpenConfirm},
		}, {
			name:    string(MsgChannelCloseInit),
			msgType: MsgTypeBitsChannelCloseInit,
			want:    []MsgType{MsgChannelCloseInit},
		}, {
			name:    string(MsgChannelCloseConfirm),
			msgType: MsgTypeBitsChannelCloseConfirm,
			want:    []MsgType{MsgChannelCloseConfirm},
		}, {
			name:    string(MsgRecvPacket),
			msgType: MsgTypeBitsRecvPacket,
			want:    []MsgType{MsgRecvPacket},
		}, {
			name:    string(MsgTimeout),
			msgType: MsgTypeBitsTimeout,
			want:    []MsgType{MsgTimeout},
		}, {
			name:    string(MsgTimeoutOnClose),
			msgType: MsgTypeBitsTimeoutOnClose,
			want:    []MsgType{MsgTimeoutOnClose},
		}, {
			name:    string(MsgAcknowledgement),
			msgType: MsgTypeBitsAcknowledgement,
			want:    []MsgType{MsgAcknowledgement},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mask := MsgTypeBits{
				Bits: NewEmptyBits(),
			}
			mask.SetBit(tt.msgType)
			require.Equal(t, tt.want, mask.Names())
		})
	}
}

func TestNewMsgTypeBitMask(t *testing.T) {
	tests := []struct {
		name   string
		values []MsgType
		want   MsgTypeBits
	}{
		{
			name:   "test 0",
			values: nil,
			want:   MsgTypeBits{NewBits(0)},
		}, {
			name:   "test 1",
			values: []MsgType{MsgUnknown},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUnknown)},
		}, {
			name:   "test 2",
			values: []MsgType{MsgSetWithdrawAddress},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSetWithdrawAddress)},
		}, {
			name:   "test 3",
			values: []MsgType{MsgWithdrawDelegatorReward},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsWithdrawDelegatorReward)},
		}, {
			name:   "test 4",
			values: []MsgType{MsgWithdrawValidatorCommission},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsWithdrawValidatorCommission)},
		}, {
			name:   "test 5",
			values: []MsgType{MsgFundCommunityPool},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsFundCommunityPool)},
		}, {
			name:   "test 6",
			values: []MsgType{MsgCreateValidator},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreateValidator)},
		}, {
			name:   "test 7",
			values: []MsgType{MsgEditValidator},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsEditValidator)},
		}, {
			name:   "test 8",
			values: []MsgType{MsgDelegate},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsDelegate)},
		}, {
			name:   "test 9",
			values: []MsgType{MsgBeginRedelegate},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsBeginRedelegate)},
		}, {
			name:   "test 10",
			values: []MsgType{MsgUndelegate},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUndelegate)},
		}, {
			name:   "test 11",
			values: []MsgType{MsgCancelUnbondingDelegation},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCancelUnbondingDelegation)},
		}, {
			name:   "test 12",
			values: []MsgType{MsgUnjail},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUnjail)},
		}, {
			name:   "test 13",
			values: []MsgType{MsgSend},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSend)},
		}, {
			name:   "test 14",
			values: []MsgType{MsgMultiSend},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsMultiSend)},
		}, {
			name:   "test 15",
			values: []MsgType{MsgCreateVestingAccount},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreateVestingAccount)},
		}, {
			name:   "test 16",
			values: []MsgType{MsgCreatePermanentLockedAccount},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreatePermanentLockedAccount)},
		}, {
			name:   "test 17",
			values: []MsgType{MsgCreatePeriodicVestingAccount},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreatePeriodicVestingAccount)},
		}, {
			name:   "test 18",
			values: []MsgType{MsgPayForBlobs},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsPayForBlobs)},
		}, {
			name:   "test 19",
			values: []MsgType{MsgGrant},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsGrant)},
		}, {
			name:   "test 20",
			values: []MsgType{MsgExec},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsExec)},
		}, {
			name:   "test 21",
			values: []MsgType{MsgRevoke},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRevoke)},
		}, {
			name:   "test 22",
			values: []MsgType{MsgGrantAllowance},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsGrantAllowance)},
		}, {
			name:   "test 23",
			values: []MsgType{MsgRevokeAllowance},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRevokeAllowance)},
		}, {
			name:   "test 24",
			values: []MsgType{MsgRegisterEVMAddress},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRegisterEVMAddress)},
		}, {
			name:   "test 25",
			values: []MsgType{MsgSubmitProposal},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSubmitProposal)},
		}, {
			name:   "test 26",
			values: []MsgType{MsgExecLegacyContent},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsExecLegacyContent)},
		}, {
			name:   "test 27",
			values: []MsgType{MsgVote},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsVote)},
		}, {
			name:   "test 28",
			values: []MsgType{MsgVoteWeighted},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsVoteWeighted)},
		}, {
			name:   "test 29",
			values: []MsgType{MsgDeposit},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsDeposit)},
		}, {
			name:   "test 30",
			values: []MsgType{IBCTransfer},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsIBCTransfer)},
		}, {
			name:   "test 31",
			values: []MsgType{MsgVerifyInvariant},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsVerifyInvariant)},
		}, {
			name:   "test 32",
			values: []MsgType{MsgSubmitEvidence},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSubmitEvidence)},
		}, {
			name:   "test 33",
			values: []MsgType{MsgSendNFT},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSendNFT)},
		}, {
			name:   "test 34",
			values: []MsgType{MsgCreateGroup},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreateGroup)},
		}, {
			name:   "test 35",
			values: []MsgType{MsgUpdateGroupMembers},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateGroupMembers)},
		}, {
			name:   "test 36",
			values: []MsgType{MsgUpdateGroupAdmin},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateGroupAdmin)},
		}, {
			name:   "test 37",
			values: []MsgType{MsgUpdateGroupMetadata},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateGroupMetadata)},
		}, {
			name:   "test 38",
			values: []MsgType{MsgCreateGroupPolicy},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreateGroupPolicy)},
		}, {
			name:   "test 39",
			values: []MsgType{MsgUpdateGroupPolicyAdmin},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateGroupPolicyAdmin)},
		}, {
			name:   "test 40",
			values: []MsgType{MsgCreateGroupWithPolicy},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreateGroupWithPolicy)},
		}, {
			name:   "test 41",
			values: []MsgType{MsgUpdateGroupPolicyDecisionPolicy},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateGroupPolicyDecisionPolicy)},
		}, {
			name:   "test 42",
			values: []MsgType{MsgUpdateGroupPolicyMetadata},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateGroupPolicyMetadata)},
		}, {
			name:   "test 43",
			values: []MsgType{MsgSubmitProposalGroup},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSubmitProposalGroup)},
		}, {
			name:   "test 44",
			values: []MsgType{MsgWithdrawProposal},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsWithdrawProposal)},
		}, {
			name:   "test 45",
			values: []MsgType{MsgVoteGroup},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsVoteGroup)},
		}, {
			name:   "test 46",
			values: []MsgType{MsgExecGroup},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsExecGroup)},
		}, {
			name:   "test 47",
			values: []MsgType{MsgLeaveGroup},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsLeaveGroup)},
		}, {
			name:   "test 48",
			values: []MsgType{MsgSoftwareUpgrade},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSoftwareUpgrade)},
		}, {
			name:   "test 49",
			values: []MsgType{MsgCancelUpgrade},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCancelUpgrade)},
		}, {
			name:   "test 50",
			values: []MsgType{MsgRegisterInterchainAccount},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRegisterInterchainAccount)},
		}, {
			name:   "test 51",
			values: []MsgType{MsgSendTx},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSendTx)},
		}, {
			name:   "test 52",
			values: []MsgType{MsgRegisterPayee},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRegisterPayee)},
		}, {
			name:   "test 53",
			values: []MsgType{MsgRegisterCounterpartyPayee},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRegisterCounterpartyPayee)},
		}, {
			name:   "test 54",
			values: []MsgType{MsgPayPacketFee},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsPayPacketFee)},
		}, {
			name:   "test 55",
			values: []MsgType{MsgPayPacketFeeAsync},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsPayPacketFeeAsync)},
		}, {
			name:   "test 56",
			values: []MsgType{MsgTransfer},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsTransfer)},
		}, {
			name:   "test 57",
			values: []MsgType{MsgCreateClient},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreateClient)},
		}, {
			name:   "test 58",
			values: []MsgType{MsgUpdateClient},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateClient)},
		}, {
			name:   "test 59",
			values: []MsgType{MsgUpgradeClient},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpgradeClient)},
		}, {
			name:   "test 60",
			values: []MsgType{MsgSubmitMisbehaviour},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSubmitMisbehaviour)},
		}, {
			name:   "test 61",
			values: []MsgType{MsgConnectionOpenInit},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsConnectionOpenInit)},
		}, {
			name:   "test 62",
			values: []MsgType{MsgConnectionOpenTry},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsConnectionOpenTry)},
		}, {
			name:   "test 63",
			values: []MsgType{MsgConnectionOpenAck},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsConnectionOpenAck)},
		}, {
			name:   "test 64",
			values: []MsgType{MsgConnectionOpenConfirm},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsConnectionOpenConfirm)},
		},

		{
			name:   "test 65",
			values: []MsgType{MsgChannelOpenInit},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsChannelOpenInit)},
		}, {
			name:   "test 66",
			values: []MsgType{MsgChannelOpenTry},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsChannelOpenTry)},
		}, {
			name:   "test 67",
			values: []MsgType{MsgChannelOpenAck},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsChannelOpenAck)},
		}, {
			name:   "test 68",
			values: []MsgType{MsgChannelOpenConfirm},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsChannelOpenConfirm)},
		}, {
			name:   "test 69",
			values: []MsgType{MsgChannelCloseInit},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsChannelCloseInit)},
		}, {
			name:   "test 70",
			values: []MsgType{MsgChannelCloseConfirm},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsChannelCloseConfirm)},
		}, {
			name:   "test 71",
			values: []MsgType{MsgRecvPacket},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRecvPacket)},
		}, {
			name:   "test 72",
			values: []MsgType{MsgTimeout},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsTimeout)},
		}, {
			name:   "test 73",
			values: []MsgType{MsgTimeoutOnClose},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsTimeoutOnClose)},
		}, {
			name:   "test 74",
			values: []MsgType{MsgAcknowledgement},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsAcknowledgement)},
		},

		{
			name:   "test combo",
			values: []MsgType{MsgWithdrawDelegatorReward, MsgBeginRedelegate},
			want:   MsgTypeBits{NewBits(260)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.want, NewMsgTypeBitMask(tt.values...))
		})
	}
}

func TestMsgTypeBits_SetByMsgType(t *testing.T) {
	tests := []struct {
		name  string
		value MsgType
		want  MsgTypeBits
	}{
		{
			name:  "test 1",
			value: MsgBeginRedelegate,
			want:  NewMsgTypeBitMask(MsgBeginRedelegate),
		}, {
			name:  "test 2",
			value: MsgType("unknown"),
			want:  NewMsgTypeBitMask(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mask := &MsgTypeBits{
				Bits: NewEmptyBits(),
			}
			mask.SetByMsgType(tt.value)
			require.EqualValues(t, tt.want.Bits, mask.Bits)
		})
	}
}

func TestMsgTypeBits_HasOne(t *testing.T) {
	tests := []struct {
		name  string
		mask  MsgTypeBits
		value MsgTypeBits
		want  bool
	}{
		{
			name:  "test 1",
			mask:  NewMsgTypeBitMask(MsgBeginRedelegate),
			value: NewMsgTypeBitMask(MsgBeginRedelegate),
			want:  true,
		}, {
			name:  "test 2",
			mask:  NewMsgTypeBitMask(MsgBeginRedelegate, MsgDelegate, MsgSend),
			value: NewMsgTypeBitMask(MsgBeginRedelegate),
			want:  true,
		}, {
			name:  "test 3",
			mask:  NewMsgTypeBitMask(MsgBeginRedelegate),
			value: NewMsgTypeBitMask(MsgBeginRedelegate, MsgDelegate, MsgSend),
			want:  true,
		}, {
			name:  "test 4",
			mask:  NewMsgTypeBitMask(MsgBeginRedelegate),
			value: NewMsgTypeBitMask(MsgDelegate, MsgSend),
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			has := tt.mask.HasOne(tt.value)
			require.Equal(t, tt.want, has)
		})
	}
}
