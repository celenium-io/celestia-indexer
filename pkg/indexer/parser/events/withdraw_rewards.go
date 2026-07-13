// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func parseWithdrawRewards(ctx *context.Context, msg *storage.Message, data map[string]string) error {
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
		amount, err := storageTypes.NumericFromString(rewards.Amount.Amount.String())
		if err != nil {
			return errors.Wrap(err, "parse withdraw rewards amount")
		}
		validator.Rewards = amount.Neg()
	}

	ctx.AddValidator(validator)

	rewardReceiver := &storage.Address{
		Address:    rewards.Delegator,
		Height:     ctx.Block.Height,
		LastHeight: ctx.Block.Height,
		Balances:   []storage.Balance{storage.EmptyBalance()},
	}
	ctx.AddStakingLog(storage.StakingLog{
		Height:    msg.Height,
		Time:      msg.Time,
		Validator: &validator,
		Address:   rewardReceiver,
		Change:    validator.Rewards.Copy(),
		Type:      storageTypes.StakingLogTypeRewards,
	})

	if err := ctx.AddAddress(rewardReceiver); err != nil {
		return err
	}

	return nil
}

func handleWithdrawDelegatorRewards(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	return processWithdrawDelegatorRewards(ctx, c, msg)
}

func processWithdrawDelegatorRewards(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	first, _ := c.Peek()
	msgIdx := decoder.StringFromMap(first.Data, "msg_index")
	newFormat := msgIdx != ""

	for {
		event, ok := c.Next()
		if !ok {
			return nil
		}
		switch event.Type {
		case storageTypes.EventTypeMessage:
			if newFormat {
				continue
			}
			if module := decoder.StringFromMap(event.Data, "module"); module == storageTypes.ModuleNameDistribution.String() {
				return nil
			}
		case storageTypes.EventTypeWithdrawRewards:
			if err := parseWithdrawRewards(ctx, msg, event.Data); err != nil {
				return err
			}
			if newFormat {
				return nil
			}
		}
	}
}
