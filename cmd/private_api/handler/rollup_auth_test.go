// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// AuthTestSuite -
type AuthTestSuite struct {
	suite.Suite
	address   *mock.MockIAddress
	namespace *mock.MockINamespace
	rollups   *mock.MockIRollup
	echo      *echo.Echo
	ctrl      *gomock.Controller
}

// SetupSuite -
func (s *AuthTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.address = mock.NewMockIAddress(s.ctrl)
	s.namespace = mock.NewMockINamespace(s.ctrl)
	s.rollups = mock.NewMockIRollup(s.ctrl)
}

// TearDownSuite -
func (s *AuthTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteAuth_Run(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) TestBulk() {
	body := `
	{
		"rollups": [
			{
				"id": 1,
				"vm": "evm",
				"providers": [{
					"address": "celestia1kywuhlvslyt0qy8yr4p5lgkzz74qryujkjgprx"
				}]
			}, {
				"name": "Test",
				"providers": [{
					"address": "celestia1kywuhlvslyt0qy8yr4p5lgkzz74qryujkjgprx"
				}],
				"vm": "svm"
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.Set(ApiKeyName, storage.ApiKey{
		Admin:       true,
		Description: "test",
		Key:         "test",
	})
	c.SetPath("/v1/bulk")

	txUpdate := mock.NewMockTransaction(s.ctrl)
	txCreate := mock.NewMockTransaction(s.ctrl)

	// first rollup: update existing
	txUpdate.EXPECT().
		DeleteProviders(gomock.Any(), uint64(1)).
		Return(nil).
		Times(1)

	txUpdate.EXPECT().
		SaveProviders(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	txUpdate.EXPECT().
		UpdateRollup(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	txUpdate.EXPECT().
		Flush(gomock.Any()).
		Return(nil).
		Times(1)

	// second rollup: create new
	txCreate.EXPECT().
		SaveRollup(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, r *storage.Rollup) error {
			r.Id = 2
			return nil
		}).
		Times(1)

	txCreate.EXPECT().
		SaveProviders(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	txCreate.EXPECT().
		Flush(gomock.Any()).
		Return(nil).
		Times(1)

	s.rollups.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&storage.Rollup{
			Id:   1,
			Name: "First",
		}, nil).
		Times(1)

	s.address.EXPECT().
		ByHash(gomock.Any(), []byte{177, 29, 203, 253, 144, 249, 22, 240, 16, 228, 29, 67, 79, 162, 194, 23, 170, 1, 147, 146}).
		Return(storage.Address{
			Id:      100,
			Address: "celestia1kywuhlvslyt0qy8yr4p5lgkzz74qryujkjgprx",
		}, nil).
		Times(2)

	callCount := 0
	txBeginner := func(_ context.Context, _ sdk.Transactable) (storage.Transaction, error) {
		callCount++
		if callCount == 1 {
			return txUpdate, nil
		}
		return txCreate, nil
	}
	handler := NewRollupAuthHandler(s.rollups, s.address, s.namespace, nil, txBeginner)

	s.Require().NoError(handler.Bulk(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []bulkResultItem
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 2)
	s.Require().EqualValues(1, results[0].Id)
	s.Require().Empty(results[0].Error)
	s.Require().EqualValues(2, results[1].Id)
	s.Require().Empty(results[1].Error)
}

func (s *AuthTestSuite) TestBulkPartialError() {
	body := `
	{
		"rollups": [
			{
				"name": "Good Rollup",
				"description": "works fine",
				"logo": "https://example.com/logo.png",
				"providers": [{
					"address": "celestia1kywuhlvslyt0qy8yr4p5lgkzz74qryujkjgprx"
				}],
				"vm": "svm"
			},
			{
				"name": "Bad Rollup",
				"description": "will fail",
				"logo": "https://example.com/logo.png",
				"providers": [{
					"address": "celestia1kywuhlvslyt0qy8yr4p5lgkzz74qryujkjgprx"
				}],
				"vm": "evm"
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.Set(ApiKeyName, storage.ApiKey{
		Admin:       true,
		Description: "test",
		Key:         "test",
	})
	c.SetPath("/v1/bulk")

	txSuccess := mock.NewMockTransaction(s.ctrl)
	txFail := mock.NewMockTransaction(s.ctrl)

	// first rollup succeeds
	txSuccess.EXPECT().
		SaveRollup(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, r *storage.Rollup) error {
			r.Id = 1
			return nil
		}).
		Times(1)

	txSuccess.EXPECT().
		SaveProviders(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	txSuccess.EXPECT().
		Flush(gomock.Any()).
		Return(nil).
		Times(1)

	// second rollup fails on SaveRollup
	txFail.EXPECT().
		SaveRollup(gomock.Any(), gomock.Any()).
		Return(errors.New("duplicate slug")).
		Times(1)

	txFail.EXPECT().
		HandleError(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	s.address.EXPECT().
		ByHash(gomock.Any(), []byte{177, 29, 203, 253, 144, 249, 22, 240, 16, 228, 29, 67, 79, 162, 194, 23, 170, 1, 147, 146}).
		Return(storage.Address{
			Id:      100,
			Address: "celestia1kywuhlvslyt0qy8yr4p5lgkzz74qryujkjgprx",
		}, nil).
		Times(1)

	callCount := 0
	txBeginner := func(_ context.Context, _ sdk.Transactable) (storage.Transaction, error) {
		callCount++
		if callCount == 1 {
			return txSuccess, nil
		}
		return txFail, nil
	}
	handler := NewRollupAuthHandler(s.rollups, s.address, s.namespace, nil, txBeginner)

	s.Require().NoError(handler.Bulk(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []bulkResultItem
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 2)
	s.Require().EqualValues(1, results[0].Id)
	s.Require().Empty(results[0].Error)
	s.Require().Zero(results[1].Id)
	s.Require().Contains(results[1].Error, "duplicate slug")
}
