// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var testIbcClient = storage.IbcClient{
	Id:                   "client-1",
	Height:               100,
	CreatedAt:            testTime,
	UpdatedAt:            testTime,
	LatestRevisionHeight: 100,
	LatestRevisionNumber: 0,
	FrozenRevisionHeight: 0,
	FrozenRevisionNumber: 0,
	TrustingPeriod:       time.Hour,
	UnbondingPeriod:      time.Microsecond,
	MaxClockDrift:        time.Minute,
	ConnectionCount:      10,
	ChainId:              "osmosis-1",
	Tx: &storage.Tx{
		Hash: testTxHashBytes,
	},
	Creator: &storage.Address{
		Address: testAddress,
	},
}

var testIbcConn = storage.IbcConnection{
	ConnectionId:             "conn-1",
	Height:                   100,
	CreatedAt:                testTime,
	ConnectedAt:              testTime,
	ClientId:                 "client-1",
	CreateTxId:               1,
	ConnectionTxId:           2,
	ConnectionHeight:         101,
	CreateTx:                 &testTx,
	CounterpartyConnectionId: "cc-1",
	CounterpartyClientId:     "cc-2",
	ChannelsCount:            1,
	ConnectionTx:             &testTx,
	Client:                   &testIbcClient,
}

var testIbcChannel = storage.IbcChannel{
	ConnectionId:          "conn-1",
	PortId:                "transfer",
	Ordering:              true,
	Version:               "ics20-3",
	Status:                types.IbcChannelStatusOpened,
	Height:                100,
	CreatedAt:             testTime,
	ConfirmedAt:           testTime,
	ClientId:              "client-1",
	CreateTxId:            1,
	ConfirmationTxId:      2,
	ConfirmationHeight:    101,
	CreateTx:              &testTx,
	CounterpartyPortId:    "transfer",
	CounterpartyChannelId: "channel-2",
	ConfirmationTx:        &testTx,
	Client:                &testIbcClient,
	CreatorId:             1,
	Creator: &storage.Address{
		Id:      1,
		Address: testAddress,
	},
}

// IbcTestSuite -
type IbcTestSuite struct {
	suite.Suite
	echo      *echo.Echo
	address   *mock.MockIAddress
	clients   *mock.MockIIbcClient
	conns     *mock.MockIIbcConnection
	channels  *mock.MockIIbcChannel
	transfers *mock.MockIIbcTransfer
	txs       *mock.MockITx
	handler   *IbcHandler
	ctrl      *gomock.Controller
}

