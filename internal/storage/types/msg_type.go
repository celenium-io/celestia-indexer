package types

// MsgType -
type MsgType string

// supported message types
const (
	MsgTypeUnknown                      MsgType = "MsgUnknown"
	MsgTypeWithdrawValidatorCommission  MsgType = "MsgWithdrawValidatorCommission"
	MsgTypeWithdrawDelegatorReward      MsgType = "MsgWithdrawDelegatorReward"
	MsgTypeEditValidator                MsgType = "MsgEditValidator"
	MsgTypeBeginRedelegate              MsgType = "MsgBeginRedelegate"
	MsgTypeCreateValidator              MsgType = "MsgCreateValidator"
	MsgTypeDelegate                     MsgType = "MsgDelegate"
	MsgTypeUndelegate                   MsgType = "MsgUndelegate"
	MsgTypeUnjail                       MsgType = "MsgUnjail"
	MsgTypeSend                         MsgType = "MsgSend"
	MsgTypeCreateVestingAccount         MsgType = "MsgCreateVestingAccount"
	MsgTypeCreatePeriodicVestingAccount MsgType = "MsgCreatePeriodicVestingAccount"
	MsgTypePayForBlobs                  MsgType = "MsgPayForBlobs"
	MsgTypeGrantAllowance               MsgType = "MsgGrantAllowance"
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
	string(MsgTypeGrantAllowance):               {},
}

func IsMsgType(val string) bool {
	_, ok := availiableMsgTypes[val]
	return ok
}
