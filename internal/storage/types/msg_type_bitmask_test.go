// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgTypeBits_Names(t *testing.T) {
	tests := []struct {
		name string
		Bits Bits
		want []MsgType
	}{
		{
			name: string(MsgUnknown),
			Bits: Bits(MsgTypeBitsUnknown),
			want: []MsgType{MsgUnknown},
		}, {
			name: string(MsgSetWithdrawAddress),
			Bits: Bits(MsgTypeBitsSetWithdrawAddress),
			want: []MsgType{MsgSetWithdrawAddress},
		}, {
			name: string(MsgWithdrawDelegatorReward),
			Bits: Bits(MsgTypeBitsWithdrawDelegatorReward),
			want: []MsgType{MsgWithdrawDelegatorReward},
		}, {
			name: string(MsgWithdrawValidatorCommission),
			Bits: Bits(MsgTypeBitsWithdrawValidatorCommission),
			want: []MsgType{MsgWithdrawValidatorCommission},
		}, {
			name: string(MsgFundCommunityPool),
			Bits: Bits(MsgTypeBitsFundCommunityPool),
			want: []MsgType{MsgFundCommunityPool},
		}, {
			name: string(MsgCreateValidator),
			Bits: Bits(MsgTypeBitsCreateValidator),
			want: []MsgType{MsgCreateValidator},
		}, {
			name: string(MsgEditValidator),
			Bits: Bits(MsgTypeBitsEditValidator),
			want: []MsgType{MsgEditValidator},
		}, {
			name: string(MsgDelegate),
			Bits: Bits(MsgTypeBitsDelegate),
			want: []MsgType{MsgDelegate},
		}, {
			name: string(MsgBeginRedelegate),
			Bits: Bits(MsgTypeBitsBeginRedelegate),
			want: []MsgType{MsgBeginRedelegate},
		}, {
			name: string(MsgUndelegate),
			Bits: Bits(MsgTypeBitsUndelegate),
			want: []MsgType{MsgUndelegate},
		}, {
			name: string(MsgCancelUnbondingDelegation),
			Bits: Bits(MsgTypeBitsCancelUnbondingDelegation),
			want: []MsgType{MsgCancelUnbondingDelegation},
		}, {
			name: string(MsgUnjail),
			Bits: Bits(MsgTypeBitsUnjail),
			want: []MsgType{MsgUnjail},
		}, {
			name: string(MsgSend),
			Bits: Bits(MsgTypeBitsSend),
			want: []MsgType{MsgSend},
		}, {
			name: string(MsgMultiSend),
			Bits: Bits(MsgTypeBitsMultiSend),
			want: []MsgType{MsgMultiSend},
		}, {
			name: string(MsgCreateVestingAccount),
			Bits: Bits(MsgTypeBitsCreateVestingAccount),
			want: []MsgType{MsgCreateVestingAccount},
		}, {
			name: string(MsgCreatePermanentLockedAccount),
			Bits: Bits(MsgTypeBitsCreatePermanentLockedAccount),
			want: []MsgType{MsgCreatePermanentLockedAccount},
		}, {
			name: string(MsgCreatePeriodicVestingAccount),
			Bits: Bits(MsgTypeBitsCreatePeriodicVestingAccount),
			want: []MsgType{MsgCreatePeriodicVestingAccount},
		}, {
			name: string(MsgPayForBlobs),
			Bits: Bits(MsgTypeBitsPayForBlobs),
			want: []MsgType{MsgPayForBlobs},
		}, {
			name: string(MsgGrant),
			Bits: Bits(MsgTypeBitsGrant),
			want: []MsgType{MsgGrant},
		}, {
			name: string(MsgExec),
			Bits: Bits(MsgTypeBitsExec),
			want: []MsgType{MsgExec},
		}, {
			name: string(MsgRevoke),
			Bits: Bits(MsgTypeBitsRevoke),
			want: []MsgType{MsgRevoke},
		}, {
			name: string(MsgGrantAllowance),
			Bits: Bits(MsgTypeBitsGrantAllowance),
			want: []MsgType{MsgGrantAllowance},
		}, {
			name: string(MsgRevokeAllowance),
			Bits: Bits(MsgTypeBitsRevokeAllowance),
			want: []MsgType{MsgRevokeAllowance},
		}, {
			name: string(MsgRegisterEVMAddress),
			Bits: Bits(MsgTypeBitsRegisterEVMAddress),
			want: []MsgType{MsgRegisterEVMAddress},
		}, {
			name: string(MsgSubmitProposal),
			Bits: Bits(MsgTypeBitsSubmitProposal),
			want: []MsgType{MsgSubmitProposal},
		}, {
			name: string(MsgExecLegacyContent),
			Bits: Bits(MsgTypeBitsExecLegacyContent),
			want: []MsgType{MsgExecLegacyContent},
		}, {
			name: string(MsgVote),
			Bits: Bits(MsgTypeBitsVote),
			want: []MsgType{MsgVote},
		}, {
			name: string(MsgVoteWeighted),
			Bits: Bits(MsgTypeBitsVoteWeighted),
			want: []MsgType{MsgVoteWeighted},
		}, {
			name: string(MsgDeposit),
			Bits: Bits(MsgTypeBitsDeposit),
			want: []MsgType{MsgDeposit},
		}, {
			name: string(IBCTransfer),
			Bits: Bits(MsgTypeBitsIBCTransfer),
			want: []MsgType{IBCTransfer},
		}, {
			name: string(MsgVerifyInvariant),
			Bits: Bits(MsgTypeBitsVerifyInvariant),
			want: []MsgType{MsgVerifyInvariant},
		}, {
			name: string(MsgSubmitEvidence),
			Bits: Bits(MsgTypeBitsSubmitEvidence),
			want: []MsgType{MsgSubmitEvidence},
		}, {
			name: string(MsgCreateGroup),
			Bits: Bits(MsgTypeBitsCreateGroup),
			want: []MsgType{MsgCreateGroup},
		}, {
			name: string(MsgUpdateGroupMembers),
			Bits: Bits(MsgTypeBitsUpdateGroupMembers),
			want: []MsgType{MsgUpdateGroupMembers},
		}, {
			name: string(MsgUpdateGroupAdmin),
			Bits: Bits(MsgTypeBitsUpdateGroupAdmin),
			want: []MsgType{MsgUpdateGroupAdmin},
		}, {
			name: string(MsgUpdateGroupMetadata),
			Bits: Bits(MsgTypeBitsUpdateGroupMetadata),
			want: []MsgType{MsgUpdateGroupMetadata},
		}, {
			name: string(MsgCreateGroupPolicy),
			Bits: Bits(MsgTypeBitsCreateGroupPolicy),
			want: []MsgType{MsgCreateGroupPolicy},
		}, {
			name: string(MsgUpdateGroupPolicyAdmin),
			Bits: Bits(MsgTypeBitsUpdateGroupPolicyAdmin),
			want: []MsgType{MsgUpdateGroupPolicyAdmin},
		}, {
			name: string(MsgCreateGroupWithPolicy),
			Bits: Bits(MsgTypeBitsCreateGroupWithPolicy),
			want: []MsgType{MsgCreateGroupWithPolicy},
		}, {
			name: string(MsgUpdateGroupPolicyDecisionPolicy),
			Bits: Bits(MsgTypeBitsUpdateGroupPolicyDecisionPolicy),
			want: []MsgType{MsgUpdateGroupPolicyDecisionPolicy},
		}, {
			name: string(MsgUpdateGroupPolicyMetadata),
			Bits: Bits(MsgTypeBitsUpdateGroupPolicyMetadata),
			want: []MsgType{MsgUpdateGroupPolicyMetadata},
		}, {
			name: string(MsgSubmitProposalGroup),
			Bits: Bits(MsgTypeBitsSubmitProposalGroup),
			want: []MsgType{MsgSubmitProposalGroup},
		}, {
			name: string(MsgWithdrawProposal),
			Bits: Bits(MsgTypeBitsWithdrawProposal),
			want: []MsgType{MsgWithdrawProposal},
		}, {
			name: string(MsgVoteGroup),
			Bits: Bits(MsgTypeBitsVoteGroup),
			want: []MsgType{MsgVoteGroup},
		}, {
			name: string(MsgExecGroup),
			Bits: Bits(MsgTypeBitsExecGroup),
			want: []MsgType{MsgExecGroup},
		}, {
			name: string(MsgLeaveGroup),
			Bits: Bits(MsgTypeBitsLeaveGroup),
			want: []MsgType{MsgLeaveGroup},
		}, {
			name: string(MsgSoftwareUpgrade),
			Bits: Bits(MsgTypeBitsSoftwareUpgrade),
			want: []MsgType{MsgSoftwareUpgrade},
		}, {
			name: string(MsgCancelUpgrade),
			Bits: Bits(MsgTypeBitsCancelUpgrade),
			want: []MsgType{MsgCancelUpgrade},
		}, {
			name: string(MsgRegisterInterchainAccount),
			Bits: Bits(MsgTypeBitsRegisterInterchainAccount),
			want: []MsgType{MsgRegisterInterchainAccount},
		}, {
			name: string(MsgSendTx),
			Bits: Bits(MsgTypeBitsSendTx),
			want: []MsgType{MsgSendTx},
		}, {
			name: string(MsgRegisterPayee),
			Bits: Bits(MsgTypeBitsRegisterPayee),
			want: []MsgType{MsgRegisterPayee},
		}, {
			name: string(MsgRegisterCounterpartyPayee),
			Bits: Bits(MsgTypeBitsRegisterCounterpartyPayee),
			want: []MsgType{MsgRegisterCounterpartyPayee},
		}, {
			name: string(MsgPayPacketFee),
			Bits: Bits(MsgTypeBitsPayPacketFee),
			want: []MsgType{MsgPayPacketFee},
		}, {
			name: string(MsgPayPacketFeeAsync),
			Bits: Bits(MsgTypeBitsPayPacketFeeAsync),
			want: []MsgType{MsgPayPacketFeeAsync},
		}, {
			name: string(MsgTransfer),
			Bits: Bits(MsgTypeBitsTransfer),
			want: []MsgType{MsgTransfer},
		}, {
			name: string(MsgCreateClient),
			Bits: Bits(MsgTypeBitsCreateClient),
			want: []MsgType{MsgCreateClient},
		}, {
			name: string(MsgUpdateClient),
			Bits: Bits(MsgTypeBitsUpdateClient),
			want: []MsgType{MsgUpdateClient},
		}, {
			name: string(MsgUpgradeClient),
			Bits: Bits(MsgTypeBitsUpgradeClient),
			want: []MsgType{MsgUpgradeClient},
		}, {
			name: string(MsgSubmitMisbehaviour),
			Bits: Bits(MsgTypeBitsSubmitMisbehaviour),
			want: []MsgType{MsgSubmitMisbehaviour},
		}, {
			name: string(MsgConnectionOpenInit),
			Bits: Bits(MsgTypeBitsConnectionOpenInit),
			want: []MsgType{MsgConnectionOpenInit},
		}, {
			name: string(MsgConnectionOpenTry),
			Bits: Bits(MsgTypeBitsConnectionOpenTry),
			want: []MsgType{MsgConnectionOpenTry},
		}, {
			name: string(MsgConnectionOpenAck),
			Bits: Bits(MsgTypeBitsConnectionOpenAck),
			want: []MsgType{MsgConnectionOpenAck},
		}, {
			name: string(MsgConnectionOpenConfirm),
			Bits: Bits(MsgTypeBitsConnectionOpenConfirm),
			want: []MsgType{MsgConnectionOpenConfirm},
		},
		// {
		// 	name: string(MsgChannelOpenInit),
		// 	Bits: Bits(MsgTypeBitsChannelOpenInit),
		// 	want: []MsgType{MsgChannelOpenInit},
		// }, {
		// 	name: string(MsgChannelOpenTry),
		// 	Bits: Bits(MsgTypeBitsChannelOpenTry),
		// 	want: []MsgType{MsgChannelOpenTry},
		// }, {
		// 	name: string(MsgChannelOpenAck),
		// 	Bits: Bits(MsgTypeBitsChannelOpenAck),
		// 	want: []MsgType{MsgChannelOpenAck},
		// }, {
		// 	name: string(MsgChannelOpenConfirm),
		// 	Bits: Bits(MsgTypeBitsChannelOpenConfirm),
		// 	want: []MsgType{MsgChannelOpenConfirm},
		// }, {
		// 	name: string(MsgChannelCloseInit),
		// 	Bits: Bits(MsgTypeBitsChannelCloseInit),
		// 	want: []MsgType{MsgChannelCloseInit},
		// }, {
		// 	name: string(MsgChannelCloseConfirm),
		// 	Bits: Bits(MsgTypeBitsChannelCloseConfirm),
		// 	want: []MsgType{MsgChannelCloseConfirm},
		// }, {
		// 	name: string(MsgRecvPacket),
		// 	Bits: Bits(MsgTypeBitsRecvPacket),
		// 	want: []MsgType{MsgRecvPacket},
		// }, {
		// 	name: string(MsgTimeout),
		// 	Bits: Bits(MsgTypeBitsTimeout),
		// 	want: []MsgType{MsgTimeout},
		// }, {
		// 	name: string(MsgTimeoutOnClose),
		// 	Bits: Bits(MsgTypeBitsTimeoutOnClose),
		// 	want: []MsgType{MsgTimeoutOnClose},
		// }, {
		// 	name: string(MsgAcknowledgement),
		// 	Bits: Bits(MsgTypeBitsAcknowledgement),
		// 	want: []MsgType{MsgAcknowledgement},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mask := MsgTypeBits{
				Bits: tt.Bits,
			}
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
			want:   MsgTypeBits{Bits(0)},
		}, {
			name:   "test 1",
			values: []MsgType{MsgUnknown},
			want:   MsgTypeBits{Bits(MsgTypeBitsUnknown)},
		}, {
			name:   "test 2",
			values: []MsgType{MsgSetWithdrawAddress},
			want:   MsgTypeBits{Bits(MsgTypeBitsSetWithdrawAddress)},
		}, {
			name:   "test 3",
			values: []MsgType{MsgWithdrawDelegatorReward},
			want:   MsgTypeBits{Bits(MsgTypeBitsWithdrawDelegatorReward)},
		}, {
			name:   "test 4",
			values: []MsgType{MsgWithdrawValidatorCommission},
			want:   MsgTypeBits{Bits(MsgTypeBitsWithdrawValidatorCommission)},
		}, {
			name:   "test 5",
			values: []MsgType{MsgFundCommunityPool},
			want:   MsgTypeBits{Bits(MsgTypeBitsFundCommunityPool)},
		}, {
			name:   "test 6",
			values: []MsgType{MsgCreateValidator},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateValidator)},
		}, {
			name:   "test 7",
			values: []MsgType{MsgEditValidator},
			want:   MsgTypeBits{Bits(MsgTypeBitsEditValidator)},
		}, {
			name:   "test 8",
			values: []MsgType{MsgDelegate},
			want:   MsgTypeBits{Bits(MsgTypeBitsDelegate)},
		}, {
			name:   "test 9",
			values: []MsgType{MsgBeginRedelegate},
			want:   MsgTypeBits{Bits(MsgTypeBitsBeginRedelegate)},
		}, {
			name:   "test 10",
			values: []MsgType{MsgUndelegate},
			want:   MsgTypeBits{Bits(MsgTypeBitsUndelegate)},
		}, {
			name:   "test 11",
			values: []MsgType{MsgCancelUnbondingDelegation},
			want:   MsgTypeBits{Bits(MsgTypeBitsCancelUnbondingDelegation)},
		}, {
			name:   "test 12",
			values: []MsgType{MsgUnjail},
			want:   MsgTypeBits{Bits(MsgTypeBitsUnjail)},
		}, {
			name:   "test 13",
			values: []MsgType{MsgSend},
			want:   MsgTypeBits{Bits(MsgTypeBitsSend)},
		}, {
			name:   "test 14",
			values: []MsgType{MsgMultiSend},
			want:   MsgTypeBits{Bits(MsgTypeBitsMultiSend)},
		}, {
			name:   "test 15",
			values: []MsgType{MsgCreateVestingAccount},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateVestingAccount)},
		}, {
			name:   "test 16",
			values: []MsgType{MsgCreatePermanentLockedAccount},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreatePermanentLockedAccount)},
		}, {
			name:   "test 17",
			values: []MsgType{MsgCreatePeriodicVestingAccount},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreatePeriodicVestingAccount)},
		}, {
			name:   "test 18",
			values: []MsgType{MsgPayForBlobs},
			want:   MsgTypeBits{Bits(MsgTypeBitsPayForBlobs)},
		}, {
			name:   "test 19",
			values: []MsgType{MsgGrant},
			want:   MsgTypeBits{Bits(MsgTypeBitsGrant)},
		}, {
			name:   "test 20",
			values: []MsgType{MsgExec},
			want:   MsgTypeBits{Bits(MsgTypeBitsExec)},
		}, {
			name:   "test 21",
			values: []MsgType{MsgRevoke},
			want:   MsgTypeBits{Bits(MsgTypeBitsRevoke)},
		}, {
			name:   "test 22",
			values: []MsgType{MsgGrantAllowance},
			want:   MsgTypeBits{Bits(MsgTypeBitsGrantAllowance)},
		}, {
			name:   "test 23",
			values: []MsgType{MsgRevokeAllowance},
			want:   MsgTypeBits{Bits(MsgTypeBitsRevokeAllowance)},
		}, {
			name:   "test 24",
			values: []MsgType{MsgRegisterEVMAddress},
			want:   MsgTypeBits{Bits(MsgTypeBitsRegisterEVMAddress)},
		}, {
			name:   "test 25",
			values: []MsgType{MsgSubmitProposal},
			want:   MsgTypeBits{Bits(MsgTypeBitsSubmitProposal)},
		}, {
			name:   "test 26",
			values: []MsgType{MsgExecLegacyContent},
			want:   MsgTypeBits{Bits(MsgTypeBitsExecLegacyContent)},
		}, {
			name:   "test 27",
			values: []MsgType{MsgVote},
			want:   MsgTypeBits{Bits(MsgTypeBitsVote)},
		}, {
			name:   "test 28",
			values: []MsgType{MsgVoteWeighted},
			want:   MsgTypeBits{Bits(MsgTypeBitsVoteWeighted)},
		}, {
			name:   "test 29",
			values: []MsgType{MsgDeposit},
			want:   MsgTypeBits{Bits(MsgTypeBitsDeposit)},
		}, {
			name:   "test 30",
			values: []MsgType{IBCTransfer},
			want:   MsgTypeBits{Bits(MsgTypeBitsIBCTransfer)},
		}, {
			name:   "test 31",
			values: []MsgType{MsgVerifyInvariant},
			want:   MsgTypeBits{Bits(MsgTypeBitsVerifyInvariant)},
		}, {
			name:   "test 32",
			values: []MsgType{MsgSubmitEvidence},
			want:   MsgTypeBits{Bits(MsgTypeBitsSubmitEvidence)},
		}, {
			name:   "test 33",
			values: []MsgType{MsgSendNFT},
			want:   MsgTypeBits{Bits(MsgTypeBitsSendNFT)},
		}, {
			name:   "test 34",
			values: []MsgType{MsgCreateGroup},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateGroup)},
		}, {
			name:   "test 35",
			values: []MsgType{MsgUpdateGroupMembers},
			want:   MsgTypeBits{Bits(MsgTypeBitsUpdateGroupMembers)},
		}, {
			name:   "test 36",
			values: []MsgType{MsgUpdateGroupAdmin},
			want:   MsgTypeBits{Bits(MsgTypeBitsUpdateGroupAdmin)},
		}, {
			name:   "test 37",
			values: []MsgType{MsgUpdateGroupMetadata},
			want:   MsgTypeBits{Bits(MsgTypeBitsUpdateGroupMetadata)},
		}, {
			name:   "test 38",
			values: []MsgType{MsgCreateGroupPolicy},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateGroupPolicy)},
		}, {
			name:   "test 39",
			values: []MsgType{MsgUpdateGroupPolicyAdmin},
			want:   MsgTypeBits{Bits(MsgTypeBitsUpdateGroupPolicyAdmin)},
		}, {
			name:   "test 40",
			values: []MsgType{MsgCreateGroupWithPolicy},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateGroupWithPolicy)},
		}, {
			name:   "test 41",
			values: []MsgType{MsgUpdateGroupPolicyDecisionPolicy},
			want:   MsgTypeBits{Bits(MsgTypeBitsUpdateGroupPolicyDecisionPolicy)},
		}, {
			name:   "test 42",
			values: []MsgType{MsgUpdateGroupPolicyMetadata},
			want:   MsgTypeBits{Bits(MsgTypeBitsUpdateGroupPolicyMetadata)},
		}, {
			name:   "test 43",
			values: []MsgType{MsgSubmitProposalGroup},
			want:   MsgTypeBits{Bits(MsgTypeBitsSubmitProposalGroup)},
		}, {
			name:   "test 44",
			values: []MsgType{MsgWithdrawProposal},
			want:   MsgTypeBits{Bits(MsgTypeBitsWithdrawProposal)},
		}, {
			name:   "test 45",
			values: []MsgType{MsgVoteGroup},
			want:   MsgTypeBits{Bits(MsgTypeBitsVoteGroup)},
		}, {
			name:   "test 46",
			values: []MsgType{MsgExecGroup},
			want:   MsgTypeBits{Bits(MsgTypeBitsExecGroup)},
		}, {
			name:   "test 47",
			values: []MsgType{MsgLeaveGroup},
			want:   MsgTypeBits{Bits(MsgTypeBitsLeaveGroup)},
		}, {
			name:   "test 48",
			values: []MsgType{MsgSoftwareUpgrade},
			want:   MsgTypeBits{Bits(MsgTypeBitsSoftwareUpgrade)},
		}, {
			name:   "test 49",
			values: []MsgType{MsgCancelUpgrade},
			want:   MsgTypeBits{Bits(MsgTypeBitsCancelUpgrade)},
		}, {
			name:   "test 50",
			values: []MsgType{MsgRegisterInterchainAccount},
			want:   MsgTypeBits{Bits(MsgTypeBitsRegisterInterchainAccount)},
		}, {
			name:   "test 51",
			values: []MsgType{MsgSendTx},
			want:   MsgTypeBits{Bits(MsgTypeBitsSendTx)},
		}, {
			name:   "test 52",
			values: []MsgType{MsgRegisterPayee},
			want:   MsgTypeBits{Bits(MsgTypeBitsRegisterPayee)},
		}, {
			name:   "test 53",
			values: []MsgType{MsgRegisterCounterpartyPayee},
			want:   MsgTypeBits{Bits(MsgTypeBitsRegisterCounterpartyPayee)},
		}, {
			name:   "test 54",
			values: []MsgType{MsgPayPacketFee},
			want:   MsgTypeBits{Bits(MsgTypeBitsPayPacketFee)},
		}, {
			name:   "test 55",
			values: []MsgType{MsgPayPacketFeeAsync},
			want:   MsgTypeBits{Bits(MsgTypeBitsPayPacketFeeAsync)},
		}, {
			name:   "test 56",
			values: []MsgType{MsgTransfer},
			want:   MsgTypeBits{Bits(MsgTypeBitsTransfer)},
		}, {
			name:   "test 57",
			values: []MsgType{MsgCreateClient},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateClient)},
		}, {
			name:   "test 58",
			values: []MsgType{MsgUpdateClient},
			want:   MsgTypeBits{Bits(MsgTypeBitsUpdateClient)},
		}, {
			name:   "test 59",
			values: []MsgType{MsgUpgradeClient},
			want:   MsgTypeBits{Bits(MsgTypeBitsUpgradeClient)},
		}, {
			name:   "test 60",
			values: []MsgType{MsgSubmitMisbehaviour},
			want:   MsgTypeBits{Bits(MsgTypeBitsSubmitMisbehaviour)},
		}, {
			name:   "test 61",
			values: []MsgType{MsgConnectionOpenInit},
			want:   MsgTypeBits{Bits(MsgTypeBitsConnectionOpenInit)},
		}, {
			name:   "test 62",
			values: []MsgType{MsgConnectionOpenTry},
			want:   MsgTypeBits{Bits(MsgTypeBitsConnectionOpenTry)},
		}, {
			name:   "test 63",
			values: []MsgType{MsgConnectionOpenAck},
			want:   MsgTypeBits{Bits(MsgTypeBitsConnectionOpenAck)},
		}, {
			name:   "test 64",
			values: []MsgType{MsgConnectionOpenConfirm},
			want:   MsgTypeBits{Bits(MsgTypeBitsConnectionOpenConfirm)},
		},

		// {
		// 	name:   "test 65",
		// 	values: []MsgType{MsgChannelOpenInit},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsChannelOpenInit)},
		// }, {
		// 	name:   "test 66",
		// 	values: []MsgType{MsgChannelOpenTry},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsChannelOpenTry)},
		// }, {
		// 	name:   "test 67",
		// 	values: []MsgType{MsgChannelOpenAck},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsChannelOpenAck)},
		// }, {
		// 	name:   "test 68",
		// 	values: []MsgType{MsgChannelOpenConfirm},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsChannelOpenConfirm)},
		// }, {
		// 	name:   "test 69",
		// 	values: []MsgType{MsgChannelCloseInit},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsChannelCloseInit)},
		// }, {
		// 	name:   "test 70",
		// 	values: []MsgType{MsgChannelCloseConfirm},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsChannelCloseConfirm)},
		// }, {
		// 	name:   "test 71",
		// 	values: []MsgType{MsgRecvPacket},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsRecvPacket)},
		// }, {
		// 	name:   "test 72",
		// 	values: []MsgType{MsgTimeout},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsTimeout)},
		// }, {
		// 	name:   "test 73",
		// 	values: []MsgType{MsgTimeoutOnClose},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsTimeoutOnClose)},
		// }, {
		// 	name:   "test 74",
		// 	values: []MsgType{MsgAcknowledgement},
		// 	want:   MsgTypeBits{Bits(MsgTypeBitsAcknowledgement)},
		// },

		{
			name:   "test combo",
			values: []MsgType{MsgWithdrawDelegatorReward, MsgBeginRedelegate},
			want:   MsgTypeBits{Bits(260)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.want, NewMsgTypeBitMask(tt.values...))
		})
	}
}

func TestMsgTypeBits_SetBit(t *testing.T) {
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
				Bits: 0,
			}
			mask.SetBit(tt.value)
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
