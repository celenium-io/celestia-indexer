package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

const (
	testAddress = "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"
)

// AddressTestSuite -
type AddressTestSuite struct {
	suite.Suite
	address *mock.MockIAddress
	txs     *mock.MockITx
	echo    *echo.Echo
	handler *AddressHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *AddressTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.address = mock.NewMockIAddress(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.handler = NewAddressHandler(s.address, s.txs)
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

	hash, err := responses.DecodeAddress(testAddress)
	s.Require().NoError(err)

	s.address.EXPECT().
		ByHash(gomock.Any(), hash).
		Return(storage.Address{
			Id:      1,
			Hash:    hash,
			Height:  100,
			Balance: decimal.RequireFromString("100"),
		}, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var address responses.Address
	err = json.NewDecoder(rec.Body).Decode(&address)
	s.Require().NoError(err)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(100, address.Height)
	s.Require().Equal("100", address.Balance)
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

	hash, err := responses.DecodeAddress(testAddress)
	s.Require().NoError(err)

	s.address.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.Address{
			{
				Id:      1,
				Hash:    hash,
				Height:  100,
				Balance: decimal.RequireFromString("100"),
			},
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var address []responses.Address
	err = json.NewDecoder(rec.Body).Decode(&address)
	s.Require().NoError(err)
	s.Require().Len(address, 1)
	s.Require().EqualValues(1, address[0].Id)
	s.Require().EqualValues(100, address[0].Height)
	s.Require().Equal("100", address[0].Balance)
	s.Require().Equal(testAddress, address[0].Hash)
}

func (s *AddressTestSuite) TestListHeight() {
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

	hash, err := responses.DecodeAddress(testAddress)
	s.Require().NoError(err)

	s.address.EXPECT().
		ByHash(gomock.Any(), hash).
		Return(storage.Address{
			Id:      1,
			Hash:    hash,
			Height:  100,
			Balance: decimal.RequireFromString("100"),
		}, nil)

	s.txs.EXPECT().
		ByAddress(gomock.Any(), uint64(1), gomock.Any()).
		Return([]storage.Tx{
			testTx,
		}, nil)

	s.Require().NoError(s.handler.Transactions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err = json.NewDecoder(rec.Body).Decode(&txs)
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
