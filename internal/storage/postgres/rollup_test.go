// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestRollupLeaderboard() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	for _, column := range []string{
		sizeColumn, blobsCountColumn, timeColumn, feeColumn, "",
	} {

		rollups, err := s.storage.Rollup.Leaderboard(ctx, storage.LeaderboardFilters{
			SortField: column,
			Sort:      sdk.SortOrderDesc,
			Limit:     10,
			Offset:    0,
		})
		s.Require().NoError(err, column)
		s.Require().Len(rollups, 3, column)

		rollup := rollups[0]
		s.Require().EqualValues("Rollup 3", rollup.Name, column)
		s.Require().EqualValues("The third", rollup.Description, column)
		s.Require().EqualValues(34, rollup.Size, column)
		s.Require().EqualValues(3, rollup.BlobsCount, column)
		s.Require().False(rollup.LastActionTime.IsZero())
		s.Require().False(rollup.FirstActionTime.IsZero())
		s.Require().Equal("7000", rollup.Fee.String())
		s.Require().EqualValues(0.6363636363636364, rollup.FeePct)
		s.Require().EqualValues(0.42857142857142855, rollup.BlobsCountPct)
		s.Require().EqualValues(0.3953488372093023, rollup.SizePct)
	}
}

func (s *StorageTestSuite) TestRollupLeaderboardWithCategory() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	for _, column := range []string{
		sizeColumn, blobsCountColumn, timeColumn, feeColumn, "",
	} {

		rollups, err := s.storage.Rollup.Leaderboard(ctx, storage.LeaderboardFilters{
			SortField: column,
			Sort:      sdk.SortOrderDesc,
			Limit:     10,
			Offset:    0,
			Category:  []types.RollupCategory{types.RollupCategoryNft},
		})
		s.Require().NoError(err, column)
		s.Require().Len(rollups, 1, column)

		rollup := rollups[0]
		s.Require().EqualValues("Rollup 3", rollup.Name, column)
		s.Require().EqualValues("The third", rollup.Description, column)
		s.Require().EqualValues(34, rollup.Size, column)
		s.Require().EqualValues(3, rollup.BlobsCount, column)
		s.Require().False(rollup.LastActionTime.IsZero())
		s.Require().False(rollup.FirstActionTime.IsZero())
		s.Require().Equal("7000", rollup.Fee.String())
		s.Require().EqualValues(0.6363636363636364, rollup.FeePct)
		s.Require().EqualValues(0.42857142857142855, rollup.BlobsCountPct)
		s.Require().EqualValues(0.3953488372093023, rollup.SizePct)
		s.Require().EqualValues("nft", rollup.Category)
	}
}

func (s *StorageTestSuite) TestRollupLeaderboardWithType() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	for _, column := range []string{
		sizeColumn, blobsCountColumn, timeColumn, feeColumn, "",
	} {

		rollups, err := s.storage.Rollup.Leaderboard(ctx, storage.LeaderboardFilters{
			SortField: column,
			Sort:      sdk.SortOrderDesc,
			Limit:     10,
			Offset:    0,
			Type:      []types.RollupType{types.RollupTypeSovereign},
		})
		s.Require().NoError(err, column)
		s.Require().Len(rollups, 1, column)

		rollup := rollups[0]
		s.Require().EqualValues("Rollup 3", rollup.Name, column)
		s.Require().EqualValues("The third", rollup.Description, column)
		s.Require().EqualValues(34, rollup.Size, column)
		s.Require().EqualValues(3, rollup.BlobsCount, column)
		s.Require().False(rollup.LastActionTime.IsZero())
		s.Require().False(rollup.FirstActionTime.IsZero())
		s.Require().Equal("7000", rollup.Fee.String())
		s.Require().EqualValues(0.6363636363636364, rollup.FeePct)
		s.Require().EqualValues(0.42857142857142855, rollup.BlobsCountPct)
		s.Require().EqualValues(0.3953488372093023, rollup.SizePct)
		s.Require().EqualValues("sovereign", rollup.Type)
	}
}

func (s *StorageTestSuite) TestRollupLeaderboardDay() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard_day;")
	s.Require().NoError(err)

	for _, column := range []string{
		"avg_size", blobsCountColumn, "total_size", "total_fee", "throughput", "namespace_count", "pfb_count", "mb_price", "",
	} {

		rollups, err := s.storage.Rollup.LeaderboardDay(ctx, storage.LeaderboardFilters{
			SortField: column,
			Sort:      sdk.SortOrderDesc,
			Limit:     10,
			Offset:    0,
		})
		s.Require().NoError(err, column)
		s.Require().Len(rollups, 0, column)
	}
}

