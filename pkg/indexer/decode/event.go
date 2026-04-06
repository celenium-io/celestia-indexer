// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

type CoinReceived struct {
	Amount   *types.Coin
	Receiver string
}

func NewCoinReceived(m map[string]string) (body CoinReceived, err error) {
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

func NewCoinSpent(m map[string]string) (body CoinSpent, err error) {
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

func NewCompleteRedelegation(m map[string]string) (body CompleteRedelegation, err error) {
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

func NewCompleteUnbonding(m map[string]string) (body CompleteUnbonding, err error) {
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
	Amount    storageTypes.Numeric
	Validator string
}

func NewCommission(m map[string]string) (body Commission, err error) {
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.NumericAmountFromMap(m, "amount")
	return
}

type Rewards struct {
	Amount    storageTypes.Numeric
	Validator string
}

func NewRewards(m map[string]string) (body Rewards, err error) {
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.Amount, err = decoder.NumericAmountFromMap(m, "amount")
	return
}

type WithdrawRewards struct {
	Amount    *types.Coin
	Validator string
	Delegator string
}

func NewWithdrawRewards(m map[string]string) (body WithdrawRewards, err error) {
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

func NewWithdrawCommission(m map[string]string) (body WithdrawCommission, err error) {
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type Redelegate struct {
	Amount         *types.Coin
	DestValidator  string
	SrcValidator   string
	CompletionTime time.Time
}

func NewRedelegate(m map[string]string) (body Redelegate, err error) {
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

func NewUnbond(m map[string]string) (body Unbond, err error) {
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
	NewShares storageTypes.Numeric
	Validator string
}

func NewDelegate(m map[string]string) (body Delegate, err error) {
	body.Validator = decoder.StringFromMap(m, "validator")
	if body.Validator == "" {
		err = errors.Errorf("validator key not found in %##v", m)
		return
	}
	body.NewShares = decoder.NumericFromMap(m, "new_shares")
	body.Amount, err = decoder.BalanceFromMap(m, "amount")
	return
}

type CancelUnbondingDelegation struct {
	Amount         *types.Coin
	Validator      string
	Delegator      string
	CreationHeight int64
}

func NewCancelUnbondingDelegation(m map[string]string) (body CancelUnbondingDelegation, err error) {
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
	Power       storageTypes.Numeric
	Jailed      string
	Reason      string
	Address     string
	BurnedCoins storageTypes.Numeric
}

func NewSlash(m map[string]string) (body Slash, err error) {
	body.Power = decoder.NumericFromMap(m, "power")
	body.BurnedCoins = decoder.NumericFromMap(m, "burned_coins")
	body.Reason = decoder.StringFromMap(m, "reason")
	body.Jailed = decoder.StringFromMap(m, "jailed")
	body.Address = decoder.StringFromMap(m, "address")
	return
}

type ProposalStatus struct {
	Id     uint64
	Result string
	Log    string
}

func NewProposalStatus(m map[string]string) (body ProposalStatus, err error) {
	body.Result = decoder.StringFromMap(m, "proposal_result")
	body.Log = decoder.StringFromMap(m, "proposal_log")
	body.Id, err = decoder.Uint64FromMap(m, "proposal_id")
	return
}

type UpdateClient struct {
	Id              string
	Type            string
	ConsensusHeight uint64
	Revision        uint64
}

func NewUpdateClient(m map[string]string) (cc UpdateClient, err error) {
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

func NewConnectionOpen(m map[string]string) (cc ConnectionChange) {
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

func NewChannelChange(m map[string]string) (cc ChannelChange) {
	cc.ChannelId = decoder.StringFromMap(m, "channel_id")
	cc.ConnectionId = decoder.StringFromMap(m, "connection_id")
	cc.CounterpartyChannelId = decoder.StringFromMap(m, "counterparty_channel_id")
	cc.CounterpartyPortId = decoder.StringFromMap(m, "counterparty_port_id")
	cc.PortId = decoder.StringFromMap(m, "port_id")
	return
}

type FungibleTokenPacket struct {
	Amount   storageTypes.Numeric
	Denom    string
	Memo     string
	Module   string
	Receiver string
	Sender   string
	Success  string
	Error    string
}

func NewFungibleTokenPacket(m map[string]string) (ftp FungibleTokenPacket) {
	ftp.Amount = decoder.NumericFromMap(m, "amount")
	ftp.Denom = decoder.StringFromMap(m, "denom")
	ftp.Memo = decoder.StringFromMap(m, "memo")
	ftp.Module = decoder.StringFromMap(m, "module")
	ftp.Receiver = decoder.StringFromMap(m, "receiver")
	ftp.Sender = decoder.StringFromMap(m, "sender")
	ftp.Success = decoder.StringFromMap(m, "success")
	ftp.Error = decoder.StringFromMap(m, "error")
	return
}

type AcknowledgementPacket struct {
	ConnectionID          string
	MsgIndex              string
	PacketChannelOrdering string
	PacketConnection      string
	PacketDstChannel      string
	PacketDstPort         string
	PacketSequence        uint64
	PacketSrcChannel      string
	PacketSrcPort         string
	Timeout               time.Time
	TimeoutHeight         uint64
}

func NewAcknowledgementPacket(m map[string]string) (ap AcknowledgementPacket, err error) {
	ap.ConnectionID = decoder.StringFromMap(m, "connection_id")
	ap.MsgIndex = decoder.StringFromMap(m, "msg_index")
	ap.PacketChannelOrdering = decoder.StringFromMap(m, "packet_channel_ordering")
	ap.PacketConnection = decoder.StringFromMap(m, "packet_connection")
	ap.PacketDstChannel = decoder.StringFromMap(m, "packet_dst_channel")
	ap.PacketDstPort = decoder.StringFromMap(m, "packet_dst_port")
	ap.PacketSequence, err = decoder.Uint64FromMap(m, "packet_sequence")
	if err != nil {
		return ap, errors.Wrap(err, "packet_sequence")
	}
	ap.PacketSrcChannel = decoder.StringFromMap(m, "packet_src_channel")
	ap.PacketSrcPort = decoder.StringFromMap(m, "packet_src_port")
	_, height, err := decoder.RevisionHeightFromMap(m, "packet_timeout_height")
	if err != nil {
		return ap, errors.Wrap(err, "packet_timeout_height")
	}
	ap.TimeoutHeight = height
	ap.Timeout = decoder.UnixNanoFromMap(m, "packet_timeout_timestamp")
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

func NewRecvPacket(m map[string]string) (rp RecvPacket, err error) {
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

func parseUnquoteOptional(s string) (string, error) {
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return strconv.Unquote(s)
	}
	return s, nil
}

type CreateMailbox struct {
	MailboxId    string
	Owner        string
	DefaultIsm   string
	DefaultHook  string
	RequiredHook string
	LocalDomain  uint64
}

func NewCreateMailbox(m map[string]string) (cm CreateMailbox, err error) {
	cm.MailboxId, err = parseUnquoteOptional(decoder.StringFromMap(m, "mailbox_id"))
	if err != nil {
		return cm, errors.Wrap(err, "mailbox_id")
	}
	cm.Owner, err = parseUnquoteOptional(decoder.StringFromMap(m, "owner"))
	if err != nil {
		return cm, errors.Wrap(err, "mailbox_id")
	}
	cm.DefaultIsm, err = parseUnquoteOptional(decoder.StringFromMap(m, "default_ism"))
	if err != nil {
		return cm, errors.Wrap(err, "default_ism")
	}
	cm.DefaultHook, err = parseUnquoteOptional(decoder.StringFromMap(m, "default_hook"))
	if err != nil {
		return cm, errors.Wrap(err, "default_hook")
	}
	cm.RequiredHook, err = parseUnquoteOptional(decoder.StringFromMap(m, "required_hook"))
	if err != nil {
		return cm, errors.Wrap(err, "required_hook")
	}
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

func NewSetMailbox(m map[string]string) (sm SetMailbox, err error) {
	sm.MailboxId, err = parseUnquoteOptional(decoder.StringFromMap(m, "mailbox_id"))
	if err != nil {
		return sm, errors.Wrap(err, "mailbox_id")
	}
	sm.Owner, err = parseUnquoteOptional(decoder.StringFromMap(m, "owner"))
	if err != nil {
		return sm, errors.Wrap(err, "mailbox_id")
	}
	sm.DefaultIsm, err = parseUnquoteOptional(decoder.StringFromMap(m, "default_ism"))
	if err != nil {
		return sm, errors.Wrap(err, "default_ism")
	}
	sm.DefaultHook, err = parseUnquoteOptional(decoder.StringFromMap(m, "default_hook"))
	if err != nil {
		return sm, errors.Wrap(err, "default_hook")
	}
	sm.NewOwner, err = parseUnquoteOptional(decoder.StringFromMap(m, "new_owner"))
	if err != nil {
		return sm, errors.Wrap(err, "new_owner")
	}
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

func NewHyperlaneProcessEvent(m map[string]string) (hpe HyperlaneProcessEvent, err error) {
	hpe.OriginMailboxId, err = parseUnquoteOptional(decoder.StringFromMap(m, "origin_mailbox_id"))
	if err != nil {
		return hpe, errors.Wrap(err, "origin_mailbox_id")
	}
	hpe.Sender, err = parseUnquoteOptional(decoder.StringFromMap(m, "sender"))
	if err != nil {
		return hpe, errors.Wrap(err, "sender")
	}
	hpe.Recipient, err = parseUnquoteOptional(decoder.StringFromMap(m, "recipient"))
	if err != nil {
		return hpe, errors.Wrap(err, "recipient")
	}
	hpe.MessageId, err = parseUnquoteOptional(decoder.StringFromMap(m, "message_id"))
	if err != nil {
		return hpe, errors.Wrap(err, "message_id")
	}
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

func NewHyperlaneDispatchEvent(m map[string]string) (hde HyperlaneDispatchEvent, err error) {
	hde.OriginMailboxId, err = parseUnquoteOptional(decoder.StringFromMap(m, "origin_mailbox_id"))
	if err != nil {
		return hde, errors.Wrap(err, "origin_mailbox_id")
	}
	hde.Sender, err = parseUnquoteOptional(decoder.StringFromMap(m, "sender"))
	if err != nil {
		return hde, errors.Wrap(err, "sender")
	}
	hde.Recipient, err = parseUnquoteOptional(decoder.StringFromMap(m, "recipient"))
	if err != nil {
		return hde, errors.Wrap(err, "recipient")
	}
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

func NewCreateCollateralToken(m map[string]string) (cct CreateCollateralToken, err error) {
	cct.MailboxId, err = parseUnquoteOptional(decoder.StringFromMap(m, "origin_mailbox"))
	if err != nil {
		return cct, errors.Wrap(err, "origin_mailbox_id")
	}
	cct.Owner, err = parseUnquoteOptional(decoder.StringFromMap(m, "owner"))
	if err != nil {
		return cct, errors.Wrap(err, "owner")
	}
	cct.TokenId, err = parseUnquoteOptional(decoder.StringFromMap(m, "token_id"))
	if err != nil {
		return cct, errors.Wrap(err, "token_id")
	}
	cct.Denom, err = parseUnquoteOptional(decoder.StringFromMap(m, "origin_denom"))
	if err != nil {
		return cct, errors.Wrap(err, "origin_denom")
	}
	return
}

type CreateSyntheticToken struct {
	MailboxId string
	Owner     string
	TokenId   string
	Denom     string
}

func NewCreateSyntheticToken(m map[string]string) (cst CreateSyntheticToken, err error) {
	cst.MailboxId, err = parseUnquoteOptional(decoder.StringFromMap(m, "origin_mailbox"))
	if err != nil {
		return cst, errors.Wrap(err, "origin_mailbox_id")
	}
	cst.Owner, err = parseUnquoteOptional(decoder.StringFromMap(m, "owner"))
	if err != nil {
		return cst, errors.Wrap(err, "owner")
	}
	cst.TokenId, err = parseUnquoteOptional(decoder.StringFromMap(m, "token_id"))
	if err != nil {
		return cst, errors.Wrap(err, "token_id")
	}
	cst.Denom, err = parseUnquoteOptional(decoder.StringFromMap(m, "origin_denom"))
	if err != nil {
		return cst, errors.Wrap(err, "origin_denom")
	}
	return
}

type HyperlaneReceiveTransferEvent struct {
	Amount       storageTypes.Numeric
	Denom        string
	OriginDomain uint64
	Recipient    string
	Sender       string
	TokenId      string
}

func NewHyperlaneReceiveTransferEvent(m map[string]string) (hrte HyperlaneReceiveTransferEvent, err error) {
	hrte.Sender, err = parseUnquoteOptional(decoder.StringFromMap(m, "sender"))
	if err != nil {
		return hrte, errors.Wrap(err, "sender")
	}
	hrte.Recipient, err = parseUnquoteOptional(decoder.StringFromMap(m, "recipient"))
	if err != nil {
		return hrte, errors.Wrap(err, "recipient")
	}
	hrte.TokenId, err = parseUnquoteOptional(decoder.StringFromMap(m, "token_id"))
	if err != nil {
		return hrte, errors.Wrap(err, "token_id")
	}
	hrte.OriginDomain, err = decoder.Uint64FromMap(m, "origin_domain")
	if err != nil {
		return hrte, errors.Wrap(err, "origin")
	}
	amount, err := parseUnquoteOptional(decoder.StringFromMap(m, "amount"))
	if err != nil {
		return hrte, errors.Wrap(err, "amount")
	}
	if amount != "" {
		coin, err := types.ParseCoinNormalized(amount)
		if err != nil {
			return hrte, errors.Wrap(err, amount)
		}
		hrte.Amount = storageTypes.NumericFromBigInt(coin.Amount.BigInt(), 0)
		hrte.Denom = coin.GetDenom()
	} else {
		hrte.Amount = storageTypes.NumericZero()
		hrte.Denom = currency.Utia
	}

	return hrte, nil
}

type HyperlaneSendTransferEvent struct {
	Amount            storageTypes.Numeric
	Denom             string
	DestinationDomain uint64
	Recipient         string
	Sender            string
	TokenId           string
}

func NewHyperlaneSendTransferEvent(m map[string]string) (hste HyperlaneSendTransferEvent, err error) {
	hste.Sender, err = parseUnquoteOptional(decoder.StringFromMap(m, "sender"))
	if err != nil {
		return hste, errors.Wrap(err, "sender")
	}
	hste.Recipient, err = parseUnquoteOptional(decoder.StringFromMap(m, "recipient"))
	if err != nil {
		return hste, errors.Wrap(err, "recipient")
	}
	hste.TokenId, err = parseUnquoteOptional(decoder.StringFromMap(m, "token_id"))
	if err != nil {
		return hste, errors.Wrap(err, "token_id")
	}
	hste.DestinationDomain, err = decoder.Uint64FromMap(m, "destination_domain")
	if err != nil {
		return hste, errors.Wrap(err, "origin")
	}
	amount, err := parseUnquoteOptional(decoder.StringFromMap(m, "amount"))
	if err != nil {
		return hste, errors.Wrap(err, "amount")
	}
	if amount != "" {
		coin, err := types.ParseCoinNormalized(amount)
		if err != nil {
			return hste, errors.Wrap(err, amount)
		}
		hste.Amount = storageTypes.NumericFromBigInt(coin.Amount.BigInt(), 0)
		hste.Denom = coin.GetDenom()
	} else {
		hste.Amount = storageTypes.NumericZero()
		hste.Denom = currency.Utia
	}
	return
}

type SetToken struct {
	IsmId             string
	TokenId           string
	NewOwner          string
	Owner             string
	RenounceOwnership bool
}

func NewSetToken(m map[string]string) (st SetToken, err error) {
	st.NewOwner, err = parseUnquoteOptional(decoder.StringFromMap(m, "new_owner"))
	if err != nil {
		return st, errors.Wrap(err, "new_owner")
	}
	st.Owner, err = parseUnquoteOptional(decoder.StringFromMap(m, "owner"))
	if err != nil {
		return st, errors.Wrap(err, "owner")
	}
	st.TokenId, err = parseUnquoteOptional(decoder.StringFromMap(m, "token_id"))
	if err != nil {
		return st, errors.Wrap(err, "token_id")
	}
	st.IsmId, err = parseUnquoteOptional(decoder.StringFromMap(m, "ism_id"))
	if err != nil {
		return st, errors.Wrap(err, "ism_id")
	}
	st.RenounceOwnership, err = decoder.BoolFromMap(m, "renounce_ownership")
	if err != nil {
		return st, errors.Wrap(err, "renounce_ownership")
	}
	return
}

type HyperlaneCreateIgpEvent struct {
	IgpId string
	Owner string
	Denom string
}

func NewHyperlaneCreateIgpEvent(m map[string]string) (hcie HyperlaneCreateIgpEvent, err error) {
	hcie.IgpId, err = parseUnquoteOptional(decoder.StringFromMap(m, "igp_id"))
	if err != nil {
		return hcie, errors.Wrap(err, "igp_id")
	}
	hcie.Owner, err = parseUnquoteOptional(decoder.StringFromMap(m, "owner"))
	if err != nil {
		return hcie, errors.Wrap(err, "owner")
	}
	hcie.Denom, err = parseUnquoteOptional(decoder.StringFromMap(m, "denom"))
	if err != nil {
		return hcie, errors.Wrap(err, "denom")
	}
	return
}

type HyperlaneSetDestinationGasConfig struct {
	GasOverhead       storageTypes.Numeric
	GasPrice          storageTypes.Numeric
	IgpId             string
	Owner             string
	RemoteDomain      uint64
	TokenExchangeRate string
}

func NewHyperlaneSetDestinationGasConfig(m map[string]string) (hsdgc HyperlaneSetDestinationGasConfig, err error) {
	price, err := parseUnquoteOptional(decoder.StringFromMap(m, "gas_price"))
	if err != nil {
		return hsdgc, errors.Wrap(err, "gas_price")
	}
	hsdgc.GasPrice, err = storageTypes.NumericFromString(price)
	if err != nil {
		return hsdgc, errors.Wrap(err, "parse gas_price")
	}
	overhead, err := parseUnquoteOptional(decoder.StringFromMap(m, "gas_overhead"))
	if err != nil {
		return hsdgc, errors.Wrap(err, "gas_overhead")
	}
	hsdgc.GasOverhead, err = storageTypes.NumericFromString(overhead)
	if err != nil {
		return hsdgc, errors.Wrap(err, "parse gas_overhead")
	}
	hsdgc.IgpId, err = parseUnquoteOptional(decoder.StringFromMap(m, "igp_id"))
	if err != nil {
		return hsdgc, errors.Wrap(err, "igp_id")
	}
	hsdgc.Owner, err = parseUnquoteOptional(decoder.StringFromMap(m, "owner"))
	if err != nil {
		return hsdgc, errors.Wrap(err, "owner")
	}
	hsdgc.RemoteDomain, err = decoder.Uint64FromMap(m, "remote_domain")
	if err != nil {
		return hsdgc, errors.Wrap(err, "remote_domain")
	}
	hsdgc.TokenExchangeRate, err = parseUnquoteOptional(decoder.StringFromMap(m, "token_exchange_rate"))
	if err != nil {
		return hsdgc, errors.Wrap(err, "token_exchange_rate")
	}

	return
}

type HyperlaneGasPaymentEvent struct {
	Amount    storageTypes.Numeric
	GasAmount storageTypes.Numeric
	IgpId     string
}

func NewHyperlaneGasPaymentEvent(m map[string]string) (hgpe HyperlaneGasPaymentEvent, err error) {
	gasAmount, err := parseUnquoteOptional(decoder.StringFromMap(m, "gas_amount"))
	if err != nil {
		return hgpe, errors.Wrap(err, "gas_amount")
	}
	hgpe.GasAmount, err = storageTypes.NumericFromString(gasAmount)
	if err != nil {
		return hgpe, errors.Wrap(err, "parse gas_amount")
	}

	payment, err := parseUnquoteOptional(decoder.StringFromMap(m, "payment"))
	if err != nil {
		return hgpe, errors.Wrap(err, "payment")
	}
	coin, err := types.ParseCoinNormalized(payment)
	if err != nil {
		return hgpe, errors.Wrap(err, payment)
	}
	hgpe.Amount = storageTypes.NumericFromBigInt(coin.Amount.BigInt(), 0)

	hgpe.IgpId, err = parseUnquoteOptional(decoder.StringFromMap(m, "igp_id"))
	if err != nil {
		return hgpe, errors.Wrap(err, "igp_id")
	}

	return
}

type HyperlaneSetIgpEvent struct {
	IgpId             string
	Owner             string
	NewOwner          string
	RenounceOwnership bool
}

func NewHyperlaneSetIgpEvent(m map[string]string) (hsie HyperlaneSetIgpEvent, err error) {
	hsie.IgpId, err = parseUnquoteOptional(decoder.StringFromMap(m, "igp_id"))
	if err != nil {
		return hsie, errors.Wrap(err, "igp_id")
	}
	hsie.Owner, err = parseUnquoteOptional(decoder.StringFromMap(m, "owner"))
	if err != nil {
		return hsie, errors.Wrap(err, "owner")
	}
	hsie.NewOwner, err = parseUnquoteOptional(decoder.StringFromMap(m, "new_owner"))
	if err != nil {
		return hsie, errors.Wrap(err, "new_owner")
	}
	hsie.RenounceOwnership, err = decoder.BoolFromMap(m, "renounce_ownership")
	if err != nil {
		return hsie, errors.Wrap(err, "renounce_ownership")
	}
	return
}

type EventForwardingComplete struct {
	ForwardAddress       string
	DestinationDomain    uint64
	DestinationRecipient []byte
	SuccessfulCount      uint64
	FailedCount          uint64
}

func NewEventForwardingComplete(m map[string]string) (efc EventForwardingComplete, err error) {
	efc.ForwardAddress, err = parseUnquoteOptional(decoder.StringFromMap(m, "forward_addr"))
	if err != nil {
		return efc, errors.Wrap(err, "forward_addr")
	}
	efc.DestinationDomain, err = decoder.Uint64FromMap(m, "dest_domain")
	if err != nil {
		return efc, errors.Wrap(err, "dest_domain")
	}
	efc.DestinationRecipient, err = decoder.BytesFromMap(m, "dest_recipient")
	if err != nil {
		return efc, errors.Wrap(err, "dest_recipient")
	}
	efc.SuccessfulCount, err = decoder.Uint64FromMap(m, "tokens_forwarded")
	if err != nil {
		return efc, errors.Wrap(err, "tokens_forwarded")
	}
	efc.FailedCount, err = decoder.Uint64FromMap(m, "tokens_failed")
	if err != nil {
		return efc, errors.Wrap(err, "tokens_failed")
	}
	return
}

type EventTokenForwarded struct {
	ForwardAddress string
	Denom          string
	Amount         string
	MessageId      string
	Success        bool
	Error          string
}

func NewEventTokenForwarded(m map[string]string) (etf EventTokenForwarded, err error) {
	etf.ForwardAddress, err = parseUnquoteOptional(decoder.StringFromMap(m, "forward_addr"))
	if err != nil {
		return etf, errors.Wrap(err, "forward_addr")
	}
	etf.Denom, err = parseUnquoteOptional(decoder.StringFromMap(m, "denom"))
	if err != nil {
		return etf, errors.Wrap(err, "denom")
	}
	etf.Amount, err = parseUnquoteOptional(decoder.StringFromMap(m, "amount"))
	if err != nil {
		return etf, errors.Wrap(err, "amount")
	}
	etf.MessageId, err = parseUnquoteOptional(decoder.StringFromMap(m, "message_id"))
	if err != nil {
		return etf, errors.Wrap(err, "message_id")
	}
	etf.Success, err = decoder.BoolFromMap(m, "success")
	if err != nil {
		return etf, errors.Wrap(err, "success")
	}
	etf.Error, err = parseUnquoteOptional(decoder.StringFromMap(m, "error"))
	if err != nil {
		return etf, errors.Wrap(err, "error")
	}
	return
}

type ZkISMCreateEvent struct {
	Id                  []byte
	Creator             string
	State               []byte
	MerkleTreeAddress   []byte
	Groth16VKey         []byte
	StateTransitionVKey []byte
	StateMembershipVKey []byte
}

func NewZkISMCreateEvent(m map[string]string) (e ZkISMCreateEvent, err error) {
	e.Id, err = decoder.BytesFromMap(m, "id")
	if err != nil {
		return e, errors.Wrap(err, "id")
	}
	e.Creator, err = parseUnquoteOptional(decoder.StringFromMap(m, "owner"))
	if err != nil {
		return e, errors.Wrap(err, "creator")
	}
	e.State, err = decoder.BytesFromMap(m, "state")
	if err != nil {
		return e, errors.Wrap(err, "state")
	}
	e.MerkleTreeAddress, err = decoder.BytesFromMap(m, "merkle_tree_address")
	if err != nil {
		return e, errors.Wrap(err, "merkle_tree_address")
	}
	e.Groth16VKey, err = decoder.BytesFromMap(m, "groth16_vkey")
	if err != nil {
		return e, errors.Wrap(err, "groth16_vkey")
	}
	e.StateTransitionVKey, err = decoder.BytesFromMap(m, "state_transition_vkey")
	if err != nil {
		return e, errors.Wrap(err, "state_transition_vkey")
	}
	e.StateMembershipVKey, err = decoder.BytesFromMap(m, "state_membership_vkey")
	if err != nil {
		return e, errors.Wrap(err, "state_membership_vkey")
	}
	return
}

type ZkISMUpdateEvent struct {
	Id       []byte
	NewState []byte
}

func NewZkISMUpdateEvent(m map[string]string) (e ZkISMUpdateEvent, err error) {
	e.Id, err = decoder.BytesFromMap(m, "id")
	if err != nil {
		return e, errors.Wrap(err, "id")
	}
	e.NewState, err = decoder.BytesFromMap(m, "state")
	if err != nil {
		return e, errors.Wrap(err, "state")
	}
	return
}

type ZkISMSubmitMessagesEvent struct {
	Id         []byte
	StateRoot  []byte
	MessageIds [][]byte
}

func NewZkISMSubmitMessagesEvent(m map[string]string) (e ZkISMSubmitMessagesEvent, err error) {
	e.Id, err = decoder.BytesFromMap(m, "id")
	if err != nil {
		return e, errors.Wrap(err, "id")
	}
	e.StateRoot, err = decoder.BytesFromMap(m, "state_root")
	if err != nil {
		return e, errors.Wrap(err, "state_root")
	}
	messageIds := decoder.StringFromMap(m, "messages")
	messageIds = strings.Trim(messageIds, "[]")
	if messageIds != "" {
		for _, id := range strings.Split(messageIds, ",") {
			str, err := parseUnquoteOptional(id)
			if err != nil {
				return e, errors.Wrap(err, "failed to unquote string")
			}
			idBytes, err := hex.DecodeString(strings.TrimPrefix(str, "0x"))
			if err != nil {
				return e, errors.Wrap(err, "failed to decode hex string")
			}
			e.MessageIds = append(e.MessageIds, idBytes)
		}
	}
	return
}
