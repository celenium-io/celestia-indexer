// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type CoinReceived struct {
	Amount   *types.Coin
	Receiver string
}

func NewCoinReceived(m map[string]any) (body CoinReceived, err error) {
	body.Receiver = decoder.StringFromMap(m, "receiver")
	if body.Receiver == "" {
		err = errors.Errorf("receiver key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type CoinSpent struct {
	Amount  *types.Coin
	Spender string
}

func NewCoinSpent(m map[string]any) (body CoinSpent, err error) {
	body.Spender = decoder.StringFromMap(m, "spender")
	if body.Spender == "" {
		err = errors.Errorf("spender key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type CompleteRedelegation struct {
	Amount        *types.Coin
	Delegator     string
	DestValidator string
	SrcValidator  string
}

func NewCompleteRedelegation(m map[string]any) (body CompleteRedelegation, err error) {
	body.Delegator = decoder.StringFromMap(m, "delegator")
	if body.Delegator == "" {
		err = errors.Errorf("delegator key not found in %##v", m)
		return
	}
	body.DestValidator = decoder.StringFromMap(m, "destination_validator")
	if body.DestValidator == "" {
		err = errors.Errorf("destination_validator key not found in %##v", m)
		return
	}
	body.SrcValidator = decoder.StringFromMap(m, "source_validator")
	if body.SrcValidator == "" {
		err = errors.Errorf("source_validator key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type CompleteUnbonding struct {
	Amount    *types.Coin
	Delegator string
	Validator string
}

func NewCompleteUnbonding(m map[string]any) (body CompleteUnbonding, err error) {
	body.Delegator = decoder.StringFromMap(m, "delegator")
	if body.Delegator == "" {
		err = errors.Errorf("delegator key not found in %##v", m)
		return
	}
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type Commission struct {
	Amount    decimal.Decimal
	Validator string
}

func NewCommission(m map[string]any) (body Commission, err error) {
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.Amount = decoder.AmountFromMap(m, "amount")
	return
}

type Rewards struct {
	Amount    decimal.Decimal
	Validator string
}

func NewRewards(m map[string]any) (body Rewards, err error) {
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.Amount = decoder.AmountFromMap(m, "amount")
	return
}

type WithdrawRewards struct {
	Amount    *types.Coin
	Validator string
	Delegator string
}

func NewWithdrawRewards(m map[string]any) (body WithdrawRewards, err error) {
	body.Delegator = decoder.StringFromMap(m, "delegator")
	if body.Delegator == "" {
		err = errors.Errorf("delegator key not found in %##v", m)
		return
	}
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type WithdrawCommission struct {
	Amount *types.Coin
}

func NewWithdrawCommission(m map[string]any) (body WithdrawCommission, err error) {
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type Redelegate struct {
	Amount         *types.Coin
	DestValidator  string
	SrcValidator   string
	CompletionTime time.Time
}

func NewRedelegate(m map[string]any) (body Redelegate, err error) {
	body.CompletionTime, err = decoder.TimeFromMap(m, "completion_time")
	if err != nil {
		err = errors.Wrap(err, "completion_time")
		return
	}
	body.DestValidator = decoder.StringFromMap(m, "destination_validator")
	if body.DestValidator == "" {
		err = errors.Errorf("destination_validator key not found in %##v", m)
		return
	}
	body.SrcValidator = decoder.StringFromMap(m, "source_validator")
	if body.SrcValidator == "" {
		err = errors.Errorf("source_validator key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type Unbond struct {
	Amount         *types.Coin
	Validator      string
	CompletionTime time.Time
}

func NewUnbond(m map[string]any) (body Unbond, err error) {
	body.CompletionTime, err = decoder.TimeFromMap(m, "completion_time")
	if err != nil {
		err = errors.Wrap(err, "completion_time")
		return
	}
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type Delegate struct {
	Amount    *types.Coin
	NewShares decimal.Decimal
	Validator string
}

func NewDelegate(m map[string]any) (body Delegate, err error) {
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.NewShares = decoder.DecimalFromMap(m, "new_shares")
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type CancelUnbondingDelegation struct {
	Amount         *types.Coin
	Validator      string
	Delegator      string
	CreationHeight int64
}

func NewCancelUnbondingDelegation(m map[string]any) (body CancelUnbondingDelegation, err error) {
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.Delegator = decoder.StringFromMap(m, "delegator")
	if body.Delegator == "" {
		err = errors.Errorf("delegator key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	if err != nil {
		return
	}
	body.CreationHeight, err = decoder.Int64FromMap(m, "creation_height")
	return
}

type Slash struct {
	Power       decimal.Decimal
	Jailed      string
	Reason      string
	Address     string
	BurnedCoins decimal.Decimal
}

func NewSlash(m map[string]any) (body Slash, err error) {
	body.Power = decoder.DecimalFromMap(m, "power")
	body.BurnedCoins = decoder.DecimalFromMap(m, "burned_coins")
	body.Reason = decoder.StringFromMap(m, "reason")
	body.Jailed = decoder.StringFromMap(m, "jailed")
	body.Address = decoder.StringFromMap(m, "address")
	return
}

type ProposalStatus struct {
	Id     uint64
	Result string
}

func NewProposalStatus(m map[string]any) (body ProposalStatus, err error) {
	body.Result = decoder.StringFromMap(m, "proposal_result")
	body.Id, err = decoder.Uint64FromMap(m, "proposal_id")
	return
}

type UpdateClient struct {
	Id              string
	Type            string
	ConsensusHeight uint64
	Revision        uint64
}

func NewUpdateClient(m map[string]any) (cc UpdateClient, err error) {
	cc.Id = decoder.StringFromMap(m, "client_id")
	cc.Type = decoder.StringFromMap(m, "client_type")
	revision, height, err := decoder.RevisionHeightFromMap(m, "consensus_height")
	if err != nil {
		return cc, errors.Wrap(err, "consensus_height")
	}
	cc.ConsensusHeight = height
	cc.Revision = revision
	return
}

type ConnectionChange struct {
	ClientId                 string
	ConnectionId             string
	CounterpartyClientId     string
	CounterpartyConnectionId string
}

func NewConnectionOpen(m map[string]any) (cc ConnectionChange) {
	cc.ClientId = decoder.StringFromMap(m, "client_id")
	cc.ConnectionId = decoder.StringFromMap(m, "connection_id")
	cc.CounterpartyClientId = decoder.StringFromMap(m, "counterparty_client_id")
	cc.CounterpartyConnectionId = decoder.StringFromMap(m, "counterparty_connection_id")
	return
}

type ChannelChange struct {
	ChannelId             string
	ConnectionId          string
	CounterpartyChannelId string
	CounterpartyPortId    string
	PortId                string
}

func NewChannelChange(m map[string]any) (cc ChannelChange) {
	cc.ChannelId = decoder.StringFromMap(m, "channel_id")
	cc.ConnectionId = decoder.StringFromMap(m, "connection_id")
	cc.CounterpartyChannelId = decoder.StringFromMap(m, "counterparty_channel_id")
	cc.CounterpartyPortId = decoder.StringFromMap(m, "counterparty_port_id")
	cc.PortId = decoder.StringFromMap(m, "port_id")
	return
}

type FungibleTokenPacket struct {
	Amount   decimal.Decimal
	Denom    string
	Memo     string
	Module   string
	Receiver string
	Sender   string
	Success  string
}

func NewFungibleTokenPacket(m map[string]any) (ftp FungibleTokenPacket) {
	ftp.Amount = decoder.DecimalFromMap(m, "amount")
	ftp.Denom = decoder.StringFromMap(m, "denom")
	ftp.Memo = decoder.StringFromMap(m, "memo")
	ftp.Module = decoder.StringFromMap(m, "module")
	ftp.Receiver = decoder.StringFromMap(m, "receiver")
	ftp.Sender = decoder.StringFromMap(m, "sender")
	ftp.Success = decoder.StringFromMap(m, "success")
	return
}

type RecvPacket struct {
	Ordering      string
	Connection    string
	Data          string
	DstChannel    string
	DstPort       string
	SrcChannel    string
	SrcPort       string
	Sequence      uint64
	Timeout       time.Time
	TimeoutHeight uint64
}

func NewRecvPacket(m map[string]any) (rp RecvPacket, err error) {
	rp.Ordering = decoder.StringFromMap(m, "packet_channel_ordering")
	rp.Connection = decoder.StringFromMap(m, "packet_connection")
	rp.Data = decoder.StringFromMap(m, "packet_data")
	rp.DstChannel = decoder.StringFromMap(m, "packet_dst_channel")
	rp.DstPort = decoder.StringFromMap(m, "packet_dst_port")
	rp.Sequence, err = decoder.Uint64FromMap(m, "packet_sequence")
	if err != nil {
		return rp, errors.Wrap(err, "packet_sequence")
	}
	rp.SrcChannel = decoder.StringFromMap(m, "packet_src_channel")
	rp.SrcPort = decoder.StringFromMap(m, "packet_src_port")
	_, height, err := decoder.RevisionHeightFromMap(m, "packet_timeout_height")
	if err != nil {
		return rp, errors.Wrap(err, "packet_timeout_height")
	}
	rp.TimeoutHeight = height
	rp.Timeout = decoder.UnixNanoFromMap(m, "packet_timeout_timestamp")
	return
}
