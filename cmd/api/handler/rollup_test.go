// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
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
		Tags:        []string{"ai"},
		Category:    types.RollupCategoryNft,
		Type:        types.RollupTypeSettled,
		Stack:       "stack 1, stack 2",
		Provider:    "Provider 1",
		Color:       "#123456",
	}
	testRollupWithStats = storage.RollupWithStats{
		Rollup: testRollup,
		RollupStats: storage.RollupStats{
			BlobsCount:      100,
			Size:            1000,
			LastActionTime:  testTime,
			FirstActionTime: testTime,
			BlobsCountPct:   0.1,
			FeePct:          0.2,
			SizePct:         0.3,
		},
		DAChange: storage.DAChange{
			DAPct: 0.1,
		},
	}
	testRollupWithGroupedStats = storage.RollupGroupedStats{
		Fee:        0.1,
		Size:       0.2,
		BlobsCount: 3,
		Group:      "stack",
	}
	testRollupActivity = true
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
		q.Add("type", "sovereign")
		q.Add("category", "nft,gaming")
		q.Add("tags", "ai")
		q.Add("stack", "stack 1,stack 2")
		q.Add("provider", "provider 1")
		q.Add("is_active", "true")

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/rollup")

		s.rollups.EXPECT().
			Leaderboard(gomock.Any(), storage.LeaderboardFilters{
				SortField: sort,
				Sort:      sdk.SortOrderDesc,
				Limit:     10,
				Offset:    0,
				Category: []types.RollupCategory{
					types.RollupCategoryNft,
					types.RollupCategoryGaming,
				},
				Type: []types.RollupType{
					types.RollupTypeSovereign,
				},
				Tags:     []string{"ai"},
				Stack:    []string{"stack 1", "stack 2"},
				Provider: []string{"provider 1"},
				IsActive: &testRollupActivity,
			}).
			Return([]storage.RollupWithStats{testRollupWithStats}, nil).
			Times(1)

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
		s.Require().EqualValues(0.1, rollup.BlobsCountPct)
		s.Require().EqualValues(0.2, rollup.FeePct)
		s.Require().EqualValues(0.3, rollup.SizePct)
		s.Require().EqualValues(0.1, rollup.DAPct)
		s.Require().EqualValues("#123456", rollup.Color)
	}
}

