// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
func MsgCreateVestingAccount(ctx *context.Context, status storageTypes.Status, txId, msgId uint64, m *cosmosVestingTypes.MsgCreateVestingAccount) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreateVestingAccount
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	v := &storage.VestingAccount{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		Address: &storage.Address{
			Address: m.ToAddress,
		},
		TxId: &txId,
	}

	amount := storageTypes.NewNumeric(decimal.NewFromBigInt(m.Amount.AmountOf(currency.Utia).BigInt(), 0))
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

	ctx.AddVestingAccount(v)
	return msgType, err
}

// MsgCreatePermanentLockedAccount defines a message that enables creating a permanent
// locked account.
//
// Since: cosmos-sdk 0.46
func MsgCreatePermanentLockedAccount(ctx *context.Context, status storageTypes.Status, txId, msgId uint64, m *cosmosVestingTypes.MsgCreatePermanentLockedAccount) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreatePermanentLockedAccount
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	v := &storage.VestingAccount{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		Address: &storage.Address{
			Address: m.ToAddress,
		},
		Type: storageTypes.VestingTypePermanent,
		TxId: &txId,
	}

	amount := storageTypes.NewNumeric(decimal.NewFromBigInt(m.Amount.AmountOf(currency.Utia).BigInt(), 0))
	v.Amount = v.Amount.Add(amount)

	ctx.AddVestingAccount(v)
	return msgType, err
}

// MsgCreateVestingAccount defines a message that enables creating a vesting
// account.
//
// Since: cosmos-sdk 0.46
func MsgCreatePeriodicVestingAccount(ctx *context.Context, status storageTypes.Status, txId, msgId uint64, m *cosmosVestingTypes.MsgCreatePeriodicVestingAccount) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgCreatePeriodicVestingAccount
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, ctx.Block.Height, msgId)
	if err != nil || status == storageTypes.StatusFailed {
		return msgType, err
	}

	v := &storage.VestingAccount{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		Address: &storage.Address{
			Address: m.ToAddress,
		},
		Type:   storageTypes.VestingTypePeriodic,
		Amount: storageTypes.NewNumeric(decimal.Zero),
		TxId:   &txId,
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
		period.Amount = storageTypes.NewNumeric(decimal.NewFromBigInt(m.VestingPeriods[i].Amount.AmountOf(currency.Utia).BigInt(), 0))
		v.Amount = v.Amount.Add(period.Amount)
		periodTime = periodTime.Add(time.Second * time.Duration(m.VestingPeriods[i].Length))
		period.Time = periodTime
		v.VestingPeriods = append(v.VestingPeriods, period)
	}

	ctx.AddVestingAccount(v)
	return msgType, err
}
