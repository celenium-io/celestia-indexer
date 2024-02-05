// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testIndexerName = "test_indexer"
	testState       = storage.State{
		Name:            testIndexerName,
		LastHeight:      80000,
		LastTime:        testTime,
		TotalTx:         14149240,
		TotalAccounts:   123123,
		TotalNamespaces: 123,
		TotalBlobsSize:  1000,
	}
	testAddress     = "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"
	testHashAddress = []byte{0x96, 0xa, 0xa0, 0x36, 0x6b, 0x25, 0x4e, 0x1e, 0xa7, 0x9b, 0xda, 0x46, 0x7e, 0xb3, 0xaa, 0x5c, 0x97, 0xcb, 0xa5, 0xae}
)

// AddressTestSuite -
type AddressTestSuite struct {
	suite.Suite
	address  *mock.MockIAddress
	txs      *mock.MockITx
	blobLogs *mock.MockIBlobLog
	messages *mock.MockIMessage
	state    *mock.MockIState
	echo     *echo.Echo
	handler  *AddressHandler
	ctrl     *gomock.Controller
}

// SetupSuite -
func (s *AddressTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.address = mock.NewMockIAddress(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.blobLogs = mock.NewMockIBlobLog(s.ctrl)
	s.messages = mock.NewMockIMessage(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewAddressHandler(s.address, s.txs, s.blobLogs, s.messages, s.state, testIndexerName)
}

// TearDownSuite -
func (s *AddressTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteAddress_Run(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}

func (s *AddressTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		ByHash(gomock.Any(), testHashAddress).
		Return(storage.Address{
			Id:         1,
			Hash:       testHashAddress,
			Address:    testAddress,
			Height:     100,
			LastHeight: 100,
		}, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var address responses.Address
	err := json.NewDecoder(rec.Body).Decode(&address)
	s.Require().NoError(err)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(100, address.Height)
	s.Require().EqualValues(100, address.LastHeight)
	s.Require().Equal(testAddress, address.Hash)
}

func (s *AddressTestSuite) TestGetInvalidAddress() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash")
	c.SetParamNames("hash")
	c.SetParamValues("invalid")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *AddressTestSuite) TestGetBadAddress() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash")
	c.SetParamNames("hash")
	c.SetParamValues("celestia111111111111111111111111111111111111111")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *AddressTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address")

	s.address.EXPECT().
		ListWithBalance(gomock.Any(), storage.AddressListFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("asc"),
		}).
		Return([]storage.Address{
			{
				Id:         1,
				Hash:       testHashAddress,
				Address:    testAddress,
				Height:     100,
				LastHeight: 100,
				Balance: storage.Balance{
					Currency: "utia",
					Total:    decimal.RequireFromString("100"),
				},
			},
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var address []responses.Address
	err := json.NewDecoder(rec.Body).Decode(&address)
	s.Require().NoError(err)
	s.Require().Len(address, 1)
	s.Require().EqualValues(1, address[0].Id)
	s.Require().EqualValues(100, address[0].Height)
	s.Require().EqualValues(100, address[0].LastHeight)
	s.Require().Equal(testAddress, address[0].Hash)
	s.Require().Equal("100", address[0].Balance.Value)
	s.Require().Equal("utia", address[0].Balance.Currency)
}

func (s *AddressTestSuite) TestTransactions() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("msg_type", "MsgSend")
	q.Set("height", "1000")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/txs")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		ByHash(gomock.Any(), testHashAddress).
		Return(storage.Address{
			Id:         1,
			Hash:       testHashAddress,
			Address:    testAddress,
			Height:     100,
			LastHeight: 100,
		}, nil)

	s.txs.EXPECT().
		ByAddress(gomock.Any(), uint64(1), gomock.Any()).
		Return([]storage.Tx{
			testTx,
		}, nil)

	s.Require().NoError(s.handler.Transactions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&txs)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, strings.ToUpper(tx.Hash))
	s.Require().EqualValues(2, tx.Position)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().Equal("80410", tx.Fee)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(10, tx.EventsCount)
	s.Require().EqualValues(2, tx.MessagesCount)
	s.Require().Equal("memo", tx.Memo)
	s.Require().Equal("sdk", tx.Codespace)
	s.Require().Equal(types.StatusSuccess, tx.Status)
}

func (s *AddressTestSuite) TestMessages() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("sort", "desc")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/messages")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		ByHash(gomock.Any(), testHashAddress).
		Return(storage.Address{
			Id:      1,
			Hash:    testHashAddress,
			Address: testAddress,
		}, nil)

	s.messages.EXPECT().
		ByAddress(gomock.Any(), uint64(1), gomock.Any()).
		Return([]storage.AddressMessageWithTx{
			{
				MsgAddress: storage.MsgAddress{
					AddressId: 1,
					MsgId:     1,
					Type:      types.MsgAddressTypeDelegator,
				},
				Msg: &storage.Message{
					Id:       1,
					Height:   1000,
					Position: 0,
					Type:     types.MsgWithdrawDelegatorReward,
					TxId:     1,
					Data:     nil,
				},
				Tx: &storage.Tx{
					Id:            1,
					MessageTypes:  types.NewMsgTypeBitMask(types.MsgWithdrawDelegatorReward),
					MessagesCount: 1,
					Status:        types.StatusSuccess,
				},
			},
		}, nil)

	s.Require().NoError(s.handler.Messages(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var msgs []responses.Message
	err := json.NewDecoder(rec.Body).Decode(&msgs)
	s.Require().NoError(err)
	s.Require().Len(msgs, 1)

	msg := msgs[0]
	s.Require().EqualValues(1, msg.Id)
	s.Require().EqualValues(1000, msg.Height)
	s.Require().Equal(int64(0), msg.Position)
	s.Require().EqualValues(types.MsgWithdrawDelegatorReward, msg.Type)
	s.Require().NotNil(msg.Tx)
}

func (s *AddressTestSuite) TestBlobs() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("sort", "desc")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/blobs")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		ByHash(gomock.Any(), testHashAddress).
		Return(storage.Address{
			Id:      1,
			Hash:    testHashAddress,
			Address: testAddress,
		}, nil)

	s.blobLogs.EXPECT().
		BySigner(gomock.Any(), uint64(1), storage.BlobLogFilters{
			Limit: 10,
			Sort:  "desc",
		}).
		Return([]storage.BlobLog{
			{
				NamespaceId: testNamespace.Id,
				MsgId:       1,
				TxId:        1,
				SignerId:    1,
				Signer: &storage.Address{
					Address: testAddress,
				},
				Commitment: "test_commitment",
				Size:       1000,
				Height:     10000,
				Time:       testTime,
			},
		}, nil)

	s.Require().NoError(s.handler.Blobs(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var logs []responses.BlobLog
	err := json.NewDecoder(rec.Body).Decode(&logs)
	s.Require().NoError(err)
	s.Require().Len(logs, 1)

	l := logs[0]
	s.Require().EqualValues(10000, l.Height)
	s.Require().Equal(testTime, l.Time)
	s.Require().Equal(testAddress, l.Signer)
	s.Require().Equal("test_commitment", l.Commitment)
	s.Require().EqualValues(1000, l.Size)
	s.Require().Nil(l.Namespace)
}

func (s *AddressTestSuite) TestCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/count")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(testState, nil)

	s.Require().NoError(s.handler.Count(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count uint64
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(123123, count)
}
