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

	"github.com/celestiaorg/celestia-app/v3/pkg/proof"

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
	tendermintTypes "github.com/tendermint/tendermint/types"
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
	testTxBytes         = []byte{
		0xa, 0xcb, 0x2, 0xa, 0xa0, 0x1, 0xa, 0x9d, 0x1, 0xa, 0x20, 0x2f, 0x63, 0x65, 0x6c, 0x65, 0x73, 0x74, 0x69, 0x61,
		0x2e, 0x62, 0x6c, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67, 0x50, 0x61, 0x79, 0x46, 0x6f, 0x72,
		0x42, 0x6c, 0x6f, 0x62, 0x73, 0x12, 0x79, 0xa, 0x2f, 0x63, 0x65, 0x6c, 0x65, 0x73, 0x74, 0x69, 0x61, 0x31, 0x73,
		0x73, 0x77, 0x7a, 0x77, 0x33, 0x6d, 0x74, 0x66, 0x75, 0x65, 0x38, 0x70, 0x7a, 0x79, 0x32, 0x79, 0x6d, 0x63,
		0x37, 0x6a, 0x61, 0x78, 0x64, 0x30, 0x79, 0x72, 0x79, 0x37, 0x72, 0x6b, 0x76, 0x34, 0x34, 0x6e, 0x32, 0x75,
		0x73, 0x12, 0x1d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x56, 0x49, 0x41, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1a, 0x2, 0xab, 0x5, 0x22, 0x20, 0x86, 0x8c, 0xd6, 0x1d,
		0x4b, 0x1c, 0x90, 0xc5, 0xdd, 0x6, 0x10, 0x4e, 0xd7, 0x94, 0xcb, 0xae, 0x26, 0xf8, 0xac, 0xeb, 0x69, 0x80, 0xdb,
		0xbf, 0x1f, 0x2b, 0xb6, 0xdb, 0x98, 0x23, 0x14, 0xb2, 0x42, 0x1, 0x0, 0x12, 0x64, 0xa, 0x4e, 0xa, 0x46, 0xa,
		0x1f, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2e, 0x73, 0x65,
		0x63, 0x70, 0x32, 0x35, 0x36, 0x6b, 0x31, 0x2e, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x12, 0x23, 0xa, 0x21, 0x2,
		0x1e, 0x73, 0x47, 0xc3, 0x17, 0x85, 0x7d, 0xe, 0xa4, 0x47, 0x8f, 0x71, 0xb, 0x3, 0xbe, 0x97, 0x18, 0x84, 0x3,
		0xee, 0x91, 0x88, 0xef, 0xfe, 0x89, 0x28, 0x2c, 0x27, 0x13, 0xdc, 0x84, 0x66, 0x12, 0x4, 0xa, 0x2, 0x8, 0x1,
		0x12, 0x12, 0xa, 0xc, 0xa, 0x4, 0x75, 0x74, 0x69, 0x61, 0x12, 0x4, 0x39, 0x32, 0x32, 0x39, 0x10, 0xf9, 0xd0,
		0x5, 0x1a, 0x40, 0xe7, 0x8d, 0x57, 0x47, 0xeb, 0x49, 0xd0, 0xd4, 0x5b, 0x14, 0x98, 0xa, 0xcd, 0xcb, 0xc3, 0x5c,
		0x89, 0x88, 0xe1, 0x69, 0x4c, 0x2a, 0x64, 0x76, 0xb7, 0x74, 0x56, 0x3a, 0x81, 0x7a, 0x5c, 0x20, 0x51, 0x88,
		0x37, 0x35, 0xc7, 0x22, 0xa8, 0x0, 0x30, 0x47, 0x91, 0xce, 0xf3, 0xba, 0xf1, 0x7f, 0xa, 0xed, 0xe1, 0xb4, 0x38,
		0xce, 0xe8, 0x69, 0xe, 0x74, 0xcf, 0x24, 0x86, 0xbf, 0x7, 0x48, 0x12, 0xcc, 0x5, 0xa, 0x1c, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x56, 0x49, 0x41, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x12, 0xab, 0x5, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80, 0x1, 0x1, 0x16, 0x96, 0xe5, 0xba, 0x4e, 0x9a, 0x5, 0x5a,
		0xaf, 0xe8, 0x86, 0x47, 0xaf, 0xe3, 0xa1, 0x38, 0x90, 0x65, 0x40, 0x54, 0x8f, 0xba, 0xc9, 0x3a, 0x5a, 0x7f,
		0xdf, 0x97, 0xc6, 0x32, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x1, 0x0, 0x2, 0x42, 0x4, 0x0, 0xc, 0x71, 0xe9, 0x17, 0x21, 0xf9, 0x91, 0x85, 0x76, 0xd7, 0x60, 0xf0,
		0x2f, 0x3, 0xca, 0xc4, 0x7c, 0x6f, 0x40, 0x3, 0x31, 0x60, 0x31, 0x84, 0x8e, 0x3c, 0x1d, 0x99, 0xe6, 0xe8, 0x3a,
		0x47, 0x43, 0x29, 0x2, 0x54, 0x39, 0xaa, 0xc0, 0x2d, 0x5a, 0x69, 0x62, 0xcc, 0xce, 0xe5, 0xd4, 0xad, 0xb4, 0x8a,
		0x36, 0xbb, 0xbf, 0x44, 0x3a, 0x53, 0x17, 0x21, 0x48, 0x43, 0x81, 0x12, 0x59, 0x37, 0xf3, 0x0, 0x1a, 0xc5, 0xff,
		0x87, 0x5b, 0x19, 0x28, 0x7b, 0x98, 0x8d, 0x61, 0x7e, 0xc0, 0x5a, 0xcb, 0xbf, 0x5f, 0xe2, 0x45, 0x29, 0xa6,
		0x4b, 0x23, 0x85, 0xa9, 0x6a, 0xad, 0x43, 0xf0, 0x9b, 0xe1, 0xad, 0xa9, 0x2c, 0x70, 0x40, 0x31, 0xdc, 0xc1,
		0x48, 0x1b, 0x29, 0x2, 0x54, 0x11, 0x2f, 0x28, 0xa2, 0x54, 0x20, 0xc1, 0xd9, 0xd7, 0x5, 0x35, 0x8c, 0x13, 0x4c,
		0xc6, 0x1, 0xd9, 0xd1, 0x84, 0xcb, 0x4d, 0xfd, 0xde, 0x7e, 0x1c, 0xac, 0x2b, 0xc3, 0xd4, 0xd3, 0x8b, 0xf9, 0xec,
		0x44, 0xe6, 0x9, 0x64, 0x12, 0x3b, 0xaf, 0xc5, 0x86, 0xf7, 0x77, 0x64, 0x48, 0x8c, 0xd2, 0x4c, 0x6a, 0x77, 0x54,
		0x6e, 0x5a, 0xf, 0xe8, 0xbd, 0xfb, 0x4f, 0xa2, 0x3, 0xcf, 0xaf, 0xfc, 0x36, 0xcc, 0xe4, 0xdd, 0x5b, 0x89, 0x1,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x67, 0x50, 0x70, 0x0, 0x8e, 0x7d, 0xd0, 0x6a, 0xc5,
		0xb7, 0x3b, 0x47, 0x3b, 0xe6, 0xbc, 0x5a, 0x51, 0x3, 0xf, 0x4c, 0x74, 0x37, 0x65, 0x7c, 0xb7, 0xb2, 0x9b, 0xf3,
		0x76, 0xc5, 0x64, 0xb8, 0xd1, 0x67, 0x5a, 0x5e, 0x89, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x67, 0x50, 0x70, 0x1, 0x4b, 0xa8, 0x4e, 0x1f, 0x37, 0xd0, 0x41, 0xbc, 0x6e, 0x55, 0xba, 0x39, 0x68,
		0x26, 0xcc, 0x49, 0x4e, 0x84, 0xd4, 0x81, 0x5b, 0x6d, 0xb5, 0x26, 0x90, 0x42, 0x2e, 0xea, 0x73, 0x86, 0x31,
		0x4f, 0x0, 0xe8, 0xe7, 0x76, 0x26, 0x58, 0x6f, 0x73, 0xb9, 0x55, 0x36, 0x4c, 0x7b, 0x4b, 0xbf, 0xb, 0xb7, 0xf7,
		0x68, 0x5e, 0xbd, 0x40, 0xe8, 0x52, 0xb1, 0x64, 0x63, 0x3a, 0x4a, 0xcb, 0xd3, 0x24, 0x4c, 0x3d, 0xe2, 0x20,
		0x2c, 0xcb, 0x62, 0x6a, 0xd3, 0x87, 0xd7, 0x7, 0x22, 0xe6, 0x4f, 0xbe, 0x44, 0x56, 0x2e, 0x2f, 0x23, 0x1a, 0x29,
		0xc, 0x8, 0x53, 0x2b, 0x8d, 0x6a, 0xba, 0x40, 0x2f, 0xf5, 0x0, 0x96, 0xf9, 0x7a, 0x20, 0x84, 0x4a, 0x4f, 0xba,
		0x58, 0x29, 0x5, 0x1, 0x37, 0xa, 0xa1, 0x44, 0xb7, 0x7d, 0x26, 0x25, 0x3a, 0x38, 0x4d, 0x27, 0x4a, 0xec, 0x57,
		0x22, 0xce, 0x21, 0xf8, 0xf1, 0x79, 0x9, 0x35, 0x88, 0xd0, 0xe8, 0x47, 0xef, 0xa7, 0x3a, 0x10, 0xce, 0x20, 0xe4,
		0x79, 0x9f, 0xb1, 0xe4, 0x66, 0x42, 0xd6, 0x56, 0x17, 0xc7, 0xe5, 0x21, 0x3f, 0xa0, 0x49, 0x89, 0xd9, 0x2d,
		0x89, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x67, 0x50, 0x70, 0x1, 0x87, 0xde, 0xd2,
		0x47, 0xe1, 0x66, 0xf, 0x82, 0x70, 0x71, 0xc7, 0xf1, 0x37, 0x19, 0x34, 0x58, 0x97, 0x51, 0x8, 0x53, 0x84, 0xfc,
		0x9f, 0x44, 0x62, 0xc1, 0xf1, 0x89, 0x7c, 0x5c, 0x3e, 0xef, 0x89, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x86, 0x24, 0x81, 0x93, 0xeb, 0x4d, 0xd2, 0xa8, 0xce, 0x81, 0x5f, 0x87,
		0x6c, 0x12, 0x4d, 0x48, 0x35, 0x95, 0x22, 0xf0, 0x85, 0x4d, 0x95, 0xd8, 0x7, 0x2e, 0xaf, 0xf0, 0xd3, 0x7d, 0x55,
		0xbd, 0x11, 0x3, 0x20, 0x91, 0x1d, 0xd2, 0xad, 0x74, 0x3f, 0xf2, 0x37, 0xd4, 0x11, 0x64, 0x8a, 0xf, 0xe3, 0x2c,
		0x6d, 0x74, 0xee, 0xc0, 0x60, 0x71, 0x6a, 0x2a, 0x74, 0x35, 0x2f, 0x6b, 0x1c, 0x43, 0x5b, 0x5d, 0x67, 0x0, 0xa8,
		0xba, 0x1e, 0x76, 0xb4, 0xd8, 0x56, 0xb, 0xee, 0x8c, 0x3, 0x9c, 0x91, 0xd8, 0x8b, 0x48, 0x60, 0xca, 0xe7, 0x9b,
		0x5, 0xd4, 0xb2, 0xd6, 0x56, 0xf0, 0xd3, 0x8d, 0x78, 0x56, 0xa2, 0xdd, 0x1a, 0x4, 0x42, 0x4c, 0x4f, 0x42}
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
	s.node = nodeMock.NewMockApi(s.ctrl)
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
			To:    testNamespace.LastMessageTime,
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
	s.Require().Equal(testAddress, l.Signer.Hash)
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
			To:      testNamespace.LastMessageTime,
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
	s.Require().Equal(testAddress, l.Signer.Hash)
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
			To:         testNamespace.LastMessageTime,
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
	s.Require().Equal(testAddress, l.Signer.Hash)
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
	s.Require().EqualValues(testAddress, blob.Signer.Hash)
	s.Require().EqualValues(testNamespace.Hash(), blob.Namespace)
}

