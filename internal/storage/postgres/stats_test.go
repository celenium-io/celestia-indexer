package postgres

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
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
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := database.NewPostgreSQLContainer(ctx, database.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb:latest-pg15",
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
	})
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
			want:  "3",
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
		defer ctxCancel()

		count, err := s.storage.Stats.Count(ctx, storage.CountRequest{
			Table: tests[i].table,
			From:  1672573739,
		})
		s.Require().NoError(err)
		s.Require().EqualValues(tests[i].want, count)
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
		defer ctxCancel()

		count, err := s.storage.Stats.Count(ctx, storage.CountRequest{
			Table: tests[i].table,
			From:  1693324139,
		})
		s.Require().NoError(err)
		s.Require().EqualValues("0", count)
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
			want:     "80410",
		}, {
			table:    "tx",
			column:   "fee",
			function: "min",
			want:     "80410",
		}, {
			table:    "tx",
			column:   "fee",
			function: "max",
			want:     "80410",
		},
	}

	for i := range tests {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer ctxCancel()

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
	}
}

func (s *StatsTestSuite) TestHistogram() {

	type test struct {
		timeframe storage.Timeframe
		table     string
		column    string
		function  string
		wantDate  time.Time
		want      string
	}

	tests := []test{
		// Block tests
		{
			timeframe: storage.TimeframeHour,
			table:     "block_stats",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "4599819996",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "block_stats",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "4599819996",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "block_stats",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "4599819996",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "block_stats",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "4599819996",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "block_stats",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "4599819996",
		}, {
			timeframe: storage.TimeframeHour,
			table:     "block_stats",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "2299909998",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "block_stats",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "2299909998",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "block_stats",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "2299909998",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "block_stats",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "2299909998",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "block_stats",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "2299909998",
		}, {
			timeframe: storage.TimeframeHour,
			table:     "block_stats",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "1726351723",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "block_stats",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "1726351723",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "block_stats",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "1726351723",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "block_stats",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "1726351723",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "block_stats",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "1726351723",
		}, {
			timeframe: storage.TimeframeHour,
			table:     "block_stats",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "2873468273",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "block_stats",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "2873468273",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "block_stats",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "2873468273",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "block_stats",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "2873468273",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "block_stats",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "2873468273",
		},
		// Tx tests
		{
			timeframe: storage.TimeframeHour,
			table:     "tx",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "241230",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "tx",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "241230",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "tx",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "241230",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "tx",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "241230",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "tx",
			column:    "fee",
			function:  "sum",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "241230",
		}, {
			timeframe: storage.TimeframeHour,
			table:     "tx",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "tx",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "tx",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "tx",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "tx",
			column:    "fee",
			function:  "avg",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeHour,
			table:     "tx",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "tx",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "tx",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "tx",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "tx",
			column:    "fee",
			function:  "min",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeHour,
			table:     "tx",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "tx",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "tx",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "tx",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "tx",
			column:    "fee",
			function:  "max",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "80410",
		},
	}

	for i := range tests {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer ctxCancel()

		histogram, err := s.storage.Stats.Histogram(ctx,
			storage.HistogramRequest{
				SummaryRequest: storage.SummaryRequest{
					CountRequest: storage.CountRequest{
						Table: tests[i].table,
						From:  1672573739,
					},
					Function: tests[i].function,
					Column:   tests[i].column,
				},
				Timeframe: tests[i].timeframe,
			})
		s.Require().NoError(err)
		s.Require().Len(histogram, 1)

		item := histogram[0]
		parts := strings.Split(item.Value, ".")
		s.Require().Equal(tests[i].want, parts[0])
		s.Require().True(item.Time.Equal(tests[i].wantDate))
	}
}

func (s *StatsTestSuite) TestHistogramCount() {
	type test struct {
		timeframe storage.Timeframe
		table     string
		wantDate  time.Time
		want      string
	}

	tests := []test{
		{
			timeframe: storage.TimeframeHour,
			table:     "block_stats",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "2",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "block_stats",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "2",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "block_stats",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "2",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "block_stats",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "2",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "block_stats",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "2",
		}, {
			timeframe: storage.TimeframeHour,
			table:     "tx",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "3",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "tx",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "3",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "tx",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "3",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "tx",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "3",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "tx",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "3",
		}, {
			timeframe: storage.TimeframeHour,
			table:     "message",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "5",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "message",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "5",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "message",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "5",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "message",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "5",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "message",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "5",
		}, {
			timeframe: storage.TimeframeHour,
			table:     "event",
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
			want:      "3",
		}, {
			timeframe: storage.TimeframeDay,
			table:     "event",
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
			want:      "3",
		}, {
			timeframe: storage.TimeframeWeek,
			table:     "event",
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
			want:      "3",
		}, {
			timeframe: storage.TimeframeMonth,
			table:     "event",
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			want:      "3",
		}, {
			timeframe: storage.TimeframeYear,
			table:     "event",
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      "3",
		},
	}

	for i := range tests {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer ctxCancel()

		histogram, err := s.storage.Stats.HistogramCount(ctx,
			storage.HistogramCountRequest{
				CountRequest: storage.CountRequest{
					Table: tests[i].table,
					From:  1672573739,
				},
				Timeframe: tests[i].timeframe,
			})
		s.Require().NoError(err)
		s.Require().Len(histogram, 1)

		item := histogram[0]
		s.Require().Equal(tests[i].want, item.Value)
		s.Require().True(item.Time.Equal(tests[i].wantDate))
	}
}

func TestSuiteStats_Run(t *testing.T) {
	suite.Run(t, new(StatsTestSuite))
}
