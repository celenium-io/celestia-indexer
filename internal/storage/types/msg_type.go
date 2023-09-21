package types

// swagger:enum MsgType
/*
	ENUM(
		MsgUnknown,
		MsgWithdrawValidatorCommission,
		MsgWithdrawDelegatorReward,
		MsgEditValidator,
		MsgBeginRedelegate,
		MsgCreateValidator,
		MsgDelegate,
		MsgUndelegate,
		MsgUnjail,

		MsgSend,
		MsgMultiSend,

		MsgCreateVestingAccount,
		MsgCreatePeriodicVestingAccount,
		MsgPayForBlobs,

		MsgGrant,
		MsgExec,
		MsgRevoke,

		MsgGrantAllowance,

		MsgRegisterEVMAddress,
		MsgSetWithdrawAddress,

		MsgVote,
		MsgVoteWeighted,
		MsgSubmitProposal,
	)
*/
//go:generate go-enum --marshal --sql --values --noprefix
type MsgType string
