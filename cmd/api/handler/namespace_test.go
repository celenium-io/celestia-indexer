// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"

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
		Id:              1,
		Version:         0,
		NamespaceID:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xfc, 0x74, 0x43, 0xb1, 0x55, 0x92, 0x1, 0x56},
		Size:            100,
		PfbCount:        12,
		LastHeight:      100,
		LastMessageTime: testTime,
	}
	testNamespaceId     = "0000000000000000000000000000000000000000fc7443b155920156"
	testNamespaceBase64 = "AAAAAAAAAAAAAAAAAAAAAAAAAAAA/HRDsVWSAVY="
)

// NamespaceTestSuite -
type NamespaceTestSuite struct {
	suite.Suite
	namespaces   *mock.MockINamespace
	blobLogs     *mock.MockIBlobLog
	rollups      *mock.MockIRollup
	address      *mock.MockIAddress
	state        *mock.MockIState
	blobReceiver *nodeMock.MockDalApi
	node         *nodeMock.MockApi
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
	s.blobLogs = mock.NewMockIBlobLog(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
	s.rollups = mock.NewMockIRollup(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.blobReceiver = nodeMock.NewMockDalApi(s.ctrl)
	s.handler = NewNamespaceHandler(
		s.namespaces,
		s.blobLogs,
		s.rollups,
		s.address,
		s.state,
		testIndexerName,
		s.blobReceiver,
		s.node,
	)
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
	s.Require().EqualValues(0, namespace[0].Version)
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
		ListWithSort(gomock.Any(), "", sdk.SortOrderDesc, 10, 0).
		Return([]storage.Namespace{
			testNamespace,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var namespaces []responses.Namespace
	err := json.NewDecoder(rec.Body).Decode(&namespaces)
	s.Require().NoError(err)
	s.Require().Len(namespaces, 1)
	s.Require().EqualValues(1, namespaces[0].ID)
	s.Require().EqualValues(100, namespaces[0].Size)
	s.Require().EqualValues(0, namespaces[0].Version)
	s.Require().EqualValues(12, namespaces[0].PfbCount)
	s.Require().Equal(testNamespaceId, namespaces[0].NamespaceID)
	s.Require().Equal(testNamespaceBase64, namespaces[0].Hash)
}

func (s *NamespaceTestSuite) TestListWithSort() {
	for _, request := range []namespaceList{
		{
			Sort:   "asc",
			SortBy: "size",
		}, {
			Sort:   "desc",
			SortBy: "size",
		}, {
			Sort:   "asc",
			SortBy: "pfb_count",
		}, {
			Sort:   "asc",
			SortBy: "pfb_count",
		}, {
			Sort:   "asc",
			SortBy: "time",
		}, {
			Sort:   "asc",
			SortBy: "time",
		},
	} {
		q := make(url.Values)
		q.Set("sort", request.Sort)
		q.Set("sort_by", request.SortBy)

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/namespace")

		s.namespaces.EXPECT().
			ListWithSort(gomock.Any(), request.SortBy, pgSort(request.Sort), 10, 0).
			Return([]storage.Namespace{
				testNamespace,
			}, nil)

		s.Require().NoError(s.handler.List(c))
		s.Require().Equal(http.StatusOK, rec.Code)

		var namespaces []responses.Namespace
		err := json.NewDecoder(rec.Body).Decode(&namespaces)
		s.Require().NoError(err)
		s.Require().Len(namespaces, 1)
		s.Require().EqualValues(1, namespaces[0].ID)
		s.Require().EqualValues(100, namespaces[0].Size)
		s.Require().EqualValues(0, namespaces[0].Version)
		s.Require().EqualValues(12, namespaces[0].PfbCount)
		s.Require().Equal(testNamespaceId, namespaces[0].NamespaceID)
		s.Require().Equal(testNamespaceBase64, namespaces[0].Hash)
	}
}

func (s *NamespaceTestSuite) TestGetWithVersion() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id/:version")
	c.SetParamNames("id", "version")
	c.SetParamValues(testNamespaceId, "0")

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(0)).
		Return(testNamespace, nil)

	s.Require().NoError(s.handler.GetWithVersion(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var namespace responses.Namespace
	err := json.NewDecoder(rec.Body).Decode(&namespace)
	s.Require().NoError(err)
	s.Require().EqualValues(1, namespace.ID)
	s.Require().EqualValues(100, namespace.Size)
	s.Require().EqualValues(0, namespace.Version)
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
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(0)).
		Return(testNamespace, nil)

	s.Require().NoError(s.handler.GetByHash(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var namespace responses.Namespace
	err := json.NewDecoder(rec.Body).Decode(&namespace)
	s.Require().NoError(err)
	s.Require().EqualValues(1, namespace.ID)
	s.Require().EqualValues(100, namespace.Size)
	s.Require().EqualValues(0, namespace.Version)
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

func (s *NamespaceTestSuite) TestGetMessages() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id/:version/messages")
	c.SetParamNames("id", "version")
	c.SetParamValues(testNamespaceId, "0")

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(0)).
		Return(testNamespace, nil)

	s.namespaces.EXPECT().
		Messages(gomock.Any(), testNamespace.Id, 10, 0).
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
		ListWithSort(gomock.Any(), "time", sdk.SortOrderDesc, 5, 0).
		Return([]storage.Namespace{
			testNamespace,
		}, nil)

	s.Require().NoError(s.handler.GetActive(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var ns []responses.Namespace
	err := json.NewDecoder(rec.Body).Decode(&ns)
	s.Require().NoError(err)
	s.Require().Len(ns, 1)

	namespace := ns[0]
	s.Require().Equal("0000000000000000000000000000000000000000fc7443b155920156", namespace.NamespaceID)
	s.Require().EqualValues(100, namespace.LastHeight)
	s.Require().EqualValues(100, namespace.Size)
	s.Require().Equal(testTime, namespace.LastMessageTime)
}

func (s *NamespaceTestSuite) TestGetActiveWithSort() {
	for _, field := range []string{"pfb_count", "time", "size"} {
		q := make(url.Values)
		q.Set("sort", field)

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/namespace/active")

		s.namespaces.EXPECT().
			ListWithSort(gomock.Any(), field, sdk.SortOrderDesc, 5, 0).
			Return([]storage.Namespace{
				testNamespace,
			}, nil)

		s.Require().NoError(s.handler.GetActive(c))
		s.Require().Equal(http.StatusOK, rec.Code)

		var ns []responses.Namespace
		err := json.NewDecoder(rec.Body).Decode(&ns)
		s.Require().NoError(err)
		s.Require().Len(ns, 1)

		namespace := ns[0]
		s.Require().Equal("0000000000000000000000000000000000000000fc7443b155920156", namespace.NamespaceID)
		s.Require().EqualValues(100, namespace.LastHeight)
		s.Require().EqualValues(100, namespace.Size)
		s.Require().Equal(testTime, namespace.LastMessageTime)

	}
}

func (s *NamespaceTestSuite) TestBlob() {
	commitment := "ZeKGjIwsIkFsACD0wtEh/jbzzW+zIPP716VihNpm9T0="

	blobReq := map[string]any{
		"hash":       testNamespaceBase64,
		"height":     1000,
		"commitment": commitment,
	}
	stream := new(bytes.Buffer)
	err := json.NewEncoder(stream).Encode(blobReq)
	s.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/", stream)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/blob")

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	data := make([]byte, 88)
	_, err = rand.Read(data)
	s.Require().NoError(err)

	result := nodeTypes.Blob{
		Namespace:    testNamespaceBase64,
		Data:         base64.StdEncoding.EncodeToString(data),
		Commitment:   commitment,
		ShareVersion: 0,
	}

	s.blobReceiver.EXPECT().
		Blob(gomock.Any(), pkgTypes.Level(1000), testNamespaceBase64, commitment).
		Return(result, nil).
		MaxTimes(1).
		MinTimes(1)

	s.Require().NoError(s.handler.Blob(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blob responses.Blob
	err = json.NewDecoder(rec.Body).Decode(&blob)
	s.Require().NoError(err)

	s.Require().EqualValues(0, blob.ShareVersion)
	s.Require().Equal(testNamespaceBase64, blob.Namespace)
	s.Require().Equal(result.Data, blob.Data)
	s.Require().Equal(commitment, blob.Commitment)

}

func (s *NamespaceTestSuite) TestGetLogs() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id/:version/logs")
	c.SetParamNames("id", "version")
	c.SetParamValues(testNamespaceId, "0")

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(0)).
		Return(testNamespace, nil)

	s.blobLogs.EXPECT().
		ByNamespace(gomock.Any(), testNamespace.Id, storage.BlobLogFilters{
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
				Rollup:     &testRollup,
			},
		}, nil)

	s.Require().NoError(s.handler.GetBlobLogs(c))
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
	s.Require().NotNil(l.Rollup)
}

func (s *NamespaceTestSuite) TestGetLogsBySigner() {
	args := make(url.Values)
	args.Set("signers", "celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r,celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6")

	req := httptest.NewRequest(http.MethodGet, "/?"+args.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id/:version/logs")
	c.SetParamNames("id", "version")
	c.SetParamValues(testNamespaceId, "0")

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(0)).
		Return(testNamespace, nil)

	_, h1, err := pkgTypes.Address("celestia12zs7e3n8pjd8y8ex0cyv67ethv30mekgqu665r").Decode()
	s.Require().NoError(err)
	_, h2, err := pkgTypes.Address("celestia1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3y3clr6").Decode()
	s.Require().NoError(err)

	s.address.EXPECT().
		IdByHash(gomock.Any(), h1, h2).
		Return([]uint64{1, 2}, nil).
		Times(1)

	s.blobLogs.EXPECT().
		ByNamespace(gomock.Any(), testNamespace.Id, storage.BlobLogFilters{
			Limit:   10,
			Sort:    "desc",
			Joins:   true,
			Signers: []uint64{1, 2},
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
				Rollup:     &testRollup,
			},
		}, nil)

	s.Require().NoError(s.handler.GetBlobLogs(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var logs []responses.BlobLog
	err = json.NewDecoder(rec.Body).Decode(&logs)
	s.Require().NoError(err)
	s.Require().Len(logs, 1)

	l := logs[0]
	s.Require().EqualValues(10000, l.Height)
	s.Require().Equal(testTime, l.Time)
	s.Require().Equal(testAddress, l.Signer)
	s.Require().Equal("test_commitment", l.Commitment)
	s.Require().EqualValues(1000, l.Size)
	s.Require().Nil(l.Namespace)
	s.Require().NotNil(l.Rollup)
}

func (s *NamespaceTestSuite) TestGetLogsWithCommitment() {
	cm := "T1EPYi3jq6hC3ueLOZRtWB7LUsAC4DcnAX_oSwDopps="
	args := make(url.Values)
	args.Set("commitment", cm)

	req := httptest.NewRequest(http.MethodGet, "/?"+args.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id/:version/logs")
	c.SetParamNames("id", "version")
	c.SetParamValues(testNamespaceId, "0")

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(0)).
		Return(testNamespace, nil)

	s.blobLogs.EXPECT().
		ByNamespace(gomock.Any(), testNamespace.Id, storage.BlobLogFilters{
			Limit:      10,
			Sort:       "desc",
			Commitment: "T1EPYi3jq6hC3ueLOZRtWB7LUsAC4DcnAX/oSwDopps=",
			Joins:      true,
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
				Commitment: "T1EPYi3jq6hC3ueLOZRtWB7LUsAC4DcnAX/oSwDopps=",
				Size:       1000,
				Height:     10000,
				Time:       testTime,
			},
		}, nil)

	s.Require().NoError(s.handler.GetBlobLogs(c))
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())

	var logs []responses.BlobLog
	err := json.NewDecoder(rec.Body).Decode(&logs)
	s.Require().NoError(err)
	s.Require().Len(logs, 1)

	l := logs[0]
	s.Require().EqualValues(10000, l.Height)
	s.Require().Equal(testTime, l.Time)
	s.Require().Equal(testAddress, l.Signer)
	s.Require().Equal("T1EPYi3jq6hC3ueLOZRtWB7LUsAC4DcnAX/oSwDopps=", l.Commitment)
	s.Require().EqualValues(1000, l.Size)
	s.Require().Nil(l.Namespace)
}

func (s *NamespaceTestSuite) TestRollups() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/namespace/:id/:version/rollups")
	c.SetParamNames("id", "version")
	c.SetParamValues(testNamespaceId, "0")

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), testNamespace.NamespaceID, byte(0)).
		Return(testNamespace, nil)

	s.rollups.EXPECT().
		RollupsByNamespace(gomock.Any(), testNamespace.Id, 10, 0).
		Return([]storage.Rollup{testRollup}, nil)

	s.Require().NoError(s.handler.Rollups(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollups []responses.Rollup
	err := json.NewDecoder(rec.Body).Decode(&rollups)
	s.Require().NoError(err)
	s.Require().Len(rollups, 1)

	rollup := rollups[0]
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues("test rollup", rollup.Name)
	s.Require().EqualValues("image.png", rollup.Logo)
	s.Require().EqualValues("test-rollup", rollup.Slug)
}

func (s *NamespaceTestSuite) TestBlobMetadata() {
	commitment := "ZeKGjIwsIkFsACD0wtEh/jbzzW+zIPP716VihNpm9T1="

	blobReq := map[string]any{
		"hash":       testNamespaceBase64,
		"height":     1000,
		"commitment": commitment,
	}
	stream := new(bytes.Buffer)
	err := json.NewEncoder(stream).Encode(blobReq)
	s.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/", stream)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/blob/metadata")

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	data := make([]byte, 88)
	_, err = rand.Read(data)
	s.Require().NoError(err)

	s.blobLogs.EXPECT().
		Blob(gomock.Any(), pkgTypes.Level(1000), uint64(1), "ZeKGjIwsIkFsACD0wtEh/jbzzW+zIPP716VihNpm9T1=").
		Return(storage.BlobLog{
			NamespaceId: testNamespace.Id,
			MsgId:       1,
			TxId:        1,
			SignerId:    1,
			Signer: &storage.Address{
				Address: testAddress,
			},
			Commitment: "test_commitment",
			Size:       1000,
			Height:     1000,
			Time:       testTime,
			Tx:         &testTx,
			Namespace:  &testNamespace,
		}, nil).
		Times(1)

	s.namespaces.EXPECT().
		ByNamespaceIdAndVersion(gomock.Any(), gomock.Any(), uint8(0)).
		Return(testNamespace, nil).
		Times(1)

	s.Require().NoError(s.handler.BlobMetadata(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blob responses.BlobLog
	err = json.NewDecoder(rec.Body).Decode(&blob)
	s.Require().NoError(err)
	s.Require().NotNil(blob.Namespace)
	s.Require().NotNil(blob.Tx)
	s.Require().NotNil(blob.Signer)
}

func (s *NamespaceTestSuite) TestBlobs() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/blobs")

	s.blobLogs.EXPECT().
		ListBlobs(gomock.Any(), gomock.Any()).
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
				Height:     1000,
				Time:       testTime,
				Tx:         &testTx,
				Namespace:  &testNamespace,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Blobs(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blobs []responses.LightBlobLog
	err := json.NewDecoder(rec.Body).Decode(&blobs)
	s.Require().NoError(err)
	s.Require().Len(blobs, 1)

	blob := blobs[0]
	s.Require().EqualValues(1000, blob.Size)
	s.Require().EqualValues(1000, blob.Height)
	s.Require().EqualValues(testTime, blob.Time)
	s.Require().EqualValues("test_commitment", blob.Commitment)
	s.Require().EqualValues(testAddress, blob.Signer)
	s.Require().EqualValues(testNamespace.Hash(), blob.Namespace)
}
