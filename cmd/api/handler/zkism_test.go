// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

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

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// ──────────────────────────────────────────────────────────
// Test fixtures
// ──────────────────────────────────────────────────────────

var (
	testZkISMStateRoot           = testsuite.RandomBytes(32)
	testZkISMMerkleTreeAddress   = testsuite.RandomBytes(32)
	testZkISMStateTransitionVKey = testsuite.RandomBytes(32)
	testZkISMStateMembershipVKey = testsuite.RandomBytes(32)

	testZkISM = storage.ZkISM{
		Id:                  1,
		ExternalId:          []byte{0x42},
		Height:              1000,
		Time:                testTime,
		TxId:                1,
		CreatorId:           1,
		StateRoot:           testZkISMStateRoot,
		MerkleTreeAddress:   testZkISMMerkleTreeAddress,
		StateTransitionVKey: testZkISMStateTransitionVKey,
		StateMembershipVKey: testZkISMStateMembershipVKey,
		Creator: &storage.Address{
			Address: testAddress,
			Hash:    testHashAddress,
		},
		Tx: &testTx,
	}

	testZkISMUpdateNewStateRoot = testsuite.RandomBytes(32)
	testZkISMUpdate             = storage.ZkISMUpdate{
		Id:           10,
		Height:       1001,
		Time:         testTime,
		ZkISMId:      1,
		NewStateRoot: testZkISMUpdateNewStateRoot,
		Signer: &storage.Address{
			Address: testAddress,
			Hash:    testHashAddress,
		},
		Tx: &testTx,
	}

	testZkISMMessageStateRoot = testsuite.RandomBytes(32)
	testZkISMMessageId        = testsuite.RandomBytes(32)
	testZkISMMsg              = storage.ZkISMMessage{
		Id:        20,
		Height:    1002,
		Time:      testTime,
		ZkISMId:   1,
		StateRoot: testZkISMMessageStateRoot,
		MessageId: testZkISMMessageId,
		Signer: &storage.Address{
			Address: testAddress,
			Hash:    testHashAddress,
		},
		Tx: &testTx,
	}
)

// ──────────────────────────────────────────────────────────
// Test suite
// ──────────────────────────────────────────────────────────

type ZkISMTestSuite struct {
	suite.Suite
	echo    *echo.Echo
	zkism   *mock.MockIZkISM
	address *mock.MockIAddress
	txs     *mock.MockITx
	handler *ZkISMHandler
	ctrl    *gomock.Controller
}

func (s *ZkISMTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.zkism = mock.NewMockIZkISM(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.handler = NewZkISMHandler(s.zkism, s.address, s.txs)
}

func (s *ZkISMTestSuite) TearDownSuite() {
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteZkISM_Run(t *testing.T) {
	suite.Run(t, new(ZkISMTestSuite))
}

// ──────────────────────────────────────────────────────────
// List
// ──────────────────────────────────────────────────────────

func (s *ZkISMTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism")

	s.zkism.EXPECT().
		List(gomock.Any(), storage.ZkISMFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
		}).
		Return([]storage.ZkISM{testZkISM}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISM
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 1)

	ism := response[0]
	s.Require().EqualValues(testZkISM.Id, ism.Id)
	s.Require().EqualValues(hex.EncodeToString(testZkISM.ExternalId), ism.ExternalId)
	s.Require().EqualValues(testZkISM.Height, ism.Height)
	s.Require().Equal(testZkISM.Time, ism.Time)
	s.Require().Equal(hex.EncodeToString(testZkISMStateRoot), ism.StateRoot)
	s.Require().Equal(hex.EncodeToString(testZkISMMerkleTreeAddress), ism.MerkleTreeAddress)
	s.Require().Equal(hex.EncodeToString(testZkISMStateTransitionVKey), ism.StateTransitionVKey)
	s.Require().Equal(hex.EncodeToString(testZkISMStateMembershipVKey), ism.StateMembershipVKey)
	s.Require().Equal(strings.ToLower(testTxHash), ism.TxHash)
	s.Require().NotNil(ism.Creator)
	s.Require().Equal(testAddress, ism.Creator.Hash)
}

func (s *ZkISMTestSuite) TestListEmpty() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism")

	s.zkism.EXPECT().
		List(gomock.Any(), storage.ZkISMFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
		}).
		Return([]storage.ZkISM{}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISM
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 0)
}