// SetupSuite -
func (s *IbcTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.clients = mock.NewMockIIbcClient(s.ctrl)
	s.conns = mock.NewMockIIbcConnection(s.ctrl)
	s.channels = mock.NewMockIIbcChannel(s.ctrl)
	s.transfers = mock.NewMockIIbcTransfer(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
	s.handler = NewIbcHandler(s.clients, s.conns, s.channels, s.transfers, s.address, s.txs)
}

// TearDownSuite -
func (s *IbcTestSuite) TearDownSuite() {
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteIbc_Run(t *testing.T) {
	suite.Run(t, new(IbcTestSuite))
}

func (s *IbcTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/client/:id")
	c.SetParamNames("id")
	c.SetParamValues("client-1")

	s.clients.EXPECT().
		ById(gomock.Any(), "client-1").
		Return(testIbcClient, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.IbcClient
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(testIbcClient.Id, response.Id)
	s.Require().EqualValues(testIbcClient.Height, response.Height)
	s.Require().EqualValues(testIbcClient.CreatedAt, response.CreatedAt)
	s.Require().EqualValues(testIbcClient.LatestRevisionHeight, response.LatestRevisionHeight)
	s.Require().EqualValues(testIbcClient.LatestRevisionNumber, response.LatestRevisionNumber)
	s.Require().EqualValues(testIbcClient.FrozenRevisionHeight, response.FrozenRevisionHeight)
	s.Require().EqualValues(testIbcClient.FrozenRevisionNumber, response.FrozenRevisionNumber)
	s.Require().EqualValues(testIbcClient.TrustingPeriod, response.TrustingPeriod)
	s.Require().EqualValues(testIbcClient.UnbondingPeriod, response.UnbondingPeriod)
	s.Require().EqualValues(testIbcClient.MaxClockDrift, response.MaxClockDrift)
	s.Require().EqualValues(testIbcClient.ConnectionCount, response.ConnectionCount)
	s.Require().EqualValues(testIbcClient.ChainId, response.ChainId)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.TxHash)
	s.Require().EqualValues(testAddress, response.Creator.Hash)
}

func (s *IbcTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/client")

	s.clients.EXPECT().
		List(gomock.Any(), storage.ListIbcClientsFilters{
			Limit: 10,
			Sort:  sdk.SortOrderDesc,
		}).
		Return([]storage.IbcClient{testIbcClient}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.IbcClient
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	client := response[0]
	s.Require().EqualValues(testIbcClient.Id, client.Id)
	s.Require().EqualValues(testIbcClient.Height, client.Height)
	s.Require().EqualValues(testIbcClient.CreatedAt, client.CreatedAt)
	s.Require().EqualValues(testIbcClient.LatestRevisionHeight, client.LatestRevisionHeight)
	s.Require().EqualValues(testIbcClient.LatestRevisionNumber, client.LatestRevisionNumber)
	s.Require().EqualValues(testIbcClient.FrozenRevisionHeight, client.FrozenRevisionHeight)
	s.Require().EqualValues(testIbcClient.FrozenRevisionNumber, client.FrozenRevisionNumber)
	s.Require().EqualValues(testIbcClient.TrustingPeriod, client.TrustingPeriod)
	s.Require().EqualValues(testIbcClient.UnbondingPeriod, client.UnbondingPeriod)
	s.Require().EqualValues(testIbcClient.MaxClockDrift, client.MaxClockDrift)
	s.Require().EqualValues(testIbcClient.ConnectionCount, client.ConnectionCount)
	s.Require().EqualValues(testIbcClient.ChainId, client.ChainId)
	s.Require().EqualValues(strings.ToLower(testTxHash), client.TxHash)
	s.Require().EqualValues(testAddress, client.Creator.Hash)

}

func (s *IbcTestSuite) TestGetConnection() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/connection/:id")
	c.SetParamNames("id")
	c.SetParamValues("conn-1")

	s.conns.EXPECT().
		ById(gomock.Any(), "conn-1").
		Return(testIbcConn, nil).
		Times(1)

	s.Require().NoError(s.handler.GetConnection(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.IbcConnection
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(testIbcConn.ConnectionId, response.Id)
	s.Require().EqualValues(testIbcConn.Height, response.Height)
	s.Require().EqualValues(testIbcConn.ChannelsCount, response.ChannelsCount)
	s.Require().EqualValues(testIbcConn.ConnectionHeight, response.ConnectedHeight)
	s.Require().EqualValues(testIbcConn.CounterpartyClientId, response.CounterpartyClientId)
	s.Require().EqualValues(testIbcConn.CounterpartyConnectionId, response.CounterpartyConnId)
	s.Require().EqualValues(testIbcConn.CreatedAt, response.CreatedAt)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.CreatedTxHash)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.ConnectedTxHash)
	s.Require().NotNil(response.Client)
	s.Require().EqualValues(testIbcConn.Client.Id, response.Client.Id)
}

func (s *IbcTestSuite) TestListConnection() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/connection")

	s.conns.EXPECT().
		List(gomock.Any(), storage.ListConnectionFilters{
			Limit: 10,
			Sort:  sdk.SortOrderDesc,
		}).
		Return([]storage.IbcConnection{testIbcConn}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListConnections(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var conns []responses.IbcConnection
	err := json.NewDecoder(rec.Body).Decode(&conns)
	s.Require().NoError(err)
	s.Require().Len(conns, 1)

	response := conns[0]
	s.Require().EqualValues(testIbcConn.ConnectionId, response.Id)
	s.Require().EqualValues(testIbcConn.Height, response.Height)
	s.Require().EqualValues(testIbcConn.ChannelsCount, response.ChannelsCount)
	s.Require().EqualValues(testIbcConn.ConnectionHeight, response.ConnectedHeight)
	s.Require().EqualValues(testIbcConn.CounterpartyClientId, response.CounterpartyClientId)
	s.Require().EqualValues(testIbcConn.CounterpartyConnectionId, response.CounterpartyConnId)
	s.Require().EqualValues(testIbcConn.CreatedAt, response.CreatedAt)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.CreatedTxHash)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.ConnectedTxHash)
	s.Require().NotNil(response.Client)
	s.Require().EqualValues(testIbcConn.Client.Id, response.Client.Id)
}

func (s *IbcTestSuite) TestGetChannel() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/channel/:id")
	c.SetParamNames("id")
	c.SetParamValues("channel-1")

	s.channels.EXPECT().
		ById(gomock.Any(), "channel-1").
		Return(testIbcChannel, nil).
		Times(1)

	s.Require().NoError(s.handler.GetChannel(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.IbcChannel
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(testIbcChannel.Id, response.Id)
	s.Require().EqualValues(testIbcChannel.Height, response.Height)
	s.Require().EqualValues(testIbcChannel.Status.String(), response.Status)
	s.Require().EqualValues(testIbcChannel.ConfirmationHeight, response.ConfirmationHeight)
	s.Require().EqualValues(testIbcChannel.CounterpartyChannelId, response.CounterpartyChannelId)
	s.Require().EqualValues(testIbcChannel.CounterpartyPortId, response.CounterpartyPortId)
	s.Require().EqualValues(testIbcChannel.CreatedAt, response.CreatedAt)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.CreatedTxHash)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.ConfirmationTxHash)
	s.Require().NotNil(response.Client)
	s.Require().EqualValues(testIbcChannel.Client.Id, response.Client.Id)
	s.Require().NotNil(response.Creator)
	s.Require().EqualValues(testAddress, response.Creator.Hash)
	s.Require().True(response.Ordering)
}

