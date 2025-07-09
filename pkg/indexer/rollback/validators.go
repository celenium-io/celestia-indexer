// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"
	"fmt"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	st "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

type rollbackedValidators struct {
	count int
	stake decimal.Decimal
}

func rollbackValidators(
	ctx context.Context,
	tx storage.Transaction,
	height types.Level,
) (result rollbackedValidators, err error) {
	removedValidators, err := tx.RollbackValidators(ctx, height)
	if err != nil {
		return result, err
	}
	result.count = len(removedValidators)

	var (
		removedIds    = make([]uint64, len(removedValidators))
		mapRemovedIds = make(map[uint64]struct{}, 0)
		updated       = make(map[uint64]*storage.Validator)
		balances      = make(map[uint64]*storage.Balance)
		delegations   = make(map[string]*storage.Delegation)
	)

	for i := range removedValidators {
		removedIds[i] = removedValidators[i].Id
		mapRemovedIds[removedValidators[i].Id] = struct{}{}
	}

	if err := tx.RollbackUndelegations(ctx, height); err != nil {
		return result, err
	}
	if err := tx.RollbackRedelegations(ctx, height); err != nil {
		return result, err
	}
	jails, err := tx.RollbackJails(ctx, height)
	if err != nil {
		return result, err
	}
	for i := range jails {
		jailed := false
		if val, ok := updated[jails[i].ValidatorId]; ok {
			val.Jailed = &jailed
		} else {
			updated[jails[i].ValidatorId] = &storage.Validator{
				Id:     jails[i].ValidatorId,
				Jailed: &jailed,
			}
		}
	}

	logs, err := tx.RollbackStakingLogs(ctx, height)
	if err != nil {
		return result, err
	}

	for i := range logs {
		_, removed := mapRemovedIds[logs[i].ValidatorId]

		switch logs[i].Type {
		case st.StakingLogTypeDelegation:
			if logs[i].AddressId != nil {
				addressId := *logs[i].AddressId
				if val, ok := balances[addressId]; ok {
					val.Delegated = val.Delegated.Sub(logs[i].Change)
				} else {
					balances[addressId] = &storage.Balance{
						Id:        addressId,
						Delegated: logs[i].Change.Copy().Neg(),
						Currency:  currency.DefaultCurrency,
					}
				}

				if !removed {
					dId := delegationId(addressId, logs[i].ValidatorId)
					if val, ok := delegations[dId]; ok {
						val.Amount = val.Amount.Sub(logs[i].Change.Copy())
					} else {
						delegations[dId] = &storage.Delegation{
							AddressId:   addressId,
							ValidatorId: logs[i].ValidatorId,
							Amount:      logs[i].Change.Copy().Neg(),
						}
					}
				}
			}

			if !removed {
				if val, ok := updated[logs[i].ValidatorId]; ok {
					val.Stake = val.Stake.Sub(logs[i].Change)
				} else {
					updated[logs[i].ValidatorId] = &storage.Validator{
						Id:    logs[i].ValidatorId,
						Stake: logs[i].Change.Copy().Neg(),
					}
				}
			}

			result.stake = result.stake.Sub(logs[i].Change)

		case st.StakingLogTypeUnbonding:
			if logs[i].AddressId != nil {
				addressId := *logs[i].AddressId
				if val, ok := balances[addressId]; ok {
					val.Delegated = val.Delegated.Sub(logs[i].Change)
					val.Unbonding = val.Unbonding.Add(logs[i].Change)
				} else {
					balances[addressId] = &storage.Balance{
						Id:        addressId,
						Delegated: logs[i].Change.Copy().Neg(),
						Unbonding: logs[i].Change.Copy(),
						Currency:  currency.DefaultCurrency,
					}
				}

				if !removed {
					dId := delegationId(addressId, logs[i].ValidatorId)
					if val, ok := delegations[dId]; ok {
						val.Amount = val.Amount.Sub(logs[i].Change.Copy())
					} else {
						delegations[dId] = &storage.Delegation{
							AddressId:   addressId,
							ValidatorId: logs[i].ValidatorId,
							Amount:      logs[i].Change.Copy().Neg(),
						}
					}
				}
			}

			if !removed {
				if val, ok := updated[logs[i].ValidatorId]; ok {
					val.Stake = val.Stake.Add(logs[i].Change)
				} else {
					updated[logs[i].ValidatorId] = &storage.Validator{
						Id:    logs[i].ValidatorId,
						Stake: logs[i].Change.Copy(),
					}
				}
			}

			result.stake = result.stake.Add(logs[i].Change)

		case st.StakingLogTypeCommissions:
			if removed {
				continue
			}
			if val, ok := updated[logs[i].ValidatorId]; ok {
				val.Commissions = val.Commissions.Sub(logs[i].Change)
			} else {
				updated[logs[i].ValidatorId] = &storage.Validator{
					Id:          logs[i].ValidatorId,
					Commissions: logs[i].Change.Copy().Neg(),
				}
			}

		case st.StakingLogTypeRewards:
			if removed {
				continue
			}
			if val, ok := updated[logs[i].ValidatorId]; ok {
				val.Rewards = val.Rewards.Sub(logs[i].Change)
			} else {
				updated[logs[i].ValidatorId] = &storage.Validator{
					Id:      logs[i].ValidatorId,
					Rewards: logs[i].Change.Copy().Neg(),
				}
			}
		}
	}

	if len(updated) > 0 {
		arr := make([]*storage.Validator, 0)
		for _, value := range updated {
			arr = append(arr, value)
		}
		if err = tx.UpdateValidators(ctx, arr...); err != nil {
			return
		}
	}

	if len(balances) > 0 {
		arr := make([]storage.Balance, 0)
		for _, value := range balances {
			arr = append(arr, *value)
		}
		if err = tx.SaveBalances(ctx, arr...); err != nil {
			return
		}
	}

	if len(delegations) > 0 {
		arr := make([]storage.Delegation, 0)
		for _, value := range delegations {
			arr = append(arr, *value)
		}
		if err = tx.SaveDelegations(ctx, arr...); err != nil {
			return
		}
	}

	if len(removedIds) > 0 {
		if err = tx.DeleteDelegationsByValidator(ctx, removedIds...); err != nil {
			return
		}
	}

	return
}

func delegationId(addrId, valId uint64) string {
	return fmt.Sprintf("%d_%d", addrId, valId)
}
