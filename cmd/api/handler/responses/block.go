package responses

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

type Block struct {
	Id                 uint64    `example:"321"                                                              json:"id"                   swaggertype:"integer"`
	Height             uint64    `example:"100"                                                              json:"height"               swaggertype:"integer"`
	Time               time.Time `example:"2023-07-04T03:10:57+00:00"                                        json:"time"                 swaggertype:"string"`
	VersionBlock       string    `example:"11"                                                               json:"version_block"        swaggertype:"string"`
	VersionApp         string    `example:"1"                                                                json:"version_app"          swaggertype:"string"`
	TxCount            uint64    `example:"12"                                                               json:"tx_count"             swaggertype:"integer"`
	EventsCount        uint64    `example:"18"                                                               json:"events_count"         swaggertype:"integer"`
	BlobsSize          uint64    `example:"12354"                                                            json:"blobs_size"           swaggertype:"integer"`
	Fee                string    `example:"28347628346"                                                      json:"fee"                  swaggertype:"string"`
	Hash               string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"hash"                 swaggertype:"string"`
	ParentHash         string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"parent_hash"          swaggertype:"string"`
	LastCommitHash     string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"last_commit_hash"     swaggertype:"string"`
	DataHash           string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"data_hash"            swaggertype:"string"`
	ValidatorsHash     string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"validators_hash"      swaggertype:"string"`
	NextValidatorsHash string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"next_validators_hash" swaggertype:"string"`
	ConsensusHash      string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"consensus_hash"       swaggertype:"string"`
	AppHash            string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"app_hash"             swaggertype:"string"`
	LastResultsHash    string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"last_results_hash"    swaggertype:"string"`
	EvidenceHash       string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"evidence_hash"        swaggertype:"string"`
	ProposerAddress    string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"proposer_address"     swaggertype:"string"`
	MessageTypes       []string  `example:"MsgSend,MsgUnjail"                                                json:"message_types"        swaggertype:"array,string"`
}

func NewBlock(block storage.Block) Block {
	return Block{
		Id:                 block.Id,
		Height:             uint64(block.Height),
		Time:               block.Time,
		VersionBlock:       strconv.FormatUint(block.VersionBlock, 10),
		VersionApp:         strconv.FormatUint(block.VersionApp, 10),
		TxCount:            block.TxCount,
		EventsCount:        block.EventsCount,
		BlobsSize:          block.BlobsSize,
		Fee:                block.Fee.String(),
		Hash:               hex.EncodeToString(block.Hash),
		ParentHash:         hex.EncodeToString(block.ParentHash),
		LastCommitHash:     hex.EncodeToString(block.LastCommitHash),
		DataHash:           hex.EncodeToString(block.DataHash),
		ValidatorsHash:     hex.EncodeToString(block.ValidatorsHash),
		NextValidatorsHash: hex.EncodeToString(block.NextValidatorsHash),
		ConsensusHash:      hex.EncodeToString(block.ConsensusHash),
		AppHash:            hex.EncodeToString(block.AppHash),
		LastResultsHash:    hex.EncodeToString(block.LastResultsHash),
		EvidenceHash:       hex.EncodeToString(block.EvidenceHash),
		ProposerAddress:    hex.EncodeToString(block.ProposerAddress),
		MessageTypes:       block.MessageTypes.Names(),
	}
}

func (Block) SearchType() string {
	return "block"
}
