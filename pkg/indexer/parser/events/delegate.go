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

func handleDelegate(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.staking.v1beta1.MsgDelegate" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processDelegate(ctx, events, msg, idx)
}

func processDelegate(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	var (
		validator  = storage.EmptyValidator()
		delegation = storage.Delegation{
			Validator: &validator,
		}
		msgIdx    = decoder.StringFromMap(events[*idx].Data, "msg_index")
		newFormat = msgIdx != ""

		endDelegation = func(event storage.Event, key string) error {
			delegator := decoder.StringFromMap(event.Data, key)

			address := &storage.Address{
				Address:    delegator,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balance: storage.Balance{
					Currency:  currency.DefaultCurrency,
					Delegated: delegation.Amount,
				},
			}

			if err := ctx.AddAddress(address); err != nil {
				return err
			}
			delegation.Address = address

			ctx.AddDelegation(delegation)

			ctx.AddStakingLog(storage.StakingLog{
				Height:    msg.Height,
				Time:      msg.Time,
				Address:   address,
				Validator: delegation.Validator,
				Change:    delegation.Amount,
				Type:      storageTypes.StakingLogTypeDelegation,
			})
			return nil
		}
	)

	for i := *idx; i < len(events); i++ {
		switch events[i].Type {
		case storageTypes.EventTypeMessage:
			if module := decoder.StringFromMap(events[i].Data, "module"); module == storageTypes.ModuleNameStaking.String() {
				if err := endDelegation(events[i], "sender"); err != nil {
					return errors.Wrap(err, "end delegation")
				}
				*idx = i + 1
				return nil
			}
		case storageTypes.EventTypeWithdrawRewards:
			if err := parseWithdrawRewards(ctx, msg, events[i].Data); err != nil {
				return err
			}
		case storageTypes.EventTypeDelegate:
			delegate, err := decode.NewDelegate(events[i].Data)
			if err != nil {
				return err
			}
			delegation.Amount = decimal.RequireFromString(delegate.Amount.Amount.String())
			prefix, hash, err := types.Address(delegate.Validator).Decode()
			if err != nil {
				return errors.Wrap(err, "decode validator address")
			}
			if prefix == types.AddressPrefixCelestia {
				addr, err := types.NewValoperAddressFromBytes(hash)
				if err != nil {
					return errors.Wrap(err, "encode validator address")
				}
				delegation.Validator.Address = addr.String()
			} else {
				delegation.Validator.Address = delegate.Validator
			}
			delegation.Validator.Stake = delegation.Amount
			ctx.AddValidator(*delegation.Validator)

			if newFormat {
				if err := endDelegation(events[i], "delegator"); err != nil {
					return errors.Wrap(err, "end delegation")
				}
				*idx = i + 1
				return nil
			}
		}

	}

	return nil
}
