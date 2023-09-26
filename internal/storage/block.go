package storage

import (
	"context"
	"time"

	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"

	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlock interface {
	storage.Table[*Block]

	Last(ctx context.Context) (Block, error)
	ByHeight(ctx context.Context, height uint64) (Block, error)
	ByIdWithRelations(ctx context.Context, id uint64) (Block, error)
	ByHeightWithStats(ctx context.Context, height uint64) (Block, error)
	ByHash(ctx context.Context, hash []byte) (Block, error)
	ListWithStats(ctx context.Context, limit, offset uint64, order storage.SortOrder) ([]*Block, error)
}

// Block -
type Block struct {
	bun.BaseModel `bun:"table:block" comment:"Table with celestia blocks."`

	Id           uint64         `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height       pkgTypes.Level `bun:"height"                    comment:"The number (height) of this block" stats:"func:min max,filterable"`
	Time         time.Time      `bun:"time,pk,notnull"           comment:"The time of block"                 stats:"func:min max,filterable"`
	VersionBlock uint64         `bun:"version_block"             comment:"Block version"`
	VersionApp   uint64         `bun:"version_app"               comment:"App version"`

	MessageTypes types.MsgTypeBits `bun:"message_types,type:int8" comment:"Bit mask with containing messages"`

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
	ProposerAddress    pkgTypes.Hex `bun:"proposer_address"     comment:"Proposer address"`

	ChainId   string    `bun:"-"` // internal field for filling state
	Addresses []Address `bun:"-"` // internal field for balance passing

	Txs    []Tx       `bun:"rel:has-many"`
	Events []Event    `bun:"rel:has-many"`
	Stats  BlockStats `bun:"rel:has-one,join:height=height"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
