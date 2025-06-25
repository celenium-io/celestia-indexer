// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestAddressByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash := []byte{0xde, 0xce, 0x42, 0x5b, 0x75, 0xd6, 0x71, 0x15, 0xbd, 0xa8, 0x77, 0xe1, 0xe7, 0xa1, 0xf2, 0x62, 0xf6, 0xfa, 0x51, 0xd6}
	address, err := s.storage.Address.ByHash(ctx, hash)
	s.Require().NoError(err)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(100, address.Height)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", address.Address)
	s.Require().NotNil(address.Celestials)
	s.Require().EqualValues("name 1", address.Celestials.Id)
}

func (s *StorageTestSuite) TestIcaAddressByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash := []byte{48, 160, 217, 234, 151, 39, 229, 141, 230, 127, 49, 37, 182, 237, 136, 189, 218, 247, 87, 139, 87, 173, 20, 154, 154, 144, 84, 29, 23, 55, 212, 7}
	address, err := s.storage.Address.ByHash(ctx, hash)
	s.Require().NoError(err)
	s.Require().EqualValues(4, address.Id)
	s.Require().EqualValues(101, address.Height)
	s.Require().Equal("celestia1xzsdn65hyljcmenlxyjmdmvghhd0w4ut27k3fx56jp2p69eh6srs8p3rss", address.Address)
	s.Require().NotNil(address.Celestials)
	s.Require().EqualValues("name 4", address.Celestials.Id)
}

func (s *StorageTestSuite) TestAddressByHashWithoutCelestials() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash := []byte{0xC0, 0xD4, 0xB7, 0x92, 0x82, 0xE1, 0x60, 0x4E, 0xEB, 0xCB, 0x2C, 0x02, 0xF6, 0x2C, 0xB8, 0x37, 0x37, 0x57, 0x9C, 0x9C}
	address, err := s.storage.Address.ByHash(ctx, hash)
	s.Require().NoError(err)
	s.Require().EqualValues(3, address.Id)
	s.Require().EqualValues(102, address.Height)
	s.Require().Equal("celestia1cr2t0y5zu9sya67t9sp0vt9cxum408yuphkhex", address.Address)
	s.Require().Nil(address.Celestials)
}

func (s *StorageTestSuite) TestAddressList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	addresses, err := s.storage.Address.ListWithBalance(ctx, storage.AddressListFilter{
		Limit:  10,
		Offset: 0,
		Sort:   sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(addresses, 4)

	s.Require().EqualValues(1, addresses[0].Id)
	s.Require().EqualValues(100, addresses[0].Height)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", addresses[0].Address)
	s.Require().Equal("432", addresses[0].Balance.Spendable.String())
	s.Require().Equal("utia", addresses[0].Balance.Currency)
	s.Require().NotNil(addresses[0].Celestials)

	s.Require().EqualValues(2, addresses[1].Id)
	s.Require().EqualValues(101, addresses[1].Height)
	s.Require().Equal("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", addresses[1].Address)
	s.Require().Equal("321", addresses[1].Balance.Spendable.String())
	s.Require().Equal("utia", addresses[1].Balance.Currency)
	s.Require().NotNil(addresses[1].Celestials)

	s.Require().EqualValues(3, addresses[2].Id)
	s.Require().EqualValues(102, addresses[2].Height)
	s.Require().Equal("celestia1cr2t0y5zu9sya67t9sp0vt9cxum408yuphkhex", addresses[2].Address)
	s.Require().Equal("555", addresses[2].Balance.Spendable.String())
	s.Require().Equal("utia", addresses[2].Balance.Currency)
	s.Require().Nil(addresses[2].Celestials)
}

func (s *StorageTestSuite) TestAddressListWithSortAscHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, field := range []string{"first_height", "last_height"} {
		addresses, err := s.storage.Address.ListWithBalance(ctx, storage.AddressListFilter{
			Limit:     10,
			Offset:    0,
			Sort:      sdk.SortOrderAsc,
			SortField: field,
		})
		s.Require().NoError(err)
		s.Require().Len(addresses, 4)

		s.Require().EqualValues(1, addresses[0].Id)
		s.Require().EqualValues(100, addresses[0].Height)
		s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", addresses[0].Address)
		s.Require().Equal("432", addresses[0].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[0].Balance.Currency)

		s.Require().EqualValues(2, addresses[1].Id)
		s.Require().EqualValues(101, addresses[1].Height)
		s.Require().Equal("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", addresses[1].Address)
		s.Require().Equal("321", addresses[1].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[1].Balance.Currency)
	}
}

