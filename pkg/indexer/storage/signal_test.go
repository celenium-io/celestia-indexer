// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPostProcessingSignal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)

	votingPower := decimal.RequireFromString("123456")
	signals := []*storage.SignalVersion{
		{
			Height:      100,
			Time:        time.Now(),
			Version:     11,
			TxId:        1,
			MsgId:       1,
			ValidatorId: 1,
			VotingPower: votingPower,
		},
	}
	upgrades := sync.NewMap[uint64, *storage.Upgrade]()
	upgrades.Set(11, &storage.Upgrade{
		Version: 11,
	})
	validators := []storage.Validator{
		{
			Id:      1,
			Stake:   votingPower,
			Version: 11,
		},
	}

	tx.EXPECT().
		UpdateSignalsAfterUpgrade(gomock.Any(), uint64(11)).
		Return(nil).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := postProcessingSignal(ctx, tx, signals, upgrades, votingPower, validators)
	require.NoError(t, err)
}
