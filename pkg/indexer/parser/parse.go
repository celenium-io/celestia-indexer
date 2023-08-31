package parser

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/shopspring/decimal"
	"strconv"
)

func (p *Parser) parse(ctx context.Context, resultBlock types.ResultBlock) error {
	p.log.Info().Int64("height", resultBlock.Block.Height).Msg("parsing block...")

	block := storage.Block{
		Height:       storage.Level(resultBlock.Block.Height),
		Time:         resultBlock.Block.Time,
		VersionBlock: strconv.FormatUint(resultBlock.Block.Version.Block, 10), // should we use string in storage type?
		VersionApp:   strconv.FormatUint(resultBlock.Block.Version.App, 10),   // should we use string in storage type?

		TxCount:       0, // TODO
		EventsCount:   0, // TODO
		MessageTypes:  storageTypes.MsgTypeBits{},
		NamespaceSize: 0, // "Summary block namespace size from pay for blob"` // should it be in block?

		Hash:               []byte(resultBlock.BlockID.Hash), // create a Hex type for common usage through indexer app
		ParentHash:         []byte(resultBlock.Block.LastBlockID.Hash),
		LastCommitHash:     []byte(resultBlock.Block.LastCommitHash),
		DataHash:           []byte(resultBlock.Block.DataHash),
		ValidatorsHash:     []byte(resultBlock.Block.ValidatorsHash),
		NextValidatorsHash: []byte(resultBlock.Block.NextValidatorsHash),
		ConsensusHash:      []byte(resultBlock.Block.ConsensusHash),
		AppHash:            []byte(resultBlock.Block.AppHash),
		LastResultsHash:    []byte(resultBlock.Block.LastResultsHash),
		EvidenceHash:       []byte(resultBlock.Block.EvidenceHash),
		ProposerAddress:    []byte(resultBlock.Block.ProposerAddress),

		Fee:     decimal.Zero, // TODO
		ChainId: resultBlock.Block.ChainID,

		Txs:    make([]storage.Tx, 0),    // TODO
		Events: make([]storage.Event, 0), // TODO
	}

	p.output.Push(block)
	return nil
}