func (s *RollupTestSuite) TestLeaderboardDay() {
	for _, sort := range []string{
		"avg_size", "blobs_count", "total_size", "total_fee", "throughput", "namespace_count", "pfb_count", "mb_price",
	} {
		q := make(url.Values)
		q.Add("sort_by", sort)
		q.Add("category", "nft,gaming")
		q.Add("type", "sovereign")
		q.Add("tags", "ai")
		q.Add("stack", "stack 1,stack 2")
		q.Add("provider", "provider 1")

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/rollup/day")

		s.rollups.EXPECT().
			LeaderboardDay(gomock.Any(), storage.LeaderboardFilters{
				SortField: sort,
				Sort:      sdk.SortOrderDesc,
				Limit:     10,
				Offset:    0,
				Category: []types.RollupCategory{
					types.RollupCategoryNft,
					types.RollupCategoryGaming,
				},
				Type: []types.RollupType{
					types.RollupTypeSovereign,
				},
				Tags:     []string{"ai"},
				Stack:    []string{"stack 1", "stack 2"},
				Provider: []string{"provider 1"},
			}).
			Return([]storage.RollupWithDayStats{
				{
					Rollup: testRollup,
					RolluDayStats: storage.RolluDayStats{
						BlobsCount: 100,
					},
				},
			}, nil).
			Times(1)

		s.Require().NoError(s.handler.LeaderboardDay(c))
		s.Require().Equal(http.StatusOK, rec.Code)

		var rollups []responses.RollupWithDayStats
		err := json.NewDecoder(rec.Body).Decode(&rollups)
		s.Require().NoError(err)
		s.Require().Len(rollups, 1)

		rollup := rollups[0]
		s.Require().EqualValues(1, rollup.Id)
		s.Require().EqualValues("test rollup", rollup.Name)
		s.Require().EqualValues("image.png", rollup.Logo)
		s.Require().EqualValues("test-rollup", rollup.Slug)
		s.Require().EqualValues(100, rollup.BlobsCount)
		s.Require().EqualValues("#123456", rollup.Color)
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
		ById(gomock.Any(), uint64(1)).
		Return(testRollupWithStats, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollup responses.RollupWithStats
	err := json.NewDecoder(rec.Body).Decode(&rollup)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues("test rollup", rollup.Name)
	s.Require().EqualValues("image.png", rollup.Logo)
	s.Require().EqualValues("test-rollup", rollup.Slug)
	s.Require().EqualValues(100, rollup.BlobsCount)
	s.Require().EqualValues(1000, rollup.Size)
	s.Require().EqualValues(testTime, rollup.LastAction)
	s.Require().EqualValues(testTime, rollup.FirstAction)
	s.Require().EqualValues(0.1, rollup.BlobsCountPct)
	s.Require().EqualValues(0.2, rollup.FeePct)
	s.Require().EqualValues(0.3, rollup.SizePct)
	s.Require().EqualValues("#123456", rollup.Color)
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

	s.rollups.EXPECT().
		ById(gomock.Any(), uint64(1)).
		Return(storage.RollupWithStats{
			Rollup: testRollup,
			RollupStats: storage.RollupStats{
				LastActionTime:  testTime,
				FirstActionTime: testTime,
			},
		}, nil).
		Times(1)

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
		for _, tf := range []storage.Timeframe{storage.TimeframeHour, storage.TimeframeDay, storage.TimeframeMonth} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/rollup/:id/stats/:name/:timeframe")
			c.SetParamNames("id", "name", "timeframe")
			c.SetParamValues("1", name, string(tf))

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
		for _, tf := range []storage.Timeframe{storage.TimeframeHour, storage.TimeframeDay} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/rollup/:id/distribution/:name/:timeframe")
			c.SetParamNames("id", "name", "timeframe")
			c.SetParamValues("1", name, string(tf))

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
		Return(testRollupWithStats, nil).
		Times(1)

	s.Require().NoError(s.handler.BySlug(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollup responses.RollupWithStats
	err := json.NewDecoder(rec.Body).Decode(&rollup)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues("test rollup", rollup.Name)
	s.Require().EqualValues("test-rollup", rollup.Slug)
	s.Require().EqualValues("image.png", rollup.Logo)
	s.Require().EqualValues(1000, rollup.Size)
	s.Require().EqualValues(100, rollup.BlobsCount)
	s.Require().EqualValues(testTime, rollup.LastAction)
}

func (s *RollupTestSuite) TestByExportBlobs() {
	q := make(url.Values)
	q.Set("from", "1")
	q.Set("to", "2")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:id/export")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.rollups.EXPECT().
		Providers(gomock.Any(), uint64(1)).
		Return([]storage.RollupProvider{
			{
				RollupId:    1,
				NamespaceId: 2,
				AddressId:   3,
			},
		}, nil)

	from := time.Unix(1, 0).UTC()
	to := time.Unix(2, 0).UTC()
	s.blobs.EXPECT().
		ExportByProviders(gomock.Any(), []storage.RollupProvider{
			{
				RollupId:    1,
				NamespaceId: 2,
				AddressId:   3,
			},
		}, from, to, gomock.Any()).
		Return(nil)

	s.Require().NoError(s.handler.ExportBlobs(c))
	s.Require().Equal(http.StatusOK, rec.Code)
}

func (s *RollupTestSuite) TestAllSeries() {
	for _, tf := range []storage.Timeframe{
		storage.TimeframeHour,
		storage.TimeframeDay,
		storage.TimeframeMonth,
	} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/rollup/stats/series/:timeframe")
		c.SetParamNames("timeframe")
		c.SetParamValues(string(tf))

		s.rollups.EXPECT().
			AllSeries(gomock.Any(), tf).
			Return([]storage.RollupHistogramItem{
				{
					Name:       testRollup.Name,
					Logo:       testRollup.Logo,
					Time:       testTime,
					BlobsCount: 1,
					Size:       2,
					Fee:        "3",
				},
			}, nil).
			Times(1)

		s.Require().NoError(s.handler.AllSeries(c))
		s.Require().Equal(http.StatusOK, rec.Code)

		var items []responses.RollupAllSeriesResponse
		err := json.NewDecoder(rec.Body).Decode(&items)
		s.Require().NoError(err)
		s.Require().Len(items, 1)

		for _, item := range items {
			s.Require().Equal(testTime.String(), item.Time.String())
			s.Require().Len(item.Items, 1)
			s.Require().EqualValues("test rollup", item.Items[0].Name)
			s.Require().EqualValues("3", item.Items[0].Fee)
			s.Require().EqualValues("image.png", item.Items[0].Logo)
			s.Require().EqualValues(2, item.Items[0].Size)
			s.Require().EqualValues(1, item.Items[0].BlobsCount)
		}
	}
}

func (s *RollupTestSuite) TestRollupStatsGrouping() {
	for _, funcName := range []string{
		"sum",
		"avg",
	} {
		for _, groupName := range []string{
			"stack",
			"type",
			"category",
			"vm",
			"provider",
		} {
			q := make(url.Values)
			q.Add("func", funcName)
			q.Add("column", groupName)

			req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/rollup/group")

			s.rollups.EXPECT().
				RollupStatsGrouping(gomock.Any(), storage.RollupGroupStatsFilters{
					Func:   funcName,
					Column: groupName,
				}).
				Return([]storage.RollupGroupedStats{testRollupWithGroupedStats}, nil).
				Times(1)

			s.Require().NoError(s.handler.RollupGroupedStats(c))
			s.Require().Equal(http.StatusOK, rec.Code)
			var stats []responses.RollupGroupedStats
			err := json.NewDecoder(rec.Body).Decode(&stats)
			s.Require().NoError(err)
			s.Require().Len(stats, 1)

			groupedStats := stats[0]

			s.Require().EqualValues(0.1, groupedStats.Fee)
			s.Require().EqualValues(0.2, groupedStats.Size)
			s.Require().EqualValues(3, groupedStats.BlobsCount)
			s.Require().EqualValues("stack", groupedStats.Group)
		}
	}
}
