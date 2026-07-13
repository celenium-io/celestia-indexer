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
)

const msgWithdrawValidatorCommission = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"

func handleWithdrawValidatorCommission(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != msgWithdrawValidatorCommission {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	return processWithdrawValidatorCommission(ctx, c, msg)
}

func processWithdrawValidatorCommission(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	var validator = storage.EmptyValidator()

	if validatorAddress := msg.Data.GetStringOrDefault("ValidatorAddress"); validatorAddress != "" {
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

	for {
		event, ok := c.Peek()
		if !ok {
			break
		}
		if event.Type != storageTypes.EventTypeWithdrawCommission {
			c.Next()
			continue
		}

		commission, err := decode.NewWithdrawCommission(event.Data)
		if err != nil {
			return err
		}
		if commission.Amount == nil {
			c.Next()
			continue
		}

		amount, err := storageTypes.NumericFromString(commission.Amount.Amount.String())
		if err != nil {
			return errors.Wrap(err, "parse withdraw commission amount")
		}
		validator.Commissions = amount.Neg()

		ctx.AddValidator(validator)

		ctx.AddStakingLog(storage.StakingLog{
			Height:    msg.Height,
			Time:      msg.Time,
			Validator: &validator,
			Change:    amount.Neg().Copy(),
			Type:      storageTypes.StakingLogTypeCommissions,
		})
		break
	}
	return nil
}
