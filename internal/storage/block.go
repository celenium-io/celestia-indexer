package storage

import (
	"context"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
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
	bun.BaseModel `bun:"table:block" comment:"Table with celestia blocks."`

	Id           uint64         `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height       pkgTypes.Level `bun:"height"                    comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time         time.Time      `bun:"time,pk,notnull"           comment:"The time of block"                 stats:"func:min max,filterable"`
	VersionBlock uint64         `bun:"version_block"             comment:"Block version"`
	VersionApp   uint64         `bun:"version_app"               comment:"App version"`

	TxCount      uint64            `bun:"tx_count"                comment:"Count of transactions in block"            stats:"func:min max sum avg"`
	EventsCount  uint64            `bun:"events_count"            comment:"Count of events in begin and end of block" stats:"func:min max sum avg"`
	BlobsSize    uint64            `bun:"blobs_size"              comment:"Summary blocks size from pay for blob"     stats:"func:min max sum avg"`
	MessageTypes types.MsgTypeBits `bun:"message_types,type:int8" comment:"Bit mask with containing messages"`

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

	Fee     decimal.Decimal `bun:"fee,type:numeric" comment:"Summary block fee" stats:"func:min max sum avg"`
	ChainId string          `bun:"-"` // internal field for filling state

	Txs    []Tx    `bun:"rel:has-many"`
	Events []Event `bun:"rel:has-many"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
