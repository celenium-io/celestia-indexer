// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package genesis

import (
	"context"
	"os"
	"testing"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	dipdupConfig "github.com/dipdup-io/go-lib/config"
	"github.com/dipdup-io/go-lib/testhelpers"
	"github.com/stretchr/testify/suite"
)

// SaveTestSuite spins up a real TimescaleDB container, since Module.save opens a
// genuine DB transaction (postgres.BeginTransaction) and can't be exercised against
// a zero-value postgres.Storage.
type SaveTestSuite struct {
	suite.Suite
	psqlContainer *testhelpers.PostgreSQLContainer
	storage       postgres.Storage
}

func (s *SaveTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer ctxCancel()

	psqlContainer, err := testhelpers.NewPostgreSQLContainer(ctx, testhelpers.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb-ha:pg15.8-ts2.17.0-all",
	})
	s.Require().NoError(err)
	s.psqlContainer = psqlContainer
	s.T().Cleanup(func() {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer ctxCancel()
		s.Require().NoError(s.psqlContainer.Terminate(ctx))
	})

	strg, err := postgres.Create(ctx, dipdupConfig.Database{
		Kind:     dipdupConfig.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}, "../../../database", false)
	s.Require().NoError(err)
	s.storage = strg
	s.T().Cleanup(func() {
		s.Require().NoError(s.storage.Close())
	})
}

// parsedGenesisContext builds a decode context by running the real parser over the
// genesis.json fixture, so it is a known-good, DB-writable context before any test
// deliberately corrupts a part of it.
func (s *SaveTestSuite) parsedGenesisContext(module Module) *decodeContext.Context {
	f, err := os.Open("../../../test/json/genesis.json")
	s.Require().NoError(err)
	defer f.Close()

	var g types.Genesis
	s.Require().NoError(json.ConfigFastest.NewDecoder(f).Decode(&g))

	dCtx, err := module.parse(types.GenesisOutput{Genesis: g})
	s.Require().NoError(err)
	return dCtx
}

// Regression guard: a nil entry (or one with a nil Address) in decodeCtx.VestingAccounts
// used to reach decodeCtx.VestingAccounts[i].Address.Address directly and panic with a
// nil pointer dereference instead of failing the block gracefully.
func (s *SaveTestSuite) TestSave_NilVestingAccountReturnsErrorInsteadOfPanicking() {
	module := NewModule(s.storage, config.Indexer{Name: "test_indexer"})
	dCtx := s.parsedGenesisContext(module)

	dCtx.VestingAccounts = append(dCtx.VestingAccounts, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.Require().NotPanics(func() {
		err := module.save(ctx, dCtx)
		s.Require().Error(err)
	})
}

func TestSuiteSave_Run(t *testing.T) {
	suite.Run(t, new(SaveTestSuite))
}
