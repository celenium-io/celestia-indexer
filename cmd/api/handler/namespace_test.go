package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/blob"
	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
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
	namespaces   *mock.MockINamespace
	blobReceiver *blob.MockReceiver
	echo         *echo.Echo
	handler      *NamespaceHandler
	ctrl         *gomock.Controller
}

// SetupSuite -
func (s *NamespaceTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.namespaces = mock.NewMockINamespace(s.ctrl)
	s.blobReceiver = blob.NewMockReceiver(s.ctrl)
	s.handler = NewNamespaceHandler(s.namespaces, s.blobReceiver)
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

func (s *NamespaceTestSuite) TestGetBlob() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace_by_hash/:hash/:height")
	c.SetParamNames("hash", "height")
	c.SetParamValues(testNamespaceBase64, "1000")

	result := make([]blob.Blob, 2)

	for i := 0; i < len(result); i++ {
		result[i].Namespace = testNamespaceBase64

		data := make([]byte, 88)
		_, err := rand.Read(data)
		s.Require().NoError(err)
		result[i].Data = base64.URLEncoding.EncodeToString(data)

		commitment := make([]byte, 32)
		_, err = rand.Read(commitment)
		s.Require().NoError(err)
		result[i].Commitment = base64.URLEncoding.EncodeToString(commitment)

		result[i].ShareVersion = 0
	}

	s.blobReceiver.EXPECT().
		Blobs(gomock.Any(), uint64(1000), testNamespaceBase64).
		Return(result, nil).
		MaxTimes(1).
		MinTimes(1)

	s.Require().NoError(s.handler.GetBlob(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blobs []blob.Blob
	err := json.NewDecoder(rec.Body).Decode(&blobs)
	s.Require().NoError(err)

	s.Require().Len(blobs, 2)

	blob := blobs[0]
	s.Require().EqualValues(result[0].ShareVersion, blob.ShareVersion)
	s.Require().Equal(result[0].Namespace, blob.Namespace)
	s.Require().Equal(result[0].Data, blob.Data)
	s.Require().Equal(result[0].Commitment, blob.Commitment)

}

func (s *NamespaceTestSuite) TestGetMessages() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id/:version/messages")
	c.SetParamNames("id", "version")
	c.SetParamValues(testNamespaceId, "1")

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(1)).
		Return(testNamespace, nil)

	s.namespaces.EXPECT().
		Messages(gomock.Any(), testNamespace.ID, 0, 0).
		Return([]storage.NamespaceMessage{
			{
				NamespaceId: testNamespace.ID,
				MsgId:       1,
				Message: &storage.Message{
					Id:       1,
					TxId:     2,
					Position: 3,
					Type:     types.MsgTypeBeginRedelegate,
					Height:   100,
					Time:     testTime,
				},
				TxId: 1,
				Tx:   &testTx,
			},
		}, nil)

	s.Require().NoError(s.handler.GetMessages(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var msgs []responses.NamespaceMessage
	err := json.NewDecoder(rec.Body).Decode(&msgs)
	s.Require().NoError(err)
	s.Require().Len(msgs, 1)

	msg := msgs[0]
	s.Require().EqualValues(1, msg.Id)
	s.Require().EqualValues(100, msg.Height)
	s.Require().EqualValues(3, msg.Position)
	s.Require().Equal(testTime, msg.Time)
	s.Require().EqualValues(string(types.MsgTypeBeginRedelegate), msg.Type)
	s.Require().EqualValues(1, msg.Tx.Id)
}
