// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"
)

func (s *StorageTestSuite) TestHyperlaneMailboxByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	mailbox, err := s.storage.HLMailbox.ByHash(ctx, []byte("mailbox"))
	s.Require().NoError(err)

	s.Require().EqualValues(1, mailbox.Id)
	s.Require().EqualValues(1000, mailbox.Height)
	s.Require().EqualValues(12, mailbox.SentMessages)
	s.Require().EqualValues(21, mailbox.ReceivedMessages)
	s.Require().EqualValues(123, mailbox.Domain)
	s.Require().EqualValues([]byte("mailbox"), mailbox.Mailbox)
	s.Require().NotNil(mailbox.DefaultHook)
	s.Require().NotNil(mailbox.DefaultIsm)
	s.Require().NotNil(mailbox.RequiredHook)

	s.Require().NotNil(mailbox.Tx)
	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, mailbox.Tx.Hash)

	s.Require().NotNil(mailbox.Owner)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", mailbox.Owner.Address)
	s.Require().NotNil(mailbox.Owner.Celestials)
}

func (s *StorageTestSuite) TestHyperlaneMailboxList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.HLMailbox.List(ctx, 1, 0)
	s.Require().NoError(err)
	s.Require().Len(items, 1)

	mailbox := items[0]
	s.Require().EqualValues(1, mailbox.Id)
	s.Require().EqualValues(1000, mailbox.Height)
	s.Require().EqualValues(12, mailbox.SentMessages)
	s.Require().EqualValues(21, mailbox.ReceivedMessages)
	s.Require().EqualValues(123, mailbox.Domain)
	s.Require().EqualValues([]byte("mailbox"), mailbox.Mailbox)
	s.Require().NotNil(mailbox.DefaultHook)
	s.Require().NotNil(mailbox.DefaultIsm)
	s.Require().NotNil(mailbox.RequiredHook)

	s.Require().NotNil(mailbox.Tx)
	txHash, err := hex.DecodeString("652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF")
	s.Require().NoError(err)
	s.Require().Equal(txHash, mailbox.Tx.Hash)

	s.Require().NotNil(mailbox.Owner)
	s.Require().EqualValues("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", mailbox.Owner.Address)
	s.Require().NotNil(mailbox.Owner.Celestials)
}
