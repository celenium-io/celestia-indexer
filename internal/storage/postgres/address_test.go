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
	s.Require().Len(addresses, 2)

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
		s.Require().Len(addresses, 2)

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
		s.Require().Len(addresses, 2)

		s.Require().EqualValues(2, addresses[0].Id, field)
		s.Require().EqualValues(101, addresses[0].Height, field)
		s.Require().Equal("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", addresses[0].Address, field)
		s.Require().Equal("321", addresses[0].Balance.Spendable.String(), field)
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
		s.Require().Len(addresses, 2)

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
		s.Require().Len(addresses, 2)

		s.Require().EqualValues(1, addresses[1].Id, field)
		s.Require().EqualValues(100, addresses[1].Height)
		s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", addresses[1].Address)
		s.Require().Equal("432", addresses[1].Balance.Spendable.String())
		s.Require().Equal("utia", addresses[1].Balance.Currency)

		s.Require().EqualValues(2, addresses[0].Id)
		s.Require().EqualValues(101, addresses[0].Height)
		s.Require().Equal("celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60", addresses[0].Address)
		s.Require().Equal("321", addresses[0].Balance.Spendable.String())
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
	s.Require().EqualValues(1, id)
}
