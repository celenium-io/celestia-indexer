// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	nodeMock "github.com/celenium-io/celestia-indexer/pkg/node/mock"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	tmTypes "github.com/cometbft/cometbft/types"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testBlock = storage.Block{
		Id:           1,
		Hash:         []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		Height:       100,
		VersionBlock: 11,
		VersionApp:   1,
		Time:         testTime,
		MessageTypes: types.NewMsgTypeBitMask(types.MsgSend),
	}
	testBlockStats = storage.BlockStats{
		TxCount:     1,
		EventsCount: 2,
		Time:        testTime,
		Height:      100,
		BlockTime:   11043,
	}
	testValidator = storage.Validator{
		Id:          1,
		Moniker:     "moniker",
		ConsAddress: "012345",
		Jailed:      testsuite.Ptr(false),
	}
	testBlockWithStats = storage.Block{
		Id:           1,
		Hash:         []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		Height:       100,
		VersionBlock: 11,
		VersionApp:   1,
		Time:         testTime,
		MessageTypes: types.NewMsgTypeBitMask(types.MsgSend),
		Stats:        testBlockStats,
		Proposer:     testValidator,
	}

	testTime = time.Date(2023, 8, 1, 1, 1, 0, 0, time.UTC)
)

// BlockTestSuite -
type BlockTestSuite struct {
	suite.Suite
	blocks     *mock.MockIBlock
	blockStats *mock.MockIBlockStats
	events     *mock.MockIEvent
	message    *mock.MockIMessage
	namespace  *mock.MockINamespace
	blobLogs   *mock.MockIBlobLog
	state      *mock.MockIState
	node       *nodeMock.MockApi
	echo       *echo.Echo
	handler    *BlockHandler
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *BlockTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.blocks = mock.NewMockIBlock(s.ctrl)
	s.blockStats = mock.NewMockIBlockStats(s.ctrl)
	s.events = mock.NewMockIEvent(s.ctrl)
	s.namespace = mock.NewMockINamespace(s.ctrl)
	s.blobLogs = mock.NewMockIBlobLog(s.ctrl)
	s.message = mock.NewMockIMessage(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.node = nodeMock.NewMockApi(s.ctrl)
	s.handler = NewBlockHandler(s.blocks, s.blockStats, s.events, s.namespace, s.message, s.blobLogs, s.state, s.node, testIndexerName)
}

// TearDownSuite -
func (s *BlockTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteBlock_Run(t *testing.T) {
	suite.Run(t, new(BlockTestSuite))
}

func (s *BlockTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100)).
		Return(testBlock, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var block responses.Block
	err := json.NewDecoder(rec.Body).Decode(&block)
	s.Require().NoError(err)
	s.Require().EqualValues(1, block.Id)
	s.Require().EqualValues(100, block.Height)
	s.Require().Equal("1", block.VersionApp)
	s.Require().Equal("11", block.VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", block.Hash.String())
	s.Require().Equal(testTime, block.Time)
	s.Require().Equal([]types.MsgType{types.MsgSend}, block.MessageTypes)
	s.Require().Nil(block.Stats)
}

func (s *BlockTestSuite) TestGetNoContent() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100)).
		Return(storage.Block{}, sql.ErrNoRows)

	s.blocks.EXPECT().
		IsNoRows(gomock.Any()).
		Return(true)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusNoContent, rec.Code)
}

