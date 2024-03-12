// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
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

func handleWithdrawDelegatorRewards(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events hanler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processWithdrawDelegatorRewards(ctx, events, msg, idx)
}

func processWithdrawDelegatorRewards(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	for i := *idx; i < len(events); i++ {
		switch events[i].Type {
		case storageTypes.EventTypeMessage:
			if module := decoder.StringFromMap(events[i].Data, "module"); module == storageTypes.ModuleNameDistribution.String() {
				*idx = i + 1
				return nil
			}
		case storageTypes.EventTypeWithdrawRewards:
			if err := parseWithdrawRewards(ctx, msg, events[i].Data); err != nil {
				return err
			}
		}
	}
	*idx = len(events) - 1
	return nil
}
