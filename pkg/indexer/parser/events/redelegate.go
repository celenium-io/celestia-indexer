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

func handleRedelegate(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.staking.v1beta1.MsgBeginRedelegate" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processRedelegate(ctx, events, msg, idx)
}

func processRedelegate(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
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
		case storageTypes.EventTypeRedelegate:
			redelegate, err := decode.NewRedelegate(events[i].Data)
			if err != nil {
				return err
			}
			amount := decimal.RequireFromString(redelegate.Amount.Amount.String())

			source := storage.EmptyValidator()
			prefix, hash, err := types.Address(redelegate.SrcValidator).Decode()
			if err != nil {
				return errors.Wrap(err, "decode validator address")
			}
			if prefix == types.AddressPrefixCelestia {
				addr, err := types.NewValoperAddressFromBytes(hash)
				if err != nil {
					return errors.Wrap(err, "encode validator address")
				}
				source.Address = addr.String()
			} else {
				source.Address = redelegate.SrcValidator
			}
			source.Stake = amount.Copy().Neg()
			ctx.AddValidator(source)

			dest := storage.EmptyValidator()
			prefix, hash, err = types.Address(redelegate.DestValidator).Decode()
			if err != nil {
				return errors.Wrap(err, "decode validator address")
			}
			if prefix == types.AddressPrefixCelestia {
				addr, err := types.NewValoperAddressFromBytes(hash)
				if err != nil {
					return errors.Wrap(err, "encode validator address")
				}
				dest.Address = addr.String()
			} else {
				dest.Address = redelegate.DestValidator
			}
			dest.Stake = amount
			ctx.AddValidator(dest)

			delegator := decoder.StringFromMap(msg.Data, "DelegatorAddress")

			address := &storage.Address{
				Address:    delegator,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balance: storage.Balance{
					Currency:  currency.DefaultCurrency,
					Delegated: decimal.Zero,
				},
			}
			if err := ctx.AddAddress(address); err != nil {
				return err
			}

			ctx.AddDelegation(storage.Delegation{
				Validator: &source,
				Address:   address,
				Amount:    amount.Copy().Neg(),
			})
			ctx.AddDelegation(storage.Delegation{
				Validator: &dest,
				Address:   address,
				Amount:    amount.Copy(),
			})
			ctx.AddRedelegation(storage.Redelegation{
				Source:         &source,
				Destination:    &dest,
				Address:        address,
				Amount:         amount,
				Time:           msg.Time,
				Height:         msg.Height,
				CompletionTime: redelegate.CompletionTime,
			})
			ctx.AddStakingLog(storage.StakingLog{
				Height:    msg.Height,
				Time:      msg.Time,
				Address:   address,
				Validator: &source,
				Change:    amount.Copy().Neg(),
				Type:      storageTypes.StakingLogTypeDelegation,
			})
			ctx.AddStakingLog(storage.StakingLog{
				Height:    msg.Height,
				Time:      msg.Time,
				Address:   address,
				Validator: &dest,
				Change:    amount.Copy(),
				Type:      storageTypes.StakingLogTypeDelegation,
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