func (s *ZkISMTestSuite) TestListWithLimitOffset() {
	q := make(url.Values)
	q.Set("limit", "5")
	q.Set("offset", "10")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism")

	s.zkism.EXPECT().
		List(gomock.Any(), storage.ZkISMFilter{
			Limit:  5,
			Offset: 10,
			Sort:   pgSort("desc"),
		}).
		Return([]storage.ZkISM{testZkISM}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISM
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 1)
}

func (s *ZkISMTestSuite) TestListWithAddress() {
	q := make(url.Values)
	q.Set("address", testAddress)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism")

	addressId := uint64(1)
	s.address.EXPECT().
		IdByAddress(gomock.Any(), testAddress).
		Return(addressId, nil).
		Times(1)

	s.zkism.EXPECT().
		List(gomock.Any(), storage.ZkISMFilter{
			Limit:     10,
			Offset:    0,
			Sort:      pgSort("desc"),
			CreatorId: &addressId,
		}).
		Return([]storage.ZkISM{testZkISM}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISM
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 1)
}

func (s *ZkISMTestSuite) TestListWithTxHash() {
	q := make(url.Values)
	q.Set("tx_hash", strings.ToLower(testTxHash))

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism")

	txId := uint64(1)
	s.txs.EXPECT().
		IdAndTimeByHash(gomock.Any(), testTxHashBytes).
		Return(txId, testTime, nil).
		Times(1)

	s.zkism.EXPECT().
		List(gomock.Any(), storage.ZkISMFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
			TxId:   &txId,
		}).
		Return([]storage.ZkISM{testZkISM}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISM
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 1)
}

func (s *ZkISMTestSuite) TestListValidationError() {
	q := make(url.Values)
	q.Set("limit", "200")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism")

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)
}

// ──────────────────────────────────────────────────────────
// Get
// ──────────────────────────────────────────────────────────

func (s *ZkISMTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.zkism.EXPECT().
		ById(gomock.Any(), uint64(1)).
		Return(testZkISM, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.ZkISM
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))

	s.Require().EqualValues(testZkISM.Id, response.Id)
	s.Require().EqualValues(hex.EncodeToString(testZkISM.ExternalId), response.ExternalId)
	s.Require().EqualValues(testZkISM.Height, response.Height)
	s.Require().Equal(testZkISM.Time, response.Time)
	s.Require().Equal(hex.EncodeToString(testZkISMStateRoot), response.StateRoot)
	s.Require().Equal(hex.EncodeToString(testZkISMMerkleTreeAddress), response.MerkleTreeAddress)
	s.Require().Equal(hex.EncodeToString(testZkISMStateTransitionVKey), response.StateTransitionVKey)
	s.Require().Equal(hex.EncodeToString(testZkISMStateMembershipVKey), response.StateMembershipVKey)
	s.Require().Equal(strings.ToLower(testTxHash), response.TxHash)
	s.Require().NotNil(response.Creator)
	s.Require().Equal(testAddress, response.Creator.Hash)
}

func (s *ZkISMTestSuite) TestGetValidationError() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id")
	c.SetParamNames("id")
	c.SetParamValues("0") // min=1, must fail

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)
}

// ──────────────────────────────────────────────────────────
// GetUpdates
// ──────────────────────────────────────────────────────────

func (s *ZkISMTestSuite) TestGetUpdates() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id/updates")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.zkism.EXPECT().
		Updates(gomock.Any(), uint64(1), storage.ZkISMUpdatesFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
		}).
		Return([]storage.ZkISMUpdate{testZkISMUpdate}, nil).
		Times(1)

	s.Require().NoError(s.handler.GetUpdates(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISMUpdate
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 1)

	upd := response[0]
	s.Require().EqualValues(testZkISMUpdate.Id, upd.Id)
	s.Require().EqualValues(testZkISMUpdate.Height, upd.Height)
	s.Require().Equal(testZkISMUpdate.Time, upd.Time)
	s.Require().Equal(hex.EncodeToString(testZkISMUpdateNewStateRoot), upd.NewStateRoot)
	s.Require().Equal(strings.ToLower(testTxHash), upd.TxHash)
	s.Require().NotNil(upd.Signer)
	s.Require().Equal(testAddress, upd.Signer.Hash)
}

func (s *ZkISMTestSuite) TestGetUpdatesEmpty() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id/updates")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.zkism.EXPECT().
		Updates(gomock.Any(), uint64(1), storage.ZkISMUpdatesFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
		}).
		Return([]storage.ZkISMUpdate{}, nil).
		Times(1)

	s.Require().NoError(s.handler.GetUpdates(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISMUpdate
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 0)
}

