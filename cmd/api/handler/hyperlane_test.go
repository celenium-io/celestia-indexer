// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/cmd/api/hyperlane"
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	hl "github.com/celenium-io/celestia-indexer/pkg/node/hyperlane"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testMailbox = storage.HLMailbox{
		Id:               1,
		Height:           1000,
		Time:             testTime,
		TxId:             testTx.Id,
		Mailbox:          []byte{1, 2, 3},
		OwnerId:          1,
		DefaultIsm:       testsuite.RandomBytes(20),
		DefaultHook:      testsuite.RandomBytes(20),
		RequiredHook:     testsuite.RandomBytes(20),
		Domain:           456,
		SentMessages:     10,
		ReceivedMessages: 10,
		Owner: &storage.Address{
			Address:    testAddress,
			Hash:       testHashAddress,
			Height:     1000,
			LastHeight: 1000,
		},
		Tx: &testTx,
	}

	testHyperlaneToken = storage.HLToken{
		Id:        1,
		Height:    1000,
		Time:      testTime,
		TxId:      testTx.Id,
		Tx:        &testTx,
		MailboxId: 1,
		Mailbox:   &testMailbox,
		OwnerId:   1,
		Owner: &storage.Address{
			Address:    testAddress,
			Hash:       testHashAddress,
			Height:     1000,
			LastHeight: 1000,
		},
		Type:             types.HLTokenTypeCollateral,
		Denom:            currency.Utia,
		TokenId:          testsuite.RandomBytes(32),
		SentTransfers:    10,
		ReceiveTransfers: 10,
		Sent:             decimal.RequireFromString("1000"),
		Received:         decimal.RequireFromString("1000"),
	}

	testChainMetadata = hl.ChainMetadata{
		DomainId:    123,
		DisplayName: "test chain",
		BlockExplorers: []hl.BlockExplorer{
			{
				Name:   "test explorer",
				ApiUrl: "https://api.test.url.io",
				Family: "testscan",
				Url:    "https://test.url.io",
			},
		},
		NativeToken: hl.NativeToken{
			Decimals: 18,
			Name:     "Test coin",
			Symbol:   "TEST",
		},
	}

	testChainStore = map[uint64]hl.ChainMetadata{
		testChainMetadata.DomainId: testChainMetadata,
	}
)

