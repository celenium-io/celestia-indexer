package postgres

import (
	"context"
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
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

func (s *StorageTestSuite) TestBlockLast() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.Last(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().EqualValues(0, block.Stats.TxCount)

	hash, err := hex.DecodeString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")
	s.Require().NoError(err)
	s.Require().Equal(hash, block.Hash.Bytes())
}

func (s *StorageTestSuite) TestBlockByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.ByHeight(ctx, 1000)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().Equal(storage.BlockStats{}, block.Stats)

	hash, err := hex.DecodeString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")
	s.Require().NoError(err)
	s.Require().Equal(hash, block.Hash.Bytes())
}

func (s *StorageTestSuite) TestBlockByHeightWithStats() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.ByHeightWithStats(ctx, 1000)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)

	loc := &time.Location{}
	expectedStats := storage.BlockStats{
		Id:            2,
		Height:        1000,
		Time:          time.Date(2023, 07, 04, 03, 10, 57, 0, loc).UTC(),
		TxCount:       0,
		EventsCount:   0,
		BlobsSize:     0,
		BlockTime:     11000,
		SupplyChange:  decimal.NewFromInt(30930476),
		InflationRate: decimal.NewFromFloat(0.08),
		Fee:           decimal.NewFromInt(2873468273),
		MessagesCounts: map[types.MsgType]int64{
			types.MsgDelegate:                1,
			types.MsgPayForBlobs:             1,
			types.MsgUnjail:                  1,
			types.MsgWithdrawDelegatorReward: 1,
		},
	}
	s.Require().Equal(expectedStats, block.Stats)

	hash, err := hex.DecodeString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")
	s.Require().NoError(err)
	s.Require().Equal(hash, block.Hash.Bytes())
}

func (s *StorageTestSuite) TestBlockByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("6A30C94091DA7C436D64E62111D6890D772E351823C41496B4E52F28F5B000BF")
	s.Require().NoError(err)

	block, err := s.storage.Blocks.ByHash(ctx, hash)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().EqualValues(0, block.Stats.TxCount)
	s.Require().Equal(hash, block.Hash.Bytes())
}

func (s *StorageTestSuite) TestBlockListWithStats() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	blocks, err := s.storage.Blocks.ListWithStats(ctx, 10, 0, sdk.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(blocks, 2)

	block := blocks[0]
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().EqualValues(0, block.Stats.TxCount)
	s.Require().EqualValues(11000, block.Stats.BlockTime)
	s.Require().EqualValues(map[types.MsgType]int64{
		types.MsgWithdrawDelegatorReward: 1,
		types.MsgDelegate:                1,
		types.MsgUnjail:                  1,
		types.MsgPayForBlobs:             1,
	}, block.Stats.MessagesCounts)

	blocks, err = s.storage.Blocks.List(ctx, 10, 0, sdk.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(blocks, 2)

	block = blocks[0]
	s.Require().EqualValues(1000, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().EqualValues(storage.BlockStats{}, block.Stats)
}

func (s *StorageTestSuite) TestAddressByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash := []byte{0xde, 0xce, 0x42, 0x5b, 0x75, 0xd6, 0x71, 0x15, 0xbd, 0xa8, 0x77, 0xe1, 0xe7, 0xa1, 0xf2, 0x62, 0xf6, 0xfa, 0x51, 0xd6}
	address, err := s.storage.Address.ByHash(ctx, hash)
	s.Require().NoError(err)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(100, address.Height)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", address.Address)
}

func (s *StorageTestSuite) TestAddressList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	addresses, err := s.storage.Address.ListWithBalance(ctx, storage.AddressListFilter{
		Limit:  10,
		Offset: 0,
		Sort:   sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(addresses, 2)

	s.Require().EqualValues(1, addresses[0].Id)
	s.Require().EqualValues(100, addresses[0].Height)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", addresses[0].Address)
	s.Require().Equal("123", addresses[0].Balance.Total.String())
	s.Require().Equal("utia", addresses[0].Balance.Currency)

	s.Require().EqualValues(2, addresses[1].Id)
	s.Require().EqualValues(101, addresses[1].Height)
	s.Require().Equal("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", addresses[1].Address)
	s.Require().Equal("321", addresses[1].Balance.Total.String())
	s.Require().Equal("utia", addresses[1].Balance.Currency)
}

func (s *StorageTestSuite) TestEventByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	events, err := s.storage.Event.ByTxId(ctx, 1)
	s.Require().NoError(err)
	s.Require().Len(events, 1)
	s.Require().EqualValues(2, events[0].Id)
	s.Require().EqualValues(1000, events[0].Height)
	s.Require().EqualValues(1, events[0].Position)
	s.Require().Equal(types.EventTypeMint, events[0].Type)
}

func (s *StorageTestSuite) TestEventByBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	events, err := s.storage.Event.ByBlock(ctx, 1000)
	s.Require().NoError(err)
	s.Require().Len(events, 1)
	s.Require().EqualValues(1, events[0].Id)
	s.Require().EqualValues(1000, events[0].Height)
	s.Require().EqualValues(0, events[0].Position)
	s.Require().Equal(types.EventTypeBurn, events[0].Type)
}

func (s *StorageTestSuite) TestMessageByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	msgs, err := s.storage.Message.ByTxId(ctx, 1)
	s.Require().NoError(err)
	s.Require().Len(msgs, 2)
	s.Require().EqualValues(1, msgs[0].Id)
	s.Require().EqualValues(1000, msgs[0].Height)
	s.Require().EqualValues(0, msgs[0].Position)
	s.Require().Equal(types.MsgWithdrawDelegatorReward, msgs[0].Type)
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

