// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/hex"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
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
		amount, err := decimal.NewFromString(coinSpent.Amount.Amount.String())
		if err != nil {
			return err
		}
		address.Balance.Spendable = amount.Neg()
		address.Balance.Currency = coinSpent.Amount.GetDenom()
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
		amount, err := decimal.NewFromString(coinReceived.Amount.Amount.String())
		if err != nil {
			return err
		}
		address.Balance.Spendable = amount
		address.Balance.Currency = coinReceived.Amount.GetDenom()
	}

	return ctx.AddAddress(address)
}

func parseCompleteUnbonding(ctx *context.Context, data map[string]any) error {
	unbonding, err := decode.NewCompleteUnbonding(data)
	if err != nil {
		return err
	}

	if unbonding.Validator == "" || unbonding.Delegator == "" {
		return nil
	}

	address := &storage.Address{
		Address:    unbonding.Delegator,
		Height:     ctx.Block.Height,
		LastHeight: ctx.Block.Height,
		Balance: storage.Balance{
			Currency:  currency.DefaultCurrency,
			Spendable: decimal.Zero,
			Delegated: decimal.Zero,
			Unbonding: decimal.Zero,
		},
	}

	if unbonding.Amount != nil {
		amount, err := decimal.NewFromString(unbonding.Amount.Amount.String())
		if err != nil {
			return err
		}
		address.Balance.Unbonding = amount.Copy().Neg()
		address.Balance.Currency = unbonding.Amount.GetDenom()

		validator := &storage.Validator{
			Address: unbonding.Validator,
		}
		ctx.AddValidator(*validator)
		ctx.AddStakingLog(storage.StakingLog{
			Time:      ctx.Block.Time,
			Height:    ctx.Block.Height,
			Address:   address,
			Validator: validator,
			Change:    amount.Copy().Neg(),
			Type:      types.StakingLogTypeUnbonded,
		})
	}
	return ctx.AddAddress(address)
}

func parseCompleteRedelegation(ctx *context.Context, data map[string]any) error {
	redelegation, err := decode.NewCompleteRedelegation(data)
	if err != nil {
		return err
	}

	if redelegation.SrcValidator == "" || redelegation.DestValidator == "" || redelegation.Delegator == "" {
		return nil
	}

	address := &storage.Address{
		Address:    redelegation.Delegator,
		Height:     ctx.Block.Height,
		LastHeight: ctx.Block.Height,
		Balance:    storage.EmptyBalance(),
	}

	if redelegation.Amount != nil {
		amount, err := decimal.NewFromString(redelegation.Amount.Amount.String())
		if err != nil {
			return err
		}

		validator := &storage.Validator{
			Address: redelegation.SrcValidator,
		}
		ctx.AddValidator(*validator)
		ctx.AddStakingLog(storage.StakingLog{
			Time:      ctx.Block.Time,
			Height:    ctx.Block.Height,
			Address:   address,
			Validator: validator,
			Change:    amount.Copy().Neg(),
			Type:      types.StakingLogTypeUnbonded,
		})
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

	validator := storage.EmptyValidator()
	validator.Address = commission.Validator

	if !commission.Amount.IsZero() {
		validator.Commissions = commission.Amount.Copy()
		ctx.Block.Stats.Commissions = ctx.Block.Stats.Commissions.Add(commission.Amount.Copy())

		ctx.AddStakingLog(storage.StakingLog{
			Height:    ctx.Block.Height,
			Time:      ctx.Block.Time,
			Validator: &validator,
			Change:    commission.Amount.Copy(),
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

	validator := storage.EmptyValidator()
	validator.Address = rewards.Validator

	if !rewards.Amount.IsZero() {
		validator.Rewards = rewards.Amount.Copy()
		ctx.Block.Stats.Rewards = ctx.Block.Stats.Rewards.Add(rewards.Amount)

		ctx.AddStakingLog(storage.StakingLog{
			Height:    ctx.Block.Height,
			Time:      ctx.Block.Time,
			Validator: &validator,
			Change:    rewards.Amount.Copy(),
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

		jailed := true
		ctx.AddJail(storage.Jail{
			Height: ctx.Block.Height,
			Time:   ctx.Block.Time,
			Reason: slash.Reason,
			Burned: slash.BurnedCoins,
			Validator: &storage.Validator{
				ConsAddress: consAddress,
				Stake:       slash.BurnedCoins.Copy(),
				Jailed:      &jailed,
			},
		})
	}

	return nil
}

func parseProposal(ctx *context.Context, data map[string]any) error {
	status, err := decode.NewProposalStatus(data)
	if err != nil {
		return err
	}

	if status.Id > 0 && status.Result != "" {
		proposal := &storage.Proposal{
			Id: status.Id,
		}
		switch status.Result {
		case "proposal_rejected":
			proposal.Status = types.ProposalStatusRejected
		case "proposal_passed":
			proposal.Status = types.ProposalStatusApplied
		case "proposal_dropped":
			proposal.Status = types.ProposalStatusRemoved
		case "proposal_failed":
			proposal.Status = types.ProposalStatusFailed
			proposal.Error = status.Log
		default:
			return errors.Errorf("unknown proposal status: %s", status.Result)
		}

		ctx.AddProposal(proposal)
	}

	return nil
}

func parseCreateIgp(ctx *context.Context, data map[string]any) error {
	igp, err := decode.NewHyperlaneCreateIgpEvent(data)
	if err != nil {
		return err
	}

	if igp.IgpId == "" {
		return nil
	}

	igpId, err := util.DecodeHexAddress(igp.IgpId)
	if err != nil {
		return errors.Wrap(err, "decode igp id")
	}

	newIgp := storage.HLIGP{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		Denom:  igp.Denom,
		IgpId:  igpId.Bytes(),
		Owner: &storage.Address{
			Address: igp.Owner,
		},
	}
	ctx.AddIgp(igp.IgpId, &newIgp)

	return nil
}

func parseSetDestinationGasConfig(ctx *context.Context, data map[string]any) error {
	igp, err := decode.NewHyperlaneSetDestinationGasConfig(data)
	if err != nil {
		return err
	}

	if igp.IgpId == "" {
		return nil
	}

	config := storage.HLIGPConfig{
		Height:            ctx.Block.Height,
		Time:              ctx.Block.Time,
		GasPrice:          igp.GasPrice,
		GasOverhead:       igp.GasOverhead,
		RemoteDomain:      igp.RemoteDomain,
		TokenExchangeRate: igp.TokenExchangeRate,
	}
	ctx.AddIgpConfig(igp.IgpId, &config)

	return nil
}

func parseSetIgp(ctx *context.Context, data map[string]any) error {
	igp, err := decode.NewHyperlaneSetIgpEvent(data)
	if err != nil {
		return err
	}

	if igp.IgpId == "" {
		return nil
	}

	igpId, err := util.DecodeHexAddress(igp.IgpId)
	if err != nil {
		return errors.Wrap(err, "decode igp id")
	}

	newIgp := storage.HLIGP{
		Height: ctx.Block.Height,
		Time:   ctx.Block.Time,
		IgpId:  igpId.Bytes(),
		Owner: &storage.Address{
			Address: igp.NewOwner,
		},
	}
	ctx.AddIgp(igp.IgpId, &newIgp)

	return nil
}
