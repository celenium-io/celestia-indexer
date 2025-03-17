// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

const testIndexerName = "test_indexer"

// StorageTestSuite -
type StorageTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *StorageTestSuite) SetupSuite() {
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

	strg, err := Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}, "../../../database", false)
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
func (s *StorageTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *StorageTestSuite) TestStateGetByName() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	state, err := s.storage.State.ByName(ctx, testIndexerName)
	s.Require().NoError(err)
	s.Require().EqualValues(1, state.Id)
	s.Require().EqualValues(1000, state.LastHeight)
	s.Require().EqualValues(394067, state.TotalTx)
	s.Require().EqualValues(12512357, state.TotalAccounts)
	s.Require().Equal("172635712635813", state.TotalFee.String())
	s.Require().EqualValues(324234, state.TotalBlobsSize)
	s.Require().Equal(testIndexerName, state.Name)
}

func (s *StorageTestSuite) TestStateGetByNameFailed() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.State.ByName(ctx, "unknown")
	s.Require().Error(err)
}

func (s *StorageTestSuite) TestMessageByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	msgs, err := s.storage.Message.ByTxId(ctx, 1, 1, 0)
	s.Require().NoError(err)
	s.Require().Len(msgs, 1)
	s.Require().EqualValues(1, msgs[0].Id)
	s.Require().EqualValues(1000, msgs[0].Height)
	s.Require().EqualValues(0, msgs[0].Position)
	s.Require().Equal(types.MsgWithdrawDelegatorReward, msgs[0].Type)
}

func (s *StorageTestSuite) TestMessageListWithTx() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	msgs, err := s.storage.Message.ListWithTx(ctx, storage.MessageListWithTxFilters{
		Limit:                10,
		Offset:               0,
		Height:               1000,
		MessageTypes:         []string{types.MsgWithdrawDelegatorReward.String(), types.MsgUnjail.String()},
		ExcludedMessageTypes: []string{types.MsgUnjail.String()},
	})
	s.Require().NoError(err)
	s.Require().Len(msgs, 1)

	s.Require().EqualValues(1, msgs[0].Id)
	s.Require().EqualValues(1000, msgs[0].Height)
	s.Require().EqualValues(0, msgs[0].Position)
	s.Require().Equal(types.MsgWithdrawDelegatorReward, msgs[0].Type)
	s.Require().NotNil(msgs[0].Tx)

	tx := msgs[0].Tx

	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(0, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(2, tx.MessagesCount)

	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, tx.Hash)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo", tx.Memo)
	s.Require().Equal("sdk", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())
}

func (s *StorageTestSuite) TestNamespaceId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	namespaceId, err := hex.DecodeString("5F7A8DDFE6136FE76B65B9066D4F816D707F")
	s.Require().NoError(err)

	namespaces, err := s.storage.Namespace.ByNamespaceId(ctx, namespaceId)
	s.Require().NoError(err)
	s.Require().Len(namespaces, 2)

	s.Require().EqualValues(1, namespaces[0].Id)
	s.Require().EqualValues(0, namespaces[0].Version)
	s.Require().EqualValues(1234, namespaces[0].Size)
	s.Require().Equal(namespaceId, namespaces[0].NamespaceID)

	s.Require().EqualValues(2, namespaces[1].Id)
	s.Require().EqualValues(1, namespaces[1].Version)
	s.Require().EqualValues(1255, namespaces[1].Size)
	s.Require().Equal(namespaceId, namespaces[1].NamespaceID)
}

func (s *StorageTestSuite) TestNamespaceIdAndVersion() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	namespaceId, err := hex.DecodeString("5F7A8DDFE6136FE76B65B9066D4F816D707F")
	s.Require().NoError(err)

	namespace, err := s.storage.Namespace.ByNamespaceIdAndVersion(ctx, namespaceId, 1)
	s.Require().NoError(err)

	s.Require().EqualValues(2, namespace.Id)
	s.Require().EqualValues(1, namespace.Version)
	s.Require().EqualValues(1255, namespace.Size)
	s.Require().Equal(namespaceId, namespace.NamespaceID)
}

