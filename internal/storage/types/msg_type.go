package types

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

var availiableMsgTypes = map[string]struct{}{
	string(MsgTypeUnknown):                      {},
	string(MsgTypeWithdrawValidatorCommission):  {},
	string(MsgTypeWithdrawDelegatorReward):      {},
	string(MsgTypeEditValidator):                {},
	string(MsgTypeBeginRedelegate):              {},
	string(MsgTypeCreateValidator):              {},
	string(MsgTypeDelegate):                     {},
	string(MsgTypeUndelegate):                   {},
	string(MsgTypeUnjail):                       {},
	string(MsgTypeSend):                         {},
	string(MsgTypeCreateVestingAccount):         {},
	string(MsgTypeCreatePeriodicVestingAccount): {},
	string(MsgTypePayForBlobs):                  {},
}

func IsMsgType(val string) bool {
	_, ok := availiableMsgTypes[val]
	return ok
}
