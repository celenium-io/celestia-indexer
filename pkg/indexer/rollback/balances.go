// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (module *Module) rollbackBalances(
	ctx context.Context,
	tx storage.Transaction,
	deletedEvents []storage.Event,
	deletedAddresses []storage.Address,
) error {
	var (
		ids     = make([]uint64, len(deletedAddresses))
		deleted = make(map[string]struct{}, len(deletedAddresses))
	)
	for i := range deletedAddresses {
		ids[i] = deletedAddresses[i].Id
		deleted[deletedAddresses[i].Address] = struct{}{}
	}

	if err := tx.DeleteBalances(ctx, ids); err != nil {
		return err
	}

	if len(deletedEvents) == 0 {
		return nil
	}

	updates, err := getBalanceUpdates(ctx, tx, deleted, deletedEvents)
	if err != nil {
		return err
	}

	_, err = tx.SaveAddresses(ctx, updates...)
	return err
}

func getBalanceUpdates(
	ctx context.Context,
	tx storage.Transaction,
	deletedAddress map[string]struct{},
	deletedEvents []storage.Event,
) ([]*storage.Address, error) {
	updates := make(map[string]*storage.Address)

	for _, event := range deletedEvents {
		var (
			address *storage.Address
			err     error
		)

		switch event.Type {
		case types.EventTypeCoinSpent:
			address, err = coinSpent(event.Data)
		case types.EventTypeCoinReceived:
			address, err = coinReceived(event.Data)
		default:
			continue
		}

		if err != nil {
			return nil, err
		}

		if _, ok := deletedAddress[address.Address]; ok {
			continue
		}

		if addr, ok := updates[address.Address]; ok {
			addr.Balance.Spendable = addr.Balance.Spendable.Add(address.Balance.Spendable)
		} else {
			lastHeight, err := tx.LastAddressAction(ctx, address.Hash)
			if err != nil {
				return nil, err
			}

			//nolint:gosec
			address.LastHeight = pkgTypes.Level(lastHeight)
			updates[address.Address] = address
		}
	}

	result := make([]*storage.Address, 0, len(updates))
	for i := range updates {
		result = append(result, updates[i])
	}
	return result, nil
}

func coinSpent(data map[string]string) (*storage.Address, error) {
	coinSpent, err := decode.NewCoinSpent(data)
	if err != nil {
		return nil, err
	}

	_, hash, err := pkgTypes.Address(coinSpent.Spender).Decode()
	if err != nil {
		return nil, errors.Wrapf(err, "decode spender: %s", coinSpent.Spender)
	}
	balance := storage.Balance{
		Currency:  currency.DefaultCurrency,
		Spendable: types.NumericZero(),
	}
	for i := range coinSpent.Amount {
		if coinSpent.Amount[i] == nil || coinSpent.Amount[i].IsZero() {
			continue
		}
		if coinSpent.Amount[i].GetDenom() == currency.DefaultCurrency { // TODO: support other currencies
			amount := types.NumericFromBigInt(coinSpent.Amount[i].Amount.BigInt(), 0)
			balance.Spendable = amount
			balance.Currency = coinSpent.Amount[i].GetDenom()
		}
	}
	return &storage.Address{
		Address: coinSpent.Spender,
		Hash:    hash,
		Balance: balance,
	}, nil
}

func coinReceived(data map[string]string) (*storage.Address, error) {
	coinReceived, err := decode.NewCoinReceived(data)
	if err != nil {
		return nil, err
	}

	_, hash, err := pkgTypes.Address(coinReceived.Receiver).Decode()
	if err != nil {
		return nil, errors.Wrapf(err, "decode receiver: %s", coinReceived.Receiver)
	}

	balance := storage.Balance{
		Currency:  currency.DefaultCurrency,
		Spendable: types.NumericZero(),
	}
	for i := range coinReceived.Amount {
		if coinReceived.Amount[i] == nil || coinReceived.Amount[i].IsZero() {
			continue
		}
		if coinReceived.Amount[i].GetDenom() == currency.DefaultCurrency { // TODO: support other currencies
			amount := types.NumericFromBigInt(coinReceived.Amount[i].Amount.Neg().BigInt(), 0)
			balance.Spendable = amount
			balance.Currency = coinReceived.Amount[i].GetDenom()
		}
	}

	return &storage.Address{
		Address: coinReceived.Receiver,
		Hash:    hash,
		Balance: balance,
	}, nil
}
