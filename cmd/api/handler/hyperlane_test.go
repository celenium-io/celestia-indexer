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
		DomainId:    1,
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

	testTransfer = storage.HLTransfer{
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
		Counterparty:        1,
		CounterpartyAddress: testsuite.RandomText(20),
		Nonce:               12,
		Version:             1,
		Body:                testsuite.RandomBytes(32),
		Metadata:            testsuite.RandomBytes(32),
		Amount:              decimal.RequireFromString("125678"),
		Denom:               currency.Utia,
		Type:                types.HLTransferTypeReceive,
	}

	testTransactionHash         = "123452A670018D678CC116E510BA88C1CABE061336661B1F3D206D248BD558GH"
	testTransactionHashBytes, _ = hex.DecodeString(testTransactionHash)
	testTransaction             = storage.Tx{
		Id:            2,
		Hash:          testTransactionHashBytes,
		Height:        200,
		Time:          testTime,
		Position:      3,
		GasWanted:     80410,
		GasUsed:       77483,
		TimeoutHeight: 0,
		EventsCount:   11,
		MessagesCount: 3,
		Fee:           decimal.RequireFromString("80410"),
		Status:        types.StatusSuccess,
		Codespace:     "sdk",
		Memo:          "memo",
		MessageTypes:  types.NewMsgTypeBitMask(types.MsgSend),
		Messages: []storage.Message{
			{
				Id:   1,
				Type: types.MsgSend,
			},
		},
		Signers: []storage.Address{
			{
				Address: testAddress,
			},
		},
	}

	testIgp = storage.HLIGP{
		Id:      1,
		Height:  1488,
		Time:    testTime,
		OwnerId: 1,
		Denom:   currency.Utia,
		IgpId:   testsuite.RandomBytes(32),
		Owner: &storage.Address{
			Address:    testAddress,
			Hash:       testHashAddress,
			Height:     1488,
			LastHeight: 1488,
		},
		Config: &storage.HLIGPConfig{
			GasPrice:          decimal.RequireFromString("1"),
			GasOverhead:       decimal.RequireFromString("100000"),
			TokenExchangeRate: "1234",
			RemoteDomain:      4321,
		},
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
	txs        *mock.MockITx
	igp        *mock.MockIHLIGP
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
	s.txs = mock.NewMockITx(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
	s.igp = mock.NewMockIHLIGP(s.ctrl)
	s.chainStore = hyperlane.NewMockIChainStore(s.ctrl)
	s.handler = NewHyperlaneHandler(s.mailbox, s.token, s.transfer, s.txs, s.address, s.igp, s.chainStore)
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

func (s *HyperlaneTestSuite) TestGetTransfer() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/transfer/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.chainStore.EXPECT().
		Set(testChainStore).
		Times(1)

	s.chainStore.EXPECT().
		Get(uint64(1)).
		Return(testChainMetadata, true)

	s.transfer.EXPECT().
		ById(gomock.Any(), uint64(1)).
		Return(testTransfer, nil).
		Times(1)

	s.chainStore.Set(testChainStore)

	s.Require().NoError(s.handler.GetTransfer(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.HyperlaneTransfer
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(testTransfer.Id, response.Id)
	s.Require().EqualValues(testTransfer.Height, response.Height)
	s.Require().EqualValues(testTransfer.Time, response.Time)
	s.Require().EqualValues(hex.EncodeToString(testTransfer.Mailbox.Mailbox), response.Mailbox)
	s.Require().EqualValues(testTransfer.Type.String(), response.Type)
	s.Require().EqualValues(testTransfer.Denom, response.Denom)
	s.Require().EqualValues(testTransfer.Counterparty, response.Counterparty.Domain)
	s.Require().EqualValues(testTransfer.CounterpartyAddress, response.Counterparty.Hash)
	s.Require().EqualValues(testTransfer.Version, response.Version)
	s.Require().EqualValues(testTransfer.Nonce, response.Nonce)
	s.Require().EqualValues(testTransfer.Amount.String(), response.Amount)
	s.Require().EqualValues(hex.EncodeToString(testTransfer.Token.TokenId), response.TokenId)
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

func (s *HyperlaneTestSuite) TestListTransferWithHash() {
	q := make(url.Values)
	q.Add("hash", testTxHash)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/transfer")

	s.chainStore.EXPECT().
		Set(testChainStore).
		Times(1)

	s.chainStore.EXPECT().
		Get(gomock.Any()).
		Return(testChainMetadata, true).
		Times(1)

	s.txs.EXPECT().
		ByHash(gomock.Any(), testTxHashBytes).
		Return(testTx, nil)

	s.transfer.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return([]storage.HLTransfer{
			testTransfer,
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
	s.Require().EqualValues(testTransfer.Id, response.Id)
	s.Require().EqualValues(testTransfer.Height, response.Height)
	s.Require().EqualValues(testTransfer.Time, response.Time)
	s.Require().EqualValues(hex.EncodeToString(testTransfer.Mailbox.Mailbox), response.Mailbox)
	s.Require().EqualValues(testTransfer.Type.String(), response.Type)
	s.Require().EqualValues(testTransfer.Denom, response.Denom)
	s.Require().EqualValues(testTransfer.Counterparty, response.Counterparty.Domain)
	s.Require().EqualValues(testTransfer.CounterpartyAddress, response.Counterparty.Hash)
	s.Require().EqualValues(testTransfer.Version, response.Version)
	s.Require().EqualValues(testTransfer.Nonce, response.Nonce)
	s.Require().EqualValues(testTransfer.Amount.String(), response.Amount)
	s.Require().EqualValues(hex.EncodeToString(testTransfer.Token.TokenId), response.TokenId)
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

func (s *HyperlaneTestSuite) TestListTransfer() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/transfer")

	transfers := []storage.HLTransfer{
		testTransfer,
		{
			Id:        2,
			Height:    1001,
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
			Tx:                  &testTransaction,
			Counterparty:        1,
			CounterpartyAddress: testsuite.RandomText(20),
			Nonce:               12,
			Version:             1,
			Body:                testsuite.RandomBytes(32),
			Metadata:            testsuite.RandomBytes(32),
			Amount:              decimal.RequireFromString("102030"),
			Denom:               currency.Utia,
			Type:                types.HLTransferTypeReceive,
		},
	}

	s.chainStore.EXPECT().
		Set(testChainStore).
		Times(1)

	s.chainStore.EXPECT().
		Get(gomock.Any()).
		Return(testChainMetadata, true).
		Times(len(transfers))

	s.transfer.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(transfers, nil).
		Times(1)

	s.chainStore.Set(testChainStore)
	s.Require().NoError(s.handler.ListTransfers(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.HyperlaneTransfer
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 2)

	response := items[0]
	s.Require().EqualValues(testTransfer.Id, response.Id)
	s.Require().EqualValues(testTransfer.Height, response.Height)
	s.Require().EqualValues(testTransfer.Time, response.Time)
	s.Require().EqualValues(hex.EncodeToString(testTransfer.Mailbox.Mailbox), response.Mailbox)
	s.Require().EqualValues(testTransfer.Type.String(), response.Type)
	s.Require().EqualValues(testTransfer.Denom, response.Denom)
	s.Require().EqualValues(testTransfer.Counterparty, response.Counterparty.Domain)
	s.Require().EqualValues(testTransfer.CounterpartyAddress, response.Counterparty.Hash)
	s.Require().EqualValues(testTransfer.Version, response.Version)
	s.Require().EqualValues(testTransfer.Nonce, response.Nonce)
	s.Require().EqualValues(testTransfer.Amount.String(), response.Amount)
	s.Require().EqualValues(hex.EncodeToString(testTransfer.Token.TokenId), response.TokenId)
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
	s.handler = NewHyperlaneHandler(s.mailbox, s.token, s.transfer, s.txs, s.address, s.igp, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/transfer")

	s.transfer.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return([]storage.HLTransfer{
			testTransfer,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListTransfers(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items []responses.HyperlaneTransfer
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	response := items[0]
	s.Require().EqualValues(testTransfer.Id, response.Id)
	s.Require().EqualValues(testTransfer.Height, response.Height)
	s.Require().EqualValues(testTransfer.Time, response.Time)
	s.Require().EqualValues(hex.EncodeToString(testTransfer.Mailbox.Mailbox), response.Mailbox)
	s.Require().EqualValues(testTransfer.Type.String(), response.Type)
	s.Require().EqualValues(testTransfer.Denom, response.Denom)
	s.Require().EqualValues(testTransfer.Counterparty, response.Counterparty.Domain)
	s.Require().EqualValues(testTransfer.CounterpartyAddress, response.Counterparty.Hash)
	s.Require().EqualValues(testTransfer.Version, response.Version)
	s.Require().EqualValues(testTransfer.Nonce, response.Nonce)
	s.Require().EqualValues(testTransfer.Amount.String(), response.Amount)
	s.Require().EqualValues(hex.EncodeToString(testTransfer.Token.TokenId), response.TokenId)
	s.Require().EqualValues(strings.ToLower(testTxHash), response.TxHash)
	s.Require().NotNil(response.Address)
	s.Require().Equal(testAddress, response.Address.Hash)
	s.Require().NotNil(response.Relayer)
	s.Require().Equal(testAddress, response.Relayer.Hash)
	s.Require().NotNil(response.Body)
	s.Require().NotNil(response.Metadata)
	s.Require().Nil(response.Counterparty.ChainMetadata)
}

func (s *HyperlaneTestSuite) TestListDomains() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/domains")

	s.chainStore.EXPECT().
		Set(testChainStore).
		Times(1)

	s.chainStore.EXPECT().
		All().
		Return(testChainStore).
		Times(1)

	s.chainStore.EXPECT().
		Get(gomock.Any()).
		Return(testChainMetadata, true).
		Times(len(testChainStore))

	s.chainStore.Set(testChainStore)
	s.Require().NoError(s.handler.ListDomains(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var items map[uint64]responses.DomainMetadata
	err := json.NewDecoder(rec.Body).Decode(&items)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	var response responses.DomainMetadata
	for _, v := range items {
		response = v
		break
	}
	s.Require().EqualValues(testChainMetadata.DomainId, response.Domain)
	s.Require().EqualValues(testChainMetadata.DisplayName, response.Name)
	s.Require().NotNil(response.BlockExplorers)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].Name, response.BlockExplorers[0].Name)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].ApiUrl, response.BlockExplorers[0].ApiUrl)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].Family, response.BlockExplorers[0].Family)
	s.Require().EqualValues(testChainMetadata.BlockExplorers[0].Url, response.BlockExplorers[0].Url)
	s.Require().NotNil(response.NativeToken)
	s.Require().EqualValues(testChainMetadata.NativeToken.Name, response.NativeToken.Name)
	s.Require().EqualValues(testChainMetadata.NativeToken.Decimals, response.NativeToken.Decimals)
	s.Require().EqualValues(testChainMetadata.NativeToken.Symbol, response.NativeToken.Symbol)
}

func (s *HyperlaneTestSuite) TestGetIgp() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/igp/:id")
	c.SetParamNames("id")
	c.SetParamValues("010203")

	s.igp.EXPECT().
		ByHash(gomock.Any(), []byte{1, 2, 3}).
		Return(testIgp, nil).
		Times(1)

	s.Require().NoError(s.handler.GetIgp(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response responses.HyperlaneIgp
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)

	s.Require().EqualValues(testIgp.Id, response.Id)
	s.Require().EqualValues(testIgp.Height, response.Height)
	s.Require().EqualValues(testIgp.Time, response.Time)
	s.Require().EqualValues(testIgp.Denom, response.Denom)
	s.Require().EqualValues(hex.EncodeToString(testIgp.IgpId), response.IgpId)
	s.Require().NotNil(response.Owner)
	s.Require().Equal(testAddress, response.Owner.Hash)
	s.Require().NotNil(response.Config)
	s.Require().EqualValues(testIgp.Config.RemoteDomain, response.Config.RemoteDomain)
	s.Require().EqualValues(testIgp.Config.GasPrice.String(), response.Config.GasPrice)
	s.Require().EqualValues(testIgp.Config.GasOverhead.String(), response.Config.GasOverhead)
	s.Require().EqualValues(testIgp.Config.TokenExchangeRate, response.Config.TokenExchangeRate)
}

func (s *HyperlaneTestSuite) TestListIgps() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/hyperlane/igp")

	s.igp.EXPECT().
		List(gomock.Any(), 10, 0).
		Return([]storage.HLIGP{testIgp}, nil).
		Times(1)

	s.Require().NoError(s.handler.ListIgps(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var response []responses.HyperlaneIgp
	err := json.NewDecoder(rec.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	igp := response[0]
	s.Require().EqualValues(testIgp.Id, igp.Id)
	s.Require().EqualValues(testIgp.Height, igp.Height)
	s.Require().EqualValues(testIgp.Time, igp.Time)
	s.Require().EqualValues(testIgp.Denom, igp.Denom)
	s.Require().EqualValues(hex.EncodeToString(testIgp.IgpId), igp.IgpId)
	s.Require().NotNil(igp.Owner)
	s.Require().Equal(testAddress, igp.Owner.Hash)
	s.Require().NotNil(igp.Config)
	s.Require().EqualValues(testIgp.Config.RemoteDomain, igp.Config.RemoteDomain)
	s.Require().EqualValues(testIgp.Config.GasPrice.String(), igp.Config.GasPrice)
	s.Require().EqualValues(testIgp.Config.GasOverhead.String(), igp.Config.GasOverhead)
	s.Require().EqualValues(testIgp.Config.TokenExchangeRate, igp.Config.TokenExchangeRate)
}
