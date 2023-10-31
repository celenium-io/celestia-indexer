// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/pkg/errors"
)

type MsgTypeBits struct {
	Bits
}

func NewMsgTypeBits() MsgTypeBits {
	return MsgTypeBits{NewEmptyBits()}
}

const (
	MsgTypeBitsUnknown int = iota

	MsgTypeBitsSetWithdrawAddress
	MsgTypeBitsWithdrawDelegatorReward
	MsgTypeBitsWithdrawValidatorCommission
	MsgTypeBitsFundCommunityPool

	MsgTypeBitsCreateValidator
	MsgTypeBitsEditValidator
	MsgTypeBitsDelegate
	MsgTypeBitsBeginRedelegate
	MsgTypeBitsUndelegate
	MsgTypeBitsCancelUnbondingDelegation

	MsgTypeBitsUnjail

	MsgTypeBitsSend
	MsgTypeBitsMultiSend

	MsgTypeBitsCreateVestingAccount
	MsgTypeBitsCreatePermanentLockedAccount
	MsgTypeBitsCreatePeriodicVestingAccount

	MsgTypeBitsPayForBlobs

	MsgTypeBitsGrant
	MsgTypeBitsExec
	MsgTypeBitsRevoke

	MsgTypeBitsGrantAllowance
	MsgTypeBitsRevokeAllowance

	MsgTypeBitsRegisterEVMAddress

	MsgTypeBitsSubmitProposal
	MsgTypeBitsExecLegacyContent
	MsgTypeBitsVote
	MsgTypeBitsVoteWeighted
	MsgTypeBitsDeposit

	MsgTypeBitsIBCTransfer

	MsgTypeBitsVerifyInvariant

	MsgTypeBitsSubmitEvidence

	MsgTypeBitsSendNFT

	MsgTypeBitsCreateGroup
	MsgTypeBitsUpdateGroupMembers
	MsgTypeBitsUpdateGroupAdmin
	MsgTypeBitsUpdateGroupMetadata
	MsgTypeBitsCreateGroupPolicy
	MsgTypeBitsUpdateGroupPolicyAdmin
	MsgTypeBitsCreateGroupWithPolicy
	MsgTypeBitsUpdateGroupPolicyDecisionPolicy
	MsgTypeBitsUpdateGroupPolicyMetadata
	MsgTypeBitsSubmitProposalGroup
	MsgTypeBitsWithdrawProposal
	MsgTypeBitsVoteGroup
	MsgTypeBitsExecGroup
	MsgTypeBitsLeaveGroup

	MsgTypeBitsSoftwareUpgrade
	MsgTypeBitsCancelUpgrade

	MsgTypeBitsRegisterInterchainAccount
	MsgTypeBitsSendTx

	MsgTypeBitsRegisterPayee
	MsgTypeBitsRegisterCounterpartyPayee
	MsgTypeBitsPayPacketFee
	MsgTypeBitsPayPacketFeeAsync

	MsgTypeBitsTransfer

	MsgTypeBitsCreateClient
	MsgTypeBitsUpdateClient
	MsgTypeBitsUpgradeClient
	MsgTypeBitsSubmitMisbehaviour

	MsgTypeBitsConnectionOpenInit
	MsgTypeBitsConnectionOpenTry
	MsgTypeBitsConnectionOpenAck
	MsgTypeBitsConnectionOpenConfirm

	MsgTypeBitsChannelOpenInit
	MsgTypeBitsChannelOpenTry
	MsgTypeBitsChannelOpenAck
	MsgTypeBitsChannelOpenConfirm
	MsgTypeBitsChannelCloseInit
	MsgTypeBitsChannelCloseConfirm
	MsgTypeBitsRecvPacket
	MsgTypeBitsTimeout
	MsgTypeBitsTimeoutOnClose
	MsgTypeBitsAcknowledgement
)

