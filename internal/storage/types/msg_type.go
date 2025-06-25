// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum MsgType
/*
	ENUM(
		MsgUnknown,

		MsgSetWithdrawAddress,
		MsgWithdrawDelegatorReward,
		MsgWithdrawValidatorCommission,
		MsgFundCommunityPool,

		MsgCreateValidator,
		MsgEditValidator,
		MsgDelegate,
		MsgBeginRedelegate,
		MsgUndelegate,
		MsgCancelUnbondingDelegation,

		MsgUnjail,

		MsgSend,
		MsgMultiSend,

		MsgCreateVestingAccount,
		MsgCreatePermanentLockedAccount,
		MsgCreatePeriodicVestingAccount,

		MsgPayForBlobs,

		MsgGrant,
		MsgExec,
		MsgRevoke,

		MsgGrantAllowance,
		MsgRevokeAllowance,

		MsgRegisterEVMAddress,

		MsgSubmitProposal,
		MsgExecLegacyContent,
		MsgVote,
		MsgVoteWeighted,
		MsgDeposit,

		IBCTransfer,

		MsgVerifyInvariant,

		MsgSubmitEvidence,

		MsgSendNFT,

		MsgCreateGroup,
		MsgUpdateGroupMembers,
		MsgUpdateGroupAdmin,
		MsgUpdateGroupMetadata,
		MsgCreateGroupPolicy,
		MsgUpdateGroupPolicyAdmin,
		MsgCreateGroupWithPolicy,
		MsgUpdateGroupPolicyDecisionPolicy,
		MsgUpdateGroupPolicyMetadata,
		MsgSubmitProposalGroup,
		MsgWithdrawProposal,
		MsgVoteGroup,
		MsgExecGroup,
		MsgLeaveGroup,

		MsgSoftwareUpgrade,
		MsgCancelUpgrade,

		MsgRegisterInterchainAccount,
		MsgSendTx,

		MsgRegisterPayee,
		MsgRegisterCounterpartyPayee,
		MsgPayPacketFee,
		MsgPayPacketFeeAsync,

		MsgTransfer,

		MsgCreateClient,
		MsgUpdateClient,
		MsgUpgradeClient,
		MsgSubmitMisbehaviour,
		MsgRecoverClient,
		MsgIBCSoftwareUpgrade,
		MsgUpdateParams,

		MsgConnectionOpenInit,
		MsgConnectionOpenTry,
		MsgConnectionOpenAck,
		MsgConnectionOpenConfirm,

		MsgChannelOpenInit,
		MsgChannelOpenTry,
		MsgChannelOpenAck,
		MsgChannelOpenConfirm,
		MsgChannelCloseInit,
		MsgChannelCloseConfirm,
		MsgRecvPacket,
		MsgTimeout,
		MsgTimeoutOnClose,
		MsgAcknowledgement,

		MsgSignalVersion,
		MsgTryUpgrade,

		MsgCreateMailbox,
		MsgProcessMessage,
		MsgSetMailbox,
		MsgCreateCollateralToken,
		MsgCreateSyntheticToken,
		MsgSetToken,
		MsgEnrollRemoteRouter,
		MsgUnrollRemoteRouter,
		MsgRemoteTransfer,

		MsgUpdateMinfeeParams
	)
*/
//go:generate go-enum --marshal --sql --values --noprefix --names
type MsgType string
