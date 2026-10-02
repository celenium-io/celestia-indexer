// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// The outgoing transfer of a MsgForward has a sender but no relayer.
func TestSaveHlTransfers_Forward(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	mailboxRepo := mock.NewMockIHLMailbox(ctrl)
	tokenRepo := mock.NewMockIHLToken(ctrl)
	igpRepo := mock.NewMockIHLIGP(ctrl)

	const forwardAddr = "celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60"
	var (
		mailbox = util.CreateMockHexAddress("mailbox", 1)
		token   = util.CreateMockHexAddress("token", 1)
		igp     = util.CreateMockHexAddress("igp", 1)
	)

	transfer := &storage.HLTransfer{
		Height:  100,
		TxId:    7,
		Type:    types.HLTransferTypeSend,
		Address: &storage.Address{Address: forwardAddr},
		Mailbox: &storage.HLMailbox{InternalId: mailbox.GetInternalId()},
		Token:   &storage.HLToken{TokenId: token.Bytes()},
		GasPayment: &storage.HLGasPayment{
			Height: 100,
			Amount: types.NumericFromInt64(4379232),
			Igp:    &storage.HLIGP{IgpId: igp.Bytes()},
		},
	}

	mailboxRepo.EXPECT().
		ByInternalId(gomock.Any(), mailbox.GetInternalId()).
		Return(storage.HLMailbox{Id: 2}, nil).
		Times(1)
	tokenRepo.EXPECT().
		IdByTokenId(gomock.Any(), token.Bytes()).
		Return(uint64(3), nil).
		Times(1)
	igpRepo.EXPECT().
		IdByHash(gomock.Any(), igp.Bytes()).
		Return(uint64(4), nil).
		Times(1)

	tx.EXPECT().
		SaveHyperlaneTransfers(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, transfers ...*storage.HLTransfer) error {
			require.Len(t, transfers, 1)
			require.EqualValues(t, 5, transfers[0].AddressId)
			require.Zero(t, transfers[0].RelayerId)
			require.EqualValues(t, 2, transfers[0].MailboxId)
			require.EqualValues(t, 3, transfers[0].TokenId)
			transfers[0].Id = 11 // filled by RETURNING
			return nil
		}).
		Times(1)
	tx.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, models any) error {
			payments, ok := models.(*[]*storage.HLGasPayment)
			require.True(t, ok)
			require.Len(t, *payments, 1)
			require.EqualValues(t, 11, (*payments)[0].TransferId)
			require.EqualValues(t, 4, (*payments)[0].IgpId)
			return nil
		}).
		Times(1)

	err := saveHlTransfers(t.Context(), tx, mailboxRepo, tokenRepo, igpRepo, []*storage.HLTransfer{transfer}, map[string]uint64{forwardAddr: 5})
	require.NoError(t, err)
}
