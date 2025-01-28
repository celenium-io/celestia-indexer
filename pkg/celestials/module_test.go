// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package celestials

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/celestials"
	celestialsMock "github.com/celenium-io/celestia-indexer/internal/celestials/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

const testIndexerName = "indexer"
const network = "celestia"

// ModuleTestSuite -
type ModuleTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       postgres.Storage
	ctrl          *gomock.Controller
	api           *celestialsMock.MockAPI
}

// SetupSuite -
func (s *ModuleTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer ctxCancel()

	psqlContainer, err := database.NewPostgreSQLContainer(ctx, database.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb-ha:pg15.8-ts2.17.0-all",
	})
	s.Require().NoError(err)
	s.psqlContainer = psqlContainer

	strg, err := postgres.Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}, "../../database", false)
	s.Require().NoError(err)
	s.storage = strg

	s.ctrl = gomock.NewController(s.T())
	s.api = celestialsMock.NewMockAPI(s.ctrl)
}

// TearDownSuite -
func (s *ModuleTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
	s.ctrl.Finish()
}

func (s *ModuleTestSuite) TestSync() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("timescaledb"),
		testfixtures.Directory("../../test/data"),
		testfixtures.UseAlterConstraint(),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	s.api.EXPECT().
		Changes(
			gomock.Any(),
			network,
			gomock.Any(),
		).
		Times(1).
		Return(celestials.Changes{
			Head: 1,
			Changes: []celestials.Change{
				{
					CelestialID: "test",
					Address:     "celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8",
					ImageURL:    "image_url",
					ChangeID:    4,
				},
			},
		}, nil)

	cfgDs := config.DataSource{
		Kind:              "celestials",
		URL:               "base_url",
		Timeout:           10,
		RequestsPerSecond: 10,
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	m := New(
		cfgDs,
		s.storage.Address,
		s.storage.CelestialState,
		s.storage.Transactable,
		testIndexerName,
		network,
		WithAddressPrefix(types.AddressPrefixCelestia),
		WithLimit(10),
	)
	m.celestials = s.api

	err = m.getState(ctx)
	s.Require().NoError(err)

	err = m.sync(ctx)
	s.Require().NoError(err)

	st, err := s.storage.CelestialState.ByName(ctx, testIndexerName)
	s.Require().NoError(err)
	s.Require().EqualValues(4, st.ChangeId)
	s.Require().EqualValues(testIndexerName, st.Name)

	item, err := s.storage.Celestials.ById(ctx, "test")
	s.Require().NoError(err)
	s.Require().EqualValues("image_url", item.ImageUrl)
	s.Require().EqualValues("test", item.Id)
	s.Require().EqualValues(4, item.ChangeId)
	s.Require().EqualValues(1, item.AddressId)
}

func TestSuiteModule_Run(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
