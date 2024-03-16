// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/shopspring/decimal"
)

// MsgCreateVestingAccount defines a message that enables creating a vesting
// account.
func MsgCreateVestingAccount(ctx *context.Context, status storageTypes.Status, m *cosmosVestingTypes.MsgCreateVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, *storage.VestingAccount, error) {
	msgType := storageTypes.MsgCreateVestingAccount
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, addresses, nil, err
	}

	v := &storage.VestingAccount{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		Address: &storage.Address{
			Address: m.ToAddress,
		},
	}

	amount := decimal.NewFromBigInt(m.Amount.AmountOf(currency.Utia).BigInt(), 0)
	v.Amount = v.Amount.Add(amount)

	if m.EndTime > 0 {
		t := time.Unix(m.EndTime, 0).UTC()
		v.EndTime = &t
	}
	if m.StartTime > 0 {
		t := time.Unix(m.StartTime, 0).UTC()
		v.StartTime = &t
	}

	if m.Delayed {
		v.Type = storageTypes.VestingTypeDelayed
	} else {
		v.Type = storageTypes.VestingTypeContinuous
	}

	return msgType, addresses, v, err
}

// MsgCreatePermanentLockedAccount defines a message that enables creating a permanent
// locked account.
//
// Since: cosmos-sdk 0.46
func MsgCreatePermanentLockedAccount(ctx *context.Context, status storageTypes.Status, m *cosmosVestingTypes.MsgCreatePermanentLockedAccount) (storageTypes.MsgType, []storage.AddressWithType, *storage.VestingAccount, error) {
	msgType := storageTypes.MsgCreatePermanentLockedAccount
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, addresses, nil, err
	}

	v := &storage.VestingAccount{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		Address: &storage.Address{
			Address: m.ToAddress,
		},
		Type: storageTypes.VestingTypePermanent,
	}

	amount := decimal.NewFromBigInt(m.Amount.AmountOf(currency.Utia).BigInt(), 0)
	v.Amount = v.Amount.Add(amount)

	return msgType, addresses, v, err
}

// MsgCreateVestingAccount defines a message that enables creating a vesting
// account.
//
// Since: cosmos-sdk 0.46
func MsgCreatePeriodicVestingAccount(ctx *context.Context, status storageTypes.Status, m *cosmosVestingTypes.MsgCreatePeriodicVestingAccount) (storageTypes.MsgType, []storage.AddressWithType, *storage.VestingAccount, error) {
	msgType := storageTypes.MsgCreatePeriodicVestingAccount
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, addresses, nil, err
	}

	v := &storage.VestingAccount{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		Address: &storage.Address{
			Address: m.ToAddress,
		},
		Type:   storageTypes.VestingTypePeriodic,
		Amount: decimal.Zero,
	}

	periodTime := v.Time
	if m.StartTime > 0 {
		t := time.Unix(m.StartTime, 0).UTC()
		v.StartTime = &t
		periodTime = t
	}

	for i := range m.VestingPeriods {
		period := storage.VestingPeriod{
			Height: v.Height,
		}
		period.Amount = decimal.NewFromBigInt(m.VestingPeriods[i].Amount.AmountOf(currency.Utia).BigInt(), 0)
		v.Amount = v.Amount.Add(period.Amount)
		periodTime = periodTime.Add(time.Second * time.Duration(m.VestingPeriods[i].Length))
		period.Time = periodTime
		v.VestingPeriods = append(v.VestingPeriods, period)
	}

	return msgType, addresses, v, err
}