// HyperlaneTestSuite -
type HyperlaneTestSuite struct {
	suite.Suite
	echo       *echo.Echo
	address    *mock.MockIAddress
	mailbox    *mock.MockIHLMailbox
	token      *mock.MockIHLToken
	transfer   *mock.MockIHLTransfer
	handler    *HyperlaneHandler
	chainStore *hyperlane.MockIChainStore
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *HyperlaneTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.mailbox = mock.NewMockIHLMailbox(s.ctrl)
	s.token = mock.NewMockIHLToken(s.ctrl)
	s.transfer = mock.NewMockIHLTransfer(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
	s.chainStore = hyperlane.NewMockIChainStore(s.ctrl)
	s.handler = NewHyperlaneHandler(s.mailbox, s.token, s.transfer, s.address, s.chainStore)
}

// TearDownSuite -
func (s *HyperlaneTestSuite) TearDownSuite() {
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteHyperlane_Run(t *testing.T) {
	suite.Run(t, new(HyperlaneTestSuite))
}

func (s *HyperlaneTestSuite) TestGetMailbox() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/mailbox/:id")
	c.SetParamNames("id")
	c.SetParamValues("010203")

	s.mailbox.EXPECT().
		ByHash(gomock.Any(), []byte{1, 2, 3}).
		Return(testMailbox, nil).
		Times(1)

	s.Require().NoError(s.handler.GetMailbox(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.HyperlaneMailbox
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(testMailbox.Id, response.Id)
	s.Require().EqualValues(testMailbox.Height, response.Height)
	s.Require().EqualValues(testMailbox.Time, response.Time)
	s.Require().EqualValues(testMailbox.Domain, response.Domain)
	s.Require().EqualValues(testMailbox.SentMessages, response.SentMessages)
	s.Require().EqualValues(testMailbox.ReceivedMessages, response.ReceivedMessages)
	s.Require().EqualValues(hex.EncodeToString(testMailbox.Mailbox), response.Mailbox)
	s.Require().EqualValues(hex.EncodeToString(testMailbox.DefaultHook), response.DefaultHook)
	s.Require().EqualValues(hex.EncodeToString(testMailbox.DefaultIsm), response.DefaultIsm)
	s.Require().EqualValues(hex.EncodeToString(testMailbox.RequiredHook), response.RequiredHook)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.TxHash)
	s.Require().NotNil(response.Owner)
	s.Require().Equal(testAddress, response.Owner.Hash)
}

func (s *HyperlaneTestSuite) TestListMailboxes() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/mailbox")

	s.mailbox.EXPECT().
		List(gomock.Any(), 10, 0).
		Return([]storage.HLMailbox{testMailbox}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListMailboxes(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.HyperlaneMailbox
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	mailbox := response[0]
	s.Require().EqualValues(testMailbox.Id, mailbox.Id)
	s.Require().EqualValues(testMailbox.Height, mailbox.Height)
	s.Require().EqualValues(testMailbox.Time, mailbox.Time)
	s.Require().EqualValues(testMailbox.Domain, mailbox.Domain)
	s.Require().EqualValues(testMailbox.SentMessages, mailbox.SentMessages)
	s.Require().EqualValues(testMailbox.ReceivedMessages, mailbox.ReceivedMessages)
	s.Require().EqualValues(hex.EncodeToString(testMailbox.Mailbox), mailbox.Mailbox)
	s.Require().EqualValues(hex.EncodeToString(testMailbox.DefaultHook), mailbox.DefaultHook)
	s.Require().EqualValues(hex.EncodeToString(testMailbox.DefaultIsm), mailbox.DefaultIsm)
	s.Require().EqualValues(hex.EncodeToString(testMailbox.RequiredHook), mailbox.RequiredHook)
	s.Require().EqualValues(strings.ToLower(testTxHash), mailbox.TxHash)
	s.Require().NotNil(mailbox.Owner)
	s.Require().Equal(testAddress, mailbox.Owner.Hash)
}

func (s *HyperlaneTestSuite) TestGetToken() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/token/:id")
	c.SetParamNames("id")
	c.SetParamValues("010203")

	s.token.EXPECT().
		ByHash(gomock.Any(), []byte{1, 2, 3}).
		Return(testHyperlaneToken, nil).
		Times(1)

	s.Require().NoError(s.handler.GetToken(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.HyperlaneToken
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(testHyperlaneToken.Id, response.Id)
	s.Require().EqualValues(testHyperlaneToken.Height, response.Height)
	s.Require().EqualValues(testHyperlaneToken.Time, response.Time)
	s.Require().EqualValues(hex.EncodeToString(testHyperlaneToken.Mailbox.Mailbox), response.Mailbox)
	s.Require().EqualValues(testHyperlaneToken.Type.String(), response.Type)
	s.Require().EqualValues(testHyperlaneToken.Denom, response.Denom)
	s.Require().EqualValues(testHyperlaneToken.SentTransfers, response.SentTransfers)
	s.Require().EqualValues(testHyperlaneToken.ReceiveTransfers, response.ReceiveTransfers)
	s.Require().EqualValues(testHyperlaneToken.Sent.String(), response.Sent)
	s.Require().EqualValues(testHyperlaneToken.Received.String(), response.Received)
	s.Require().EqualValues(hex.EncodeToString(testHyperlaneToken.TokenId), response.TokenId)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.TxHash)
	s.Require().NotNil(response.Owner)
	s.Require().Equal(testAddress, response.Owner.Hash)
}

func (s *HyperlaneTestSuite) TestListToken() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/token")

	s.token.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return([]storage.HLToken{testHyperlaneToken}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListTokens(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.HyperlaneToken
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().EqualValues(testHyperlaneToken.Id, response.Id)
	s.Require().EqualValues(testHyperlaneToken.Height, response.Height)
	s.Require().EqualValues(testHyperlaneToken.Time, response.Time)
	s.Require().EqualValues(hex.EncodeToString(testHyperlaneToken.Mailbox.Mailbox), response.Mailbox)
	s.Require().EqualValues(testHyperlaneToken.Type.String(), response.Type)
	s.Require().EqualValues(testHyperlaneToken.Denom, response.Denom)
	s.Require().EqualValues(testHyperlaneToken.SentTransfers, response.SentTransfers)
	s.Require().EqualValues(testHyperlaneToken.ReceiveTransfers, response.ReceiveTransfers)
	s.Require().EqualValues(testHyperlaneToken.Sent.String(), response.Sent)
	s.Require().EqualValues(testHyperlaneToken.Received.String(), response.Received)
	s.Require().EqualValues(hex.EncodeToString(testHyperlaneToken.TokenId), response.TokenId)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.TxHash)
	s.Require().NotNil(response.Owner)
	s.Require().Equal(testAddress, response.Owner.Hash)
}

func (s *HyperlaneTestSuite) TestListTransfer() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/transfer")

	transfer := storage.HLTransfer{
		Id:        1,
		Height:    1000,
		Time:      testTime,
		MailboxId: 1,
		Mailbox:   &testMailbox,
		AddressId: 1,
		Address: &storage.Address{
			Address: testAddress,
		},
		RelayerId: 1,
		Relayer: &storage.Address{
			Address: testAddress,
		},
		TokenId:             1,
		Token:               &testHyperlaneToken,
		TxId:                1,
		Tx:                  &testTx,
		Counterparty:        123,
		CounterpartyAddress: testsuite.RandomText(20),
		Nonce:               12,
		Version:             1,
		Body:                testsuite.RandomBytes(32),
		Metadata:            testsuite.RandomBytes(32),
		Amount:              decimal.RequireFromString("125678"),
		Denom:               currency.Utia,
		Type:                types.HLTransferTypeReceive,
	}

	s.chainStore.EXPECT().
		Set(testChainStore).
		Times(1)

	s.chainStore.EXPECT().
		Get(uint64(123)).
		Return(testChainMetadata, true)

	s.transfer.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return([]storage.HLTransfer{
			transfer,
		}, nil).
		Times(1)

	s.chainStore.Set(testChainStore)

	s.Require().NoError(s.handler.ListTransfers(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.HyperlaneTransfer
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().EqualValues(transfer.Id, response.Id)
	s.Require().EqualValues(transfer.Height, response.Height)
	s.Require().EqualValues(transfer.Time, response.Time)
	s.Require().EqualValues(hex.EncodeToString(transfer.Mailbox.Mailbox), response.Mailbox)
	s.Require().EqualValues(transfer.Type.String(), response.Type)
	s.Require().EqualValues(transfer.Denom, response.Denom)
	s.Require().EqualValues(transfer.Counterparty, response.Counterparty.Domain)
	s.Require().EqualValues(transfer.CounterpartyAddress, response.Counterparty.Hash)
	s.Require().EqualValues(transfer.Version, response.Version)
	s.Require().EqualValues(transfer.Nonce, response.Nonce)
	s.Require().EqualValues(transfer.Amount.String(), response.Amount)
	s.Require().EqualValues(hex.EncodeToString(transfer.Token.TokenId), response.TokenId)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.TxHash)
	s.Require().NotNil(response.Address)
	s.Require().Equal(testAddress, response.Address.Hash)
	s.Require().NotNil(response.Relayer)
	s.Require().Equal(testAddress, response.Relayer.Hash)
	s.Require().NotNil(response.Body)
	s.Require().NotNil(response.Metadata)
	s.Require().EqualValues(testChainMetadata.DisplayName, response.Counterparty.ChainMetadata.Name)
	s.Require().EqualValues(testChainMetadata.DomainId, response.Counterparty.Domain)
	s.Require().EqualValues(testChainMetadata.NativeToken.Decimals, response.Counterparty.ChainMetadata.NativeToken.Decimals)
	s.Require().EqualValues(testChainMetadata.NativeToken.Name, response.Counterparty.ChainMetadata.NativeToken.Name)
	s.Require().EqualValues(testChainMetadata.NativeToken.Symbol, response.Counterparty.ChainMetadata.NativeToken.Symbol)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].Name, response.Counterparty.ChainMetadata.BlockExplorers[0].Name)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].ApiUrl, response.Counterparty.ChainMetadata.BlockExplorers[0].ApiUrl)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].Url, response.Counterparty.ChainMetadata.BlockExplorers[0].Url)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].Family, response.Counterparty.ChainMetadata.BlockExplorers[0].Family)
}