func (s *StorageTestSuite) TestNamespaceMessagesByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	msgs, err := s.storage.Namespace.MessagesByHeight(ctx, 1000, 2, 0)
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
	s.Require().EqualValues(1255, msg.Namespace.Size)
}

func (s *StorageTestSuite) TestNamespaceCountMessagesByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	count, err := s.storage.Namespace.CountMessagesByHeight(ctx, 1000)
	s.Require().NoError(err)
	s.Require().EqualValues(count, 2)
}

func (s *StorageTestSuite) TestNamespaceActive() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	ns, err := s.storage.Namespace.Active(ctx, 2)
	s.Require().NoError(err)
	s.Require().Len(ns, 2)

	namespace := ns[0]
	s.Require().EqualValues(1000, namespace.Height)
	s.Require().EqualValues(2, namespace.Id)
	s.Require().NotNil(namespace.Namespace)
	s.Require().EqualValues(1255, namespace.Namespace.Size)
}

func (s *StorageTestSuite) TestTxByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)

	tx, err := s.storage.Tx.ByHash(ctx, txHash)
	s.Require().NoError(err)

	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(0, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(2, tx.MessagesCount)
	s.Require().Equal(txHash, tx.Hash)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo", tx.Memo)
	s.Require().Equal("sdk", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())
}

func (s *StorageTestSuite) TestTxFilterSuccessUnjailAsc() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Sort:         sdk.SortOrderAsc,
		Limit:        10,
		Offset:       0,
		MessageTypes: types.NewMsgTypeBitMask(types.MsgUnjail),
		Status:       []string{string(types.StatusSuccess)},
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]

	s.Require().EqualValues(2, tx.Id)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().EqualValues(256, tx.MessageTypes.Bits)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo2", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())
}

func (s *StorageTestSuite) TestTxFilterSuccessDesc() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Sort:   sdk.SortOrderDesc,
		Limit:  10,
		Offset: 0,
		Status: []string{string(types.StatusSuccess)},
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 4)

	tx := txs[1]

	s.Require().EqualValues(3, tx.Id)
	s.Require().EqualValues(0, tx.Position)
	s.Require().EqualValues(999, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(0, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().EqualValues(32, tx.MessageTypes.Bits)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())
}

func (s *StorageTestSuite) TestTxFilterHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Sort:   sdk.SortOrderDesc,
		Limit:  10,
		Offset: 0,
		Status: []string{string(types.StatusSuccess)},
		Height: 1000,
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 2)

	tx := txs[0]

	s.Require().EqualValues(2, tx.Id)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().EqualValues(256, tx.MessageTypes.Bits)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo2", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())
}

func (s *StorageTestSuite) TestTxFilterTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Limit:    10,
		TimeFrom: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 4)

	txs, err = s.storage.Tx.Filter(ctx, storage.TxFilter{
		Limit:  10,
		TimeTo: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 0)

	txs, err = s.storage.Tx.Filter(ctx, storage.TxFilter{
		Limit: 10,

		TimeFrom: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
		TimeTo:   time.Date(2023, 7, 5, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 4)
}

func (s *StorageTestSuite) TestTxByIdWithRelations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := s.storage.Tx.ByIdWithRelations(ctx, 2)
	s.Require().NoError(err)

	s.Require().EqualValues(2, tx.Id)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo2", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())
	s.Require().EqualValues(256, tx.MessageTypes.Bits)

	s.Require().Len(tx.Messages, 2)
}

func (s *StorageTestSuite) TestTxGenesis() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Genesis(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(4, tx.Id)
	s.Require().EqualValues(0, tx.Position)
	s.Require().EqualValues(0, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(0, tx.GasWanted)
	s.Require().EqualValues(0, tx.GasUsed)
	s.Require().EqualValues(0, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("34499b1ac473fbb03894c883178ecc83f0d6eaf6@64.227.18.169:26656", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("0", tx.Fee.String())
	s.Require().EqualValues(32, tx.MessageTypes.Bits)
}

func (s *StorageTestSuite) TestTxByAddressAndTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.ByAddress(ctx, 1, storage.TxFilter{
		Limit:    10,
		TimeFrom: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	txs, err = s.storage.Tx.ByAddress(ctx, 1, storage.TxFilter{
		Limit:  10,
		TimeTo: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 0)

	txs, err = s.storage.Tx.ByAddress(ctx, 1, storage.TxFilter{
		Limit: 10,

		TimeFrom: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
		TimeTo:   time.Date(2023, 7, 5, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 1)
}

func (s *StorageTestSuite) TestValidatorByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	validator, err := s.storage.Validator.ByAddress(ctx, "celestiavaloper17vmk8m246t648hpmde2q7kp4ft9uwrayy09dmw")
	s.Require().NoError(err)

	s.Require().Equal("celestiavaloper17vmk8m246t648hpmde2q7kp4ft9uwrayy09dmw", validator.Address)
	s.Require().Equal("celestia17vmk8m246t648hpmde2q7kp4ft9uwrayps85dg", validator.Delegator)
	s.Require().Equal("Conqueror", validator.Moniker)
	s.Require().Equal("https://github.com/DasRasyo", validator.Website)
	s.Require().Equal("EAD22B173DE57E6A", validator.Identity)
	s.Require().Equal("https://t.me/DasRasyo || conqueror.prime", validator.Contacts)
	s.Require().Equal("1", validator.MinSelfDelegation.String())
	s.Require().Equal("0.2", validator.MaxRate.String())
	s.Require().EqualValues(4, validator.MsgId)
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
