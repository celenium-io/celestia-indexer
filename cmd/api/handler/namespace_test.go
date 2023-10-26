// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	nodeMock "github.com/celenium-io/celestia-indexer/pkg/node/mock"
	nodeTypes "github.com/celenium-io/celestia-indexer/pkg/node/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testNamespace = storage.Namespace{
		Id:          1,
		Version:     1,
		NamespaceID: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7},
		Size:        100,
		PfbCount:    12,
	}
	testNamespaceId     = "00010203040506070809000102030405060708090001020304050607"
	testNamespaceBase64 = "AQABAgMEBQYHCAkAAQIDBAUGBwgJAAECAwQFBgc="
)

// NamespaceTestSuite -
type NamespaceTestSuite struct {
	suite.Suite
	namespaces   *mock.MockINamespace
	state        *mock.MockIState
	blobReceiver *nodeMock.MockDalApi
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
	s.state = mock.NewMockIState(s.ctrl)
	s.blobReceiver = nodeMock.NewMockDalApi(s.ctrl)
	s.handler = NewNamespaceHandler(s.namespaces, s.state, testIndexerName, s.blobReceiver)
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
	s.Require().EqualValues(12, namespace[0].PfbCount)
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
	s.Require().EqualValues(12, namespaces[0].PfbCount)
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
	s.Require().EqualValues(12, namespace.PfbCount)
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
	s.Require().EqualValues(12, namespace.PfbCount)
	s.Require().Equal(testNamespaceId, namespace.NamespaceID)
	s.Require().Equal(testNamespaceBase64, namespace.Hash)
}

func (s *NamespaceTestSuite) TestGetBlobs() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace_by_hash/:hash/:height")
	c.SetParamNames("hash", "height")
	c.SetParamValues(testNamespaceBase64, "1000")

	result := make([]nodeTypes.Blob, 2)

	for i := 0; i < len(result); i++ {
		result[i].Namespace = testNamespaceBase64

		data := make([]byte, 88)
		_, err := rand.Read(data)
		s.Require().NoError(err)
		result[i].Data = base64.StdEncoding.EncodeToString(data)

		commitment := make([]byte, 32)
		_, err = rand.Read(commitment)
		s.Require().NoError(err)
		result[i].Commitment = base64.StdEncoding.EncodeToString(commitment)

		result[i].ShareVersion = 0
	}

	s.blobReceiver.EXPECT().
		Blobs(gomock.Any(), pkgTypes.Level(1000), testNamespaceBase64).
		Return(result, nil).
		MaxTimes(1).
		MinTimes(1)

	s.Require().NoError(s.handler.GetBlobs(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blobs []nodeTypes.Blob
	err := json.NewDecoder(rec.Body).Decode(&blobs)
	s.Require().NoError(err)

	s.Require().Len(blobs, 2)

	blob := blobs[0]
	s.Require().EqualValues(result[0].ShareVersion, blob.ShareVersion)
	s.Require().Equal(result[0].Namespace, blob.Namespace)
	s.Require().Equal(result[0].Data, blob.Data)
	s.Require().Equal(result[0].Commitment, blob.Commitment)

}

func (s *NamespaceTestSuite) TestGetBlob() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace_by_hash/:hash/:height/:commitment")
	c.SetParamNames("hash", "height", "commitment")
	c.SetParamValues(testNamespaceBase64, "1000", "Bw==")

	data := make([]byte, 88)
	_, err := rand.Read(data)
	s.Require().NoError(err)

	result := nodeTypes.Blob{
		Namespace:    testNamespaceBase64,
		Data:         base64.StdEncoding.EncodeToString(data),
		Commitment:   "Bw==",
		ShareVersion: 0,
	}

	s.blobReceiver.EXPECT().
		Blob(gomock.Any(), pkgTypes.Level(1000), testNamespaceBase64, "Bw==").
		Return(result, nil).
		MaxTimes(1).
		MinTimes(1)

	s.Require().NoError(s.handler.GetBlob(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blob nodeTypes.Blob
	err = json.NewDecoder(rec.Body).Decode(&blob)
	s.Require().NoError(err)

	s.Require().EqualValues(0, blob.ShareVersion)
	s.Require().Equal(testNamespaceBase64, blob.Namespace)
	s.Require().Equal(result.Data, blob.Data)
	s.Require().Equal("Bw==", blob.Commitment)

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
		Messages(gomock.Any(), testNamespace.Id, 0, 0).
		Return([]storage.NamespaceMessage{
			{
				NamespaceId: testNamespace.Id,
				MsgId:       1,
				Message: &storage.Message{
					Id:       1,
					TxId:     2,
					Position: 3,
					Type:     types.MsgBeginRedelegate,
					Height:   100,
					Time:     testTime,
				},
				TxId:      1,
				Tx:        &testTx,
				Namespace: &testNamespace,
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
	s.Require().EqualValues(string(types.MsgBeginRedelegate), msg.Type)
	s.Require().EqualValues(1, msg.Tx.Id)
}

func (s *NamespaceTestSuite) TestCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/count")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(testState, nil)

	s.Require().NoError(s.handler.Count(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count uint64
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(123, count)
}

func (s *NamespaceTestSuite) TestGetActive() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/active")

	s.namespaces.EXPECT().
		Active(gomock.Any(), 5).
		Return([]storage.ActiveNamespace{
			{
				Height:    100,
				Time:      testTime,
				Namespace: testNamespace,
			},
		}, nil)

	s.Require().NoError(s.handler.GetActive(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var ns []responses.ActiveNamespace
	err := json.NewDecoder(rec.Body).Decode(&ns)
	s.Require().NoError(err)
	s.Require().Len(ns, 1)

	namespace := ns[0]
	s.Require().Equal("00010203040506070809000102030405060708090001020304050607", namespace.NamespaceID)
	s.Require().EqualValues(100, namespace.Height)
	s.Require().EqualValues(100, namespace.Size)
	s.Require().Equal(testTime, namespace.Time)
}
