// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/hex"
	"strings"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

func parseCoinSpent(ctx *context.Context, data map[string]any, height pkgTypes.Level) error {
	coinSpent, err := decode.NewCoinSpent(data)
	if err != nil {
		return err
	}

	if coinSpent.Spender == "" {
		return nil
	}

	address := &storage.Address{
		Address:    coinSpent.Spender,
		Height:     height,
		LastHeight: height,
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.Zero,
			Delegated: decimal.Zero,
			Unbonding: decimal.Zero,
		},
	}

	if coinSpent.Amount != nil {
		address.Balance.Currency = coinSpent.Amount.Denom
		amount, err := decimal.NewFromString(coinSpent.Amount.Amount.String())
		if err != nil {
			return err
		}
		address.Balance.Spendable = amount.Neg()
	}

	return ctx.AddAddress(address)
}

func parseCoinReceived(ctx *context.Context, data map[string]any, height pkgTypes.Level) error {
	coinReceived, err := decode.NewCoinReceived(data)
	if err != nil {
		return err
	}

	if coinReceived.Receiver == "" {
		return nil
	}

	address := &storage.Address{
		Address:    coinReceived.Receiver,
		Height:     height,
		LastHeight: height,
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.Zero,
			Delegated: decimal.Zero,
			Unbonding: decimal.Zero,
		},
	}

	if coinReceived.Amount != nil {
		address.Balance.Currency = coinReceived.Amount.Denom
		amount, err := decimal.NewFromString(coinReceived.Amount.Amount.String())
		if err != nil {
			return err
		}
		address.Balance.Spendable = amount
	}

	return ctx.AddAddress(address)
}

func parseCompleteUnbonding(ctx *context.Context, data map[string]any, height pkgTypes.Level) error {
	unbonding, err := decode.NewCompleteUnbonding(data)
	if err != nil {
		return err
	}

	if unbonding.Validator == "" || unbonding.Delegator == "" {
		return nil
	}

	address := &storage.Address{
		Address:    unbonding.Delegator,
		Height:     height,
		LastHeight: height,
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.Zero,
			Delegated: decimal.Zero,
			Unbonding: decimal.Zero,
		},
	}

	if unbonding.Amount != nil {
		address.Balance.Currency = unbonding.Amount.Denom
		amount, err := decimal.NewFromString(unbonding.Amount.Amount.String())
		if err != nil {
			return err
		}
		address.Balance.Unbonding = amount.Neg()
	}
	return ctx.AddAddress(address)
}

func parseCommission(ctx *context.Context, data map[string]any) error {
	commission, err := decode.NewCommission(data)
	if err != nil {
		return err
	}

	if commission.Validator == "" {
		return nil
	}

	validator := storage.Validator{
		Address:     commission.Validator,
		Rewards:     decimal.Zero,
		Commissions: decimal.Zero,
		Stake:       decimal.Zero,
	}

	if commission.Amount != nil {
		amount, err := decimal.NewFromString(commission.Amount.Amount.String())
		if err != nil {
			return err
		}
		validator.Commissions = amount
		ctx.Block.Stats.Commissions = ctx.Block.Stats.Commissions.Add(amount)

		ctx.AddStakingLog(storage.StakingLog{
			Height:    ctx.Block.Height,
			Time:      ctx.Block.Time,
			Validator: &validator,
			Change:    amount.Copy(),
			Type:      types.StakingLogTypeCommissions,
		})
	}

	ctx.AddValidator(validator)
	return nil
}

func parseRewards(ctx *context.Context, data map[string]any) error {
	rewards, err := decode.NewRewards(data)
	if err != nil {
		return err
	}

	if rewards.Validator == "" {
		return nil
	}

	validator := storage.Validator{
		Address:     rewards.Validator,
		Rewards:     decimal.Zero,
		Commissions: decimal.Zero,
		Stake:       decimal.Zero,
	}

	if rewards.Amount != nil {
		amount, err := decimal.NewFromString(rewards.Amount.Amount.String())
		if err != nil {
			return err
		}
		validator.Rewards = amount
		ctx.Block.Stats.Rewards = ctx.Block.Stats.Rewards.Add(amount)

		ctx.AddStakingLog(storage.StakingLog{
			Height:    ctx.Block.Height,
			Time:      ctx.Block.Time,
			Validator: &validator,
			Change:    amount.Copy(),
			Type:      types.StakingLogTypeRewards,
		})
	}

	ctx.AddValidator(validator)
	return nil
}

func parseSlash(ctx *context.Context, data map[string]any) error {
	slash, err := decode.NewSlash(data)
	if err != nil {
		return err
	}

	if slash.Address != "" && slash.Reason != "" {
		_, hash, err := pkgTypes.Address(slash.Address).Decode()
		if err != nil {
			return err
		}
		consAddress := strings.ToUpper(hex.EncodeToString(hash))
		ctx.AddJailedValidator(consAddress)

		ctx.AddJail(storage.Jail{
			Height: ctx.Block.Height,
			Time:   ctx.Block.Time,
			Reason: slash.Reason,
			Burned: slash.BurnedCoins,
			Validator: &storage.Validator{
				ConsAddress: consAddress,
			},
		})
	}

	return nil
}