func (s *IbcTestSuite) TestListChannels() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/channel")

	s.channels.EXPECT().
		List(gomock.Any(), storage.ListChannelFilters{
			Limit: 10,
			Sort:  sdk.SortOrderDesc,
		}).
		Return([]storage.IbcChannel{testIbcChannel}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListChannels(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var conns []responses.IbcChannel
	err := json.NewDecoder(rec.Body).Decode(&conns)
	s.Require().NoError(err)
	s.Require().Len(conns, 1)

	response := conns[0]
	s.Require().EqualValues(testIbcChannel.Id, response.Id)
	s.Require().EqualValues(testIbcChannel.Height, response.Height)
	s.Require().EqualValues(testIbcChannel.Status.String(), response.Status)
	s.Require().EqualValues(testIbcChannel.ConfirmationHeight, response.ConfirmationHeight)
	s.Require().EqualValues(testIbcChannel.CounterpartyChannelId, response.CounterpartyChannelId)
	s.Require().EqualValues(testIbcChannel.CounterpartyPortId, response.CounterpartyPortId)
	s.Require().EqualValues(testIbcChannel.CreatedAt, response.CreatedAt)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.CreatedTxHash)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.ConfirmationTxHash)
	s.Require().NotNil(response.Client)
	s.Require().EqualValues(testIbcChannel.Client.Id, response.Client.Id)
	s.Require().NotNil(response.Creator)
	s.Require().EqualValues(testAddress, response.Creator.Hash)
	s.Require().True(response.Ordering)
}

