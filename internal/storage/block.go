// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlock interface {
	storage.Table[*Block]

	Last(ctx context.Context) (Block, error)
	ByIdWithRelations(ctx context.Context, id uint64) (Block, error)
	ByHeight(ctx context.Context, height pkgTypes.Level) (Block, error)
	ByHeightWithStats(ctx context.Context, height pkgTypes.Level) (Block, error)
	ByHash(ctx context.Context, hash []byte) (Block, error)
	ByProposer(ctx context.Context, proposerId uint64, limit, offset int) ([]Block, error)
	ListWithStats(ctx context.Context, limit, offset uint64, order storage.SortOrder) ([]*Block, error)
	Time(ctx context.Context, height pkgTypes.Level) (time.Time, error)
}

// Block -
type Block struct {
	bun.BaseModel `bun:"table:block" comment:"Table with celestia blocks."`

	Id           uint64         `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height       pkgTypes.Level `bun:"height"                    comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time         time.Time      `bun:"time,pk,notnull"           comment:"The time of block"                 stats:"func:min max,filterable"`
	VersionBlock uint64         `bun:"version_block"             comment:"Block version"`
	VersionApp   uint64         `bun:"version_app"               comment:"App version"`

	MessageTypes types.MsgTypeBits `bun:"message_types,type:bit(76)" comment:"Bit mask with containing messages"`

	Hash               pkgTypes.Hex `bun:"hash"                 comment:"Block hash"`
	ParentHash         pkgTypes.Hex `bun:"parent_hash"          comment:"Hash of parent block"`
	LastCommitHash     pkgTypes.Hex `bun:"last_commit_hash"     comment:"Last commit hash"`
	DataHash           pkgTypes.Hex `bun:"data_hash"            comment:"Data hash"`
	ValidatorsHash     pkgTypes.Hex `bun:"validators_hash"      comment:"Validators hash"`
	NextValidatorsHash pkgTypes.Hex `bun:"next_validators_hash" comment:"Next validators hash"`
	ConsensusHash      pkgTypes.Hex `bun:"consensus_hash"       comment:"Consensus hash"`
	AppHash            pkgTypes.Hex `bun:"app_hash"             comment:"App hash"`
	LastResultsHash    pkgTypes.Hex `bun:"last_results_hash"    comment:"Last results hash"`
	EvidenceHash       pkgTypes.Hex `bun:"evidence_hash"        comment:"Evidence hash"`
	ProposerId         uint64       `bun:"proposer_id,nullzero" comment:"Proposer internal id"`

	ChainId         string           `bun:"-" json:"-"` // internal field for filling state
	ProposerAddress string           `bun:"-" json:"-"` // internal field for proposer
	BlockSignatures []BlockSignature `bun:"-" json:"-"` // internal field for block signature

	Txs      []Tx       `bun:"rel:has-many"                   json:"-"`
	Events   []Event    `bun:"rel:has-many"                   json:"-"`
	Stats    BlockStats `bun:"rel:has-one,join:height=height"`
	Proposer Validator  `bun:"rel:belongs-to"                 json:"-"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
