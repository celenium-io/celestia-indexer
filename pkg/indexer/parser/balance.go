package parser

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/decode"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func parseCoinSpent(data map[string]any, height pkgTypes.Level) (*storage.Address, error) {
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
		Address:    coinSpent.Spender,
		Hash:       hash,
		Height:     height,
		LastHeight: height,
		Balance: storage.Balance{
			Currency: coinSpent.Amount.Denom,
			Total:    decimal.NewFromBigInt(coinSpent.Amount.Amount.Neg().BigInt(), 0),
		},
	}, nil
}

func parseCoinReceived(data map[string]any, height pkgTypes.Level) (*storage.Address, error) {
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
		Address:    coinReceived.Receiver,
		Hash:       hash,
		Height:     height,
		LastHeight: height,
		Balance: storage.Balance{
			Currency: coinReceived.Amount.Denom,
			Total:    decimal.NewFromBigInt(coinReceived.Amount.Amount.BigInt(), 0),
		},
	}, nil
}
