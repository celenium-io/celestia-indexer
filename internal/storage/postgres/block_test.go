// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

func (s *StorageTestSuite) TestBlockLast() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.Last(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().EqualValues(0, block.Stats.TxCount)

	hash, err := hex.DecodeString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")
	s.Require().NoError(err)
	s.Require().Equal(hash, block.Hash.Bytes())
	s.Require().Equal("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", block.Proposer.ConsAddress)
}

func (s *StorageTestSuite) TestBlockByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.ByHeight(ctx, 1000)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().Equal(storage.BlockStats{}, block.Stats)

	hash, err := hex.DecodeString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")
	s.Require().NoError(err)
	s.Require().Equal(hash, block.Hash.Bytes())
	s.Require().Equal("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", block.Proposer.ConsAddress)
}

func (s *StorageTestSuite) TestBlockByHeightWithStats() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.ByHeightWithStats(ctx, 1000)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().Equal("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", block.Proposer.ConsAddress)

	expectedStats := storage.BlockStats{
		Id:            2,
		Height:        1000,
		TxCount:       2,
		EventsCount:   0,
		BlobsSize:     1234,
		BlobsCount:    4,
		BlockTime:     11000,
		SupplyChange:  decimal.NewFromInt(30930476),
		InflationRate: decimal.NewFromFloat(0.08),
		Fee:           decimal.NewFromInt(2873468273),
	}
	s.Require().EqualValues(expectedStats.Id, block.Stats.Id)
	s.Require().EqualValues(expectedStats.Height, block.Stats.Height)
	s.Require().EqualValues(expectedStats.TxCount, block.Stats.TxCount)
	s.Require().EqualValues(expectedStats.EventsCount, block.Stats.EventsCount)
	s.Require().EqualValues(expectedStats.BlobsSize, block.Stats.BlobsSize)
	s.Require().EqualValues(expectedStats.BlobsCount, block.Stats.BlobsCount)
	s.Require().EqualValues(expectedStats.BlockTime, block.Stats.BlockTime)
	s.Require().EqualValues(expectedStats.SupplyChange.String(), block.Stats.SupplyChange.String())
	s.Require().EqualValues(expectedStats.InflationRate.String(), block.Stats.InflationRate.String())
	s.Require().EqualValues(expectedStats.Fee.String(), block.Stats.Fee.String())

	hash, err := hex.DecodeString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")
	s.Require().NoError(err)
	s.Require().Equal(hash, block.Hash.Bytes())
}

func (s *StorageTestSuite) TestBlockByIdWithRelations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.ByIdWithRelations(ctx, 2)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().Equal("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", block.Proposer.ConsAddress)

	expectedStats := storage.BlockStats{
		Id:            2,
		Height:        1000,
		TxCount:       2,
		EventsCount:   0,
		BlobsSize:     1234,
		BlockTime:     11000,
		SupplyChange:  decimal.NewFromInt(30930476),
		InflationRate: decimal.NewFromFloat(0.08),
		Fee:           decimal.NewFromInt(2873468273),
		BlobsCount:    4,
	}
	s.Require().EqualValues(expectedStats.Id, block.Stats.Id)
	s.Require().EqualValues(expectedStats.Height, block.Stats.Height)
	s.Require().EqualValues(expectedStats.TxCount, block.Stats.TxCount)
	s.Require().EqualValues(expectedStats.EventsCount, block.Stats.EventsCount)
	s.Require().EqualValues(expectedStats.BlobsSize, block.Stats.BlobsSize)
	s.Require().EqualValues(expectedStats.BlobsCount, block.Stats.BlobsCount)
	s.Require().EqualValues(expectedStats.BlockTime, block.Stats.BlockTime)
	s.Require().EqualValues(expectedStats.SupplyChange.String(), block.Stats.SupplyChange.String())
	s.Require().EqualValues(expectedStats.InflationRate.String(), block.Stats.InflationRate.String())
	s.Require().EqualValues(expectedStats.Fee.String(), block.Stats.Fee.String())

	hash, err := hex.DecodeString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")
	s.Require().NoError(err)
	s.Require().Equal(hash, block.Hash.Bytes())
}

func (s *StorageTestSuite) TestBlockByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")
	s.Require().NoError(err)

	block, err := s.storage.Blocks.ByHash(ctx, hash)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().EqualValues(2, block.Stats.TxCount)
	s.Require().Equal(hash, block.Hash.Bytes())
	s.Require().Equal("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", block.Proposer.ConsAddress)
}

func (s *StorageTestSuite) TestBlockListWithStats() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	blocks, err := s.storage.Blocks.ListWithStats(ctx, 10, 0, sdk.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(blocks, 2)

	block := blocks[0]
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().EqualValues(2, block.Stats.TxCount)
	s.Require().EqualValues(11000, block.Stats.BlockTime)
	s.Require().EqualValues(4, block.Stats.BlobsCount)
	s.Require().Equal("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", block.Proposer.ConsAddress)
}

func (s *StorageTestSuite) TestBlockListWithStatsAsc() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	blocks, err := s.storage.Blocks.ListWithStats(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(blocks, 2)

	block := blocks[1]
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().EqualValues(2, block.Stats.TxCount)
	s.Require().EqualValues(11000, block.Stats.BlockTime)
	s.Require().EqualValues(4, block.Stats.BlobsCount)
	s.Require().Equal("81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8", block.Proposer.ConsAddress)
}

func (s *StorageTestSuite) TestBlockByProposer() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	blocks, err := s.storage.Blocks.ByProposer(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(blocks, 2)

	block := blocks[0]
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().NotNil(block.Stats)
	s.Require().EqualValues(2, block.Stats.TxCount)
	s.Require().EqualValues(11000, block.Stats.BlockTime)
	s.Require().EqualValues(4, block.Stats.BlobsCount)
}

func (s *StorageTestSuite) TestBlockTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	expected, err := time.Parse(time.RFC3339, "2023-07-04T03:10:57Z")
	s.Require().NoError(err)

	blockTime, err := s.storage.Blocks.Time(ctx, 1000)
	s.Require().NoError(err)
	s.Require().Equal(expected.UTC(), blockTime.UTC())
}
