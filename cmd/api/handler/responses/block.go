// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"strconv"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
)

type Block struct {
	Id                 uint64       `example:"321"                                                              json:"id"                   swaggertype:"integer"`
	Height             uint64       `example:"100"                                                              json:"height"               swaggertype:"integer"`
	Time               time.Time    `example:"2023-07-04T03:10:57+00:00"                                        json:"time"                 swaggertype:"string"`
	VersionBlock       string       `example:"11"                                                               json:"version_block"        swaggertype:"string"`
	VersionApp         string       `example:"1"                                                                json:"version_app"          swaggertype:"string"`
	Hash               pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"hash"                 swaggertype:"string"`
	ParentHash         pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"parent_hash"          swaggertype:"string"`
	LastCommitHash     pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"last_commit_hash"     swaggertype:"string"`
	DataHash           pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"data_hash"            swaggertype:"string"`
	ValidatorsHash     pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"validators_hash"      swaggertype:"string"`
	NextValidatorsHash pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"next_validators_hash" swaggertype:"string"`
	ConsensusHash      pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"consensus_hash"       swaggertype:"string"`
	AppHash            pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"app_hash"             swaggertype:"string"`
	LastResultsHash    pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"last_results_hash"    swaggertype:"string"`
	EvidenceHash       pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"evidence_hash"        swaggertype:"string"`
	ProposerAddress    pkgTypes.Hex `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"proposer_address"     swaggertype:"string"`

	MessageTypes []types.MsgType `example:"MsgSend,MsgUnjail" json:"message_types" swaggertype:"array,string"`

	Stats *BlockStats `json:"stats,omitempty"`
}

func NewBlock(block storage.Block, withStats bool) Block {
	result := Block{
		Id:                 block.Id,
		Height:             uint64(block.Height),
		Time:               block.Time,
		VersionBlock:       strconv.FormatUint(block.VersionBlock, 10),
		VersionApp:         strconv.FormatUint(block.VersionApp, 10),
		Hash:               block.Hash,
		ParentHash:         block.ParentHash,
		LastCommitHash:     block.LastCommitHash,
		DataHash:           block.DataHash,
		ValidatorsHash:     block.ValidatorsHash,
		NextValidatorsHash: block.NextValidatorsHash,
		ConsensusHash:      block.ConsensusHash,
		AppHash:            block.AppHash,
		LastResultsHash:    block.LastResultsHash,
		EvidenceHash:       block.EvidenceHash,
		ProposerAddress:    block.ProposerAddress,
		MessageTypes:       block.MessageTypes.Names(),
	}

	if withStats {
		result.Stats = NewBlockStats(block.Stats)
	}
	return result
}

func (Block) SearchType() string {
	return "block"
}

type BlockStats struct {
	TxCount        int64                   `example:"12"                              json:"tx_count"        swaggertype:"integer"`
	EventsCount    int64                   `example:"18"                              json:"events_count"    swaggertype:"integer"`
	BlobsSize      int64                   `example:"12354"                           json:"blobs_size"      swaggertype:"integer"`
	Fee            string                  `example:"28347628346"                     json:"fee"             swaggertype:"string"`
	SupplyChange   string                  `example:"8635234"                         json:"supply_change"   swaggertype:"string"`
	InflationRate  string                  `example:"0.0800000"                       json:"inflation_rate"  swaggertype:"string"`
	BlockTime      uint64                  `example:"12354"                           json:"block_time"      swaggertype:"integer"`
	MessagesCounts map[types.MsgType]int64 `example:"{MsgPayForBlobs:10,MsgUnjail:1}" json:"messages_counts" swaggertype:"string"`
}

func NewBlockStats(stats storage.BlockStats) *BlockStats {
	return &BlockStats{
		TxCount:        stats.TxCount,
		EventsCount:    stats.EventsCount,
		BlobsSize:      stats.BlobsSize,
		Fee:            stats.Fee.String(),
		SupplyChange:   stats.SupplyChange.String(),
		InflationRate:  stats.InflationRate.String(),
		BlockTime:      stats.BlockTime,
		MessagesCounts: stats.MessagesCounts,
	}
}