func (s *StorageTestSuite) TestRollupNamespaces() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	nsIds, err := s.storage.Rollup.Namespaces(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(nsIds, 2)
}

func (s *StorageTestSuite) TestRollupProviders() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	providers, err := s.storage.Rollup.Providers(ctx, 1)
	s.Require().NoError(err)
	s.Require().Len(providers, 2)
}

func (s *StorageTestSuite) TestRollupSeries() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, tf := range []string{
		"day", "hour", "month",
	} {
		for _, column := range []string{
			"size", "blobs_count", "size_per_blob", "fee",
		} {
			series, err := s.storage.Rollup.Series(ctx, 1, tf, column, storage.SeriesRequest{
				From: time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC),
			})
			s.Require().NoError(err)
			s.Require().Len(series, 2)

		}
	}
}

func (s *StorageTestSuite) TestRollupBySlug() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	rollup, err := s.storage.Rollup.BySlug(ctx, "rollup_3")
	s.Require().NoError(err)

	s.Require().EqualValues("Rollup 3", rollup.Name)
	s.Require().EqualValues("The third", rollup.Description)
	s.Require().EqualValues(34, rollup.Size)
	s.Require().EqualValues(3, rollup.BlobsCount)
	s.Require().False(rollup.LastActionTime.IsZero())
	s.Require().False(rollup.FirstActionTime.IsZero())
	s.Require().Equal("7000", rollup.Fee.String())
	s.Require().EqualValues(0.6363636363636364, rollup.FeePct)
	s.Require().EqualValues(0.42857142857142855, rollup.BlobsCountPct)
	s.Require().EqualValues(0.3953488372093023, rollup.SizePct)
}

func (s *StorageTestSuite) TestRollupById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	rollup, err := s.storage.Rollup.ById(ctx, 3)
	s.Require().NoError(err)

	s.Require().EqualValues("Rollup 3", rollup.Name)
	s.Require().EqualValues("The third", rollup.Description)
	s.Require().EqualValues(34, rollup.Size)
	s.Require().EqualValues(3, rollup.BlobsCount)
	s.Require().False(rollup.LastActionTime.IsZero())
	s.Require().False(rollup.FirstActionTime.IsZero())
	s.Require().Equal("7000", rollup.Fee.String())
	s.Require().EqualValues(0.6363636363636364, rollup.FeePct)
	s.Require().EqualValues(0.42857142857142855, rollup.BlobsCountPct)
	s.Require().EqualValues(0.3953488372093023, rollup.SizePct)
}

func (s *StorageTestSuite) TestRollupsByNamespace() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	rollups, err := s.storage.Rollup.RollupsByNamespace(ctx, 2, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(rollups, 2)

	rollup := rollups[0]
	s.Require().Greater(rollup.Id, uint64(0))
	s.Require().NotEmpty(rollup.Name)
}

func (s *StorageTestSuite) TestRollupDistribution() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, groupBy := range []string{
		"day", "hour",
	} {
		for _, series := range []string{
			"size", "blobs_count", "size_per_blob", "fee_per_blob",
		} {
			items, err := s.storage.Rollup.Distribution(ctx, 1, series, groupBy)
			s.Require().NoError(err)
			s.Require().Len(items, 1, groupBy, series)

		}
	}
}

func (s *StorageTestSuite) TestRollupAllSeries() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Rollup.AllSeries(ctx)
	s.Require().NoError(err)
	s.Require().Len(items, 5)
}

func (s *StorageTestSuite) TestRollupStatsGrouping() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	column := "stack"
	rollups, err := s.storage.Rollup.RollupStatsGrouping(ctx, storage.RollupGroupStatsFilters{
		Func:   "sum",
		Column: column,
	})
	s.Require().NoError(err, column)
	s.Require().Len(rollups, 2, column)

	rollup := rollups[1]
	s.Require().EqualValues(4000, rollup.Fee, column)
	s.Require().EqualValues(52, rollup.Size, column)
	s.Require().EqualValues(4, rollup.BlobsCount, column)
	s.Require().EqualValues("OP Stack", rollup.Group, column)
}