func (s *HyperlaneTestSuite) TestListTransferWithoutChainStore() {
	s.chainStore = nil
	s.handler = NewHyperlaneHandler(s.mailbox, s.token, s.transfer, s.address, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/transfer")

	transfer := storage.HLTransfer{
		Id:        1,
		Height:    1000,
		Time:      testTime,
		MailboxId: 1,
		Mailbox:   &testMailbox,
		AddressId: 1,
		Address: &storage.Address{
			Address: testAddress,
		},
		RelayerId: 1,
		Relayer: &storage.Address{
			Address: testAddress,
		},
		TokenId:             1,
		Token:               &testHyperlaneToken,
		TxId:                1,
		Tx:                  &testTx,
		Counterparty:        123,
		CounterpartyAddress: testsuite.RandomText(20),
		Nonce:               12,
		Version:             1,
		Body:                testsuite.RandomBytes(32),
		Metadata:            testsuite.RandomBytes(32),
		Amount:              decimal.RequireFromString("125678"),
		Denom:               currency.Utia,
		Type:                types.HLTransferTypeReceive,
	}

	s.transfer.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return([]storage.HLTransfer{
			transfer,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListTransfers(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.HyperlaneTransfer
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().EqualValues(transfer.Id, response.Id)
	s.Require().EqualValues(transfer.Height, response.Height)
	s.Require().EqualValues(transfer.Time, response.Time)
	s.Require().EqualValues(hex.EncodeToString(transfer.Mailbox.Mailbox), response.Mailbox)
	s.Require().EqualValues(transfer.Type.String(), response.Type)
	s.Require().EqualValues(transfer.Denom, response.Denom)
	s.Require().EqualValues(transfer.Counterparty, response.Counterparty.Domain)
	s.Require().EqualValues(transfer.CounterpartyAddress, response.Counterparty.Hash)
	s.Require().EqualValues(transfer.Version, response.Version)
	s.Require().EqualValues(transfer.Nonce, response.Nonce)
	s.Require().EqualValues(transfer.Amount.String(), response.Amount)
	s.Require().EqualValues(hex.EncodeToString(transfer.Token.TokenId), response.TokenId)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.TxHash)
	s.Require().NotNil(response.Address)
	s.Require().Equal(testAddress, response.Address.Hash)
	s.Require().NotNil(response.Relayer)
	s.Require().Equal(testAddress, response.Relayer.Hash)
	s.Require().NotNil(response.Body)
	s.Require().NotNil(response.Metadata)
	s.Require().Nil(response.Counterparty.ChainMetadata)
}
