// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"io"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
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
	&Validator{},
	&Delegation{},
	&Redelegation{},
	&Undelegation{},
	&StakingLog{},
	&Jail{},
	&BlobLog{},
	&Price{},
	&Rollup{},
	&RollupProvider{},
	&Grant{},
	&ApiKey{},
	&Tvl{},
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
	RetentionBlockSignatures(ctx context.Context, height types.Level) error
	CancelUnbondings(ctx context.Context, cancellations ...Undelegation) error
	RetentionCompletedUnbondings(ctx context.Context, blockTime time.Time) error
	RetentionCompletedRedelegations(ctx context.Context, blockTime time.Time) error
	Jail(ctx context.Context, validators ...*Validator) error

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
	RollbackUndelegations(ctx context.Context, height types.Level) (err error)
	RollbackRedelegations(ctx context.Context, height types.Level) (err error)
	RollbackStakingLogs(ctx context.Context, height types.Level) ([]StakingLog, error)
	RollbackJails(ctx context.Context, height types.Level) ([]Jail, error)
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
	Delegation(ctx context.Context, validatorId, addressId uint64) (val Delegation, err error)
	RefreshLeaderboard(ctx context.Context) error
}

const (
	ChannelHead  = "head"
	ChannelBlock = "block"
)

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
