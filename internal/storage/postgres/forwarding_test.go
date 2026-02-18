// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestForwardingById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	fwd, prevTime, err := s.storage.Forwardings.ById(ctx, 2)
	s.Require().NoError(err)
	s.Require().EqualValues(2, fwd.Id)
	s.Require().EqualValues(10000, fwd.Height)
	s.Require().Equal(fwd.Time.Unix(), prevTime.Unix())

	s.Require().EqualValues(5, fwd.AddressId)
	s.Require().NotNil(fwd.Address)
	s.Require().EqualValues("celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt", fwd.Address.Address)

	s.Require().EqualValues(5, fwd.TxId)
	s.Require().NotNil(fwd.Tx)
	s.Require().EqualValues("d764fea03c8d8dbf0608d0e24ab0b600adb15149b465356cc73d78b2278e38d5", hex.EncodeToString(fwd.Tx.Hash))

	s.Require().NotNil(fwd.Transfers)
}

func (s *StorageTestSuite) TestForwardingByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	forwards, err := s.storage.Forwardings.Filter(ctx, storage.ForwardingFilter{
		Height: testsuite.Ptr(uint64(10000)),
		Limit:  1,
		Sort:   sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(forwards, 1)

	fwd := forwards[0]
	s.Require().EqualValues(1, fwd.Id)
	s.Require().EqualValues(10000, fwd.Height)

	s.Require().EqualValues(5, fwd.AddressId)
	s.Require().NotNil(fwd.Address)
	s.Require().EqualValues("celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt", fwd.Address.Address)

	s.Require().EqualValues(5, fwd.TxId)
	s.Require().NotNil(fwd.Tx)
	s.Require().EqualValues("d764fea03c8d8dbf0608d0e24ab0b600adb15149b465356cc73d78b2278e38d5", hex.EncodeToString(fwd.Tx.Hash))

	s.Require().NotNil(fwd.Transfers)
}

func (s *StorageTestSuite) TestForwardingByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	forwards, err := s.storage.Forwardings.Filter(ctx, storage.ForwardingFilter{
		TxId:  testsuite.Ptr(uint64(5)),
		Limit: 1,
		Sort:  sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(forwards, 1)

	fwd := forwards[0]
	s.Require().EqualValues(1, fwd.Id)
	s.Require().EqualValues(10000, fwd.Height)

	s.Require().EqualValues(5, fwd.AddressId)
	s.Require().NotNil(fwd.Address)
	s.Require().EqualValues("celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt", fwd.Address.Address)

	s.Require().EqualValues(5, fwd.TxId)
	s.Require().NotNil(fwd.Tx)
	s.Require().EqualValues("d764fea03c8d8dbf0608d0e24ab0b600adb15149b465356cc73d78b2278e38d5", hex.EncodeToString(fwd.Tx.Hash))

	s.Require().NotNil(fwd.Transfers)
}

func (s *StorageTestSuite) TestForwardingByAddressId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	forwards, err := s.storage.Forwardings.Filter(ctx, storage.ForwardingFilter{
		AddressId: testsuite.Ptr(uint64(5)),
		Limit:     1,
		Sort:      sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(forwards, 1)

	fwd := forwards[0]
	s.Require().EqualValues(1, fwd.Id)
	s.Require().EqualValues(10000, fwd.Height)

	s.Require().EqualValues(5, fwd.AddressId)
	s.Require().NotNil(fwd.Address)
	s.Require().EqualValues("celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt", fwd.Address.Address)

	s.Require().EqualValues(5, fwd.TxId)
	s.Require().NotNil(fwd.Tx)
	s.Require().EqualValues("d764fea03c8d8dbf0608d0e24ab0b600adb15149b465356cc73d78b2278e38d5", hex.EncodeToString(fwd.Tx.Hash))

	s.Require().NotNil(fwd.Transfers)
}

func (s *StorageTestSuite) TestForwardingByFrom() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	forwards, err := s.storage.Forwardings.Filter(ctx, storage.ForwardingFilter{
		From:  time.Unix(1600000000, 0),
		Limit: 1,
		Sort:  sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(forwards, 1)

	fwd := forwards[0]
	s.Require().EqualValues(1, fwd.Id)
	s.Require().EqualValues(10000, fwd.Height)

	s.Require().EqualValues(5, fwd.AddressId)
	s.Require().NotNil(fwd.Address)
	s.Require().EqualValues("celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt", fwd.Address.Address)

	s.Require().EqualValues(5, fwd.TxId)
	s.Require().NotNil(fwd.Tx)
	s.Require().EqualValues("d764fea03c8d8dbf0608d0e24ab0b600adb15149b465356cc73d78b2278e38d5", hex.EncodeToString(fwd.Tx.Hash))

	s.Require().NotNil(fwd.Transfers)
}

func (s *StorageTestSuite) TestForwardingByTo() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	forwards, err := s.storage.Forwardings.Filter(ctx, storage.ForwardingFilter{
		To:    time.Unix(1771334044, 0),
		Limit: 1,
		Sort:  sdk.SortOrderAsc,
	})
	s.Require().NoError(err)
	s.Require().Len(forwards, 1)

	fwd := forwards[0]
	s.Require().EqualValues(1, fwd.Id)
	s.Require().EqualValues(10000, fwd.Height)

	s.Require().EqualValues(5, fwd.AddressId)
	s.Require().NotNil(fwd.Address)
	s.Require().EqualValues("celestia1ccqy2wlzf2zndn4vspmuksw5frqq0ufsgw4gmt", fwd.Address.Address)

	s.Require().EqualValues(5, fwd.TxId)
	s.Require().NotNil(fwd.Tx)
	s.Require().EqualValues("d764fea03c8d8dbf0608d0e24ab0b600adb15149b465356cc73d78b2278e38d5", hex.EncodeToString(fwd.Tx.Hash))

	s.Require().NotNil(fwd.Transfers)
}

func (s *StorageTestSuite) TestForwardingInputs() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	inputs, err := s.storage.Forwardings.Inputs(ctx, 5, time.Unix(1600000000, 0), time.Unix(1771334044, 0))
	s.Require().NoError(err)
	s.Require().Len(inputs, 1)

	input1 := inputs[0]
	s.Require().EqualValues(1000, input1.Height)
	s.Require().EqualValues("652452a670011d629cc116e510ba88c1cabe061336661b1f3d206d248bd55811", hex.EncodeToString(input1.TxHash))
	s.Require().Equal("1234567890abcdef", input1.From)
	s.Require().Equal("1000", input1.Amount)
	s.Require().Equal("utia", input1.Denom)
	s.Require().EqualValues(123450, input1.Counterparty)
}
