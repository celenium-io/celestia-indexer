package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// IBlock -
type IBlock interface {
	storage.Table[*Block]

	Last(ctx context.Context) (Block, error)
}

// Block -
type Block struct {
	bun.BaseModel `bun:"table:block" comment:"Table with celestia blocks."`

	Id           uint64    `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height       uint64    `bun:"height"                    comment:"The number (height) of this block"`
	Time         time.Time `bun:"time"                      comment:"The time of block"`
	VersionBlock string    `bun:"version_block"             comment:"Block version"`
	VersionApp   string    `bun:"version_app"               comment:"App version"`

	TxCount uint64 `bun:"tx_count" comment:"Count of transactions in block"`

	Hash               []byte `bun:"hash"                 comment:"Block hash"`
	ParentHash         []byte `bun:"parent_hash"          comment:"Hash of parent block"`
	LastCommitHash     []byte `bun:"last_commit_hash"     comment:"Last commit hash"`
	DataHash           []byte `bun:"data_hash"            comment:"Data hash"`
	ValidatorsHash     []byte `bun:"validators_hash"      comment:"Validators hash"`
	NextValidatorsHash []byte `bun:"next_validators_hash" comment:"Next validators hash"`
	ConsensusHash      []byte `bun:"consensus_hash"       comment:"Consensus hash"`
	AppHash            []byte `bun:"app_hash"             comment:"App hash"`
	LastResultsHash    []byte `bun:"last_results_hash"    comment:"Last results hash"`
	EvidenceHash       []byte `bun:"evidence_hash"        comment:"Evidence hash"`
	ProposerAddress    []byte `bun:"proposer_address"     comment:"Proposer address"`

	Txs []Tx `bun:"rel:has-many"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