func (s *IbcTestSuite) TestListTransfers() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/transfer")

	s.transfers.EXPECT().
		List(gomock.Any(), storage.ListIbcTransferFilters{
			Limit: 10,
			Sort:  sdk.SortOrderDesc,
		}).
		Return([]storage.IbcTransfer{
			{
				Id:              1,
				Time:            testTime,
				Height:          1000,
				Timeout:         &testTime,
				ChannelId:       "channel-1",
				ConnectionId:    "connection-1",
				Amount:          decimal.RequireFromString("101"),
				Denom:           currency.Utia,
				Memo:            "memo",
				ReceiverAddress: testsuite.Ptr("osmo1mj37s3mmv78tj0ke3yely7zwmzl5rkh9gx9ma2"),
				Sender: &storage.Address{
					Hash:    testHashAddress,
					Address: testAddress,
				},
				Sequence: 123456,
				Tx:       &testTx,
				Connection: &storage.IbcConnection{
					Client: &storage.IbcClient{
						ChainId: "chain-id",
					},
				},
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListTransfers(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var transfers []responses.IbcTransfer
	err := json.NewDecoder(rec.Body).Decode(&transfers)
	s.Require().NoError(err)
	s.Require().Len(transfers, 1)

	transfer := transfers[0]
	s.Require().EqualValues(1, transfer.Id)
	s.Require().EqualValues(1000, transfer.Height)
	s.Require().EqualValues(testTime, transfer.Time)
	s.Require().NotNil(transfer.Timeout)
	s.Require().EqualValues(testTime, *transfer.Timeout)
	s.Require().EqualValues(0, transfer.TimeoutHeight)
	s.Require().EqualValues("101", transfer.Amount)
	s.Require().EqualValues("utia", transfer.Denom)
	s.Require().EqualValues("channel-1", transfer.ChannelId)
	s.Require().EqualValues("connection-1", transfer.ConnectionId)
	s.Require().EqualValues("memo", transfer.Memo)
	s.Require().EqualValues(strings.ToLower(testTxHash), transfer.TxHash)
	s.Require().NotNil(transfer.Receiver)
	s.Require().EqualValues("osmo1mj37s3mmv78tj0ke3yely7zwmzl5rkh9gx9ma2", transfer.Receiver.Hash)
	s.Require().NotNil(transfer.Sender)
	s.Require().EqualValues(testAddress, transfer.Sender.Hash)
	s.Require().Equal("chain-id", transfer.ChainId)
}

func (s *IbcTestSuite) TestListTransfersByChainId() {
	q := make(url.Values)
	q.Set("chain_id", "test")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/transfer")

	s.clients.EXPECT().
		ByChainId(gomock.Any(), "test").
		Return([]string{"client"}, nil).
		Times(1)

	s.conns.EXPECT().
		IdsByClients(gomock.Any(), "client").
		Return([]string{"connection-1"}, nil).
		Times(1)

	s.transfers.EXPECT().
		List(gomock.Any(), storage.ListIbcTransferFilters{
			Limit:         10,
			Sort:          sdk.SortOrderDesc,
			ConnectionIds: []string{"connection-1"},
		}).
		Return([]storage.IbcTransfer{
			{
				Id:              1,
				Time:            testTime,
				Height:          1000,
				Timeout:         &testTime,
				ChannelId:       "channel-1",
				ConnectionId:    "connection-1",
				Amount:          decimal.RequireFromString("101"),
				Denom:           currency.Utia,
				Memo:            "memo",
				ReceiverAddress: testsuite.Ptr("osmo1mj37s3mmv78tj0ke3yely7zwmzl5rkh9gx9ma2"),
				Sender: &storage.Address{
					Hash:    testHashAddress,
					Address: testAddress,
				},
				Sequence: 123456,
				Tx:       &testTx,
				Connection: &storage.IbcConnection{
					Client: &storage.IbcClient{
						ChainId: "chain-id",
					},
				},
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListTransfers(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var transfers []responses.IbcTransfer
	err := json.NewDecoder(rec.Body).Decode(&transfers)
	s.Require().NoError(err)
	s.Require().Len(transfers, 1)

	transfer := transfers[0]
	s.Require().EqualValues(1, transfer.Id)
	s.Require().EqualValues(1000, transfer.Height)
	s.Require().EqualValues(testTime, transfer.Time)
	s.Require().NotNil(transfer.Timeout)
	s.Require().EqualValues(testTime, *transfer.Timeout)
	s.Require().EqualValues(0, transfer.TimeoutHeight)
	s.Require().EqualValues("101", transfer.Amount)
	s.Require().EqualValues("utia", transfer.Denom)
	s.Require().EqualValues("channel-1", transfer.ChannelId)
	s.Require().EqualValues("connection-1", transfer.ConnectionId)
	s.Require().EqualValues("memo", transfer.Memo)
	s.Require().EqualValues(strings.ToLower(testTxHash), transfer.TxHash)
	s.Require().NotNil(transfer.Receiver)
	s.Require().EqualValues("osmo1mj37s3mmv78tj0ke3yely7zwmzl5rkh9gx9ma2", transfer.Receiver.Hash)
	s.Require().NotNil(transfer.Sender)
	s.Require().EqualValues(testAddress, transfer.Sender.Hash)
	s.Require().Equal("chain-id", transfer.ChainId)
}

func (s *IbcTestSuite) TestListTransfersByChainIdUnknownChain() {
	q := make(url.Values)
	q.Set("chain_id", "test")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/transfer")

	s.clients.EXPECT().
		ByChainId(gomock.Any(), "test").
		Return([]string{}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListTransfers(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var transfers []responses.IbcTransfer
	err := json.NewDecoder(rec.Body).Decode(&transfers)
	s.Require().NoError(err)
	s.Require().Len(transfers, 0)
}

func (s *IbcTestSuite) TestGetTransfer() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/transfer/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.transfers.EXPECT().
		ById(gomock.Any(), uint64(1)).
		Return(storage.IbcTransfer{
			Id:              1,
			Time:            testTime,
			Height:          1000,
			Timeout:         &testTime,
			ChannelId:       "channel-1",
			ConnectionId:    "connection-1",
			Amount:          decimal.RequireFromString("101"),
			Denom:           currency.Utia,
			Memo:            "memo",
			ReceiverAddress: testsuite.Ptr("osmo1mj37s3mmv78tj0ke3yely7zwmzl5rkh9gx9ma2"),
			Sender: &storage.Address{
				Hash:    testHashAddress,
				Address: testAddress,
			},
			Sequence: 123456,
			Tx:       &testTx,
			Connection: &storage.IbcConnection{
				Client: &storage.IbcClient{
					ChainId: "chain-id",
				},
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.GetIbcTransfer(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.IbcTransfer
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(1, response.Id)
	s.Require().EqualValues(1000, response.Height)
	s.Require().EqualValues(testTime, response.Time)
	s.Require().NotNil(response.Timeout)
	s.Require().EqualValues(testTime, *response.Timeout)
	s.Require().EqualValues(0, response.TimeoutHeight)
	s.Require().EqualValues("101", response.Amount)
	s.Require().EqualValues("utia", response.Denom)
	s.Require().EqualValues("channel-1", response.ChannelId)
	s.Require().EqualValues("connection-1", response.ConnectionId)
	s.Require().EqualValues("memo", response.Memo)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.TxHash)
	s.Require().NotNil(response.Receiver)
	s.Require().EqualValues("osmo1mj37s3mmv78tj0ke3yely7zwmzl5rkh9gx9ma2", response.Receiver.Hash)
	s.Require().NotNil(response.Sender)
	s.Require().EqualValues(testAddress, response.Sender.Hash)
	s.Require().Equal("chain-id", response.ChainId)
}

func (s *IbcTestSuite) TestListTransferWithHash() {
	q := make(url.Values)
	q.Add("hash", testTxHash)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/transfer")

	s.txs.EXPECT().
		ByHash(gomock.Any(), testTxHashBytes).
		Return(testTx, nil)

	s.transfers.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return([]storage.IbcTransfer{
			{
				Id:              1,
				Time:            testTime,
				Height:          1000,
				Timeout:         &testTime,
				ChannelId:       "channel-1",
				ConnectionId:    "connection-1",
				Amount:          decimal.RequireFromString("101"),
				Denom:           currency.Utia,
				Memo:            "memo",
				ReceiverAddress: testsuite.Ptr("osmo1mj37s3mmv78tj0ke3yely7zwmzl5rkh9gx9ma2"),
				Sender: &storage.Address{
					Hash:    testHashAddress,
					Address: testAddress,
				},
				Sequence: 123456,
				Tx:       &testTx,
				Connection: &storage.IbcConnection{
					Client: &storage.IbcClient{
						ChainId: "chain-id",
					},
				},
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListTransfers(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var transfers []responses.IbcTransfer
	err := json.NewDecoder(rec.Body).Decode(&transfers)
	s.Require().NoError(err)
	s.Require().Len(transfers, 1)

	transfer := transfers[0]
	s.Require().EqualValues(1, transfer.Id)
	s.Require().EqualValues(1000, transfer.Height)
	s.Require().EqualValues(testTime, transfer.Time)
	s.Require().NotNil(transfer.Timeout)
	s.Require().EqualValues(testTime, *transfer.Timeout)
	s.Require().EqualValues(0, transfer.TimeoutHeight)
	s.Require().EqualValues("101", transfer.Amount)
	s.Require().EqualValues("utia", transfer.Denom)
	s.Require().EqualValues("channel-1", transfer.ChannelId)
	s.Require().EqualValues("connection-1", transfer.ConnectionId)
	s.Require().EqualValues("memo", transfer.Memo)
	s.Require().EqualValues(strings.ToLower(testTxHash), transfer.TxHash)
	s.Require().NotNil(transfer.Receiver)
	s.Require().EqualValues("osmo1mj37s3mmv78tj0ke3yely7zwmzl5rkh9gx9ma2", transfer.Receiver.Hash)
	s.Require().NotNil(transfer.Sender)
	s.Require().EqualValues(testAddress, transfer.Sender.Hash)
	s.Require().Equal("chain-id", transfer.ChainId)
}
