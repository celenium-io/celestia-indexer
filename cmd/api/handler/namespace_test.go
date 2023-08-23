package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testNamespace = storage.Namespace{
		ID:          1,
		Version:     1,
		NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
		Size:        100,
	}
	testNamespaceId     = "00010203040506070809000102030405060708090001020304050607"
	testNamespaceBase64 = "AQABAgMEBQYHCAkAAQIDBAUGBwgJAAECAwQFBgc="
)

// NamespaceTestSuite -
type NamespaceTestSuite struct {
	suite.Suite
	namespaces *mock.MockINamespace
	echo       *echo.Echo
	handler    *NamespaceHandler
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *NamespaceTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.namespaces = mock.NewMockINamespace(s.ctrl)
	s.handler = NewNamespaceHandler(s.namespaces)
}

// TearDownSuite -
func (s *NamespaceTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteNamespace_Run(t *testing.T) {
	suite.Run(t, new(NamespaceTestSuite))
}

func (s *NamespaceTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id")
	c.SetParamNames("id")
	c.SetParamValues(testNamespaceId)

	s.namespaces.EXPECT().
		ByNamespaceId(gomock.Any(), testNamespace.NamespaceID).
		Return([]storage.Namespace{testNamespace}, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var namespace []responses.Namespace
	err := json.NewDecoder(rec.Body).Decode(&namespace)
	s.Require().NoError(err)
	s.Require().Len(namespace, 1)
	s.Require().EqualValues(1, namespace[0].ID)
	s.Require().EqualValues(100, namespace[0].Size)
	s.Require().EqualValues(1, namespace[0].Version)
	s.Require().Equal(testNamespaceId, namespace[0].NamespaceID)
	s.Require().Equal(testNamespaceBase64, namespace[0].Hash)
}

func (s *NamespaceTestSuite) TestGetInvalidNamespaceHeight() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id")
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *NamespaceTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace")

	s.namespaces.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.Namespace{
			&testNamespace,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var namespaces []responses.Namespace
	err := json.NewDecoder(rec.Body).Decode(&namespaces)
	s.Require().NoError(err)
	s.Require().Len(namespaces, 1)
	s.Require().EqualValues(1, namespaces[0].ID)
	s.Require().EqualValues(100, namespaces[0].Size)
	s.Require().EqualValues(1, namespaces[0].Version)
	s.Require().Equal(testNamespaceId, namespaces[0].NamespaceID)
	s.Require().Equal(testNamespaceBase64, namespaces[0].Hash)
}

func (s *NamespaceTestSuite) TestGetWithVersion() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id/:version")
	c.SetParamNames("id", "version")
	c.SetParamValues(testNamespaceId, "1")

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(1)).
		Return(testNamespace, nil)

	s.Require().NoError(s.handler.GetWithVersion(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var namespace responses.Namespace
	err := json.NewDecoder(rec.Body).Decode(&namespace)
	s.Require().NoError(err)
	s.Require().EqualValues(1, namespace.ID)
	s.Require().EqualValues(100, namespace.Size)
	s.Require().EqualValues(1, namespace.Version)
	s.Require().Equal(testNamespaceId, namespace.NamespaceID)
	s.Require().Equal(testNamespaceBase64, namespace.Hash)
}

func (s *NamespaceTestSuite) TestGetByHash() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace_by_hash/:hash")
	c.SetParamNames("hash")
	c.SetParamValues(testNamespaceBase64)

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(1)).
		Return(testNamespace, nil)

	s.Require().NoError(s.handler.GetByHash(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var namespace responses.Namespace
	err := json.NewDecoder(rec.Body).Decode(&namespace)
	s.Require().NoError(err)
	s.Require().EqualValues(1, namespace.ID)
	s.Require().EqualValues(100, namespace.Size)
	s.Require().EqualValues(1, namespace.Version)
	s.Require().Equal(testNamespaceId, namespace.NamespaceID)
	s.Require().Equal(testNamespaceBase64, namespace.Hash)
}
