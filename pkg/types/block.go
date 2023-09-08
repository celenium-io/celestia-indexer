package types

import (
	"time"

	"github.com/tendermint/tendermint/types"
)

// ResultBlock is a single block (with meta)
type ResultBlock struct {
	BlockID types.BlockID `json:"block_id"`
	Block   *Block        `json:"block"`
}

// Block defines the atomic unit of a CometBFT blockchain.
type Block struct {
	Header     `json:"header"`
	Data       `json:"data"`
	Evidence   types.EvidenceData `json:"evidence"`
	LastCommit *Commit            `json:"last_commit"`
}

// Consensus captures the consensus rules for processing a block in the blockchain,
// including all blockchain data structures and the rules of the application's
// state transition machine.
type Consensus struct {
	Block uint64 `json:"block,omitempty,string" protobuf:"varint,1,opt,name=block,proto3"`
	App   uint64 `json:"app,omitempty,string"   protobuf:"varint,2,opt,name=app,proto3"`
}

// Header defines the structure of a CometBFT block header.
type Header struct {
	// basic block info
	Version Consensus `json:"version"`
	ChainID string    `json:"chain_id"`
	Height  int64     `json:"height,string"`
	Time    time.Time `json:"time"`

	// prev block info
	LastBlockID types.BlockID `json:"last_block_id"` // TODO check if hash unmarshales correctly with HexBytes type

	// hashes of block data
	LastCommitHash Hex `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       Hex `json:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash     Hex `json:"validators_hash"`      // validators for the current block
	NextValidatorsHash Hex `json:"next_validators_hash"` // validators for the next block
	ConsensusHash      Hex `json:"consensus_hash"`       // consensus params for current block
	AppHash            Hex `json:"app_hash"`             // state after txs from the previous block
	// root hash of all results from the txs from the previous block
	// see `deterministicResponseDeliverTx` to understand which parts of a tx are hashed into here
	LastResultsHash Hex `json:"last_results_hash"`

	// consensus info
	EvidenceHash    Hex `json:"evidence_hash"`    // evidence included in the block
	ProposerAddress Hex `json:"proposer_address"` // original proposer of the block
}

// Data contains all the available Data of the block.
// Data with reserved namespaces (Txs, IntermediateStateRoots, Evidence) and
// Celestia application-specific Blobs.
type Data struct {
	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Txs types.Txs `json:"txs"`

	// SquareSize is the size of the square after splitting all the block data
	// into shares. The erasure data is discarded after generation, and keeping this
	// value avoids unnecessarily regenerating all the shares when returning
	// proofs that some element was included in the block
	SquareSize uint64 `json:"square_size,string"`
}

// Commit contains the evidence that a block was committed by a set of validators.
// NOTE: Commit is empty for height 1, but never nil.
type Commit struct {
	// NOTE: The signatures are in order of address to preserve the bonded
	// ValidatorSet order.
	// Any peer with a block can gossip signatures by index with a peer without
	// recalculating the active ValidatorSet.
	Height     int64             `json:"height,string"`
	Round      int32             `json:"round"`
	BlockID    types.BlockID     `json:"block_id"`
	Signatures []types.CommitSig `json:"signatures"`
}
