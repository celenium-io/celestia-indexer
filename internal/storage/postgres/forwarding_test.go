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

	fwd, err := s.storage.Forwardings.ById(ctx, 1)
	s.Require().NoError(err)
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
