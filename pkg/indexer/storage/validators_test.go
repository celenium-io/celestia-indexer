// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
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

const testValAddress = "celestiavaloper17vmk8m246t648hpmde2q7kp4ft9uwrayy09dmw"

func powerUpdate(consAddress string, power int64) storage.ValidatorBondUpdate {
	val := storage.EmptyValidator()
	val.ConsAddress = consAddress
	val.BondUpdatesCount = 1
	p := types.NumericFromInt64(power)
	val.Power = &p
	return storage.ValidatorBondUpdate{
		Validator: &val,
		Power:     &p,
		Height:    1000,
		Time:      time.Now(),
	}
}

func newPowerModule(ctrl *gomock.Controller) Module {
	module := newStorageModule(ctrl)
	module.validatorsByAddress[testValAddress] = 7
	return module
}

// A power update carries only the consensus address; id and operator address come from the cache.
func Test_processValidatorUpdates_ResolvesValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newPowerModule(ctrl)

	dCtx := decodeContext.NewContext()
	dCtx.AddValidatorUpdate(powerUpdate(testConsAddress, 42))

	require.NoError(t, module.processValidatorBondUpdates(t.Context(), tx, dCtx))

	val, ok := dCtx.Validators.Get(testValAddress)
	require.True(t, ok)
	require.EqualValues(t, 7, val.Id)
	require.Equal(t, testConsAddress, val.ConsAddress)
	require.NotNil(t, val.Power)
	require.Equal(t, "42", val.Power.String())
}

// The update merges into the validator already touched in this block instead of replacing it.
func Test_processValidatorUpdates_MergesWithBlockValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newPowerModule(ctrl)

	dCtx := decodeContext.NewContext()
	touched := storage.EmptyValidator()
	touched.Address = testValAddress
	touched.Stake = types.NumericFromInt64(1000)
	dCtx.AddValidator(touched)
	dCtx.AddValidatorUpdate(powerUpdate(testConsAddress, 3))

	require.NoError(t, module.processValidatorBondUpdates(t.Context(), tx, dCtx))

	require.Equal(t, 1, dCtx.Validators.Len())
	val, ok := dCtx.Validators.Get(testValAddress)
	require.True(t, ok)
	require.Equal(t, "1000", val.Stake.String())
	require.Equal(t, "3", val.Power.String())
}

// MsgCreateValidator with enough stake bonds the validator in the same block's EndBlock.
func Test_processValidatorUpdates_ValidatorCreatedInBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newStorageModule(ctrl)

	const newConsAddress = "AE216C2EF5247A3782C135EFA279A3E4CDC61094"
	const newAddress = "celestiavaloper1new"

	dCtx := decodeContext.NewContext()
	created := storage.EmptyValidator()
	created.Address = newAddress
	created.ConsAddress = newConsAddress
	dCtx.AddValidator(created)
	dCtx.AddValidatorUpdate(powerUpdate(newConsAddress, 5))

	require.NoError(t, module.processValidatorBondUpdates(t.Context(), tx, dCtx))

	require.Equal(t, 1, dCtx.Validators.Len())
	val, ok := dCtx.Validators.Get(newAddress)
	require.True(t, ok)
	require.NotNil(t, val.Power)
	require.Equal(t, "5", val.Power.String())
}

func Test_processValidatorUpdates_UnknownValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newPowerModule(ctrl)

	tx.EXPECT().GetProposerId(gomock.Any(), "DEAD").Return(uint64(0), sql.ErrNoRows).Times(1)

	dCtx := decodeContext.NewContext()
	dCtx.AddValidatorUpdate(powerUpdate("DEAD", 5))

	require.Error(t, module.processValidatorBondUpdates(t.Context(), tx, dCtx))
}

func Test_processValidatorUpdates_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newPowerModule(ctrl)

	dCtx := decodeContext.NewContext()
	require.NoError(t, module.processValidatorBondUpdates(t.Context(), tx, dCtx))
	require.Equal(t, 0, dCtx.Validators.Len())
}

// The bond update counter reaches the validator row only once per block.
func Test_processValidatorUpdates_BondUpdatesCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newPowerModule(ctrl)

	dCtx := decodeContext.NewContext()
	dCtx.AddValidatorUpdate(powerUpdate(testConsAddress, 3))
	dCtx.AddValidatorUpdate(powerUpdate(testConsAddress, 0))

	require.NoError(t, module.processValidatorBondUpdates(t.Context(), tx, dCtx))

	val, ok := dCtx.Validators.Get(testValAddress)
	require.True(t, ok)
	require.EqualValues(t, 1, val.BondUpdatesCount)
	require.True(t, val.Power.IsZero())
}

// The cache knows the consensus address but lost the operator address: it is loaded from the DB.
func Test_processValidatorUpdates_NoOperatorAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newStorageModule(ctrl)

	tx.EXPECT().GetProposerId(gomock.Any(), testConsAddress).Return(uint64(7), nil).Times(1)
	tx.EXPECT().Validator(gomock.Any(), uint64(7)).Return(storage.Validator{
		Id:          7,
		Address:     testValAddress,
		ConsAddress: testConsAddress,
	}, nil).Times(1)

	dCtx := decodeContext.NewContext()
	dCtx.AddValidatorUpdate(powerUpdate(testConsAddress, 5))

	require.NoError(t, module.processValidatorBondUpdates(t.Context(), tx, dCtx))

	val, ok := dCtx.Validators.Get(testValAddress)
	require.True(t, ok)
	require.EqualValues(t, 7, val.Id)
	require.Equal(t, "5", val.Power.String())
}

