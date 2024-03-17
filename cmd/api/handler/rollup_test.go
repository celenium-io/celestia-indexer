// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testRollup = storage.Rollup{
		Id:          1,
		Name:        "test rollup",
		Description: "loooooooooooooooooong description",
		Website:     "https://website.com",
		GitHub:      "https://githib.com",
		Twitter:     "https://x.com",
		Logo:        "image.png",
		Slug:        "test-rollup",
	}
	testRollupWithStats = storage.RollupWithStats{
		Rollup: testRollup,
		RollupStats: storage.RollupStats{
			BlobsCount:      100,
			Size:            1000,
			LastActionTime:  testTime,
			FirstActionTime: testTime,
		},
	}
)

// RollupTestSuite -
type RollupTestSuite struct {
	suite.Suite
	namespace *mock.MockINamespace
	rollups   *mock.MockIRollup
	blobs     *mock.MockIBlobLog
	echo      *echo.Echo
	handler   RollupHandler
	ctrl      *gomock.Controller
}

// SetupSuite -
func (s *RollupTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.namespace = mock.NewMockINamespace(s.ctrl)
	s.rollups = mock.NewMockIRollup(s.ctrl)
	s.blobs = mock.NewMockIBlobLog(s.ctrl)
	s.handler = NewRollupHandler(s.rollups, s.namespace, s.blobs)
}

