// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func handleCancelUnbonding(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processCancelUnbonding(ctx, events, msg, idx)
}

func processCancelUnbonding(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	var (
		msgIdx    = decoder.StringFromMap(events[*idx].Data, "msg_index")
		newFormat = msgIdx != ""
	)

	for i := *idx; i < len(events); i++ {
		switch events[i].Type {
		case storageTypes.EventTypeMessage:
			if module := decoder.StringFromMap(events[i].Data, "module"); module == storageTypes.ModuleNameStaking.String() {
				*idx = i + 1
				return nil
			}
		case storageTypes.EventTypeWithdrawRewards:
			if err := parseWithdrawRewards(ctx, msg, events[i].Data); err != nil {
				return err
			}
		case storageTypes.EventTypeCancelUnbondingDelegation:
			cancel, err := decode.NewCancelUnbondingDelegation(events[i].Data)
			if err != nil {
				return err
			}

			amount := decimal.RequireFromString(cancel.Amount.Amount.String())
			validator := storage.EmptyValidator()
			prefix, hash, err := types.Address(cancel.Validator).Decode()
			if err != nil {
				return errors.Wrap(err, "decode validator address")
			}
			if prefix == types.AddressPrefixCelestia {
				addr, err := types.NewValoperAddressFromBytes(hash)
				if err != nil {
					return errors.Wrap(err, "encode validator address")
				}
				validator.Address = addr.String()
			} else {
				validator.Address = cancel.Validator
			}
			validator.Stake = amount.Copy()
			ctx.AddValidator(validator)

			address := &storage.Address{
				Address:    cancel.Delegator,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balance: storage.Balance{
					Currency:  currency.DefaultCurrency,
					Delegated: amount.Copy(),
					Unbonding: amount.Copy().Neg(),
				},
			}
			if err := ctx.AddAddress(address); err != nil {
				return err
			}

			ctx.AddDelegation(storage.Delegation{
				Validator: &validator,
				Address:   address,
				Amount:    amount,
			})

			ctx.AddCancelUndelegation(storage.Undelegation{
				Validator: &validator,
				Address:   address,
				Height:    msg.Height,
				Time:      msg.Time,
				Amount:    amount,
			})

			ctx.AddStakingLog(storage.StakingLog{
				Height:    msg.Height,
				Time:      msg.Time,
				Address:   address,
				Validator: &validator,
				Change:    amount.Copy(),
				Type:      storageTypes.StakingLogTypeUnbonding,
			})

			if newFormat {
				*idx = i + 1
				return nil
			}
		}
	}

	toTheNextAction(events, idx)
	return nil
}
