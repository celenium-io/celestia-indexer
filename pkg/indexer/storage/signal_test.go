// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPostProcessingSignal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)
	constants := mock.NewMockIConstant(ctrl)
	validators := mock.NewMockIValidator(ctrl)

	module := NewModule(nil, constants, validators, nil, config.Indexer{})

	signals := []*storage.SignalVersion{
		{
			Height:      100,
			Time:        time.Now(),
			Version:     11,
			TxId:        1,
			MsgId:       1,
			ValidatorId: 1,
			VotingPower: decimal.RequireFromString("123456"),
		},
	}
	upgrades := sync.NewMap[uint64, *storage.Upgrade]()
	upgrades.Set(11, &storage.Upgrade{
		Version: 11,
	})

	constants.
		EXPECT().
		Get(gomock.Any(), types.ModuleNameStaking, "max_validators").
		Return(storage.Constant{
			Value: "100",
		}, nil).
		Times(1)

	tx.EXPECT().
		BondedValidators(gomock.Any(), 100).
		Return([]storage.Validator{
			{
				Id:      1,
				Version: 11,
				Stake:   decimal.RequireFromString("123456"),
			},
		}, nil).
		Times(1)

	tx.EXPECT().
		UpdateSignalsAfterUpgrade(gomock.Any(), uint64(11)).
		Return(nil).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := module.postProcessingSignal(ctx, tx, signals, upgrades)
	require.NoError(t, err)
}
