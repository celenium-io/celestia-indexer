// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestForwardingById() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	fwd, prevTime, err := s.storage.Forwardings.ById(ctx, 2)
	s.Require().NoError(err)
	s.Require().EqualValues(2, fwd.Id)
	s.Require().EqualValues(3, fwd.MsgId)
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
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
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
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
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
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
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
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
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
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
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
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
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
	s.Require().Equal("hyperlane", input1.Type)
	s.Require().Empty(input1.ChainId)
}

// TestForwardingByIdCorrectRecord verifies that ById returns the record matching the
// requested id even when it sits at a different time than the preceding records.
// An incorrect ORDER BY in the outer query (using Order instead of OrderExpr) would
// cause PostgreSQL to return rows in heap order, yielding id=2 instead of id=3.
func (s *StorageTestSuite) TestForwardingByIdCorrectRecord() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
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
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	// id=100 does not exist in the forwarding fixture (only 1..5 exist).
	// Without the fwds[0].Id != id guard the query would silently return id=5.
	_, _, err := s.storage.Forwardings.ById(ctx, 100)
	s.Require().ErrorIs(err, sql.ErrNoRows)
}

// The previous forwarding is taken from the same address: id=4 (address 2) sits between ids 3 and 5 (address 5).
func (s *StorageTestSuite) TestForwardingByIdIgnoresOtherAddresses() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	fwd, prevTime, err := s.storage.Forwardings.ById(ctx, 5)
	s.Require().NoError(err)
	s.Require().EqualValues(5, fwd.Id)
	s.Require().EqualValues(5, fwd.AddressId)
	s.Require().Equal(time.Date(2024, 7, 4, 3, 7, 0, 0, time.UTC), prevTime.UTC())
}

// The first forwarding of an address has no lower bound, even if other addresses forwarded earlier.
func (s *StorageTestSuite) TestForwardingByIdFirstOfAddress() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	fwd, prevTime, err := s.storage.Forwardings.ById(ctx, 4)
	s.Require().NoError(err)
	s.Require().EqualValues(4, fwd.Id)
	s.Require().EqualValues(2, fwd.AddressId)
	s.Require().True(prevTime.IsZero())
}

// MsgSend inputs: only successful txs inside (from, to] count.
func (s *StorageTestSuite) TestForwardingInputsMsgSend() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	// Fixtures are shared by the suite: insert in a transaction and roll it back.
	tx, err := s.storage.Connection().DB().BeginTx(ctx, nil)
	s.Require().NoError(err)
	defer tx.Rollback() //nolint:errcheck

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)

	cases := []struct {
		id     uint64
		time   time.Time
		status types.Status
	}{
		{id: 100, time: from, status: types.StatusSuccess},                      // previous forwarding's block
		{id: 101, time: from.Add(5 * time.Minute), status: types.StatusSuccess}, // the only input
		{id: 102, time: from.Add(6 * time.Minute), status: types.StatusFailed},
		{id: 103, time: to.Add(time.Minute), status: types.StatusSuccess},
	}
	for i, c := range cases {
		hash := make([]byte, 32)
		hash[31] = byte(c.id)
		_, err = tx.NewInsert().Model(&storage.Tx{
			Id:            c.id,
			Height:        pkgTypes.Level(20000 + i),
			Time:          c.time,
			MessagesCount: 1,
			Status:        c.status,
			Hash:          hash,
		}).Column("id", "height", "time", "messages_count", "status", "hash").Exec(ctx)
		s.Require().NoError(err)

		_, err = tx.NewInsert().Model(&storage.Message{
			Id:     c.id,
			Height: pkgTypes.Level(20000 + i),
			Time:   c.time,
			Type:   types.MsgSend,
			TxId:   c.id,
			Data: types.PackedBytes{
				"FromAddress": "celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8",
				"Amount":      []any{map[string]any{"Denom": "utia", "Amount": strconv.FormatUint(c.id, 10)}},
			},
		}).Exec(ctx)
		s.Require().NoError(err)

		_, err = tx.NewInsert().Model(&storage.MsgAddress{
			AddressId: 5,
			MsgId:     c.id,
			Type:      types.MsgAddressTypeToAddress,
		}).Exec(ctx)
		s.Require().NoError(err)
	}

	inputs, err := NewForwarding(tx).Inputs(ctx, 5, from, to)
	s.Require().NoError(err)
	s.Require().Len(inputs, 1)

	s.Require().EqualValues(20001, inputs[0].Height)
	s.Require().Equal("101", inputs[0].Amount)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", inputs[0].From)
	s.Require().EqualValues(101, inputs[0].TxHash[31])
	s.Require().Equal("send", inputs[0].Type)
}

