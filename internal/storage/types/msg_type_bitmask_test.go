// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgTypeBits_Names(t *testing.T) {
	tests := []struct {
		name    string
		msgType []int
		want    []MsgType
	}{
		{
			name:    string(MsgUnknown),
			msgType: []int{MsgTypeBitsUnknown},
			want:    []MsgType{MsgUnknown},
		}, {
			name:    string(MsgSetWithdrawAddress),
			msgType: []int{MsgTypeBitsSetWithdrawAddress},
			want:    []MsgType{MsgSetWithdrawAddress},
		}, {
			name:    string(MsgWithdrawDelegatorReward),
			msgType: []int{MsgTypeBitsWithdrawDelegatorReward},
			want:    []MsgType{MsgWithdrawDelegatorReward},
		}, {
			name:    string(MsgWithdrawValidatorCommission),
			msgType: []int{MsgTypeBitsWithdrawValidatorCommission},
			want:    []MsgType{MsgWithdrawValidatorCommission},
		}, {
			name:    string(MsgFundCommunityPool),
			msgType: []int{MsgTypeBitsFundCommunityPool},
			want:    []MsgType{MsgFundCommunityPool},
		}, {
			name:    string(MsgCreateValidator),
			msgType: []int{MsgTypeBitsCreateValidator},
			want:    []MsgType{MsgCreateValidator},
		}, {
			name:    string(MsgEditValidator),
			msgType: []int{MsgTypeBitsEditValidator},
			want:    []MsgType{MsgEditValidator},
		}, {
			name:    string(MsgDelegate),
			msgType: []int{MsgTypeBitsDelegate},
			want:    []MsgType{MsgDelegate},
		}, {
			name:    string(MsgBeginRedelegate),
			msgType: []int{MsgTypeBitsBeginRedelegate},
			want:    []MsgType{MsgBeginRedelegate},
		}, {
			name:    string(MsgUndelegate),
			msgType: []int{MsgTypeBitsUndelegate},
			want:    []MsgType{MsgUndelegate},
		}, {
			name:    string(MsgCancelUnbondingDelegation),
			msgType: []int{MsgTypeBitsCancelUnbondingDelegation},
			want:    []MsgType{MsgCancelUnbondingDelegation},
		}, {
			name:    string(MsgUnjail),
			msgType: []int{MsgTypeBitsUnjail},
			want:    []MsgType{MsgUnjail},
		}, {
			name:    string(MsgSend),
			msgType: []int{MsgTypeBitsSend},
			want:    []MsgType{MsgSend},
		}, {
			name:    string(MsgMultiSend),
			msgType: []int{MsgTypeBitsMultiSend},
			want:    []MsgType{MsgMultiSend},
		}, {
			name:    string(MsgCreateVestingAccount),
			msgType: []int{MsgTypeBitsCreateVestingAccount},
			want:    []MsgType{MsgCreateVestingAccount},
		}, {
			name:    string(MsgCreatePermanentLockedAccount),
			msgType: []int{MsgTypeBitsCreatePermanentLockedAccount},
			want:    []MsgType{MsgCreatePermanentLockedAccount},
		}, {
			name:    string(MsgCreatePeriodicVestingAccount),
			msgType: []int{MsgTypeBitsCreatePeriodicVestingAccount},
			want:    []MsgType{MsgCreatePeriodicVestingAccount},
		}, {
			name:    string(MsgPayForBlobs),
			msgType: []int{MsgTypeBitsPayForBlobs},
			want:    []MsgType{MsgPayForBlobs},
		}, {
			name:    string(MsgGrant),
			msgType: []int{MsgTypeBitsGrant},
			want:    []MsgType{MsgGrant},
		}, {
			name:    string(MsgExec),
			msgType: []int{MsgTypeBitsExec},
			want:    []MsgType{MsgExec},
		}, {
			name:    string(MsgRevoke),
			msgType: []int{MsgTypeBitsRevoke},
			want:    []MsgType{MsgRevoke},
		}, {
			name:    string(MsgGrantAllowance),
			msgType: []int{MsgTypeBitsGrantAllowance},
			want:    []MsgType{MsgGrantAllowance},
		}, {
			name:    string(MsgRevokeAllowance),
			msgType: []int{MsgTypeBitsRevokeAllowance},
			want:    []MsgType{MsgRevokeAllowance},
		}, {
			name:    string(MsgRegisterEVMAddress),
			msgType: []int{MsgTypeBitsRegisterEVMAddress},
			want:    []MsgType{MsgRegisterEVMAddress},
		}, {
			name:    string(MsgSubmitProposal),
			msgType: []int{MsgTypeBitsSubmitProposal},
			want:    []MsgType{MsgSubmitProposal},
		}, {
			name:    string(MsgExecLegacyContent),
			msgType: []int{MsgTypeBitsExecLegacyContent},
			want:    []MsgType{MsgExecLegacyContent},
		}, {
			name:    string(MsgVote),
			msgType: []int{MsgTypeBitsVote},
			want:    []MsgType{MsgVote},
		}, {
			name:    string(MsgVoteWeighted),
			msgType: []int{MsgTypeBitsVoteWeighted},
			want:    []MsgType{MsgVoteWeighted},
		}, {
			name:    string(MsgDeposit),
			msgType: []int{MsgTypeBitsDeposit},
			want:    []MsgType{MsgDeposit},
		}, {
			name:    string(IBCTransfer),
			msgType: []int{MsgTypeBitsIBCTransfer},
			want:    []MsgType{IBCTransfer},
		}, {
			name:    string(MsgVerifyInvariant),
			msgType: []int{MsgTypeBitsVerifyInvariant},
			want:    []MsgType{MsgVerifyInvariant},
		}, {
			name:    string(MsgSubmitEvidence),
			msgType: []int{MsgTypeBitsSubmitEvidence},
			want:    []MsgType{MsgSubmitEvidence},
		}, {
			name:    string(MsgCreateGroup),
			msgType: []int{MsgTypeBitsCreateGroup},
			want:    []MsgType{MsgCreateGroup},
		}, {
			name:    string(MsgUpdateGroupMembers),
			msgType: []int{MsgTypeBitsUpdateGroupMembers},
			want:    []MsgType{MsgUpdateGroupMembers},
		}, {
			name:    string(MsgUpdateGroupAdmin),
			msgType: []int{MsgTypeBitsUpdateGroupAdmin},
			want:    []MsgType{MsgUpdateGroupAdmin},
		}, {
			name:    string(MsgUpdateGroupMetadata),
			msgType: []int{MsgTypeBitsUpdateGroupMetadata},
			want:    []MsgType{MsgUpdateGroupMetadata},
		}, {
			name:    string(MsgCreateGroupPolicy),
			msgType: []int{MsgTypeBitsCreateGroupPolicy},
			want:    []MsgType{MsgCreateGroupPolicy},
		}, {
			name:    string(MsgUpdateGroupPolicyAdmin),
			msgType: []int{MsgTypeBitsUpdateGroupPolicyAdmin},
			want:    []MsgType{MsgUpdateGroupPolicyAdmin},
		}, {
			name:    string(MsgCreateGroupWithPolicy),
			msgType: []int{MsgTypeBitsCreateGroupWithPolicy},
			want:    []MsgType{MsgCreateGroupWithPolicy},
		}, {
			name:    string(MsgUpdateGroupPolicyDecisionPolicy),
			msgType: []int{MsgTypeBitsUpdateGroupPolicyDecisionPolicy},
			want:    []MsgType{MsgUpdateGroupPolicyDecisionPolicy},
		}, {
			name:    string(MsgUpdateGroupPolicyMetadata),
			msgType: []int{MsgTypeBitsUpdateGroupPolicyMetadata},
			want:    []MsgType{MsgUpdateGroupPolicyMetadata},
		}, {
			name:    string(MsgSubmitProposalGroup),
			msgType: []int{MsgTypeBitsSubmitProposalGroup},
			want:    []MsgType{MsgSubmitProposalGroup},
		}, {
			name:    string(MsgWithdrawProposal),
			msgType: []int{MsgTypeBitsWithdrawProposal},
			want:    []MsgType{MsgWithdrawProposal},
		}, {
			name:    string(MsgVoteGroup),
			msgType: []int{MsgTypeBitsVoteGroup},
			want:    []MsgType{MsgVoteGroup},
		}, {
			name:    string(MsgExecGroup),
			msgType: []int{MsgTypeBitsExecGroup},
			want:    []MsgType{MsgExecGroup},
		}, {
			name:    string(MsgLeaveGroup),
			msgType: []int{MsgTypeBitsLeaveGroup},
			want:    []MsgType{MsgLeaveGroup},
		}, {
			name:    string(MsgSoftwareUpgrade),
			msgType: []int{MsgTypeBitsSoftwareUpgrade},
			want:    []MsgType{MsgSoftwareUpgrade},
		}, {
			name:    string(MsgCancelUpgrade),
			msgType: []int{MsgTypeBitsCancelUpgrade},
			want:    []MsgType{MsgCancelUpgrade},
		}, {
			name:    string(MsgRegisterInterchainAccount),
			msgType: []int{MsgTypeBitsRegisterInterchainAccount},
			want:    []MsgType{MsgRegisterInterchainAccount},
		}, {
			name:    string(MsgSendTx),
			msgType: []int{MsgTypeBitsSendTx},
			want:    []MsgType{MsgSendTx},
		}, {
			name:    string(MsgRegisterPayee),
			msgType: []int{MsgTypeBitsRegisterPayee},
			want:    []MsgType{MsgRegisterPayee},
		}, {
			name:    string(MsgRegisterCounterpartyPayee),
			msgType: []int{MsgTypeBitsRegisterCounterpartyPayee},
			want:    []MsgType{MsgRegisterCounterpartyPayee},
		}, {
			name:    string(MsgPayPacketFee),
			msgType: []int{MsgTypeBitsPayPacketFee},
			want:    []MsgType{MsgPayPacketFee},
		}, {
			name:    string(MsgPayPacketFeeAsync),
			msgType: []int{MsgTypeBitsPayPacketFeeAsync},
			want:    []MsgType{MsgPayPacketFeeAsync},
		}, {
			name:    string(MsgTransfer),
			msgType: []int{MsgTypeBitsTransfer},
			want:    []MsgType{MsgTransfer},
		}, {
			name:    string(MsgCreateClient),
			msgType: []int{MsgTypeBitsCreateClient},
			want:    []MsgType{MsgCreateClient},
		}, {
			name:    string(MsgUpdateClient),
			msgType: []int{MsgTypeBitsUpdateClient},
			want:    []MsgType{MsgUpdateClient},
		}, {
			name:    string(MsgUpgradeClient),
			msgType: []int{MsgTypeBitsUpgradeClient},
			want:    []MsgType{MsgUpgradeClient},
		}, {
			name:    string(MsgSubmitMisbehaviour),
			msgType: []int{MsgTypeBitsSubmitMisbehaviour},
			want:    []MsgType{MsgSubmitMisbehaviour},
		}, {
			name:    string(MsgConnectionOpenInit),
			msgType: []int{MsgTypeBitsConnectionOpenInit},
			want:    []MsgType{MsgConnectionOpenInit},
		}, {
			name:    string(MsgConnectionOpenTry),
			msgType: []int{MsgTypeBitsConnectionOpenTry},
			want:    []MsgType{MsgConnectionOpenTry},
		}, {
			name:    string(MsgConnectionOpenAck),
			msgType: []int{MsgTypeBitsConnectionOpenAck},
			want:    []MsgType{MsgConnectionOpenAck},
		}, {
			name:    string(MsgConnectionOpenConfirm),
			msgType: []int{MsgTypeBitsConnectionOpenConfirm},
			want:    []MsgType{MsgConnectionOpenConfirm},
		},
		{
			name:    string(MsgChannelOpenInit),
			msgType: []int{MsgTypeBitsChannelOpenInit},
			want:    []MsgType{MsgChannelOpenInit},
		}, {
			name:    string(MsgChannelOpenTry),
			msgType: []int{MsgTypeBitsChannelOpenTry},
			want:    []MsgType{MsgChannelOpenTry},
		}, {
			name:    string(MsgChannelOpenAck),
			msgType: []int{MsgTypeBitsChannelOpenAck},
			want:    []MsgType{MsgChannelOpenAck},
		}, {
			name:    string(MsgChannelOpenConfirm),
			msgType: []int{MsgTypeBitsChannelOpenConfirm},
			want:    []MsgType{MsgChannelOpenConfirm},
		}, {
			name:    string(MsgChannelCloseInit),
			msgType: []int{MsgTypeBitsChannelCloseInit},
			want:    []MsgType{MsgChannelCloseInit},
		}, {
			name:    string(MsgChannelCloseConfirm),
			msgType: []int{MsgTypeBitsChannelCloseConfirm},
			want:    []MsgType{MsgChannelCloseConfirm},
		}, {
			name:    string(MsgRecvPacket),
			msgType: []int{MsgTypeBitsRecvPacket},
			want:    []MsgType{MsgRecvPacket},
		}, {
			name:    string(MsgTimeout),
			msgType: []int{MsgTypeBitsTimeout},
			want:    []MsgType{MsgTimeout},
		}, {
			name:    string(MsgTimeoutOnClose),
			msgType: []int{MsgTypeBitsTimeoutOnClose},
			want:    []MsgType{MsgTimeoutOnClose},
		}, {
			name:    string(MsgAcknowledgement),
			msgType: []int{MsgTypeBitsAcknowledgement},
			want:    []MsgType{MsgAcknowledgement},
		}, {
			name:    "MsgMultiSend x2 and MsgSetWithdrawalAccess",
			msgType: []int{MsgTypeBitsMultiSend, MsgTypeBitsMultiSend, MsgTypeBitsSetWithdrawAddress},
			want:    []MsgType{MsgSetWithdrawAddress, MsgMultiSend},
		}, {
			name:    string(MsgTryUpgrade),
			msgType: []int{MsgTypeBitsTryUpgrade},
			want:    []MsgType{MsgTryUpgrade},
		}, {
			name:    string(MsgSignalVersion),
			msgType: []int{MsgTypeBitsSignalVersion},
			want:    []MsgType{MsgSignalVersion},
		}, {
			name:    string(MsgIBCSoftwareUpgrade),
			msgType: []int{MsgTypeBitsIBCSoftwareUpgrade},
			want:    []MsgType{MsgIBCSoftwareUpgrade},
		}, {
			name:    string(MsgRecoverClient),
			msgType: []int{MsgTypeBitsRecoverClient},
			want:    []MsgType{MsgRecoverClient},
		}, {
			name:    string(MsgUpdateParams),
			msgType: []int{MsgTypeBitsUpdateParams},
			want:    []MsgType{MsgUpdateParams},
		}, {
			name:    string(MsgCreateMailbox),
			msgType: []int{MsgTypeBitsCreateMailbox},
			want:    []MsgType{MsgCreateMailbox},
		}, {
			name:    string(MsgProcessMessage),
			msgType: []int{MsgTypeBitsProcessMessage},
			want:    []MsgType{MsgProcessMessage},
		}, {
			name:    string(MsgSetMailbox),
			msgType: []int{MsgTypeBitsSetMailbox},
			want:    []MsgType{MsgSetMailbox},
		}, {
			name:    string(MsgCreateCollateralToken),
			msgType: []int{MsgTypeBitsCreateCollateralToken},
			want:    []MsgType{MsgCreateCollateralToken},
		}, {
			name:    string(MsgCreateSyntheticToken),
			msgType: []int{MsgTypeBitsCreateSyntheticToken},
			want:    []MsgType{MsgCreateSyntheticToken},
		}, {
			name:    string(MsgSetToken),
			msgType: []int{MsgTypeBitsSetToken},
			want:    []MsgType{MsgSetToken},
		}, {
			name:    string(MsgEnrollRemoteRouter),
			msgType: []int{MsgTypeBitsEnrollRemoteRouter},
			want:    []MsgType{MsgEnrollRemoteRouter},
		}, {
			name:    string(MsgUnrollRemoteRouter),
			msgType: []int{MsgTypeBitsUnrollRemoteRouter},
			want:    []MsgType{MsgUnrollRemoteRouter},
		}, {
			name:    string(MsgRemoteTransfer),
			msgType: []int{MsgTypeBitsRemoteTransfer},
			want:    []MsgType{MsgRemoteTransfer},
		}, {
			name:    string(MsgUpdateMinfeeParams),
			msgType: []int{MsgTypeBitsUpdateMinfeeParams},
			want:    []MsgType{MsgUpdateMinfeeParams},
		}, {
			name:    string(MsgCreateIgp),
			msgType: []int{MsgTypeBitsCreateIgp},
			want:    []MsgType{MsgCreateIgp},
		}, {
			name:    string(MsgSetIgpOwner),
			msgType: []int{MsgTypeBitsSetIgpOwner},
			want:    []MsgType{MsgSetIgpOwner},
		}, {
			name:    string(MsgSetDestinationGasConfig),
			msgType: []int{MsgTypeBitsSetDestinationGasConfig},
			want:    []MsgType{MsgSetDestinationGasConfig},
		}, {
			name:    string(MsgPayForGas),
			msgType: []int{MsgTypeBitsPayForGas},
			want:    []MsgType{MsgPayForGas},
		}, {
			name:    string(MsgClaim),
			msgType: []int{MsgTypeBitsClaim},
			want:    []MsgType{MsgClaim},
		}, {
			name:    string(MsgCreateMerkleTreeHook),
			msgType: []int{MsgTypeBitsCreateMerkleTreeHook},
			want:    []MsgType{MsgCreateMerkleTreeHook},
		}, {
			name:    string(MsgCreateNoopHook),
			msgType: []int{MsgTypeBitsCreateNoopHook},
			want:    []MsgType{MsgCreateNoopHook},
		}, {
			name:    string(MsgCreateMessageIdMultisigIsm),
			msgType: []int{MsgTypeBitsCreateMessageIdMultisigIsm},
			want:    []MsgType{MsgCreateMessageIdMultisigIsm},
		}, {
			name:    string(MsgCreateMerkleRootMultisigIsm),
			msgType: []int{MsgTypeBitsCreateMerkleRootMultisigIsm},
			want:    []MsgType{MsgCreateMerkleRootMultisigIsm},
		}, {
			name:    string(MsgCreateNoopIsm),
			msgType: []int{MsgTypeBitsCreateNoopIsm},
			want:    []MsgType{MsgCreateNoopIsm},
		}, {
			name:    string(MsgAnnounceValidator),
			msgType: []int{MsgTypeBitsAnnounceValidator},
			want:    []MsgType{MsgAnnounceValidator},
		}, {
			name:    string(MsgCreateRoutingIsm),
			msgType: []int{MsgTypeBitsCreateRoutingIsm},
			want:    []MsgType{MsgCreateRoutingIsm},
		}, {
			name:    string(MsgSetRoutingIsmDomain),
			msgType: []int{MsgTypeBitsSetRoutingIsmDomain},
			want:    []MsgType{MsgSetRoutingIsmDomain},
		}, {
			name:    string(MsgRemoveRoutingIsmDomain),
			msgType: []int{MsgTypeBitsRemoveRoutingIsmDomain},
			want:    []MsgType{MsgRemoveRoutingIsmDomain},
		}, {
			name:    string(MsgUpdateRoutingIsmOwner),
			msgType: []int{MsgTypeBitsUpdateRoutingIsmOwner},
			want:    []MsgType{MsgUpdateRoutingIsmOwner},
		}, {
			name:    string(MsgUpdateBlobParams),
			msgType: []int{MsgTypeBitsUpdateBlobParams},
			want:    []MsgType{MsgUpdateBlobParams},
		}, {
			name:    string(MsgPruneExpiredGrants),
			msgType: []int{MsgTypeBitsPruneExpiredGrants},
			want:    []MsgType{MsgPruneExpiredGrants},
		}, {
			name:    string(MsgSetSendEnabled),
			msgType: []int{MsgTypeBitsSetSendEnabled},
			want:    []MsgType{MsgSetSendEnabled},
		}, {
			name:    string(MsgAuthorizeCircuitBreaker),
			msgType: []int{MsgTypeBitsAuthorizeCircuitBreaker},
			want:    []MsgType{MsgAuthorizeCircuitBreaker},
		}, {
			name:    string(MsgModuleQuerySafe),
			msgType: []int{MsgTypeBitsModuleQuerySafe},
			want:    []MsgType{MsgModuleQuerySafe},
		}, {
			name:    string(MsgResetCircuitBreaker),
			msgType: []int{MsgTypeBitsResetCircuitBreaker},
			want:    []MsgType{MsgResetCircuitBreaker},
		}, {
			name:    string(MsgTripCircuitBreaker),
			msgType: []int{MsgTypeBitsTripCircuitBreaker},
			want:    []MsgType{MsgTripCircuitBreaker},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mask := MsgTypeBits{
				Bits: NewEmptyBits(),
			}
			for i := range tt.msgType {
				mask.SetBit(tt.msgType[i])
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
		}, {
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
		}, {
			name:   "test 75",
			values: []MsgType{MsgTryUpgrade},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsTryUpgrade)},
		}, {
			name:   "test 76",
			values: []MsgType{MsgSignalVersion},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSignalVersion)},
		}, {
			name:   "test 77",
			values: []MsgType{MsgIBCSoftwareUpgrade},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsIBCSoftwareUpgrade)},
		}, {
			name:   "test 78",
			values: []MsgType{MsgUpdateParams},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateParams)},
		}, {
			name:   "test 79",
			values: []MsgType{MsgRecoverClient},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRecoverClient)},
		}, {
			name:   "test 80",
			values: []MsgType{MsgCreateMailbox},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreateMailbox)},
		}, {
			name:   "test 81",
			values: []MsgType{MsgProcessMessage},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsProcessMessage)},
		}, {
			name:   "test 82",
			values: []MsgType{MsgSetMailbox},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSetMailbox)},
		}, {
			name:   "test 83",
			values: []MsgType{MsgCreateCollateralToken},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreateCollateralToken)},
		}, {
			name:   "test 84",
			values: []MsgType{MsgCreateSyntheticToken},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsCreateSyntheticToken)},
		}, {
			name:   "test 85",
			values: []MsgType{MsgSetToken},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSetToken)},
		}, {
			name:   "test 86",
			values: []MsgType{MsgEnrollRemoteRouter},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsEnrollRemoteRouter)},
		}, {
			name:   "test 87",
			values: []MsgType{MsgUnrollRemoteRouter},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUnrollRemoteRouter)},
		}, {
			name:   "test 88",
			values: []MsgType{MsgRemoteTransfer},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRemoteTransfer)},
		}, {
			name:   "test 89",
			values: []MsgType{MsgUpdateMinfeeParams},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateMinfeeParams)},
		}, {
			name:   "test 90",
			values: []MsgType{MsgSetRoutingIsmDomain},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSetRoutingIsmDomain)},
		}, {
			name:   "test 91",
			values: []MsgType{MsgRemoveRoutingIsmDomain},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsRemoveRoutingIsmDomain)},
		}, {
			name:   "test 92",
			values: []MsgType{MsgUpdateRoutingIsmOwner},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateRoutingIsmOwner)},
		}, {
			name:   "test 93",
			values: []MsgType{MsgUpdateBlobParams},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsUpdateBlobParams)},
		}, {
			name:   "test 94",
			values: []MsgType{MsgPruneExpiredGrants},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsPruneExpiredGrants)},
		}, {
			name:   "test 95",
			values: []MsgType{MsgSetSendEnabled},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsSetSendEnabled)},
		}, {
			name:   "test 96",
			values: []MsgType{MsgAuthorizeCircuitBreaker},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsAuthorizeCircuitBreaker)},
		}, {
			name:   "test 97",
			values: []MsgType{MsgResetCircuitBreaker},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsResetCircuitBreaker)},
		}, {
			name:   "test 98",
			values: []MsgType{MsgTripCircuitBreaker},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsTripCircuitBreaker)},
		}, {
			name:   "test 99",
			values: []MsgType{MsgModuleQuerySafe},
			want:   MsgTypeBits{NewBitsWithPosition(MsgTypeBitsModuleQuerySafe)},
		}, {
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

func TestMarshall(t *testing.T) {
	tests := []struct {
		name string
		mask MsgTypeBits
	}{
		{
			name: "test 1",
			mask: NewMsgTypeBitMask(MsgBeginRedelegate),
		}, {
			name: "test 2",
			mask: NewMsgTypeBitMask(MsgBeginRedelegate, MsgDelegate, MsgSend),
		}, {
			name: "test 3",
			mask: NewMsgTypeBitMask(MsgAcknowledgement, MsgCancelUpgrade),
		}, {
			name: "test 4",
			mask: NewMsgTypeBitMask(MsgChannelOpenInit, MsgConnectionOpenConfirm),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.mask)
			require.NoError(t, err)

			var newMask MsgTypeBits
			err = json.Unmarshal(data, &newMask)
			require.NoError(t, err)
			require.Equal(t, tt.mask, newMask)
		})
	}
}
