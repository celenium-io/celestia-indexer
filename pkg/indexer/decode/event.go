// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"time"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
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
	Error    string
}

func NewFungibleTokenPacket(m map[string]any) (ftp FungibleTokenPacket) {
	ftp.Amount = decoder.DecimalFromMap(m, "amount")
	ftp.Denom = decoder.StringFromMap(m, "denom")
	ftp.Memo = decoder.StringFromMap(m, "memo")
	ftp.Module = decoder.StringFromMap(m, "module")
	ftp.Receiver = decoder.StringFromMap(m, "receiver")
	ftp.Sender = decoder.StringFromMap(m, "sender")
	ftp.Success = decoder.StringFromMap(m, "success")
	ftp.Error = decoder.StringFromMap(m, "error")
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

type CreateMailbox struct {
	MailboxId    string
	Owner        string
	DefaultIsm   string
	DefaultHook  string
	RequiredHook string
	LocalDomain  uint64
}

func NewCreateMailbox(m map[string]any) (cm CreateMailbox, err error) {
	cm.MailboxId = decoder.StringFromMap(m, "mailbox_id")
	cm.Owner = decoder.StringFromMap(m, "owner")
	cm.DefaultIsm = decoder.StringFromMap(m, "default_ism")
	cm.DefaultHook = decoder.StringFromMap(m, "default_hook")
	cm.RequiredHook = decoder.StringFromMap(m, "required_hook")
	cm.LocalDomain, err = decoder.Uint64FromMap(m, "local_domain")
	if err != nil {
		return cm, errors.Wrap(err, "local_domain")
	}
	return
}

type SetMailbox struct {
	MailboxId         string
	Owner             string
	DefaultIsm        string
	DefaultHook       string
	NewOwner          string
	RenounceOwnership bool
}

func NewSetMailbox(m map[string]any) (sm SetMailbox, err error) {
	sm.MailboxId = decoder.StringFromMap(m, "mailbox_id")
	sm.Owner = decoder.StringFromMap(m, "owner")
	sm.DefaultIsm = decoder.StringFromMap(m, "default_ism")
	sm.DefaultHook = decoder.StringFromMap(m, "default_hook")
	sm.NewOwner = decoder.StringFromMap(m, "new_owner")
	sm.RenounceOwnership, err = decoder.BoolFromMap(m, "renounce_ownership")
	if err != nil {
		return sm, errors.Wrap(err, "renounce_ownership")
	}
	return
}

type HyperlaneProcessEvent struct {
	OriginMailboxId string
	Sender          string
	Recipient       string
	MessageId       string
	Origin          uint64
	Message         *util.HyperlaneMessage
}

func NewHyperlaneProcessEvent(m map[string]any) (hpe HyperlaneProcessEvent, err error) {
	hpe.OriginMailboxId = decoder.StringFromMap(m, "origin_mailbox_id")
	hpe.Sender = decoder.StringFromMap(m, "sender")
	hpe.Recipient = decoder.StringFromMap(m, "recipient")
	hpe.MessageId = decoder.StringFromMap(m, "message_id")
	hpe.Origin, err = decoder.Uint64FromMap(m, "origin")
	if err != nil {
		return hpe, errors.Wrap(err, "origin")
	}
	hpe.Message, err = decoder.HyperlaneMessageFromMap(m, "message")
	if err != nil {
		return hpe, errors.Wrap(err, "message")
	}

	return
}

type HyperlaneDispatchEvent struct {
	OriginMailboxId string
	Sender          string
	Recipient       string
	Destination     uint64
	Message         *util.HyperlaneMessage
}

func NewHyperlaneDispatchEvent(m map[string]any) (hde HyperlaneDispatchEvent, err error) {
	hde.OriginMailboxId = decoder.StringFromMap(m, "origin_mailbox_id")
	hde.Sender = decoder.StringFromMap(m, "sender")
	hde.Recipient = decoder.StringFromMap(m, "recipient")
	hde.Destination, err = decoder.Uint64FromMap(m, "destination")
	if err != nil {
		return hde, errors.Wrap(err, "destination")
	}
	hde.Message, err = decoder.HyperlaneMessageFromMap(m, "message")
	if err != nil {
		return hde, errors.Wrap(err, "message")
	}

	return
}

type CreateCollateralToken struct {
	MailboxId string
	Owner     string
	TokenId   string
	Denom     string
}

func NewCreateCollateralToken(m map[string]any) (cct CreateCollateralToken, err error) {
	cct.MailboxId = decoder.StringFromMap(m, "origin_mailbox")
	cct.Owner = decoder.StringFromMap(m, "owner")
	cct.TokenId = decoder.StringFromMap(m, "token_id")
	cct.Denom = decoder.StringFromMap(m, "origin_denom")
	return
}

type CreateSyntheticToken struct {
	MailboxId string
	Owner     string
	TokenId   string
	Denom     string
}

func NewCreateSyntheticToken(m map[string]any) (cst CreateSyntheticToken, err error) {
	cst.MailboxId = decoder.StringFromMap(m, "origin_mailbox")
	cst.Owner = decoder.StringFromMap(m, "owner")
	cst.TokenId = decoder.StringFromMap(m, "token_id")
	cst.Denom = decoder.StringFromMap(m, "origin_denom")
	return
}

type HyperlaneReceiveTransferEvent struct {
	Amount       decimal.Decimal
	Denom        string
	OriginDomain uint64
	Recipient    string
	Sender       string
	TokenId      string
}

func NewHyperlaneReceiveTransferEvent(m map[string]any) (hrte HyperlaneReceiveTransferEvent, err error) {
	hrte.Sender = decoder.StringFromMap(m, "sender")
	hrte.Recipient = decoder.StringFromMap(m, "recipient")
	hrte.TokenId = decoder.StringFromMap(m, "token_id")
	hrte.OriginDomain, err = decoder.Uint64FromMap(m, "origin_domain")
	if err != nil {
		return hrte, errors.Wrap(err, "origin")
	}
	coin, err := decoder.CoinFromMap(m, "amount")
	if err != nil {
		return hrte, errors.Wrap(err, "amount")
	}
	hrte.Amount = decimal.RequireFromString(coin.Amount.String())
	hrte.Denom = coin.GetDenom()
	return
}

type HyperlaneSendTransferEvent struct {
	Amount            decimal.Decimal
	Denom             string
	DestinationDomain uint64
	Recipient         string
	Sender            string
	TokenId           string
}

func NewHyperlaneSendTransferEvent(m map[string]any) (hste HyperlaneSendTransferEvent, err error) {
	hste.Sender = decoder.StringFromMap(m, "sender")
	hste.Recipient = decoder.StringFromMap(m, "recipient")
	hste.TokenId = decoder.StringFromMap(m, "token_id")
	hste.DestinationDomain, err = decoder.Uint64FromMap(m, "destination_domain")
	if err != nil {
		return hste, errors.Wrap(err, "origin")
	}
	coin, err := decoder.CoinFromMap(m, "amount")
	if err != nil {
		return hste, errors.Wrap(err, "amount")
	}
	hste.Amount = decimal.RequireFromString(coin.Amount.String())
	hste.Denom = coin.GetDenom()
	return
}

type SetToken struct {
	IsmId             string
	TokenId           string
	NewOwner          string
	Owner             string
	RenounceOwnership bool
}

func NewSetToken(m map[string]any) (st SetToken, err error) {
	st.NewOwner = decoder.StringFromMap(m, "new_owner")
	st.Owner = decoder.StringFromMap(m, "owner")
	st.TokenId = decoder.StringFromMap(m, "token_id")
	st.IsmId = decoder.StringFromMap(m, "ism_id")
	st.RenounceOwnership, err = decoder.BoolFromMap(m, "renounce_ownership")
	if err != nil {
		return st, errors.Wrap(err, "renounce_ownership")
	}
	return
}
