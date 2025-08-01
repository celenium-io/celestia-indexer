// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	celestialMock "github.com/celenium-io/celestial-module/pkg/storage/mock"
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
	address       *mock.MockIAddress
	txs           *mock.MockITx
	blobLogs      *mock.MockIBlobLog
	messages      *mock.MockIMessage
	delegations   *mock.MockIDelegation
	undelegations *mock.MockIUndelegation
	redelegations *mock.MockIRedelegation
	vestings      *mock.MockIVestingAccount
	grants        *mock.MockIGrant
	celestials    *celestialMock.MockICelestial
	votes         *mock.MockIVote
	state         *mock.MockIState
	echo          *echo.Echo
	handler       *AddressHandler
	ctrl          *gomock.Controller
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
	s.delegations = mock.NewMockIDelegation(s.ctrl)
	s.undelegations = mock.NewMockIUndelegation(s.ctrl)
	s.redelegations = mock.NewMockIRedelegation(s.ctrl)
	s.vestings = mock.NewMockIVestingAccount(s.ctrl)
	s.grants = mock.NewMockIGrant(s.ctrl)
	s.celestials = celestialMock.NewMockICelestial(s.ctrl)
	s.votes = mock.NewMockIVote(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewAddressHandler(s.address, s.txs, s.blobLogs, s.messages, s.delegations, s.undelegations, s.redelegations, s.vestings, s.grants, s.celestials, s.votes, s.state, testIndexerName)
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
			Celestials: &celestials.Celestial{
				Id:       "name",
				ImageUrl: "image",
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var address responses.Address
	err := json.NewDecoder(rec.Body).Decode(&address)
	s.Require().NoError(err)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(100, address.Height)
	s.Require().EqualValues(100, address.LastHeight)
	s.Require().Equal(testAddress, address.Hash)
	s.Require().NotNil(address.Celestials)
	s.Require().EqualValues("name", address.Celestials.Name)
	s.Require().EqualValues("image", address.Celestials.ImageUrl)
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
					Currency:  "utia",
					Spendable: decimal.RequireFromString("100"),
					Delegated: decimal.RequireFromString("1"),
					Unbonding: decimal.RequireFromString("2"),
				},
				Celestials: &celestials.Celestial{
					Id:       "name",
					ImageUrl: "image",
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
	s.Require().Equal("100", address[0].Balance.Spendable)
	s.Require().Equal("utia", address[0].Balance.Currency)
	s.Require().Equal("1", address[0].Balance.Delegated)
	s.Require().Equal("2", address[0].Balance.Unbonding)
	s.Require().NotNil(address[0].Celestials)
	s.Require().EqualValues("name", address[0].Celestials.Name)
	s.Require().EqualValues("image", address[0].Celestials.ImageUrl)
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
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

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
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

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
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

	s.blobLogs.EXPECT().
		BySigner(gomock.Any(), uint64(1), storage.BlobLogFilters{
			Limit: 10,
			Sort:  "desc",
			Joins: true,
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
	s.Require().Equal(testAddress, l.Signer.Hash)
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

func (s *AddressTestSuite) TestDelegations() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("show_zero", "true")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/delegations")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

	s.delegations.EXPECT().
		ByAddress(gomock.Any(), uint64(1), 10, 0, true).
		Return([]storage.Delegation{
			{
				AddressId:   1,
				ValidatorId: 1,
				Amount:      decimal.RequireFromString("100"),
				Validator:   &testValidator,
				Address: &storage.Address{
					Address: testAddress,
					Id:      1,
				},
			},
		}, nil)

	s.Require().NoError(s.handler.Delegations(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var delegations []responses.Delegation
	err := json.NewDecoder(rec.Body).Decode(&delegations)
	s.Require().NoError(err)
	s.Require().Len(delegations, 1)

	d := delegations[0]
	s.Require().Equal("100", d.Amount)
	s.Require().Equal(testAddress, d.Delegator.Hash)
	s.Require().NotNil(d.Validator)
	s.Require().Equal(testValidator.ConsAddress, d.Validator.ConsAddress)
}

func (s *AddressTestSuite) TestUndelegations() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/undelegations")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

	s.undelegations.EXPECT().
		ByAddress(gomock.Any(), uint64(1), 10, 0).
		Return([]storage.Undelegation{
			{
				Time:        testTime,
				Height:      1000,
				AddressId:   1,
				ValidatorId: 1,
				Amount:      decimal.RequireFromString("100"),
				Validator:   &testValidator,
				Address: &storage.Address{
					Address: testAddress,
					Id:      1,
				},
				CompletionTime: testTime.Add(time.Hour),
			},
		}, nil)

	s.Require().NoError(s.handler.Undelegations(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var undelegations []responses.Undelegation
	err := json.NewDecoder(rec.Body).Decode(&undelegations)
	s.Require().NoError(err)
	s.Require().Len(undelegations, 1)

	d := undelegations[0]
	s.Require().Equal("100", d.Amount)
	s.Require().EqualValues(1000, d.Height)
	s.Require().Equal(testTime, d.Time)
	s.Require().Equal(testTime.Add(time.Hour), d.CompletionTime)
	s.Require().Equal(testAddress, d.Delegator.Hash)
	s.Require().NotNil(d.Validator)
	s.Require().Equal(testValidator.ConsAddress, d.Validator.ConsAddress)
}

func (s *AddressTestSuite) TestRedelegations() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/redelegations")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

	s.redelegations.EXPECT().
		ByAddress(gomock.Any(), uint64(1), 10, 0).
		Return([]storage.Redelegation{
			{
				Time:        testTime,
				Height:      1000,
				AddressId:   1,
				SrcId:       1,
				DestId:      1,
				Amount:      decimal.RequireFromString("100"),
				Source:      &testValidator,
				Destination: &testValidator,
				Address: &storage.Address{
					Address: testAddress,
					Id:      1,
				},
				CompletionTime: testTime.Add(time.Hour),
			},
		}, nil)

	s.Require().NoError(s.handler.Redelegations(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var redelegations []responses.Redelegation
	err := json.NewDecoder(rec.Body).Decode(&redelegations)
	s.Require().NoError(err)
	s.Require().Len(redelegations, 1)

	d := redelegations[0]
	s.Require().Equal("100", d.Amount)
	s.Require().EqualValues(1000, d.Height)
	s.Require().Equal(testTime, d.Time)
	s.Require().Equal(testTime.Add(time.Hour), d.CompletionTime)
	s.Require().Equal(testAddress, d.Delegator.Hash)
	s.Require().NotNil(d.Source)
	s.Require().Equal(testValidator.ConsAddress, d.Source.ConsAddress)
	s.Require().NotNil(d.Destination)
	s.Require().Equal(testValidator.ConsAddress, d.Destination.ConsAddress)
}

func (s *AddressTestSuite) TestVestings() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/vestings")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

	s.vestings.EXPECT().
		ByAddress(gomock.Any(), uint64(1), 10, 0, false).
		Return([]storage.VestingAccount{
			{
				Time:      testTime,
				Height:    1000,
				AddressId: 1,
				Amount:    decimal.RequireFromString("100"),
				Address: &storage.Address{
					Address: testAddress,
					Id:      1,
				},
				Tx:      &testTx,
				EndTime: &testTime,
				Type:    types.VestingTypeDelayed,
			},
		}, nil)

	s.Require().NoError(s.handler.Vestings(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var vestings []responses.Vesting
	err := json.NewDecoder(rec.Body).Decode(&vestings)
	s.Require().NoError(err)
	s.Require().Len(vestings, 1)

	v := vestings[0]
	s.Require().Equal("100", v.Amount)
	s.Require().EqualValues(1000, v.Height)
	s.Require().Equal(testTime, v.Time)
	s.Require().NotNil(v.EndTime)
	s.Require().Equal(testTime, *v.EndTime)
	s.Require().Nil(v.StartTime)
}

func (s *AddressTestSuite) TestGrants() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/grants")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

	s.grants.EXPECT().
		ByGranter(gomock.Any(), uint64(1), 10, 0).
		Return([]storage.Grant{
			{
				Time:   testTime,
				Height: 1000,
				Grantee: &storage.Address{
					Address: testAddress,
					Id:      1,
				},
				Expiration: nil,
				Params: map[string]any{
					"test": "key",
				},
				Authorization: "test_msg",
			},
		}, nil)

	s.Require().NoError(s.handler.Grants(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var grants []responses.Grant
	err := json.NewDecoder(rec.Body).Decode(&grants)
	s.Require().NoError(err)
	s.Require().Len(grants, 1)

	g := grants[0]
	s.Require().Equal(testTime, g.Time)
	s.Require().EqualValues(1000, g.Height)
	s.Require().Equal(testAddress, g.Grantee.Hash)
	s.Require().Equal("test_msg", g.Authorization)
	s.Require().NotNil(g.Params)
	s.Require().False(g.Revoked)
}

func (s *AddressTestSuite) TestGrantee() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/grantee")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

	s.grants.EXPECT().
		ByGrantee(gomock.Any(), uint64(1), 10, 0).
		Return([]storage.Grant{
			{
				Time:   testTime,
				Height: 1000,
				Granter: &storage.Address{
					Address: testAddress,
					Id:      1,
				},
				Expiration: nil,
				Params: map[string]any{
					"test": "key",
				},
				Authorization: "test_msg",
			},
		}, nil)

	s.Require().NoError(s.handler.Grantee(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var grants []responses.Grant
	err := json.NewDecoder(rec.Body).Decode(&grants)
	s.Require().NoError(err)
	s.Require().Len(grants, 1)

	g := grants[0]
	s.Require().Equal(testTime, g.Time)
	s.Require().EqualValues(1000, g.Height)
	s.Require().Equal(testAddress, g.Granter.Hash)
	s.Require().Equal("test_msg", g.Authorization)
	s.Require().NotNil(g.Params)
	s.Require().False(g.Revoked)
}

func (s *AddressTestSuite) TestStats() {
	for _, name := range []string{"tx_count", "fee", "gas_used", "gas_wanted"} {
		for _, tf := range []string{"hour", "day", "month"} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/address/:hash/stats/:name/:timeframe")
			c.SetParamNames("hash", "name", "timeframe")
			c.SetParamValues(testAddress, name, tf)

			s.address.EXPECT().
				IdByHash(gomock.Any(), testHashAddress).
				Return([]uint64{1}, nil).
				Times(1)

			s.address.EXPECT().
				Series(gomock.Any(), uint64(1), storage.Timeframe(tf), name, gomock.Any()).
				Return([]storage.HistogramItem{
					{
						Time:  testTime,
						Value: "1000",
					},
				}, nil).
				Times(1)

			s.Require().NoError(s.handler.Stats(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var items []responses.HistogramItem
			err := json.NewDecoder(rec.Body).Decode(&items)
			s.Require().NoError(err)
			s.Require().Len(items, 1)

			g := items[0]
			s.Require().Equal(testTime, g.Time)
			s.Require().Equal("1000", g.Value)
		}
	}
}

func (s *AddressTestSuite) TestCelestials() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/celestials")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		IdByHash(gomock.Any(), testHashAddress).
		Return([]uint64{1}, nil).
		Times(1)

	storageResponse := make([]celestials.Celestial, 0)
	for i := 0; i < 10; i++ {
		storageResponse = append(storageResponse, celestials.Celestial{
			Id:       testsuite.RandomText(i + 10),
			ImageUrl: testsuite.RandomText(2*i + 1),
			Status:   celestials.StatusVERIFIED,
		})
	}

	s.celestials.EXPECT().
		ByAddressId(gomock.Any(), uint64(1), 10, 0).
		Return(storageResponse, nil)

	s.Require().NoError(s.handler.Celestials(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var celestials []responses.Celestial
	err := json.NewDecoder(rec.Body).Decode(&celestials)
	s.Require().NoError(err)
	s.Require().Len(celestials, len(storageResponse))

	for i := range celestials {
		s.Require().EqualValues(storageResponse[i].Id, celestials[i].Name)
		s.Require().EqualValues(storageResponse[i].ImageUrl, celestials[i].ImageUrl)
		s.Require().EqualValues("VERIFIED", celestials[i].Status)
	}
}

func (s *AddressTestSuite) TestVotes() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/votes")
	c.SetParamNames("hash")
	c.SetParamValues(testAddress)

	s.address.EXPECT().
		IdByAddress(gomock.Any(), testAddress).
		Return(123, nil)

	s.votes.EXPECT().
		ByVoterId(gomock.Any(), uint64(123), storage.VoteFilters{
			Limit:  10,
			Offset: 0,
		}).
		Return([]storage.Vote{
			{
				Id:      1,
				Height:  1000,
				Weight:  decimal.NewFromFloat(1),
				Option:  types.VoteOptionYes,
				VoterId: 1,
				Voter: &storage.Address{
					Id:         111,
					Hash:       testHashAddress,
					Address:    testAddress,
					Height:     333,
					LastHeight: 333,
					Balance: storage.Balance{
						Currency:  "utia",
						Spendable: decimal.RequireFromString("100"),
						Delegated: decimal.RequireFromString("1"),
						Unbonding: decimal.RequireFromString("2"),
					},
					Celestials: &celestials.Celestial{
						Id:       "name",
						ImageUrl: "image",
					},
				},
			},
		}, nil)

	s.Require().NoError(s.handler.Votes(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var votes []responses.Vote
	err := json.NewDecoder(rec.Body).Decode(&votes)
	s.Require().NoError(err)
	s.Require().Len(votes, 1)
	s.Require().EqualValues(1, votes[0].Id)
	s.Require().EqualValues(1000, votes[0].Height)
	s.Require().EqualValues(decimal.NewFromFloat(1), votes[0].Weight)
	s.Require().EqualValues(types.VoteOptionYes, votes[0].Option)
	s.Require().EqualValues(1, votes[0].VoterId)
	s.Require().Nil(votes[0].Validator)
	s.Require().NotNil(votes[0].Voter)
	s.Require().NotNil(votes[0].Voter.Celestials)
	s.Require().EqualValues("image", votes[0].Voter.Celestials.ImageUrl)
}
