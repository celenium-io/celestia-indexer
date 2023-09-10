package handler

import (
	"context"
	"database/sql"
	"encoding/hex"
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
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// SearchTestSuite -
type SearchTestSuite struct {
	suite.Suite

	address   *mock.MockIAddress
	block     *mock.MockIBlock
	namespace *mock.MockINamespace
	tx        *mock.MockITx

	echo    *echo.Echo
	handler SearchHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *SearchTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.address = mock.NewMockIAddress(s.ctrl)
	s.block = mock.NewMockIBlock(s.ctrl)
	s.namespace = mock.NewMockINamespace(s.ctrl)
	s.tx = mock.NewMockITx(s.ctrl)
	s.handler = NewSearchHandler(s.address, s.block, s.namespace, s.tx)
}

// TearDownSuite -
func (s *SearchTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteSearch_Run(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}

func (s *SearchTestSuite) TestSearchAddress() {
	q := make(url.Values)
	q.Set("query", testAddress)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.address.EXPECT().
		ByHash(gomock.Any(), testHashAddress).
		Return(storage.Address{
			Id:      1,
			Hash:    testHashAddress,
			Address: testAddress,
			Height:  100,
		}, nil)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.SearchResponse[responses.Address]
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Equal("address", response.Type)
	s.Require().EqualValues(1, response.Result.Id)
	s.Require().EqualValues(100, response.Result.Height)
	s.Require().Equal(testAddress, response.Result.Hash)
}

func (s *SearchTestSuite) TestSearchBlock() {
	q := make(url.Values)
	q.Set("query", hex.EncodeToString(testBlock.Hash))

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.tx.EXPECT().
		ByHash(gomock.Any(), testBlock.Hash).
		Return(storage.Tx{}, sql.ErrNoRows)

	s.tx.EXPECT().
		IsNoRows(sql.ErrNoRows).
		Return(true)

	s.block.EXPECT().
		ByHash(gomock.Any(), testBlock.Hash).
		Return(testBlock, nil)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.SearchResponse[responses.Block]
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Equal("block", response.Type)
	s.Require().EqualValues(1, response.Result.Id)
	s.Require().EqualValues(100, response.Result.Height)
	s.Require().Equal("1", response.Result.VersionApp)
	s.Require().Equal("11", response.Result.VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", response.Result.Hash.String())
	s.Require().Equal(testTime, response.Result.Time)
	s.Require().Equal([]types.MsgType{types.MsgSend}, response.Result.MessageTypes)
}

func (s *SearchTestSuite) TestSearchTx() {
	q := make(url.Values)
	q.Set("query", testTxHash)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.tx.EXPECT().
		ByHash(gomock.Any(), testTx.Hash).
		Return(testTx, nil)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.SearchResponse[responses.Tx]
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Equal("tx", response.Type)
	s.Require().EqualValues(1, response.Result.Id)
	s.Require().EqualValues(100, response.Result.Height)
	s.Require().Equal(testTime, response.Result.Time)
	s.Require().Equal(testTxHash, strings.ToUpper(response.Result.Hash))
	s.Require().EqualValues(2, response.Result.Position)
	s.Require().EqualValues(80410, response.Result.GasWanted)
	s.Require().EqualValues(77483, response.Result.GasUsed)
	s.Require().Equal("80410", response.Result.Fee)
	s.Require().EqualValues(0, response.Result.TimeoutHeight)
	s.Require().EqualValues(10, response.Result.EventsCount)
	s.Require().EqualValues(2, response.Result.MessagesCount)
	s.Require().Equal("memo", response.Result.Memo)
	s.Require().Equal("sdk", response.Result.Codespace)
	s.Require().Equal(types.StatusSuccess, response.Result.Status)
}

func (s *SearchTestSuite) TestSearchNamespaceById() {
	q := make(url.Values)
	q.Set("query", "01"+testNamespaceId)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.namespace.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, testNamespace.Version).
		Return(testNamespace, nil)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.SearchResponse[responses.Namespace]
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Equal("namespace", response.Type)
	s.Require().EqualValues(1, response.Result.ID)
	s.Require().EqualValues(100, response.Result.Size)
	s.Require().EqualValues(1, response.Result.Version)
	s.Require().Equal(testNamespaceId, response.Result.NamespaceID)
	s.Require().Equal(testNamespaceBase64, response.Result.Hash)
}

func (s *SearchTestSuite) TestSearchNamespaceByBase64() {
	q := make(url.Values)
	q.Set("query", testNamespaceBase64)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.namespace.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, testNamespace.Version).
		Return(testNamespace, nil)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.SearchResponse[responses.Namespace]
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Equal("namespace", response.Type)
	s.Require().EqualValues(1, response.Result.ID)
	s.Require().EqualValues(100, response.Result.Size)
	s.Require().EqualValues(1, response.Result.Version)
	s.Require().Equal(testNamespaceId, response.Result.NamespaceID)
	s.Require().Equal(testNamespaceBase64, response.Result.Hash)
}

func (s *SearchTestSuite) TestSearchNoResult() {
	q := make(url.Values)
	q.Set("query", "unknown")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusNoContent, rec.Code)
}