func NewMsgTypeBitMask(values ...MsgType) MsgTypeBits {
	mask := NewMsgTypeBits()
	for i := range values {
		mask.SetByMsgType(values[i])
	}
	return mask
}

func (mask *MsgTypeBits) SetByMsgType(value MsgType) {
	switch value {
	case MsgUnknown:
		mask.SetBit(MsgTypeBitsUnknown)

	case MsgSetWithdrawAddress:
		mask.SetBit(MsgTypeBitsSetWithdrawAddress)
	case MsgWithdrawDelegatorReward:
		mask.SetBit(MsgTypeBitsWithdrawDelegatorReward)
	case MsgWithdrawValidatorCommission:
		mask.SetBit(MsgTypeBitsWithdrawValidatorCommission)
	case MsgFundCommunityPool:
		mask.SetBit(MsgTypeBitsFundCommunityPool)

	case MsgCreateValidator:
		mask.SetBit(MsgTypeBitsCreateValidator)
	case MsgEditValidator:
		mask.SetBit(MsgTypeBitsEditValidator)
	case MsgDelegate:
		mask.SetBit(MsgTypeBitsDelegate)
	case MsgBeginRedelegate:
		mask.SetBit(MsgTypeBitsBeginRedelegate)
	case MsgUndelegate:
		mask.SetBit(MsgTypeBitsUndelegate)
	case MsgCancelUnbondingDelegation:
		mask.SetBit(MsgTypeBitsCancelUnbondingDelegation)

	case MsgUnjail:
		mask.SetBit(MsgTypeBitsUnjail)

	case MsgSend:
		mask.SetBit(MsgTypeBitsSend)
	case MsgMultiSend:
		mask.SetBit(MsgTypeBitsMultiSend)

	case MsgCreateVestingAccount:
		mask.SetBit(MsgTypeBitsCreateVestingAccount)
	case MsgCreatePermanentLockedAccount:
		mask.SetBit(MsgTypeBitsCreatePermanentLockedAccount)
	case MsgCreatePeriodicVestingAccount:
		mask.SetBit(MsgTypeBitsCreatePeriodicVestingAccount)

	case MsgPayForBlobs:
		mask.SetBit(MsgTypeBitsPayForBlobs)

	case MsgGrant:
		mask.SetBit(MsgTypeBitsGrant)
	case MsgExec:
		mask.SetBit(MsgTypeBitsExec)
	case MsgRevoke:
		mask.SetBit(MsgTypeBitsRevoke)

	case MsgGrantAllowance:
		mask.SetBit(MsgTypeBitsGrantAllowance)
	case MsgRevokeAllowance:
		mask.SetBit(MsgTypeBitsRevokeAllowance)

	case MsgRegisterEVMAddress:
		mask.SetBit(MsgTypeBitsRegisterEVMAddress)

	case MsgSubmitProposal:
		mask.SetBit(MsgTypeBitsSubmitProposal)
	case MsgExecLegacyContent:
		mask.SetBit(MsgTypeBitsExecLegacyContent)
	case MsgVote:
		mask.SetBit(MsgTypeBitsVote)
	case MsgVoteWeighted:
		mask.SetBit(MsgTypeBitsVoteWeighted)
	case MsgDeposit:
		mask.SetBit(MsgTypeBitsDeposit)
	case IBCTransfer:
		mask.SetBit(MsgTypeBitsIBCTransfer)
	case MsgVerifyInvariant:
		mask.SetBit(MsgTypeBitsVerifyInvariant)
	case MsgSubmitEvidence:
		mask.SetBit(MsgTypeBitsSubmitEvidence)
	case MsgSendNFT:
		mask.SetBit(MsgTypeBitsSendNFT)

	case MsgCreateGroup:
		mask.SetBit(MsgTypeBitsCreateGroup)
	case MsgUpdateGroupMembers:
		mask.SetBit(MsgTypeBitsUpdateGroupMembers)
	case MsgUpdateGroupAdmin:
		mask.SetBit(MsgTypeBitsUpdateGroupAdmin)
	case MsgUpdateGroupMetadata:
		mask.SetBit(MsgTypeBitsUpdateGroupMetadata)
	case MsgCreateGroupPolicy:
		mask.SetBit(MsgTypeBitsCreateGroupPolicy)
	case MsgUpdateGroupPolicyAdmin:
		mask.SetBit(MsgTypeBitsUpdateGroupPolicyAdmin)
	case MsgCreateGroupWithPolicy:
		mask.SetBit(MsgTypeBitsCreateGroupWithPolicy)
	case MsgUpdateGroupPolicyDecisionPolicy:
		mask.SetBit(MsgTypeBitsUpdateGroupPolicyDecisionPolicy)
	case MsgUpdateGroupPolicyMetadata:
		mask.SetBit(MsgTypeBitsUpdateGroupPolicyMetadata)
	case MsgSubmitProposalGroup:
		mask.SetBit(MsgTypeBitsSubmitProposalGroup)
	case MsgWithdrawProposal:
		mask.SetBit(MsgTypeBitsWithdrawProposal)
	case MsgVoteGroup:
		mask.SetBit(MsgTypeBitsVoteGroup)
	case MsgExecGroup:
		mask.SetBit(MsgTypeBitsExecGroup)
	case MsgLeaveGroup:
		mask.SetBit(MsgTypeBitsLeaveGroup)

	case MsgSoftwareUpgrade:
		mask.SetBit(MsgTypeBitsSoftwareUpgrade)
	case MsgCancelUpgrade:
		mask.SetBit(MsgTypeBitsCancelUpgrade)
	case MsgRegisterInterchainAccount:
		mask.SetBit(MsgTypeBitsRegisterInterchainAccount)
	case MsgSendTx:
		mask.SetBit(MsgTypeBitsSendTx)
	case MsgRegisterPayee:
		mask.SetBit(MsgTypeBitsRegisterPayee)
	case MsgRegisterCounterpartyPayee:
		mask.SetBit(MsgTypeBitsRegisterCounterpartyPayee)
	case MsgPayPacketFee:
		mask.SetBit(MsgTypeBitsPayPacketFee)
	case MsgPayPacketFeeAsync:
		mask.SetBit(MsgTypeBitsPayPacketFeeAsync)
	case MsgTransfer:
		mask.SetBit(MsgTypeBitsTransfer)
	case MsgCreateClient:
		mask.SetBit(MsgTypeBitsCreateClient)
	case MsgUpdateClient:
		mask.SetBit(MsgTypeBitsUpdateClient)
	case MsgUpgradeClient:
		mask.SetBit(MsgTypeBitsUpgradeClient)
	case MsgSubmitMisbehaviour:
		mask.SetBit(MsgTypeBitsSubmitMisbehaviour)
	case MsgConnectionOpenInit:
		mask.SetBit(MsgTypeBitsConnectionOpenInit)
	case MsgConnectionOpenTry:
		mask.SetBit(MsgTypeBitsConnectionOpenTry)
	case MsgConnectionOpenAck:
		mask.SetBit(MsgTypeBitsConnectionOpenAck)
	case MsgConnectionOpenConfirm:
		mask.SetBit(MsgTypeBitsConnectionOpenConfirm)

	case MsgChannelOpenInit:
		mask.SetBit(MsgTypeBitsChannelOpenInit)
	case MsgChannelOpenTry:
		mask.SetBit(MsgTypeBitsChannelOpenTry)
	case MsgChannelOpenAck:
		mask.SetBit(MsgTypeBitsChannelOpenAck)
	case MsgChannelOpenConfirm:
		mask.SetBit(MsgTypeBitsChannelOpenConfirm)
	case MsgChannelCloseInit:
		mask.SetBit(MsgTypeBitsChannelCloseInit)
	case MsgChannelCloseConfirm:
		mask.SetBit(MsgTypeBitsChannelCloseConfirm)
	case MsgRecvPacket:
		mask.SetBit(MsgTypeBitsRecvPacket)
	case MsgTimeout:
		mask.SetBit(MsgTypeBitsTimeout)
	case MsgTimeoutOnClose:
		mask.SetBit(MsgTypeBitsTimeoutOnClose)
	case MsgAcknowledgement:
		mask.SetBit(MsgTypeBitsAcknowledgement)
	}
}