func (s *NamespaceTestSuite) TestBlobProofs() {
	commitment := "hozWHUsckMXdBhBO15TLrib4rOtpgNu/Hyu225gjFLI="
	namespaceBase64 := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAVklBAAAAAAA="
	blockHeight := pkgTypes.Level(2892352)
	txs := [][]byte{testTxBytes}

	proofReq := map[string]any{
		"hash":       namespaceBase64,
		"height":     blockHeight,
		"commitment": commitment,
	}
	stream := new(bytes.Buffer)
	err := json.NewEncoder(stream).Encode(proofReq)
	s.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/", stream)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/blob/proofs")

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	s.node.EXPECT().
		Block(gomock.Any(), blockHeight).
		Return(pkgTypes.ResultBlock{
			Block: &pkgTypes.Block{
				Data: pkgTypes.Data{
					Txs: tendermintTypes.ToTxs(txs),
				},
			},
		}, nil)

	s.Require().NoError(s.handler.BlobProofs(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var proofs []*proof.NMTProof
	err = json.NewDecoder(rec.Body).Decode(&proofs)
	s.Require().NoError(err)
	s.Require().EqualValues(2, len(proofs))
	s.Require().EqualValues(1, proofs[0].Start)
	s.Require().EqualValues(2, proofs[0].End)
	s.Require().EqualValues(1, proofs[1].End)
	s.Require().EqualValues(
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABJKmGf0yBpvzj9FXhOZbyVHBq1LBKP5eI34r38GvSPZ9",
		base64.StdEncoding.EncodeToString(proofs[0].Nodes[0]),
	)
	s.Require().EqualValues(
		"/////////////////////////////////////////////////////////////////////////////55b1xq+xbf5Olhxqwf4N8vvgExxFkjtgX4j/uCbESK9",
		base64.StdEncoding.EncodeToString(proofs[0].Nodes[1]),
	)
	s.Require().EqualValues(
		"//////////////////////////////////////7//////////////////////////////////////plEqgR/c4IAVkNdYRWOYOAESD4whneKR54Dz5Dfe4p2",
		base64.StdEncoding.EncodeToString(proofs[1].Nodes[0]),
	)
	s.Require().EqualValues(
		"/////////////////////////////////////////////////////////////////////////////9DUGO+QswnUItJQnpHTEz6nj13KfN9iuS9pG2tYF5tI",
		base64.StdEncoding.EncodeToString(proofs[1].Nodes[1]),
	)
}
