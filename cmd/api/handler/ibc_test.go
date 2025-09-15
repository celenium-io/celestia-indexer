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
	mockIbc "github.com/celenium-io/celestia-indexer/cmd/api/ibc_relayer"
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
	CreatorId:            1,
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

var relayersMap = map[uint64]responses.Relayer{
	1: {
		Name: "Test name 1",
		Logo: "https://example.com/logo1.png",
		Contact: &responses.Contact{
			Website: "https://test1.io",
			Github:  "https://github.com/testrepo1",
			Twitter: "https://twitter.com/test1",
		},
		Addresses: []string{"celestia1xyz1488"},
	},
	2: {
		Name: "Test name 2",
		Logo: "https://example.com/logo2.png",
		Contact: &responses.Contact{
			Website: "https://test2.io",
			Github:  "https://github.com/testrepo2",
			Twitter: "https://twitter.com/test2",
		},
		Addresses: []string{"celestia1xyz2222"},
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
	relayers  *mockIbc.MockIRelayerStore
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
	s.relayers = mockIbc.NewMockIRelayerStore(s.ctrl)

	s.handler = NewIbcHandler(s.clients, s.conns, s.channels, s.transfers, s.address, s.txs, s.relayers)
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

	s.relayers.EXPECT().
		List().
		Return(relayersMap).
		Times(1)

	s.transfers.EXPECT().
		List(gomock.Any(), storage.ListIbcTransferFilters{
			Limit: 10,
			Sort:  sdk.SortOrderDesc,
		}).
		Return([]storage.IbcTransferWithSigner{
			{
				IbcTransfer: storage.IbcTransfer{
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
				SignerId: testsuite.Ptr(uint64(1)),
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
	s.Require().NotNil(transfer.Relayer)
	s.Require().EqualValues("Test name 1", transfer.Relayer.Name)
	s.Require().EqualValues("https://example.com/logo1.png", transfer.Relayer.Logo)
	s.Require().NotNil(transfer.Relayer.Contact)
	s.Require().EqualValues("https://test1.io", transfer.Relayer.Contact.Website)
	s.Require().EqualValues("https://github.com/testrepo1", transfer.Relayer.Contact.Github)
	s.Require().EqualValues("https://twitter.com/test1", transfer.Relayer.Contact.Twitter)
	s.Require().Len(transfer.Relayer.Addresses, 1)
	s.Require().EqualValues("celestia1xyz1488", transfer.Relayer.Addresses[0])
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

	s.relayers.EXPECT().
		List().
		Return(relayersMap).
		Times(1)

	s.transfers.EXPECT().
		List(gomock.Any(), storage.ListIbcTransferFilters{
			Limit:         10,
			Sort:          sdk.SortOrderDesc,
			ConnectionIds: []string{"connection-1"},
		}).
		Return([]storage.IbcTransferWithSigner{
			{
				IbcTransfer: storage.IbcTransfer{
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
				SignerId: testsuite.Ptr(uint64(1)),
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
	s.Require().NotNil(transfer.Relayer)
	s.Require().EqualValues("Test name 1", transfer.Relayer.Name)
	s.Require().EqualValues("https://example.com/logo1.png", transfer.Relayer.Logo)
	s.Require().NotNil(transfer.Relayer.Contact)
	s.Require().EqualValues("https://test1.io", transfer.Relayer.Contact.Website)
	s.Require().EqualValues("https://github.com/testrepo1", transfer.Relayer.Contact.Github)
	s.Require().EqualValues("https://twitter.com/test1", transfer.Relayer.Contact.Twitter)
	s.Require().Len(transfer.Relayer.Addresses, 1)
	s.Require().EqualValues("celestia1xyz1488", transfer.Relayer.Addresses[0])
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
		Return(storage.IbcTransferWithSigner{
			IbcTransfer: storage.IbcTransfer{
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
						ChainId:   "chain-id",
						CreatorId: 1,
					},
				},
			},
			SignerId: testsuite.Ptr(uint64(1)),
		}, nil).
		Times(1)

	s.relayers.EXPECT().
		List().
		Return(relayersMap).
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
	s.Require().NotNil(response.Relayer)
	s.Require().EqualValues("Test name 1", response.Relayer.Name)
	s.Require().EqualValues("https://example.com/logo1.png", response.Relayer.Logo)
	s.Require().NotNil(response.Relayer.Contact)
	s.Require().EqualValues("https://test1.io", response.Relayer.Contact.Website)
	s.Require().EqualValues("https://github.com/testrepo1", response.Relayer.Contact.Github)
	s.Require().EqualValues("https://twitter.com/test1", response.Relayer.Contact.Twitter)
	s.Require().Len(response.Relayer.Addresses, 1)
	s.Require().EqualValues("celestia1xyz1488", response.Relayer.Addresses[0])
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

	s.relayers.EXPECT().
		List().
		Return(relayersMap).
		Times(1)

	s.transfers.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return([]storage.IbcTransferWithSigner{
			{
				IbcTransfer: storage.IbcTransfer{
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
							ChainId:   "chain-id",
							CreatorId: 1,
						},
					},
				},
				SignerId: testsuite.Ptr(uint64(1)),
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
	s.Require().NotNil(transfer.Relayer)
	s.Require().EqualValues("Test name 1", transfer.Relayer.Name)
	s.Require().EqualValues("https://example.com/logo1.png", transfer.Relayer.Logo)
	s.Require().NotNil(transfer.Relayer.Contact)
	s.Require().EqualValues("https://test1.io", transfer.Relayer.Contact.Website)
	s.Require().EqualValues("https://github.com/testrepo1", transfer.Relayer.Contact.Github)
	s.Require().EqualValues("https://twitter.com/test1", transfer.Relayer.Contact.Twitter)
	s.Require().Len(transfer.Relayer.Addresses, 1)
	s.Require().EqualValues("celestia1xyz1488", transfer.Relayer.Addresses[0])
}

func (s *IbcTestSuite) TestAllRelayers() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/ibc/relayers")

	values := make([]responses.Relayer, 0, len(relayersMap))
	for _, v := range relayersMap {
		values = append(values, v)
	}

	s.relayers.EXPECT().
		All().
		Return(values).
		Times(1)

	s.Require().NoError(s.handler.IbcRelayers(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.Relayer
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 2)
	s.Require().EqualValues("Test name 1", response[0].Name)
	s.Require().EqualValues("https://example.com/logo1.png", response[0].Logo)
	s.Require().Len(response[0].Addresses, 1)
	s.Require().EqualValues(response[0].Addresses[0], "celestia1xyz1488")
	s.Require().NotNil(response[0].Contact)
	s.Require().EqualValues(response[0].Contact.Website, "https://test1.io")
	s.Require().EqualValues(response[0].Contact.Github, "https://github.com/testrepo1")
	s.Require().EqualValues(response[0].Contact.Twitter, "https://twitter.com/test1")
}
