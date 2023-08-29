package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgTypeBits_Names(t *testing.T) {
	tests := []struct {
		name string
		Bits Bits
		want []string
	}{
		{
			name: string(MsgTypeBeginRedelegate),
			Bits: Bits(MsgTypeBitsBeginRedelegate),
			want: []string{string(MsgTypeBeginRedelegate)},
		}, {
			name: string(MsgTypeWithdrawValidatorCommission),
			Bits: Bits(MsgTypeBitsWithdrawValidatorCommission),
			want: []string{string(MsgTypeWithdrawValidatorCommission)},
		}, {
			name: string(MsgTypeCreatePeriodicVestingAccount),
			Bits: Bits(MsgTypeBitsCreatePeriodicVestingAccount),
			want: []string{string(MsgTypeCreatePeriodicVestingAccount)},
		}, {
			name: string(MsgTypeCreateValidator),
			Bits: Bits(MsgTypeBitsCreateValidator),
			want: []string{string(MsgTypeCreateValidator)},
		}, {
			name: string(MsgTypeCreateVestingAccount),
			Bits: Bits(MsgTypeBitsCreateVestingAccount),
			want: []string{string(MsgTypeCreateVestingAccount)},
		}, {
			name: string(MsgTypeDelegate),
			Bits: Bits(MsgTypeBitsDelegate),
			want: []string{string(MsgTypeDelegate)},
		}, {
			name: string(MsgTypeEditValidator),
			Bits: Bits(MsgTypeBitsEditValidator),
			want: []string{string(MsgTypeEditValidator)},
		}, {
			name: string(MsgTypePayForBlobs),
			Bits: Bits(MsgTypeBitsPayForBlobs),
			want: []string{string(MsgTypePayForBlobs)},
		}, {
			name: string(MsgTypeSend),
			Bits: Bits(MsgTypeBitsSend),
			want: []string{string(MsgTypeSend)},
		}, {
			name: string(MsgTypeUndelegate),
			Bits: Bits(MsgTypeBitsUndelegate),
			want: []string{string(MsgTypeUndelegate)},
		}, {
			name: string(MsgTypeUnjail),
			Bits: Bits(MsgTypeBitsUnjail),
			want: []string{string(MsgTypeUnjail)},
		}, {
			name: string(MsgTypeUnknown),
			Bits: Bits(MsgTypeBitsUnknown),
			want: []string{string(MsgTypeUnknown)},
		}, {
			name: string(MsgTypeWithdrawDelegatorReward),
			Bits: Bits(MsgTypeBitsWithdrawDelegatorReward),
			want: []string{string(MsgTypeWithdrawDelegatorReward)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mask := MsgTypeBits{
				Bits: tt.Bits,
			}
			require.Equal(t, tt.want, mask.Names())
		})
	}
}

func TestNewMsgTypeBitMask(t *testing.T) {
	tests := []struct {
		name   string
		values []MsgType
		want   MsgTypeBits
	}{
		{
			name:   "test 0",
			values: nil,
			want:   MsgTypeBits{Bits(0)},
		}, {
			name:   "test 1",
			values: []MsgType{MsgTypeBeginRedelegate},
			want:   MsgTypeBits{Bits(MsgTypeBitsBeginRedelegate)},
		}, {
			name:   "test 2",
			values: []MsgType{MsgTypeWithdrawValidatorCommission},
			want:   MsgTypeBits{Bits(MsgTypeBitsWithdrawValidatorCommission)},
		}, {
			name:   "test 3",
			values: []MsgType{MsgTypeCreatePeriodicVestingAccount},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreatePeriodicVestingAccount)},
		}, {
			name:   "test 4",
			values: []MsgType{MsgTypeCreateValidator},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateValidator)},
		}, {
			name:   "test 5",
			values: []MsgType{MsgTypeCreateVestingAccount},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateVestingAccount)},
		}, {
			name:   "test 6",
			values: []MsgType{MsgTypeDelegate},
			want:   MsgTypeBits{Bits(MsgTypeBitsDelegate)},
		}, {
			name:   "test 7",
			values: []MsgType{MsgTypeEditValidator},
			want:   MsgTypeBits{Bits(MsgTypeBitsEditValidator)},
		}, {
			name:   "test 8",
			values: []MsgType{MsgTypePayForBlobs},
			want:   MsgTypeBits{Bits(MsgTypeBitsPayForBlobs)},
		}, {
			name:   "test 9",
			values: []MsgType{MsgTypeSend},
			want:   MsgTypeBits{Bits(MsgTypeBitsSend)},
		}, {
			name:   "test 10",
			values: []MsgType{MsgTypeUndelegate},
			want:   MsgTypeBits{Bits(MsgTypeBitsUndelegate)},
		}, {
			name:   "test 11",
			values: []MsgType{MsgTypeUnjail},
			want:   MsgTypeBits{Bits(MsgTypeBitsUnjail)},
		}, {
			name:   "test 12",
			values: []MsgType{MsgTypeUnknown},
			want:   MsgTypeBits{Bits(MsgTypeBitsUnknown)},
		}, {
			name:   "test 13",
			values: []MsgType{MsgTypeWithdrawDelegatorReward},
			want:   MsgTypeBits{Bits(MsgTypeBitsWithdrawDelegatorReward)},
		}, {
			name:   "test 14",
			values: []MsgType{MsgTypeWithdrawDelegatorReward, MsgTypeBeginRedelegate},
			want:   MsgTypeBits{Bits(20)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.want, NewMsgTypeBitMask(tt.values...))
		})
	}
}

func TestMsgTypeBits_SetBit(t *testing.T) {
	tests := []struct {
		name  string
		value MsgType
		want  MsgTypeBits
	}{
		{
			name:  "test 1",
			value: MsgTypeBeginRedelegate,
			want:  NewMsgTypeBitMask(MsgTypeBeginRedelegate),
		}, {
			name:  "test 2",
			value: MsgType("unknown"),
			want:  NewMsgTypeBitMask(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mask := &MsgTypeBits{
				Bits: 0,
			}
			mask.SetBit(tt.value)
			require.EqualValues(t, tt.want.Bits, mask.Bits)
		})
	}
}

func TestMsgTypeBits_HasOne(t *testing.T) {
	tests := []struct {
		name  string
		mask  MsgTypeBits
		value MsgTypeBits
		want  bool
	}{
		{
			name:  "test 1",
			mask:  NewMsgTypeBitMask(MsgTypeBeginRedelegate),
			value: NewMsgTypeBitMask(MsgTypeBeginRedelegate),
			want:  true,
		}, {
			name:  "test 2",
			mask:  NewMsgTypeBitMask(MsgTypeBeginRedelegate, MsgTypeDelegate, MsgTypeSend),
			value: NewMsgTypeBitMask(MsgTypeBeginRedelegate),
			want:  true,
		}, {
			name:  "test 3",
			mask:  NewMsgTypeBitMask(MsgTypeBeginRedelegate),
			value: NewMsgTypeBitMask(MsgTypeBeginRedelegate, MsgTypeDelegate, MsgTypeSend),
			want:  true,
		}, {
			name:  "test 4",
			mask:  NewMsgTypeBitMask(MsgTypeBeginRedelegate),
			value: NewMsgTypeBitMask(MsgTypeDelegate, MsgTypeSend),
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			has := tt.mask.HasOne(tt.value)
			require.Equal(t, tt.want, has)
		})
	}
}
