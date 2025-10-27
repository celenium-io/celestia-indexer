// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const msgWithdrawValidatorCommission = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"

func handleWithdrawValidatorCommission(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != msgWithdrawValidatorCommission {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	return processWithdrawValidatorCommission(ctx, events, msg, idx)
}

func processWithdrawValidatorCommission(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	msgIdx := decoder.StringFromMap(events[*idx].Data, "msg_index")
	newFormat := msgIdx != ""

	var validator = storage.EmptyValidator()

	if validatorAddress := decoder.StringFromMap(msg.Data, "ValidatorAddress"); validatorAddress != "" {
		prefix, hash, err := types.Address(validatorAddress).Decode()
		if err != nil {
			return errors.Wrap(err, "decoding sender in WithdrawValidatorCommission")
		}

		switch prefix {
		case types.AddressPrefixCelestia:
			address, err := types.NewValoperAddressFromBytes(hash)
			if err != nil {
				return errors.Wrap(err, "creating valoper address in WithdrawValidatorCommission")
			}
			validator.Address = address.String()
		case types.AddressPrefixValoper:
			validator.Address = validatorAddress
		default:
			return errors.Errorf("unexpected sender address prefix in WithdrawValidatorCommission: %s", prefix)
		}
	}

	if validator.Address == "" {
		return errors.Errorf("empty validator address in WithdrawValidatorCommission: %##v", msg.Data)
	}

	for ; *idx < len(events); *idx++ {
		if events[*idx].Type == storageTypes.EventTypeWithdrawCommission {
			commission, err := decode.NewWithdrawCommission(events[*idx].Data)
			if err != nil {
				return err
			}
			if commission.Amount == nil {
				continue
			}

			amount := decimal.RequireFromString(commission.Amount.Amount.String())
			validator.Commissions = amount.Neg()

			ctx.AddValidator(validator)

			ctx.AddStakingLog(storage.StakingLog{
				Height:    msg.Height,
				Time:      msg.Time,
				Validator: &validator,
				Change:    amount.Neg().Copy(),
				Type:      storageTypes.StakingLogTypeCommissions,
			})

			if !newFormat {
				*idx++
			}
			break
		}
	}

	return nil
}
