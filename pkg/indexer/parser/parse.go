package parser

import (
	"context"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (p *Module) parse(ctx context.Context, b types.BlockData) error {
	start := time.Now()
	p.Log.Info().
		Int64("height", b.Block.Height).
		Msg("parsing block...")

	txs, err := parseTxs(b)
	if err != nil {
		return errors.Wrapf(err, "while parsing block on level=%d", b.Height)
	}

	block := storage.Block{
		Height:       b.Height,
		Time:         b.Block.Time,
		VersionBlock: b.Block.Version.Block,
		VersionApp:   b.Block.Version.App,

		MessageTypes: storageTypes.NewMsgTypeBitMask(),

		Hash:               []byte(b.BlockID.Hash),
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

		ChainId: b.Block.ChainID,

		Txs:    txs,
		Events: nil,

		Stats: storage.BlockStats{
			Height:        b.Height,
			Time:          b.Block.Time,
			TxCount:       uint64(len(b.Block.Data.Txs)),
			EventsCount:   uint64(len(b.BeginBlockEvents) + len(b.EndBlockEvents)),
			BlobsSize:     0,
			Fee:           decimal.Zero,
			SupplyChange:  decimal.Zero,
			InflationRate: decimal.Zero,
		},
	}

	for _, tx := range txs {
		block.Stats.Fee = block.Stats.Fee.Add(tx.Fee)
		block.MessageTypes.Set(tx.MessageTypes.Bits)
		block.Stats.BlobsSize += tx.BlobsSize
	}

	block.Events = parseEvents(b, b.ResultBlockResults.BeginBlockEvents)

	var bbeResult beginBlockEventsResult
	if err := bbeResult.Fill(block.Events); err != nil {
		return err
	}

	block.Stats.InflationRate = bbeResult.InflationRate
	block.Stats.SupplyChange = bbeResult.SupplyChange

	endEvents := parseEvents(b, b.ResultBlockResults.EndBlockEvents)
	block.Events = append(block.Events, endEvents...)

	p.Log.Debug().
		Uint64("height", uint64(block.Height)).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block parsed")

	output := p.MustOutput(OutputName)
	output.Push(block)
	return nil
}
