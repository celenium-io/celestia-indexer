package handler

import (
	"context"
	"encoding/json"
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

	tx := mock.NewMockTransaction(s.ctrl)

	tx.EXPECT().
		Flush(gomock.Any()).
		Return(nil).
		Times(1)

	tx.EXPECT().
		SaveProviders(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(2)

	tx.EXPECT().
		DeleteProviders(gomock.Any(), uint64(1)).
		Return(nil).
		Times(1)

	tx.EXPECT().
		SaveRollup(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, r *storage.Rollup) error {
			r.Id = 2
			return nil
		}).
		Times(1)

	tx.EXPECT().
		UpdateRollup(gomock.Any(), gomock.Any()).
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

	txBeginner := func(_ context.Context, _ sdk.Transactable) (storage.Transaction, error) {
		return tx, nil
	}
	handler := NewRollupAuthHandler(s.rollups, s.address, s.namespace, nil, txBeginner)

	s.Require().NoError(handler.Bulk(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var ids struct {
		Values []uint64 `json:"ids"`
	}
	err := json.NewDecoder(rec.Body).Decode(&ids)
	s.Require().NoError(err)
	s.Require().Len(ids.Values, 2)
	s.Require().EqualValues(1, ids.Values[0])
	s.Require().EqualValues(2, ids.Values[1])
}
