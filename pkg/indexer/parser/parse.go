// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	dCtx "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (p *Module) parse(b *types.BlockData) error {
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
			Fee:           storageTypes.NumericZero(),
			SupplyChange:  storageTypes.NumericZero(),
			InflationRate: storageTypes.NumericZero(),
			Commissions:   storageTypes.NumericZero(),
			Rewards:       storageTypes.NumericZero(),
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
	if err := p.parseValidatorUpdates(decodeCtx, b.ValidatorUpdates); err != nil {
		return errors.Wrap(err, "parseValidatorUpdates")
	}

	blockEvents, err := parseBlockEvents(decodeCtx, b, b.FinalizeBlockEvents, getFirstTxEvent(b.TxsResults))
	if err != nil {
		return errors.Wrap(err, "parsing begin end events")
	}
	decodeCtx.AddEvents(blockEvents...)

	p.Log.Info().
		Uint64("height", uint64(decodeCtx.Block.Height)).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block parsed")

	output := p.MustOutput(OutputName)
	output.Push(decodeCtx)

	return nil
}

func (p *Module) parseBlockSignatures(commit *types.Commit) []storage.BlockSignature {
	signs := make([]storage.BlockSignature, 0, len(commit.Signatures))
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
	if params.Evidence != nil {
		ctx.AddConstant(storageTypes.ModuleNameConsensus, "evidence_max_age_num_blocks", strconv.FormatInt(params.Evidence.MaxAgeNumBlocks, 10))
		ctx.AddConstant(storageTypes.ModuleNameConsensus, "evidence_max_age_duration", strconv.FormatInt(params.Evidence.MaxAgeDuration.Nanoseconds(), 10))
		ctx.AddConstant(storageTypes.ModuleNameConsensus, "evidence_max_bytes", strconv.FormatInt(params.Evidence.MaxBytes, 10))
	}
	if params.Block != nil {
		ctx.AddConstant(storageTypes.ModuleNameConsensus, "block_max_bytes", strconv.FormatInt(params.Block.MaxBytes, 10))
		ctx.AddConstant(storageTypes.ModuleNameConsensus, "block_max_gas", strconv.FormatInt(params.Block.MaxGas, 10))
	}
}

func getFirstTxEvent(results []types.ResponseDeliverTx) *types.Event {
	if len(results) == 0 {
		return nil
	}

	for i := 0; i < len(results)-1; i++ {
		if len(results[i].Events) > 0 {
			return &results[i].Events[0]
		}
	}

	return nil
}

func (p *Module) parseValidatorUpdates(ctx *dCtx.Context, updates []types.ValidatorUpdate) error {
	for i := range updates {
		if updates[i].PubKey.Sum.Type != "tendermint.crypto.PublicKey_Ed25519" {
			p.Log.Warn().Str("typ", updates[i].PubKey.Sum.Type).Msg("unknown pubkey validator type")
			continue
		}
		if len(updates[i].PubKey.Sum.Value.Ed25519) != 32 {
			p.Log.Warn().
				Int("length", len(updates[i].PubKey.Sum.Value.Ed25519)).
				Msg("invalid length of ed25519 pub key")
			continue
		}
		pk := ed25519.PubKey{Key: updates[i].PubKey.Sum.Value.Ed25519}
		consAddr := sdk.ConsAddress(pk.Address())
		hexAddr := strings.ToUpper(hex.EncodeToString(consAddr))
		// cmtjson omits power 0: the validator left the active set
		power := storageTypes.NumericZero()
		if updates[i].Power != nil {
			p, err := storageTypes.NumericFromString(*updates[i].Power)
			if err != nil {
				return errors.Wrapf(err, "validator update power: %s", *updates[i].Power)
			}
			power = p
		}

		// operator address and id are resolved by the storage module
		validator := storage.EmptyValidator()
		validator.ConsAddress = hexAddr
		validator.Power = &power
		validator.BondUpdatesCount = 1

		ctx.AddValidatorUpdate(storage.ValidatorBondUpdate{
			Height:    ctx.Block.Height,
			Time:      ctx.Block.Time,
			Power:     &power,
			Validator: &validator,
		})
	}
	return nil
}
