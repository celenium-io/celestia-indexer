package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"testing"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

// TransactionTestSuite -
type TransactionTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *TransactionTestSuite) SetupSuite() {
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
}

// TearDownSuite -
func (s *TransactionTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *StorageTestSuite) TestSaveNamespaces() {

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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	namespaceId := []byte{0x5F, 0x7A, 0x8D, 0xDF, 0xE6, 0x13, 0x6F, 0xE7, 0x6B, 0x65, 0xB9, 0x06, 0x6D, 0x4F, 0x81, 0x6D, 0x70, 0x7F}
	namespaces := []storage.Namespace{
		{
			Version:     0,
			NamespaceID: namespaceId,
			PfbCount:    2,
			Size:        100,
		}, {
			Version:     2,
			NamespaceID: namespaceId,
			PfbCount:    1,
			Size:        11,
		},
	}

	err = tx.SaveNamespaces(ctx, namespaces...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	s.Require().Greater(namespaces[0].Id, uint64(0))
	s.Require().Greater(namespaces[1].Id, uint64(0))

	ns1, err := s.storage.Namespace.ByNamespaceIdAndVersion(ctx, namespaceId, 0)
	s.Require().NoError(err)

	s.Require().EqualValues(1, ns1.Id)
	s.Require().EqualValues(0, ns1.Version)
	s.Require().EqualValues(5, ns1.PfbCount)
	s.Require().EqualValues(1334, ns1.Size)
	s.Require().Equal(namespaceId, ns1.NamespaceID)

	ns2, err := s.storage.Namespace.ByNamespaceIdAndVersion(ctx, namespaceId, 2)
	s.Require().NoError(err)

	s.Require().Greater(ns2.Id, uint64(0))
	s.Require().EqualValues(2, ns2.Version)
	s.Require().EqualValues(1, ns2.PfbCount)
	s.Require().EqualValues(11, ns2.Size)
	s.Require().Equal(namespaceId, ns2.NamespaceID)
}

func (s *StorageTestSuite) TestSaveAddresses() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	addresses := make([]storage.Address, 5)
	for i := 0; i < 5; i++ {
		addresses[i].Height = storage.Level(10000 + i)
		addresses[i].Balance = decimal.NewFromInt(int64(i * 100))
		hash := make([]byte, 32)
		_, err := rand.Read(hash)
		s.NoError(err)
		addresses[i].Hash = hash
	}

	err = tx.SaveAddresses(ctx, addresses...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	s.Require().Greater(addresses[0].Id, uint64(0))
	s.Require().Greater(addresses[1].Id, uint64(0))

	for i := 0; i < 5; i++ {
		address, err := s.storage.Address.GetByID(ctx, addresses[i].Id)
		s.Require().NoError(err)

		s.Require().EqualValues(10000+i, address.Height)
		s.Require().EqualValues(i*100, address.Balance.IntPart())
		s.Require().Len(address.Hash, 32)
	}
}

func (s *StorageTestSuite) TestRollbackBlock() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	oldBlock, err := tx.RollbackBlock(ctx, 1000)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, oldBlock.Height)

	newHead, err := tx.LastBlock(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(999, newHead.Height)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

}

func (s *StorageTestSuite) TestRollbackAddress() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	deleted, err := tx.RollbackAddresses(ctx, 101)
	s.Require().NoError(err)
	s.Require().Len(deleted, 1)
	s.Require().Equal("321", deleted[0].Balance.String())

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Address.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func (s *StorageTestSuite) TestRollbackTxs() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	deleted, err := tx.RollbackTxs(ctx, 1000)
	s.Require().NoError(err)
	s.Require().Len(deleted, 2)
	s.Require().EqualValues(80410, deleted[0].GasWanted)
	s.Require().EqualValues(80410, deleted[1].GasWanted)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Tx.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *StorageTestSuite) TestRollbackEvents() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	deleted, err := tx.RollbackEvents(ctx, 1000)
	s.Require().NoError(err)
	s.Require().Len(deleted, 3)
	s.Require().Equal(types.EventTypeBurn, deleted[0].Type)
	s.Require().Equal(types.EventTypeMint, deleted[1].Type)
	s.Require().Equal(types.EventTypeMint, deleted[2].Type)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Event.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *StorageTestSuite) TestRollbackMessages() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	deleted, err := tx.RollbackMessages(ctx, 1000)
	s.Require().NoError(err)
	s.Require().Len(deleted, 4)
	s.Require().Equal(types.MsgWithdrawDelegatorReward, deleted[0].Type)
	s.Require().Equal(types.MsgDelegate, deleted[1].Type)
	s.Require().Equal(types.MsgUnjail, deleted[2].Type)
	s.Require().Equal(types.MsgPayForBlobs, deleted[3].Type)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Message.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *StorageTestSuite) TestRollbackNamespaces() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	deleted, err := tx.RollbackNamespaces(ctx, 1000)
	s.Require().NoError(err)
	s.Require().Len(deleted, 3)
	s.Require().EqualValues(1234, deleted[0].Size)
	s.Require().EqualValues(1255, deleted[1].Size)
	s.Require().EqualValues(12, deleted[2].Size)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Namespace.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *StorageTestSuite) TestRollbackNamespaceMessages() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	deleted, err := tx.RollbackNamespaceMessages(ctx, 1000)
	s.Require().NoError(err)
	s.Require().Len(deleted, 1)
	s.Require().EqualValues(2, deleted[0].NamespaceId)

	ns, err := tx.Namespace(ctx, 2)
	s.Require().NoError(err)
	s.Require().EqualValues(2, ns.Id)

	state, err := tx.State(ctx, testIndexerName)
	s.Require().NoError(err)
	s.Require().EqualValues(1, state.Id)
	s.Require().EqualValues(1000, state.LastHeight)
	s.Require().EqualValues(394067, state.TotalTx)
	s.Require().EqualValues(12512357, state.TotalAccounts)
	s.Require().Equal("172635712635813", state.TotalFee.String())
	s.Require().EqualValues(324234, state.TotalBlobsSize)
	s.Require().Equal(testIndexerName, state.Name)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func TestSuiteTransaction_Run(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
