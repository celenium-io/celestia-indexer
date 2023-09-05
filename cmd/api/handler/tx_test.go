package handler

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testTxHash         = "652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF"
	testTxHashBytes, _ = hex.DecodeString(testTxHash)
	testTx             = storage.Tx{
		Id:            1,
		Hash:          testTxHashBytes,
		Height:        100,
		Time:          testTime,
		Position:      2,
		GasWanted:     80410,
		GasUsed:       77483,
		TimeoutHeight: 0,
		EventsCount:   10,
		MessagesCount: 2,
		Fee:           decimal.RequireFromString("80410"),
		Status:        types.StatusSuccess,
		Codespace:     "sdk",
		Memo:          "memo",
	}
)

// TxTestSuite -
type TxTestSuite struct {
	suite.Suite
	tx       *mock.MockITx
	events   *mock.MockIEvent
	messages *mock.MockIMessage
	echo     *echo.Echo
	handler  *TxHandler
	ctrl     *gomock.Controller
}

// SetupSuite -
func (s *TxTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.tx = mock.NewMockITx(s.ctrl)
	s.events = mock.NewMockIEvent(s.ctrl)
	s.messages = mock.NewMockIMessage(s.ctrl)
	s.handler = NewTxHandler(s.tx, s.events, s.messages)
}

// TearDownSuite -
func (s *TxTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteTx_Run(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}

func (s *TxTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/:hash")
	c.SetParamNames("hash")
	c.SetParamValues(testTxHash)

	s.tx.EXPECT().
		ByHash(gomock.Any(), testTxHashBytes).
		Return(testTx, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var tx responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&tx)
	s.Require().NoError(err)
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, strings.ToUpper(tx.Hash))
	s.Require().EqualValues(2, tx.Position)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().Equal("80410", tx.Fee)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(10, tx.EventsCount)
	s.Require().EqualValues(2, tx.MessagesCount)
	s.Require().Equal("memo", tx.Memo)
	s.Require().Equal("sdk", tx.Codespace)
	s.Require().Equal(string(types.StatusSuccess), tx.Status)
}

func (s *TxTestSuite) TestGetInvalidTx() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/:hash")
	c.SetParamNames("hash")
	c.SetParamValues("invalid")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *TxTestSuite) TestList() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("msg_type", "MsgUnjail,MsgSend")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.tx.EXPECT().
		Filter(gomock.Any(), gomock.Any()).
		Return([]storage.Tx{
			testTx,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&txs)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, strings.ToUpper(tx.Hash))
	s.Require().EqualValues(2, tx.Position)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().Equal("80410", tx.Fee)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(10, tx.EventsCount)
	s.Require().EqualValues(2, tx.MessagesCount)
	s.Require().Equal("memo", tx.Memo)
	s.Require().Equal("sdk", tx.Codespace)
	s.Require().Equal(string(types.StatusSuccess), tx.Status)
}

func (s *TxTestSuite) TestListValidationStatusError() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "invalid")
	q.Set("msg_type", "MsgUnjail,MsgSend")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *TxTestSuite) TestListValidationMsgTypeError() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("msg_type", "invalid,MsgSend")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *TxTestSuite) TestListTime() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("msg_type", "MsgSend")
	q.Set("from", "1692880000")
	q.Set("to", "1692890000")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.tx.EXPECT().
		Filter(gomock.Any(), gomock.Any()).
		Return([]storage.Tx{
			testTx,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&txs)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, strings.ToUpper(tx.Hash))
	s.Require().EqualValues(2, tx.Position)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().Equal("80410", tx.Fee)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(10, tx.EventsCount)
	s.Require().EqualValues(2, tx.MessagesCount)
	s.Require().Equal("memo", tx.Memo)
	s.Require().Equal("sdk", tx.Codespace)
	s.Require().Equal(string(types.StatusSuccess), tx.Status)
}

func (s *TxTestSuite) TestListHeight() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("msg_type", "MsgSend")
	q.Set("height", "1000")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.tx.EXPECT().
		Filter(gomock.Any(), gomock.Any()).
		Return([]storage.Tx{
			testTx,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&txs)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, strings.ToUpper(tx.Hash))
	s.Require().EqualValues(2, tx.Position)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().Equal("80410", tx.Fee)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(10, tx.EventsCount)
	s.Require().EqualValues(2, tx.MessagesCount)
	s.Require().Equal("memo", tx.Memo)
	s.Require().Equal("sdk", tx.Codespace)
	s.Require().Equal(string(types.StatusSuccess), tx.Status)
}

func (s *TxTestSuite) TestGetEvents() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/:hash/events")
	c.SetParamNames("hash")
	c.SetParamValues(testTxHash)

	s.tx.EXPECT().
		ByHash(gomock.Any(), testTxHashBytes).
		Return(testTx, nil)

	s.events.EXPECT().
		ByTxId(gomock.Any(), uint64(1)).
		Return([]storage.Event{
			{
				Id:       1,
				Height:   100,
				Time:     testTime,
				Position: 2,
				Type:     types.EventTypeBurn,
				TxId:     testsuite.Ptr(uint64(1)),
				Data: map[string]any{
					"test": "value",
				},
			},
		}, nil)

	s.Require().NoError(s.handler.GetEvents(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var events []responses.Event
	err := json.NewDecoder(rec.Body).Decode(&events)
	s.Require().NoError(err)
	s.Require().Len(events, 1)
	s.Require().EqualValues(1, events[0].Id)
	s.Require().EqualValues(100, events[0].Height)
	s.Require().EqualValues(2, events[0].Position)
	s.Require().Equal(testTime, events[0].Time)
	s.Require().EqualValues(1, events[0].TxId)
	s.Require().EqualValues(string(types.EventTypeBurn), events[0].Type)
}

func (s *TxTestSuite) TestGetMessage() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/:hash/events")
	c.SetParamNames("hash")
	c.SetParamValues(testTxHash)

	s.tx.EXPECT().
		ByHash(gomock.Any(), testTxHashBytes).
		Return(testTx, nil)

	s.messages.EXPECT().
		ByTxId(gomock.Any(), uint64(1)).
		Return([]storage.Message{
			{
				Id:       1,
				Height:   100,
				Time:     testTime,
				Position: 2,
				Type:     types.MsgTypeBeginRedelegate,
				TxId:     1,
				Data: map[string]any{
					"test": "value",
				},
			},
		}, nil)

	s.Require().NoError(s.handler.GetMessages(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var msgs []responses.Message
	err := json.NewDecoder(rec.Body).Decode(&msgs)
	s.Require().NoError(err)
	s.Require().Len(msgs, 1)
	s.Require().EqualValues(1, msgs[0].Id)
	s.Require().EqualValues(100, msgs[0].Height)
	s.Require().EqualValues(2, msgs[0].Position)
	s.Require().Equal(testTime, msgs[0].Time)
	s.Require().EqualValues(1, msgs[0].TxId)
	s.Require().EqualValues(string(types.MsgTypeBeginRedelegate), msgs[0].Type)
}
