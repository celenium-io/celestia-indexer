// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	dCtx "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (p *Module) parse(b types.BlockData) error {
	start := time.Now()
	p.Log.Info().
		Int64("height", b.Block.Height).
		Msg("parsing block...")

	decodeCtx := dCtx.NewContext()

	decodeCtx.Block = &storage.Block{
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

		Events: nil,

		Stats: storage.BlockStats{
			Height:        b.Height,
			Time:          b.Block.Time,
			TxCount:       int64(len(b.Block.Txs)),
			EventsCount:   int64(len(b.FinalizeBlockEvents)),
			BlobsSize:     0,
			Fee:           decimal.Zero,
			SupplyChange:  decimal.Zero,
			InflationRate: decimal.Zero,
			Commissions:   decimal.Zero,
			Rewards:       decimal.Zero,
			SquareSize:    b.Block.SquareSize,
		},
	}

	txs, err := p.parseTxs(decodeCtx, b)
	if err != nil {
		return errors.Wrapf(err, "while parsing block on level=%d", b.Height)
	}
	decodeCtx.Block.Txs = txs

	for i := range b.Block.Txs {
		decodeCtx.Block.Stats.BytesInBlock += int64(len(b.Block.Txs[i]))
	}

	decodeCtx.Block.BlockSignatures = p.parseBlockSignatures(b.Block.LastCommit)
	p.parseConsensusParamUpdates(decodeCtx, b.ConsensusParamUpdates)

	if isEventDuplicated(b.ResultBlockResults) {
		b.FinalizeBlockEvents = b.FinalizeBlockEvents[:len(b.FinalizeBlockEvents)-decodeCtx.TxEventsCount]
	}

	decodeCtx.Block.Events, err = parseEvents(decodeCtx, b, b.FinalizeBlockEvents)
	if err != nil {
		return errors.Wrap(err, "parsing begin end events")
	}

	p.Log.Info().
		Uint64("height", uint64(decodeCtx.Block.Height)).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block parsed")

	output := p.MustOutput(OutputName)
	output.Push(decodeCtx)

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

func (p *Module) parseConsensusParamUpdates(ctx *dCtx.Context, params *types.ConsensusParams) {
	ctx.AddConstant(storageTypes.ModuleNameConsensus, "evidence_max_age_num_blocks", fmt.Sprintf("%d", params.Evidence.MaxAgeNumBlocks))
	ctx.AddConstant(storageTypes.ModuleNameConsensus, "evidence_max_age_duration", params.Evidence.MaxAgeDuration.String())
}

func isEventDuplicated(results types.ResultBlockResults) bool {
	if len(results.TxsResults) == 0 || len(results.FinalizeBlockEvents) == 0 {
		return false
	}
	lastEvents := results.TxsResults[len(results.TxsResults)-1].Events
	return lastEvents[len(lastEvents)-1].Compare(results.FinalizeBlockEvents[len(results.FinalizeBlockEvents)-1])
}
