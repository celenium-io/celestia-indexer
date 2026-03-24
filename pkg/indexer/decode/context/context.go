// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package context

import (
	"encoding/hex"
	"fmt"
	"sync/atomic"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Context struct {
	Validators        *sync.Map[string, *storage.Validator]
	Addresses         *sync.Map[string, *storage.Address]
	Delegations       *sync.Map[string, *storage.Delegation]
	Jails             *sync.Map[string, *storage.Jail]
	Proposals         *sync.Map[uint64, *storage.Proposal]
	Constants         *sync.Map[string, *storage.Constant]
	Igps              *sync.Map[string, *storage.HLIGP]
	IgpConfigs        *sync.Map[string, *storage.HLIGPConfig]
	ZkISMs            *sync.Map[string, *storage.ZkISM]
	Upgrades          *sync.Map[uint64, *storage.Upgrade]
	Grants            *sync.Map[string, *storage.Grant]
	IbcClients        *sync.Map[string, *storage.IbcClient]
	IbcConnections    *sync.Map[string, *storage.IbcConnection]
	IbcChannels       *sync.Map[string, *storage.IbcChannel]
	HlMailboxes       *sync.Map[uint64, *storage.HLMailbox]
	HlTokens          *sync.Map[string, *storage.HLToken]
	Namespaces        *sync.Map[string, *storage.Namespace]
	NamespaceMessages *sync.Map[string, *storage.NamespaceMessage]
	AddressMessages   *sync.Map[string, *storage.MsgAddress]

	Messages        []*storage.Message
	Events          []storage.Event
	Redelegations   []storage.Redelegation
	Undelegations   []storage.Undelegation
	CancelUnbonding []storage.Undelegation
	StakingLogs     []storage.StakingLog
	Votes           []*storage.Vote
	VestingAccounts []*storage.VestingAccount
	Forwardings     []*storage.Forwarding
	ZkIsmUpdates    []*storage.ZkISMUpdate
	ZkIsmMessages   []*storage.ZkISMMessage
	HlTransfers     []*storage.HLTransfer
	IbcTransfers    []*storage.IbcTransfer
	BlobLogs        []*storage.BlobLog
	Signals         []*storage.SignalVersion

	Block         *storage.Block
	TryUpgrade    *storage.Upgrade
	TxEventsCount int

	msgCounter *atomic.Int64
}

func NewContext() *Context {
	return &Context{
		Validators:        sync.NewMap[string, *storage.Validator](),
		Addresses:         sync.NewMap[string, *storage.Address](),
		Delegations:       sync.NewMap[string, *storage.Delegation](),
		Jails:             sync.NewMap[string, *storage.Jail](),
		Proposals:         sync.NewMap[uint64, *storage.Proposal](),
		Constants:         sync.NewMap[string, *storage.Constant](),
		Igps:              sync.NewMap[string, *storage.HLIGP](),
		IgpConfigs:        sync.NewMap[string, *storage.HLIGPConfig](),
		Upgrades:          sync.NewMap[uint64, *storage.Upgrade](),
		ZkISMs:            sync.NewMap[string, *storage.ZkISM](),
		Grants:            sync.NewMap[string, *storage.Grant](),
		IbcClients:        sync.NewMap[string, *storage.IbcClient](),
		IbcConnections:    sync.NewMap[string, *storage.IbcConnection](),
		IbcChannels:       sync.NewMap[string, *storage.IbcChannel](),
		HlMailboxes:       sync.NewMap[uint64, *storage.HLMailbox](),
		HlTokens:          sync.NewMap[string, *storage.HLToken](),
		Namespaces:        sync.NewMap[string, *storage.Namespace](),
		NamespaceMessages: sync.NewMap[string, *storage.NamespaceMessage](),
		AddressMessages:   sync.NewMap[string, *storage.MsgAddress](),

		Messages:        make([]*storage.Message, 0, 100),
		Events:          make([]storage.Event, 0, 1000),
		Redelegations:   make([]storage.Redelegation, 0),
		Undelegations:   make([]storage.Undelegation, 0),
		CancelUnbonding: make([]storage.Undelegation, 0),
		StakingLogs:     make([]storage.StakingLog, 0),
		Votes:           make([]*storage.Vote, 0),
		VestingAccounts: make([]*storage.VestingAccount, 0),
		Forwardings:     make([]*storage.Forwarding, 0),
		ZkIsmUpdates:    make([]*storage.ZkISMUpdate, 0),
		ZkIsmMessages:   make([]*storage.ZkISMMessage, 0),
		HlTransfers:     make([]*storage.HLTransfer, 0),
		IbcTransfers:    make([]*storage.IbcTransfer, 0),
		Signals:         make([]*storage.SignalVersion, 0),

		msgCounter: new(atomic.Int64),
	}
}

func (ctx *Context) AddAddress(address *storage.Address) error {
	if address == nil {
		return nil
	}
	if addr, ok := ctx.Addresses.Get(address.String()); ok {
		if address.Balance.Currency == currency.DefaultCurrency { // hotfix for balance updates with non-default currency
			addr.Balance.Spendable = addr.Balance.Spendable.Add(address.Balance.Spendable)
			addr.Balance.Delegated = addr.Balance.Delegated.Add(address.Balance.Delegated)
			addr.Balance.Unbonding = addr.Balance.Unbonding.Add(address.Balance.Unbonding)
		}

		if address.IsForwarding {
			addr.IsForwarding = true
		}
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
		if validator.Version > 0 {
			val.Version = validator.Version
		}
		if validator.MessagesCount > 0 {
			val.MessagesCount += validator.MessagesCount
		}
	} else {
		ctx.Validators.Set(validator.Address, &validator)
	}
}

func (ctx *Context) AddSupply(data map[string]string) {
	coin, err := decoder.CoinFromMap(data, "amount")
	if err == nil {
		if coin.GetDenom() == currency.DefaultCurrency {
			amount := decimal.NewFromBigInt(coin.Amount.BigInt(), 0)
			ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Add(amount)
		}
	} else {
		amount := decoder.DecimalFromMap(data, "amount")
		ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Add(amount)
	}
}

func (ctx *Context) SubSupply(data map[string]string) {
	coin, err := decoder.CoinFromMap(data, "amount")
	if err == nil {
		if coin.GetDenom() == currency.DefaultCurrency {
			amount := decimal.NewFromBigInt(coin.Amount.BigInt(), 0)
			ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Sub(amount)
		}
	} else {
		amount := decoder.DecimalFromMap(data, "amount")
		ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Sub(amount)
	}
}

func (ctx *Context) SetInflation(data map[string]string) {
	ctx.Block.Stats.InflationRate = decoder.DecimalFromMap(data, "inflation_rate")
}

func (ctx *Context) AddDelegation(d storage.Delegation) {
	if val, ok := ctx.Delegations.Get(d.String()); ok {
		val.Amount = val.Amount.Add(d.Amount)
	} else {
		ctx.Delegations.Set(d.String(), &d)
	}
}

func (ctx *Context) AddMessage(msg *storage.Message) {
	ctx.Messages = append(ctx.Messages, msg)
}

func (ctx *Context) AddEvents(events ...storage.Event) {
	ctx.Events = append(ctx.Events, events...)
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

func (ctx *Context) AddProposal(proposal *storage.Proposal) {
	if p, ok := ctx.Proposals.Get(proposal.Id); ok {
		if proposal.Status.GreaterThan(p.Status) {
			p.Status = proposal.Status
		}
		if proposal.ActivationTime != nil {
			p.ActivationTime = proposal.ActivationTime
		}
		if proposal.Deposit.IsPositive() {
			p.Deposit = p.Deposit.Add(proposal.Deposit)
		}
	} else {
		ctx.Proposals.Set(proposal.Id, proposal)
	}
}

func (ctx *Context) AddVote(vote *storage.Vote) {
	ctx.Votes = append(ctx.Votes, vote)
}

func (ctx *Context) AddConstant(module types.ModuleName, name, value string) {
	key := fmt.Sprintf("%s_%s", module, name)
	ctx.Constants.Set(key, &storage.Constant{
		Module: module,
		Name:   name,
		Value:  value,
	})
}

func (ctx *Context) AddIgp(igpId string, igp *storage.HLIGP) {
	if val, ok := ctx.Igps.Get(igpId); ok {
		val.Owner = igp.Owner
	} else {
		ctx.Igps.Set(igpId, igp)
	}
}

func (ctx *Context) AddIgpConfig(igpId string, config *storage.HLIGPConfig) {
	if val, ok := ctx.Igps.Get(igpId); ok {
		val.Configs = append(val.Configs, config)
	} else {
		ctx.IgpConfigs.Set(igpId, config)
	}
}

func (ctx *Context) AddUpgrade(upgrade storage.Upgrade) {
	if val, ok := ctx.Upgrades.Get(upgrade.Version); ok {
		val.SignalsCount += upgrade.SignalsCount
	} else {
		ctx.Upgrades.Set(upgrade.Version, &upgrade)
	}
}

func (ctx *Context) AddSignal(signal *storage.SignalVersion) {
	ctx.Signals = append(ctx.Signals, signal)
}

func (ctx *Context) AddZkISM(ism *storage.ZkISM) {
	key := hex.EncodeToString(ism.ExternalId)
	if value, ok := ctx.ZkISMs.Get(key); ok {
		value.State = ism.State
	} else {
		ctx.ZkISMs.Set(key, ism)
	}
}

func (ctx *Context) GetMsgPosition() int64 {
	return ctx.msgCounter.Add(1) - 1
}

func (ctx *Context) AddVestingAccount(acc *storage.VestingAccount) {
	ctx.VestingAccounts = append(ctx.VestingAccounts, acc)
}

func (ctx *Context) AddGrants(grants ...*storage.Grant) {
	for i := range grants {
		ctx.Grants.Set(grants[i].String(), grants[i])
	}
}

func (ctx *Context) AddIbcClient(client *storage.IbcClient) {
	if item, ok := ctx.IbcClients.Get(client.Id); ok {
		item.ConnectionCount += client.ConnectionCount
	} else {
		ctx.IbcClients.Set(client.Id, client)
	}
}

func (ctx *Context) AddIbcConnection(conn *storage.IbcConnection) {
	if item, ok := ctx.IbcConnections.Get(conn.ConnectionId); ok {
		item.ChannelsCount += conn.ChannelsCount
	} else {
		ctx.IbcConnections.Set(conn.ConnectionId, conn)
	}
}

func (ctx *Context) AddIbcChannel(channel *storage.IbcChannel) {
	if ch, ok := ctx.IbcChannels.Get(channel.Id); ok {
		ch.Received = ch.Received.Add(channel.Received)
		ch.Sent = ch.Sent.Add(channel.Sent)
		ch.TransfersCount += channel.TransfersCount
	} else {
		ctx.IbcChannels.Set(channel.Id, channel)
	}
}

func (ctx *Context) DeleteIbcChannel(chanId string) {
	ctx.IbcChannels.Delete(chanId)
}

func (ctx *Context) AddIbcTransfer(transfer *storage.IbcTransfer) {
	ctx.IbcTransfers = append(ctx.IbcTransfers, transfer)
}

func (ctx *Context) GetLastIbcTransfer() *storage.IbcTransfer {
	if len(ctx.IbcTransfers) == 0 {
		return nil
	}
	return ctx.IbcTransfers[len(ctx.IbcTransfers)-1]
}

func (ctx *Context) RemoveLastIbcTransfer() {
	if len(ctx.IbcTransfers) == 0 {
		return
	}
	ctx.IbcTransfers = ctx.IbcTransfers[:len(ctx.IbcTransfers)-1]
}

func (ctx *Context) AddForwarding(fwd *storage.Forwarding) {
	ctx.Forwardings = append(ctx.Forwardings, fwd)
}

func (ctx *Context) AddZkIsmUpdate(upd *storage.ZkISMUpdate) {
	ctx.ZkIsmUpdates = append(ctx.ZkIsmUpdates, upd)
}

func (ctx *Context) AddZkIsmMessage(msg *storage.ZkISMMessage) {
	ctx.ZkIsmMessages = append(ctx.ZkIsmMessages, msg)
}

func (ctx *Context) AddHlMailbox(mailbox *storage.HLMailbox) {
	if item, ok := ctx.HlMailboxes.Get(mailbox.InternalId); ok {
		if mailbox.ReceivedMessages > 0 {
			item.ReceivedMessages += mailbox.ReceivedMessages
		}
		if mailbox.SentMessages > 0 {
			item.SentMessages += mailbox.SentMessages
		}
		if len(mailbox.DefaultHook) > 0 {
			item.DefaultHook = mailbox.DefaultHook
		}
		if len(mailbox.RequiredHook) > 0 {
			item.RequiredHook = mailbox.RequiredHook
		}
		if len(mailbox.DefaultIsm) > 0 {
			item.DefaultIsm = mailbox.DefaultIsm
		}
		if mailbox.Owner != nil {
			item.Owner = mailbox.Owner
		}
	} else {
		ctx.HlMailboxes.Set(mailbox.InternalId, mailbox)
	}
}

func (ctx *Context) AddHlToken(token *storage.HLToken) {
	key := token.String()
	if item, ok := ctx.HlTokens.Get(key); ok {
		item.ReceiveTransfers += token.ReceiveTransfers
		item.SentTransfers += token.SentTransfers
		if !token.Sent.IsZero() {
			item.Sent = item.Sent.Add(token.Sent)
		}
		if !token.Received.IsZero() {
			item.Received = item.Received.Add(token.Received)
		}
	} else {
		ctx.HlTokens.Set(key, token)
	}
}

func (ctx *Context) AddHlTransfer(transfer *storage.HLTransfer) {
	ctx.HlTransfers = append(ctx.HlTransfers, transfer)

	if transfer.Token != nil {
		if transfer.Mailbox != nil {
			ctx.AddHlMailbox(&storage.HLMailbox{
				InternalId:       transfer.Mailbox.InternalId,
				ReceivedMessages: transfer.Token.ReceiveTransfers,
				SentMessages:     transfer.Token.SentTransfers,
			})
		}
		ctx.AddHlToken(transfer.Token)
	}
}

func (ctx *Context) AddNamespace(namespace *storage.Namespace) *storage.Namespace {
	key := namespace.String()
	if ns, ok := ctx.Namespaces.Get(key); ok {
		ns.PfbCount += namespace.PfbCount
		ns.Size += namespace.Size
		ns.BlobsCount += namespace.BlobsCount
		return ns
	}
	ctx.Namespaces.Set(key, namespace)
	return namespace
}

func (ctx *Context) AddNamespaceMessage(msg *storage.NamespaceMessage) {
	if msg.Namespace == nil {
		return
	}
	key := fmt.Sprintf("%d-%s", msg.MsgId, msg.Namespace.String())
	ctx.NamespaceMessages.Set(key, msg)
}

func (ctx *Context) AddBlobLogs(logs ...*storage.BlobLog) {
	ctx.BlobLogs = append(ctx.BlobLogs, logs...)
}

func (ctx *Context) AddAddressMessage(msg *storage.MsgAddress) {
	ctx.AddressMessages.Set(msg.String(), msg)
}
