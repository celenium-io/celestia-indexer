// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

//go:generate mockgen -destination=mock/$GOFILE -package=mock -typed . Transaction

// Inserter inserts a pointer to a slice of sdk.Model; use storage.Insert.
type Inserter interface {
	Insert(ctx context.Context, models any) error
}

type Transaction interface {
	sdk.Transaction

	Inserter
	RollbackByHeight(ctx context.Context, height pkgTypes.Level, models ...sdk.Model) error

	BlockTx
	AccountTx
	NamespaceTx
	StakingTx
	GovTx
	IbcTx
	HyperlaneTx
	ZkIsmTx
	RollupTx
}

// BlockTx covers block, transactions, messages and events.
type BlockTx interface {
	Inserter

	SaveConstants(ctx context.Context, constants ...Constant) error
	SaveTransactions(ctx context.Context, txs ...Tx) error
	SaveMessages(ctx context.Context, msgs ...*Message) error
	SaveMsgAddresses(ctx context.Context, addresses ...*MsgAddress) error
	SaveEvents(ctx context.Context, events ...Event) error
	RetentionBlockSignatures(ctx context.Context, height pkgTypes.Level) error
	RollbackBlockStats(ctx context.Context, height pkgTypes.Level) (stats BlockStats, err error)
	RollbackTxs(ctx context.Context, height pkgTypes.Level) (txs []Tx, err error)
	RollbackEvents(ctx context.Context, height pkgTypes.Level) (events []Event, err error)
	RollbackMessages(ctx context.Context, height pkgTypes.Level) (msgs []Message, err error)
	RollbackSigners(ctx context.Context, txIds []uint64) (err error)
	RollbackMessageAddresses(ctx context.Context, msgIds []uint64) (err error)
}

// AccountTx covers addresses, balances, vesting and grants.
type AccountTx interface {
	Inserter

	SaveAddresses(ctx context.Context, addresses ...*Address) (int64, error)
	SaveBalances(ctx context.Context, balances ...Balance) error
	SaveVestingAccounts(ctx context.Context, accounts ...*VestingAccount) error
	SaveGrants(ctx context.Context, grants ...*Grant) error
	RollbackAddresses(ctx context.Context, height pkgTypes.Level) (address []Address, err error)
	RollbackGrants(ctx context.Context, height pkgTypes.Level) error
	DeleteBalances(ctx context.Context, ids []uint64) error
}

// NamespaceTx covers namespaces and blobs.
type NamespaceTx interface {
	Inserter

	SaveNamespaces(ctx context.Context, namespaces ...*Namespace) (int64, error)
	SaveBlobLogs(ctx context.Context, logs ...*BlobLog) error
	RollbackNamespaceMessages(ctx context.Context, height pkgTypes.Level) (msgs []NamespaceMessage, err error)
	RollbackNamespaces(ctx context.Context, height pkgTypes.Level) (ns []Namespace, err error)
}

// StakingTx covers validators, delegations and jails.
type StakingTx interface {
	Inserter

	SaveValidators(ctx context.Context, validators ...*Validator) (int, error)
	UpdateValidators(ctx context.Context, validators ...*Validator) error
	SaveDelegations(ctx context.Context, delegations ...Delegation) error
	UpdateSlashedDelegations(ctx context.Context, validatorId uint64, burned types.Numeric) ([]Balance, error)
	DeleteDelegationsByValidator(ctx context.Context, ids ...uint64) error
	SaveStakingLogs(ctx context.Context, logs ...StakingLog) error
	CancelUnbondings(ctx context.Context, cancellations ...*Undelegation) error
	RetentionCompletedUnbondings(ctx context.Context, blockTime time.Time) error
	RetentionCompletedRedelegations(ctx context.Context, blockTime time.Time) error
	Jail(ctx context.Context, validators ...*Validator) error
	RollbackValidators(ctx context.Context, height pkgTypes.Level) ([]Validator, error)
	RollbackStakingLogs(ctx context.Context, height pkgTypes.Level) ([]StakingLog, error)
	RollbackJails(ctx context.Context, height pkgTypes.Level) ([]Jail, error)
	RollbackBondUpdates(ctx context.Context, height pkgTypes.Level) ([]ValidatorBondUpdate, error)
}

// GovTx covers proposals, votes, signals and upgrades.
type GovTx interface {
	Inserter

	SaveProposals(ctx context.Context, proposals ...*Proposal) (int64, error)
	SaveVotes(ctx context.Context, votes ...*Vote) (map[uint64]*VotesCount, error)
	SaveUpgrades(ctx context.Context, upgrades ...*Upgrade) error
	UpdateSignalsAfterUpgrade(ctx context.Context, version uint64) (types.Numeric, error)
}

// IbcTx covers IBC clients, connections and channels.
type IbcTx interface {
	Inserter

	SaveIbcClients(ctx context.Context, clients ...*IbcClient) (int64, error)
	RecoverIbcClient(ctx context.Context, subjectId, substituteId string, updatedAt time.Time) error
	SaveIbcConnections(ctx context.Context, connections ...*IbcConnection) error
	SaveIbcChannels(ctx context.Context, channels ...*IbcChannel) error
	RollbackIbcClients(ctx context.Context, height pkgTypes.Level) (int64, error)
	RollbackIbcConnections(ctx context.Context, height pkgTypes.Level) (int64, error)
	RollbackIbcChannels(ctx context.Context, height pkgTypes.Level) (int64, error)
}

// HyperlaneTx covers Hyperlane mailboxes, tokens, transfers and IGPs.
type HyperlaneTx interface {
	Inserter

	SaveHyperlaneMailbox(ctx context.Context, mailbox ...*HLMailbox) error
	SaveHyperlaneTokens(ctx context.Context, tokens ...*HLToken) error
	SaveHyperlaneTransfers(ctx context.Context, transfers ...*HLTransfer) error
	SaveHyperlaneIgps(ctx context.Context, igps ...*HLIGP) error
	SaveHyperlaneIgpConfigs(ctx context.Context, configs ...HLIGPConfig) error
}

// ZkIsmTx covers ZK ISMs.
type ZkIsmTx interface {
	Inserter

	SaveZkISMs(ctx context.Context, items ...*ZkISM) error
}

// RollupTx covers rollups and their providers.
type RollupTx interface {
	Inserter

	UpdateRollup(ctx context.Context, rollup *Rollup) error
	DeleteProviders(ctx context.Context, rollupId uint64) error
	DeleteRollup(ctx context.Context, rollupId uint64) error
}

// Insert is a type-safe wrapper around Inserter.Insert.
func Insert[T sdk.Model](ctx context.Context, tx Inserter, items ...T) error {
	if len(items) == 0 {
		return nil
	}
	return tx.Insert(ctx, &items)
}

type TxRepos struct {
	Address          IAddress
	Blocks           IBlock
	BondUpdates      IValidatorBondUpdate
	Constants        IConstant
	Delegation       IDelegation
	HyperlaneMailbox IHLMailbox
	HyperlaneToken   IHLToken
	HyperlaneIgp     IHLIGP
	IbcConnections   IIbcConnection
	Namespace        INamespace
	Proposals        IProposal
	State            IState
	Validators       IValidator
	Votes            IVote
	ZkIsm            IZkISM
}

type TxReposFactory func(tx Transaction) TxRepos
