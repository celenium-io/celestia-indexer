// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/shopspring/decimal"
)

func parseWithdrawRewards(ctx *context.Context, msg *storage.Message, data map[string]any) error {
	rewards, err := decode.NewWithdrawRewards(data)
	if err != nil {
		return err
	}

	if rewards.Validator == "" || rewards.Delegator != "" {
		return nil
	}

	validator := storage.EmptyValidator()
	validator.Address = rewards.Validator

	if rewards.Amount != nil {
		amount, err := decimal.NewFromString(rewards.Amount.Amount.String())
		if err != nil {
			return err
		}
		validator.Rewards = amount.Neg()
	}

	ctx.AddValidator(validator)

	ctx.AddStakingLog(storage.StakingLog{
		Height:    msg.Height,
		Time:      msg.Time,
		Validator: &validator,
		Change:    validator.Rewards.Copy(),
		Type:      storageTypes.StakingLogTypeRewards,
	})

	return nil
}
