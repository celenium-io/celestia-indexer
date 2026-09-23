// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	st "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// rollbackTx wires the calls rollbackValidators always makes and captures what it writes back.
type rollbackTx struct {
	tx          *mock.MockTransaction
	validators  []*storage.Validator
	balances    []storage.Balance
	delegations []storage.Delegation
}

func newRollbackTx(t *testing.T, jails []storage.Jail, logs []storage.StakingLog) *rollbackTx {
	ctrl := gomock.NewController(t)
	tx := mock.NewMockTransaction(ctrl)
	captured := &rollbackTx{tx: tx}

	tx.EXPECT().RollbackValidators(gomock.Any(), gomock.Any()).Return(nil, nil)
	tx.EXPECT().RollbackUndelegations(gomock.Any(), gomock.Any()).Return(nil)
	tx.EXPECT().RollbackRedelegations(gomock.Any(), gomock.Any()).Return(nil)
	tx.EXPECT().RollbackJails(gomock.Any(), gomock.Any()).Return(jails, nil)
	tx.EXPECT().RollbackStakingLogs(gomock.Any(), gomock.Any()).Return(logs, nil)

	tx.EXPECT().UpdateValidators(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, validators ...*storage.Validator) error {
			captured.validators = validators
			return nil
		}).AnyTimes()
	tx.EXPECT().SaveBalances(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, balances ...storage.Balance) error {
			captured.balances = balances
			return nil
		}).AnyTimes()
	tx.EXPECT().SaveDelegations(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, delegations ...storage.Delegation) error {
			captured.delegations = delegations
			return nil
		}).AnyTimes()

	return captured
}

func (r *rollbackTx) validator(t *testing.T, id uint64) *storage.Validator {
	t.Helper()
	for i := range r.validators {
		if r.validators[i].Id == id {
			return r.validators[i]
		}
	}
	require.Failf(t, "validator was not updated", "id %d", id)
	return nil
}

// Undelegate took the stake away, so its rollback has to put it back. Both the first log of
// a validator and every next one must move the stake in the same direction.
func Test_rollbackValidators_UnbondingRestoresStake(t *testing.T) {
	addressId := uint64(1)
	logs := []storage.StakingLog{
		{ValidatorId: 5, AddressId: &addressId, Change: st.NumericFromInt64(-300), Type: st.StakingLogTypeUnbonding},
		{ValidatorId: 5, AddressId: &addressId, Change: st.NumericFromInt64(-200), Type: st.StakingLogTypeUnbonding},
	}
	captured := newRollbackTx(t, nil, logs)

	result, err := rollbackValidators(t.Context(), captured.tx, types.Level(100))
	require.NoError(t, err)
	require.Equal(t, 0, result.count)

	require.Len(t, captured.validators, 1)
	require.Equal(t, "500", captured.validator(t, 5).Stake.String())

	// the delegator gets the stake back and leaves the unbonding queue
	require.Len(t, captured.balances, 1)
	require.Equal(t, "500", captured.balances[0].Delegated.String())
	require.Equal(t, "-500", captured.balances[0].Unbonding.String())

	require.Len(t, captured.delegations, 1)
	require.Equal(t, "500", captured.delegations[0].Amount.String())
}

// A cancelled unbonding is stored as a positive unbonding log, so its rollback runs the
// other way: the stake given back to the validator has to be taken away again.
func Test_rollbackValidators_CancelledUnbondingTakesStake(t *testing.T) {
	addressId := uint64(1)
	logs := []storage.StakingLog{
		{ValidatorId: 5, AddressId: &addressId, Change: st.NumericFromInt64(300), Type: st.StakingLogTypeUnbonding},
	}
	captured := newRollbackTx(t, nil, logs)

	_, err := rollbackValidators(t.Context(), captured.tx, types.Level(100))
	require.NoError(t, err)
	require.Equal(t, "-300", captured.validator(t, 5).Stake.String())
}

func Test_rollbackValidators_DelegationTakesStake(t *testing.T) {
	addressId := uint64(1)
	logs := []storage.StakingLog{
		{ValidatorId: 5, AddressId: &addressId, Change: st.NumericFromInt64(300), Type: st.StakingLogTypeDelegation},
		{ValidatorId: 5, AddressId: &addressId, Change: st.NumericFromInt64(200), Type: st.StakingLogTypeDelegation},
	}
	captured := newRollbackTx(t, nil, logs)

	_, err := rollbackValidators(t.Context(), captured.tx, types.Level(100))
	require.NoError(t, err)
	require.Equal(t, "-500", captured.validator(t, 5).Stake.String())
	require.Equal(t, "-500", captured.balances[0].Delegated.String())
	require.Equal(t, "-500", captured.delegations[0].Amount.String())
}

// Only a rolled back jail carries a flag. Validators that got into the update through their
// rewards logs leave it nil, and UpdateValidators keeps whatever the column holds.
func Test_rollbackValidators_JailedFlagOnlyFromJails(t *testing.T) {
	jails := []storage.Jail{{ValidatorId: 5}}
	logs := []storage.StakingLog{
		{ValidatorId: 5, Change: st.NumericFromInt64(10), Type: st.StakingLogTypeRewards},
		{ValidatorId: 7, Change: st.NumericFromInt64(20), Type: st.StakingLogTypeRewards},
		{ValidatorId: 7, Change: st.NumericFromInt64(3), Type: st.StakingLogTypeCommissions},
	}
	captured := newRollbackTx(t, jails, logs)

	_, err := rollbackValidators(t.Context(), captured.tx, types.Level(100))
	require.NoError(t, err)
	require.Len(t, captured.validators, 2)

	jailed := captured.validator(t, 5)
	require.NotNil(t, jailed.Jailed)
	require.False(t, *jailed.Jailed)
	require.Equal(t, "-10", jailed.Rewards.String())

	fromLogs := captured.validator(t, 7)
	require.Nil(t, fromLogs.Jailed)
	require.Equal(t, "-20", fromLogs.Rewards.String())
	require.Equal(t, "-3", fromLogs.Commissions.String())
}
