// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// ValidatorTestSuite -
type ValidatorTestSuite struct {
	suite.Suite
	validators *mock.MockIValidator
	echo       *echo.Echo
	handler    *ValidatorHandler
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *ValidatorTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.validators = mock.NewMockIValidator(s.ctrl)
	s.handler = NewValidatorHandler(s.validators)
}

// TearDownSuite -
func (s *ValidatorTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteValidator_Run(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}

func (s *ValidatorTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/validators/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.validators.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&testValidator, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var validator responses.Validator
	err := json.NewDecoder(rec.Body).Decode(&validator)
	s.Require().NoError(err)
	s.Require().EqualValues(1, validator.Id)
	s.Require().EqualValues("moniker", validator.Moniker)
	s.Require().EqualValues("012345", validator.ConsAddress)
}

func (s *ValidatorTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address")

	s.validators.EXPECT().
		List(gomock.Any(), uint64(10), uint64(0), pgSort("asc")).
		Return([]*storage.Validator{
			&testValidator,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var validators []responses.Validator
	err := json.NewDecoder(rec.Body).Decode(&validators)
	s.Require().NoError(err)
	s.Require().Len(validators, 1)
	s.Require().EqualValues(1, validators[0].Id)
	s.Require().EqualValues("moniker", validators[0].Moniker)
	s.Require().EqualValues("012345", validators[0].ConsAddress)
}