func (s *StorageTestSuite) TestAddressListWithSortDescHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, field := range []string{"first_height", "last_height"} {
		addresses, err := s.storage.Address.ListWithBalance(ctx, storage.AddressListFilter{
			Limit:     10,
			Offset:    0,
			Sort:      sdk.SortOrderDesc,
			SortField: field,
		})
		s.Require().NoError(err)
		s.Require().Len(addresses, 4)

		s.Require().EqualValues(3, addresses[0].Id, field)
		s.Require().EqualValues(102, addresses[0].Height, field)
		s.Require().Equal("celestia1cr2t0y5zu9sya67t9sp0vt9cxum408yuphkhex", addresses[0].Address, field)
		s.Require().Equal("555", addresses[0].Balance.Spendable.String(), field)
		s.Require().Equal("utia", addresses[0].Balance.Currency, field)
	}
}

func (s *StorageTestSuite) TestAddressListWithSortDesc() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, field := range []string{"delegated", "spendable", "unbonding"} {
		addresses, err := s.storage.Address.ListWithBalance(ctx, storage.AddressListFilter{
			Limit:     10,
			Offset:    0,
			Sort:      sdk.SortOrderDesc,
			SortField: field,
		})
		s.Require().NoError(err)
		s.Require().Len(addresses, 4)

		s.Require().EqualValues(3, addresses[0].Id)
		s.Require().EqualValues(102, addresses[0].Height)
		s.Require().Equal("celestia1cr2t0y5zu9sya67t9sp0vt9cxum408yuphkhex", addresses[0].Address)
		s.Require().Equal("555", addresses[0].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[0].Balance.Currency)
		s.Require().Nil(addresses[0].Celestials)

		s.Require().EqualValues(1, addresses[1].Id)
		s.Require().EqualValues(100, addresses[1].Height)
		s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", addresses[1].Address)
		s.Require().Equal("432", addresses[1].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[1].Balance.Currency)
		s.Require().NotNil(addresses[1].Celestials)

		s.Require().EqualValues(2, addresses[2].Id)
		s.Require().EqualValues(101, addresses[2].Height)
		s.Require().Equal("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", addresses[2].Address)
		s.Require().Equal("321", addresses[2].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[2].Balance.Currency)
		s.Require().NotNil(addresses[2].Celestials)

		s.Require().EqualValues(4, addresses[3].Id)
		s.Require().EqualValues(101, addresses[3].Height)
		s.Require().Equal("celestia1xzsdn65hyljcmenlxyjmdmvghhd0w4ut27k3fx56jp2p69eh6srs8p3rss", addresses[3].Address)
		s.Require().Equal("210", addresses[3].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[3].Balance.Currency)
		s.Require().NotNil(addresses[3].Celestials)
	}
}

func (s *StorageTestSuite) TestAddressListWithSortAsc() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, field := range []string{"delegated", "spendable", "unbonding"} {
		addresses, err := s.storage.Address.ListWithBalance(ctx, storage.AddressListFilter{
			Limit:     10,
			Offset:    0,
			Sort:      sdk.SortOrderAsc,
			SortField: field,
		})
		s.Require().NoError(err)
		s.Require().Len(addresses, 4)

		s.Require().EqualValues(1, addresses[2].Id, field)
		s.Require().EqualValues(100, addresses[2].Height)
		s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", addresses[2].Address)
		s.Require().Equal("432", addresses[2].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[2].Balance.Currency)

		s.Require().EqualValues(2, addresses[1].Id)
		s.Require().EqualValues(101, addresses[1].Height)
		s.Require().Equal("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", addresses[1].Address)
		s.Require().Equal("321", addresses[1].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[1].Balance.Currency)

		s.Require().EqualValues(101, addresses[0].Height)
		s.Require().Equal("celestia1xzsdn65hyljcmenlxyjmdmvghhd0w4ut27k3fx56jp2p69eh6srs8p3rss", addresses[0].Address)
		s.Require().Equal("210", addresses[0].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[0].Balance.Currency)
	}
}

