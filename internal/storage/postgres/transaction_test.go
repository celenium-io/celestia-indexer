// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
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
}

// TearDownSuite -
func (s *TransactionTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *TransactionTestSuite) BeforeTest(suiteName, testName string) {
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

func (s *TransactionTestSuite) TestSaveNamespaces() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	testTime := time.Now().UTC()
	existedNamespace := testsuite.MustHexDecode("62491A45621ABEA79EBA193FD2944B5B9EBD")
	namespaceId := []byte{0x5F, 0x7A, 0x8D, 0xDF, 0xE6, 0x13, 0x6F, 0xE7, 0x6B, 0x65, 0xB9, 0x06, 0x6D, 0x4F, 0x81, 0x6D, 0x70, 0x7F}
	namespaces := []*storage.Namespace{
		{
			Version:         0,
			NamespaceID:     namespaceId,
			PfbCount:        2,
			Size:            100,
			LastHeight:      1001,
			LastMessageTime: testTime,
		}, {
			Version:         2,
			NamespaceID:     namespaceId,
			PfbCount:        1,
			Size:            11,
			LastHeight:      1001,
			LastMessageTime: testTime,
		}, {
			Version:         0,
			NamespaceID:     existedNamespace,
			PfbCount:        1,
			Size:            12,
			LastHeight:      1001,
			LastMessageTime: testTime,
		},
	}

	countAddedNamespaces, err := tx.SaveNamespaces(ctx, namespaces...)
	s.Require().NoError(err)
	s.Require().Equal(int64(1), countAddedNamespaces)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	s.Require().Greater(namespaces[0].Id, uint64(0))
	s.Require().Greater(namespaces[1].Id, uint64(0))
	s.Require().Greater(namespaces[2].Id, uint64(0))

	ns1, err := s.storage.Namespace.ByNamespaceIdAndVersion(ctx, namespaceId, 0)
	s.Require().NoError(err)

	s.Require().EqualValues(1, ns1.Id)
	s.Require().EqualValues(0, ns1.Version)
	s.Require().EqualValues(5, ns1.PfbCount)
	s.Require().EqualValues(1334, ns1.Size)
	s.Require().EqualValues(1001, ns1.LastHeight)
	s.Require().Equal(testTime.Unix(), ns1.LastMessageTime.Unix())
	s.Require().Equal(namespaceId, ns1.NamespaceID)

	ns2, err := s.storage.Namespace.ByNamespaceIdAndVersion(ctx, namespaceId, 2)
	s.Require().NoError(err)

	s.Require().Greater(ns2.Id, uint64(0))
	s.Require().EqualValues(2, ns2.Version)
	s.Require().EqualValues(1, ns2.PfbCount)
	s.Require().EqualValues(11, ns2.Size)
	s.Require().EqualValues(1001, ns2.LastHeight)
	s.Require().Equal(testTime.Unix(), ns2.LastMessageTime.Unix())
	s.Require().Equal(namespaceId, ns2.NamespaceID)

	ns3, err := s.storage.Namespace.ByNamespaceIdAndVersion(ctx, existedNamespace, 0)
	s.Require().NoError(err)

	s.Require().EqualValues(ns3.Id, 3)
	s.Require().EqualValues(0, ns3.Version)
	s.Require().EqualValues(2, ns3.PfbCount)
	s.Require().EqualValues(24, ns3.Size)
	s.Require().EqualValues(1001, ns3.LastHeight)
	s.Require().Equal(testTime.Unix(), ns3.LastMessageTime.Unix())
	s.Require().Equal(existedNamespace, ns3.NamespaceID)
}

func (s *TransactionTestSuite) TestSaveAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	replyAddress := storage.Address{}
	addresses := make([]*storage.Address, 0, 5)
	for i := 0; i < 5; i++ {
		hash := make([]byte, 20)
		for j := 0; j < 19; j++ {
			hash[j] = byte(j)
		}
		hash[19] = byte(i)
		s.NoError(err)

		addr, err := bech32.ConvertAndEncode(pkgTypes.AddressPrefixCelestia, hash)
		s.NoError(err)

		addresses = append(addresses, &storage.Address{
			Height:     pkgTypes.Level(10000 + i),
			LastHeight: pkgTypes.Level(10000 + i),
			Hash:       hash,
			Address:    addr,
			Id:         uint64(i),
		})

		if i == 2 {
			replyAddress.Address = addresses[i].Address
			replyAddress.Hash = addresses[i].Hash
			replyAddress.Height = addresses[i].Height + 1
			replyAddress.LastHeight = addresses[i].Height + 1
		}
	}

	count1, err := tx.SaveAddresses(ctx, addresses...)
	s.Require().NoError(err)
	s.Require().EqualValues(5, count1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	s.Require().Greater(addresses[0].Id, uint64(0))
	s.Require().Greater(addresses[1].Id, uint64(0))

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	count2, err := tx2.SaveAddresses(ctx, &replyAddress)
	s.Require().NoError(err)
	s.Require().EqualValues(0, count2)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
	s.Require().Equal(replyAddress.Id, addresses[2].Id)
}

