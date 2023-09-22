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
	)
*/
//go:generate go-enum --marshal --sql --values --noprefix
type MsgType string
