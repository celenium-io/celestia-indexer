// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upValidatorRate, downValidatorRate)
}

func upValidatorRate(ctx context.Context, db *bun.DB) error {
	limit := 100
	offset := 0
	end := false

	validators := make(map[string]*storage.Validator)

	for !end {
		var msgs []storage.Message
		err := db.NewSelect().Model(&msgs).
			Where("type IN ('MsgCreateValidator', 'MsgEditValidator')").
			Limit(limit).
			Offset(offset).
			Order("id asc").
			Scan(ctx)
		if err != nil {
			return err
		}

		for i := range msgs {
			var tx storage.Tx
			err := db.NewSelect().Model(&tx).
				Where("id = ?", msgs[i].TxId).
				Scan(ctx)
			if err != nil {
				return err
			}

			if tx.Status != types.StatusSuccess {
				continue
			}

			switch msgs[i].Type {
			case types.MsgCreateValidator:
				validatorAddressVal, ok := msgs[i].Data["ValidatorAddress"]
				if !ok {
					continue
				}
				validatorAddress, ok := validatorAddressVal.(string)
				if !ok {
					continue
				}
				minSelfDelegation, ok := msgs[i].Data["MinSelfDelegation"]
				if !ok {
					minSelfDelegation = "0"
				}
				commission, ok := msgs[i].Data["Commission"]
				if !ok {
					continue
				}
				commissionMap, ok := commission.(map[string]any)
				if !ok {
					commissionMap = make(map[string]any)
				}
				rate, ok := commissionMap["Rate"]
				if !ok {
					rate = "0"
				}
				validators[validatorAddress] = &storage.Validator{
					Address:           validatorAddress,
					MinSelfDelegation: decimal.RequireFromString(minSelfDelegation.(string)),
					Rate:              decimal.RequireFromString(rate.(string)),
				}
			case types.MsgEditValidator:
				validatorAddressVal, ok := msgs[i].Data["ValidatorAddress"]
				if !ok {
					continue
				}
				validatorAddress, ok := validatorAddressVal.(string)
				if !ok {
					continue
				}
				validator, ok := validators[validatorAddress]
				if !ok {
					continue
				}
				minSelfDelegation, ok := msgs[i].Data["MinSelfDelegation"]
				if ok && minSelfDelegation != nil {
					validator.MinSelfDelegation = decimal.RequireFromString(minSelfDelegation.(string))
				}
				rate, ok := msgs[i].Data["CommissionRate"]
				if ok && rate != nil {
					validator.Rate = decimal.RequireFromString(rate.(string))
				}
			}
		}

		offset += len(msgs)
		end = len(msgs) < limit
	}

	for address, validator := range validators {
		if validator.Rate.IsZero() && validator.MinSelfDelegation.IsZero() {
			continue
		}

		query := db.NewUpdate().
			Model((*storage.Validator)(nil)).
			Where("address = ?", address)

		if !validator.Rate.IsZero() {
			query.Set("rate = ?", validator.Rate)
		}
		if !validator.MinSelfDelegation.IsZero() {
			query.Set("min_self_delegation = ?", validator.MinSelfDelegation)
		}

		if _, err := query.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}
func downValidatorRate(ctx context.Context, db *bun.DB) error {
	return nil
}
