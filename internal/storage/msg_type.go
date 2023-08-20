package storage

// MsgType -
type MsgType string

// supported message types
const (
	MsgTypeUnknown                      MsgType = "Unknown"
	MsgTypeWithdrawValidatorCommission  MsgType = "WithdrawValidatorCommission"
	MsgTypeWithdrawDelegatorReward      MsgType = "WithdrawDelegatorReward"
	MsgTypeEditValidator                MsgType = "EditValidator"
	MsgTypeBeginRedelegate              MsgType = "BeginRedelegate"
	MsgTypeCreateValidator              MsgType = "CreateValidator"
	MsgTypeDelegate                     MsgType = "Delegate"
	MsgTypeUndelegate                   MsgType = "Undelegate"
	MsgTypeUnjail                       MsgType = "Unjail"
	MsgTypeSend                         MsgType = "Send"
	MsgTypeCreateVestingAccount         MsgType = "CreateVestingAccount"
	MsgTypeCreatePeriodicVestingAccount MsgType = "CreatePeriodicVestingAccount"
	MsgTypePayForBlobs                  MsgType = "PayForBlobs"
)
