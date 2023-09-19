package rollback

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/decode"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (module *Module) rollbackBalances(ctx context.Context, tx storage.Transaction, deletedEvents []storage.Event, deletedAddresses []storage.Address) error {
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
		}

		if err != nil {
			return err
		}
		if address != nil {
			if _, ok := deleted[address.Address]; ok {
				continue
			}
			if addr, ok := updates[address.Address]; ok {
				addr.Balance.Total = addr.Balance.Total.Add(address.Balance.Total)
			} else {
				lastHeight, err := tx.LastAddressAction(ctx, address.Hash)
				if err != nil {
					return err
				}
				address.Height = pkgTypes.Level(lastHeight)
				updates[address.Address] = address
			}
		}
	}

	result := make([]*storage.Address, 0, len(updates))
	for _, addr := range updates {
		result = append(result, addr)
	}

	_, err := tx.SaveAddresses(ctx, result...)
	return err
}

func coinSpent(data map[string]any) (*storage.Address, error) {
	coinSpent, err := decode.NewCoinSpent(data)
	if err != nil {
		return nil, err
	}

	if coinSpent.Spender == "" {
		return nil, nil
	}
	_, hash, err := pkgTypes.Address(coinSpent.Spender).Decode()
	if err != nil {
		return nil, errors.Wrapf(err, "decode spender: %s", coinSpent.Spender)
	}
	return &storage.Address{
		Address: coinSpent.Spender,
		Hash:    hash,
		Balance: storage.Balance{
			Currency: coinSpent.Amount.Denom,
			Total:    decimal.NewFromBigInt(coinSpent.Amount.Amount.BigInt(), 0),
		},
	}, nil
}

func coinReceived(data map[string]any) (*storage.Address, error) {
	coinReceived, err := decode.NewCoinReceived(data)
	if err != nil {
		return nil, err
	}

	if coinReceived.Receiver == "" {
		return nil, nil
	}

	_, hash, err := pkgTypes.Address(coinReceived.Receiver).Decode()
	if err != nil {
		return nil, errors.Wrapf(err, "decode receiver: %s", coinReceived.Receiver)
	}
	return &storage.Address{
		Address: coinReceived.Receiver,
		Hash:    hash,
		Balance: storage.Balance{
			Currency: coinReceived.Amount.Denom,
			Total:    decimal.NewFromBigInt(coinReceived.Amount.Amount.Neg().BigInt(), 0),
		},
	}, nil
}
