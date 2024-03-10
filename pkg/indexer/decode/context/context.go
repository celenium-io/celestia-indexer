// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package context

import (
	"cosmossdk.io/errors"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
)

type Context struct {
	Validators  *sync.Map[string, *storage.Validator]
	Addresses   *sync.Map[string, *storage.Address]
	Delegations *sync.Map[string, *storage.Delegation]
	Jails       *sync.Map[string, *storage.Jail]

	Redelegations   []storage.Redelegation
	Undelegations   []storage.Undelegation
	CancelUnbonding []storage.Undelegation
	StakingLogs     []storage.StakingLog

	Block *storage.Block
}

func NewContext() *Context {
	return &Context{
		Validators:      sync.NewMap[string, *storage.Validator](),
		Addresses:       sync.NewMap[string, *storage.Address](),
		Delegations:     sync.NewMap[string, *storage.Delegation](),
		Jails:           sync.NewMap[string, *storage.Jail](),
		Redelegations:   make([]storage.Redelegation, 0),
		Undelegations:   make([]storage.Undelegation, 0),
		CancelUnbonding: make([]storage.Undelegation, 0),
		StakingLogs:     make([]storage.StakingLog, 0),
	}
}

func (ctx *Context) AddAddress(address *storage.Address) error {
	if addr, ok := ctx.Addresses.Get(address.String()); ok {
		addr.Balance.Spendable = addr.Balance.Spendable.Add(address.Balance.Spendable)
		addr.Balance.Delegated = addr.Balance.Delegated.Add(address.Balance.Delegated)
		addr.Balance.Unbonding = addr.Balance.Unbonding.Add(address.Balance.Unbonding)
	} else {
		if len(address.Hash) == 0 {
			_, hash, err := pkgTypes.Address(address.Address).Decode()
			if err != nil {
				return errors.Wrap(err, address.Address)
			}
			address.Hash = hash
		}
		ctx.Addresses.Set(address.String(), address)
	}
	return nil
}

func (ctx *Context) AddValidator(validator storage.Validator) {
	if val, ok := ctx.Validators.Get(validator.Address); ok {
		if !validator.Stake.IsZero() {
			val.Stake = val.Stake.Add(validator.Stake)
		}
		if !validator.Commissions.IsZero() {
			val.Commissions = val.Commissions.Add(validator.Commissions)
		}
		if !validator.Rewards.IsZero() {
			val.Rewards = val.Rewards.Add(validator.Rewards)
		}
		if !validator.MaxChangeRate.IsZero() {
			val.MaxChangeRate = validator.MaxChangeRate.Copy()
		}
		if !validator.MaxRate.IsZero() {
			val.MaxRate = validator.MaxRate.Copy()
		}
		if !validator.MinSelfDelegation.IsZero() {
			val.MinSelfDelegation = validator.MinSelfDelegation.Copy()
		}
		if !validator.Rate.IsZero() {
			val.Rate = validator.Rate.Copy()
		}
		if validator.Delegator != "" {
			val.Delegator = validator.Delegator
		}
		if validator.Contacts != storage.DoNotModify {
			val.Contacts = validator.Contacts
		}
		if validator.Details != storage.DoNotModify {
			val.Details = validator.Details
		}
		if validator.Identity != storage.DoNotModify {
			val.Identity = validator.Identity
		}
		if validator.Moniker != storage.DoNotModify {
			val.Moniker = validator.Moniker
		}
		if validator.Website != storage.DoNotModify {
			val.Website = validator.Website
		}
	} else {
		ctx.Validators.Set(validator.Address, &validator)
	}
}

func (ctx *Context) AddSupply(data map[string]any) {
	ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Add(decoder.Amount(data))
}

func (ctx *Context) SubSupply(data map[string]any) {
	ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Sub(decoder.Amount(data))
}

func (ctx *Context) SetInflation(data map[string]any) {
	ctx.Block.Stats.InflationRate = decoder.DecimalFromMap(data, "inflation_rate")
}

func (ctx *Context) GetValidators() []*storage.Validator {
	validators := make([]*storage.Validator, 0)
	_ = ctx.Validators.Range(func(_ string, value *storage.Validator) (error, bool) {
		validators = append(validators, value)
		return nil, false
	})
	return validators
}

func (ctx *Context) GetAddresses() []*storage.Address {
	addresses := make([]*storage.Address, 0)
	_ = ctx.Addresses.Range(func(_ string, value *storage.Address) (error, bool) {
		addresses = append(addresses, value)
		return nil, false
	})
	return addresses
}

func (ctx *Context) AddDelegation(d storage.Delegation) {
	if val, ok := ctx.Delegations.Get(d.String()); ok {
		val.Amount = val.Amount.Add(d.Amount)
	} else {
		ctx.Delegations.Set(d.String(), &d)
	}
}

func (ctx *Context) AddRedelegation(r storage.Redelegation) {
	ctx.Redelegations = append(ctx.Redelegations, r)
}

func (ctx *Context) AddUndelegation(u storage.Undelegation) {
	ctx.Undelegations = append(ctx.Undelegations, u)
}

func (ctx *Context) AddCancelUndelegation(u storage.Undelegation) {
	ctx.CancelUnbonding = append(ctx.CancelUnbonding, u)
}

func (ctx *Context) AddJail(jail storage.Jail) {
	if j, ok := ctx.Jails.Get(jail.Validator.ConsAddress); ok {
		if jail.Reason != "" {
			j.Reason = jail.Reason
		}
		if !jail.Burned.IsZero() {
			j.Validator.Stake = j.Validator.Stake.Sub(jail.Burned)
			j.Burned = j.Burned.Add(jail.Burned)
		}
		if jail.Validator.Jailed != nil {
			j.Validator.Jailed = jail.Validator.Jailed
		}
	} else {
		ctx.Jails.Set(jail.Validator.ConsAddress, &jail)
	}
}

func (ctx *Context) AddStakingLog(l storage.StakingLog) {
	ctx.StakingLogs = append(ctx.StakingLogs, l)
}
