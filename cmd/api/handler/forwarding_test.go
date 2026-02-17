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

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var testForwarding = storage.Forwarding{
	Id:            1,
	Height:        100,
	Time:          testTime,
	AddressId:     1,
	DestDomain:    123456789,
	DestRecipient: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf1},
	SuccessCount:  10,
	FailedCount:   2,
	TxId:          1,
	Transfers:     json.RawMessage(`[{"denom":"utia","amount":"1000"}]`),
	Address: &storage.Address{
		Address: testAddress,
	},
	Tx: &storage.Tx{
		Hash: testTxHashBytes,
	},
}

// ForwardingTestSuite -
type ForwardingTestSuite struct {
	suite.Suite
	echo        *echo.Echo
	forwardings *mock.MockIForwarding
	address     *mock.MockIAddress
	txs         *mock.MockITx
	handler     ForwardingsHandler
	ctrl        *gomock.Controller
}

// SetupSuite -
func (s *ForwardingTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.forwardings = mock.NewMockIForwarding(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.handler = NewForwardingsHandler(s.forwardings, s.address, s.txs)
}

// TearDownSuite -
func (s *ForwardingTestSuite) TearDownSuite() {
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteForwarding_Run(t *testing.T) {
	suite.Run(t, new(ForwardingTestSuite))
}

func (s *ForwardingTestSuite) TestList() {
	q := make(url.Values)
	q.Set("limit", "5")
	q.Set("offset", "0")
	q.Set("sort", "desc")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/forwarding")

	s.forwardings.EXPECT().
		Filter(gomock.Any(), storage.ForwardingFilter{
			Limit:  5,
			Offset: 0,
			Sort:   pgSort("desc"),
		}).
		Return([]storage.Forwarding{testForwarding}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.Forwarding
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	fwd := response[0]
	s.Require().EqualValues(1, fwd.Id)
	s.Require().EqualValues(100, fwd.Height)
	s.Require().Equal(testTime, fwd.Time)
	s.Require().EqualValues(strings.ToLower(testTxHash), fwd.TxHash)
	s.Require().EqualValues(123456789, fwd.DestDomain)
	s.Require().Equal(testAddress, fwd.ForwardAddress)
	s.Require().EqualValues(10, fwd.SuccessCount)
	s.Require().EqualValues(2, fwd.FailedCount)
}

func (s *ForwardingTestSuite) TestListDefaults() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/forwarding")

	s.forwardings.EXPECT().
		Filter(gomock.Any(), storage.ForwardingFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
		}).
		Return([]storage.Forwarding{}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.Forwarding
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 0)
}

func (s *ForwardingTestSuite) TestListWithTxHash() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("tx_hash", testTxHash)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/forwarding")

	txId := uint64(1)

	s.txs.EXPECT().
		IdAndTimeByHash(gomock.Any(), testTxHashBytes).
		Return(txId, testTime, nil).
		Times(1)

	s.forwardings.EXPECT().
		Filter(gomock.Any(), storage.ForwardingFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
			TxId:   &txId,
			From:   testTime,
		}).
		Return([]storage.Forwarding{testForwarding}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.Forwarding
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)
}

func (s *ForwardingTestSuite) TestListWithAddress() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("address", testAddress)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/forwarding")

	addressId := uint64(1)

	s.address.EXPECT().
		IdByAddress(gomock.Any(), testAddress).
		Return(addressId, nil).
		Times(1)

	s.forwardings.EXPECT().
		Filter(gomock.Any(), storage.ForwardingFilter{
			Limit:     10,
			Offset:    0,
			Sort:      pgSort("desc"),
			AddressId: &addressId,
		}).
		Return([]storage.Forwarding{testForwarding}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.Forwarding
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)
}

func (s *ForwardingTestSuite) TestListWithHeight() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("height", "100")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/forwarding")

	height := uint64(100)

	s.forwardings.EXPECT().
		Filter(gomock.Any(), storage.ForwardingFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
			Height: &height,
		}).
		Return([]storage.Forwarding{testForwarding}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.Forwarding
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)
}

func (s *ForwardingTestSuite) TestListValidationError() {
	q := make(url.Values)
	q.Set("limit", "101")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/forwarding")

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)
}

func (s *ForwardingTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/forwarding/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.forwardings.EXPECT().
		ById(gomock.Any(), uint64(1)).
		Return(testForwarding, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.Forwarding
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().EqualValues(1, response.Id)
	s.Require().EqualValues(100, response.Height)
	s.Require().Equal(testTime, response.Time)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.TxHash)
	s.Require().EqualValues(123456789, response.DestDomain)
	s.Require().Equal(testAddress, response.ForwardAddress)
	s.Require().EqualValues(10, response.SuccessCount)
	s.Require().EqualValues(2, response.FailedCount)
}

func (s *ForwardingTestSuite) TestGetValidationError() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/forwarding/:id")
	c.SetParamNames("id")
	c.SetParamValues("0")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)
}
