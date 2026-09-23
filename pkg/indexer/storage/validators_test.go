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
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const testConsAddress = "A5B3C7D1E2F40912345678901234567890ABCDEF"

func jailsWith(jail *storage.Jail) *sdkSync.Map[string, *storage.Jail] {
	jails := sdkSync.NewMap[string, *storage.Jail]()
	jails.Set(jail.Validator.ConsAddress, jail)
	return jails
}

func newJail(burned int64) *storage.Jail {
	jailed := true
	return &storage.Jail{
		Height: pkgTypes.Level(100),
		Time:   time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC),
		Reason: "double_sign",
		Burned: types.NumericFromInt64(burned),
		Validator: &storage.Validator{
			ConsAddress: testConsAddress,
			Stake:       types.NumericFromInt64(-burned),
			Jailed:      &jailed,
		},
	}
}

func newStorageModule(ctrl *gomock.Controller) Module {
	module := NewModule(nil, mock.NewMockIConstant(ctrl), mock.NewMockIValidator(ctrl), nil, config.Indexer{})
	module.validatorsByConsAddress[testConsAddress] = 7
	return module
}

// A downtime jail burns nothing on Celestia, but it still has to reach the jail table.
func Test_saveValidators_ZeroBurnJailIsStored(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newStorageModule(ctrl)

	jail := newJail(0)

	tx.EXPECT().Jail(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, validators ...*storage.Validator) error {
		require.Len(t, validators, 1)
		require.EqualValues(t, 7, validators[0].Id)
		require.True(t, validators[0].Stake.IsZero())
		return nil
	})
	tx.EXPECT().SaveJails(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, jails ...storage.Jail) error {
		require.Len(t, jails, 1)
		require.EqualValues(t, 7, jails[0].ValidatorId)
		require.True(t, jails[0].Burned.IsZero())
		return nil
	})
	// no burn means nothing to take from the delegators: UpdateSlashedDelegations must not run

	count, err := module.saveValidators(t.Context(), tx, nil, jailsWith(jail))
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

// A double sign burns stake: the delegations are cut by the magnitude and the validator
// stake is moved by the negative delta.
func Test_saveValidators_SlashCutsDelegations(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newStorageModule(ctrl)

	jail := newJail(1000)
	balances := []storage.Balance{{Id: 1, Currency: "utia", Delegated: types.NumericFromInt64(-1000)}}

	tx.EXPECT().UpdateSlashedDelegations(gomock.Any(), uint64(7), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ uint64, burned types.Numeric) ([]storage.Balance, error) {
			require.Equal(t, "1000", burned.String())
			return balances, nil
		})
	tx.EXPECT().SaveBalances(gomock.Any(), balances[0]).Return(nil)
	tx.EXPECT().Jail(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, validators ...*storage.Validator) error {
		require.Len(t, validators, 1)
		require.Equal(t, "-1000", validators[0].Stake.String())
		return nil
	})
	tx.EXPECT().SaveJails(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, jails ...storage.Jail) error {
		require.Len(t, jails, 1)
		require.Equal(t, "1000", jails[0].Burned.String())
		return nil
	})

	count, err := module.saveValidators(t.Context(), tx, nil, jailsWith(jail))
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func Test_saveValidators_UnknownValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newStorageModule(ctrl)

	jail := newJail(1000)
	jail.Validator.ConsAddress = "DEAD"

	_, err := module.saveValidators(t.Context(), tx, nil, jailsWith(jail))
	require.Error(t, err)
}
