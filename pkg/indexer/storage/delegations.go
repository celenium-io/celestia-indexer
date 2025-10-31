// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	decodeContext "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/pkg/errors"
)

func (module *Module) saveDelegations(
	ctx context.Context,
	tx storage.Transaction,
	dCtx *decodeContext.Context,
	addrToId map[string]uint64,
) error {
	if len(dCtx.StakingLogs) > 0 {
		for i := range dCtx.StakingLogs {
			if dCtx.StakingLogs[i].Address != nil {
				addressId, ok := addrToId[dCtx.StakingLogs[i].Address.Address]
				if !ok {
					return errors.Wrapf(errCantFindAddress, "staking log address %s", dCtx.StakingLogs[i].Address.Address)
				}
				dCtx.StakingLogs[i].AddressId = &addressId
			}

			validatorId, ok := module.validatorsByAddress[dCtx.StakingLogs[i].Validator.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "staking log validator address %s", dCtx.StakingLogs[i].Validator.Address)
			}
			dCtx.StakingLogs[i].ValidatorId = validatorId
		}

		if err := tx.SaveStakingLogs(ctx, dCtx.StakingLogs...); err != nil {
			return err
		}
	}

	if dCtx.Delegations.Len() > 0 {
		delegations := make([]storage.Delegation, 0)
		if err := dCtx.Delegations.Range(func(_ string, value *storage.Delegation) (error, bool) {
			addressId, ok := addrToId[value.Address.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "delegation address %s", value.Address.Address), false
			}
			value.AddressId = addressId

			validatorId, ok := module.validatorsByAddress[value.Validator.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "delegation validator address %s", value.Validator.Address), false
			}
			value.ValidatorId = validatorId

			delegations = append(delegations, *value)
			return nil, false
		}); err != nil {
			return err
		}

		if err := tx.SaveDelegations(ctx, delegations...); err != nil {
			return err
		}
	}

	if len(dCtx.Redelegations) > 0 {
		for i := range dCtx.Redelegations {
			addressId, ok := addrToId[dCtx.Redelegations[i].Address.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "delegation address %s", dCtx.Redelegations[i].Address.Address)
			}
			dCtx.Redelegations[i].AddressId = addressId

			srcId, ok := module.validatorsByAddress[dCtx.Redelegations[i].Source.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "source validator address %s", dCtx.Redelegations[i].Source.Address)
			}
			dCtx.Redelegations[i].SrcId = srcId

			destId, ok := module.validatorsByAddress[dCtx.Redelegations[i].Destination.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "dest validator address %s", dCtx.Redelegations[i].Destination.Address)
			}
			dCtx.Redelegations[i].DestId = destId
		}

		if err := tx.SaveRedelegations(ctx, dCtx.Redelegations...); err != nil {
			return err
		}
	}

	if len(dCtx.Undelegations) > 0 {
		for i := range dCtx.Undelegations {
			addressId, ok := addrToId[dCtx.Undelegations[i].Address.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "delegation address %s", dCtx.Undelegations[i].Address.Address)
			}
			dCtx.Undelegations[i].AddressId = addressId

			validatorId, ok := module.validatorsByAddress[dCtx.Undelegations[i].Validator.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "validator address %s", dCtx.Undelegations[i].Validator.Address)
			}
			dCtx.Undelegations[i].ValidatorId = validatorId
		}

		if err := tx.SaveUndelegations(ctx, dCtx.Undelegations...); err != nil {
			return err
		}
	}

	if len(dCtx.CancelUnbonding) > 0 {
		for i := range dCtx.CancelUnbonding {
			validatorId, ok := module.validatorsByAddress[dCtx.CancelUnbonding[i].Validator.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "cancel undelegation validator address %s", dCtx.CancelUnbonding[i].Validator.Address)
			}
			dCtx.CancelUnbonding[i].ValidatorId = validatorId

			addressId, ok := addrToId[dCtx.CancelUnbonding[i].Address.Address]
			if !ok {
				return errors.Wrapf(errCantFindAddress, "cancel undelegation address %s", dCtx.CancelUnbonding[i].Address.Address)
			}
			dCtx.CancelUnbonding[i].AddressId = addressId
		}

		if err := tx.CancelUnbondings(ctx, dCtx.CancelUnbonding...); err != nil {
			return err
		}
	}

	if err := tx.RetentionCompletedRedelegations(ctx, dCtx.Block.Time); err != nil {
		return errors.Wrap(err, "retention completed redelegations")
	}
	if err := tx.RetentionCompletedUnbondings(ctx, dCtx.Block.Time); err != nil {
		return errors.Wrap(err, "retention completed unbondings")
	}

	return nil
}