func (s *TransactionTestSuite) TestSaveTxAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	addresses := make([]storage.Signer, 5)
	for i := uint64(0); i < 5; i++ {
		addresses[i].AddressId = i + 1
		addresses[i].TxId = 5 - i
	}

	err = tx.SaveSigners(ctx, addresses...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveMsgAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	addresses := make([]storage.MsgAddress, 5)
	for i := 0; i < 5; i++ {
		addresses[i].AddressId = uint64(i + 1)
		addresses[i].MsgId = uint64(5 - i)
		addresses[i].Type = types.MsgAddressTypeValues()[i]
	}

	err = tx.SaveMsgAddresses(ctx, addresses...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveBalances() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	balances := make([]storage.Balance, 5)
	for i := 0; i < 5; i++ {
		balances[i].Id = uint64(i + 1)
		balances[i].Spendable = decimal.RequireFromString("1000")
		balances[i].Currency = "utia"
	}

	err = tx.SaveBalances(ctx, balances...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveNamespaceMessages() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	nsMsgs := make([]storage.NamespaceMessage, 5)
	for i := 0; i < 5; i++ {
		nsMsgs[i].MsgId = uint64(i + 1)
		nsMsgs[i].NamespaceId = uint64(5 - i)
		nsMsgs[i].TxId = uint64((i + 1) * 2)
	}

	err = tx.SaveNamespaceMessage(ctx, nsMsgs...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveBlobLogs() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	blobLogs := make([]storage.BlobLog, 5)
	for i := 0; i < 5; i++ {
		blobLogs[i].MsgId = uint64(i + 1)
		blobLogs[i].NamespaceId = uint64(5 - i)
		blobLogs[i].TxId = uint64((i + 1) * 2)
		blobLogs[i].Time = time.Now()
		blobLogs[i].Height = 1000
	}

	err = tx.SaveBlobLogs(ctx, blobLogs...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveProposals() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	count, err := tx.SaveProposals(ctx, &storage.Proposal{
		Id:          3,
		ProposerId:  2,
		Status:      types.ProposalStatusInactive,
		Type:        types.ProposalTypeText,
		CreatedAt:   time.Now(),
		VotingPower: decimal.Zero,
	}, &storage.Proposal{
		Id:          1,
		Status:      types.ProposalStatusActive,
		VotingPower: decimal.Zero,
	})
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Proposals.ListWithFilters(ctx, storage.ListProposalFilters{
		Limit: 10,
		Sort:  sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 3)
}

func (s *TransactionTestSuite) TestSaveVotes() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	var vote = &storage.Vote{
		Option:     types.VoteOptionAbstain,
		ProposalId: 1,
		VoterId:    1,
		Height:     1001,
		Time:       time.Now(),
		Weight:     decimal.RequireFromString("1.0"),
	}

	err = tx.SaveVotes(ctx, vote)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Votes.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 4)

	for i := range items {
		if items[i].VoterId == vote.VoterId && items[i].ProposalId == vote.ProposalId {
			s.Require().Equal(vote.Option, items[i].Option)
		}
	}
}

func (s *TransactionTestSuite) TestSaveBlockSignatures() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	bs := make([]storage.BlockSignature, 5)
	for i := 0; i < 5; i++ {
		bs[i].ValidatorId = uint64(i + 1)
		bs[i].Height = 10000
		bs[i].Time = time.Now()
	}

	err = tx.SaveBlockSignatures(ctx, bs...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackBlockSignatures() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackBlockSignatures(ctx, 7965)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackBlock(ctx, 1000)
	s.Require().NoError(err)

	newHead, err := tx.LastBlock(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(999, newHead.Height)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

}

func (s *TransactionTestSuite) TestRollbackBlockStats() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	stats, err := tx.RollbackBlockStats(ctx, 1000)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, stats.Height)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

}

func (s *TransactionTestSuite) TestRollbackAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	deleted, err := tx.RollbackAddresses(ctx, 101)
	s.Require().NoError(err)
	s.Require().Len(deleted, 2)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Address.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 2)
}