// TearDownSuite -
func (s *RollupTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteRollup_Run(t *testing.T) {
	suite.Run(t, new(RollupTestSuite))
}

func (s *RollupTestSuite) TestLeaderboard() {
	for _, sort := range []string{
		"fee",
		"blobs_count",
		"time",
		"size",
	} {
		q := make(url.Values)
		q.Add("sort_by", sort)

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/rollup")

		s.rollups.EXPECT().
			Leaderboard(gomock.Any(), sort, sdk.SortOrderDesc, 10, 0).
			Return([]storage.RollupWithStats{testRollupWithStats}, nil)

		s.Require().NoError(s.handler.Leaderboard(c))
		s.Require().Equal(http.StatusOK, rec.Code)

		var rollups []responses.RollupWithStats
		err := json.NewDecoder(rec.Body).Decode(&rollups)
		s.Require().NoError(err)
		s.Require().Len(rollups, 1)

		rollup := rollups[0]
		s.Require().EqualValues(1, rollup.Id)
		s.Require().EqualValues("test rollup", rollup.Name)
		s.Require().EqualValues("image.png", rollup.Logo)
		s.Require().EqualValues("test-rollup", rollup.Slug)
		s.Require().EqualValues(100, rollup.BlobsCount)
		s.Require().EqualValues(1000, rollup.Size)
		s.Require().EqualValues(testTime, rollup.LastAction)
		s.Require().EqualValues(testTime, rollup.FirstAction)
	}
}

func (s *RollupTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.rollups.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&testRollup, nil)

	s.rollups.EXPECT().
		Stats(gomock.Any(), uint64(1)).
		Return(storage.RollupStats{
			Size:            10,
			BlobsCount:      11,
			LastActionTime:  testTime,
			FirstActionTime: testTime,
		}, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollup responses.RollupWithStats
	err := json.NewDecoder(rec.Body).Decode(&rollup)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues("test rollup", rollup.Name)
	s.Require().EqualValues("image.png", rollup.Logo)
	s.Require().EqualValues("test-rollup", rollup.Slug)
	s.Require().EqualValues(10, rollup.Size)
	s.Require().EqualValues(11, rollup.BlobsCount)
	s.Require().EqualValues(testTime, rollup.LastAction)
	s.Require().EqualValues(testTime, rollup.FirstAction)
}

func (s *RollupTestSuite) TestGetNamespaces() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:id/namespaces")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.rollups.EXPECT().
		Namespaces(gomock.Any(), uint64(1), 10, 0).
		Return([]uint64{1}, nil)

	s.namespace.EXPECT().
		GetByIds(gomock.Any(), uint64(1)).
		Return([]storage.Namespace{testNamespace}, nil)

	s.Require().NoError(s.handler.GetNamespaces(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var ns []responses.Namespace
	err := json.NewDecoder(rec.Body).Decode(&ns)
	s.Require().NoError(err)
	s.Require().Len(ns, 1)

	namespace := ns[0]
	s.Require().EqualValues(1, namespace.ID)
	s.Require().EqualValues(100, namespace.Size)
	s.Require().EqualValues(0, namespace.Version)
	s.Require().EqualValues(12, namespace.PfbCount)
	s.Require().Equal(testNamespaceId, namespace.NamespaceID)
	s.Require().Equal(testNamespaceBase64, namespace.Hash)
}

func (s *RollupTestSuite) TestGetBlobs() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:id/blobs")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.rollups.EXPECT().
		Providers(gomock.Any(), uint64(1)).
		Return([]storage.RollupProvider{
			{
				AddressId: 1,
				RollupId:  1,
			}, {
				NamespaceId: 1,
				AddressId:   2,
				RollupId:    1,
			},
		}, nil)

	s.blobs.EXPECT().
		ByProviders(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]storage.BlobLog{
			{
				Height:      100,
				Time:        testTime,
				Size:        1000,
				SignerId:    1,
				NamespaceId: 1,
				MsgId:       1,
				TxId:        1,

				Namespace: &testNamespace,
				Tx:        &testTx,
				Signer: &storage.Address{
					Address: testAddress,
				},
			},
		}, nil)

	s.Require().NoError(s.handler.GetBlobs(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var logs []responses.BlobLog
	err := json.NewDecoder(rec.Body).Decode(&logs)
	s.Require().NoError(err)
	s.Require().Len(logs, 1)
}

func (s *RollupTestSuite) TestStats() {
	for _, name := range []string{"blobs_count", "size", "size_per_blob", "fee"} {
		for _, tf := range []string{"hour", "day", "month"} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/rollup/:id/stats/:name/:timeframe")
			c.SetParamNames("id", "name", "timeframe")
			c.SetParamValues("1", name, tf)

			s.rollups.EXPECT().
				Series(gomock.Any(), uint64(1), tf, name, storage.NewSeriesRequest(0, 0)).
				Return([]storage.HistogramItem{
					{
						Value: "0.1",
						Time:  testTime,
					},
				}, nil)

			s.Require().NoError(s.handler.Stats(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var histogram []responses.HistogramItem
			err := json.NewDecoder(rec.Body).Decode(&histogram)
			s.Require().NoError(err)
			s.Require().Len(histogram, 1)
		}
	}
}

func (s *RollupTestSuite) TestDistribution() {
	for _, name := range []string{"blobs_count", "size", "size_per_blob", "fee_per_blob"} {
		for _, tf := range []string{"hour", "day"} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/rollup/:id/distribution/:name/:timeframe")
			c.SetParamNames("id", "name", "timeframe")
			c.SetParamValues("1", name, tf)

			s.rollups.EXPECT().
				Distribution(gomock.Any(), uint64(1), name, tf).
				Return([]storage.DistributionItem{
					{
						Value: "0.1",
						Name:  10,
					},
				}, nil)

			s.Require().NoError(s.handler.Distribution(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var distr []responses.DistributionItem
			err := json.NewDecoder(rec.Body).Decode(&distr)
			s.Require().NoError(err)
			s.Require().Len(distr, 1)
		}
	}
}

func (s *RollupTestSuite) TestBySlug() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/slug/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("test")

	s.rollups.EXPECT().
		BySlug(gomock.Any(), "test").
		Return(testRollup, nil)

	s.rollups.EXPECT().
		Stats(gomock.Any(), uint64(1)).
		Return(storage.RollupStats{
			Size:           10,
			BlobsCount:     11,
			LastActionTime: testTime,
		}, nil)

	s.Require().NoError(s.handler.BySlug(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollup responses.RollupWithStats
	err := json.NewDecoder(rec.Body).Decode(&rollup)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues("test rollup", rollup.Name)
	s.Require().EqualValues("test-rollup", rollup.Slug)
	s.Require().EqualValues("image.png", rollup.Logo)
	s.Require().EqualValues(10, rollup.Size)
	s.Require().EqualValues(11, rollup.BlobsCount)
	s.Require().EqualValues(testTime, rollup.LastAction)
}