func (s *BlockTestSuite) TestGetWithoutStats() {
	q := make(url.Values)
	q.Set("stats", "false")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100)).
		Return(testBlock, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var block responses.Block
	err := json.NewDecoder(rec.Body).Decode(&block)
	s.Require().NoError(err)
	s.Require().EqualValues(1, block.Id)
	s.Require().EqualValues(100, block.Height)
	s.Require().Equal("1", block.VersionApp)
	s.Require().Equal("11", block.VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", block.Hash.String())
	s.Require().Equal(testTime, block.Time)
	s.Require().Equal([]types.MsgType{types.MsgSend}, block.MessageTypes)
	s.Require().Nil(block.Stats)
}

func (s *BlockTestSuite) TestGetWithStats() {
	q := make(url.Values)
	q.Set("stats", "true")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		ByHeightWithStats(gomock.Any(), pkgTypes.Level(100)).
		Return(testBlockWithStats, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var block responses.Block
	err := json.NewDecoder(rec.Body).Decode(&block)
	s.Require().NoError(err)
	s.Require().EqualValues(1, block.Id)
	s.Require().EqualValues(100, block.Height)
	s.Require().Equal("1", block.VersionApp)
	s.Require().Equal("11", block.VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", block.Hash.String())
	s.Require().Equal(testTime, block.Time)
	s.Require().Equal([]types.MsgType{types.MsgSend}, block.MessageTypes)
	s.Require().NotNil(block.Stats)
	s.Require().EqualValues(1, block.Stats.TxCount)
	s.Require().EqualValues(2, block.Stats.EventsCount)
}

func (s *BlockTestSuite) TestGetInvalidBlockHeight() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("invalid")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "parsing")
}

func (s *BlockTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block")

	s.blocks.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.Block{
			&testBlock,
		}, nil).
		MaxTimes(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blocks []responses.Block
	err := json.NewDecoder(rec.Body).Decode(&blocks)
	s.Require().NoError(err)
	s.Require().Len(blocks, 1)
	s.Require().EqualValues(1, blocks[0].Id)
	s.Require().EqualValues(100, blocks[0].Height)
	s.Require().Equal("1", blocks[0].VersionApp)
	s.Require().Equal("11", blocks[0].VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", blocks[0].Hash.String())
	s.Require().Equal(testTime, blocks[0].Time)
	s.Require().Equal([]types.MsgType{types.MsgSend}, blocks[0].MessageTypes)
}

func (s *BlockTestSuite) TestListWithStats() {
	q := make(url.Values)
	q.Set("stats", "true")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block")

	s.blocks.EXPECT().
		ListWithStats(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.Block{
			&testBlockWithStats,
		}, nil).
		MaxTimes(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blocks []responses.Block
	err := json.NewDecoder(rec.Body).Decode(&blocks)
	s.Require().NoError(err)
	s.Require().Len(blocks, 1)
	s.Require().EqualValues(1, blocks[0].Id)
	s.Require().EqualValues(100, blocks[0].Height)
	s.Require().Equal("1", blocks[0].VersionApp)
	s.Require().Equal("11", blocks[0].VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", blocks[0].Hash.String())
	s.Require().Equal(testTime, blocks[0].Time)
	s.Require().Equal([]types.MsgType{types.MsgSend}, blocks[0].MessageTypes)
}

func (s *BlockTestSuite) TestGetEvents() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/events")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		Time(gomock.Any(), pkgTypes.Level(100)).
		Return(testTime, nil).
		Times(1)

	s.events.EXPECT().
		ByBlock(gomock.Any(), pkgTypes.Level(100), gomock.Any()).
		Return([]storage.Event{
			{
				Id:       1,
				Height:   100,
				Time:     testTime,
				Position: 2,
				Type:     types.EventTypeBurn,
				TxId:     nil,
				Data: map[string]any{
					"test": "value",
				},
			},
		}, nil).
		Times(1)

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
	s.Require().Equal(types.EventTypeBurn, events[0].Type)
}

func (s *BlockTestSuite) TestGetStats() {
	req := httptest.NewRequest(http.MethodGet, "/?", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/stats")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blockStats.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100)).
		Return(testBlockStats, nil)

	s.Require().NoError(s.handler.GetStats(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var stats responses.BlockStats
	err := json.NewDecoder(rec.Body).Decode(&stats)
	s.Require().NoError(err)
	s.Require().EqualValues(1, stats.TxCount)
	s.Require().EqualValues(2, stats.EventsCount)
	s.Require().EqualValues(11043, stats.BlockTime)
}

func (s *BlockTestSuite) TestBlobs() {
	req := httptest.NewRequest(http.MethodGet, "/?", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/blobs")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		Time(gomock.Any(), pkgTypes.Level(100)).
		Return(testBlock.Time, nil).
		Times(1)

	s.blobLogs.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100), gomock.Any()).
		Return([]storage.BlobLog{
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

	s.Require().NoError(s.handler.Blobs(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var logs []responses.BlobLog
	err := json.NewDecoder(rec.Body).Decode(&logs)
	s.Require().NoError(err)
	s.Require().Len(logs, 1)
}

func (s *BlockTestSuite) TestGetBlobsCount() {
	req := httptest.NewRequest(http.MethodGet, "/?", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/blobs/count")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blockStats.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100)).
		Return(storage.BlockStats{
			BlobsCount: 12,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.BlobsCount(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count int
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(count, 12)
}

func (s *BlockTestSuite) TestCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/count")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(testState, nil)

	s.Require().NoError(s.handler.Count(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count uint64
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(80001, count)
}

func (s *BlockTestSuite) TestGetMessages() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("msg_type", "MsgSend")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/messages")
	c.SetParamNames("height")
	c.SetParamValues("1000")

	s.message.EXPECT().
		ListWithTx(gomock.Any(), storage.MessageListWithTxFilters{
			Limit:                10,
			Offset:               0,
			MessageTypes:         []string{"MsgSend"},
			ExcludedMessageTypes: nil,
			Height:               1000,
		}).
		Return([]storage.MessageWithTx{
			{
				Message: storage.Message{
					Id:       1,
					TxId:     2,
					Position: 3,
					Type:     types.MsgSend,
					Height:   100,
					Time:     testTime,
				},
				Tx: &testTx,
			},
		}, nil)

	s.Require().NoError(s.handler.GetMessages(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var msgs []responses.Message
	err := json.NewDecoder(rec.Body).Decode(&msgs)
	s.Require().NoError(err)
	s.Require().Len(msgs, 1)

	msg := msgs[0]
	s.Require().EqualValues(1, msg.Id)
	s.Require().EqualValues(100, msg.Height)
	s.Require().EqualValues(3, msg.Position)
	s.Require().Equal(testTime, msg.Time)
	s.Require().EqualValues(string(types.MsgSend), msg.Type)
	s.Require().EqualValues(1, msg.Tx.Id)
	s.Require().NotNil(msg.Tx)

	tx := msg.Tx
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
	s.Require().Equal(types.StatusSuccess, tx.Status)
}

func (s *BlockTestSuite) TestBlockODS() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/ods")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blockStats.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100)).
		Return(testBlockStats, nil).
		Times(1)

	rawTxs := []string{
		"CpkBCpYBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnYKL2NlbGVzdGlhMW5lNGhsMHZqczJmcjVtbXprYWF2YXFlOWtrODdlbmU5YXp5bjVsEi9jZWxlc3RpYTFram15N3Zoeno5NTRycWUyOTNmeGdjNjN1anF5eTBubmUwemtrahoSCgR1dGlhEgoxOTQzMTMwNDI0EmgKUQpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQOyjTDvB2Ur253W9bSWnQjNXUSEKQLqXmGJAIbEsGlBmxIECgIIARjNExITCg0KBHV0aWESBTI1MDAwEJChDxpAImupFuuJjoWgZXJvdIvRhuIPUcu4d4NfScVAnlH2974PIOnVljFM4rmYcSQt3U9IHDxg1HN4cF5h6lZrXVURoQ==",
		"CpcBCpQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnQKL2NlbGVzdGlhMWR6bnN6dDVrMjBqODdwa3VmYXdydXpyMnJ6a3JqZ3hsY210cjM5Ei9jZWxlc3RpYTF0d2dwYzdjbXg0dGczeXZtOXFubDB2YWs3a3l4MG1leDA0enR2YRoQCgR1dGlhEgg0MTUwMDAwMBJkCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECOYBKFkZUA9pt+sASM0MZmm9wCO00o9MnPz46qcJKjd0SBAoCCH8SEgoMCgR1dGlhEgQ5MDAwEJC/BRpAdBYZN/811Z5V/uKRC5MLV83TF/vcWxKGLVvNbA3CG04M7PLM5xJn/YLqWkIgEM6k8XUvd/4Uf2z4T6k+M/5MQw==",
		"Cs0CCqABCp0BCiAvY2VsZXN0aWEuYmxvYi52MS5Nc2dQYXlGb3JCbG9icxJ5Ci9jZWxlc3RpYTF1bTlxdDB6dWVsMGY0a3JlenZ1ejRrYTQ0M2syaHc1cGd4NDJmbhIdAAAAAAAAAAAAAAAAAAAAAAAAAAAAY2NsYWJzMDEaAt4EIiABYhBkLI4ANqUrn22XSeofPm3//nJBIe4RSOk/3tm3RUIBABJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDKvAge1/isZG9nUxUjNWtsjHf0t3EAIwrXEONbjGBiVYSBAoCCAEYBRISCgwKBHV0aWESBDgzOTAQtI8FGkCaSY/c4H0h2HOeNv7C5Uu8YiuRShagt02sn8DDgSQ2qiVEQhZ+/cyh101cU8B+fFvojCjJ41gFIqC5gmplAg27Ev8EChwAAAAAAAAAAAAAAAAAAAAAAAAAAGNjbGFiczAxEt4EewogICJuYW1lIjogIjUzMjk3ODY1MTQ1XzJhZTk1ODg4MzFfNmsuanBnIiwKICAibWltZVR5cGUiOiAiaW1hZ2UvanBlZyIsCiAgInNpemUiOiA0NjIxMDk2LAogICJjaHVua3MiOiBbCiAgICB7ICJibG9iIjogIjE4OTcyLzYzNjM2YzYxNjI3MzMwMzEvNzcwMmU5ZmJjYTFlYmQxYjViMDkzZTZiOTcyMTZhMDAyYTkyNjhiN2ZiODlhNzg0ODgxZWU4NTEwZGI2NTM1MiIsICJzaXplIjogMTUwMDAwMCB9LAogICAgeyAiYmxvYiI6ICIxODk3My82MzYzNmM2MTYyNzMzMDMxLzUzMDBhOWJkZTdiZGE4MTgwOGYxZjAyNTkyYzRlZDUxMWMwN2EyYTlhNjQzM2I5Zjc4YmM4OTMxM2NkOTU3YmUiLCAic2l6ZSI6IDE1MDAwMDAgfSwKICAgIHsgImJsb2IiOiAiMTg5NzQvNjM2MzZjNjE2MjczMzAzMS9lYTA0NmIxZmQ0ZGU2OTI3Yzg0MWViNjA3YjFmZjhmZGQ2ZWJlYTQxZWM0ZTk4Zjk4NTIwMjkzNTIyN2M2YmM4IiwgInNpemUiOiAxNTAwMDAwIH0sCiAgICB7ICJibG9iIjogIjE4OTc1LzYzNjM2YzYxNjI3MzMwMzEvODUwZjg5Y2M1MmM5ODI4NzM2MTRhYWQzYmQyOGM2MTEzNmE1MWY3NTAwNjcwZGEzODhiOTEzMzk0YzQwODM3OCIsICJzaXplIjogMTIxMDk2IH0KICBdCn0KGgRCTE9C",
	}

	txs := make(tmTypes.Txs, 0)
	for i := range rawTxs {
		t, err := base64.StdEncoding.DecodeString(rawTxs[i])
		s.Require().NoError(err)
		txs = append(txs, t)
	}

	s.node.EXPECT().
		Block(gomock.Any(), pkgTypes.Level(100)).
		Return(pkgTypes.ResultBlock{
			Block: &pkgTypes.Block{
				Data: pkgTypes.Data{
					Txs: txs,
				},
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.BlockODS(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var ods responses.ODS
	err := json.NewDecoder(rec.Body).Decode(&ods)
	s.Require().NoError(err)

	s.Require().EqualValues(4, ods.Width)
	s.Require().Len(ods.Items, 4)
}

func (s *BlockTestSuite) TestEmptyBlockODS() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/ods")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blockStats.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100)).
		Return(storage.BlockStats{}, nil).
		Times(1)

	s.Require().NoError(s.handler.BlockODS(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var ods responses.ODS
	err := json.NewDecoder(rec.Body).Decode(&ods)
	s.Require().NoError(err)

	s.Require().EqualValues(1, ods.Width)
	s.Require().Len(ods.Items, 1)
}
