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
	MsgTypeBitsWithdrawValidatorCommission
	MsgTypeBitsWithdrawDelegatorReward
	MsgTypeBitsEditValidator
	MsgTypeBitsBeginRedelegate
	MsgTypeBitsCreateValidator
	MsgTypeBitsDelegate
	MsgTypeBitsUndelegate
	MsgTypeBitsUnjail
	MsgTypeBitsSend
	MsgTypeBitsCreateVestingAccount
	MsgTypeBitsCreatePeriodicVestingAccount
	MsgTypeBitsPayForBlobs
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
	case MsgWithdrawValidatorCommission:
		mask.Set(Bits(MsgTypeBitsWithdrawValidatorCommission))
	case MsgWithdrawDelegatorReward:
		mask.Set(Bits(MsgTypeBitsWithdrawDelegatorReward))
	case MsgEditValidator:
		mask.Set(Bits(MsgTypeBitsEditValidator))
	case MsgBeginRedelegate:
		mask.Set(Bits(MsgTypeBitsBeginRedelegate))
	case MsgCreateValidator:
		mask.Set(Bits(MsgTypeBitsCreateValidator))
	case MsgUndelegate:
		mask.Set(Bits(MsgTypeBitsUndelegate))
	case MsgUnjail:
		mask.Set(Bits(MsgTypeBitsUnjail))
	case MsgSend:
		mask.Set(Bits(MsgTypeBitsSend))
	case MsgCreateVestingAccount:
		mask.Set(Bits(MsgTypeBitsCreateVestingAccount))
	case MsgCreatePeriodicVestingAccount:
		mask.Set(Bits(MsgTypeBitsCreatePeriodicVestingAccount))
	case MsgPayForBlobs:
		mask.Set(Bits(MsgTypeBitsPayForBlobs))
	case MsgDelegate:
		mask.Set(Bits(MsgTypeBitsDelegate))
	}
}

func (mask MsgTypeBits) Names() []MsgType {
	names := make([]MsgType, mask.CountBits())
	var i int

	if mask.Has(Bits(MsgTypeBitsUnknown)) {
		names[i] = MsgUnknown
		i++
	}
	if mask.Has(Bits(MsgTypeBitsDelegate)) {
		names[i] = MsgDelegate
		i++
	}
	if mask.Has(Bits(MsgTypeBitsWithdrawValidatorCommission)) {
		names[i] = MsgWithdrawValidatorCommission
		i++
	}
	if mask.Has(Bits(MsgTypeBitsWithdrawDelegatorReward)) {
		names[i] = MsgWithdrawDelegatorReward
		i++
	}
	if mask.Has(Bits(MsgTypeBitsEditValidator)) {
		names[i] = MsgEditValidator
		i++
	}
	if mask.Has(Bits(MsgTypeBitsBeginRedelegate)) {
		names[i] = MsgBeginRedelegate
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreateValidator)) {
		names[i] = MsgCreateValidator
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUndelegate)) {
		names[i] = MsgUndelegate
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
	if mask.Has(Bits(MsgTypeBitsCreateVestingAccount)) {
		names[i] = MsgCreateVestingAccount
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreatePeriodicVestingAccount)) {
		names[i] = MsgCreatePeriodicVestingAccount
		i++
	}
	if mask.Has(Bits(MsgTypeBitsPayForBlobs)) {
		names[i] = MsgPayForBlobs
	}

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
