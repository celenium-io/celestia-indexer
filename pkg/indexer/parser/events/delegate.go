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
)

func handleDelegate(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/cosmos.staking.v1beta1.MsgDelegate" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processDelegate(ctx, c, msg)
}

func processDelegate(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	first, _ := c.Peek()
	var (
		validator  = storage.EmptyValidator()
		delegation = storage.Delegation{
			Validator: &validator,
		}
		msgIdx    = decoder.StringFromMap(first.Data, "msg_index")
		newFormat = msgIdx != ""

		endDelegation = func(event storage.Event, key string) error {
			delegator := decoder.StringFromMap(event.Data, key)

			address := &storage.Address{
				Address:    delegator,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balances: []storage.Balance{
					{
						Currency:  currency.DefaultCurrency,
						Delegated: delegation.Amount,
					},
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

	for {
		event, ok := c.Next()
		if !ok {
			return nil
		}
		switch event.Type {
		case storageTypes.EventTypeMessage:
			if module := decoder.StringFromMap(event.Data, "module"); module == storageTypes.ModuleNameStaking.String() {
				return errors.Wrap(endDelegation(event, "sender"), "end delegation")
			}
		case storageTypes.EventTypeWithdrawRewards:
			if err := parseWithdrawRewards(ctx, msg, event.Data); err != nil {
				return err
			}
		case storageTypes.EventTypeDelegate:
			delegate, err := decode.NewDelegate(event.Data)
			if err != nil {
				return err
			}
			amount, err := storageTypes.NumericFromString(delegate.Amount.Amount.String())
			if err != nil {
				return errors.Wrap(err, "parse delegation amount")
			}
			delegation.Amount = amount
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
				return errors.Wrap(endDelegation(event, "delegator"), "end delegation")
			}
		}
	}
}
