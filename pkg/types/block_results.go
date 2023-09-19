package types

import (
	"time"

	"github.com/goccy/go-json"
)

// ResultBlockResults is an ABCI results from a block
// origin: github.com/celestiaorg/celestia-core@v1.26.2-tm-v0.34.28/rpc/core/types/responses.go
type ResultBlockResults struct {
	Height                Level                `json:"height,string"`
	TxsResults            []*ResponseDeliverTx `json:"txs_results"`
	BeginBlockEvents      []Event              `json:"begin_block_events"`
	EndBlockEvents        []Event              `json:"end_block_events"`
	ValidatorUpdates      []ValidatorUpdate    `json:"validator_updates"`
	ConsensusParamUpdates *ConsensusParams     `json:"consensus_param_updates"`
}

type ResponseDeliverTx struct {
	Code      uint32          `json:"code,omitempty"              protobuf:"varint,1,opt,name=code,proto3"`
	Data      json.RawMessage `json:"data,omitempty"              protobuf:"bytes,2,opt,name=data,proto3"`
	Log       string          `json:"log,omitempty"               protobuf:"bytes,3,opt,name=log,proto3"`
	Info      string          `json:"info,omitempty"              protobuf:"bytes,4,opt,name=info,proto3"`
	GasWanted int64           `json:"gas_wanted,omitempty,string" protobuf:"varint,5,opt,name=gas_wanted,proto3"`
	GasUsed   int64           `json:"gas_used,omitempty,string"   protobuf:"varint,6,opt,name=gas_used,proto3"`
	Events    []Event         `json:"events,omitempty"            protobuf:"bytes,7,rep,name=events,proto3"`
	Codespace string          `json:"codespace,omitempty"         protobuf:"bytes,8,opt,name=codespace,proto3"`
}

func (tx *ResponseDeliverTx) IsFailed() bool {
	return tx.Code != 0
}

// Event allows application developers to attach additional information to
// ResponseBeginBlock, ResponseEndBlock, ResponseCheckTx and ResponseDeliverTx.
// Later transactions may be queried using these events.
type Event struct {
	Type       string           `json:"type,omitempty"       protobuf:"bytes,1,opt,name=type,proto3"`
	Attributes []EventAttribute `json:"attributes,omitempty" protobuf:"bytes,2,rep,name=attributes,proto3"`
}

// EventAttribute is a single key-value pair, associated with an event.
type EventAttribute struct {
	Key   string `json:"key,omitempty"   protobuf:"bytes,1,opt,name=key,proto3"`
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value,proto3"`
	Index bool   `json:"index,omitempty" protobuf:"varint,3,opt,name=index,proto3"`
}

// ValidatorUpdate
type ValidatorUpdate struct {
	// PubKey any   `json:"pub_key"                protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3"` // crypto.PublicKey
	Power int64 `json:"power,omitempty,string" protobuf:"varint,2,opt,name=power,proto3"`
}

// ConsensusParams contains all consensus-relevant parameters
// that can be adjusted by the abci app
type ConsensusParams struct {
	Block     *BlockParams     `json:"block"     protobuf:"bytes,1,opt,name=block,proto3"`
	Evidence  *EvidenceParams  `json:"evidence"  protobuf:"bytes,2,opt,name=evidence,proto3"`
	Validator *ValidatorParams `json:"validator" protobuf:"bytes,3,opt,name=validator,proto3"`
	Version   *VersionParams   `json:"version"   protobuf:"bytes,4,opt,name=version,proto3"`
}

// BlockParams contains limits on the block size.
type BlockParams struct {
	// Note: must be greater than 0
	MaxBytes int64 `json:"max_bytes,omitempty,string" protobuf:"varint,1,opt,name=max_bytes,json=maxBytes,proto3"`
	// Note: must be greater or equal to -1
	MaxGas int64 `json:"max_gas,omitempty,string" protobuf:"varint,2,opt,name=max_gas,json=maxGas,proto3"`
}

// EvidenceParams determine how we handle evidence of malfeasance.
type EvidenceParams struct {
	// Max age of evidence, in blocks.
	//
	// The basic formula for calculating this is: MaxAgeDuration / {average block
	// time}.
	MaxAgeNumBlocks int64 `json:"max_age_num_blocks,omitempty,string" protobuf:"varint,1,opt,name=max_age_num_blocks,json=maxAgeNumBlocks,proto3"`
	// Max age of evidence, in time.
	//
	// It should correspond with an app's "unbonding period" or other similar
	// mechanism for handling [Nothing-At-Stake
	// attacks](https://github.com/ethereum/wiki/wiki/Proof-of-Stake-FAQ#what-is-the-nothing-at-stake-problem-and-how-can-it-be-fixed).
	MaxAgeDuration time.Duration `json:"max_age_duration,string" protobuf:"bytes,2,opt,name=max_age_duration,json=maxAgeDuration,proto3,stdduration"`
	// This sets the maximum size of total evidence in bytes that can be committed in a single block.
	// And should fall comfortably under the max block bytes.
	// Default is 1048576 or 1MB
	MaxBytes int64 `json:"max_bytes,omitempty,string" protobuf:"varint,3,opt,name=max_bytes,json=maxBytes,proto3"`
}

// ValidatorParams restrict the public key types validators can use.
// NOTE: uses ABCI pubkey naming, not Amino names.
type ValidatorParams struct {
	PubKeyTypes []string `json:"pub_key_types,omitempty" protobuf:"bytes,1,rep,name=pub_key_types,json=pubKeyTypes,proto3"`
}

// VersionParams contains the ABCI application version.
type VersionParams struct {
	AppVersion uint64 `json:"app_version,omitempty,string" protobuf:"varint,1,opt,name=app_version,json=appVersion,proto3"`
}
