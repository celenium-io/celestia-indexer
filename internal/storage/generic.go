// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"io"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var Models = []any{
	&State{},
	&Constant{},
	&DenomMetadata{},
	&Balance{},
	&Address{},
	&VestingAccount{},
	&VestingPeriod{},
	&Block{},
	&BlockStats{},
	&BlockSignature{},
	&Tx{},
	&Message{},
	&Event{},
	&Namespace{},
	&NamespaceMessage{},
	&Signer{},
	&MsgAddress{},
	&MsgValidator{},
	&Validator{},
	&Delegation{},
	&Redelegation{},
	&Undelegation{},
	&StakingLog{},
	&Jail{},
	&BlobLog{},
	&Rollup{},
	&RollupProvider{},
	&Grant{},
	&ApiKey{},
	&celestials.Celestial{},
	&celestials.CelestialState{},
	&Proposal{},
	&Vote{},
	&IbcClient{},
	&IbcConnection{},
	&IbcChannel{},
	&IbcTransfer{},
	&HLMailbox{},
	&HLToken{},
	&HLTransfer{},
	&SignalVersion{},
	&Upgrade{},
	&HLIGP{},
	&HLIGPConfig{},
	&HLGasPayment{},
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Notificator interface {
	Notify(ctx context.Context, channel string, payload string) error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Listener interface {
	io.Closer

	Subscribe(ctx context.Context, channels ...string) error
	Listen() chan *pq.Notification
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ListenerFactory interface {
	CreateListener() Listener
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Transaction interface {
	sdk.Transaction

	SaveConstants(ctx context.Context, constants ...Constant) error
	SaveTransactions(ctx context.Context, txs ...Tx) error
	SaveNamespaces(ctx context.Context, namespaces ...*Namespace) (int64, error)
	SaveAddresses(ctx context.Context, addresses ...*Address) (int64, error)
	SaveVestingAccounts(ctx context.Context, accounts ...*VestingAccount) error
	SaveVestingPeriods(ctx context.Context, periods ...VestingPeriod) error
	SaveBalances(ctx context.Context, balances ...Balance) error
	SaveMessages(ctx context.Context, msgs ...*Message) error
	SaveSigners(ctx context.Context, addresses ...Signer) error
	SaveMsgAddresses(ctx context.Context, addresses ...MsgAddress) error
	SaveMsgValidator(ctx context.Context, validatorMsgs ...MsgValidator) error
	SaveNamespaceMessage(ctx context.Context, nsMsgs ...NamespaceMessage) error
	SaveBlobLogs(ctx context.Context, logs ...BlobLog) error
	SaveValidators(ctx context.Context, validators ...*Validator) (int, error)
	SaveEvents(ctx context.Context, events ...Event) error
	SaveRollup(ctx context.Context, rollup *Rollup) error
	SaveGrants(ctx context.Context, grants ...Grant) error
	UpdateRollup(ctx context.Context, rollup *Rollup) error
	SaveProviders(ctx context.Context, providers ...RollupProvider) error
	SaveUndelegations(ctx context.Context, undelegations ...Undelegation) error
	SaveRedelegations(ctx context.Context, redelegations ...Redelegation) error
	SaveDelegations(ctx context.Context, delegations ...Delegation) error
	UpdateSlashedDelegations(ctx context.Context, validatorId uint64, fraction decimal.Decimal) ([]Balance, error)
	SaveStakingLogs(ctx context.Context, logs ...StakingLog) error
	SaveJails(ctx context.Context, jails ...Jail) error
	SaveBlockSignatures(ctx context.Context, signs ...BlockSignature) error
	SaveProposals(ctx context.Context, proposals ...*Proposal) (int64, error)
	SaveVotes(ctx context.Context, votes ...*Vote) error
	SaveIbcClients(ctx context.Context, clients ...*IbcClient) (int64, error)
	SaveIbcConnections(ctx context.Context, connections ...*IbcConnection) error
	SaveIbcChannels(ctx context.Context, channels ...*IbcChannel) error
	SaveIbcTransfers(ctx context.Context, transfers ...*IbcTransfer) error
	SaveHyperlaneMailbox(ctx context.Context, mailbox ...*HLMailbox) error
	SaveHyperlaneTokens(ctx context.Context, tokens ...*HLToken) error
	SaveHyperlaneTransfers(ctx context.Context, transfers ...*HLTransfer) error
	RetentionBlockSignatures(ctx context.Context, height types.Level) error
	CancelUnbondings(ctx context.Context, cancellations ...Undelegation) error
	RetentionCompletedUnbondings(ctx context.Context, blockTime time.Time) error
	RetentionCompletedRedelegations(ctx context.Context, blockTime time.Time) error
	Jail(ctx context.Context, validators ...*Validator) error
	SaveSignals(ctx context.Context, signals ...*SignalVersion) error
	SaveUpgrades(ctx context.Context, signals ...*Upgrade) error
	UpdateSignalsAfterUpgrade(ctx context.Context, version uint64) error
	SaveIgps(ctx context.Context, igps ...HLIGP) error
	SaveIgpConfigs(ctx context.Context, configs ...HLIGPConfig) error
	SaveHyperlaneGasPayments(ctx context.Context, payments ...*HLGasPayment) error

	RollbackBlock(ctx context.Context, height types.Level) error
	RollbackBlockStats(ctx context.Context, height types.Level) (stats BlockStats, err error)
	RollbackAddresses(ctx context.Context, height types.Level) (address []Address, err error)
	RollbackVestingAccounts(ctx context.Context, height types.Level) error
	RollbackVestingPeriods(ctx context.Context, height types.Level) error
	RollbackTxs(ctx context.Context, height types.Level) (txs []Tx, err error)
	RollbackEvents(ctx context.Context, height types.Level) (events []Event, err error)
	RollbackMessages(ctx context.Context, height types.Level) (msgs []Message, err error)
	RollbackNamespaceMessages(ctx context.Context, height types.Level) (msgs []NamespaceMessage, err error)
	RollbackNamespaces(ctx context.Context, height types.Level) (ns []Namespace, err error)
	RollbackValidators(ctx context.Context, height types.Level) ([]Validator, error)
	RollbackBlobLog(ctx context.Context, height types.Level) error
	RollbackGrants(ctx context.Context, height types.Level) error
	RollbackBlockSignatures(ctx context.Context, height types.Level) (err error)
	RollbackSigners(ctx context.Context, txIds []uint64) (err error)
	RollbackMessageAddresses(ctx context.Context, msgIds []uint64) (err error)
	RollbackMessageValidators(ctx context.Context, height types.Level) (err error)
	RollbackUndelegations(ctx context.Context, height types.Level) (err error)
	RollbackRedelegations(ctx context.Context, height types.Level) (err error)
	RollbackStakingLogs(ctx context.Context, height types.Level) ([]StakingLog, error)
	RollbackJails(ctx context.Context, height types.Level) ([]Jail, error)
	RollbackProposals(ctx context.Context, height types.Level) error
	RollbackVotes(ctx context.Context, height types.Level) error
	RollbackIbcClients(ctx context.Context, height types.Level) error
	RollbackIbcConnections(ctx context.Context, height types.Level) error
	RollbackIbcChannels(ctx context.Context, height types.Level) error
	RollbackIbcTransfers(ctx context.Context, height types.Level) error
	RollbackHyperlaneMailbox(ctx context.Context, height types.Level) error
	RollbackHyperlaneTokens(ctx context.Context, height types.Level) error
	RollbackHyperlaneTransfers(ctx context.Context, height types.Level) error
	RollbackSignals(ctx context.Context, height types.Level) error
	RollbackUpgrades(ctx context.Context, height types.Level) error
	DeleteBalances(ctx context.Context, ids []uint64) error
	DeleteProviders(ctx context.Context, rollupId uint64) error
	DeleteRollup(ctx context.Context, rollupId uint64) error
	DeleteDelegationsByValidator(ctx context.Context, ids ...uint64) error
	UpdateValidators(ctx context.Context, validators ...*Validator) error

	State(ctx context.Context, name string) (state State, err error)
	LastBlock(ctx context.Context) (block Block, err error)
	Namespace(ctx context.Context, id uint64) (ns Namespace, err error)
	LastNamespaceMessage(ctx context.Context, nsId uint64) (msg NamespaceMessage, err error)
	LastAddressAction(ctx context.Context, address []byte) (uint64, error)
	GetProposerId(ctx context.Context, address string) (uint64, error)
	Validator(ctx context.Context, id uint64) (val Validator, err error)
	BondedValidators(ctx context.Context, limit int) ([]Validator, error)
	Delegation(ctx context.Context, validatorId, addressId uint64) (val Delegation, err error)
	AddressDelegations(ctx context.Context, addressId uint64) (val []Delegation, err error)
	ActiveProposals(ctx context.Context) ([]Proposal, error)
	ProposalVotes(ctx context.Context, proposalId uint64, limit, offset int) ([]Vote, error)
	Proposal(ctx context.Context, id uint64) (Proposal, error)
	RefreshLeaderboard(ctx context.Context) error
	IbcConnection(ctx context.Context, id string) (IbcConnection, error)
	HyperlaneMailbox(ctx context.Context, internalId uint64) (HLMailbox, error)
	HyperlaneToken(ctx context.Context, id []byte) (HLToken, error)
	SignalVersions(ctx context.Context) ([]Signal, error)
	HyperlaneIgp(ctx context.Context, id []byte) (HLIGP, error)
	HyperlaneIgpConfig(ctx context.Context, id []byte) (HLIGPConfig, error)
}

const (
	ChannelHead  = "head"
	ChannelBlock = "block"
)

type Signal struct {
	VotingPower decimal.Decimal `bun:"voting_power"`
	Version     uint64          `bun:"version"`
}

type SearchResult struct {
	Id    uint64 `bun:"id"`
	Value string `bun:"value"`
	Type  string `bun:"type"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ISearch interface {
	Search(ctx context.Context, query []byte) ([]SearchResult, error)
	SearchText(ctx context.Context, text string) ([]SearchResult, error)
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Export interface {
	ToCsv(ctx context.Context, writer io.Writer, query string) error
	Close() error
}