func (s *ZkISMTestSuite) TestGetUpdatesWithAddress() {
	q := make(url.Values)
	q.Set("address", testAddress)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id/updates")
	c.SetParamNames("id")
	c.SetParamValues("1")

	signerId := uint64(1)
	s.address.EXPECT().
		IdByAddress(gomock.Any(), testAddress).
		Return(signerId, nil).
		Times(1)

	s.zkism.EXPECT().
		Updates(gomock.Any(), uint64(1), storage.ZkISMUpdatesFilter{
			Limit:    10,
			Offset:   0,
			Sort:     pgSort("desc"),
			SignerId: &signerId,
		}).
		Return([]storage.ZkISMUpdate{testZkISMUpdate}, nil).
		Times(1)

	s.Require().NoError(s.handler.GetUpdates(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISMUpdate
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 1)
}

func (s *ZkISMTestSuite) TestGetUpdatesValidationError() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id/updates")
	c.SetParamNames("id")
	c.SetParamValues("0") // min=1

	s.Require().NoError(s.handler.GetUpdates(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)
}

// ──────────────────────────────────────────────────────────
// GetMessages
// ──────────────────────────────────────────────────────────

func (s *ZkISMTestSuite) TestGetMessages() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id/messages")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.zkism.EXPECT().
		Messages(gomock.Any(), uint64(1), storage.ZkISMUpdatesFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
		}).
		Return([]storage.ZkISMMessage{testZkISMMsg}, nil).
		Times(1)

	s.Require().NoError(s.handler.GetMessages(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISMMessage
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 1)

	msg := response[0]
	s.Require().EqualValues(testZkISMMsg.Id, msg.Id)
	s.Require().EqualValues(testZkISMMsg.Height, msg.Height)
	s.Require().Equal(testZkISMMsg.Time, msg.Time)
	s.Require().Equal(hex.EncodeToString(testZkISMMessageStateRoot), msg.StateRoot)
	s.Require().Equal(hex.EncodeToString(testZkISMMessageId), msg.MessageId)
	s.Require().Equal(strings.ToLower(testTxHash), msg.TxHash)
	s.Require().NotNil(msg.Signer)
	s.Require().Equal(testAddress, msg.Signer.Hash)
}

func (s *ZkISMTestSuite) TestGetMessagesEmpty() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id/messages")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.zkism.EXPECT().
		Messages(gomock.Any(), uint64(1), storage.ZkISMUpdatesFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
		}).
		Return([]storage.ZkISMMessage{}, nil).
		Times(1)

	s.Require().NoError(s.handler.GetMessages(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISMMessage
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 0)
}

func (s *ZkISMTestSuite) TestGetMessagesWithAddress() {
	q := make(url.Values)
	q.Set("address", testAddress)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id/messages")
	c.SetParamNames("id")
	c.SetParamValues("1")

	signerId := uint64(1)
	s.address.EXPECT().
		IdByAddress(gomock.Any(), testAddress).
		Return(signerId, nil).
		Times(1)

	s.zkism.EXPECT().
		Messages(gomock.Any(), uint64(1), storage.ZkISMUpdatesFilter{
			Limit:    10,
			Offset:   0,
			Sort:     pgSort("desc"),
			SignerId: &signerId,
		}).
		Return([]storage.ZkISMMessage{testZkISMMsg}, nil).
		Times(1)

	s.Require().NoError(s.handler.GetMessages(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISMMessage
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 1)
}

func (s *ZkISMTestSuite) TestGetMessagesWithTxHash() {
	q := make(url.Values)
	q.Set("tx_hash", strings.ToLower(testTxHash))

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id/messages")
	c.SetParamNames("id")
	c.SetParamValues("1")

	txId := uint64(1)
	s.txs.EXPECT().
		IdAndTimeByHash(gomock.Any(), testTxHashBytes).
		Return(txId, testTime, nil).
		Times(1)

	s.zkism.EXPECT().
		Messages(gomock.Any(), uint64(1), storage.ZkISMUpdatesFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("desc"),
			TxId:   &txId,
		}).
		Return([]storage.ZkISMMessage{testZkISMMsg}, nil).
		Times(1)

	s.Require().NoError(s.handler.GetMessages(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.ZkISMMessage
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(&response))
	s.Require().Len(response, 1)
}

func (s *ZkISMTestSuite) TestGetMessagesValidationError() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/zkism/:id/messages")
	c.SetParamNames("id")
	c.SetParamValues("0") // min=1

	s.Require().NoError(s.handler.GetMessages(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)
}