func (s *StorageTestSuite) TestNamespaceMessages() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	msgs, err := s.storage.Namespace.Messages(ctx, 2, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(msgs, 2)

	msg := msgs[0]
	s.Require().EqualValues(3, msg.MsgId)
	s.Require().EqualValues(2, msg.NamespaceId)
	s.Require().NotNil(msg.Namespace)
	s.Require().NotNil(msg.Message)
	s.Require().NotNil(msg.Tx)
	s.Require().Equal(types.MsgUnjail, msg.Message.Type)
	s.Require().EqualValues(2, msg.Tx.Id)
}

func (s *StorageTestSuite) TestNamespaceActive() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	ns, err := s.storage.Namespace.ListWithSort(ctx, "", sdk.SortOrderDesc, 2, 0)
	s.Require().NoError(err)
	s.Require().Len(ns, 2)

	namespace := ns[0]
	s.Require().EqualValues(1000, namespace.LastHeight)
	s.Require().EqualValues(3, namespace.Id)
	s.Require().EqualValues(12, namespace.Size)
}

func (s *StorageTestSuite) TestNamespaceActiveByPfbCount() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	ns, err := s.storage.Namespace.ListWithSort(ctx, pfbCountColumn, sdk.SortOrderDesc, 2, 0)
	s.Require().NoError(err)
	s.Require().Len(ns, 2)

	namespace := ns[0]
	s.Require().EqualValues(1000, namespace.LastHeight)
	s.Require().EqualValues(1, namespace.Id)
	s.Require().EqualValues(1234, namespace.Size)
}

func (s *StorageTestSuite) TestNamespaceActiveBySize() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	ns, err := s.storage.Namespace.ListWithSort(ctx, "size", sdk.SortOrderDesc, 2, 0)
	s.Require().NoError(err)
	s.Require().Len(ns, 2)

	namespace := ns[0]
	s.Require().EqualValues(1000, namespace.LastHeight)
	s.Require().EqualValues(2, namespace.Id)
	s.Require().EqualValues(1255, namespace.Size)
}

func (s *StorageTestSuite) TestNamespaceGetByIds() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	ns, err := s.storage.Namespace.GetByIds(ctx, 1, 2, 3)
	s.Require().NoError(err)
	s.Require().Len(ns, 3)
}

func (s *StorageTestSuite) TestConstantGet() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	c, err := s.storage.Constants.Get(ctx, types.ModuleNameBlob, "gas_per_blob_byte")
	s.Require().NoError(err)

	s.Require().EqualValues("8", c.Value)
}

func (s *StorageTestSuite) TestConstantByModule() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	c, err := s.storage.Constants.ByModule(ctx, types.ModuleNameAuth)
	s.Require().NoError(err)
	s.Require().Len(c, 2)

	s.Require().EqualValues("256", c[0].Value)
	s.Require().EqualValues("10", c[1].Value)
}

func (s *StorageTestSuite) TestDenomMetadata() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	metadata, err := s.storage.DenomMetadata.All(ctx)
	s.Require().NoError(err)
	s.Require().Len(metadata, 1)

	m := metadata[0]
	s.Require().EqualValues("utia", m.Base)
	s.Require().EqualValues("TIA", m.Display)
	s.Require().EqualValues("TIA", m.Symbol)
	s.Require().EqualValues("TIA", m.Name)
	s.Require().EqualValues("The native staking token of the Celestia network.", m.Description)
	s.Require().Greater(len(m.Units), 0)
}

func (s *StorageTestSuite) TestNotify() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.storage.Notificator.Subscribe(ctx, "test")
	s.Require().NoError(err)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var ticks int
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-s.storage.Notificator.Listen():
			log.Info().Str("msg", msg.Extra).Str("channel", msg.Channel).Msg("new message")
			s.Require().Equal("test", msg.Channel)
			s.Require().Equal("message", msg.Extra)
			if ticks == 2 {
				return
			}
		case <-ticker.C:
			ticks++
			err = s.storage.Notificator.Notify(ctx, "test", "message")
			s.Require().NoError(err)
		}
	}
}

func TestSuiteStorage_Run(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
