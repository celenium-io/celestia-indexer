// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
)

// StatsTestSuite -
type StatsTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *StatsTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer ctxCancel()

	psqlContainer, err := database.NewPostgreSQLContainer(ctx, database.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb-ha:pg15-latest",
	})
	s.Require().NoError(err)
	s.psqlContainer = psqlContainer

	strg, err := Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}, "../../../database")
	s.Require().NoError(err)
	s.storage = strg

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("timescaledb"),
		testfixtures.Directory("../../../test/data"),
		testfixtures.UseAlterConstraint(),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

// TearDownSuite -
func (s *StatsTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *StatsTestSuite) TestCount() {
	type test struct {
		table string
		want  string
	}

	tests := []test{
		{
			table: "block_stats",
			want:  "2",
		}, {
			table: "tx",
			want:  "4",
		}, {
			table: "event",
			want:  "3",
		}, {
			table: "message",
			want:  "5",
		},
	}

	for i := range tests {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)

		count, err := s.storage.Stats.Count(ctx, storage.CountRequest{
			Table: tests[i].table,
			From:  1672573739,
		})
		s.Require().NoError(err)
		s.Require().EqualValues(tests[i].want, count)

		ctxCancel()
	}
}

func (s *StatsTestSuite) TestCountNoData() {
	type test struct {
		table string
	}

	tests := []test{
		{
			table: "block_stats",
		}, {
			table: "tx",
		}, {
			table: "event",
		}, {
			table: "message",
		},
	}

	for i := range tests {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)

		count, err := s.storage.Stats.Count(ctx, storage.CountRequest{
			Table: tests[i].table,
			From:  1693324139,
		})
		s.Require().NoError(err)
		s.Require().EqualValues("0", count)

		ctxCancel()
	}
}

func (s *StatsTestSuite) TestSummaryBlock() {
	type test struct {
		table    string
		column   string
		function string
		want     string
	}

	tests := []test{
		// Block tests
		{
			table:    "block_stats",
			column:   "fee",
			function: "sum",
			want:     "4599819996",
		}, {
			table:    "block_stats",
			column:   "fee",
			function: "avg",
			want:     "2299909998",
		}, {
			table:    "block_stats",
			column:   "fee",
			function: "min",
			want:     "1726351723",
		}, {
			table:    "block_stats",
			column:   "fee",
			function: "max",
			want:     "2873468273",
		},
		// Tx tests
		{
			table:    "tx",
			column:   "fee",
			function: "sum",
			want:     "241230",
		}, {
			table:    "tx",
			column:   "fee",
			function: "avg",
			want:     "60307",
		}, {
			table:    "tx",
			column:   "fee",
			function: "min",
			want:     "0",
		}, {
			table:    "tx",
			column:   "fee",
			function: "max",
			want:     "80410",
		},
	}

	for i := range tests {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)

		summary, err := s.storage.Stats.Summary(ctx, storage.SummaryRequest{
			CountRequest: storage.CountRequest{
				Table: tests[i].table,
				From:  1672573739,
			},
			Function: tests[i].function,
			Column:   tests[i].column,
		})
		s.Require().NoError(err)

		parts := strings.Split(summary, ".")
		s.Require().Equal(tests[i].want, parts[0])

		ctxCancel()
	}
}

func (s *StatsTestSuite) TestTPS() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tps, err := s.storage.Stats.TPS(ctx)
	s.Require().NoError(err)

	s.Require().EqualValues(0, tps.High)
	s.Require().EqualValues(0, tps.Low)
	s.Require().EqualValues(0, tps.Current)
	s.Require().EqualValues(0, tps.ChangeLastHourPct)
}

func (s *StatsTestSuite) TestSeries() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Stats.Series(ctx, storage.TimeframeHour, storage.SeriesBlobsSize, storage.SeriesRequest{})
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func (s *StatsTestSuite) TestSeriesWithFrom() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Stats.Series(ctx, storage.TimeframeDay, storage.SeriesBlobsSize, storage.SeriesRequest{
		From: time.Unix(1701192801, 0).UTC(),
	})
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *StatsTestSuite) TestCumulativeSeries() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Stats.CumulativeSeries(ctx, storage.TimeframeDay, storage.SeriesBlobsSize, storage.SeriesRequest{})
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func (s *StatsTestSuite) TestCumulativeSeriesWithFrom() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Stats.CumulativeSeries(ctx, storage.TimeframeDay, storage.SeriesBlobsSize, storage.SeriesRequest{
		From: time.Unix(1701192801, 0).UTC(),
	})
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *StatsTestSuite) TestNamespaceSeries() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Stats.NamespaceSeries(ctx, storage.TimeframeHour, storage.SeriesNsSize, 1, storage.SeriesRequest{})
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func (s *StatsTestSuite) TestSquareSize() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	from := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	items, err := s.storage.Stats.SquareSize(ctx, &from, nil)
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func TestSuiteStats_Run(t *testing.T) {
	suite.Run(t, new(StatsTestSuite))
}
