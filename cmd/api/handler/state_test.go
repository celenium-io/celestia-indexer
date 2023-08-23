package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// StateTestSuite -
type StateTestSuite struct {
	suite.Suite
	state   *mock.MockIState
	echo    *echo.Echo
	handler *StateHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *StateTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewStateHandler(s.state)
}

// TearDownSuite -
func (s *StateTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteState_Run(t *testing.T) {
	suite.Run(t, new(StateTestSuite))
}

func (s *StateTestSuite) TestHead() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/head")

	s.state.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.State{
			{
				ID:                 1,
				Name:               "test",
				LastHeight:         100,
				LastTime:           testTime,
				TotalTx:            1234,
				TotalAccounts:      123,
				TotalNamespaces:    2,
				TotalNamespaceSize: 30,
			},
		}, nil)

	s.Require().NoError(s.handler.Head(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var state responses.State
	err := json.NewDecoder(rec.Body).Decode(&state)
	s.Require().NoError(err)
	s.Require().EqualValues(1, state.Id)
	s.Require().EqualValues("test", state.Name)
	s.Require().EqualValues(100, state.LastHeight)
	s.Require().EqualValues(1234, state.TotalTx)
	s.Require().EqualValues(123, state.TotalAccounts)
	s.Require().EqualValues(2, state.TotalNamespaces)
	s.Require().EqualValues(30, state.TotalNamespaceSize)
	s.Require().Equal(testTime, state.LastTime)
}

func (s *StateTestSuite) TestNoHead() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/head")

	s.state.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.State{}, nil)

	s.Require().NoError(s.handler.Head(c))
	s.Require().Equal(http.StatusNoContent, rec.Code)
}
