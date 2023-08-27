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
		switch values[i] {
		case MsgTypeUnknown:
			mask.Set(Bits(MsgTypeBitsUnknown))
		case MsgTypeWithdrawValidatorCommission:
			mask.Set(Bits(MsgTypeBitsWithdrawValidatorCommission))
		case MsgTypeWithdrawDelegatorReward:
			mask.Set(Bits(MsgTypeBitsWithdrawDelegatorReward))
		case MsgTypeEditValidator:
			mask.Set(Bits(MsgTypeBitsEditValidator))
		case MsgTypeBeginRedelegate:
			mask.Set(Bits(MsgTypeBitsBeginRedelegate))
		case MsgTypeCreateValidator:
			mask.Set(Bits(MsgTypeBitsCreateValidator))
		case MsgTypeUndelegate:
			mask.Set(Bits(MsgTypeBitsUndelegate))
		case MsgTypeUnjail:
			mask.Set(Bits(MsgTypeBitsUnjail))
		case MsgTypeSend:
			mask.Set(Bits(MsgTypeBitsSend))
		case MsgTypeCreateVestingAccount:
			mask.Set(Bits(MsgTypeBitsCreateVestingAccount))
		case MsgTypeCreatePeriodicVestingAccount:
			mask.Set(Bits(MsgTypeBitsCreatePeriodicVestingAccount))
		case MsgTypePayForBlobs:
			mask.Set(Bits(MsgTypeBitsPayForBlobs))
		case MsgTypeDelegate:
			mask.Set(Bits(MsgTypeBitsDelegate))
		}
	}
	return mask
}

func (mask MsgTypeBits) Names() []string {
	names := make([]string, mask.CountBits())
	var i int

	if mask.Has(Bits(MsgTypeBitsUnknown)) {
		names[i] = string(MsgTypeUnknown)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsDelegate)) {
		names[i] = string(MsgTypeDelegate)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsWithdrawValidatorCommission)) {
		names[i] = string(MsgTypeWithdrawValidatorCommission)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsWithdrawDelegatorReward)) {
		names[i] = string(MsgTypeWithdrawDelegatorReward)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsEditValidator)) {
		names[i] = string(MsgTypeEditValidator)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsBeginRedelegate)) {
		names[i] = string(MsgTypeBeginRedelegate)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreateValidator)) {
		names[i] = string(MsgTypeCreateValidator)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUndelegate)) {
		names[i] = string(MsgTypeUndelegate)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsUnjail)) {
		names[i] = string(MsgTypeUnjail)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsSend)) {
		names[i] = string(MsgTypeSend)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreateVestingAccount)) {
		names[i] = string(MsgTypeCreateVestingAccount)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsCreatePeriodicVestingAccount)) {
		names[i] = string(MsgTypeCreatePeriodicVestingAccount)
		i++
	}
	if mask.Has(Bits(MsgTypeBitsPayForBlobs)) {
		names[i] = string(MsgTypePayForBlobs)
	}

	return names
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