// On a sync from scratch genesis validators are not cached; a delegation in the same block
// puts the validator into the context without a consensus address.
func Test_processValidatorUpdates_GenesisValidatorNotCached(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := NewModule(nil, mock.NewMockIConstant(ctrl), mock.NewMockIValidator(ctrl), nil, config.Indexer{})

	tx.EXPECT().GetProposerId(gomock.Any(), testConsAddress).Return(uint64(7), nil).Times(1)
	tx.EXPECT().Validator(gomock.Any(), uint64(7)).Return(storage.Validator{
		Id:          7,
		Address:     testValAddress,
		ConsAddress: testConsAddress,
	}, nil).Times(1)

	dCtx := decodeContext.NewContext()
	delegated := storage.EmptyValidator()
	delegated.Address = testValAddress
	delegated.Stake = types.NumericFromInt64(1000)
	dCtx.AddValidator(delegated)
	dCtx.AddValidatorUpdate(powerUpdate(testConsAddress, 5))

	require.NoError(t, module.processValidatorBondUpdates(t.Context(), tx, dCtx))

	require.Equal(t, 1, dCtx.Validators.Len())
	val, ok := dCtx.Validators.Get(testValAddress)
	require.True(t, ok)
	require.Equal(t, "1000", val.Stake.String())
	require.Equal(t, "5", val.Power.String())

	// the loaded validator is cached, so the bond update is saved without a second lookup
	tx.EXPECT().SaveBondUpdates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, items ...*storage.ValidatorBondUpdate) error {
			require.Len(t, items, 1)
			require.EqualValues(t, 7, items[0].ValidatorId)
			return nil
		}).Times(1)
	require.NoError(t, module.saveValidatorBondUpdates(t.Context(), tx, dCtx.ValidatorUpdates))
}

func bondUpdatesWith(updates ...storage.ValidatorBondUpdate) *sdkSync.Map[string, *storage.ValidatorBondUpdate] {
	m := sdkSync.NewMap[string, *storage.ValidatorBondUpdate]()
	for i := range updates {
		m.Set(updates[i].Validator.ConsAddress, &updates[i])
	}
	return m
}

func Test_saveValidatorBondUpdates(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newPowerModule(ctrl)
	module.validatorsByConsAddress["B0B0"] = 8

	tx.EXPECT().SaveBondUpdates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, items ...*storage.ValidatorBondUpdate) error {
			require.Len(t, items, 2)
			ids := make(map[string]uint64)
			for _, item := range items {
				require.NotNil(t, item.Power)
				require.EqualValues(t, 1000, item.Height)
				ids[item.Validator.ConsAddress] = item.ValidatorId
			}
			require.Equal(t, map[string]uint64{testConsAddress: 7, "B0B0": 8}, ids)
			return nil
		}).Times(1)

	updates := bondUpdatesWith(powerUpdate(testConsAddress, 10), powerUpdate("B0B0", 0))
	require.NoError(t, module.saveValidatorBondUpdates(t.Context(), tx, updates))
}

func Test_saveValidatorBondUpdates_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newPowerModule(ctrl)

	require.NoError(t, module.saveValidatorBondUpdates(t.Context(), tx, bondUpdatesWith()))
}

func Test_saveValidatorBondUpdates_UnknownValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newPowerModule(ctrl)

	err := module.saveValidatorBondUpdates(t.Context(), tx, bondUpdatesWith(powerUpdate("DEAD", 1)))
	require.Error(t, err)
}

// A validator created and bonded in one block gets its id from saveValidators before the bond update is saved.
func Test_saveValidatorBondUpdates_ValidatorCreatedInBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	module := newStorageModule(ctrl)

	const newConsAddress = "AE216C2EF5247A3782C135EFA279A3E4CDC61094"

	dCtx := decodeContext.NewContext()
	created := storage.EmptyValidator()
	created.Address = "celestiavaloper1new"
	created.ConsAddress = newConsAddress
	dCtx.AddValidator(created)
	dCtx.AddValidatorUpdate(powerUpdate(newConsAddress, 5))
	require.NoError(t, module.processValidatorBondUpdates(t.Context(), tx, dCtx))

	tx.EXPECT().SaveValidators(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, validators ...*storage.Validator) (int, error) {
			require.Len(t, validators, 1)
			require.EqualValues(t, 1, validators[0].BondUpdatesCount)
			validators[0].Id = 42
			return 1, nil
		}).Times(1)
	tx.EXPECT().SaveBondUpdates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, items ...*storage.ValidatorBondUpdate) error {
			require.Len(t, items, 1)
			require.EqualValues(t, 42, items[0].ValidatorId)
			return nil
		}).Times(1)

	_, err := module.saveValidators(t.Context(), tx, dCtx.Validators.Values(), sdkSync.NewMap[string, *storage.Jail]())
	require.NoError(t, err)
	require.NoError(t, module.saveValidatorBondUpdates(t.Context(), tx, dCtx.ValidatorUpdates))
}