// IBC inputs: received by the address inside (from, to] in successful txs, sender may be foreign or celestia.
func (s *StorageTestSuite) TestForwardingInputsIbc() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	// Fixtures are shared by the suite: insert in a transaction and roll it back.
	tx, err := s.storage.Connection().DB().BeginTx(ctx, nil)
	s.Require().NoError(err)
	defer tx.Rollback() //nolint:errcheck

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	osmoSender := "osmo1m8wg4vxkefhs374qxmmqpyusgz289wmulex5qdwpfx7jnrxzer5s9cv83q"

	cases := []struct {
		id        uint64
		time      time.Time
		status    types.Status
		receiver  uint64
		sender    *string
		senderId  *uint64
		connected bool
	}{
		{id: 200, time: from.Add(5 * time.Minute), status: types.StatusSuccess, receiver: 5, sender: &osmoSender, connected: true},
		{id: 201, time: from, status: types.StatusSuccess, receiver: 5, sender: &osmoSender},                      // previous forwarding's block
		{id: 202, time: to.Add(time.Minute), status: types.StatusSuccess, receiver: 5, sender: &osmoSender},       // after the forwarding
		{id: 203, time: from.Add(6 * time.Minute), status: types.StatusSuccess, receiver: 2, sender: &osmoSender}, // other address
		{id: 204, time: from.Add(7 * time.Minute), status: types.StatusFailed, receiver: 5, sender: &osmoSender},
		{id: 205, time: from.Add(8 * time.Minute), status: types.StatusSuccess, receiver: 5, senderId: testsuite.Ptr(uint64(1))},
	}
	for i, c := range cases {
		hash := make([]byte, 32)
		hash[31] = byte(c.id)
		_, err = tx.NewInsert().Model(&storage.Tx{
			Id:            c.id,
			Height:        pkgTypes.Level(30000 + i),
			Time:          c.time,
			MessagesCount: 1,
			Status:        c.status,
			Hash:          hash,
		}).Column("id", "height", "time", "messages_count", "status", "hash").Exec(ctx)
		s.Require().NoError(err)

		connection := "connection-unknown"
		if c.connected {
			connection = "connection-1"
		}
		_, err = tx.NewInsert().Model(&storage.IbcTransfer{
			Id:            c.id,
			Time:          c.time,
			Height:        pkgTypes.Level(30000 + i),
			Amount:        types.NumericFromInt64(int64(c.id)),
			Denom:         "utia",
			ReceiverId:    testsuite.Ptr(c.receiver),
			SenderAddress: c.sender,
			SenderId:      c.senderId,
			ConnectionId:  connection,
			ChannelId:     "channel-1",
			Port:          "transfer",
			TxId:          c.id,
		}).Exec(ctx)
		s.Require().NoError(err)
	}

	inputs, err := NewForwarding(tx).Inputs(ctx, 5, from, to)
	s.Require().NoError(err)
	s.Require().Len(inputs, 2)

	// ordered by time desc
	celestiaSender := inputs[0]
	s.Require().EqualValues(30005, celestiaSender.Height)
	s.Require().Equal("ibc", celestiaSender.Type)
	s.Require().Equal("205", celestiaSender.Amount)
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", celestiaSender.From)
	s.Require().Empty(celestiaSender.ChainId, "unknown connection")

	foreignSender := inputs[1]
	s.Require().EqualValues(30000, foreignSender.Height)
	s.Require().Equal("ibc", foreignSender.Type)
	s.Require().Equal("200", foreignSender.Amount)
	s.Require().Equal("utia", foreignSender.Denom)
	s.Require().Equal(osmoSender, foreignSender.From)
	s.Require().Equal("osmosis-1", foreignSender.ChainId)
	s.Require().Equal("channel-1", foreignSender.ChannelId)
	s.Require().EqualValues(200, foreignSender.TxHash[31])
	s.Require().Zero(foreignSender.Counterparty)
}

// A transfer at the previous forwarding's time belongs to that forwarding, so the lower bound is strict.
func (s *StorageTestSuite) TestForwardingInputsExcludesLowerBound() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
	defer ctxCancel()

	// hl_transfer fixture id=3: type=receive, address_id=5, time='2023-07-04T04:11:57'
	sameTime := time.Date(2023, 7, 4, 4, 11, 57, 0, time.UTC)
	inputs, err := s.storage.Forwardings.Inputs(ctx, 5, sameTime, time.Time{})
	s.Require().NoError(err)

	for _, inp := range inputs {
		s.Require().NotEqualValues(123450, inp.Counterparty, "HL receive transfer with time == lower bound must not appear in inputs")
	}
}

// TestForwardingInputsAtSameTime verifies that an HL receive transfer whose time
// equals the upper bound is included in the inputs list.
// This requires a non-strict (<=) comparison; a strict (<) would exclude it.
func (s *StorageTestSuite) TestForwardingInputsAtSameTime() {
	ctx, ctxCancel := context.WithTimeout(s.T().Context(), 5*time.Second)
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
