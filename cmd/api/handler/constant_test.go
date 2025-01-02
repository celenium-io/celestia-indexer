// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// ConstantTestSuite -
type ConstantTestSuite struct {
	suite.Suite
	constants     *mock.MockIConstant
	denomMetadata *mock.MockIDenomMetadata
	address       *mock.MockIAddress
	echo          *echo.Echo
	handler       *ConstantHandler
	ctrl          *gomock.Controller
}

// SetupSuite -
func (s *ConstantTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.constants = mock.NewMockIConstant(s.ctrl)
	s.denomMetadata = mock.NewMockIDenomMetadata(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
}

// TearDownSuite -
func (s *ConstantTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteConstant_Run(t *testing.T) {
	suite.Run(t, new(ConstantTestSuite))
}

func (s *ConstantTestSuite) TestEnums() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/enums")

	s.Require().NoError(s.handler.Enums(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var enums responses.Enums
	err := json.NewDecoder(rec.Body).Decode(&enums)
	s.Require().NoError(err)
	s.Require().Len(enums.EventType, 59)
	s.Require().Len(enums.MessageType, 76)
	s.Require().Len(enums.Status, 2)
}
