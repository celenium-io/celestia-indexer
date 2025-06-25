// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"github.com/celenium-io/celestia-indexer/internal/currency"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/legacy"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/celestia-app/v4/app"
	"github.com/celestiaorg/celestia-app/v4/app/encoding"
	blobTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	tmTypes "github.com/cometbft/cometbft/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type DecodedTx struct {
	TimeoutHeight uint64
	Memo          string
	Messages      []cosmosTypes.Msg
	Fee           decimal.Decimal
	Signers       map[types.Address][]byte
	Blobs         []*blobTypes.Blob
}

func NewDecodedTx() DecodedTx {
	return DecodedTx{
		Signers: make(map[types.Address][]byte),
	}
}

var (
	cfg, txDecoder = createDecoder()
)

func Tx(b types.BlockData, index int) (d DecodedTx, err error) {
	raw := b.Block.Txs[index]
	if bTx, isBlob := tmTypes.UnmarshalBlobTx(raw); isBlob {
		raw = bTx.Tx
		d.Blobs = bTx.Blobs
	}
	d.Signers = make(map[types.Address][]byte)

	if err = decodeCosmosTx(txDecoder, raw, &d); err != nil {
		return
	}

	return
}

func decodeCosmosTx(decoder cosmosTypes.TxDecoder, raw tmTypes.Tx, d *DecodedTx) error {
	txDecoded, err := decoder(raw)
	if err != nil {
		return errors.Wrap(err, "decoding tx error")
	}

	if t, ok := txDecoded.(cosmosTypes.TxWithTimeoutHeight); ok {
		d.TimeoutHeight = t.GetTimeoutHeight()
	}
	if t, ok := txDecoded.(cosmosTypes.TxWithMemo); ok {
		d.Memo = t.GetMemo()
	}
	if t, ok := txDecoded.(cosmosTypes.FeeTx); ok {
		d.Fee, err = decodeFee(t.GetFee())
		if err != nil {
			return errors.Wrap(err, "decode fee")
		}
	}
	if t, ok := txDecoded.(signing.Tx); ok {
		signers, err := t.GetSigners()
		if err != nil {
			pubKeys, err := t.GetPubKeys()
			if err != nil {
				return errors.Wrap(err, "get pub keys")
			}
			for _, pk := range pubKeys {
				address, err := types.NewAddressFromBytes(pk.Address().Bytes())
				if err != nil {
					return errors.Wrap(err, "NewAddressFromBytes")
				}
				d.Signers[address] = pk.Address().Bytes()
			}
		} else {
			for i := range signers {
				address, err := types.NewAddressFromBytes(signers[i])
				if err != nil {
					return errors.Wrap(err, "NewAddressFromBytes")
				}
				d.Signers[address] = signers[i]
			}
		}
	}

	d.Messages = txDecoded.GetMsgs()

	return nil
}

func decodeFee(amount cosmosTypes.Coins) (decimal.Decimal, error) {
	if amount == nil {
		return decimal.Zero, nil
	}

	if len(amount) > 1 {
		// TODO stop indexer if tx is not in failed status
		return decimal.Zero, errors.Errorf("found fee in %d currencies", len(amount))
	}

	fee, ok := getFeeInDenom(amount, currency.Utia)
	if !ok {
		if fee, ok = getFeeInDenom(amount, currency.Tia); !ok {
			// TODO stop indexer if tx is not in failed status
			return decimal.Zero, errors.New("couldn't find fee amount in utia or in tia denom")
		}
	}

	return fee, nil
}

func getFeeInDenom(amount cosmosTypes.Coins, denom string) (decimal.Decimal, bool) {
	ok, utiaCoin := amount.Find(denom)
	if !ok {
		return decimal.Zero, false
	}

	switch denom {
	case currency.Utia:
		fee := decimal.NewFromBigInt(utiaCoin.Amount.BigInt(), 0)
		return fee, true
	case currency.Tia:
		fee := decimal.NewFromBigInt(utiaCoin.Amount.BigInt(), 6)
		return fee, true
	default:
		return decimal.Zero, false
	}
}

func createDecoder() (encoding.Config, cosmosTypes.TxDecoder) {
	cfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	cfg.InterfaceRegistry.RegisterImplementations((*cosmosTypes.Msg)(nil), &legacy.MsgRegisterEVMAddress{})
	return cfg, cfg.TxConfig.TxDecoder()
}

func JsonTx(raw []byte) (cosmosTypes.Tx, error) {
	return cfg.TxConfig.TxJSONDecoder()(raw)
}