func (s *StorageTestSuite) TestAddressMessages() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	messages, err := s.storage.Message.ByAddress(ctx, 1, storage.AddressMsgsFilter{
		Limit:  10,
		Offset: 0,
		Sort:   sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(messages, 2)

	msg := messages[0].Msg
	s.Require().EqualValues(1, msg.Id)
	s.Require().EqualValues(1000, msg.Height)
	s.Require().EqualValues(0, msg.Position)
	s.Require().EqualValues(types.MsgAddressTypeFromAddress, messages[0].Type)
	s.Require().Equal(types.MsgWithdrawDelegatorReward, msg.Type)
	s.Require().NotNil(messages[0].Tx)
	s.Require().NotNil(messages[0].Tx.Hash)
}

func (s *StorageTestSuite) TestAddressMessagesWithType() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	messages, err := s.storage.Message.ByAddress(ctx, 1, storage.AddressMsgsFilter{
		Limit:        10,
		Offset:       0,
		Sort:         sdk.SortOrderAsc,
		MessageTypes: []string{"MsgWithdrawDelegatorReward", "MsgDelegate"},
	})
	s.Require().NoError(err)
	s.Require().Len(messages, 2)

	msg := messages[0].Msg
	s.Require().EqualValues(1, msg.Id)
	s.Require().EqualValues(1000, msg.Height)
	s.Require().EqualValues(0, msg.Position)
	s.Require().EqualValues(types.MsgAddressTypeFromAddress, messages[0].Type)
	s.Require().Equal(types.MsgWithdrawDelegatorReward, msg.Type)
	s.Require().NotNil(messages[0].Tx)
}

func (s *StorageTestSuite) TestAddressStats() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, name := range []string{"tx_count", "fee", "gas_wanted", "gas_used"} {
		for _, tf := range []storage.Timeframe{storage.TimeframeHour, storage.TimeframeDay, storage.TimeframeMonth} {
			series, err := s.storage.Address.Series(ctx, 1, tf, name, storage.NewSeriesRequest(0, 0))
			s.Require().NoError(err)
			s.Require().Len(series, 1)

			item := series[0]
			s.Require().NotEqual("0", item.Value)
		}
	}
}

func (s *StorageTestSuite) TestAddressStatsError() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Address.Series(ctx, 1, storage.TimeframeDay, "invalid", storage.NewSeriesRequest(0, 0))
	s.Require().Error(err)

	_, err = s.storage.Address.Series(ctx, 1, storage.TimeframeYear, "count", storage.NewSeriesRequest(0, 0))
	s.Require().Error(err)
}

func (s *StorageTestSuite) TestAddressIdByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash := []byte{0xde, 0xce, 0x42, 0x5b, 0x75, 0xd6, 0x71, 0x15, 0xbd, 0xa8, 0x77, 0xe1, 0xe7, 0xa1, 0xf2, 0x62, 0xf6, 0xfa, 0x51, 0xd6}
	id, err := s.storage.Address.IdByHash(ctx, hash)
	s.Require().NoError(err)
	s.Require().Len(id, 1)
	s.Require().EqualValues(1, id[0])
}

func (s *StorageTestSuite) TestAddressIdByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	id, err := s.storage.Address.IdByAddress(ctx, "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", 2, 3, 4)
	s.Require().NoError(err)
	s.Require().EqualValues(2, id)
}
