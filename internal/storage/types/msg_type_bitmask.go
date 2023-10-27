// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

type MsgTypeBits struct {
	Bits
}

const (
	MsgTypeBitsUnknown uint64 = 1 << iota

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

	// MsgTypeBitsChannelOpenInit
	// MsgTypeBitsChannelOpenTry
	// MsgTypeBitsChannelOpenAck
	// MsgTypeBitsChannelOpenConfirm
	// MsgTypeBitsChannelCloseInit
	// MsgTypeBitsChannelCloseConfirm
	// MsgTypeBitsRecvPacket
	// MsgTypeBitsTimeout
	// MsgTypeBitsTimeoutOnClose
	// MsgTypeBitsAcknowledgement
)

func NewMsgTypeBitMask(values ...MsgType) MsgTypeBits {
	var mask MsgTypeBits
	for i := range values {
		mask.SetBit(values[i])
	}
	return mask
}

func (mask *MsgTypeBits) SetBit(value MsgType) {
	switch value {
	case MsgUnknown:
		mask.Set(Bits(MsgTypeBitsUnknown))

	case MsgSetWithdrawAddress:
		mask.Set(Bits(MsgTypeBitsSetWithdrawAddress))
	case MsgWithdrawDelegatorReward:
		mask.Set(Bits(MsgTypeBitsWithdrawDelegatorReward))
	case MsgWithdrawValidatorCommission:
		mask.Set(Bits(MsgTypeBitsWithdrawValidatorCommission))
	case MsgFundCommunityPool:
		mask.Set(Bits(MsgTypeBitsFundCommunityPool))

	case MsgCreateValidator:
		mask.Set(Bits(MsgTypeBitsCreateValidator))
	case MsgEditValidator:
		mask.Set(Bits(MsgTypeBitsEditValidator))
	case MsgDelegate:
		mask.Set(Bits(MsgTypeBitsDelegate))
	case MsgBeginRedelegate:
		mask.Set(Bits(MsgTypeBitsBeginRedelegate))
	case MsgUndelegate:
		mask.Set(Bits(MsgTypeBitsUndelegate))
	case MsgCancelUnbondingDelegation:
		mask.Set(Bits(MsgTypeBitsCancelUnbondingDelegation))

	case MsgUnjail:
		mask.Set(Bits(MsgTypeBitsUnjail))

	case MsgSend:
		mask.Set(Bits(MsgTypeBitsSend))
	case MsgMultiSend:
		mask.Set(Bits(MsgTypeBitsMultiSend))

	case MsgCreateVestingAccount:
		mask.Set(Bits(MsgTypeBitsCreateVestingAccount))
	case MsgCreatePermanentLockedAccount:
		mask.Set(Bits(MsgTypeBitsCreatePermanentLockedAccount))
	case MsgCreatePeriodicVestingAccount:
		mask.Set(Bits(MsgTypeBitsCreatePeriodicVestingAccount))

	case MsgPayForBlobs:
		mask.Set(Bits(MsgTypeBitsPayForBlobs))

	case MsgGrant:
		mask.Set(Bits(MsgTypeBitsGrant))
	case MsgExec:
		mask.Set(Bits(MsgTypeBitsExec))
	case MsgRevoke:
		mask.Set(Bits(MsgTypeBitsRevoke))

	case MsgGrantAllowance:
		mask.Set(Bits(MsgTypeBitsGrantAllowance))
	case MsgRevokeAllowance:
		mask.Set(Bits(MsgTypeBitsRevokeAllowance))

	case MsgRegisterEVMAddress:
		mask.Set(Bits(MsgTypeBitsRegisterEVMAddress))

	case MsgSubmitProposal:
		mask.Set(Bits(MsgTypeBitsSubmitProposal))
	case MsgExecLegacyContent:
		mask.Set(Bits(MsgTypeBitsExecLegacyContent))
	case MsgVote:
		mask.Set(Bits(MsgTypeBitsVote))
	case MsgVoteWeighted:
		mask.Set(Bits(MsgTypeBitsVoteWeighted))
	case MsgDeposit:
		mask.Set(Bits(MsgTypeBitsDeposit))
	case IBCTransfer:
		mask.Set(Bits(MsgTypeBitsIBCTransfer))
	case MsgVerifyInvariant:
		mask.Set(Bits(MsgTypeBitsVerifyInvariant))
	case MsgSubmitEvidence:
		mask.Set(Bits(MsgTypeBitsSubmitEvidence))
	case MsgSendNFT:
		mask.Set(Bits(MsgTypeBitsSendNFT))

	case MsgCreateGroup:
		mask.Set(Bits(MsgTypeBitsCreateGroup))
	case MsgUpdateGroupMembers:
		mask.Set(Bits(MsgTypeBitsUpdateGroupMembers))
	case MsgUpdateGroupAdmin:
		mask.Set(Bits(MsgTypeBitsUpdateGroupAdmin))
	case MsgUpdateGroupMetadata:
		mask.Set(Bits(MsgTypeBitsUpdateGroupMetadata))
	case MsgCreateGroupPolicy:
		mask.Set(Bits(MsgTypeBitsCreateGroupPolicy))
	case MsgUpdateGroupPolicyAdmin:
		mask.Set(Bits(MsgTypeBitsUpdateGroupPolicyAdmin))
	case MsgCreateGroupWithPolicy:
		mask.Set(Bits(MsgTypeBitsCreateGroupWithPolicy))
	case MsgUpdateGroupPolicyDecisionPolicy:
		mask.Set(Bits(MsgTypeBitsUpdateGroupPolicyDecisionPolicy))
	case MsgUpdateGroupPolicyMetadata:
		mask.Set(Bits(MsgTypeBitsUpdateGroupPolicyMetadata))
	case MsgSubmitProposalGroup:
		mask.Set(Bits(MsgTypeBitsSubmitProposalGroup))
	case MsgWithdrawProposal:
		mask.Set(Bits(MsgTypeBitsWithdrawProposal))
	case MsgVoteGroup:
		mask.Set(Bits(MsgTypeBitsVoteGroup))
	case MsgExecGroup:
		mask.Set(Bits(MsgTypeBitsExecGroup))
	case MsgLeaveGroup:
		mask.Set(Bits(MsgTypeBitsLeaveGroup))

	case MsgSoftwareUpgrade:
		mask.Set(Bits(MsgTypeBitsSoftwareUpgrade))
	case MsgCancelUpgrade:
		mask.Set(Bits(MsgTypeBitsCancelUpgrade))
	case MsgRegisterInterchainAccount:
		mask.Set(Bits(MsgTypeBitsRegisterInterchainAccount))
	case MsgSendTx:
		mask.Set(Bits(MsgTypeBitsSendTx))
	case MsgRegisterPayee:
		mask.Set(Bits(MsgTypeBitsRegisterPayee))
	case MsgRegisterCounterpartyPayee:
		mask.Set(Bits(MsgTypeBitsRegisterCounterpartyPayee))
	case MsgPayPacketFee:
		mask.Set(Bits(MsgTypeBitsPayPacketFee))
	case MsgPayPacketFeeAsync:
		mask.Set(Bits(MsgTypeBitsPayPacketFeeAsync))
	case MsgTransfer:
		mask.Set(Bits(MsgTypeBitsTransfer))
	case MsgCreateClient:
		mask.Set(Bits(MsgTypeBitsCreateClient))
	case MsgUpdateClient:
		mask.Set(Bits(MsgTypeBitsUpdateClient))
	case MsgUpgradeClient:
		mask.Set(Bits(MsgTypeBitsUpgradeClient))
	case MsgSubmitMisbehaviour:
		mask.Set(Bits(MsgTypeBitsSubmitMisbehaviour))
	case MsgConnectionOpenInit:
		mask.Set(Bits(MsgTypeBitsConnectionOpenInit))
	case MsgConnectionOpenTry:
		mask.Set(Bits(MsgTypeBitsConnectionOpenTry))
	case MsgConnectionOpenAck:
		mask.Set(Bits(MsgTypeBitsConnectionOpenAck))
	case MsgConnectionOpenConfirm:
		mask.Set(Bits(MsgTypeBitsConnectionOpenConfirm))

		// case MsgChannelOpenInit:
		// 	mask.Set(Bits(MsgTypeBitsChannelOpenInit))
		// case MsgChannelOpenTry:
		// 	mask.Set(Bits(MsgTypeBitsChannelOpenTry))
		// case MsgChannelOpenAck:
		// 	mask.Set(Bits(MsgTypeBitsChannelOpenAck))
		// case MsgChannelOpenConfirm:
		// 	mask.Set(Bits(MsgTypeBitsChannelOpenConfirm))
		// case MsgChannelCloseInit:
		// 	mask.Set(Bits(MsgTypeBitsChannelCloseInit))
		// case MsgChannelCloseConfirm:
		// 	mask.Set(Bits(MsgTypeBitsChannelCloseConfirm))
		// case MsgRecvPacket:
		// 	mask.Set(Bits(MsgTypeBitsRecvPacket))
		// case MsgTimeout:
		// 	mask.Set(Bits(MsgTypeBitsTimeout))
		// case MsgTimeoutOnClose:
		// 	mask.Set(Bits(MsgTypeBitsTimeoutOnClose))
		// case MsgAcknowledgement:
		// 	mask.Set(Bits(MsgTypeBitsAcknowledgement))
	}
}

