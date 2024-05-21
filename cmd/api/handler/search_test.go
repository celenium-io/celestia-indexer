// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/types"
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
	validator *mock.MockIValidator
	search    *mock.MockISearch
	rollup    *mock.MockIRollup

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
	s.search = mock.NewMockISearch(s.ctrl)
	s.validator = mock.NewMockIValidator(s.ctrl)
	s.rollup = mock.NewMockIRollup(s.ctrl)
	s.handler = NewSearchHandler(s.search, s.address, s.block, s.tx, s.namespace, s.validator, s.rollup)
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
			Id:         1,
			Hash:       testHashAddress,
			Address:    testAddress,
			Height:     100,
			LastHeight: 100,
		}, nil)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("address", response.Type)
	s.Require().NotNil(response.Result)
}

func (s *SearchTestSuite) TestSearchBlock() {
	searchText := hex.EncodeToString(testBlock.Hash)

	q := make(url.Values)
	q.Set("query", searchText)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), testBlock.Hash).
		Return([]storage.SearchResult{
			{
				Id:    1,
				Type:  "block",
				Value: searchText,
			},
		}, nil).
		Times(1)

	s.block.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&testBlock, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("block", response.Type)
	s.Require().NotNil(response.Result)
}

func (s *SearchTestSuite) TestSearchBlockByHeight() {
	q := make(url.Values)
	q.Set("query", "100")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.block.EXPECT().
		ByHeight(gomock.Any(), types.Level(100)).
		Return(testBlock, nil).
		Times(1)

	s.search.EXPECT().
		SearchText(gomock.Any(), "100").
		Return([]storage.SearchResult{}, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("block", response.Type)
	s.Require().NotNil(response.Result)
	block, ok := response.Result.(map[string]any)
	s.Require().True(ok)
	_, ok = block["id"]
	s.Require().True(ok)
}

func (s *SearchTestSuite) TestSearchBlockWith0x() {
	searchText := "0x" + hex.EncodeToString(testBlock.Hash)

	q := make(url.Values)
	q.Set("query", searchText)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), testBlock.Hash).
		Return([]storage.SearchResult{
			{
				Id:    1,
				Type:  "block",
				Value: searchText,
			},
		}, nil).
		Times(1)

	s.block.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&testBlock, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("block", response.Type)
	s.Require().NotNil(response.Result)
}

func (s *SearchTestSuite) TestSearchBlockWithInvalidHash() {
	q := make(url.Values)
	q.Set("query", "EDBOFE1DAA9BB1FDA0879F1EB4F285399B6F74CB1B0C420600642682043EE41E")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		SearchText(gomock.Any(), gomock.Any()).
		Return([]storage.SearchResult{}, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)
}

func (s *SearchTestSuite) TestSearchTx() {
	q := make(url.Values)
	q.Set("query", testTxHash)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), testTx.Hash).
		Return([]storage.SearchResult{
			{
				Id:    1,
				Type:  "tx",
				Value: testTxHash,
			},
		}, nil).
		Times(1)

	s.tx.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&testTx, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("tx", response.Type)
	s.Require().NotNil(response.Result)
}

func (s *SearchTestSuite) TestSearchNamespaceById() {
	q := make(url.Values)
	q.Set("query", "00"+testNamespaceId)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.namespace.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, testNamespace.Version).
		Return(testNamespace, nil)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("namespace", response.Type)
	s.Require().NotNil(response.Result)
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
		Return(testNamespace, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("namespace", response.Type)
	s.Require().NotNil(response.Result)
}

func (s *SearchTestSuite) TestSearchValidator() {
	q := make(url.Values)
	q.Set("query", "name")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		SearchText(gomock.Any(), "name").
		Return([]storage.SearchResult{
			{
				Id:    1,
				Type:  "validator",
				Value: "name",
			},
		}, nil).
		Times(1)

	s.validator.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&storage.Validator{
			Moniker: "name 1",
			Id:      1,
			Jailed:  testsuite.Ptr(false),
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("validator", response.Type)
	s.Require().NotNil(response.Result)
}

func (s *SearchTestSuite) TestSearchRollup() {
	q := make(url.Values)
	q.Set("query", "name")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		SearchText(gomock.Any(), "name").
		Return([]storage.SearchResult{
			{
				Id:    1,
				Type:  "rollup",
				Value: "name",
			},
		}, nil).
		Times(1)

	s.rollup.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&storage.Rollup{
			Name: "name 1",
			Id:   1,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("rollup", response.Type)
	s.Require().NotNil(response.Result)
}

func (s *SearchTestSuite) TestSearchTextNamespace() {
	q := make(url.Values)
	q.Set("query", "5f45")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		SearchText(gomock.Any(), "5f45").
		Return([]storage.SearchResult{
			{
				Id:    1,
				Type:  "namespace",
				Value: "5f45",
			},
		}, nil).
		Times(1)

	s.namespace.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&storage.Namespace{
			NamespaceID: testsuite.MustHexDecode("5f45"),
			Version:     1,
			Id:          1,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().Equal("namespace", response.Type)
	s.Require().NotNil(response.Result)
}

func (s *SearchTestSuite) TestSearchNoResult() {
	q := make(url.Values)
	q.Set("query", "unknown")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		SearchText(gomock.Any(), "unknown").
		Return([]storage.SearchResult{}, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)
}

func (s *SearchTestSuite) TestSearchUnknownAddress() {
	q := make(url.Values)
	q.Set("query", testAddress)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.address.EXPECT().
		ByHash(gomock.Any(), testHashAddress).
		Return(storage.Address{}, sql.ErrNoRows)

	s.address.EXPECT().
		IsNoRows(sql.ErrNoRows).
		Return(true)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.SearchItem
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}
