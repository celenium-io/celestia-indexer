package postgres

import (
	"context"
	"database/sql"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

// BlockStatsTestSuite -
type BlockStatsTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *BlockStatsTestSuite) SetupSuite() {
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
func (s *BlockStatsTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *BlockStatsTestSuite) TestByHeight() {
	tests := []struct {
		name   string
		height uint64
		want   storage.BlockStats
	}{
		{
			name:   "height 1000",
			height: 1000,
			want: storage.BlockStats{
				Id:           2,
				Height:       1000,
				SupplyChange: decimal.NewFromInt(30930476),
				MessagesCounts: map[storageTypes.MsgType]int64{
					storageTypes.MsgWithdrawDelegatorReward: 1,
					storageTypes.MsgDelegate:                1,
					storageTypes.MsgUnjail:                  1,
					storageTypes.MsgPayForBlobs:             1,
				},
			},
		},
		{
			name:   "height 999",
			height: 999,
			want: storage.BlockStats{
				Id:           1,
				Height:       999,
				SupplyChange: decimal.NewFromInt(20930476),
				MessagesCounts: map[storageTypes.MsgType]int64{
					storageTypes.MsgCreateValidator: 1,
				},
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
			s.Require().Equal(tt.want.MessagesCounts, got.MessagesCounts)
		})
	}
}

func TestSuiteBlockStats_Run(t *testing.T) {
	suite.Run(t, new(BlockStatsTestSuite))
}
