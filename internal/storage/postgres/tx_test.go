// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestTxByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)

	tx, err := s.storage.Tx.ByHash(ctx, txHash)
	s.Require().NoError(err)

	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(0, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(2, tx.MessagesCount)
	s.Require().Equal(txHash, tx.Hash)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo", tx.Memo)
	s.Require().Equal("sdk", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())

	s.Require().Len(tx.Signers, 1)
}

func (s *StorageTestSuite) TestTxIdByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)

	id, err := s.storage.Tx.IdByHash(ctx, txHash)
	s.Require().NoError(err)
	s.Require().EqualValues(1, id)
}

func (s *StorageTestSuite) TestTxFilterSuccessUnjailAsc() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Sort:         sdk.SortOrderAsc,
		Limit:        10,
		Offset:       0,
		MessageTypes: types.NewMsgTypeBitMask(types.MsgUnjail),
		Status:       []string{string(types.StatusSuccess)},
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]

	s.Require().EqualValues(2, tx.Id)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().EqualValues("2048", tx.MessageTypes.Bits.String())
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo2", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())

	s.Require().Len(tx.Signers, 2)
}

func (s *StorageTestSuite) TestTxFilterExcludedMessageTypes() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Sort:                 sdk.SortOrderAsc,
		Limit:                10,
		Offset:               0,
		ExcludedMessageTypes: types.NewMsgTypeBitMask(types.MsgUnjail),
		Height:               testsuite.Ptr(uint64(1000)),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]

	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(0, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(2, tx.MessagesCount)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo", tx.Memo)
	s.Require().Equal("sdk", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())

	s.Require().Len(tx.Signers, 1)
}

func (s *StorageTestSuite) TestTxFilterSuccessDesc() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Sort:   sdk.SortOrderDesc,
		Limit:  10,
		Offset: 0,
		Status: []string{string(types.StatusSuccess)},
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 4)

	tx := txs[1]

	s.Require().EqualValues(3, tx.Id)
	s.Require().EqualValues(0, tx.Position)
	s.Require().EqualValues(999, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(0, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().EqualValues("32", tx.MessageTypes.Bits.String())
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())

	s.Require().Len(tx.Signers, 1)
}

func (s *StorageTestSuite) TestTxFilterHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Sort:   sdk.SortOrderDesc,
		Limit:  10,
		Offset: 0,
		Status: []string{string(types.StatusSuccess)},
		Height: testsuite.Ptr(uint64(1000)),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 2)

	tx := txs[0]

	s.Require().EqualValues(2, tx.Id)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().EqualValues("2048", tx.MessageTypes.Bits.String())
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo2", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())

	s.Require().Len(tx.Signers, 2)
}

func (s *StorageTestSuite) TestTxFilterTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Limit:    10,
		TimeFrom: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 4)

	txs, err = s.storage.Tx.Filter(ctx, storage.TxFilter{
		Limit:  10,
		TimeTo: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 0)

	txs, err = s.storage.Tx.Filter(ctx, storage.TxFilter{
		Limit: 10,

		TimeFrom: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
		TimeTo:   time.Date(2023, 7, 5, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 4)
}

func (s *StorageTestSuite) TestTxFilterWithRelations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Filter(ctx, storage.TxFilter{
		Limit:        1,
		WithMessages: true,
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().Len(tx.Messages, 2)
	s.Require().EqualValues(1, tx.Messages[0].Id)
	s.Require().EqualValues(2, tx.Messages[1].Id)
}

func (s *StorageTestSuite) TestTxByIdWithRelations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := s.storage.Tx.ByIdWithRelations(ctx, 2)
	s.Require().NoError(err)

	s.Require().EqualValues(2, tx.Id)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(1000, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(80410, tx.GasWanted)
	s.Require().EqualValues(77483, tx.GasUsed)
	s.Require().EqualValues(1, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("memo2", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("80410", tx.Fee.String())
	s.Require().EqualValues("2048", tx.MessageTypes.Bits.String())

	s.Require().Len(tx.Messages, 2)
	s.Require().Len(tx.Signers, 2)
}

func (s *StorageTestSuite) TestTxGenesis() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.Genesis(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(4, tx.Id)
	s.Require().EqualValues(0, tx.Position)
	s.Require().EqualValues(0, tx.Height)
	s.Require().EqualValues(0, tx.TimeoutHeight)
	s.Require().EqualValues(0, tx.GasWanted)
	s.Require().EqualValues(0, tx.GasUsed)
	s.Require().EqualValues(0, tx.EventsCount)
	s.Require().EqualValues(1, tx.MessagesCount)
	s.Require().Equal(types.StatusSuccess, tx.Status)
	s.Require().Equal("34499b1ac473fbb03894c883178ecc83f0d6eaf6@64.227.18.169:26656", tx.Memo)
	s.Require().Equal("", tx.Codespace)
	s.Require().Equal("0", tx.Fee.String())
	s.Require().EqualValues("32", tx.MessageTypes.Bits.String())
}

func (s *StorageTestSuite) TestTxByAddressAndTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.storage.Tx.ByAddress(ctx, 1, storage.TxFilter{
		Limit:    10,
		TimeFrom: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 3)

	txs, err = s.storage.Tx.ByAddress(ctx, 1, storage.TxFilter{
		Limit:  10,
		TimeTo: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 0)

	txs, err = s.storage.Tx.ByAddress(ctx, 1, storage.TxFilter{
		Limit: 10,

		TimeFrom: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
		TimeTo:   time.Date(2023, 7, 5, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 3)

	txs, err = s.storage.Tx.ByAddress(ctx, 1, storage.TxFilter{
		Limit:  10,
		Offset: 1,

		TimeFrom: time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
		TimeTo:   time.Date(2023, 7, 5, 0, 0, 0, 0, time.UTC),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 2)
}

func (s *StorageTestSuite) TestTxGas() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	ts, err := time.Parse(time.RFC3339, "2023-07-04T03:10:57+00:00")
	s.Require().NoError(err)

	txs, err := s.storage.Tx.Gas(ctx, 1000, ts)
	s.Require().NoError(err)
	s.Require().Len(txs, 2)

	tx0 := txs[0]
	s.Require().EqualValues(80410, tx0.GasWanted)
	s.Require().EqualValues(77483, tx0.GasUsed)
	s.Require().EqualValues("80410", tx0.Fee.String())
	s.Require().EqualValues("1", tx0.GasPrice.String())

	tx1 := txs[1]
	s.Require().EqualValues(80410, tx1.GasWanted)
	s.Require().EqualValues(77483, tx1.GasUsed)
	s.Require().EqualValues("80410", tx1.Fee.String())
	s.Require().EqualValues("1", tx1.GasPrice.String())
}