func (mask MsgTypeBits) Names() []MsgType {
	names := make([]MsgType, mask.CountBits())
	var i int

	if mask.Has(Bits(MsgTypeBitsUnknown)) {
		names[i] = MsgUnknown
		i++
	}

	if mask.Has(Bits(MsgTypeBitsSetWithdrawAddress)) {
		names[i] = MsgSetWithdrawAddress
	}
	if mask.Has(Bits(MsgTypeBitsWithdrawDelegatorReward)) {
		names[i] = MsgWithdrawDelegatorReward
		i++
	}
	if mask.Has(Bits(MsgTypeBitsWithdrawValidatorCommission)) {
		names[i] = MsgWithdrawValidatorCommission
		i++
	}
	if mask.Has(Bits(MsgTypeBitsFundCommunityPool)) {
		names[i] = MsgFundCommunityPool
		i++
	}

	if mask.Has(Bits(MsgTypeBitsCreateValidator)) {
		names[i] = MsgCreateValidator
		i++
	}
	if mask.Has(Bits(MsgTypeBitsEditValidator)) {
		names[i] = MsgEditValidator
		i++
	}
	if mask.Has(Bits(MsgTypeBitsDelegate)) {
		names[i] = MsgDelegate
		i++
	}
	if mask.Has(Bits(MsgTypeBitsBeginRedelegate)) {
		names[i] = MsgBeginRedelegate
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUndelegate)) {
		names[i] = MsgUndelegate
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCancelUnbondingDelegation)) {
		names[i] = MsgCancelUnbondingDelegation
		i++
	}

	if mask.Has(Bits(MsgTypeBitsUnjail)) {
		names[i] = MsgUnjail
		i++
	}

	if mask.Has(Bits(MsgTypeBitsSend)) {
		names[i] = MsgSend
		i++
	}
	if mask.Has(Bits(MsgTypeBitsMultiSend)) {
		names[i] = MsgMultiSend
		i++
	}

	if mask.Has(Bits(MsgTypeBitsCreateVestingAccount)) {
		names[i] = MsgCreateVestingAccount
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreatePermanentLockedAccount)) {
		names[i] = MsgCreatePermanentLockedAccount
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreatePeriodicVestingAccount)) {
		names[i] = MsgCreatePeriodicVestingAccount
		i++
	}

	if mask.Has(Bits(MsgTypeBitsPayForBlobs)) {
		names[i] = MsgPayForBlobs
		i++
	}

	if mask.Has(Bits(MsgTypeBitsGrant)) {
		names[i] = MsgGrant
		i++
	}
	if mask.Has(Bits(MsgTypeBitsExec)) {
		names[i] = MsgExec
		i++
	}
	if mask.Has(Bits(MsgTypeBitsRevoke)) {
		names[i] = MsgRevoke
		i++
	}

	if mask.Has(Bits(MsgTypeBitsGrantAllowance)) {
		names[i] = MsgGrantAllowance
		i++
	}
	if mask.Has(Bits(MsgTypeBitsRevokeAllowance)) {
		names[i] = MsgRevokeAllowance
		i++
	}

	if mask.Has(Bits(MsgTypeBitsRegisterEVMAddress)) {
		names[i] = MsgRegisterEVMAddress
		i++
	}

	if mask.Has(Bits(MsgTypeBitsSubmitProposal)) {
		names[i] = MsgSubmitProposal
		i++
	}
	if mask.Has(Bits(MsgTypeBitsExecLegacyContent)) {
		names[i] = MsgExecLegacyContent
		i++
	}
	if mask.Has(Bits(MsgTypeBitsVote)) {
		names[i] = MsgVote
		i++
	}
	if mask.Has(Bits(MsgTypeBitsVoteWeighted)) {
		names[i] = MsgVoteWeighted
		i++
	}
	if mask.Has(Bits(MsgTypeBitsDeposit)) {
		names[i] = MsgDeposit
		i++
	}

	if mask.Has(Bits(MsgTypeBitsIBCTransfer)) {
		names[i] = IBCTransfer
		i++
	}

	if mask.Has(Bits(MsgTypeBitsVerifyInvariant)) {
		names[i] = MsgVerifyInvariant
		i++
	}

	if mask.Has(Bits(MsgTypeBitsSubmitEvidence)) {
		names[i] = MsgSubmitEvidence
		i++
	}

	if mask.Has(Bits(MsgTypeBitsSendNFT)) {
		names[i] = MsgSendNFT
		i++
	}

	if mask.Has(Bits(MsgTypeBitsCreateGroup)) {
		names[i] = MsgCreateGroup
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUpdateGroupMembers)) {
		names[i] = MsgUpdateGroupMembers
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUpdateGroupAdmin)) {
		names[i] = MsgUpdateGroupAdmin
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUpdateGroupMetadata)) {
		names[i] = MsgUpdateGroupMetadata
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreateGroupPolicy)) {
		names[i] = MsgCreateGroupPolicy
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUpdateGroupPolicyAdmin)) {
		names[i] = MsgUpdateGroupPolicyAdmin
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreateGroupWithPolicy)) {
		names[i] = MsgCreateGroupWithPolicy
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUpdateGroupPolicyDecisionPolicy)) {
		names[i] = MsgUpdateGroupPolicyDecisionPolicy
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUpdateGroupPolicyMetadata)) {
		names[i] = MsgUpdateGroupPolicyMetadata
		i++
	}
	if mask.Has(Bits(MsgTypeBitsSubmitProposalGroup)) {
		names[i] = MsgSubmitProposalGroup
		i++
	}
	if mask.Has(Bits(MsgTypeBitsWithdrawProposal)) {
		names[i] = MsgWithdrawProposal
		i++
	}
	if mask.Has(Bits(MsgTypeBitsVoteGroup)) {
		names[i] = MsgVoteGroup
		i++
	}
	if mask.Has(Bits(MsgTypeBitsExecGroup)) {
		names[i] = MsgExecGroup
		i++
	}
	if mask.Has(Bits(MsgTypeBitsLeaveGroup)) {
		names[i] = MsgLeaveGroup
		i++
	}

	if mask.Has(Bits(MsgTypeBitsSoftwareUpgrade)) {
		names[i] = MsgSoftwareUpgrade
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCancelUpgrade)) {
		names[i] = MsgCancelUpgrade
		i++
	}
	if mask.Has(Bits(MsgTypeBitsRegisterInterchainAccount)) {
		names[i] = MsgRegisterInterchainAccount
		i++
	}
	if mask.Has(Bits(MsgTypeBitsSendTx)) {
		names[i] = MsgSendTx
		i++
	}
	if mask.Has(Bits(MsgTypeBitsRegisterPayee)) {
		names[i] = MsgRegisterPayee
		i++
	}
	if mask.Has(Bits(MsgTypeBitsRegisterCounterpartyPayee)) {
		names[i] = MsgRegisterCounterpartyPayee
		i++
	}
	if mask.Has(Bits(MsgTypeBitsPayPacketFee)) {
		names[i] = MsgPayPacketFee
		i++
	}
	if mask.Has(Bits(MsgTypeBitsPayPacketFeeAsync)) {
		names[i] = MsgPayPacketFeeAsync
		i++
	}
	if mask.Has(Bits(MsgTypeBitsTransfer)) {
		names[i] = MsgTransfer
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreateClient)) {
		names[i] = MsgCreateClient
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUpdateClient)) {
		names[i] = MsgUpdateClient
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUpgradeClient)) {
		names[i] = MsgUpgradeClient
		i++
	}
	if mask.Has(Bits(MsgTypeBitsSubmitMisbehaviour)) {
		names[i] = MsgSubmitMisbehaviour
		i++
	}
	if mask.Has(Bits(MsgTypeBitsConnectionOpenInit)) {
		names[i] = MsgConnectionOpenInit
		i++
	}
	if mask.Has(Bits(MsgTypeBitsConnectionOpenTry)) {
		names[i] = MsgConnectionOpenTry
		i++
	}
	if mask.Has(Bits(MsgTypeBitsConnectionOpenAck)) {
		names[i] = MsgConnectionOpenAck
		i++
	}
	if mask.Has(Bits(MsgTypeBitsConnectionOpenConfirm)) {
		names[i] = MsgConnectionOpenConfirm
		// i++
	}

	// if mask.Has(Bits(MsgTypeBitsChannelOpenInit)) {
	// 	names[i] = MsgChannelOpenInit
	// 	i++
	// }
	// if mask.Has(Bits(MsgTypeBitsChannelOpenTry)) {
	// 	names[i] = MsgChannelOpenTry
	// 	i++
	// }
	// if mask.Has(Bits(MsgTypeBitsChannelOpenAck)) {
	// 	names[i] = MsgChannelOpenAck
	// 	i++
	// }
	// if mask.Has(Bits(MsgTypeBitsChannelOpenConfirm)) {
	// 	names[i] = MsgChannelOpenConfirm
	// 	i++
	// }
	// if mask.Has(Bits(MsgTypeBitsChannelCloseInit)) {
	// 	names[i] = MsgChannelCloseInit
	// 	i++
	// }
	// if mask.Has(Bits(MsgTypeBitsChannelCloseConfirm)) {
	// 	names[i] = MsgChannelCloseConfirm
	// 	i++
	// }
	// if mask.Has(Bits(MsgTypeBitsRecvPacket)) {
	// 	names[i] = MsgRecvPacket
	// 	i++
	// }
	// if mask.Has(Bits(MsgTypeBitsTimeout)) {
	// 	names[i] = MsgTimeout
	// 	i++
	// }
	// if mask.Has(Bits(MsgTypeBitsTimeoutOnClose)) {
	// 	names[i] = MsgTimeoutOnClose
	// 	i++
	// }
	// if mask.Has(Bits(MsgTypeBitsAcknowledgement)) {
	// 	names[i] = MsgAcknowledgement
	// 	// i++
	// }

	return names
}

func (mask MsgTypeBits) HasOne(value MsgTypeBits) bool {
	return mask.Bits&value.Bits > 0
}

var _ sql.Scanner = (*MsgTypeBits)(nil)

func (mask *MsgTypeBits) Scan(src interface{}) (err error) {
	switch val := src.(type) {
	case int64:
		mask.Bits = Bits(val)
	case nil:
		mask.Bits = 0
	default:
		return errors.Errorf("unknown bits database type: %T", src)
	}
	return nil
}

var _ driver.Valuer = (*MsgTypeBits)(nil)

func (mask MsgTypeBits) Value() (driver.Value, error) {
	return uint64(mask.Bits), nil
}
