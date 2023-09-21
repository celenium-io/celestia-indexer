package rollback

import (
	"context"
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	indexerCfg "github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/node/mock"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"

	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/stretchr/testify/suite"
)

const testIndexerName = "test_indexer"

// ModuleTestSuite -
type ModuleTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       postgres.Storage
	api           *mock.MockApi
}

// SetupSuite -
func (s *ModuleTestSuite) SetupSuite() {
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

	st, err := postgres.Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	})
	s.Require().NoError(err)
	s.storage = st
}

// TearDownSuite -
func (s *ModuleTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *ModuleTestSuite) InitDb(path string) {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("timescaledb"),
		testfixtures.Directory(path),
		testfixtures.UseAlterConstraint(),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

func (s *ModuleTestSuite) InitApi(configureApi func()) {
	ctrl := gomock.NewController(s.T())
	s.api = mock.NewMockApi(ctrl)

	if configureApi != nil {
		configureApi()
	}
}

func GetResultBlock(hash types.Hex) types.ResultBlock {
	return types.ResultBlock{
		BlockID: types.BlockId{
			Hash: hash,
		},
	}
}

func (s *ModuleTestSuite) TestModule_SuccessOnRollbackTwoBlocks() {
	s.InitDb("../../../test/data/rollback")

	expectedHash, err := hex.DecodeString("5F7A8DDFE6136FE76B65B9066D4F816D707F28C05B3362D66084664C5B39BA98")
	s.Require().NoError(err)
	s.InitApi(func() {
		s.api.EXPECT().
			Block(gomock.Any(), types.Level(1001)).
			Return(GetResultBlock(types.Hex{42}), nil). // not equal with block in storage
			MaxTimes(1)

		s.api.EXPECT().
			Block(gomock.Any(), types.Level(1000)).
			Return(GetResultBlock(types.Hex{42}), nil). // not equal with block in storage
			MaxTimes(1)

		s.api.EXPECT().
			Block(gomock.Any(), types.Level(999)).
			Return(GetResultBlock(expectedHash), nil).
			MaxTimes(1)
	})

	rollbackModule := NewModule(
		s.storage.Transactable,
		s.storage.State,
		s.storage.Blocks,
		s.api,
		indexerCfg.Indexer{Name: testIndexerName},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	stateListener := modules.New("state-listener")
	stateListener.CreateInput("state")
	err = stateListener.AttachTo(&rollbackModule, OutputName, "state")
	s.Require().NoError(err)

	rollbackModule.Start(ctx)
	defer func() {
		cancel()
		s.Require().NoError(rollbackModule.Close())
	}()

	// Act
	rollbackModule.MustInput(InputName).Push(struct{}{})

	for {
		select {
		case <-ctx.Done():
			s.T().Error("stop by cancelled context")
			return
		case msg, ok := <-stateListener.MustInput("state").Listen():
			s.Require().True(ok, "received value should be delivered by successful send operation")

			state, ok := msg.(storage.State)
			s.Require().True(ok, "got wrong type %T", msg)

			s.Require().Equal(types.Level(999), state.LastHeight)
			s.Require().Equal(expectedHash, state.LastHash)
			s.Require().Equal("2023-07-04 03:10:56", state.LastTime.Format(time.DateTime))
			s.Require().Equal(uint64(1), state.TotalTx)
			s.Require().Equal(uint64(324234-100-900), state.TotalBlobsSize)
			s.Require().Equal(uint64(1000-3), state.TotalNamespaces)
			s.Require().Equal(uint64(12512357-1), state.TotalAccounts)

			expectedFee := decimal.NewFromInt(172635712635813).
				Sub(decimal.NewFromInt(497012)).
				Sub(decimal.NewFromInt(2873468273))
			s.Require().Equal(expectedFee, state.TotalFee)

			expectedSupply := decimal.NewFromInt(263471253613).
				Sub(decimal.NewFromInt(23590834)).
				Sub(decimal.NewFromInt(30930476))
			s.Require().Equal(expectedSupply, state.TotalSupply)

			return
		}
	}
}

func (s *ModuleTestSuite) TestModule_OnClosedInput() {
	s.InitDb("../../../test/data/rollback")

	s.InitApi(func() {
		s.api.EXPECT().
			Block(gomock.Any(), gomock.Any()).
			Return(GetResultBlock(types.Hex{42}), nil).
			MaxTimes(0)
	})

	rollbackModule := NewModule(
		s.storage.Transactable,
		s.storage.State,
		s.storage.Blocks,
		s.api,
		indexerCfg.Indexer{Name: testIndexerName},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	stopperListener := modules.New("stopper-listener")
	stopperListener.CreateInput("stop")
	err := stopperListener.AttachTo(&rollbackModule, StopOutput, "stop")
	s.Require().NoError(err)

	rollbackModule.Start(ctx)
	defer func() {
		cancel()
		s.Require().NoError(rollbackModule.Close())
	}()

	// Act
	err = rollbackModule.MustInput(InputName).Close()
	s.Require().NoError(err)

	for {
		select {
		case <-ctx.Done():
			s.T().Error("stop by cancelled context")
			return
		case msg, ok := <-stopperListener.MustInput("stop").Listen():
			s.Require().True(ok, "received stop signal should be delivered by successful send operation")
			s.Require().Equal(struct{}{}, msg)
			return
		}
	}
}

func TestSuiteModule_Run(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