func (s *TransactionTestSuite) TestRollbackTxs() {
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
	s.Require().Len(items, 2)
}

func (s *TransactionTestSuite) TestRollbackEvents() {
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

func (s *TransactionTestSuite) TestRollbackMessages() {
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
	s.Require().Len(items, 1)
}

func (s *TransactionTestSuite) TestRollbackBlobLogs() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackBlobLog(ctx, 1000)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.BlobLogs.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func (s *TransactionTestSuite) TestRollbackValidators() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	validators, err := tx.RollbackValidators(ctx, 999)
	s.Require().NoError(err)
	s.Require().Len(validators, 1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Validator.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func (s *TransactionTestSuite) TestRollbackNamespaces() {
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

func (s *TransactionTestSuite) TestRollbackUndelegations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackUndelegations(ctx, 1000)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Undelegation.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *TransactionTestSuite) TestRollbackRedelegations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackRedelegations(ctx, 1000)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Redelegation.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *TransactionTestSuite) TestRollbackNamespaceMessages() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	deleted, err := tx.RollbackNamespaceMessages(ctx, 1000)
	s.Require().NoError(err)
	s.Require().Len(deleted, 2)
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

func (s *TransactionTestSuite) TestRollbackProposals() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackProposals(ctx, 1000)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Proposals.ListWithFilters(ctx, storage.ListProposalFilters{
		Limit: 10,
		Sort:  sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func (s *TransactionTestSuite) TestRollbackVotes() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackVotes(ctx, 1000)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	items, err := s.storage.Votes.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *TransactionTestSuite) TestDeleteBalances() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.DeleteBalances(ctx, []uint64{1})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestLastAddressAction() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	hash := testsuite.MustHexDecode("dece425b75d67115bda877e1e7a1f262f6fa51d6")

	height, err := tx.LastAddressAction(ctx, hash)
	s.Require().NoError(err)
	s.Require().EqualValues(1000, height)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveEvents() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	events := []storage.Event{
		{
			Height:   100,
			Position: 0,
			Type:     types.EventTypeBurn,
			TxId:     testsuite.Ptr(uint64(1)),
			Data: map[string]any{
				"address": "address",
				"value":   "value",
			},
		}, {
			Height:   100,
			Position: 1,
			Type:     types.EventTypeCoinSpent,
			TxId:     nil,
			Data: map[string]any{
				"address": "address",
				"value":   "value",
			},
		},
	}

	err = tx.SaveEvents(ctx, events...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	saved, err := s.storage.Event.List(ctx, 2, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(saved, 2)
}

func (s *TransactionTestSuite) TestSaveEventsWithCopy() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	events := make([]storage.Event, 100)
	for i := 0; i < 100; i++ {
		events[i].Height = 100
		events[i].Position = int64(i)
		events[i].Type = types.EventTypeBurn
		events[i].TxId = testsuite.Ptr(uint64(i))
		events[i].Data = map[string]any{
			"address": "address",
			"value":   "value",
		}
	}

	err = tx.SaveEvents(ctx, events...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	saved, err := s.storage.Event.List(ctx, 100, 0, sdk.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(saved, 100)

	for i := 0; i < 100; i++ {
		s.Require().EqualValues(100, saved[i].Height)
		s.Require().EqualValues(99-i, saved[i].Position)
		s.Require().EqualValues(types.EventTypeBurn, saved[i].Type)
		s.Require().NotNil(saved[i].TxId)
		s.Require().NotNil(saved[i].Data)
		s.Require().Len(saved[i].Data, 2)
	}
}

func (s *TransactionTestSuite) TestGetProposerId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	id, err := tx.GetProposerId(ctx, "81A24EE534DEFE1557A4C7C437E8E8FBC2F834E8")
	s.Require().NoError(err)
	s.Require().EqualValues(1, id)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

const testLink = "test_link"

func (s *TransactionTestSuite) TestSaveUpdateAndDeleteRollup() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	rollup := &storage.Rollup{
		Name:           "Rollup 2",
		Description:    "The second",
		Website:        "https://website.com",
		Twitter:        "https://x.com/rollup2",
		L2Beat:         testLink,
		BridgeContract: testLink,
		Explorer:       testLink,
		Stack:          "stack",
		Type:           types.RollupTypeSettled,
		Category:       types.RollupCategoryFinance,
		SettledOn:      "Ethereum",
		Tags:           []string{"tag"},
		Links:          []string{testLink},
		Color:          "#333",
	}
	err = tx.SaveRollup(ctx, rollup)
	s.Require().NoError(err)
	s.Require().Greater(rollup.Id, uint64(0))

	rollup.GitHub = "https://github.com/rollup2"
	rollup.DeFiLama = "test"
	err = tx.UpdateRollup(ctx, rollup)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	newRollup, err := s.storage.Rollup.GetByID(ctx, rollup.Id)
	s.Require().NoError(err)

	s.Require().EqualValues(rollup.Name, newRollup.Name)
	s.Require().EqualValues(rollup.Description, newRollup.Description)
	s.Require().EqualValues(rollup.Website, newRollup.Website)
	s.Require().EqualValues(rollup.GitHub, newRollup.GitHub)
	s.Require().EqualValues(rollup.Twitter, newRollup.Twitter)
	s.Require().EqualValues(testLink, newRollup.L2Beat)
	s.Require().EqualValues(testLink, newRollup.Explorer)
	s.Require().EqualValues(testLink, newRollup.BridgeContract)
	s.Require().EqualValues("test", newRollup.DeFiLama)
	s.Require().EqualValues("stack", newRollup.Stack)
	s.Require().EqualValues("finance", newRollup.Category)
	s.Require().EqualValues("Ethereum", newRollup.SettledOn)
	s.Require().EqualValues("#333", newRollup.Color)
	s.Require().Len(newRollup.Links, 1)
	s.Require().Len(newRollup.Tags, 1)

	tx, err = BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.DeleteRollup(ctx, newRollup.Id)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRetentionBlockSignatures() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RetentionBlockSignatures(ctx, 999)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	signs, err := s.storage.BlockSignatures.List(ctx, 10, 0, sdk.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(signs, 1)
}

func (s *TransactionTestSuite) TestSaveRedelegations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	redelegations := []storage.Redelegation{
		{
			Height:         1000,
			Time:           time.Now(),
			SrcId:          2,
			DestId:         3,
			AddressId:      1,
			Amount:         decimal.NewFromInt(10),
			CompletionTime: time.Now().Add(time.Hour),
		},
	}

	err = tx.SaveRedelegations(ctx, redelegations...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	saved, err := s.storage.Redelegation.List(ctx, 2, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(saved, 2)
}

func (s *TransactionTestSuite) TestSaveUndelegations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	undelegations := []storage.Undelegation{
		{
			Height:         1000,
			Time:           time.Now(),
			ValidatorId:    2,
			AddressId:      1,
			Amount:         decimal.NewFromInt(10),
			CompletionTime: time.Now().Add(time.Hour),
		},
	}

	err = tx.SaveUndelegations(ctx, undelegations...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	saved, err := s.storage.Undelegation.List(ctx, 2, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(saved, 2)
}

func (s *TransactionTestSuite) TestSaveDelegations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	delegations := []storage.Delegation{
		{
			ValidatorId: 2,
			AddressId:   1,
			Amount:      decimal.NewFromInt(10),
		}, {
			ValidatorId: 1,
			AddressId:   1,
			Amount:      decimal.NewFromInt(10),
		},
	}

	err = tx.SaveDelegations(ctx, delegations...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	saved, err := s.storage.Delegation.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(saved, 3)
}

func (s *TransactionTestSuite) TestRetentionCompletedUnbondings() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	blockTime, err := time.Parse(time.RFC3339, "2023-07-04T03:11:57+00:00")
	s.Require().NoError(err)

	err = tx.RetentionCompletedUnbondings(ctx, blockTime)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRetentionCompletedRedelegations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	blockTime, err := time.Parse(time.RFC3339, "2023-07-04T03:11:57+00:00")
	s.Require().NoError(err)

	err = tx.RetentionCompletedRedelegations(ctx, blockTime)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestJail() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.Jail(ctx, &storage.Validator{
		Id:    2,
		Stake: decimal.NewFromInt(-10),
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	val, err := s.storage.Validator.ByAddress(ctx, "celestiavaloper189ecvq5avj0wehrcfnagpd5sd8pup9aqmdglmr")
	s.Require().NoError(err)
	s.Require().NotNil(val.Jailed)
	s.Require().True(*val.Jailed)
	s.Require().Equal("1000090", val.Stake.String())
}

func (s *TransactionTestSuite) TestUpdateSlashedDelegations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	balances, err := tx.UpdateSlashedDelegations(ctx, 1, decimal.NewFromFloat(0.01))
	s.Require().NoError(err)
	s.Require().Len(balances, 2)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	s.Require().Equal("-100", balances[0].Delegated.String())
	s.Require().Equal("utia", balances[0].Currency)
	s.Require().EqualValues(1, balances[0].Id)

	s.Require().Equal("-100", balances[1].Delegated.String())
	s.Require().Equal("utia", balances[1].Currency)
	s.Require().EqualValues(2, balances[1].Id)
}

func (s *TransactionTestSuite) TestSaveValidators() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	validators := []*storage.Validator{
		{
			Address:           "celestiavaloper1cj45qyagkujxgdgv8lgjur56zjxtrsy40g3pw3",
			Delegator:         "celestia1cj45qyagkujxgdgv8lgjur56zjxtrsy42hncch",
			ConsAddress:       "95764047BDFFB5CCADFA635DC63365EEB65F00C2",
			Rate:              decimal.NewFromFloat32(0.2),
			MaxRate:           decimal.NewFromFloat32(0.2),
			MaxChangeRate:     decimal.Zero,
			MinSelfDelegation: decimal.Zero,
			Identity:          "0068ECE5E6EB5359",
			Stake:             decimal.NewFromInt(100),
			Moniker:           "Polychain",
			Commissions:       decimal.Zero,
			Rewards:           decimal.Zero,
			Height:            1001,
			Jailed:            testsuite.Ptr(false),
		}, {
			Address:     "celestiavaloper189ecvq5avj0wehrcfnagpd5sd8pup9aqmdglmr",
			Delegator:   "celestia189ecvq5avj0wehrcfnagpd5sd8pup9aq7j2xd9",
			Rate:        decimal.NewFromFloat32(0.06),
			Commissions: decimal.NewFromInt(100),
			Rewards:     decimal.NewFromInt(200),
			Website:     "test-website",
			Identity:    storage.DoNotModify,
			Contacts:    storage.DoNotModify,
			Moniker:     storage.DoNotModify,
			Details:     storage.DoNotModify,
			Stake:       decimal.Zero,
			Height:      1001,
			Jailed:      testsuite.Ptr(true),
		}, {
			Address:     "celestiavaloper17vmk8m246t648hpmde2q7kp4ft9uwrayy09dmw",
			Delegator:   "celestia17vmk8m246t648hpmde2q7kp4ft9uwrayps85dg",
			Commissions: decimal.NewFromInt(100),
			Rewards:     decimal.NewFromInt(200),
			Stake:       decimal.Zero,
			Website:     storage.DoNotModify,
			Identity:    storage.DoNotModify,
			Contacts:    storage.DoNotModify,
			Moniker:     storage.DoNotModify,
			Details:     storage.DoNotModify,
			Height:      1001,
		},
	}

	count, err := tx.SaveValidators(ctx, validators...)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	entities, err := s.storage.Validator.ListByPower(ctx, storage.ValidatorFilters{
		Limit: 10,
	})
	s.Require().NoError(err)
	s.Require().Len(entities, 3)

	first := entities[0]
	s.Require().EqualValues("101", first.Commissions.String())
	s.Require().EqualValues("201", first.Rewards.String())
	s.Require().NotNil(first.Jailed)
	s.Require().False(*first.Jailed)

	second := entities[1]
	s.Require().EqualValues("100", second.Stake.String())
	s.Require().EqualValues("0.2", second.Rate.String())
	s.Require().EqualValues("Polychain", second.Moniker)
	s.Require().NotNil(second.Jailed)
	s.Require().False(*second.Jailed)

	third := entities[2]
	s.Require().EqualValues("101", third.Commissions.String())
	s.Require().EqualValues("201", third.Rewards.String())
	s.Require().EqualValues("0.06", third.Rate.String())
	s.Require().EqualValues("test-website", third.Website)
	s.Require().NotNil(third.Jailed)
	s.Require().True(*third.Jailed)

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	validatorJailed := []*storage.Validator{
		{
			Address:     "celestiavaloper189ecvq5avj0wehrcfnagpd5sd8pup9aqmdglmr",
			Delegator:   "celestia189ecvq5avj0wehrcfnagpd5sd8pup9aq7j2xd9",
			Commissions: decimal.NewFromInt(100),
			Rewards:     decimal.NewFromInt(200),
			Identity:    storage.DoNotModify,
			Website:     storage.DoNotModify,
			Contacts:    storage.DoNotModify,
			Moniker:     storage.DoNotModify,
			Details:     storage.DoNotModify,
			Stake:       decimal.Zero,
		},
	}

	count2, err := tx2.SaveValidators(ctx, validatorJailed...)
	s.Require().NoError(err)
	s.Require().EqualValues(0, count2)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))

	v, err := s.storage.Validator.ByAddress(ctx, "celestiavaloper189ecvq5avj0wehrcfnagpd5sd8pup9aqmdglmr")
	s.Require().NoError(err)
	s.Require().NotNil(v.Jailed)
	s.Require().True(*v.Jailed)
}

func (s *TransactionTestSuite) TestValidators() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	validators, err := tx.BondedValidators(ctx, 100)
	s.Require().NoError(err)

	s.Require().Len(validators, 2)

	for i := range validators {
		s.Require().NotEmpty(validators[i].Id)
		s.Require().NotEmpty(validators[i].Stake.IntPart())
		s.Require().Empty(validators[i].Address)
	}

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestProposalVotes() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	votes, err := tx.ProposalVotes(ctx, 1, 1, 0)
	s.Require().NoError(err)

	s.Require().Len(votes, 1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestAddressDelegations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	delegations, err := tx.AddressDelegations(ctx, 1)
	s.Require().NoError(err)

	s.Require().Len(delegations, 1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestUpdateConstants() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.SaveConstants(ctx, storage.Constant{
		Module: "auth",
		Name:   "tx_size_cost_per_byte",
		Value:  "20",
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	val, err := s.storage.Constants.Get(ctx, "auth", "tx_size_cost_per_byte")
	s.Require().NoError(err)
	s.Require().Equal("20", val.Value)
}

func (s *TransactionTestSuite) TestProposal() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	proposal, err := tx.Proposal(ctx, 1)
	s.Require().NoError(err)

	s.Require().EqualValues(1, proposal.Id)
	s.Require().NotNil(proposal.Changes)
	s.Require().EqualValues(types.ProposalTypeText, proposal.Type)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestActiveProposal() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	proposals, err := tx.ActiveProposals(ctx)
	s.Require().NoError(err)
	s.Require().Len(proposals, 1)

	s.Require().EqualValues(2, proposals[0].Id)
	s.Require().NotNil(proposals[0].Changes)
	s.Require().EqualValues(types.ProposalTypeCommunityPoolSpend, proposals[0].Type)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestIbcClients() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	count, err := tx.SaveIbcClients(ctx, &storage.IbcClient{
		Id:              "client-100",
		Height:          10000,
		CreatedAt:       time.Now().UTC(),
		TrustingPeriod:  time.Second,
		ConnectionCount: 1,
	}, &storage.IbcClient{
		Id:              "client-1",
		ConnectionCount: 10,
	})
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx2.RollbackIbcClients(ctx, 10000)
	s.Require().NoError(err)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
}

func (s *TransactionTestSuite) TestIbcConnections() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.SaveIbcConnections(ctx, &storage.IbcConnection{
		ConnectionId:             "connection-100",
		Height:                   10000,
		CreatedAt:                time.Now().UTC(),
		ClientId:                 "client-1",
		CounterpartyConnectionId: "test-100",
		CounterpartyClientId:     "test-client-100",
		ChannelsCount:            1,
	}, &storage.IbcConnection{
		ConnectionId:  "connection-1",
		ChannelsCount: 10,
		ConnectedAt:   time.Now().UTC(),
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx2.RollbackIbcClients(ctx, 10000)
	s.Require().NoError(err)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
}

func (s *TransactionTestSuite) TestIbcChannels() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.SaveIbcChannels(ctx, &storage.IbcChannel{
		Id:                    "channel-199",
		ConnectionId:          "connection-1",
		Height:                10000,
		CreatedAt:             time.Now().UTC(),
		ClientId:              "client-1",
		CounterpartyPortId:    "transfer",
		CounterpartyChannelId: "channel-100",
		Ordering:              true,
		Status:                types.IbcChannelStatusClosed,
		Sent:                  decimal.RequireFromString("100"),
		TransfersCount:        10,
	}, &storage.IbcChannel{
		Id:                 "channel-1",
		ConfirmedAt:        time.Now().UTC(),
		ConfirmationHeight: 10000,
		ConfirmationTxId:   100,
		Status:             types.IbcChannelStatusInitialization,
		Sent:               decimal.RequireFromString("100"),
		Received:           decimal.RequireFromString("100"),
		TransfersCount:     12,
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx2.RollbackIbcClients(ctx, 10000)
	s.Require().NoError(err)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
}

func (s *TransactionTestSuite) TestIbcTransfers() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.SaveIbcTransfers(ctx, &storage.IbcTransfer{
		Id:              123,
		ConnectionId:    "connection-1",
		ChannelId:       "channel-1",
		Height:          10000,
		Time:            time.Now().UTC(),
		Amount:          decimal.RequireFromString("12123123"),
		Denom:           currency.Utia,
		SenderId:        testsuite.Ptr(uint64(1)),
		ReceiverAddress: testsuite.Ptr("osmo1m8wg4vxkefhs374qxmmqpyusgz289wmulex5qdwpfx7jnrxzer5s9cv83q"),
		TxId:            3,
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx2.RollbackIbcTransfers(ctx, 10000)
	s.Require().NoError(err)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
}

func (s *TransactionTestSuite) TestIbcConnection() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	conn, err := tx.IbcConnection(ctx, "connection-1")
	s.Require().NoError(err)
	s.Require().EqualValues("client-1", conn.ClientId)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestHyperlaneTransfers() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.SaveHyperlaneTransfers(ctx, &storage.HLTransfer{
		Id:                  123,
		Height:              10000,
		Time:                time.Now().UTC(),
		Amount:              decimal.RequireFromString("12123123"),
		Denom:               currency.Utia,
		TxId:                2,
		MailboxId:           1,
		TokenId:             1,
		AddressId:           1,
		Counterparty:        1,
		CounterpartyAddress: "1234567890",
		Version:             2,
		Nonce:               2,
		Type:                types.HLTransferTypeReceive,
		RelayerId:           2,
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx2.RollbackHyperlaneTransfers(ctx, 10000)
	s.Require().NoError(err)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
}

func (s *TransactionTestSuite) TestHyperlaneTokens() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.SaveHyperlaneTokens(ctx, &storage.HLToken{
		Id:               123,
		Height:           10000,
		Time:             time.Now().UTC(),
		Denom:            currency.Utia,
		TxId:             2,
		MailboxId:        1,
		TokenId:          []byte{0, 1, 2, 3, 4, 5, 6},
		Type:             types.HLTokenTypeCollateral,
		SentTransfers:    1,
		ReceiveTransfers: 1,
		Sent:             decimal.RequireFromString("123"),
		Received:         decimal.Zero,
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx2.RollbackHyperlaneTokens(ctx, 10000)
	s.Require().NoError(err)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
}

func (s *TransactionTestSuite) TestHyperlaneMailbox() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.SaveHyperlaneMailbox(ctx, &storage.HLMailbox{
		Id:               123,
		Height:           10000,
		Time:             time.Now().UTC(),
		TxId:             2,
		Mailbox:          []byte{1, 2, 3, 4, 5, 6},
		OwnerId:          1,
		Domain:           2,
		SentMessages:     3,
		RequiredHook:     []byte("hook"),
		ReceivedMessages: 44,
		DefaultIsm:       []byte("ism"),
		DefaultHook:      []byte("default_hook"),
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx2.RollbackHyperlaneMailbox(ctx, 10000)
	s.Require().NoError(err)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
}

func (s *TransactionTestSuite) TestGetHyperlaneMailbox() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	mailbox, err := tx.HyperlaneMailbox(ctx, []byte("mailbox"))
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
	s.Require().EqualValues(1, mailbox.Id)
}

func (s *TransactionTestSuite) TestGetHyperlaneToken() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	token, err := tx.HyperlaneToken(ctx, []byte("token"))
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
	s.Require().EqualValues(1, token.Id)
}

func TestSuiteTransaction_Run(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
