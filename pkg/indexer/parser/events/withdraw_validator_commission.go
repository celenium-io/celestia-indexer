// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
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

func handleWithdrawValidatorCommission(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	return processWithdrawValidatorCommission(ctx, events, msg, idx)
}

func processWithdrawValidatorCommission(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	var validator = storage.EmptyValidator()

	var newFormat bool
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action == "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission" {
		validator.Address = decoder.StringFromMap(events[*idx].Data, "sender")
		*idx += 1
		newFormat = true
	}

	for i := *idx; i < len(events); i++ {
		switch events[i].Type {
		case storageTypes.EventTypeMessage:
			if module := decoder.StringFromMap(events[i].Data, "module"); module == storageTypes.ModuleNameDistribution.String() {
				if !validator.Commissions.IsZero() {
					validator.Address = decoder.StringFromMap(events[i].Data, "sender")
					ctx.AddValidator(validator)

					ctx.AddStakingLog(storage.StakingLog{
						Height:    msg.Height,
						Time:      msg.Time,
						Validator: &validator,
						Change:    validator.Commissions.Copy(),
						Type:      storageTypes.StakingLogTypeCommissions,
					})
				}

				*idx = i + 1
				return nil
			}
		case storageTypes.EventTypeWithdrawCommission:
			commission, err := decode.NewWithdrawCommission(events[i].Data)
			if err != nil {
				return err
			}
			if commission.Amount == nil {
				continue
			}

			amount := decimal.RequireFromString(commission.Amount.Amount.String())
			validator.Commissions = amount.Neg()

			if newFormat {
				ctx.AddValidator(validator)

				ctx.AddStakingLog(storage.StakingLog{
					Height:    msg.Height,
					Time:      msg.Time,
					Validator: &validator,
					Change:    validator.Commissions.Copy(),
					Type:      storageTypes.StakingLogTypeCommissions,
				})
				*idx = i + 1
				return nil
			}
		}
	}

	toTheNextAction(events, idx)
	return nil
}
