// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package context

import (
	"encoding/hex"
	"fmt"
	"slices"
	"strconv"
	"sync/atomic"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	fibreTypes "github.com/celestiaorg/celestia-app/v10/x/fibre/types"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

type Context struct {
	Validators        *sdkSync.Map[string, *storage.Validator]
	Addresses         *sdkSync.Map[string, *storage.Address]
	Delegations       *sdkSync.Map[string, *storage.Delegation]
	Jails             *sdkSync.Map[string, *storage.Jail]
	Proposals         *sdkSync.Map[uint64, *storage.Proposal]
	Constants         *sdkSync.Map[string, *storage.Constant]
	Igps              *sdkSync.Map[string, *storage.HLIGP]
	IgpConfigs        *sdkSync.Map[string, *storage.HLIGPConfig]
	ZkISMs            *sdkSync.Map[string, *storage.ZkISM]
	Upgrades          *sdkSync.Map[uint64, *storage.Upgrade]
	Grants            *sdkSync.Map[string, *storage.Grant]
	IbcClients        *sdkSync.Map[string, *storage.IbcClient]
	IbcConnections    *sdkSync.Map[string, *storage.IbcConnection]
	IbcChannels       *sdkSync.Map[string, *storage.IbcChannel]
	HlMailboxes       *sdkSync.Map[uint64, *storage.HLMailbox]
	HlTokens          *sdkSync.Map[string, *storage.HLToken]
	Namespaces        *sdkSync.Map[string, *storage.Namespace]
	NamespaceMessages *sdkSync.Map[string, *storage.NamespaceMessage]
	AddressMessages   *sdkSync.Map[string, *storage.MsgAddress]

	Messages            []*storage.Message
	Events              []storage.Event
	Redelegations       []storage.Redelegation
	Undelegations       []storage.Undelegation
	CancelUnbonding     []storage.Undelegation
	StakingLogs         []storage.StakingLog
	Votes               []*storage.Vote
	VestingAccounts     []*storage.VestingAccount
	Forwardings         []*storage.Forwarding
	ZkIsmUpdates        []*storage.ZkISMUpdate
	ZkIsmMessages       []*storage.ZkISMMessage
	HlTransfers         []*storage.HLTransfer
	IbcTransfers        []*storage.IbcTransfer
	RecoveredIbcClients []string
	BlobLogs            []*storage.BlobLog
	Signals             []*storage.SignalVersion
	DenomMetadata       []storage.DenomMetadata

	Block         *storage.Block
	TryUpgrade    *storage.Upgrade
	TxEventsCount int

	msgCounter       *atomic.Int64
	ibcTransferByMsg map[uint64]*storage.IbcTransfer
}

func NewContext() *Context {
	return &Context{
		Validators:        sdkSync.NewMap[string, *storage.Validator](),
		Addresses:         sdkSync.NewMap[string, *storage.Address](),
		Delegations:       sdkSync.NewMap[string, *storage.Delegation](),
		Jails:             sdkSync.NewMap[string, *storage.Jail](),
		Proposals:         sdkSync.NewMap[uint64, *storage.Proposal](),
		Constants:         sdkSync.NewMap[string, *storage.Constant](),
		Igps:              sdkSync.NewMap[string, *storage.HLIGP](),
		IgpConfigs:        sdkSync.NewMap[string, *storage.HLIGPConfig](),
		Upgrades:          sdkSync.NewMap[uint64, *storage.Upgrade](),
		ZkISMs:            sdkSync.NewMap[string, *storage.ZkISM](),
		Grants:            sdkSync.NewMap[string, *storage.Grant](),
		IbcClients:        sdkSync.NewMap[string, *storage.IbcClient](),
		IbcConnections:    sdkSync.NewMap[string, *storage.IbcConnection](),
		IbcChannels:       sdkSync.NewMap[string, *storage.IbcChannel](),
		HlMailboxes:       sdkSync.NewMap[uint64, *storage.HLMailbox](),
		HlTokens:          sdkSync.NewMap[string, *storage.HLToken](),
		Namespaces:        sdkSync.NewMap[string, *storage.Namespace](),
		NamespaceMessages: sdkSync.NewMap[string, *storage.NamespaceMessage](),
		AddressMessages:   sdkSync.NewMap[string, *storage.MsgAddress](),

		Messages:            make([]*storage.Message, 0, 100),
		Events:              make([]storage.Event, 0, 1000),
		Redelegations:       make([]storage.Redelegation, 0),
		Undelegations:       make([]storage.Undelegation, 0),
		CancelUnbonding:     make([]storage.Undelegation, 0),
		StakingLogs:         make([]storage.StakingLog, 0),
		Votes:               make([]*storage.Vote, 0),
		VestingAccounts:     make([]*storage.VestingAccount, 0),
		Forwardings:         make([]*storage.Forwarding, 0),
		ZkIsmUpdates:        make([]*storage.ZkISMUpdate, 0),
		ZkIsmMessages:       make([]*storage.ZkISMMessage, 0),
		HlTransfers:         make([]*storage.HLTransfer, 0),
		IbcTransfers:        make([]*storage.IbcTransfer, 0),
		RecoveredIbcClients: make([]string, 0),
		Signals:             make([]*storage.SignalVersion, 0),
		DenomMetadata:       make([]storage.DenomMetadata, 0),

		msgCounter:       new(atomic.Int64),
		ibcTransferByMsg: make(map[uint64]*storage.IbcTransfer),
	}
}

func (ctx *Context) AddAddress(address *storage.Address) error {
	if address == nil {
		return nil
	}
	if addr, ok := ctx.Addresses.Get(address.String()); ok {
		for i := range address.Balances {
			found := false
			for j := range addr.Balances {
				if addr.Balances[j].Currency == address.Balances[i].Currency {
					if !address.Balances[i].Spendable.IsZero() {
						addr.Balances[j].Spendable = addr.Balances[j].Spendable.Add(address.Balances[i].Spendable)
					}
					if !address.Balances[i].Delegated.IsZero() {
						addr.Balances[j].Delegated = addr.Balances[j].Delegated.Add(address.Balances[i].Delegated)
					}
					if !address.Balances[i].Unbonding.IsZero() {
						addr.Balances[j].Unbonding = addr.Balances[j].Unbonding.Add(address.Balances[i].Unbonding)
					}
					found = true
					break
				}
			}
			if !found {
				addr.Balances = append(addr.Balances, address.Balances[i])
			}
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
		if validator.FibreHost != nil {
			val.FibreHost = validator.FibreHost
		}
		if validator.FibreHostHeight != nil {
			val.FibreHostHeight = validator.FibreHostHeight
		}
	} else {
		ctx.Validators.Set(validator.Address, &validator)
	}
}

func (ctx *Context) AddSupply(data map[string]string) {
	coin, err := decoder.CoinFromMap(data, "amount")
	if err == nil {
		if coin.GetDenom() == currency.DefaultCurrency {
			amount := types.NumericFromBigInt(coin.Amount.BigInt(), 0)
			ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Add(amount)
		}
	} else {
		amount := decoder.NumericFromMap(data, "amount")
		ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Add(amount)
	}
}

func (ctx *Context) SubSupply(data map[string]string) {
	coin, err := decoder.CoinFromMap(data, "amount")
	if err == nil {
		if coin.GetDenom() == currency.DefaultCurrency {
			amount := types.NumericFromBigInt(coin.Amount.BigInt(), 0)
			ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Sub(amount)
		}
	} else {
		amount := decoder.NumericFromMap(data, "amount")
		ctx.Block.Stats.SupplyChange = ctx.Block.Stats.SupplyChange.Sub(amount)
	}
}

func (ctx *Context) SetInflation(data map[string]string) {
	ctx.Block.Stats.InflationRate = decoder.NumericFromMap(data, "inflation_rate")
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

// AddFibreParams records the fibre module params as constants. Durations are
// stored in nanoseconds, like the other duration constants.
func (ctx *Context) AddFibreParams(params fibreTypes.Params) {
	ctx.AddConstant(types.ModuleNameFibre, "withdrawal_delay", strconv.FormatInt(params.WithdrawalDelay.Nanoseconds(), 10))
	ctx.AddConstant(types.ModuleNameFibre, "payment_promise_timeout", strconv.FormatInt(params.PaymentPromiseTimeout.Nanoseconds(), 10))
	ctx.AddConstant(types.ModuleNameFibre, "payment_promise_height_window", strconv.FormatUint(params.PaymentPromiseHeightWindow, 10))
	ctx.AddConstant(types.ModuleNameFibre, "shard_retention", strconv.FormatInt(params.ShardRetention.Nanoseconds(), 10))
	ctx.AddConstant(types.ModuleNameFibre, "full_stake_storage_budget", strconv.FormatUint(params.FullStakeStorageBudget, 10))
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
		if client.FrozenRevisionHeight > 0 {
			item.FrozenRevisionHeight = client.FrozenRevisionHeight
		}
		if client.FrozenRevisionNumber > 0 {
			item.FrozenRevisionNumber = client.FrozenRevisionNumber
		}
		// chain keeps max(LatestHeight) compared as (revision, height): lower updates don't move it back
		if client.LatestRevisionCompare(item) == 1 {
			item.LatestRevisionHeight = client.LatestRevisionHeight
			item.LatestRevisionNumber = client.LatestRevisionNumber
		}
		if !client.UpdatedAt.IsZero() {
			item.UpdatedAt = client.UpdatedAt
		}
		if client.ChainId != "" {
			item.ChainId = client.ChainId
		}
		if client.UnbondingPeriod > 0 {
			item.UnbondingPeriod = client.UnbondingPeriod
		}
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

// AddIbcChannelTransfer counts a transfer in its channel stats; call only once the transfer's success is confirmed by events
func (ctx *Context) AddIbcChannelTransfer(transfer *storage.IbcTransfer) {
	channel := &storage.IbcChannel{
		Id:             transfer.ChannelId,
		TransfersCount: 1,
		Status:         types.IbcChannelStatusInitialization,
	}
	if transfer.Receiver != nil {
		channel.Received = channel.Received.Add(transfer.Amount)
	}
	if transfer.Sender != nil {
		channel.Sent = channel.Sent.Add(transfer.Amount)
	}
	ctx.AddIbcChannel(channel)
}

// AddRecoveredIbcClient marks a subject client recovered by governance: unfrozen and replaced by its substitute
func (ctx *Context) AddRecoveredIbcClient(id string) {
	if id == "" || slices.Contains(ctx.RecoveredIbcClients, id) {
		return
	}
	ctx.RecoveredIbcClients = append(ctx.RecoveredIbcClients, id)
}

// AddIbcTransfer binds the transfer to its message so event handlers never touch another message's transfer
func (ctx *Context) AddIbcTransfer(msgId uint64, transfer *storage.IbcTransfer) {
	if ctx.ibcTransferByMsg == nil {
		ctx.ibcTransferByMsg = make(map[uint64]*storage.IbcTransfer)
	}
	ctx.ibcTransferByMsg[msgId] = transfer
	ctx.IbcTransfers = append(ctx.IbcTransfers, transfer)
}

func (ctx *Context) IbcTransferByMsg(msgId uint64) *storage.IbcTransfer {
	return ctx.ibcTransferByMsg[msgId]
}

func (ctx *Context) RemoveIbcTransfer(msgId uint64) {
	transfer, ok := ctx.ibcTransferByMsg[msgId]
	if !ok {
		return
	}
	delete(ctx.ibcTransferByMsg, msgId)
	ctx.IbcTransfers = slices.DeleteFunc(ctx.IbcTransfers, func(t *storage.IbcTransfer) bool {
		return t == transfer
	})
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
		ns.FibreSize += namespace.FibreSize
		ns.PffCount += namespace.PffCount
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

func (ctx *Context) AddDenomMetadata(metadata ...storage.DenomMetadata) {
	ctx.DenomMetadata = append(ctx.DenomMetadata, metadata...)
}
