package parser

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

func (p *Parser) parse(ctx context.Context, b types.BlockData) error {
	p.log.Info().Int64("height", b.Block.Height).Msg("parsing block...")

	block := storage.Block{
		Height:       b.Height,
		Time:         b.Block.Time,
		VersionBlock: b.Block.Version.Block,
		VersionApp:   b.Block.Version.App,

		TxCount:      uint64(len(b.Block.Data.Txs)),
		EventsCount:  uint64(len(b.BeginBlockEvents) + len(b.EndBlockEvents)),
		MessageTypes: storageTypes.MsgTypeBits{}, // TODO init
		BlobsSize:    0,

		Hash:               []byte(b.BlockID.Hash), // TODO create a Hex type for common usage through indexer app
		ParentHash:         []byte(b.Block.LastBlockID.Hash),
		LastCommitHash:     b.Block.LastCommitHash,
		DataHash:           b.Block.DataHash,
		ValidatorsHash:     b.Block.ValidatorsHash,
		NextValidatorsHash: b.Block.NextValidatorsHash,
		ConsensusHash:      b.Block.ConsensusHash,
		AppHash:            b.Block.AppHash,
		LastResultsHash:    b.Block.LastResultsHash,
		EvidenceHash:       b.Block.EvidenceHash,
		ProposerAddress:    b.Block.ProposerAddress,

		Fee:     decimal.Zero, // TODO sum of auth_info.fee // RESEARCH: done
		ChainId: b.Block.ChainID,

		Txs:    parseTxs(b),
		Events: nil,
	}

	block.Events = parseEvents(b, b.ResultBlockResults.BeginBlockEvents)
	endEvents := parseEvents(b, b.ResultBlockResults.EndBlockEvents)
	block.Events = append(block.Events, endEvents...)

	p.output.Push(block)
	return nil
}
