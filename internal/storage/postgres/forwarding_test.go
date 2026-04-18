// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
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

	s.Require().NotNil(fwd.Token)
	s.Require().EqualValues([]byte("token"), fwd.Token.TokenId)
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

	s.Require().NotNil(fwd.Token)
	s.Require().EqualValues([]byte("token"), fwd.Token.TokenId)
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

	s.Require().NotNil(fwd.Token)
	s.Require().EqualValues([]byte("token"), fwd.Token.TokenId)
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

	s.Require().NotNil(fwd.Token)
	s.Require().EqualValues([]byte("token"), fwd.Token.TokenId)
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

	s.Require().NotNil(fwd.Token)
	s.Require().EqualValues([]byte("token"), fwd.Token.TokenId)
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

	s.Require().NotNil(fwd.Token)
	s.Require().EqualValues([]byte("token"), fwd.Token.TokenId)
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

// TestForwardingByIdCorrectRecord verifies that ById returns the record matching the
// requested id even when it sits at a different time than the preceding records.
// An incorrect ORDER BY in the outer query (using Order instead of OrderExpr) would
// cause PostgreSQL to return rows in heap order, yielding id=2 instead of id=3.
func (s *StorageTestSuite) TestForwardingByIdCorrectRecord() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	// id=3 is at '2024-07-04T03:07:00' — 7 minutes after ids 1 and 2.
	// A broken ORDER BY would return id=2 here.
	fwd, prevTime, err := s.storage.Forwardings.ById(ctx, 3)
	s.Require().NoError(err)
	s.Require().EqualValues(3, fwd.Id)
	s.Require().EqualValues(10001, fwd.Height)
	// prevTime must equal the time of the preceding forwarding (id=2)
	s.Require().Equal(time.Date(2024, 7, 4, 3, 0, 0, 0, time.UTC), prevTime.UTC())
}

// TestForwardingByIdNotFound verifies that requesting a non-existent forwarding id
// returns sql.ErrNoRows instead of silently returning a record with a lower id.
func (s *StorageTestSuite) TestForwardingByIdNotFound() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	// id=4 does not exist in the forwarding fixture (only 1, 2, 3 exist).
	// Without the fwds[0].Id != id guard the query would silently return id=3.
	_, _, err := s.storage.Forwardings.ById(ctx, 4)
	s.Require().ErrorIs(err, sql.ErrNoRows)
}

// TestForwardingInputsAtSameTime verifies that an HL receive transfer whose time
// equals the upper bound is included in the inputs list.
// This requires a non-strict (<=) comparison; a strict (<) would exclude it.
func (s *StorageTestSuite) TestForwardingInputsAtSameTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	// hl_transfer fixture id=3: type=receive, address_id=5, time='2023-07-04T04:11:57', counterparty=123450.
	// Passing that exact time as `to` should include the transfer with <=, but exclude it with <.
	sameTime := time.Date(2023, 7, 4, 4, 11, 57, 0, time.UTC)
	inputs, err := s.storage.Forwardings.Inputs(ctx, 5, time.Time{}, sameTime)
	s.Require().NoError(err)

	found := false
	for _, inp := range inputs {
		if inp.Counterparty == 123450 {
			found = true
			s.Require().Equal("1000", inp.Amount)
			s.Require().Equal("utia", inp.Denom)
		}
	}
	s.Require().True(found, "HL receive transfer with time == upper bound must appear in inputs")
}
