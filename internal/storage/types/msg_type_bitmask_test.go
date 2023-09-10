package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgTypeBits_Names(t *testing.T) {
	tests := []struct {
		name string
		Bits Bits
		want []MsgType
	}{
		{
			name: string(MsgBeginRedelegate),
			Bits: Bits(MsgTypeBitsBeginRedelegate),
			want: []MsgType{MsgBeginRedelegate},
		}, {
			name: string(MsgWithdrawValidatorCommission),
			Bits: Bits(MsgTypeBitsWithdrawValidatorCommission),
			want: []MsgType{MsgWithdrawValidatorCommission},
		}, {
			name: string(MsgCreatePeriodicVestingAccount),
			Bits: Bits(MsgTypeBitsCreatePeriodicVestingAccount),
			want: []MsgType{MsgCreatePeriodicVestingAccount},
		}, {
			name: string(MsgCreateValidator),
			Bits: Bits(MsgTypeBitsCreateValidator),
			want: []MsgType{MsgCreateValidator},
		}, {
			name: string(MsgCreateVestingAccount),
			Bits: Bits(MsgTypeBitsCreateVestingAccount),
			want: []MsgType{MsgCreateVestingAccount},
		}, {
			name: string(MsgDelegate),
			Bits: Bits(MsgTypeBitsDelegate),
			want: []MsgType{MsgDelegate},
		}, {
			name: string(MsgEditValidator),
			Bits: Bits(MsgTypeBitsEditValidator),
			want: []MsgType{MsgEditValidator},
		}, {
			name: string(MsgPayForBlobs),
			Bits: Bits(MsgTypeBitsPayForBlobs),
			want: []MsgType{MsgPayForBlobs},
		}, {
			name: string(MsgSend),
			Bits: Bits(MsgTypeBitsSend),
			want: []MsgType{MsgSend},
		}, {
			name: string(MsgUndelegate),
			Bits: Bits(MsgTypeBitsUndelegate),
			want: []MsgType{MsgUndelegate},
		}, {
			name: string(MsgUnjail),
			Bits: Bits(MsgTypeBitsUnjail),
			want: []MsgType{MsgUnjail},
		}, {
			name: string(MsgUnknown),
			Bits: Bits(MsgTypeBitsUnknown),
			want: []MsgType{MsgUnknown},
		}, {
			name: string(MsgWithdrawDelegatorReward),
			Bits: Bits(MsgTypeBitsWithdrawDelegatorReward),
			want: []MsgType{MsgWithdrawDelegatorReward},
		}, {
			name: string(MsgGrantAllowance),
			Bits: Bits(MsgTypeBitsGrantAllowance),
			want: []MsgType{MsgGrantAllowance},
		}, {
			name: string(MsgRegisterEVMAddress),
			Bits: Bits(MsgTypeBitsRegisterEVMAddress),
			want: []MsgType{MsgRegisterEVMAddress},
		}, {
			name: string(MsgSetWithdrawAddress),
			Bits: Bits(MsgTypeBitsSetWithdrawAddress),
			want: []MsgType{MsgSetWithdrawAddress},
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
			values: []MsgType{MsgBeginRedelegate},
			want:   MsgTypeBits{Bits(MsgTypeBitsBeginRedelegate)},
		}, {
			name:   "test 2",
			values: []MsgType{MsgWithdrawValidatorCommission},
			want:   MsgTypeBits{Bits(MsgTypeBitsWithdrawValidatorCommission)},
		}, {
			name:   "test 3",
			values: []MsgType{MsgCreatePeriodicVestingAccount},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreatePeriodicVestingAccount)},
		}, {
			name:   "test 4",
			values: []MsgType{MsgCreateValidator},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateValidator)},
		}, {
			name:   "test 5",
			values: []MsgType{MsgCreateVestingAccount},
			want:   MsgTypeBits{Bits(MsgTypeBitsCreateVestingAccount)},
		}, {
			name:   "test 6",
			values: []MsgType{MsgDelegate},
			want:   MsgTypeBits{Bits(MsgTypeBitsDelegate)},
		}, {
			name:   "test 7",
			values: []MsgType{MsgEditValidator},
			want:   MsgTypeBits{Bits(MsgTypeBitsEditValidator)},
		}, {
			name:   "test 8",
			values: []MsgType{MsgPayForBlobs},
			want:   MsgTypeBits{Bits(MsgTypeBitsPayForBlobs)},
		}, {
			name:   "test 9",
			values: []MsgType{MsgSend},
			want:   MsgTypeBits{Bits(MsgTypeBitsSend)},
		}, {
			name:   "test 10",
			values: []MsgType{MsgUndelegate},
			want:   MsgTypeBits{Bits(MsgTypeBitsUndelegate)},
		}, {
			name:   "test 11",
			values: []MsgType{MsgUnjail},
			want:   MsgTypeBits{Bits(MsgTypeBitsUnjail)},
		}, {
			name:   "test 12",
			values: []MsgType{MsgUnknown},
			want:   MsgTypeBits{Bits(MsgTypeBitsUnknown)},
		}, {
			name:   "test 13",
			values: []MsgType{MsgWithdrawDelegatorReward},
			want:   MsgTypeBits{Bits(MsgTypeBitsWithdrawDelegatorReward)},
		}, {
			name:   "test 14",
			values: []MsgType{MsgWithdrawDelegatorReward, MsgBeginRedelegate},
			want:   MsgTypeBits{Bits(20)},
		}, {
			name:   "test 15",
			values: []MsgType{MsgGrantAllowance},
			want:   MsgTypeBits{Bits(MsgTypeBitsGrantAllowance)},
		}, {
			name:   "test 16",
			values: []MsgType{MsgRegisterEVMAddress},
			want:   MsgTypeBits{Bits(MsgTypeBitsRegisterEVMAddress)},
		}, {
			name:   "test 17",
			values: []MsgType{MsgSetWithdrawAddress},
			want:   MsgTypeBits{Bits(MsgTypeBitsSetWithdrawAddress)},
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
			value: MsgBeginRedelegate,
			want:  NewMsgTypeBitMask(MsgBeginRedelegate),
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
			mask:  NewMsgTypeBitMask(MsgBeginRedelegate),
			value: NewMsgTypeBitMask(MsgBeginRedelegate),
			want:  true,
		}, {
			name:  "test 2",
			mask:  NewMsgTypeBitMask(MsgBeginRedelegate, MsgDelegate, MsgSend),
			value: NewMsgTypeBitMask(MsgBeginRedelegate),
			want:  true,
		}, {
			name:  "test 3",
			mask:  NewMsgTypeBitMask(MsgBeginRedelegate),
			value: NewMsgTypeBitMask(MsgBeginRedelegate, MsgDelegate, MsgSend),
			want:  true,
		}, {
			name:  "test 4",
			mask:  NewMsgTypeBitMask(MsgBeginRedelegate),
			value: NewMsgTypeBitMask(MsgDelegate, MsgSend),
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
