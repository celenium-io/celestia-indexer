// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"
	"encoding/hex"
	"strings"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
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

	var bytesInBlock int64
	for i := range b.Block.Txs {
		bytesInBlock += int64(len(b.Block.Txs[i]))
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
		ProposerAddress:    b.Block.ProposerAddress.String(),

		ChainId: b.Block.ChainID,

		Txs:    txs,
		Events: nil,

		Stats: storage.BlockStats{
			Height:        b.Height,
			Time:          b.Block.Time,
			TxCount:       int64(len(b.Block.Data.Txs)),
			EventsCount:   int64(len(b.BeginBlockEvents) + len(b.EndBlockEvents)),
			BlobsSize:     0,
			Fee:           decimal.Zero,
			SupplyChange:  decimal.Zero,
			InflationRate: decimal.Zero,
			BytesInBlock:  bytesInBlock,
		},
	}

	block.BlockSignatures = p.parseBlockSignatures(b.Block.LastCommit)

	allEvents := make([]storage.Event, 0)

	block.Events = parseEvents(b, b.ResultBlockResults.BeginBlockEvents)
	allEvents = append(allEvents, block.Events...)

	for i := range txs {
		block.Stats.Fee = block.Stats.Fee.Add(txs[i].Fee)
		block.MessageTypes.Set(txs[i].MessageTypes.Bits)
		block.Stats.BlobsSize += txs[i].BlobsSize
		block.Stats.GasLimit += txs[i].GasWanted
		block.Stats.GasUsed += txs[i].GasUsed
		allEvents = append(allEvents, txs[i].Events...)
	}

	endEvents := parseEvents(b, b.ResultBlockResults.EndBlockEvents)
	block.Events = append(block.Events, endEvents...)
	allEvents = append(allEvents, endEvents...)

	var eventsResult eventsResult
	if err := eventsResult.Fill(allEvents); err != nil {
		return err
	}

	block.Stats.InflationRate = eventsResult.InflationRate
	block.Stats.SupplyChange = eventsResult.SupplyChange
	block.Addresses = eventsResult.Addresses

	p.Log.Info().
		Uint64("height", uint64(block.Height)).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block parsed")

	output := p.MustOutput(OutputName)
	output.Push(block)
	return nil
}

func (p *Module) parseBlockSignatures(commit *types.Commit) []storage.BlockSignature {
	signs := make([]storage.BlockSignature, 0)
	for i := range commit.Signatures {
		if commit.Signatures[i].BlockIDFlag != 2 {
			continue
		}
		signs = append(signs, storage.BlockSignature{
			Height: types.Level(commit.Height),
			Time:   commit.Signatures[i].Timestamp,
			Validator: &storage.Validator{
				ConsAddress: strings.ToUpper(hex.EncodeToString(commit.Signatures[i].ValidatorAddress)),
			},
		})
	}
	return signs
}