func (mask MsgTypeBits) Names() []MsgType {
	names := make([]MsgType, mask.CountBits())
	var i int

	if mask.HasBit(MsgTypeBitsUnknown) {
		names[i] = MsgUnknown
		i++
	}

	if mask.HasBit(MsgTypeBitsSetWithdrawAddress) {
		names[i] = MsgSetWithdrawAddress
	}
	if mask.HasBit(MsgTypeBitsWithdrawDelegatorReward) {
		names[i] = MsgWithdrawDelegatorReward
		i++
	}
	if mask.HasBit(MsgTypeBitsWithdrawValidatorCommission) {
		names[i] = MsgWithdrawValidatorCommission
		i++
	}
	if mask.HasBit(MsgTypeBitsFundCommunityPool) {
		names[i] = MsgFundCommunityPool
		i++
	}

	if mask.HasBit(MsgTypeBitsCreateValidator) {
		names[i] = MsgCreateValidator
		i++
	}
	if mask.HasBit(MsgTypeBitsEditValidator) {
		names[i] = MsgEditValidator
		i++
	}
	if mask.HasBit(MsgTypeBitsDelegate) {
		names[i] = MsgDelegate
		i++
	}
	if mask.HasBit(MsgTypeBitsBeginRedelegate) {
		names[i] = MsgBeginRedelegate
		i++
	}
	if mask.HasBit(MsgTypeBitsUndelegate) {
		names[i] = MsgUndelegate
		i++
	}
	if mask.HasBit(MsgTypeBitsCancelUnbondingDelegation) {
		names[i] = MsgCancelUnbondingDelegation
		i++
	}

	if mask.HasBit(MsgTypeBitsUnjail) {
		names[i] = MsgUnjail
		i++
	}

	if mask.HasBit(MsgTypeBitsSend) {
		names[i] = MsgSend
		i++
	}
	if mask.HasBit(MsgTypeBitsMultiSend) {
		names[i] = MsgMultiSend
		i++
	}

	if mask.HasBit(MsgTypeBitsCreateVestingAccount) {
		names[i] = MsgCreateVestingAccount
		i++
	}
	if mask.HasBit(MsgTypeBitsCreatePermanentLockedAccount) {
		names[i] = MsgCreatePermanentLockedAccount
		i++
	}
	if mask.HasBit(MsgTypeBitsCreatePeriodicVestingAccount) {
		names[i] = MsgCreatePeriodicVestingAccount
		i++
	}

	if mask.HasBit(MsgTypeBitsPayForBlobs) {
		names[i] = MsgPayForBlobs
		i++
	}

	if mask.HasBit(MsgTypeBitsGrant) {
		names[i] = MsgGrant
		i++
	}
	if mask.HasBit(MsgTypeBitsExec) {
		names[i] = MsgExec
		i++
	}
	if mask.HasBit(MsgTypeBitsRevoke) {
		names[i] = MsgRevoke
		i++
	}

	if mask.HasBit(MsgTypeBitsGrantAllowance) {
		names[i] = MsgGrantAllowance
		i++
	}
	if mask.HasBit(MsgTypeBitsRevokeAllowance) {
		names[i] = MsgRevokeAllowance
		i++
	}

	if mask.HasBit(MsgTypeBitsRegisterEVMAddress) {
		names[i] = MsgRegisterEVMAddress
		i++
	}

	if mask.HasBit(MsgTypeBitsSubmitProposal) {
		names[i] = MsgSubmitProposal
		i++
	}
	if mask.HasBit(MsgTypeBitsExecLegacyContent) {
		names[i] = MsgExecLegacyContent
		i++
	}
	if mask.HasBit(MsgTypeBitsVote) {
		names[i] = MsgVote
		i++
	}
	if mask.HasBit(MsgTypeBitsVoteWeighted) {
		names[i] = MsgVoteWeighted
		i++
	}
	if mask.HasBit(MsgTypeBitsDeposit) {
		names[i] = MsgDeposit
		i++
	}

	if mask.HasBit(MsgTypeBitsIBCTransfer) {
		names[i] = IBCTransfer
		i++
	}

	if mask.HasBit(MsgTypeBitsVerifyInvariant) {
		names[i] = MsgVerifyInvariant
		i++
	}

	if mask.HasBit(MsgTypeBitsSubmitEvidence) {
		names[i] = MsgSubmitEvidence
		i++
	}

	if mask.HasBit(MsgTypeBitsSendNFT) {
		names[i] = MsgSendNFT
		i++
	}

	if mask.HasBit(MsgTypeBitsCreateGroup) {
		names[i] = MsgCreateGroup
		i++
	}
	if mask.HasBit(MsgTypeBitsUpdateGroupMembers) {
		names[i] = MsgUpdateGroupMembers
		i++
	}
	if mask.HasBit(MsgTypeBitsUpdateGroupAdmin) {
		names[i] = MsgUpdateGroupAdmin
		i++
	}
	if mask.HasBit(MsgTypeBitsUpdateGroupMetadata) {
		names[i] = MsgUpdateGroupMetadata
		i++
	}
	if mask.HasBit(MsgTypeBitsCreateGroupPolicy) {
		names[i] = MsgCreateGroupPolicy
		i++
	}
	if mask.HasBit(MsgTypeBitsUpdateGroupPolicyAdmin) {
		names[i] = MsgUpdateGroupPolicyAdmin
		i++
	}
	if mask.HasBit(MsgTypeBitsCreateGroupWithPolicy) {
		names[i] = MsgCreateGroupWithPolicy
		i++
	}
	if mask.HasBit(MsgTypeBitsUpdateGroupPolicyDecisionPolicy) {
		names[i] = MsgUpdateGroupPolicyDecisionPolicy
		i++
	}
	if mask.HasBit(MsgTypeBitsUpdateGroupPolicyMetadata) {
		names[i] = MsgUpdateGroupPolicyMetadata
		i++
	}
	if mask.HasBit(MsgTypeBitsSubmitProposalGroup) {
		names[i] = MsgSubmitProposalGroup
		i++
	}
	if mask.HasBit(MsgTypeBitsWithdrawProposal) {
		names[i] = MsgWithdrawProposal
		i++
	}
	if mask.HasBit(MsgTypeBitsVoteGroup) {
		names[i] = MsgVoteGroup
		i++
	}
	if mask.HasBit(MsgTypeBitsExecGroup) {
		names[i] = MsgExecGroup
		i++
	}
	if mask.HasBit(MsgTypeBitsLeaveGroup) {
		names[i] = MsgLeaveGroup
		i++
	}

	if mask.HasBit(MsgTypeBitsSoftwareUpgrade) {
		names[i] = MsgSoftwareUpgrade
		i++
	}
	if mask.HasBit(MsgTypeBitsCancelUpgrade) {
		names[i] = MsgCancelUpgrade
		i++
	}
	if mask.HasBit(MsgTypeBitsRegisterInterchainAccount) {
		names[i] = MsgRegisterInterchainAccount
		i++
	}
	if mask.HasBit(MsgTypeBitsSendTx) {
		names[i] = MsgSendTx
		i++
	}
	if mask.HasBit(MsgTypeBitsRegisterPayee) {
		names[i] = MsgRegisterPayee
		i++
	}
	if mask.HasBit(MsgTypeBitsRegisterCounterpartyPayee) {
		names[i] = MsgRegisterCounterpartyPayee
		i++
	}
	if mask.HasBit(MsgTypeBitsPayPacketFee) {
		names[i] = MsgPayPacketFee
		i++
	}
	if mask.HasBit(MsgTypeBitsPayPacketFeeAsync) {
		names[i] = MsgPayPacketFeeAsync
		i++
	}
	if mask.HasBit(MsgTypeBitsTransfer) {
		names[i] = MsgTransfer
		i++
	}
	if mask.HasBit(MsgTypeBitsCreateClient) {
		names[i] = MsgCreateClient
		i++
	}
	if mask.HasBit(MsgTypeBitsUpdateClient) {
		names[i] = MsgUpdateClient
		i++
	}
	if mask.HasBit(MsgTypeBitsUpgradeClient) {
		names[i] = MsgUpgradeClient
		i++
	}
	if mask.HasBit(MsgTypeBitsSubmitMisbehaviour) {
		names[i] = MsgSubmitMisbehaviour
		i++
	}
	if mask.HasBit(MsgTypeBitsConnectionOpenInit) {
		names[i] = MsgConnectionOpenInit
		i++
	}
	if mask.HasBit(MsgTypeBitsConnectionOpenTry) {
		names[i] = MsgConnectionOpenTry
		i++
	}
	if mask.HasBit(MsgTypeBitsConnectionOpenAck) {
		names[i] = MsgConnectionOpenAck
		i++
	}
	if mask.HasBit(MsgTypeBitsConnectionOpenConfirm) {
		names[i] = MsgConnectionOpenConfirm
		i++
	}

	if mask.HasBit(MsgTypeBitsChannelOpenInit) {
		names[i] = MsgChannelOpenInit
		i++
	}
	if mask.HasBit(MsgTypeBitsChannelOpenTry) {
		names[i] = MsgChannelOpenTry
		i++
	}
	if mask.HasBit(MsgTypeBitsChannelOpenAck) {
		names[i] = MsgChannelOpenAck
		i++
	}
	if mask.HasBit(MsgTypeBitsChannelOpenConfirm) {
		names[i] = MsgChannelOpenConfirm
		i++
	}
	if mask.HasBit(MsgTypeBitsChannelCloseInit) {
		names[i] = MsgChannelCloseInit
		i++
	}
	if mask.HasBit(MsgTypeBitsChannelCloseConfirm) {
		names[i] = MsgChannelCloseConfirm
		i++
	}
	if mask.HasBit(MsgTypeBitsRecvPacket) {
		names[i] = MsgRecvPacket
		i++
	}
	if mask.HasBit(MsgTypeBitsTimeout) {
		names[i] = MsgTimeout
		i++
	}
	if mask.HasBit(MsgTypeBitsTimeoutOnClose) {
		names[i] = MsgTimeoutOnClose
		i++
	}
	if mask.HasBit(MsgTypeBitsAcknowledgement) {
		names[i] = MsgAcknowledgement
		// i++
	}

	return names
}

func (mask MsgTypeBits) HasOne(value MsgTypeBits) bool {
	return mask.value.And(mask.value, value.value).Cmp(zero) > 0
}

var _ sql.Scanner = (*MsgTypeBits)(nil)

func (mask *MsgTypeBits) Scan(src interface{}) (err error) {
	switch val := src.(type) {
	case []byte:
		mask.Bits, err = NewBitsFromString(string(val))
		if err != nil {
			return err
		}
	case nil:
		mask.Bits = NewEmptyBits()
	default:
		return errors.Errorf("unknown bits database type: %T", src)
	}
	return nil
}

var _ driver.Valuer = (*MsgTypeBits)(nil)

func (mask MsgTypeBits) Value() (driver.Value, error) {
	if mask.value == nil {
		return fmt.Sprintf("%074b", 0), nil
	}
	return fmt.Sprintf("%074b", mask.value), nil
}
