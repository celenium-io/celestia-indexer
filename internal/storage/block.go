package storage

import (
	"context"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlock interface {
	storage.Table[*Block]

	Last(ctx context.Context) (Block, error)
	ByHeight(ctx context.Context, height uint64) (Block, error)
	ByHash(ctx context.Context, hash []byte) (Block, error)
}

// Block -
type Block struct {
	bun.BaseModel `bun:"table:block" comment:"Table with celestia blocks." json:"-"`

	Id           uint64    `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height       uint64    `bun:"height"                    comment:"The number (height) of this block"`
	Time         time.Time `bun:"time"                      comment:"The time of block"`
	VersionBlock string    `bun:"version_block"             comment:"Block version"`
	VersionApp   string    `bun:"version_app"               comment:"App version"`

	TxCount       uint64            `bun:"tx_count"                comment:"Count of transactions in block"`
	EventsCount   uint64            `bun:"events_count"            comment:"Count of events in begin and end of block"`
	MessageTypes  types.MsgTypeBits `bun:"message_types,type:int8" comment:"Bit mask with containing messages"`
	NamespaceSize uint64            `bun:"namespace_size"          comment:"Summary block namespace size from pay for blob"`

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

	Fee     decimal.Decimal `bun:"fee,type:numeric" comment:"Summary block fee"`
	ChainId string          `bun:"-"` // internal field for filling state

	Txs    []Tx    `bun:"rel:has-many" json:"-"`
	Events []Event `bun:"rel:has-many" json:"-"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
