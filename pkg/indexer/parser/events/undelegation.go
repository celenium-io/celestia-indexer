// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"time"

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

func handleUndelegate(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.staking.v1beta1.MsgUndelegate" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processUndelegate(ctx, events, msg, idx)
}

func processUndelegate(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	var (
		amount         = decimal.Zero
		validator      = storage.EmptyValidator()
		completionTime = time.Now()
		msgIdx         = decoder.StringFromMap(events[*idx].Data, "msg_index")
		newFormat      = msgIdx != ""

		undelegationEnd = func(event storage.Event, key string) error {
			delegator := decoder.StringFromMap(event.Data, key)

			address := &storage.Address{
				Address:    delegator,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balance: storage.Balance{
					Currency:  currency.DefaultCurrency,
					Delegated: amount.Copy().Neg(),
					Unbonding: amount,
					Spendable: decimal.Zero,
				},
			}
			if err := ctx.AddAddress(address); err != nil {
				return err
			}

			ctx.AddDelegation(storage.Delegation{
				Address:   address,
				Validator: &validator,
				Amount:    amount.Copy().Neg(),
			})

			ctx.AddUndelegation(storage.Undelegation{
				Validator:      &validator,
				Address:        address,
				Amount:         amount,
				Time:           msg.Time,
				Height:         msg.Height,
				CompletionTime: completionTime,
			})

			ctx.AddStakingLog(storage.StakingLog{
				Height:    msg.Height,
				Time:      msg.Time,
				Address:   address,
				Validator: &validator,
				Change:    amount.Copy().Neg(),
				Type:      storageTypes.StakingLogTypeUnbonding,
			})

			return nil
		}
	)

	for i := *idx; i < len(events); i++ {
		switch events[i].Type {
		case storageTypes.EventTypeMessage:
			if module := decoder.StringFromMap(events[i].Data, "module"); module == storageTypes.ModuleNameStaking.String() {
				if err := undelegationEnd(events[i], "sender"); err != nil {
					return errors.Wrap(err, "undelegation end")
				}

				*idx = i + 1
				return nil
			}
		case storageTypes.EventTypeWithdrawRewards:
			if err := parseWithdrawRewards(ctx, msg, events[i].Data); err != nil {
				return err
			}
		case storageTypes.EventTypeUnbond:
			unbond, err := decode.NewUnbond(events[i].Data)
			if err != nil {
				return err
			}

			completionTime = unbond.CompletionTime
			amount = decimal.RequireFromString(unbond.Amount.Amount.String())
			prefix, hash, err := types.Address(unbond.Validator).Decode()
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
				validator.Address = unbond.Validator
			}
			validator.Stake = amount.Copy().Neg()
			ctx.AddValidator(validator)

			if newFormat {
				if err := undelegationEnd(events[i], "delegator"); err != nil {
					return errors.Wrap(err, "undelegation end")
				}
				*idx = i + 1
				return nil
			}
		}
	}

	toTheNextAction(events, idx)
	return nil
}
