// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

// BlockStatsTestSuite -
type BlockStatsTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *BlockStatsTestSuite) SetupSuite() {
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
func (s *BlockStatsTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *BlockStatsTestSuite) TestByHeight() {
	tests := []struct {
		name   string
		height pkgTypes.Level
		want   storage.BlockStats
	}{
		{
			name:   "height 1000",
			height: 1000,
			want: storage.BlockStats{
				Id:           2,
				Height:       1000,
				SupplyChange: decimal.NewFromInt(30930476),
			},
		},
		{
			name:   "height 999",
			height: 999,
			want: storage.BlockStats{
				Id:           1,
				Height:       999,
				SupplyChange: decimal.NewFromInt(20930476),
			},
		},
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.storage.BlockStats.ByHeight(ctx, tt.height)
			s.Require().NoError(err)
			s.Require().Equal(tt.want.Id, got.Id)
			s.Require().Equal(tt.want.Height, got.Height)
			s.Require().Equal(tt.want.SupplyChange, got.SupplyChange)
		})
	}
}

func (s *BlockStatsTestSuite) TestLastFrom() {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	got, err := s.storage.BlockStats.LastFrom(ctx, 999, 1)
	s.Require().NoError(err)
	s.Require().Len(got, 1)

	item := got[0]
	s.Require().EqualValues(1, item.Id)
	s.Require().EqualValues(999, item.Height)
}

func TestSuiteBlockStats_Run(t *testing.T) {
	suite.Run(t, new(BlockStatsTestSuite))
}
